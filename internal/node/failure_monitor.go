package node

import "time"

const (
	failureCheckInterval = 2 * time.Second
	suspectTimeout       = 10 * time.Second
	deadTimeout          = 30 * time.Second
)

func (n *Node) startFailureMonitor() {
	ticker := time.NewTicker(failureCheckInterval)
	defer ticker.Stop()

	last := time.Now()

	for now := range ticker.C {
		if now.Sub(last) > 2*failureCheckInterval {
			last = now
			continue
		}

		n.Membership.MarkSuspect(suspectTimeout)
		n.Membership.MarkDead(deadTimeout)
		last = now
	}
}
