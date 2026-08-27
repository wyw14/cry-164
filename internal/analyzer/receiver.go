package analyzer

import (
	"context"
	"fmt"
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
func (c *Client) Read(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ready {
		return nil
	}
	err := c.handshake(ctx)
	c.err = err
	c.ready = err == nil
	if err != nil {
		return err
	}
	if !c.ready {
		return fmt.Errorf("analyzer unavailable")
	}
	return nil
}
func (c *Client) Ready() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.ready }
