package protocol

type Gossip struct {
	Members []GossipMember `json:"members"`
}

type GossipMember struct {
	ID     NodeID  `json:"id"`
	Addr   Address `json:"addr"`
	Status int     `json:"status"`
}
