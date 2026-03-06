package response

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func JSONError(w http.ResponseWriter, status int, code, message, requestID string) {
	JSON(w, status, Error{
		Code:      code,
		Message:   message,
		RequestID: requestID,
	})
}