package analyzer

import (
	"context"
	"sync"
)

type Handshake func(context.Context) error

type Client struct {
	mu        sync.Mutex
	handshake Handshake
	err       error
	ready     bool
}

func NewClient(handshake Handshake) *Client { return &Client{handshake: handshake} }

// Read performs the handshake lazily. A successful handshake is remembered so
// later reads reuse the established session, but a failed handshake is never
// cached: the next read retries the handshake so a device that was offline at
// startup (e.g. the analyzer starting a minute late) can be recovered once it
// comes back online. Caching the failure under sync.Once is what previously
// pinned the first "connection refused" forever and suppressed all later
// connection attempts.
func (c *Client) Read(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ready {
		return nil
	}
	c.err = c.handshake(ctx)
	c.ready = c.err == nil
	return c.err
}

func (c *Client) Ready() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ready
}
