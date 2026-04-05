package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

type feedServiceImpl struct {
	repo domain.FeedRepository
}

func NewFeedService(repo domain.FeedRepository) domain.FeedService {
	return &feedServiceImpl{repo: repo}
}

func (s *feedServiceImpl) CreatePost(ctx context.Context, schoolID, userID uuid.UUID, title, body, postType string, courseID, topicID *uuid.UUID) (*domain.FeedPostView, error) {
	out, err := s.repo.CreatePost(ctx, schoolID, userID, title, body, postType, courseID, topicID)
	if err != nil {
		return nil, fmt.Errorf("CreatePost: %w", err)
	}
	return out, nil
}

func (s *feedServiceImpl) ListPosts(ctx context.Context, schoolID, requestingUserID uuid.UUID, limit, offset int) ([]domain.FeedPostView, error) {
	out, err := s.repo.ListPosts(ctx, schoolID, requestingUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ListPosts: %w", err)
	}
	return out, nil
}

func (s *feedServiceImpl) ToggleUpvote(ctx context.Context, postID, userID, schoolID uuid.UUID) (int, bool, error) {
	n, on, err := s.repo.ToggleUpvote(ctx, postID, userID, schoolID)
	if err != nil {
		return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
	}
	return n, on, nil
}

func (s *feedServiceImpl) CreateComment(ctx context.Context, schoolID, postID, userID uuid.UUID, body string, parentID *uuid.UUID) (*domain.FeedCommentView, error) {
	if strings.TrimSpace(body) == "" {
		return nil, &domain.ValidationError{Message: "body is required"}
	}
	return s.repo.CreateComment(ctx, schoolID, postID, userID, strings.TrimSpace(body), parentID)
}

func (s *feedServiceImpl) ListComments(ctx context.Context, schoolID, postID uuid.UUID) ([]domain.FeedCommentView, error) {
	return s.repo.ListComments(ctx, schoolID, postID)
}

func (s *feedServiceImpl) DeletePost(ctx context.Context, postID, userID, schoolID uuid.UUID) error {
	return s.repo.DeletePost(ctx, postID, userID, schoolID)
}
