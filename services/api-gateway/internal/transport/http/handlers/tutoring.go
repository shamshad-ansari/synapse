package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

type TutoringHandler struct {
	Service domain.TutoringService
}

func (h *TutoringHandler) authIDs(r *http.Request) (userID, schoolID uuid.UUID, ok bool) {
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

func parseTutorRequestListStatus(q string, incoming bool) (string, bool) {
	s := strings.TrimSpace(strings.ToLower(q))
	allowed := map[string]bool{
		"pending": true, "accepted": true, "declined": true, "completed": true, "cancelled": true, "all": true,
	}
	if s == "" {
		if incoming {
			return "pending", true
		}
		return "all", true
	}
	if !allowed[s] {
		return "", false
	}
	return s, true
}

func (h *TutoringHandler) ListIncomingRequests(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, ok := parseTutorRequestListStatus(r.URL.Query().Get("status"), true)
	if !ok {
		respond.Error(w, http.StatusBadRequest, "invalid status; use pending|accepted|declined|completed|cancelled|all")
		return
	}

	out, err := h.Service.ListRequestsForTutor(r.Context(), userID, schoolID, status)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	respond.JSON(w, http.StatusOK, out)
}

func (h *TutoringHandler) ListOutgoingRequests(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, ok := parseTutorRequestListStatus(r.URL.Query().Get("status"), false)
	if !ok {
		respond.Error(w, http.StatusBadRequest, "invalid status; use pending|accepted|declined|completed|cancelled|all")
		return
	}

	out, err := h.Service.ListRequestsByRequester(r.Context(), userID, schoolID, status)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	respond.JSON(w, http.StatusOK, out)
}

func (h *TutoringHandler) ListTeachingTopics(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	out, err := h.Service.ListTeachingTopics(r.Context(), userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	respond.JSON(w, http.StatusOK, out)
}

type createTutoringRequestBody struct {
	TutorID   uuid.UUID  `json:"tutor_id"`
	TopicName string     `json:"topic_name"`
	Message   string     `json:"message"`
	TopicID   *uuid.UUID `json:"topic_id"`
}

func (h *TutoringHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body createTutoringRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.TutorID == uuid.Nil {
		respond.Error(w, http.StatusBadRequest, "tutor_id is required")
		return
	}
	if body.TopicName == "" {
		respond.Error(w, http.StatusBadRequest, "topic_name is required")
		return
	}

	out, err := h.Service.CreateRequest(r.Context(), schoolID, userID, body.TutorID, body.TopicName, body.Message, body.TopicID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	respond.JSON(w, http.StatusCreated, out)
}

func (h *TutoringHandler) UpdateRequestStatus(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	requestID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request id")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	allowed := map[string]bool{
		"accepted":  true,
		"declined":  true,
		"completed": true,
		"cancelled": true,
	}
	if !allowed[body.Status] {
		respond.Error(w, http.StatusBadRequest, "invalid status; must be accepted|declined|completed|cancelled")
		return
	}

	out, err := h.Service.UpdateRequestStatus(r.Context(), requestID, userID, schoolID, body.Status)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	respond.JSON(w, http.StatusOK, out)
}

func (h *TutoringHandler) FindTutorMatches(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		respond.Error(w, http.StatusBadRequest, "topic query param is required")
		return
	}

	limit := 5
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	out, err := h.Service.FindTutorMatches(r.Context(), schoolID, userID, topic, limit)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	respond.JSON(w, http.StatusOK, out)
}
