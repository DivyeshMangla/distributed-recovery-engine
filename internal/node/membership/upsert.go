package membership

import (
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

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

	changed := false

	if local.Addr != in.Addr {
		local.Addr = in.Addr
		local.UpdatedAt = now
		changed = true
	}

	// Direct Alive observation always wins if newer
	if in.Status == Alive && in.LastSeen.After(local.LastSeen) {
		if local.Status != Alive {
			local.UpdatedAt = now
			changed = true
		}
		local.Status = Alive
		local.LastSeen = in.LastSeen
		if changed {
			m.dirty = true
		}
		return
	}

	// Accept gossip only if strictly newer
	if in.LastSeen.After(local.LastSeen) {
		if local.Status != in.Status {
			local.UpdatedAt = now
			changed = true
		}
		local.Status = in.Status
		local.LastSeen = in.LastSeen
		if changed {
			m.dirty = true
		}
		return
	}

	// Same timestamp: worse status wins (protocol consistency)
	if in.LastSeen.Equal(local.LastSeen) && in.Status > local.Status {
		local.Status = in.Status
		local.UpdatedAt = now
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
