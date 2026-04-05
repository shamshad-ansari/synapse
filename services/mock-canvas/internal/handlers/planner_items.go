package handlers

import (
	"net/http"

	"github.com/shamshad-ansari/synapse/services/mock-canvas/internal/seed"
)

// ListPlannerItems returns Canvas Planner items for the authenticated user.
// GET /api/v1/planner/items
func ListPlannerItems(w http.ResponseWriter, r *http.Request) {
	items := seed.GetPlannerItems()
	writeJSON(w, http.StatusOK, items)
}
