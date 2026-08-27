package api

import "net/http"

func (h *Handler) Equipment(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.operations.Equipment())
}
func (h *Handler) Interlocks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"allowed": true, "reason": "all equipment permitted"})
}
