package membership

import (
	"math/rand"
	"sync"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

type Membership struct {
	mu      sync.Mutex
	members map[protocol.NodeID]*Member // id -> member
}

func NewMembership() *Membership {
	return &Membership{
		members: make(map[protocol.NodeID]*Member),
	}
}

func (m *Membership) Upsert(id protocol.NodeID, addr protocol.Address) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if member, ok := m.members[id]; ok {
		member.Addr = addr
		member.Status = Alive
		member.LastSeen = time.Now()
		return
	}

	m.members[id] = &Member{
		ID:       id,
		Addr:     addr,
		Status:   Alive,
		LastSeen: time.Now(),
	}
}

func (m *Membership) Snapshot() []Member {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]Member, 0, len(m.members))
	for _, member := range m.members {
		out = append(out, *member)
	}
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

		if id == selfID {
			continue
		}

		if member.Status != Alive {
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
