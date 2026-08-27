package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CycleState string

const (
	Preparing   CycleState = "preparing"
	Compressing CycleState = "compressing"
	Reacting    CycleState = "reacting"
	Condensing  CycleState = "condensing"
	Separating  CycleState = "separating"
	Stable      CycleState = "stable"
)

type Operation struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Cycle struct {
	ID        uuid.UUID  `json:"id"`
	State     CycleState `json:"state"`
	Revision  uint64     `json:"revision"`
	StartedAt time.Time  `json:"started_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Equipment struct {
	Name        string  `json:"name"`
	Enabled     bool    `json:"enabled"`
	Pressure    float64 `json:"pressure"`
	Temperature float64 `json:"temperature"`
	Trip        bool    `json:"trip"`
}

type Reading struct {
	Source   string  `json:"source"`
	Nitrogen float64 `json:"nitrogen"`
	Hydrogen float64 `json:"hydrogen"`
	Ammonia  float64 `json:"ammonia"`
	Sequence uint64  `json:"sequence"`
}

type Incident struct {
	ID       uuid.UUID `json:"id"`
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	At       time.Time `json:"at"`
}

func NewOperation(name string) Operation {
	return Operation{ID: uuid.New(), Name: strings.TrimSpace(name), CreatedAt: time.Now().UTC()}
}

func NewCycle() Cycle {
	now := time.Now().UTC()
	return Cycle{ID: uuid.New(), State: Preparing, Revision: 1, StartedAt: now, UpdatedAt: now}
}

func (c Cycle) Advance(next CycleState) (Cycle, error) {
	allowed := map[CycleState]CycleState{Preparing: Compressing, Compressing: Reacting, Reacting: Condensing, Condensing: Separating, Separating: Stable}
	if allowed[c.State] != next {
		return c, errors.New("invalid cycle transition")
	}
	c.State = next
	c.Revision++
	c.UpdatedAt = time.Now().UTC()
	return c, nil
}

func (e Equipment) Healthy() bool {
	return e.Enabled && !e.Trip && e.Pressure >= 0 && e.Temperature < 520
}
