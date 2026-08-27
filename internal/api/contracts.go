package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OperationRequest struct {
	ActionID string `json:"action_id"`
	Detail   string `json:"detail"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func decodeOperationRequest(request *http.Request) (OperationRequest, error) {
	if request.Body == nil {
		return OperationRequest{}, errors.New("request body is required")
	}
	decoder := json.NewDecoder(io.LimitReader(request.Body, 64*1024))
	decoder.DisallowUnknownFields()
	var value OperationRequest
	if err := decoder.Decode(&value); err != nil {
		return OperationRequest{}, fmt.Errorf("decode operation request: %w", err)
	}
	value.ActionID = strings.TrimSpace(value.ActionID)
	value.Detail = strings.TrimSpace(value.Detail)
	if value.ActionID == "" {
		return OperationRequest{}, errors.New("action_id is required")
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return OperationRequest{}, errors.New("request body must contain one JSON object")
	}
	return value, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code string, err error) {
	message := "request failed"
	if err != nil && strings.TrimSpace(err.Error()) != "" {
		message = err.Error()
	}
	writeJSON(w, status, ErrorResponse{Code: code, Message: message})
}
