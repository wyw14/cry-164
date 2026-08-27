package service

import (
	"fmt"
	"github.com/wyw14/cry-164/internal/model"
	"time"
)

type OperationService struct{ graph *model.PlantGraph }

func NewOperationService(graph *model.PlantGraph) *OperationService {
	return &OperationService{graph: graph}
}

func (s *OperationService) Components() []string { return s.graph.Components() }

func (s *OperationService) Actions() []model.ActionDefinition { return s.graph.Operations() }

func (s *OperationService) Equipment() []model.Equipment { return s.graph.Equipment() }

func (s *OperationService) Begin(actionID string) (model.ActionStatus, error) {
	if actionID == "" {
		return model.ActionStatus{}, fmt.Errorf("action id is required")
	}
	return s.graph.Begin(actionID, time.Now().UTC())
}

func (s *OperationService) Complete(actionID, detail string) (model.ActionStatus, error) {
	return s.graph.Finish(actionID, time.Now().UTC(), detail)
}

func (s *OperationService) Fail(actionID, detail string) (model.ActionStatus, error) {
	return s.graph.Fail(actionID, time.Now().UTC(), detail)
}

func (s *OperationService) Recover() int { return s.graph.ResetRecoverable() }

func (s *OperationService) Statuses() []model.ActionStatus { return s.graph.Statuses() }

func (s *OperationService) Ready() []string { return s.graph.ReadyActions() }

func (s *OperationService) Summary() map[string]any { return s.graph.Summary() }
