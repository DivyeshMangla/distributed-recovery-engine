package node

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/node/membership"
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

	slog.Info("received hello", "from", h.ID, "addr", h.Addr, "isNew", isNew)

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

	msg := append([]byte{protocol.HelloPrefix}, payload...)
	_ = n.Transport.Dial(addr, msg)
}

func (n *Node) handleGossip(data []byte) bool {
	var g protocol.Gossip
	if err := json.Unmarshal(data, &g); err != nil || len(g.Members) == 0 {
		return false
	}

	for _, member := range g.Members {
		n.Membership.Upsert(membership.FromGossip(member))
	}

	slog.Debug("processed gossip", "memberCount", len(g.Members))
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
	msg := append([]byte{protocol.HeartbeatAckPrefix}, payload...)

	if addr != "" {
		_ = n.Transport.Dial(addr, msg)
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
