package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) Router() chi.Router {
	router := chi.NewRouter()
	router.Get("/healthz", h.Health)
	router.Get("/api/operations", h.Operations)
	router.Get("/api/equipment", h.Equipment)
	router.Get("/api/interlocks", h.Interlocks)
	router.Get("/api/incidents", h.Incidents)
	router.Get("/api/status", h.PlantStatus)
	router.Post("/api/actions/{action}", h.ExecutePlantAction)
	router.Get("/api/operations/components", h.OperationComponents)
	router.Get("/api/operations/statuses", h.OperationStatuses)
	router.Get("/api/operations/ready", h.OperationReady)
	router.Get("/api/operations/summary", h.OperationSummary)
	router.Post("/api/operations/begin", h.OperationBegin)
	router.Post("/api/operations/complete", h.OperationComplete)
	router.Post("/api/operations/fail", h.OperationFail)
	router.Post("/api/operations/recover", h.OperationRecover)
	return router
}
func (h *Handler) Incidents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []any{})
}
