package feed

import "github.com/wyw14/cry-164/internal/model"

func PermitEquipment(permit *Permit) model.Equipment {
	return model.Equipment{Name: "synthesis-feed-header", Enabled: true, Pressure: permit.Pressure(), Temperature: 30}
}
