package transport

import (
	"net"
)

type TCPTransport struct {
	ln net.Listener
	in chan []byte
}

func NewTCPTransport() *TCPTransport {
	return &TCPTransport{
		in: make(chan []byte),
	}
}

func (t *TCPTransport) Listen(addr string) (<-chan []byte, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	t.ln = ln

	go t.acceptLoop()

	return t.in, nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			return
		}

		go func(c net.Conn) {
			defer closeQuietly(c)

			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				return
			}

			t.in <- buf[:n]
		}(conn)

		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	defer closeQuietly(conn)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	_ = buf[:n]
}

func (t *TCPTransport) Dial(addr string, msg []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer closeQuietly(conn)

	_, err = conn.Write(msg)

	return err
}

func (t *TCPTransport) Close() error {
	if t.ln != nil {
		return t.ln.Close()
	}

	close(t.in)
	return nil
}

func closeQuietly(c net.Conn) {
	_ = c.Close()
}
