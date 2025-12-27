package protocol

type Gossip struct {
	Members []Hello `json:"members"`
}
