package node

import (
	"encoding/json"
	"fmt"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/transport"
)

type Node struct {
	ID         string
	Addr       string
	Seed       string
	Transport  transport.Transport
	Membership Membership
}

func NewNode(id, addr, seed string) *Node {
	return &Node{
		ID:         id,
		Addr:       addr,
		Seed:       seed,
		Membership: *NewMembership(),
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

	if n.Seed != "" {
		n.sendHello()
	}

	fmt.Printf(
		"node id: %s, listen address: %s, seed node address: %s\n",
		n.ID, n.Addr, n.Seed,
	)

	select {} // block, TODO: think of something better, surely blocking isn't the best
}

func (n *Node) handleIncoming(ch <-chan []byte) {
	for data := range ch {
		var h protocol.Hello
		if err := json.Unmarshal(data, &h); err != nil {
			continue
		}

		n.Membership.Upsert(h.ID, h.Addr)

		fmt.Printf(
			"received hello from %s (%s), members=%d\n",
			h.ID,
			h.Addr,
			len(n.Membership.Snapshot()),
		)
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
