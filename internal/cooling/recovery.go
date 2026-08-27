package cooling

import "github.com/wyw14/cry-164/internal/model"

func RecoverCooling(previous model.Equipment) model.Equipment {
	previous.Enabled = true
	previous.Trip = false
	return previous
}
