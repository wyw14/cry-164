package campaign

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-164/internal/compressor"
	"github.com/wyw14/cry-164/internal/feed"
	"github.com/wyw14/cry-164/internal/model"
	"time"
)

type StartupPlan struct {
	Lubrication  time.Duration
	StageTimeout time.Duration
}
type Coordinator struct {
	permit  *feed.Permit
	startup *compressor.Startup
	config  model.Config
}

func NewCoordinator(permit *feed.Permit, startup *compressor.Startup, config model.Config) *Coordinator {
	return &Coordinator{permit: permit, startup: startup, config: config}
}
func (c *Coordinator) Start(ctx context.Context, plan StartupPlan) error {
	if plan.StageTimeout <= 0 {
		plan.StageTimeout = c.config.StageTimeout
	}
	lubeCtx, cancel := context.WithTimeout(ctx, plan.StageTimeout)
	defer cancel()
	if err := c.permit.Confirm(lubeCtx, plan.Lubrication); err != nil {
		return fmt.Errorf("feed confirmation: %w", err)
	}
	stageCtx, stageCancel := context.WithTimeout(ctx, plan.StageTimeout)
	defer stageCancel()
	if err := c.startup.Start(stageCtx); err != nil {
		return fmt.Errorf("compressor startup: %w", err)
	}
	return nil
}
func (c *Coordinator) NewCycle() model.Cycle { return model.NewCycle() }
