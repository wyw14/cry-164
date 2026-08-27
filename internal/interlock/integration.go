package interlock

import "github.com/wyw14/cry-164/internal/model"

func Protect(cycle model.Cycle, decision Decision) model.Cycle {
	if !decision.Allowed {
		cycle.State = model.Preparing
	}
	return cycle
}
