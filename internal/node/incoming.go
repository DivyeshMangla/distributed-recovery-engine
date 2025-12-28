package node

import (
	"encoding/json"
	"fmt"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

func (n *Node) handleHello(data []byte) bool {
	var h protocol.Hello
	if err := json.Unmarshal(data, &h); err != nil || h.ID == protocol.NodeID("") {
		return false
	}

	isNew := !n.Membership.Exists(h.ID)
	n.Membership.Upsert(h.ID, h.Addr)

	fmt.Printf("received hello from %s (%s), members=%d\n", h.ID, h.Addr, len(n.Membership.Snapshot()))

	if isNew {
		n.replyHello(h.Addr)
	}

	return true
}

func (n *Node) replyHello(addr protocol.Address) {
	payload, err := json.Marshal(protocol.Hello{
		ID:   n.ID,
		Addr: n.Addr,
	})

	if err != nil {
		return
	}

	_ = n.Transport.Dial(addr, payload)
}

func (n *Node) handleGossip(data []byte) bool {
	var g protocol.Gossip
	if err := json.Unmarshal(data, &g); err != nil || len(g.Members) == 0 {
		return false
	}

	oldSize := len(n.Membership.Snapshot())

	for _, member := range g.Members {
		n.Membership.Upsert(member.ID, member.Addr)
	}

	newSize := len(n.Membership.Snapshot())

	if newSize > oldSize {
		fmt.Printf("merged gossip, members=%d\n", newSize)
	}

	return true
}
