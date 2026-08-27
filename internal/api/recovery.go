package api

import (
	"github.com/wyw14/cry-164/internal/model"
	"net/http"
)

func RecoveryResponse(w http.ResponseWriter, cycle model.Cycle) {
	writeJSON(w, http.StatusOK, map[string]any{"cycle": cycle, "recovered": true})
}
