package separator

import (
	"context"
	"sync"
	"time"
)

type Tank struct {
	mu       sync.Mutex
	cond     *sync.Cond
	pressure float64
	level    float64
}

func NewTank() *Tank { t := &Tank{}; t.cond = sync.NewCond(&t.mu); return t }
func (t *Tank) Set(pressure, level float64) {
	t.mu.Lock()
	t.pressure = pressure
	t.level = level
	t.cond.Broadcast()
	t.mu.Unlock()
}
func (t *Tank) WaitUntilReady(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			t.mu.Lock()
			t.cond.Broadcast()
			t.mu.Unlock()
		case <-done:
		}
	}()
	defer close(done)
	t.mu.Lock()
	defer t.mu.Unlock()
	for t.pressure < 10 || t.level < 0.25 {
		t.cond.Wait()
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}
func (t *Tank) Snapshot() (float64, float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.pressure, t.level
}
func Timeout() time.Duration { return 100 * time.Millisecond }
