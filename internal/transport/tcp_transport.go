package transport

import (
	"net"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
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

func (t *TCPTransport) Listen(addr protocol.Address) (<-chan []byte, error) {
	ln, err := net.Listen("tcp", string(addr))
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

	t.in <- buf[:n]
}

func (t *TCPTransport) Dial(addr protocol.Address, msg []byte) error {
	conn, err := net.Dial("tcp", string(addr))
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
