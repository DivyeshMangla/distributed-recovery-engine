package transport

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

const incomingBufferSize = 10000

type TCPTransport struct {
	ln net.Listener
	in chan []byte
}

func NewTCPTransport() *TCPTransport {
	return &TCPTransport{
		in: make(chan []byte, incomingBufferSize),
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
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	defer conn.Close()

	msg, err := readFrame(conn)
	if err != nil {
		return
	}

	t.in <- msg
}

func (t *TCPTransport) Dial(addr protocol.Address, msg []byte) error {
	conn, err := net.Dial("tcp", string(addr))
	if err != nil {
		return err
	}
	defer conn.Close()

	return writeFrame(conn, msg)
}

func (t *TCPTransport) Close() error {
	if t.ln != nil {
		return t.ln.Close()
	}
	close(t.in)
	return nil
}

func readFrame(r io.Reader) ([]byte, error) {
	var size uint32
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return nil, err
	}

	buf := make([]byte, size)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

func writeFrame(w io.Writer, msg []byte) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(msg))); err != nil {
		return err
	}
	_, err := w.Write(msg)
	return err
}
