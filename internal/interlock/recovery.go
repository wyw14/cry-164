package interlock

func Recover(state *State, safe bool) bool {
	if !safe {
		return false
	}
	state.Clear()
	return true
}
