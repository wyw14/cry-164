package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ActionKind string

const (
	ObserveAction ActionKind = "observe"
	ControlAction ActionKind = "control"
	ProtectAction ActionKind = "protect"
	RecoverAction ActionKind = "recover"
)

type Threshold struct {
	Name    string  `json:"name"`
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
	Unit    string  `json:"unit"`
}

type ActionDefinition struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Kind        ActionKind    `json:"kind"`
	Equipment   string        `json:"equipment"`
	Requires    []string      `json:"requires"`
	Thresholds  []Threshold   `json:"thresholds"`
	Timeout     time.Duration `json:"timeout"`
	Recoverable bool          `json:"recoverable"`
}

type ComponentPlan struct {
	Name    string             `json:"name"`
	Actions []ActionDefinition `json:"actions"`
}

type ActionStatus struct {
	ActionID   string    `json:"action_id"`
	Revision   uuid.UUID `json:"revision"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	State      string    `json:"state"`
	Detail     string    `json:"detail"`
}

type PlantGraph struct {
	mu       sync.RWMutex
	plans    map[string]ComponentPlan
	actions  map[string]ActionDefinition
	statuses map[string]ActionStatus
}

func NewPlantGraph(plans ...ComponentPlan) (*PlantGraph, error) {
	graph := &PlantGraph{
		plans:    make(map[string]ComponentPlan),
		actions:  make(map[string]ActionDefinition),
		statuses: make(map[string]ActionStatus),
	}
	for _, plan := range plans {
		if err := validateComponentPlan(plan); err != nil {
			return nil, err
		}
		if _, exists := graph.plans[plan.Name]; exists {
			return nil, fmt.Errorf("duplicate component plan %s", plan.Name)
		}
		graph.plans[plan.Name] = cloneComponentPlan(plan)
		for _, action := range plan.Actions {
			if _, exists := graph.actions[action.ID]; exists {
				return nil, fmt.Errorf("duplicate action %s", action.ID)
			}
			graph.actions[action.ID] = cloneAction(action)
			graph.statuses[action.ID] = ActionStatus{
				ActionID: action.ID,
				Revision: uuid.New(),
				State:    "waiting",
			}
		}
	}
	if err := graph.validateDependencies(); err != nil {
		return nil, err
	}
	return graph, nil
}

func validateComponentPlan(plan ComponentPlan) error {
	if strings.TrimSpace(plan.Name) == "" {
		return errors.New("component name is required")
	}
	if len(plan.Actions) == 0 {
		return fmt.Errorf("component %s has no actions", plan.Name)
	}
	seen := make(map[string]struct{}, len(plan.Actions))
	for _, action := range plan.Actions {
		if strings.TrimSpace(action.ID) == "" || strings.TrimSpace(action.Name) == "" {
			return fmt.Errorf("component %s has unnamed action", plan.Name)
		}
		if action.Kind != ObserveAction && action.Kind != ControlAction && action.Kind != ProtectAction && action.Kind != RecoverAction {
			return fmt.Errorf("action %s has invalid kind", action.ID)
		}
		if action.Timeout <= 0 {
			return fmt.Errorf("action %s has invalid timeout", action.ID)
		}
		if _, exists := seen[action.ID]; exists {
			return fmt.Errorf("component %s repeats action %s", plan.Name, action.ID)
		}
		seen[action.ID] = struct{}{}
		for _, threshold := range action.Thresholds {
			if threshold.Name == "" || threshold.Unit == "" || threshold.Minimum > threshold.Maximum {
				return fmt.Errorf("action %s has invalid threshold", action.ID)
			}
		}
	}
	return nil
}

func cloneAction(action ActionDefinition) ActionDefinition {
	action.Requires = append([]string(nil), action.Requires...)
	action.Thresholds = append([]Threshold(nil), action.Thresholds...)
	return action
}

func cloneComponentPlan(plan ComponentPlan) ComponentPlan {
	cloned := ComponentPlan{Name: plan.Name, Actions: make([]ActionDefinition, len(plan.Actions))}
	for index, action := range plan.Actions {
		cloned.Actions[index] = cloneAction(action)
	}
	return cloned
}

func (g *PlantGraph) validateDependencies() error {
	for _, action := range g.actions {
		for _, dependency := range action.Requires {
			if dependency == action.ID {
				return fmt.Errorf("action %s depends on itself", action.ID)
			}
			if _, exists := g.actions[dependency]; !exists {
				return fmt.Errorf("action %s depends on unknown action %s", action.ID, dependency)
			}
		}
	}
	visiting := make(map[string]bool)
	visited := make(map[string]bool)
	var visit func(string) error
	visit = func(id string) error {
		if visiting[id] {
			return fmt.Errorf("dependency cycle at %s", id)
		}
		if visited[id] {
			return nil
		}
		visiting[id] = true
		for _, dependency := range g.actions[id].Requires {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		visiting[id] = false
		visited[id] = true
		return nil
	}
	for id := range g.actions {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}

func (g *PlantGraph) Components() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	values := make([]string, 0, len(g.plans))
	for name := range g.plans {
		values = append(values, name)
	}
	sort.Strings(values)
	return values
}

func (g *PlantGraph) Operations() []ActionDefinition {
	g.mu.RLock()
	defer g.mu.RUnlock()
	values := make([]ActionDefinition, 0, len(g.actions))
	for _, action := range g.actions {
		values = append(values, cloneAction(action))
	}
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values
}

func (g *PlantGraph) Equipment() []Equipment {
	g.mu.RLock()
	defer g.mu.RUnlock()
	byName := make(map[string]Equipment)
	for _, action := range g.actions {
		if _, exists := byName[action.Equipment]; exists {
			continue
		}
		byName[action.Equipment] = Equipment{
			Name:        action.Equipment,
			Enabled:     true,
			Pressure:    1,
			Temperature: 20,
		}
	}
	values := make([]Equipment, 0, len(byName))
	for _, equipment := range byName {
		values = append(values, equipment)
	}
	sort.Slice(values, func(i, j int) bool { return values[i].Name < values[j].Name })
	return values
}

func (g *PlantGraph) Begin(actionID string, at time.Time) (ActionStatus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	action, exists := g.actions[actionID]
	if !exists {
		return ActionStatus{}, fmt.Errorf("unknown action %s", actionID)
	}
	for _, dependency := range action.Requires {
		if g.statuses[dependency].State != "complete" {
			return ActionStatus{}, fmt.Errorf("action %s waits for %s", actionID, dependency)
		}
	}
	status := g.statuses[actionID]
	status.Revision = uuid.New()
	status.StartedAt = at.UTC()
	status.FinishedAt = time.Time{}
	status.State = "running"
	status.Detail = ""
	g.statuses[actionID] = status
	return status, nil
}

func (g *PlantGraph) Finish(actionID string, at time.Time, detail string) (ActionStatus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	status, exists := g.statuses[actionID]
	if !exists {
		return ActionStatus{}, fmt.Errorf("unknown action %s", actionID)
	}
	if status.State != "running" {
		return ActionStatus{}, fmt.Errorf("action %s is not running", actionID)
	}
	status.FinishedAt = at.UTC()
	status.State = "complete"
	status.Detail = strings.TrimSpace(detail)
	g.statuses[actionID] = status
	return status, nil
}

func (g *PlantGraph) Fail(actionID string, at time.Time, detail string) (ActionStatus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	status, exists := g.statuses[actionID]
	if !exists {
		return ActionStatus{}, fmt.Errorf("unknown action %s", actionID)
	}
	status.FinishedAt = at.UTC()
	status.State = "failed"
	status.Detail = strings.TrimSpace(detail)
	status.Revision = uuid.New()
	g.statuses[actionID] = status
	return status, nil
}

func (g *PlantGraph) ResetRecoverable() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	count := 0
	for id, status := range g.statuses {
		if status.State == "failed" && g.actions[id].Recoverable {
			status.State = "waiting"
			status.Detail = ""
			status.Revision = uuid.New()
			g.statuses[id] = status
			count++
		}
	}
	return count
}

func (g *PlantGraph) Statuses() []ActionStatus {
	g.mu.RLock()
	defer g.mu.RUnlock()
	values := make([]ActionStatus, 0, len(g.statuses))
	for _, status := range g.statuses {
		values = append(values, status)
	}
	sort.Slice(values, func(i, j int) bool { return values[i].ActionID < values[j].ActionID })
	return values
}

func (g *PlantGraph) ReadyActions() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	values := make([]string, 0)
	for id, action := range g.actions {
		if g.statuses[id].State != "waiting" {
			continue
		}
		ready := true
		for _, dependency := range action.Requires {
			if g.statuses[dependency].State != "complete" {
				ready = false
				break
			}
		}
		if ready {
			values = append(values, id)
		}
	}
	sort.Strings(values)
	return values
}

func (g *PlantGraph) Summary() map[string]any {
	g.mu.RLock()
	defer g.mu.RUnlock()
	counts := map[string]int{"waiting": 0, "running": 0, "complete": 0, "failed": 0}
	for _, status := range g.statuses {
		counts[status.State]++
	}
	return map[string]any{
		"components": len(g.plans),
		"actions":    len(g.actions),
		"states":     counts,
	}
}
