package node

import (
	"encoding/json"
	"fmt"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/node/membership"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/transport"
)

type Node struct {
	ID         protocol.NodeID
	Addr       protocol.Address
	Seed       protocol.Address
	Transport  transport.Transport
	Membership membership.Membership
}

func NewNode(id, addr, seed string) *Node {
	return &Node{
		ID:         protocol.NodeID(id),
		Addr:       protocol.Address(addr),
		Seed:       protocol.Address(seed),
		Membership: *membership.NewMembership(),
	}
}

func (n *Node) Start() error {
	t := transport.NewTCPTransport()
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

	n.Membership.Upsert(n.ID, n.Addr) // need to register self

	fmt.Printf(
		"node id: %s, listen address: %s, seed node address: %s\n",
		n.ID, n.Addr, n.Seed,
	)

	select {} // block, TODO: think of something better, surely blocking isn't the best
}

func (n *Node) handleIncoming(ch <-chan []byte) {
	for data := range ch {
		if n.handleHello(data) {
			continue
		}

		if n.handleGossip(data) {
			continue
		}

		if n.handleHeartbeat(data) {
			continue
		}

		if n.handleHeartbeatAck(data) {
			continue
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

	_ = n.Transport.Dial(n.Seed, payload)
}
