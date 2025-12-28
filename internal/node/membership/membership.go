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

	local, ok := m.members[in.ID]
	if !ok {
		m.members[in.ID] = &in
		m.dirty = true
		return
	}

	// Address always updates
	local.Addr = in.Addr

	// Direct Alive always wins
	if in.Status == Alive && in.LastSeen.After(local.LastSeen) {
		local.Status = Alive
		local.LastSeen = in.LastSeen
		m.dirty = true
		return
	}

	// Gossip: accept only if strictly newer
	if in.LastSeen.After(local.LastSeen) {
		local.Status = in.Status
		local.LastSeen = in.LastSeen
		m.dirty = true
		return
	}

	// Same timestamp, worse status wins
	if in.LastSeen.Equal(local.LastSeen) && in.Status > local.Status {
		local.Status = in.Status
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
		m.dirty = true
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
