package analyzer

import (
	"github.com/wyw14/cry-164/internal/model"
	"sync"
)

type Stream struct {
	mu     sync.Mutex
	latest model.Reading
	seq    uint64
}

func (s *Stream) Publish(reading model.Reading) {
	s.mu.Lock()
	s.seq++
	reading.Sequence = s.seq
	s.latest = reading
	s.mu.Unlock()
}
func (s *Stream) Latest() model.Reading { s.mu.Lock(); defer s.mu.Unlock(); return s.latest }
