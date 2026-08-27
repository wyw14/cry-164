package compressor

import (
	"context"
	"fmt"
	"time"
)

type Startup struct {
	duration time.Duration
	running  bool
}

func NewStartup(duration time.Duration) *Startup { return &Startup{duration: duration} }
func (s *Startup) Start(ctx context.Context) error {
	if s.running {
		return fmt.Errorf("compressor already running")
	}
	timer := time.NewTimer(s.duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		s.running = true
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (s *Startup) Running() bool { return s.running }
func (s *Startup) Stop()         { s.running = false }
