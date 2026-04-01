package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/ai"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

// AIHandler serves RAG flashcard generation. AIService may be nil → 503 on generate paths.
type AIHandler struct {
	AIService *service.AIService
}

func (h *AIHandler) authIDs(r *http.Request) (userID, schoolID uuid.UUID, ok bool) {
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

type generateFlashcardsBody struct {
	CourseID    uuid.UUID  `json:"course_id"`
	NoteContent string     `json:"note_content"`
	TopicID     *uuid.UUID `json:"topic_id"`
}

type acceptGeneratedBody struct {
	CourseID uuid.UUID `json:"course_id"`
	TopicID  *uuid.UUID `json:"topic_id"`
	Cards    []struct {
		Prompt   string `json:"prompt"`
		Answer   string `json:"answer"`
		CardType string `json:"card_type"`
	} `json:"cards"`
}

// GenerateFlashcards POST /v1/flashcards/generate — returns candidates only (does not persist).
func (h *AIHandler) GenerateFlashcards(w http.ResponseWriter, r *http.Request) {
	if h.AIService == nil {
		respond.Error(w, http.StatusServiceUnavailable, "AI not configured")
		return
	}
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body generateFlashcardsBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.CourseID == uuid.Nil {
		respond.Error(w, http.StatusBadRequest, "course_id is required")
		return
	}
	if body.NoteContent == "" {
		respond.Error(w, http.StatusBadRequest, "note_content is required")
		return
	}

	course, err := h.AIService.Repo.GetCourse(r.Context(), body.CourseID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	candidates, err := h.AIService.GenerateFlashcardsFromNote(
		r.Context(), userID, schoolID, body.CourseID, body.NoteContent, course.Name, body.TopicID,
	)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"candidates": candidates})
}

// AcceptGeneratedFlashcards POST /v1/flashcards/generate/accept — persists accepted cards.
func (h *AIHandler) AcceptGeneratedFlashcards(w http.ResponseWriter, r *http.Request) {
	if h.AIService == nil {
		respond.Error(w, http.StatusServiceUnavailable, "AI not configured")
		return
	}
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body acceptGeneratedBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.CourseID == uuid.Nil {
		respond.Error(w, http.StatusBadRequest, "course_id is required")
		return
	}
	if len(body.Cards) == 0 {
		respond.Error(w, http.StatusBadRequest, "cards is required")
		return
	}

	cards := make([]ai.GeneratedCard, 0, len(body.Cards))
	for _, c := range body.Cards {
		ct := c.CardType
		if ct == "" {
			ct = "qa"
		}
		cards = append(cards, ai.GeneratedCard{
			Prompt:   c.Prompt,
			Answer:   c.Answer,
			CardType: ct,
		})
	}

	saved, err := h.AIService.SaveAcceptedCards(r.Context(), userID, schoolID, body.CourseID, cards, body.TopicID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]any{"saved": saved})
}
