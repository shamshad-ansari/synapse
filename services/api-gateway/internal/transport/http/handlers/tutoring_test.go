package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
)

// Mirrors manual checks against a running stack:
//   POST /v1/tutoring/requests (self-request) → 201 + {"data":{...},"error":null}
//   GET  /v1/tutoring/requests/incoming → 200 + {"data":[...]} (pending only)
//   GET  /v1/tutoring/match?topic=... → 200 + {"data":[]}

const testTutoringJWTSecret = "test-jwt-secret-for-tutoring-handlers"

var (
	testTutoringUserID   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testTutoringSchoolID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func testTutoringAccessToken(t *testing.T, userID, schoolID uuid.UUID) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":      "access",
		"user_id":   userID.String(),
		"school_id": schoolID.String(),
	})
	s, err := tok.SignedString([]byte(testTutoringJWTSecret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func newTutoringTestRouter(h *TutoringHandler) http.Handler {
	r := chi.NewRouter()
	r.Route("/v1", func(v1 chi.Router) {
		v1.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(testTutoringJWTSecret))
			protected.Route("/tutoring", func(tr chi.Router) {
				tr.Get("/teaching-topics", h.ListTeachingTopics)
				tr.Get("/requests/incoming", h.ListIncomingRequests)
				tr.Get("/requests/outgoing", h.ListOutgoingRequests)
				tr.Post("/requests", h.CreateRequest)
				tr.Patch("/requests/{id}", h.UpdateRequestStatus)
				tr.Get("/match", h.FindTutorMatches)
			})
		})
	})
	return r
}

type mockTutoringSvc struct {
	createOut          *domain.TutorRequest
	createErr          error
	listForTutor       []domain.TutorRequestView
	listForTutorErr    error
	listByRequester    []domain.TutorRequestView
	listByRequesterErr error
	updateOut          *domain.TutorRequest
	updateErr          error
	matchOut           []domain.TutorMatch
	matchErr           error
	teachingOut        []domain.TeachingTopic
	teachingErr        error
}

