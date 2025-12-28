package membership

import (
	"fmt"
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

func (m *Membership) Upsert(id protocol.NodeID, addr protocol.Address, status Status) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if member, ok := m.members[id]; ok {
		member.Addr = addr

		// Direct communication always wins
		if status == Alive {
			member.Status = Alive
			member.LastSeen = time.Now()
			m.dirty = true
			return
		}

		// Gossip only worsens state
		if status > member.Status {
			fmt.Printf("[gossip] updated node=%s from %d to %d\n", id, member.Status, status)
			member.Status = status
			m.dirty = true
		}
		return
	}

	m.members[id] = &Member{
		ID:       id,
		Addr:     addr,
		Status:   status,
		LastSeen: time.Now(),
	}
	m.dirty = true
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
		if id == selfID || member.Status != Alive {
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
