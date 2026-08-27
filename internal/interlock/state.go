package interlock

import "sync"

type State struct {
	mu        sync.RWMutex
	emergency bool
	reason    string
}

func (s *State) SetEmergency(reason string) {
	s.mu.Lock()
	s.emergency = true
	s.reason = reason
	s.mu.Unlock()
}
func (s *State) Clear() { s.mu.Lock(); s.emergency = false; s.reason = ""; s.mu.Unlock() }
func (s *State) Snapshot() (bool, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.emergency, s.reason
}
