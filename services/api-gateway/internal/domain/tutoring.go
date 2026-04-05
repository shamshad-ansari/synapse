package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TutorRequest struct {
	ID          uuid.UUID  `json:"id"`
	SchoolID    uuid.UUID  `json:"school_id"`
	RequesterID uuid.UUID  `json:"requester_id"`
	TutorID     uuid.UUID  `json:"tutor_id"`
	TopicID     *uuid.UUID `json:"topic_id"`
	TopicName   string     `json:"topic_name"`
	Status      string     `json:"status"`
	Message     string     `json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TutorRequestView includes requester and tutor display names.
type TutorRequestView struct {
	TutorRequest
	RequesterName string `json:"requester_name"`
	TutorName     string `json:"tutor_name"`
}

// TutorMatch is returned by the matching endpoint.
type TutorMatch struct {
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	TopicName  string    `json:"topic_name"`
	Mastery    float64   `json:"mastery"`
	Sessions   int       `json:"sessions"`
	Reputation int       `json:"reputation"`
}

// TeachingTopic is a topic the user qualifies to teach (mastery at or above the tutoring threshold).
type TeachingTopic struct {
	TopicID   uuid.UUID `json:"topic_id"`
	TopicName string    `json:"topic_name"`
	Mastery   float64   `json:"mastery"`
}

type TutoringRepository interface {
	CreateRequest(ctx context.Context, schoolID, requesterID, tutorID uuid.UUID, topicName, message string, topicID *uuid.UUID) (*TutorRequest, error)
	ListRequestsForTutor(ctx context.Context, tutorID, schoolID uuid.UUID, status string) ([]TutorRequestView, error)
	ListRequestsByRequester(ctx context.Context, requesterID, schoolID uuid.UUID, status string) ([]TutorRequestView, error)
	UpdateRequestStatus(ctx context.Context, requestID, userID, schoolID uuid.UUID, status string) (*TutorRequest, error)
	FindTutorMatches(ctx context.Context, schoolID, requesterID uuid.UUID, topicName string, limit int) ([]TutorMatch, error)
	ListTeachingTopics(ctx context.Context, userID, schoolID uuid.UUID) ([]TeachingTopic, error)
}

type TutoringService interface {
	CreateRequest(ctx context.Context, schoolID, requesterID, tutorID uuid.UUID, topicName, message string, topicID *uuid.UUID) (*TutorRequest, error)
	ListRequestsForTutor(ctx context.Context, tutorID, schoolID uuid.UUID, status string) ([]TutorRequestView, error)
	ListRequestsByRequester(ctx context.Context, requesterID, schoolID uuid.UUID, status string) ([]TutorRequestView, error)
	UpdateRequestStatus(ctx context.Context, requestID, userID, schoolID uuid.UUID, status string) (*TutorRequest, error)
	FindTutorMatches(ctx context.Context, schoolID, requesterID uuid.UUID, topicName string, limit int) ([]TutorMatch, error)
	ListTeachingTopics(ctx context.Context, userID, schoolID uuid.UUID) ([]TeachingTopic, error)
}
