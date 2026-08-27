package api

import (
	"github.com/wyw14/cry-164/internal/service"
	"net/http"
)

type Handler struct {
	runtime    *service.Runtime
	operations *service.OperationService
	plant      *service.Plant
}

func NewHandler(runtime *service.Runtime, operations *service.OperationService, plant *service.Plant) *Handler {
	return &Handler{runtime: runtime, operations: operations, plant: plant}
}
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.runtime.Health())
}
func (h *Handler) Operations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.operations.Actions())
}
