package membership

import (
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

type Member struct {
	ID       protocol.NodeID
	Addr     protocol.Address
	Status   Status
	LastSeen time.Time
}
