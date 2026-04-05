package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

// ProfileService loads the profile summary for the dashboard.
type ProfileService interface {
	GetSummary(ctx context.Context, userID, schoolID uuid.UUID) (*domain.ProfileSummary, error)
}

type profileServiceImpl struct {
	repo domain.ProfileRepository
}

func NewProfileService(repo domain.ProfileRepository) ProfileService {
	return &profileServiceImpl{repo: repo}
}

func (s *profileServiceImpl) GetSummary(ctx context.Context, userID, schoolID uuid.UUID) (*domain.ProfileSummary, error) {
	out, err := s.repo.GetSummary(ctx, userID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("GetSummary: %w", err)
	}
	return out, nil
}
