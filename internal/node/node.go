package node

import (
	"fmt"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/transport"
)

type Node struct {
	ID        string
	Addr      string
	Seed      string
	Transport transport.Transport
}

func NewNode(id, addr, seed string) *Node {
	return &Node{
		ID:   id,
		Addr: addr,
		Seed: seed,
	}
}

func (n *Node) Start() error {
	t := transport.NewTCPTransport()
	n.Transport = t

	if err := t.Listen(n.Addr); err != nil {
		return err
	}

	if n.Seed != "" {
		_ = t.Dial(n.Seed)
	}

	fmt.Printf("node id: %s, listen address: %s, seed node address: %s\n", n.ID, n.Addr, n.Seed)

	select {} // block
}
