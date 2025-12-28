package transport

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

const (
	udpBufferSize     = 65507
	udpIncomingBuffer = 1000
)

type UDPTransport struct {
	conn *net.UDPConn
	in   chan []byte
	done chan struct{}
	pool sync.Pool
}

func NewUDPTransport() *UDPTransport {
	return &UDPTransport{
		in:   make(chan []byte, udpIncomingBuffer),
		done: make(chan struct{}),
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, udpBufferSize)
			},
		},
	}
}

func (u *UDPTransport) Listen(addr protocol.Address) (<-chan []byte, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", string(addr))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	if err := conn.SetReadBuffer(4 * 1024 * 1024); err != nil {
		return nil, err
	}

	u.conn = conn
	go u.readLoop()

	return u.in, nil
}

func (u *UDPTransport) readLoop() {
	for {
		buf := u.pool.Get().([]byte)

		select {
		case <-u.done:
			u.pool.Put(buf)
			return
		default:
		}

		n, _, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			u.pool.Put(buf)
			return
		}

		msg, ok := u.parseFrame(buf[:n])
		u.pool.Put(buf)

		if !ok {
			continue
		}

		select {
		case u.in <- msg:
		case <-u.done:
			return
		}
	}
}

func (u *UDPTransport) parseFrame(data []byte) ([]byte, bool) {
	if len(data) < 4 {
		return nil, false
	}

	size := binary.BigEndian.Uint32(data[:4])
	if len(data) < int(size)+4 {
		return nil, false
	}

	msg := make([]byte, size)
	copy(msg, data[4:4+size])
	return msg, true
}

func (u *UDPTransport) Dial(addr protocol.Address, msg []byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp", string(addr))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SetWriteBuffer(4 * 1024 * 1024); err != nil {
		return err
	}

	return u.sendFrame(conn, msg)
}

func (u *UDPTransport) sendFrame(conn *net.UDPConn, msg []byte) error {
	frame := make([]byte, 4+len(msg))
	binary.BigEndian.PutUint32(frame[:4], uint32(len(msg)))
	copy(frame[4:], msg)

	_, err := conn.Write(frame)
	return err
}

func (u *UDPTransport) Close() error {
	close(u.done)
	close(u.in)

	if u.conn != nil {
		return u.conn.Close()
	}
	return nil
}
