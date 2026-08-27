package compressor

import (
	"sync"
	"time"
)

type Session struct {
	mu        sync.Mutex
	id        string
	lastFrame Frame
	updated   time.Time
}

func NewSession(id string) *Session { return &Session{id: id} }
func (s *Session) Update(frame Frame) {
	s.mu.Lock()
	s.lastFrame = frame
	s.updated = time.Now().UTC()
	s.mu.Unlock()
}
func (s *Session) Snapshot() (string, Frame, time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id, s.lastFrame, s.updated
}
