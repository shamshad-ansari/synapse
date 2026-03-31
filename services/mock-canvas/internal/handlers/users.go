package handlers

import (
	"net/http"

	"github.com/shamshad-ansari/synapse/services/mock-canvas/internal/seed"
)

// GetSelf returns the current user profile.
// GET /api/v1/users/self
func GetSelf(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, seed.CurrentUser)
}
