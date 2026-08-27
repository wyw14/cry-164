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
// Trip latches the trip flag on the current machine frame. A surge trip is a
// latching protection, so it forces Trip true regardless of the prior frame
// value and preserves the last observed pressure/temperature for diagnostics.
func (s *Session) Trip() {
	s.mu.Lock()
	s.lastFrame.Trip = true
	s.updated = time.Now().UTC()
	s.mu.Unlock()
}
// ClearTrip releases the latched trip flag. It is the recover counterpart to
// Trip, so an operator reset clears the machine trip in lockstep with the
// interlock emergency being cleared rather than leaving the unit recorded as
// tripped after recovery.
func (s *Session) ClearTrip() {
	s.mu.Lock()
	s.lastFrame.Trip = false
	s.updated = time.Now().UTC()
	s.mu.Unlock()
}
func (s *Session) Snapshot() (string, Frame, time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id, s.lastFrame, s.updated
}
