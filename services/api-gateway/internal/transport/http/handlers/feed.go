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

type FeedHandler struct {
	Service domain.FeedService
}

func (h *FeedHandler) authIDs(r *http.Request) (userID, schoolID uuid.UUID, ok bool) {
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

// ListFeedPosts GET /v1/feed
func (h *FeedHandler) ListFeedPosts(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := 20
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			respond.Error(w, http.StatusBadRequest, "invalid limit")
			return
		}
		if n > 50 {
			n = 50
		}
		limit = n
	}

	offset := 0
	if raw := r.URL.Query().Get("offset"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			respond.Error(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = n
	}

	out, err := h.Service.ListPosts(r.Context(), schoolID, userID, limit, offset)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

type createFeedPostBody struct {
	Title    string  `json:"title"`
	Body     string  `json:"body"`
	PostType string  `json:"post_type"`
	CourseID *string `json:"course_id"`
	TopicID  *string `json:"topic_id"`
}

// CreateFeedPost POST /v1/feed
func (h *FeedHandler) CreateFeedPost(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body createFeedPostBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(body.Title) == "" || strings.TrimSpace(body.Body) == "" {
		st, msg := service.MapLearningError(&domain.ValidationError{Message: "title and body are required"})
		respond.Error(w, st, msg)
		return
	}

	postType := body.PostType
	if postType == "" {
		postType = "question"
	}

	var courseID *uuid.UUID
	if body.CourseID != nil && *body.CourseID != "" {
		id, err := uuid.Parse(*body.CourseID)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid course_id")
			return
		}
		courseID = &id
	}
	var topicID *uuid.UUID
	if body.TopicID != nil && *body.TopicID != "" {
		id, err := uuid.Parse(*body.TopicID)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid topic_id")
			return
		}
		topicID = &id
	}

	out, err := h.Service.CreatePost(r.Context(), schoolID, userID, body.Title, body.Body, postType, courseID, topicID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

// ToggleFeedUpvote POST /v1/feed/{postId}/upvote
func (h *FeedHandler) ToggleFeedUpvote(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	postID, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid post id")
		return
	}

	upvotes, upvoted, err := h.Service.ToggleUpvote(r.Context(), postID, userID, schoolID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{
		"upvotes": upvotes,
		"upvoted": upvoted,
	})
}

// CreateFeedComment POST /v1/feed/{postId}/comments
func (h *FeedHandler) CreateFeedComment(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	postID, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid post id")
		return
	}
	var body struct {
		Body     string  `json:"body"`
		ParentID *string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(body.Body) == "" {
		respond.Error(w, http.StatusBadRequest, "body is required")
		return
	}
	var parentID *uuid.UUID
	if body.ParentID != nil && *body.ParentID != "" {
		pid, err := uuid.Parse(*body.ParentID)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid parent_id")
			return
		}
		parentID = &pid
	}
	out, err := h.Service.CreateComment(r.Context(), schoolID, postID, userID,
		strings.TrimSpace(body.Body), parentID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusCreated, out)
}

// ListFeedComments GET /v1/feed/{postId}/comments
func (h *FeedHandler) ListFeedComments(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	_ = userID
	postID, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid post id")
		return
	}
	out, err := h.Service.ListComments(r.Context(), schoolID, postID)
	if err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	respond.JSON(w, http.StatusOK, out)
}

// DeleteFeedPost DELETE /v1/feed/{postId}
func (h *FeedHandler) DeleteFeedPost(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	postID, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid post id")
		return
	}
	if err := h.Service.DeletePost(r.Context(), postID, userID, schoolID); err != nil {
		st, msg := service.MapLearningError(err)
		respond.Error(w, st, msg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
