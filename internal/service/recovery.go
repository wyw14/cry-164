package service

import (
	"github.com/wyw14/cry-164/internal/model"
	"time"
)

func Restore(cycle model.Cycle) model.Cycle {
	if cycle.UpdatedAt.IsZero() {
		cycle.UpdatedAt = time.Now().UTC()
	}
	return cycle
}
