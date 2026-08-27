package cooling

import (
	"context"
	"github.com/wyw14/cry-164/internal/model"
)

func Integrate(ctx context.Context, controller *Controller, telemetry *Telemetry) error {
	if err := controller.Establish(ctx); err != nil {
		return err
	}
	telemetry.Set(model.Equipment{Name: "ammonia-condenser", Enabled: true, Pressure: 12, Temperature: -18})
	return nil
}
