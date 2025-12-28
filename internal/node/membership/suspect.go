package membership

import (
	"fmt"
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
			fmt.Printf(
				"[suspect] node=%s addr=%s last_seen=%s silence=%s\n",
				member.ID,
				member.Addr,
				member.LastSeen.Format(time.RFC3339),
				now.Sub(member.LastSeen).Truncate(time.Millisecond),
			)
		}
	}
}
