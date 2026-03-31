package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type canvasError struct {
	Errors  []canvasErrorItem `json:"errors"`
	Status  string            `json:"status"`
}

type canvasErrorItem struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, canvasError{
		Errors: []canvasErrorItem{{Message: msg}},
		Status: "unauthenticated",
	})
}

// RequireToken is middleware that validates a Bearer token is present.
// It accepts any non-empty token (this is a mock).
func RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "Invalid access token.")
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "Invalid access token.")
			return
		}
		next.ServeHTTP(w, r)
	})
}
