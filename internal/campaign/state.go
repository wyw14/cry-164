package campaign

import (
	"github.com/wyw14/cry-164/internal/model"
	"sync"
)

type State struct {
	mu    sync.RWMutex
	cycle model.Cycle
}

func NewState() *State { return &State{cycle: model.NewCycle()} }
func (s *State) Advance(next model.CycleState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cycle, err := s.cycle.Advance(next)
	if err != nil {
		return err
	}
	s.cycle = cycle
	return nil
}
func (s *State) Snapshot() model.Cycle { s.mu.RLock(); defer s.mu.RUnlock(); return s.cycle }

// Reset returns the production cycle to a fresh preparing cycle. Protection
// trips can fire from any live cycle state, and the guarded Advance() chain
// cannot move an already-running campaign back to preparing, so a tripped
// machine resets the campaign rather than continuing to run on stopped kit.
func (s *State) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cycle = model.NewCycle()
}
