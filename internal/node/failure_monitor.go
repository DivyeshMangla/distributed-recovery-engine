package node

import "time"

func (n *Node) startFailureMonitor() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		n.Membership.MarkSuspect(10 * time.Second)
		n.Membership.MarkDead(30 * time.Second)
	}
}
