package converter

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-164/internal/steam"
	"time"
)

func RunHeating(ctx context.Context, valves []*steam.Valve, failAt int) error {
	var valve *steam.Valve
	for i := range valves {
		valve = valves[i]
		valve.Open()
		defer func() { valve.Close() }()
		if i == failAt {
			return fmt.Errorf("heating step %d failed", i)
		}
	}
	select {
	case <-timeAfter(ctx):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func timeAfter(ctx context.Context) <-chan time.Time { return time.After(1 * time.Millisecond) }
