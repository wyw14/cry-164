package analyzer

import (
	"github.com/wyw14/cry-164/internal/model"
	"sync"
)

type State struct {
	mu      sync.RWMutex
	reading model.Reading
}

func (s *State) Replace(reading model.Reading) {
	s.mu.Lock()
	s.reading = reading
	s.mu.Unlock()
}

func (s *State) Snapshot() model.Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.reading
}
