package converter

import "sync"

type BedState struct {
	mu           sync.RWMutex
	temperatures []float64
}

func (s *BedState) Replace(values []float64) {
	s.mu.Lock()
	s.temperatures = append([]float64(nil), values...)
	s.mu.Unlock()
}

func (s *BedState) Snapshot() []float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]float64(nil), s.temperatures...)
}
