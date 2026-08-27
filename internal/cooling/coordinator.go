package cooling

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-164/internal/separator"
	"time"
)

type AnalyzerBarrier struct{}

func (AnalyzerBarrier) Wait(ctx context.Context, sources []<-chan string) error {
	active := make([]<-chan string, 0, len(sources))
	for _, source := range sources {
		if source != nil {
			active = append(active, source)
		}
	}
	for _, source := range active {
		select {
		case <-source:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

type Controller struct{ tank *separator.Tank }

func NewController(tank *separator.Tank) *Controller { return &Controller{tank: tank} }
func (c *Controller) Establish(ctx context.Context) error {
	c.tank.Set(12, 0.3)
	select {
	case <-time.After(2 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (c *Controller) Status() string {
	pressure, level := c.tank.Snapshot()
	if pressure < 10 || level < 0.25 {
		return "cooling-awaiting"
	}
	return fmt.Sprintf("cooling-ready-%.2f", level)
}
