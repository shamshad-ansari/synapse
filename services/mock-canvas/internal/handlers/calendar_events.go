package handlers

import (
	"net/http"

	"github.com/shamshad-ansari/synapse/services/mock-canvas/internal/seed"
)

// ListCalendarEvents returns recurring lecture / recitation events.
// GET /api/v1/calendar_events
func ListCalendarEvents(w http.ResponseWriter, r *http.Request) {
	events := seed.GetCalendarEvents()
	writeJSON(w, http.StatusOK, events)
}
