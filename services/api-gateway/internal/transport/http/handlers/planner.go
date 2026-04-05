package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

// PlannerHandler serves planner endpoints for study sessions and deadlines.
type PlannerHandler struct {
	Service service.PlannerService
}

func (h *PlannerHandler) authIDs(r *http.Request) (userID, schoolID uuid.UUID, ok bool) {
	uid, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	sid, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return uid, sid, true
}

// ListSessions GET /v1/planner/sessions?start=YYYY-MM-DD&end=YYYY-MM-DD
func (h *PlannerHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if start == "" || end == "" {
		respond.Error(w, http.StatusBadRequest, "start and end date params required (YYYY-MM-DD)")
		return
	}
	sessions, err := h.Service.ListStudySessions(r.Context(), userID, schoolID, start, end)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if sessions == nil {
		respond.JSON(w, http.StatusOK, []*domain.StudySession{})
		return
	}
	respond.JSON(w, http.StatusOK, sessions)
}

// CreateSession POST /v1/planner/sessions
func (h *PlannerHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body struct {
		Title           string `json:"title"`
		ScheduledDate   string `json:"scheduled_date"`
		StartTime       string `json:"start_time"`
		DurationMinutes int    `json:"duration_minutes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.Service.CreateStudySession(r.Context(), userID, schoolID, body.Title, body.ScheduledDate, body.StartTime, body.DurationMinutes)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

// UpdateSessionStatus PATCH /v1/planner/sessions/{id}/status
func (h *PlannerHandler) UpdateSessionStatus(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid session id")
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.Service.UpdateStudySessionStatus(r.Context(), id, userID, schoolID, body.Status); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteSession DELETE /v1/planner/sessions/{id}
func (h *PlannerHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid session id")
		return
	}
	if err := h.Service.DeleteStudySession(r.Context(), id, userID, schoolID); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// MarkMissed POST /v1/planner/missed-yesterday
func (h *PlannerHandler) MarkMissed(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	count, err := h.Service.MarkMissedYesterday(r.Context(), userID, schoolID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]int{"marked": count})
}

// ListDeadlines GET /v1/planner/deadlines?limit=10
func (h *PlannerHandler) ListDeadlines(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	limit := 10
	if ls := r.URL.Query().Get("limit"); ls != "" {
		n, err := strconv.Atoi(ls)
		if err == nil && n > 0 {
			limit = n
		}
	}
	deadlines, err := h.Service.ListUpcomingDeadlines(r.Context(), userID, schoolID, limit)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if deadlines == nil {
		respond.JSON(w, http.StatusOK, []*domain.UpcomingDeadline{})
		return
	}
	respond.JSON(w, http.StatusOK, deadlines)
}

// CreateDeadline POST /v1/planner/deadlines
func (h *PlannerHandler) CreateDeadline(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body struct {
		Name       string `json:"name"`
		CourseName string `json:"course_name"`
		DueDate    string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.Service.CreateStudyDeadline(r.Context(), userID, schoolID, body.Name, body.CourseName, body.DueDate)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

// DeleteDeadline DELETE /v1/planner/deadlines/{id}
func (h *PlannerHandler) DeleteDeadline(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid deadline id")
		return
	}
	if err := h.Service.DeleteStudyDeadline(r.Context(), id, userID, schoolID); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegeneratePlan POST /v1/planner/regenerate
func (h *PlannerHandler) RegeneratePlan(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	sessions, err := h.Service.RegeneratePlan(r.Context(), userID, schoolID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, sessions)
}
