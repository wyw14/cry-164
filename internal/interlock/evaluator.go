package interlock

import "github.com/wyw14/cry-164/internal/model"

type Decision struct {
	Allowed bool
	Reason  string
}

func Evaluate(equipment []model.Equipment) Decision {
	for _, item := range equipment {
		if item.Trip {
			return Decision{Reason: item.Name + " tripped"}
		}
		if !item.Enabled {
			return Decision{Reason: item.Name + " disabled"}
		}
	}
	return Decision{Allowed: true, Reason: "all equipment permitted"}
}
