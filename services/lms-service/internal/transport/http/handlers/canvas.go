package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/lms-service/internal/canvas"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/config"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/crypto"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/domain"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/oauth"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/sync"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/transport/respond"
)

type CanvasHandler struct {
	Cfg    *config.Config
	Repo   domain.LMSRepository
	Syncer *sync.Syncer
	Redis  *redis.Client
	Logger *zap.Logger
}

type oauthState struct {
	UserID         string `json:"user_id"`
	SchoolID       string `json:"school_id"`
	State          string `json:"state"`
	InstitutionURL string `json:"institution_url"`
}

// ConnectCanvas initiates the Canvas OAuth2 flow.
// GET /v1/lms/connect/canvas?institution_url={url}
func (h *CanvasHandler) ConnectCanvas(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	schoolID, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	institutionURL := r.URL.Query().Get("institution_url")
	if institutionURL == "" {
		respond.Error(w, http.StatusBadRequest, "institution_url is required")
		return
	}

	parsed, err := url.ParseRequestURI(institutionURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		respond.Error(w, http.StatusBadRequest, "institution_url must be a valid URL")
		return
	}

	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		h.Logger.Error("failed to generate state", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	stateToken := base64.URLEncoding.EncodeToString(stateBytes)

	redisKey := fmt.Sprintf("oauth:state:%s:%s", userID.String(), stateToken)
	if err := h.Redis.Set(r.Context(), redisKey, "1", 10*time.Minute).Err(); err != nil {
		h.Logger.Error("failed to store oauth state in redis", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	statePayload := oauthState{
		UserID:         userID.String(),
		SchoolID:       schoolID.String(),
		State:          stateToken,
		InstitutionURL: institutionURL,
	}
	stateJSON, err := json.Marshal(statePayload)
	if err != nil {
		h.Logger.Error("failed to marshal state", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	encodedState := base64.URLEncoding.EncodeToString(stateJSON)

	authURL := oauth.BuildAuthURL(institutionURL, h.Cfg.CanvasClientID, h.Cfg.CanvasRedirectURI, encodedState)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// CallbackCanvas handles the Canvas OAuth2 callback.
// GET /v1/lms/callback/canvas?code={code}&state={state_json}
func (h *CanvasHandler) CallbackCanvas(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	stateParam := r.URL.Query().Get("state")

	if code == "" || stateParam == "" {
		h.redirectError(w, r)
		return
	}

	stateJSON, err := base64.URLEncoding.DecodeString(stateParam)
	if err != nil {
		h.Logger.Warn("invalid state encoding", zap.Error(err))
		h.redirectError(w, r)
		return
	}

	var state oauthState
	if err := json.Unmarshal(stateJSON, &state); err != nil {
		h.Logger.Warn("invalid state json", zap.Error(err))
		h.redirectError(w, r)
		return
	}

	userID, err := uuid.Parse(state.UserID)
	if err != nil {
		h.Logger.Warn("invalid user_id in state", zap.Error(err))
		h.redirectError(w, r)
		return
	}

	schoolID, err := uuid.Parse(state.SchoolID)
	if err != nil {
		h.Logger.Warn("invalid school_id in state", zap.Error(err))
		h.redirectError(w, r)
		return
	}

	redisKey := fmt.Sprintf("oauth:state:%s:%s", state.UserID, state.State)
	val, err := h.Redis.Get(r.Context(), redisKey).Result()
	if err != nil || val != "1" {
		h.Logger.Warn("oauth state not found or expired", zap.String("key", redisKey))
		h.redirectError(w, r)
		return
	}
	h.Redis.Del(r.Context(), redisKey)

	serverURL := canvas.ResolveServerURL(state.InstitutionURL, h.Cfg.CanvasInternalURL)
	tokenResp, err := oauth.ExchangeCode(
		r.Context(),
		serverURL,
		h.Cfg.CanvasClientID,
		h.Cfg.CanvasClientSecret,
		code,
		h.Cfg.CanvasRedirectURI,
	)
	if err != nil {
		h.Logger.Error("failed to exchange code", zap.Error(err))
		h.redirectError(w, r)
		return
	}

	encAccessToken, err := crypto.Encrypt([]byte(tokenResp.AccessToken), h.Cfg.EncryptionKey)
	if err != nil {
		h.Logger.Error("failed to encrypt access token", zap.Error(err))
		h.redirectError(w, r)
		return
	}
	encRefreshToken, err := crypto.Encrypt([]byte(tokenResp.RefreshToken), h.Cfg.EncryptionKey)
	if err != nil {
		h.Logger.Error("failed to encrypt refresh token", zap.Error(err))
		h.redirectError(w, r)
		return
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	conn := &domain.LMSConnection{
		UserID:         userID,
		SchoolID:       schoolID,
		LMSType:        "canvas",
		InstitutionURL: state.InstitutionURL,
		AccessToken:    encAccessToken,
		RefreshToken:   encRefreshToken,
		TokenExpiresAt: expiresAt,
		SyncStatus:     "pending",
	}

	if err := h.Repo.UpsertConnection(r.Context(), conn); err != nil {
		h.Logger.Error("failed to upsert connection", zap.Error(err))
		h.redirectError(w, r)
		return
	}

	redirectURL := h.Cfg.FrontendURL + "/canvas/connected?status=success"
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// ConnectToken connects via a personal access token (no OAuth required).
// POST /v1/lms/connect/token
func (h *CanvasHandler) ConnectToken(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	schoolID, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		InstitutionURL string `json:"institution_url"`
		AccessToken    string `json:"access_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.InstitutionURL == "" {
		respond.Error(w, http.StatusBadRequest, "institution_url is required")
		return
	}
	parsed, err := url.ParseRequestURI(body.InstitutionURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		respond.Error(w, http.StatusBadRequest, "institution_url must be a valid URL")
		return
	}
	if body.AccessToken == "" {
		respond.Error(w, http.StatusBadRequest, "access_token is required")
		return
	}

	// Verify the token actually works by calling Canvas /api/v1/users/self
	verifyURL := body.InstitutionURL + "/api/v1/users/self"
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, verifyURL, nil)
	if err != nil {
		h.Logger.Error("failed to build verify request", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)

	verifyResp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.Logger.Error("failed to verify canvas token", zap.Error(err))
		respond.Error(w, http.StatusBadGateway, "could not reach Canvas instance")
		return
	}
	defer verifyResp.Body.Close()

	if verifyResp.StatusCode != http.StatusOK {
		respond.Error(w, http.StatusUnauthorized, "Canvas rejected the token — verify it in Canvas > Account > Settings > Access Tokens")
		return
	}

	encAccessToken, err := crypto.Encrypt([]byte(body.AccessToken), h.Cfg.EncryptionKey)
	if err != nil {
		h.Logger.Error("failed to encrypt access token", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Personal tokens don't have refresh tokens — store a placeholder
	encRefreshToken, err := crypto.Encrypt([]byte("personal_token_no_refresh"), h.Cfg.EncryptionKey)
	if err != nil {
		h.Logger.Error("failed to encrypt refresh placeholder", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Personal tokens don't expire unless manually revoked — set far-future expiry
	conn := &domain.LMSConnection{
		UserID:         userID,
		SchoolID:       schoolID,
		LMSType:        "canvas",
		InstitutionURL: body.InstitutionURL,
		AccessToken:    encAccessToken,
		RefreshToken:   encRefreshToken,
		TokenExpiresAt: time.Now().Add(10 * 365 * 24 * time.Hour),
		SyncStatus:     "pending",
	}

	if err := h.Repo.UpsertConnection(r.Context(), conn); err != nil {
		h.Logger.Error("failed to upsert connection", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	respond.JSON(w, http.StatusCreated, conn.ToResponse())
}

// Status returns the current LMS connection status for the authenticated user.
// GET /v1/lms/status
func (h *CanvasHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	schoolID, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	conn, err := h.Repo.FindConnectionByUser(r.Context(), userID, schoolID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, "not connected")
			return
		}
		h.Logger.Error("failed to find connection", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	respond.JSON(w, http.StatusOK, conn.ToResponse())
}

// Sync runs a full Canvas → Postgres sync for the authenticated user's LMS connection.
// POST /v1/lms/sync
func (h *CanvasHandler) Sync(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	schoolID, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	conn, err := h.Repo.FindConnectionByUser(r.Context(), userID, schoolID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, "not connected")
			return
		}
		h.Logger.Error("failed to verify connection", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.Syncer.SyncUser(r.Context(), conn, sync.SyncOptions{
		InternalURL:        h.Cfg.CanvasInternalURL,
		CanvasClientID:     h.Cfg.CanvasClientID,
		CanvasClientSecret: h.Cfg.CanvasClientSecret,
	}); err != nil {
		h.Logger.Error("lms sync failed", zap.Error(err))
		respond.Error(w, http.StatusBadGateway, "sync failed")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]string{"status": "sync_complete"})
}

// ListSyncedCourses returns courses last synced from Canvas for the authenticated user.
// GET /v1/lms/courses
func (h *CanvasHandler) ListSyncedCourses(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	schoolID, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	courses, err := h.Repo.ListCoursesByUser(r.Context(), userID, schoolID)
	if err != nil {
		h.Logger.Error("failed to list synced courses", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	out := make([]domain.LMSCourseResponse, 0, len(courses))
	for _, c := range courses {
		out = append(out, domain.LMSCourseResponse{
			LMSCourseID:  c.LMSCourseID,
			CourseName:   c.LMSCourseName,
			Term:         c.LMSTerm,
			LastSyncedAt: c.LastSyncedAt,
		})
	}

	respond.JSON(w, http.StatusOK, out)
}

// Disconnect removes the authenticated user's LMS connection.
// DELETE /v1/lms/disconnect
func (h *CanvasHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	schoolID, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.Repo.DeleteConnection(context.Background(), userID, schoolID); err != nil {
		h.Logger.Error("failed to delete connection", zap.Error(err))
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	respond.JSON(w, http.StatusOK, "disconnected")
}

func (h *CanvasHandler) redirectError(w http.ResponseWriter, r *http.Request) {
	redirectURL := h.Cfg.FrontendURL + "/canvas/connected?status=error"
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
