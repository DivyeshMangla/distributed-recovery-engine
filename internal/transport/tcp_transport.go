package transport

import (
	"errors"
	"net"
)

type TCPTransport struct {
	ln net.Listener
}

func NewTCPTransport() *TCPTransport {
	return &TCPTransport{}
}

func (t *TCPTransport) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	t.ln = ln

	go t.acceptLoop()

	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Temporary() {
				continue
			}
			return
		}

		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	defer closeQuietly(conn)
	// for future
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer closeQuietly(conn)

	return nil
}

func (t *TCPTransport) Close() error {
	if t.ln != nil {
		return t.ln.Close()
	}

	return nil
}

func closeQuietly(c net.Conn) {
	_ = c.Close()
}
