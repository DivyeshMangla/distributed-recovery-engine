package node

import "fmt"

type Node struct {
	ID   string
	Addr string
	Seed string
}

func NewNode(id, addr, seed string) *Node {
	return &Node{
		ID:   id,
		Addr: addr,
		Seed: seed,
	}
}

func (n *Node) Start() {
	fmt.Printf("node id: %s, listen address: %s, seed node address: %s\n", n.ID, n.Addr, n.Seed)
}
