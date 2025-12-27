package node

import "sync"

type Membership struct {
	mu      sync.Mutex
	members map[string]*Member // id:member
}
