package membership

import (
	"time"
)

type Member struct {
	ID       string
	Addr     string
	Status   Status
	LastSeen time.Time
}
