package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

type tutoringServiceImpl struct {
	repo domain.TutoringRepository
}

// NewTutoringService returns a TutoringService backed by the repository.
func NewTutoringService(repo domain.TutoringRepository) domain.TutoringService {
	return &tutoringServiceImpl{repo: repo}
}

func (s *tutoringServiceImpl) CreateRequest(ctx context.Context, schoolID, requesterID, tutorID uuid.UUID, topicName, message string, topicID *uuid.UUID) (*domain.TutorRequest, error) {
	return s.repo.CreateRequest(ctx, schoolID, requesterID, tutorID, topicName, message, topicID)
}

func (s *tutoringServiceImpl) ListRequestsForTutor(ctx context.Context, tutorID, schoolID uuid.UUID, status string) ([]domain.TutorRequestView, error) {
	return s.repo.ListRequestsForTutor(ctx, tutorID, schoolID, status)
}

func (s *tutoringServiceImpl) ListRequestsByRequester(ctx context.Context, requesterID, schoolID uuid.UUID, status string) ([]domain.TutorRequestView, error) {
	return s.repo.ListRequestsByRequester(ctx, requesterID, schoolID, status)
}

func (s *tutoringServiceImpl) UpdateRequestStatus(ctx context.Context, requestID, userID, schoolID uuid.UUID, status string) (*domain.TutorRequest, error) {
	return s.repo.UpdateRequestStatus(ctx, requestID, userID, schoolID, status)
}

func (s *tutoringServiceImpl) FindTutorMatches(ctx context.Context, schoolID, requesterID uuid.UUID, topicName string, limit int) ([]domain.TutorMatch, error) {
	return s.repo.FindTutorMatches(ctx, schoolID, requesterID, topicName, limit)
}

func (s *tutoringServiceImpl) ListTeachingTopics(ctx context.Context, userID, schoolID uuid.UUID) ([]domain.TeachingTopic, error) {
	return s.repo.ListTeachingTopics(ctx, userID, schoolID)
}
