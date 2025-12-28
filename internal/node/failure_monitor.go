package node

import "time"

func (n *Node) startFailureMonitor() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	last := time.Now()
	interval := 2 * time.Second

	for now := range ticker.C {
		if now.Sub(last) > 2*interval {
			last = now // lagged
			continue
		}

		n.Membership.MarkSuspect(10 * time.Second)
		n.Membership.MarkDead(30 * time.Second)
		last = now
	}
}
