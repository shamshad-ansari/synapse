package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/shamshad-ansari/synapse/services/mock-canvas/internal/seed"
)

// ListDiscussionTopics returns discussion topics for a course.
// GET /api/v1/courses/{courseID}/discussion_topics
func ListDiscussionTopics(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(chi.URLParam(r, "courseID"))
	if err != nil {
		writeError(w, http.StatusNotFound, "The specified object cannot be found")
		return
	}
	items, ok := seed.DiscussionTopics[courseID]
	if !ok {
		writeJSON(w, http.StatusOK, []any{})
		return
	}
	writeJSON(w, http.StatusOK, items)
}
