package protocol

const HelloPrefix byte = 'O'

type Hello struct {
	ID   NodeID  `json:"id"`
	Addr Address `json:"addr"`
}
