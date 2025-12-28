package node

import (
	"encoding/json"
	"fmt"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/node/membership"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

func (n *Node) handleHello(data []byte) bool {
	var h protocol.Hello
	if err := json.Unmarshal(data, &h); err != nil || h.ID == "" {
		return false
	}

	isNew := !n.Membership.Exists(h.ID)

	n.Membership.Upsert(membership.Member{
		ID:       h.ID,
		Addr:     h.Addr,
		Status:   membership.Alive,
		LastSeen: time.Now(), // direct observation
	})

	fmt.Printf(
		"received hello from %s (%s), members=%d\n",
		h.ID,
		h.Addr,
		len(n.Membership.Snapshot()),
	)

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
		n.Membership.Upsert(membership.FromGossip(member))
	}

	newSize := len(n.Membership.Snapshot())

	if newSize > oldSize {
		fmt.Printf("merged gossip, members=%d\n", newSize)
	}

	fmt.Println("received gossip:", n.Membership.Snapshot())
	return true
}

func (n *Node) handleHeartbeat(data []byte) bool {
	var hb protocol.Heartbeat
	if err := json.Unmarshal(data, &hb); err != nil || hb.From == "" {
		return false
	}

	n.Membership.MarkAlive(hb.From)
	n.replyHeartbeat(n.Membership.GetAddr(hb.From))

	return true
}

func (n *Node) replyHeartbeat(addr protocol.Address) {
	ack := protocol.HeartbeatAck{
		From:      n.ID,
		Timestamp: time.Now(),
	}

	payload, _ := json.Marshal(ack)
	if addr != "" {
		_ = n.Transport.Dial(addr, payload)
	}
}

func (n *Node) handleHeartbeatAck(data []byte) bool {
	var ack protocol.HeartbeatAck
	if err := json.Unmarshal(data, &ack); err != nil || ack.From == "" {
		return false
	}

	n.Membership.MarkAlive(ack.From)
	return true
}
