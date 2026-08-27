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
