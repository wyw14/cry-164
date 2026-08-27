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
func (d *Dispatcher) Apply(plan []ValveCommand) Result {
	applied := make([]ValveState, 0, len(plan))
	for _, command := range plan {
		if err := d.actuator(command); err != nil {
			return Result{Err: fmt.Errorf("valve %s: %w", command.Name, err)}
		}
		state := ValveState{Name: command.Name, Position: command.Position}
		applied = append(applied, state)
		d.state = append(d.state, state)
	}
	return Result{Applied: applied, Snapshot: append([]ValveState(nil), d.state...)}
}
