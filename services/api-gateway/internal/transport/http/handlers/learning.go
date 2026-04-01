package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/ai"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

// LearningHandler exposes HTTP handlers for courses, topics, flashcards, reviews, and notes.
type LearningHandler struct {
	Service  service.LearningService
	Repo     domain.LearningRepository
	Embed    ai.TextEmbedder
	AIClient ai.Completer
}

func (h *LearningHandler) authIDs(r *http.Request) (userID, schoolID uuid.UUID, ok bool) {
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

// --- Courses ---

func (h *LearningHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	out, err := h.Service.ListCourses(r.Context(), userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

type createCourseBody struct {
	Name  string `json:"name"`
	Term  string `json:"term"`
	Color string `json:"color"`
}

func (h *LearningHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body createCourseBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.Service.CreateCourse(r.Context(), userID, schoolID, body.Name, body.Term, body.Color)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

func (h *LearningHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	courseID, err := uuid.Parse(chi.URLParam(r, "courseId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid course id")
		return
	}
	out, err := h.Service.GetCourse(r.Context(), courseID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

func (h *LearningHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	courseID, err := uuid.Parse(chi.URLParam(r, "courseId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid course id")
		return
	}
	if err := h.Service.DeleteCourse(r.Context(), courseID, userID, schoolID); err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type importFromLMSBody struct {
	Courses []service.LMSCourseInput `json:"courses"`
}

func (h *LearningHandler) ImportFromLMS(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body importFromLMSBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.Service.ImportCoursesFromLMS(r.Context(), userID, schoolID, body.Courses)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

// --- Notes ---

type createNoteBody struct {
	CourseID uuid.UUID  `json:"course_id"`
	TopicID  *uuid.UUID `json:"topic_id"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
}

// CreateNote POST /v1/notes — persists note text; optionally async-embeds when Embed client is configured.
func (h *LearningHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body createNoteBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.CourseID == uuid.Nil {
		respond.Error(w, http.StatusBadRequest, "course_id is required")
		return
	}
	out, err := h.Service.CreateNoteText(r.Context(), userID, schoolID, body.CourseID, body.Title, body.Content, body.TopicID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	if h.Embed != nil && h.Repo != nil {
		noteID := out.ID
		sid := schoolID
		content := out.Content
		embed := h.Embed
		repo := h.Repo
		go func() {
			bg, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()
			vec, err := embed.Embed(bg, ai.EmbedInput{Text: content})
			if err != nil {
				return
			}
			_ = repo.UpdateNoteTextEmbedding(bg, noteID, sid, vec)
		}()
	}
	respond.JSON(w, http.StatusCreated, out)
}

// ListNotes GET /v1/courses/{courseId}/notes
func (h *LearningHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	courseID, err := uuid.Parse(chi.URLParam(r, "courseId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid course id")
		return
	}
	out, err := h.Service.ListNoteTexts(r.Context(), courseID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

// ListNoteMetrics GET /v1/courses/{courseId}/notes/metrics?windowDays=
func (h *LearningHandler) ListNoteMetrics(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	courseID, err := uuid.Parse(chi.URLParam(r, "courseId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid course id")
		return
	}
	windowDays := 30
	if raw := r.URL.Query().Get("windowDays"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			respond.Error(w, http.StatusBadRequest, "invalid windowDays")
			return
		}
		windowDays = n
	}
	out, err := h.Service.ListNoteMetrics(r.Context(), courseID, userID, schoolID, windowDays)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

// GetNote GET /v1/notes/{noteId}
func (h *LearningHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	noteID, err := uuid.Parse(chi.URLParam(r, "noteId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid note id")
		return
	}
	out, err := h.Service.GetNoteText(r.Context(), noteID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

// UpdateNote PUT /v1/notes/{noteId}
func (h *LearningHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	noteID, err := uuid.Parse(chi.URLParam(r, "noteId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid note id")
		return
	}
	var body struct {
		Title   string     `json:"title"`
		Content string     `json:"content"`
		TopicID *uuid.UUID `json:"topic_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.Service.UpdateNoteText(r.Context(), noteID, userID, schoolID, body.Title, body.Content, body.TopicID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

// DeleteNote DELETE /v1/notes/{noteId}
func (h *LearningHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	noteID, err := uuid.Parse(chi.URLParam(r, "noteId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid note id")
		return
	}
	if err := h.Service.DeleteNoteText(r.Context(), noteID, userID, schoolID); err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AskNoteAI POST /v1/notes/{noteId}/ask
// Accepts: { question: string }
// Returns: { "data": { "answer": string } }
func (h *LearningHandler) AskNoteAI(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	noteID, err := uuid.Parse(chi.URLParam(r, "noteId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid note id")
		return
	}

	var body struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Question == "" {
		respond.Error(w, http.StatusBadRequest, "question required")
		return
	}

	note, err := h.Service.GetNoteText(r.Context(), noteID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	courses, _ := h.Service.ListCourses(r.Context(), userID, schoolID)
	courseNames := make([]string, 0, len(courses))
	for _, c := range courses {
		courseNames = append(courseNames, c.Name)
	}

	topics, _ := h.Service.ListTopics(r.Context(), note.CourseID, userID, schoolID)
	topicNames := make([]string, 0, len(topics))
	for _, t := range topics {
		topicNames = append(topicNames, t.Name)
	}

	if h.AIClient == nil {
		respond.Error(w, http.StatusServiceUnavailable, "AI not configured")
		return
	}

	systemPrompt := fmt.Sprintf(`You are Synapse, an intelligent learning assistant embedded inside a student's note editor.

Student context:
- Enrolled courses: %s
- Topics in this course: %s
- Current note title: "%s"
- Current note content:
%s

Your job: answer the student's question clearly and concisely, using the note content as primary context.
Be specific, educational, and encouraging. If the question is about a concept in the note, explain it
in a way that complements what the student has written. Keep answers under 300 words unless a longer
explanation is genuinely needed.`,
		strings.Join(courseNames, ", "),
		strings.Join(topicNames, ", "),
		note.Title,
		note.Content,
	)

	answer, err := h.AIClient.Complete(r.Context(), systemPrompt, body.Question)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}

	respond.JSON(w, http.StatusOK, map[string]string{"answer": answer})
}

// --- Topics ---

func (h *LearningHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	courseID, err := uuid.Parse(chi.URLParam(r, "courseId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid course id")
		return
	}
	out, err := h.Service.ListTopics(r.Context(), courseID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

type createTopicBody struct {
	Name       string   `json:"name"`
	ExamWeight *float64 `json:"exam_weight"`
}

func (h *LearningHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	courseID, err := uuid.Parse(chi.URLParam(r, "courseId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid course id")
		return
	}
	var body createTopicBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.Service.CreateTopic(r.Context(), userID, schoolID, courseID, body.Name, body.ExamWeight)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

// --- Flashcards ---

func (h *LearningHandler) ListFlashcards(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	courseID, err := uuid.Parse(chi.URLParam(r, "courseId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid course id")
		return
	}
	out, err := h.Service.ListFlashcards(r.Context(), courseID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

type createFlashcardBody struct {
	CourseID uuid.UUID  `json:"course_id"`
	TopicID  *uuid.UUID `json:"topic_id"`
	CardType string     `json:"card_type"`
	Prompt   string     `json:"prompt"`
	Answer   string     `json:"answer"`
}

func (h *LearningHandler) CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body createFlashcardBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.CourseID == uuid.Nil {
		respond.Error(w, http.StatusBadRequest, "course_id is required")
		return
	}
	out, err := h.Service.CreateFlashcard(r.Context(), userID, schoolID, body.CourseID, body.TopicID, body.CardType, body.Prompt, body.Answer)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

func (h *LearningHandler) DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	cardID, err := uuid.Parse(chi.URLParam(r, "cardId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid flashcard id")
		return
	}
	if err := h.Service.DeleteFlashcard(r.Context(), cardID, userID, schoolID); err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Review ---

func (h *LearningHandler) GetDueCards(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	q := r.URL.Query()
	limit := 20
	if ls := q.Get("limit"); ls != "" {
		n, err := strconv.Atoi(ls)
		if err != nil || n < 1 {
			respond.Error(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = n
	}
	var courseID *uuid.UUID
	if cs := q.Get("courseId"); cs != "" {
		id, err := uuid.Parse(cs)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid course id")
			return
		}
		courseID = &id
	}
	out, err := h.Service.GetDueCards(r.Context(), userID, schoolID, limit, courseID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

type submitReviewBody struct {
	SessionID      uuid.UUID `json:"session_id"`
	FlashcardID    uuid.UUID `json:"flashcard_id"`
	Correct        bool      `json:"correct"`
	Confidence     int       `json:"confidence"`
	Confused       bool      `json:"confused"`
	ResponseTimeMs int       `json:"response_time_ms"`
}

func (h *LearningHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body submitReviewBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.SessionID == uuid.Nil || body.FlashcardID == uuid.Nil {
		respond.Error(w, http.StatusBadRequest, "session_id and flashcard_id are required")
		return
	}
	out, err := h.Service.SubmitReview(r.Context(), userID, schoolID, body.FlashcardID, body.SessionID, body.Correct, body.Confidence, body.Confused, body.ResponseTimeMs)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

// GetConfusionInsights GET /v1/insights/confusion?courseId=&windowDays=
func (h *LearningHandler) GetConfusionInsights(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var (
		courseID   *uuid.UUID
		windowDays = 14
	)
	if raw := r.URL.Query().Get("courseId"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid course id")
			return
		}
		courseID = &id
	}
	if raw := r.URL.Query().Get("windowDays"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			respond.Error(w, http.StatusBadRequest, "invalid windowDays")
			return
		}
		windowDays = n
	}

	out, err := h.Service.GetConfusionInsights(r.Context(), userID, schoolID, courseID, windowDays)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}
