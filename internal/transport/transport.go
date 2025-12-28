package transport

import "github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"

type Transport interface {
	Listen(addr protocol.Address) (<-chan []byte, error)
	Dial(addr protocol.Address, msg []byte) error
	Close() error
}
