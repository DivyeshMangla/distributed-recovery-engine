package protocol

import "time"

type Heartbeat struct {
	From      NodeID    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
}

type HeartbeatAck struct {
	From      NodeID    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
}
