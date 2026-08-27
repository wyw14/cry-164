package separator

import "github.com/wyw14/cry-164/internal/model"

func Separate(reading model.Reading) model.Equipment {
	return model.Equipment{Name: "liquid-ammonia-separator", Enabled: reading.Ammonia > 0, Pressure: reading.Nitrogen + reading.Hydrogen, Temperature: 80}
}
func ProductReady(e model.Equipment) bool { return e.Healthy() && e.Pressure < 180 }