func (m *mockTutoringSvc) CreateRequest(ctx context.Context, schoolID, requesterID, tutorID uuid.UUID, topicName, message string, topicID *uuid.UUID) (*domain.TutorRequest, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createOut != nil {
		return m.createOut, nil
	}
	now := time.Now().UTC()
	return &domain.TutorRequest{
		ID:          uuid.New(),
		SchoolID:    schoolID,
		RequesterID: requesterID,
		TutorID:     tutorID,
		TopicName:   topicName,
		Message:     message,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (m *mockTutoringSvc) ListRequestsForTutor(ctx context.Context, tutorID, schoolID uuid.UUID, status string) ([]domain.TutorRequestView, error) {
	if m.listForTutorErr != nil {
		return nil, m.listForTutorErr
	}
	st := strings.TrimSpace(strings.ToLower(status))
	if st == "" {
		st = "pending"
	}
	if st == "all" {
		return m.listForTutor, nil
	}
	var out []domain.TutorRequestView
	for _, v := range m.listForTutor {
		if strings.EqualFold(v.Status, st) {
			out = append(out, v)
		}
	}
	return out, nil
}

func (m *mockTutoringSvc) ListRequestsByRequester(ctx context.Context, requesterID, schoolID uuid.UUID, status string) ([]domain.TutorRequestView, error) {
	if m.listByRequesterErr != nil {
		return nil, m.listByRequesterErr
	}
	st := strings.TrimSpace(strings.ToLower(status))
	if st == "" || st == "all" {
		return m.listByRequester, nil
	}
	var out []domain.TutorRequestView
	for _, v := range m.listByRequester {
		if strings.EqualFold(v.Status, st) {
			out = append(out, v)
		}
	}
	return out, nil
}

func (m *mockTutoringSvc) UpdateRequestStatus(ctx context.Context, requestID, userID, schoolID uuid.UUID, status string) (*domain.TutorRequest, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if m.updateOut != nil {
		return m.updateOut, nil
	}
	return &domain.TutorRequest{}, nil
}

func (m *mockTutoringSvc) FindTutorMatches(ctx context.Context, schoolID, requesterID uuid.UUID, topicName string, limit int) ([]domain.TutorMatch, error) {
	if m.matchErr != nil {
		return nil, m.matchErr
	}
	if m.matchOut != nil {
		return m.matchOut, nil
	}
	return []domain.TutorMatch{}, nil
}

func (m *mockTutoringSvc) ListTeachingTopics(ctx context.Context, userID, schoolID uuid.UUID) ([]domain.TeachingTopic, error) {
	if m.teachingErr != nil {
		return nil, m.teachingErr
	}
	return m.teachingOut, nil
}

func decodeRespondEnvelope(t *testing.T, body []byte) (data json.RawMessage, errStr *string) {
	t.Helper()
	var env struct {
		Data  json.RawMessage `json:"data"`
		Error *string         `json:"error"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("response JSON: %v", err)
	}
	return env.Data, env.Error
}

func TestTutoring_CreateRequest_SelfRequest_201Envelope(t *testing.T) {
	mock := &mockTutoringSvc{}
	h := &TutoringHandler{Service: mock}
	srv := newTutoringTestRouter(h)

	token := testTutoringAccessToken(t, testTutoringUserID, testTutoringSchoolID)
	// Self-request: tutor_id equals authenticated user (same as curl with MY_USER_ID).
	body := `{"tutor_id":"` + testTutoringUserID.String() + `","topic_name":"Recursion","message":"Need help with base cases"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tutoring/requests", bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rr.Code, rr.Body.String())
	}
	data, errPtr := decodeRespondEnvelope(t, rr.Body.Bytes())
	if errPtr != nil {
		t.Fatalf("expected error null, got %v", *errPtr)
	}
	var row domain.TutorRequest
	if err := json.Unmarshal(data, &row); err != nil {
		t.Fatalf("data: %v", err)
	}
	if row.TopicName != "Recursion" || row.Message != "Need help with base cases" {
		t.Fatalf("unexpected row: %+v", row)
	}
	if row.RequesterID != testTutoringUserID || row.TutorID != testTutoringUserID || row.SchoolID != testTutoringSchoolID {
		t.Fatalf("ids: %+v", row)
	}
}

func TestTutoring_ListIncoming_PendingOnly_200Envelope(t *testing.T) {
	rid := uuid.New()
	mock := &mockTutoringSvc{
		listForTutor: []domain.TutorRequestView{
			{
				TutorRequest: domain.TutorRequest{
					ID:          rid,
					SchoolID:    testTutoringSchoolID,
					RequesterID: uuid.New(),
					TutorID:     testTutoringUserID,
					TopicName:   "Recursion",
					Status:      "pending",
					Message:     "help",
				},
				RequesterName: "A",
				TutorName:     "B",
			},
			{
				TutorRequest: domain.TutorRequest{
					ID:     uuid.New(),
					Status: "accepted",
				},
				RequesterName: "X",
				TutorName:     "Y",
			},
		},
	}
	h := &TutoringHandler{Service: mock}
	srv := newTutoringTestRouter(h)

	token := testTutoringAccessToken(t, testTutoringUserID, testTutoringSchoolID)
	req := httptest.NewRequest(http.MethodGet, "/v1/tutoring/requests/incoming", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rr.Code, rr.Body.String())
	}
	data, errPtr := decodeRespondEnvelope(t, rr.Body.Bytes())
	if errPtr != nil {
		t.Fatalf("expected error null, got %v", *errPtr)
	}
	var list []domain.TutorRequestView
	if err := json.Unmarshal(data, &list); err != nil {
		t.Fatalf("data: %v", err)
	}
	if len(list) != 1 || list[0].Status != "pending" || list[0].ID != rid {
		t.Fatalf("expected single pending row, got %+v", list)
	}
}

func TestTutoring_FindMatches_EmptyStub_200Envelope(t *testing.T) {
	mock := &mockTutoringSvc{}
	h := &TutoringHandler{Service: mock}
	srv := newTutoringTestRouter(h)

	token := testTutoringAccessToken(t, testTutoringUserID, testTutoringSchoolID)
	req := httptest.NewRequest(http.MethodGet, "/v1/tutoring/match?topic=Recursion", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rr.Code, rr.Body.String())
	}
	data, errPtr := decodeRespondEnvelope(t, rr.Body.Bytes())
	if errPtr != nil {
		t.Fatalf("expected error null, got %v", *errPtr)
	}
	// jq-friendly: data is JSON array, empty until mastery-backed matching exists
	if string(data) != "[]" {
		t.Fatalf("expected data [], got %s", string(data))
	}
}
