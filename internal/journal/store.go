package journal

import (
	"encoding/json"
	"github.com/wyw14/cry-164/internal/model"
	"os"
	"sync"
)

type Store struct {
	mu      sync.Mutex
	path    string
	entries []model.Incident
}

func NewStore(path string) *Store { return &Store{path: path} }
func (s *Store) Append(incident model.Incident) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, incident)
	data, err := json.Marshal(s.entries)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
func (s *Store) Entries() []model.Incident {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]model.Incident(nil), s.entries...)
}
