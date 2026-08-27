package feed

import (
	"context"
	"fmt"
	"time"
)

type Permit struct {
	pressure float64
	enabled  bool
}

func NewPermit(pressure float64) *Permit { return &Permit{pressure: pressure, enabled: true} }

func (p *Permit) Confirm(ctx context.Context, delay time.Duration) error {
	if !p.enabled {
		return fmt.Errorf("feed disabled")
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Permit) Pressure() float64 { return p.pressure }
func (p *Permit) Disable()          { p.enabled = false }
