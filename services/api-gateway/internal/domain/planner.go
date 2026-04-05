package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// StudySession mirrors the study_sessions table.
type StudySession struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"-"`
	SchoolID        uuid.UUID  `json:"-"`
	Title           string     `json:"title"`
	CourseID        *uuid.UUID `json:"course_id,omitempty"`
	TopicID         *uuid.UUID `json:"topic_id,omitempty"`
	ScheduledDate   time.Time  `json:"scheduled_date"`
	StartTime       string     `json:"start_time"`
	DurationMinutes int        `json:"duration_minutes"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// StudyDeadline mirrors the study_deadlines table.
type StudyDeadline struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"-"`
	SchoolID        uuid.UUID `json:"-"`
	Name            string    `json:"name"`
	CourseName      string    `json:"course_name,omitempty"`
	DueDate         time.Time `json:"due_date"`
	Source          string    `json:"source"`
	LMSAssignmentID string   `json:"lms_assignment_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// UpcomingDeadline is a merged view of LMS deadlines + user-created ones.
type UpcomingDeadline struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CourseName string    `json:"course_name"`
	DueDate    time.Time `json:"due_date"`
	DaysUntil  int       `json:"days_until"`
	Source     string    `json:"source"` // "lms" or "manual"
	Urgency    string    `json:"urgency"` // "urgent", "soon", "safe"
}

// PlannerRepository abstracts persistence for the planner feature.
type PlannerRepository interface {
	ListStudySessions(ctx context.Context, userID, schoolID uuid.UUID, startDate, endDate string) ([]*StudySession, error)
	CreateStudySession(ctx context.Context, session *StudySession) (*StudySession, error)
	UpdateStudySessionStatus(ctx context.Context, id, userID, schoolID uuid.UUID, status string) error
	DeleteStudySession(ctx context.Context, id, userID, schoolID uuid.UUID) error
	MarkMissedSessions(ctx context.Context, userID, schoolID uuid.UUID, before string) (int, error)

	ListUpcomingDeadlines(ctx context.Context, userID, schoolID uuid.UUID, limit int) ([]*UpcomingDeadline, error)
	CreateStudyDeadline(ctx context.Context, deadline *StudyDeadline) (*StudyDeadline, error)
	DeleteStudyDeadline(ctx context.Context, id, userID, schoolID uuid.UUID) error
}
