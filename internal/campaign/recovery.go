package campaign

import "github.com/wyw14/cry-164/internal/model"

func Recover(cycle model.Cycle) model.Cycle {
	if cycle.State == model.Stable {
		cycle.Revision++
	}
	return cycle
}
