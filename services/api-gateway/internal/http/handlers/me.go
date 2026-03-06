package handlers

import (
	"net/http"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/http/response"
)

// Phase 0: stubbed identity endpoint.
// In Step 3 we return a deterministic response with request_id.
// In Step 4+ we wire real auth and read user from DB.
type MeHandler struct{}

func (h MeHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.RequestIDFromContext(r.Context())

	response.JSON(w, http.StatusOK, map[string]any{
		"id":         "stub",
		"school_id":  "stub",
		"email":      "stub@example.com",
		"request_id": reqID,
	})
}