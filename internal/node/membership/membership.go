package membership

import (
	"log/slog"
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
		if member.Status == Dead {
			slog.Info("node resurrected", "node", id, "addr", member.Addr)
		}
		member.Status = Alive
		member.LastSeen = time.Now()
		// Do NOT mark dirty — heartbeats are high-frequency noise
	}
}
