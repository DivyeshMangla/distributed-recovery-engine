package node

import (
	"encoding/json"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

func (n *Node) startGossip() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		n.gossipToPeer()
	}
}

func (n *Node) gossipToPeer() {
	peer := n.Membership.PickRandomPeer(n.ID)
	if peer == nil {
		return
	}

	payload, err := n.buildGossipPayload()
	if err != nil {
		return
	}

	_ = n.Transport.Dial(peer.Addr, payload)
}

func (n *Node) buildGossipPayload() ([]byte, error) {
	snapshot := n.Membership.Snapshot()

	members := make([]protocol.GossipMember, 0, len(snapshot))
	for _, m := range snapshot {
		members = append(members, protocol.GossipMember{
			ID:     m.ID,
			Addr:   m.Addr,
			Status: int(m.Status),
		})
	}

	return json.Marshal(protocol.Gossip{
		Members: members,
	})
}
