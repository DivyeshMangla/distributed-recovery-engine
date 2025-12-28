package membership

import (
	"log/slog"
	"time"
)

func (m *Membership) MarkSuspect(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	for _, member := range m.members {
		if member.Status != Alive {
			continue
		}

		if now.Sub(member.LastSeen) > timeout {
			member.Status = Suspect
			slog.Info("node suspected",
				"node", member.ID,
				"addr", member.Addr,
				"lastSeen", member.LastSeen.Format(time.RFC3339),
				"silence", now.Sub(member.LastSeen).Truncate(time.Millisecond),
			)
		}
	}
}

func (m *Membership) MarkDead(deadTimeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	for _, member := range m.members {
		if member.Status != Suspect {
			continue
		}

		if now.Sub(member.LastSeen) > deadTimeout {
			member.Status = Dead
			slog.Info("node declared dead",
				"node", member.ID,
				"addr", member.Addr,
				"lastSeen", member.LastSeen.Format(time.RFC3339),
				"silence", now.Sub(member.LastSeen).Truncate(time.Millisecond),
			)
		}
	}
}
