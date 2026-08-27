package campaign

import "github.com/wyw14/cry-164/internal/model"

func IntegrateProtection(state *State, equipment []model.Equipment) error {
	for _, item := range equipment {
		if item.Trip || !item.Enabled {
			return state.Advance(model.Compressing)
		}
	}
	return nil
}
