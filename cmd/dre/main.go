package main

import (
	"flag"
	"github.com/divyeshmangla/distributed-recovery-engine/internal/node"
)

func main() {

	//--id=node-1
	//--addr=127.0.0.1:8001
	//--seed=127.0.0.1:8000

	id := flag.String("id", "", "node id")
	addr := flag.String("addr", "", "listen address")
	seed := flag.String("seed", "", "seed node address")

	flag.Parse()
	node.NewNode(*id, *addr, *seed).Start()
}
