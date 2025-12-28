package protocol

type Hello struct {
	ID   NodeID  `json:"id"`
	Addr Address `json:"addr"`
}
