package node

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

const (
	gossipInterval     = 1 * time.Second
	gossipMaxMembers   = 8
	gossipRecentWindow = 30 * time.Second
)

func (n *Node) startGossip() {
	ticker := time.NewTicker(gossipInterval)
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

	msg := append([]byte{protocol.GossipPrefix}, payload...)
	_ = n.Transport.Dial(peer.Addr, msg)
	slog.Debug("sent gossip", "to", peer.ID, "payloadSize", len(msg))
}

func (n *Node) buildGossipPayload() ([]byte, error) {
	targets := n.Membership.SelectGossipTargets(n.ID, gossipMaxMembers, gossipRecentWindow)

	members := make([]protocol.GossipMember, 0, len(targets))
	for _, m := range targets {
		members = append(members, protocol.GossipMember{
			ID:       m.ID,
			Addr:     m.Addr,
			Status:   int(m.Status),
			LastSeen: m.LastSeen,
		})
	}

	return json.Marshal(protocol.Gossip{
		Members: members,
	})
}
