package service

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-164/internal/analyzer"
	"github.com/wyw14/cry-164/internal/campaign"
	"github.com/wyw14/cry-164/internal/cooling"
	"github.com/wyw14/cry-164/internal/model"
	"sync"
)

type Runtime struct {
	mu       sync.RWMutex
	state    *campaign.State
	cooling  *cooling.Controller
	analyzer *analyzer.Stream
	metrics  model.Metrics
}

func NewRuntime(state *campaign.State, controller *cooling.Controller, stream *analyzer.Stream) *Runtime {
	return &Runtime{state: state, cooling: controller, analyzer: stream}
}
func (r *Runtime) Snapshot() model.Cycle {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state.Snapshot()
}
func (r *Runtime) Start(ctx context.Context) error {
	if err := r.cooling.Establish(ctx); err != nil {
		return err
	}
	r.metrics.CountOperation()
	return nil
}
func (r *Runtime) Health() map[string]any {
	return map[string]any{"status": "ok", "cycle": r.Snapshot(), "metrics": r.metrics.Snapshot()}
}
func (r *Runtime) Analysis() string {
	reading := r.analyzer.Latest()
	return fmt.Sprintf("source=%s sequence=%d", reading.Source, reading.Sequence)
}
