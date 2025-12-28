package node

import (
	"encoding/json"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
	"time"
)

func (n *Node) startHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)
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
}
