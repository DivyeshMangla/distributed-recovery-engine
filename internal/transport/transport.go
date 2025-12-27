package transport

type Transport interface {
	Listen(addr string) error
	Dial(addr string) error
	Close() error
}
