package cooling

import (
	"github.com/wyw14/cry-164/internal/model"
	"sync"
)

type Telemetry struct {
	mu        sync.RWMutex
	equipment model.Equipment
}

func (t *Telemetry) Set(e model.Equipment) { t.mu.Lock(); t.equipment = e; t.mu.Unlock() }
func (t *Telemetry) Get() model.Equipment  { t.mu.RLock(); defer t.mu.RUnlock(); return t.equipment }
