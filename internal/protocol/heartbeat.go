package protocol

import "time"

const (
	HeartbeatPrefix    byte = 'H'
	HeartbeatAckPrefix byte = 'A'
)

type Heartbeat struct {
	From      NodeID    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
}

type HeartbeatAck struct {
	From      NodeID    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
}
