package converter

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-164/internal/analyzer"
	"github.com/wyw14/cry-164/internal/steam"
)

type Controller struct {
	analyzer *analyzer.Stream
	steam    *steam.Recovery
}

func NewController(stream *analyzer.Stream, recovery *steam.Recovery) *Controller {
	return &Controller{analyzer: stream, steam: recovery}
}
func (c *Controller) React(ctx context.Context) error {
	reading := c.analyzer.Latest()
	if reading.Source == "" {
		return fmt.Errorf("no analysis reading")
	}
	return c.steam.Capture(ctx, steam.Record(reading))
}
func (c *Controller) Summary() string { return fmt.Sprintf("recovered=%s", c.steam.Total()) }
