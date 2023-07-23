package pool

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type TCPPool struct {
	index uint
	size  uint

	mu      sync.RWMutex
	clients []*TCPClient

	network string
	laddr   *net.TCPAddr
	raddr   *net.TCPAddr

	maxRetries    int
	retryInterval time.Duration
}

// Dial returns a new *TCPPool.
//
// The new client connects to the remote address `raddr` on the network `network`,
// which must be "tcp", "tcp4", or "tcp6".
func Dial(network, addr string) (*TCPPool, error) {
	raddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	return DialTCP(network, nil, raddr)
}

// DialTCP returns a new *TCPPool.
//
// The new client connects to the remote address `raddr` on the network `network`,
// which must be "tcp", "tcp4", or "tcp6".
// If `laddr` is not nil, it is used as the local address for the connection.
func DialTCP(network string, laddr, raddr *net.TCPAddr) (*TCPPool, error) {
	return &TCPPool{
		index: 0,
		size:  1,

		mu: sync.RWMutex{},
		clients: []*TCPClient{
			newTCPClient(network, laddr, raddr),
		},

		network: network,
		laddr:   laddr,
		raddr:   raddr,

		maxRetries:    1,
		retryInterval: 1 * time.Millisecond,
	}, nil
}

func (p *TCPPool) SetPoolSize(size uint) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	previousSize := p.size

	if size > previousSize {
		for i := previousSize; i < size; i++ {
			client := newTCPClient(p.network, p.laddr, p.raddr)
			client.SetMaxRetries(p.maxRetries)
			client.SetRetryInterval(p.retryInterval)

			p.clients = append(p.clients, client)
		}
	} else if size < previousSize {
		for i := size; i < previousSize; i++ {
			// error is ignored
			_ = p.clients[i].Close()
		}

		p.clients = p.clients[:size]
	}

	p.size = size

	return nil
}

func (p *TCPPool) SetMaxRetries(maxRetries int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.maxRetries = maxRetries
	for i := range p.clients {
		p.clients[i].SetMaxRetries(maxRetries)
	}
}

func (p *TCPPool) SetRetryInterval(retryInterval time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.retryInterval = retryInterval
	for i := range p.clients {
		p.clients[i].SetRetryInterval(retryInterval)
	}
}

// implements net.Conn
func (p *TCPPool) Read(b []byte) (int, error) {
	return -1, fmt.Errorf("github.com/samber/go-tcp-pool is a write only connection")
}

// implements net.Conn
func (p *TCPPool) Write(b []byte) (int, error) {
	p.mu.RLock()

	index := p.index % p.size
	client := p.clients[index]
	p.index++

	p.mu.RUnlock()

	return client.Write(b)
}

// implements net.Conn
func (p *TCPPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := range p.clients {
		_ = p.clients[i].Close()
	}

	p.size = 0
	p.clients = []*TCPClient{}

	return nil
}

// implements net.Conn
func (p *TCPPool) LocalAddr() net.Addr {
	return p.laddr
}

// implements net.Conn
func (p *TCPPool) RemoteAddr() net.Addr {
	return p.raddr
}

// implements net.Conn
func (p *TCPPool) SetDeadline(t time.Time) error {
	// @TODO: implement
	return nil
}

// implements net.Conn
func (p *TCPPool) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("github.com/samber/go-tcp-pool is a write only connection")
}

// implements net.Conn
func (p *TCPPool) SetWriteDeadline(t time.Time) error {
	// @TODO: implement
	return nil
}
