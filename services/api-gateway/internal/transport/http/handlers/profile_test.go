package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
)

type mockProfileSvc struct {
	out *domain.ProfileSummary
	err error
}

func (m *mockProfileSvc) GetSummary(ctx context.Context, userID, schoolID uuid.UUID) (*domain.ProfileSummary, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.out, nil
}

var _ service.ProfileService = (*mockProfileSvc)(nil)

func newProfileTestRouter(h *ProfileHandler) http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(v1 chi.Router) {
		v1.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(testTutoringJWTSecret))
			protected.Get("/profile/summary", h.Summary)
		})
	})
	return r
}

func TestProfileHandler_Summary_Unauthorized(t *testing.T) {
	h := &ProfileHandler{Service: &mockProfileSvc{}}
	srv := newProfileTestRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/v1/profile/summary", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProfileHandler_Summary_OK(t *testing.T) {
	summary := &domain.ProfileSummary{
		ReputationTotal: 42,
		CardsShared:     1,
		SessionsDone:    2,
		Breakdown: []domain.ProfileRepBreakdownRow{
			{Label: "Flashcards shared", Value: 1, Percent: 100, Color: "var(--navy)"},
		},
		Mastery: []domain.ProfileMasteryRow{
			{Name: "Topic A", Mastery: 50, Band: "amber"},
		},
		Contributions: []domain.ProfileContributionRow{
			{
				Icon: "zap", IconColor: "var(--navy)", IconBg: "var(--navy-light)",
				Title: "Flashcards shared", Subtitle: "Visible", Value: "1", ValueColor: "var(--navy)",
			},
		},
	}
	h := &ProfileHandler{Service: &mockProfileSvc{out: summary}}
	srv := newProfileTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/v1/profile/summary", nil)
	req.Header.Set("Authorization", "Bearer "+testTutoringAccessToken(t, testTutoringUserID, testTutoringSchoolID))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var env struct {
		Data *domain.ProfileSummary `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if env.Data == nil || env.Data.ReputationTotal != 42 {
		t.Fatalf("unexpected payload: %+v", env.Data)
	}
}

func TestProfileHandler_Summary_InternalError(t *testing.T) {
	h := &ProfileHandler{Service: &mockProfileSvc{err: context.DeadlineExceeded}}
	srv := newProfileTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/v1/profile/summary", nil)
	req.Header.Set("Authorization", "Bearer "+testTutoringAccessToken(t, testTutoringUserID, testTutoringSchoolID))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}
