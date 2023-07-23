package pool

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type TCPClient struct {
	*net.TCPConn

	mu sync.RWMutex

	network string
	laddr   *net.TCPAddr
	raddr   *net.TCPAddr

	maxRetries    int
	retryInterval time.Duration
}

func newTCPClient(network string, laddr, raddr *net.TCPAddr) *TCPClient {
	client := &TCPClient{
		TCPConn: nil,

		mu: sync.RWMutex{},

		network: network,
		laddr:   laddr,
		raddr:   raddr,

		maxRetries:    10,
		retryInterval: 10 * time.Millisecond,
	}

	// error is ignored
	_ = client.reconnect()

	return client
}

func (c *TCPClient) SetMaxRetries(maxRetries int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.maxRetries = maxRetries
}

func (c *TCPClient) SetRetryInterval(retryInterval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.retryInterval = retryInterval
}

func (c *TCPClient) reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := net.DialTCP(c.network, c.laddr, c.raddr)
	if err != nil {
		return err
	}

	_ = conn.CloseRead()

	if c.TCPConn != nil {
		_ = c.TCPConn.Close()
	}
	c.TCPConn = conn
	return nil
}

func (c *TCPClient) Write(b []byte) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	disconnected := c.TCPConn == nil

	t := c.retryInterval
	for i := 0; i < c.maxRetries; i++ {
		// reconnect
		if disconnected {
			if i > 0 {
				time.Sleep(t)
			}

			c.mu.RUnlock()
			err := c.reconnect()
			c.mu.RLock()

			if err != nil {
				disconnected = true
				continue
			}
		}

		// won't be executed if the connection failed
		n, err := c.TCPConn.Write(b)
		if err == nil {
			return n, err
		}

		disconnected = true
	}

	return -1, fmt.Errorf("ErrMaxRetries")
}
