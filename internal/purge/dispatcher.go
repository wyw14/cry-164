package purge

import "fmt"

type ValveCommand struct {
	Name     string
	Position float64
}
type ValveState struct {
	Name     string
	Position float64
}
type Actuator func(ValveCommand) error
type Result struct {
	Applied  []ValveState
	Snapshot []ValveState
	Err      error
}
type Dispatcher struct {
	actuator Actuator
	state    []ValveState
}

func NewDispatcher(actuator Actuator) *Dispatcher { return &Dispatcher{actuator: actuator} }

// Apply executes a multi-valve plan. The actuator moves each valve physically
// as soon as it is called, so a plan cannot be applied atomically: when a
// later valve rejects its new position, the valves before it have already
// moved for real and cannot be rolled back.
//
// The dispatcher therefore records each valve's new position the moment its
// actuator accepts it. On a partial failure the committed state reflects the
// physical truth: valves that moved hold their new positions, the rejected
// valve and any valves after it keep their prior positions. The Snapshot
// surfaces that partial state so the page shows reality (not a stale "old"
// setting) and the next round computes from the actual valve positions. The
// valves that did move are also reported via Applied so the operator can
// reconcile the valve that rejected.
//
// Positions are stored per valve (last write wins) rather than appended, so
// the snapshot always holds one current entry per valve instead of
// accumulating duplicates across rounds.
func (d *Dispatcher) Apply(plan []ValveCommand) Result {
	applied := make([]ValveState, 0, len(plan))
	for _, command := range plan {
		if err := d.actuator(command); err != nil {
			return Result{
				Applied:  applied,
				Snapshot: append([]ValveState(nil), d.state...),
				Err:      partialFailureError(command.Name, err, applied),
			}
		}
		applied = append(applied, ValveState{Name: command.Name, Position: command.Position})
		d.state = upsertPosition(d.state, command)
	}
	return Result{Applied: applied, Snapshot: append([]ValveState(nil), d.state...)}
}

func upsertPosition(states []ValveState, command ValveCommand) []ValveState {
	for i := range states {
		if states[i].Name == command.Name {
			states[i].Position = command.Position
			return states
		}
	}
	return append(states, ValveState{Name: command.Name, Position: command.Position})
}

// partialFailureError wraps the rejected valve's error and names the valves
// that had already accepted their new position. Those valves moved physically
// and their new positions are reflected in the snapshot, so the operator must
// be told the plan did not apply as a whole.
func partialFailureError(valve string, cause error, applied []ValveState) error {
	moved := make([]string, 0, len(applied))
	for _, state := range applied {
		moved = append(moved, fmt.Sprintf("%s=%g", state.Name, state.Position))
	}
	if len(moved) == 0 {
		return fmt.Errorf("valve %s: %w", valve, cause)
	}
	return fmt.Errorf("valve %s: %w (already moved: %v; snapshot reflects partial state)", valve, cause, moved)
}
