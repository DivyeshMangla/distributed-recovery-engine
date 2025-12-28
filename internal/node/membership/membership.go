package membership

import (
	"math/rand"
	"sync"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

type Membership struct {
	mu      sync.Mutex
	members map[protocol.NodeID]*Member

	snapshot []Member
	dirty    bool
}

func NewMembership() *Membership {
	return &Membership{
		members: make(map[protocol.NodeID]*Member),
		dirty:   true,
	}
}

func (m *Membership) Upsert(in Member) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	local, ok := m.members[in.ID]
	if !ok {
		in.UpdatedAt = now
		m.members[in.ID] = &in
		m.dirty = true
		return
	}

	// Address always updates
	if local.Addr != in.Addr {
		local.Addr = in.Addr
		local.UpdatedAt = now
		m.dirty = true
	}

	// Direct Alive always wins
	if in.Status == Alive && in.LastSeen.After(local.LastSeen) {
		if local.Status != Alive {
			local.UpdatedAt = now
		}
		local.Status = Alive
		local.LastSeen = in.LastSeen
		m.dirty = true
		return
	}

	// Gossip: accept only if strictly newer
	if in.LastSeen.After(local.LastSeen) {
		if local.Status != in.Status {
			local.UpdatedAt = now
		}
		local.Status = in.Status
		local.LastSeen = in.LastSeen
		m.dirty = true
		return
	}

	// Same timestamp, worse status wins
	if in.LastSeen.Equal(local.LastSeen) && in.Status > local.Status {
		local.Status = in.Status
		local.UpdatedAt = now
		m.dirty = true
	}
}

func (m *Membership) Snapshot() []Member {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.dirty && m.snapshot != nil {
		return m.snapshot
	}

	out := make([]Member, 0, len(m.members))
	for _, member := range m.members {
		out = append(out, *member)
	}

	m.snapshot = out
	m.dirty = false
	return out
}

func (m *Membership) PickRandomPeer(selfID protocol.NodeID) *Member {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.members) <= 1 {
		return nil
	}

	peers := make([]*Member, 0, len(m.members))
	for id, member := range m.members {
		if id == selfID || member.Status == Dead {
			continue
		}
		peers = append(peers, member)
	}

	if len(peers) == 0 {
		return nil
	}

	return peers[rand.Intn(len(peers))]
}

func (m *Membership) Exists(id protocol.NodeID) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.members[id]
	return exists
}

func (m *Membership) GetAddr(id protocol.NodeID) protocol.Address {
	m.mu.Lock()
	defer m.mu.Unlock()

	if member, ok := m.members[id]; ok {
		return member.Addr
	}
	return ""
}

func (m *Membership) MarkAlive(id protocol.NodeID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if member, ok := m.members[id]; ok {
		member.Status = Alive
		member.LastSeen = time.Now()
		// Do NOT mark dirty — heartbeats are high-frequency noise
	}
}

func FromGossip(g protocol.GossipMember) Member {
	return Member{
		ID:       g.ID,
		Addr:     g.Addr,
		Status:   Status(g.Status),
		LastSeen: g.LastSeen,
	}
}

// SelectGossipTargets returns a hybrid selection for delta gossip:
// - Always includes self
// - Members updated within recentWindow
// - Random sample of alive peers to fill remaining slots
// Returns at most maxCount members.
func (m *Membership) SelectGossipTargets(selfID protocol.NodeID, maxCount int, recentWindow time.Duration) []Member {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	result := make([]Member, 0, maxCount)
	included := make(map[protocol.NodeID]bool)

	// 1. Always include self
	if self, ok := m.members[selfID]; ok {
		result = append(result, *self)
		included[selfID] = true
	}

	// 2. Add recently changed members
	for id, member := range m.members {
		if len(result) >= maxCount {
			break
		}
		if included[id] || member.Status == Dead {
			continue
		}
		if now.Sub(member.UpdatedAt) <= recentWindow {
			result = append(result, *member)
			included[id] = true
		}
	}

	// 3. Fill remaining slots with random alive peers
	if len(result) < maxCount {
		var candidates []*Member
		for id, member := range m.members {
			if !included[id] && member.Status != Dead {
				candidates = append(candidates, member)
			}
		}
		// Shuffle and pick
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		for _, member := range candidates {
			if len(result) >= maxCount {
				break
			}
			result = append(result, *member)
		}
	}

	return result
}
