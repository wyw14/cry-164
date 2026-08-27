package compressor

import "github.com/wyw14/cry-164/internal/model"

func EquipmentState(session *Session) model.Equipment {
	_, frame, _ := session.Snapshot()
	return model.Equipment{Name: "recycle-compressor", Enabled: true, Pressure: frame.Pressure, Temperature: frame.Temperature, Trip: frame.Trip}
}
