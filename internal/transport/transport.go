package transport

type Transport interface {
	Listen(addr string) (<-chan []byte, error)
	Dial(addr string, msg []byte) error
	Close() error
}
