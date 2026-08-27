package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) PlantStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"plant": h.plant.Status(), "analysis": h.runtime.Analysis()})
}

func (h *Handler) ExecutePlantAction(w http.ResponseWriter, r *http.Request) {
	request, err := decodeOperationRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err)
		return
	}
	result, err := h.plant.Execute(r.Context(), chi.URLParam(r, "action"), request.Detail)
	if err != nil {
		writeError(w, http.StatusConflict, "action_failed", err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) OperationComponents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.operations.Components())
}

func (h *Handler) OperationStatuses(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.operations.Statuses())
}

func (h *Handler) OperationReady(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.operations.Ready())
}

func (h *Handler) OperationSummary(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.operations.Summary())
}

func (h *Handler) OperationBegin(w http.ResponseWriter, r *http.Request) {
	request, err := decodeOperationRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err)
		return
	}
	status, err := h.operations.Begin(request.ActionID)
	if err != nil {
		writeError(w, http.StatusConflict, "begin_failed", err)
		return
	}
	writeJSON(w, http.StatusAccepted, status)
}

func (h *Handler) OperationComplete(w http.ResponseWriter, r *http.Request) {
	request, err := decodeOperationRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err)
		return
	}
	status, err := h.operations.Complete(request.ActionID, request.Detail)
	if err != nil {
		writeError(w, http.StatusConflict, "complete_failed", err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (h *Handler) OperationFail(w http.ResponseWriter, r *http.Request) {
	request, err := decodeOperationRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err)
		return
	}
	status, err := h.operations.Fail(request.ActionID, request.Detail)
	if err != nil {
		writeError(w, http.StatusConflict, "fail_failed", err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (h *Handler) OperationRecover(w http.ResponseWriter, r *http.Request) {
	h.operations.Recover()
	RecoveryResponse(w, h.runtime.Snapshot())
}
