package protocol

import "time"

type Gossip struct {
	Members []GossipMember `json:"members"`
}

type GossipMember struct {
	ID       NodeID    `json:"id"`
	Addr     Address   `json:"addr"`
	Status   int       `json:"status"`
	LastSeen time.Time `json:"lastSeen"`
}
