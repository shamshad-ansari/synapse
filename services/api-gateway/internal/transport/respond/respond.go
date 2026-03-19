package respond

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Data  any     `json:"data"`
	Error *string `json:"error"`
}

// JSON writes status and {"data": data, "error": null}.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Data: data, Error: nil})
}

// Error writes status and {"data": null, "error": message}.
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Data: nil, Error: &message})
}
