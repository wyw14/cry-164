package steam

import (
	"context"
	"github.com/wyw14/cry-164/internal/model"
	"time"
)

type Recovery struct{ recovered time.Duration }

func NewRecovery() *Recovery { return &Recovery{} }
func (r *Recovery) Capture(ctx context.Context, energy time.Duration) error {
	select {
	case <-time.After(energy / 10):
		r.recovered += energy
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (r *Recovery) Total() time.Duration { return r.recovered }
func Record(reading model.Reading) time.Duration {
	return time.Duration(reading.Ammonia*10) * time.Second
}
