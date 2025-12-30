package node

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/node/membership"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/transport"
)

type Node struct {
	ID     protocol.NodeID
	LifeID protocol.LifeID
	Addr   protocol.Address
	Seed   protocol.Address

	Transport  transport.Transport
	Membership membership.Membership
}

func NewNode(id, addr, seed string) *Node {
	return &Node{
		ID:         protocol.NodeID(id),
		LifeID:     0, // initial life starts at 0th
		Addr:       protocol.Address(addr),
		Seed:       protocol.Address(seed),
		Membership: *membership.NewMembership(),
	}
}

func (n *Node) Start() error {
	t := transport.NewUDPTransport()
	n.Transport = t

	ch, err := t.Listen(n.Addr)
	if err != nil {
		return err
	}

	go n.handleIncoming(ch)
	go n.startGossip()
	go n.startHeartbeat()
	go n.startFailureMonitor()

	if n.Seed != ("") {
		n.sendHello()
	}

	n.Membership.Upsert(membership.Member{
		ID:       n.ID,
		Addr:     n.Addr,
		Status:   membership.Alive,
		LastSeen: time.Now(),
	})

	slog.Info("node started",
		"id", n.ID,
		"addr", n.Addr,
		"seed", n.Seed,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down", "id", n.ID)
	_ = n.Transport.Close()
	return nil
}

func (n *Node) handleIncoming(ch <-chan []byte) {
	for data := range ch {
		if len(data) < 2 {
			continue
		}

		switch data[0] {
		case protocol.HelloPrefix:
			n.handleHello(data[1:])
		case protocol.GossipPrefix:
			n.handleGossip(data[1:])
		case protocol.HeartbeatPrefix:
			n.handleHeartbeat(data[1:])
		case protocol.HeartbeatAckPrefix:
			n.handleHeartbeatAck(data[1:])
		}
	}
}

func (n *Node) sendHello() {
	payload, err := json.Marshal(protocol.Hello{
		ID:   n.ID,
		Addr: n.Addr,
	})
	if err != nil {
		return
	}

	msg := append([]byte{protocol.HelloPrefix}, payload...)
	_ = n.Transport.Dial(n.Seed, msg)
}
