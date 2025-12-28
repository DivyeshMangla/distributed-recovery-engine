package node

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

const heartbeatInterval = 1 * time.Second

func (n *Node) startHeartbeat() {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		n.sendHeartbeat()
	}
}

func (n *Node) sendHeartbeat() {
	peer := n.Membership.PickRandomPeer(n.ID)
	if peer == nil {
		return
	}

	hb := protocol.Heartbeat{
		From:      n.ID,
		Timestamp: time.Now(),
	}

	payload, err := json.Marshal(hb)
	if err != nil {
		return
	}

	_ = n.Transport.Dial(peer.Addr, payload)

	n.Membership.MarkAlive(n.ID) // don't sus yourself out
	slog.Debug("sent heartbeat", "to", peer.ID)
}
