package cooling

import "github.com/wyw14/cry-164/internal/model"

func CondenserState(telemetry *Telemetry) model.Equipment {
	return telemetry.Get()
}
