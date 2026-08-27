package analyzer

import (
	"context"
	"fmt"
	"sync"
)

type Handshake func(context.Context) error

type Client struct {
	once      sync.Once
	handshake Handshake
	err       error
	ready     bool
}

func NewClient(handshake Handshake) *Client { return &Client{handshake: handshake} }

func (c *Client) Read(ctx context.Context) error {
	c.once.Do(func() {
		c.err = c.handshake(ctx)
		c.ready = c.err == nil
	})
	if c.err != nil {
		return c.err
	}
	if !c.ready {
		return fmt.Errorf("analyzer unavailable")
	}
	return nil
}

func (c *Client) Ready() bool { return c.ready }
