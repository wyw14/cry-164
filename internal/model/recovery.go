package model

import "time"

func RecoverEquipment(equipment Equipment, at time.Time) Equipment {
	equipment.Enabled = true
	equipment.Trip = false
	if equipment.Temperature < 0 {
		equipment.Temperature = 0
	}
	return equipment
}
