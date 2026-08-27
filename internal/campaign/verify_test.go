package campaign

import (
	"context"
	"testing"
	"time"

	"github.com/wyw14/cry-164/internal/compressor"
	"github.com/wyw14/cry-164/internal/feed"
	"github.com/wyw14/cry-164/internal/model"
)

func TestStartupStagesReceiveIndependentTimeoutBudgets(t *testing.T) {
	coordinator := NewCoordinator(feed.NewPermit(120), compressor.NewStartup(35*time.Millisecond), model.DefaultConfig())
	err := coordinator.Start(context.Background(), StartupPlan{Lubrication: 45 * time.Millisecond, StageTimeout: 65 * time.Millisecond})
	if err != nil {
		t.Fatalf("independent startup stages should complete: %v", err)
	}
}
