package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/shamshad-ansari/synapse/services/mock-canvas/internal/seed"
)

// ListCourses returns all enrolled courses.
// GET /api/v1/courses
func ListCourses(w http.ResponseWriter, r *http.Request) {
	// Strip syllabus_body from list responses (Canvas only includes it
	// when explicitly requested with include[]=syllabus_body on a single course).
	type courseNoSyllabus struct {
		seed.Course
		SyllabusBody string `json:"syllabus_body,omitempty"`
	}

	out := make([]courseNoSyllabus, len(seed.Courses))
	for i, c := range seed.Courses {
		c.SyllabusBody = ""
		out[i] = courseNoSyllabus{Course: c}
	}
	writeJSON(w, http.StatusOK, out)
}

// GetCourse returns a single course by ID.
// GET /api/v1/courses/{courseID}
func GetCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(chi.URLParam(r, "courseID"))
	if err != nil {
		writeError(w, http.StatusNotFound, "The specified object cannot be found")
		return
	}

	course, ok := seed.CoursesByID[courseID]
	if !ok {
		writeError(w, http.StatusNotFound, "The specified object cannot be found")
		return
	}

	includes := r.URL.Query().Get("include[]")
	if !strings.Contains(includes, "syllabus_body") {
		course.SyllabusBody = ""
	}

	writeJSON(w, http.StatusOK, course)
}
