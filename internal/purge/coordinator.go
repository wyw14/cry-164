package purge

import "sync"

type PlanStore struct {
	mu     sync.RWMutex
	latest []ValveState
}

func (s *PlanStore) Commit(states []ValveState) {
	s.mu.Lock()
	s.latest = append([]ValveState(nil), states...)
	s.mu.Unlock()
}
func (s *PlanStore) Snapshot() []ValveState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]ValveState(nil), s.latest...)
}
