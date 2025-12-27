package node

import (
	"sync"
	"time"
)

type Membership struct {
	mu      sync.Mutex
	members map[string]*Member // id -> member
}

func NewMembership() *Membership {
	return &Membership{
		members: make(map[string]*Member),
	}
}

func (m *Membership) Upsert(id, addr string) {
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
