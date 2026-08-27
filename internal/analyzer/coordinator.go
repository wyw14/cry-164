package analyzer

import (
	"context"
	"sync/atomic"
	"time"
)

type Probe func(context.Context) error
type Poller struct {
	probe  Probe
	active atomic.Int64
	max    atomic.Int64
}

func NewPoller(probe Probe) *Poller { return &Poller{probe: probe} }
func (p *Poller) Poll(ctx context.Context, attempts int) error {
	var last error
	for i := 0; i < attempts; i++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
		current := p.active.Add(1)
		for {
			old := p.max.Load()
			if current <= old || p.max.CompareAndSwap(old, current) {
				break
			}
		}
		err := p.probe(attemptCtx)
		defer cancel()
		defer p.active.Add(-1)
		if err == nil {
			return nil
		}
		last = err
	}
	return last
}
func (p *Poller) MaxOutstanding() int64 { return p.max.Load() }
