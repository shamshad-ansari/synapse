package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LMSConnection represents a user's connection to a Canvas LMS instance.
type LMSConnection struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SchoolID       uuid.UUID
	LMSType        string
	InstitutionURL string
	AccessToken    string
	RefreshToken   string
	TokenExpiresAt time.Time
	LastSyncedAt   *time.Time
	SyncStatus     string
	CreatedAt      time.Time
}

// LMSConnectionResponse is the public-facing DTO — never exposes raw tokens.
type LMSConnectionResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	SchoolID       uuid.UUID  `json:"school_id"`
	LMSType        string     `json:"lms_type"`
	InstitutionURL string     `json:"institution_url"`
	TokenExpiresAt time.Time  `json:"token_expires_at"`
	LastSyncedAt   *time.Time `json:"last_synced_at"`
	SyncStatus     string     `json:"sync_status"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ToResponse strips sensitive token fields from LMSConnection.
func (c *LMSConnection) ToResponse() *LMSConnectionResponse {
	return &LMSConnectionResponse{
		ID:             c.ID,
		UserID:         c.UserID,
		SchoolID:       c.SchoolID,
		LMSType:        c.LMSType,
		InstitutionURL: c.InstitutionURL,
		TokenExpiresAt: c.TokenExpiresAt,
		LastSyncedAt:   c.LastSyncedAt,
		SyncStatus:     c.SyncStatus,
		CreatedAt:      c.CreatedAt,
	}
}

// LMSCourse is a synced Canvas course row for a user.
type LMSCourse struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SchoolID       uuid.UUID
	LMSCourseID    string
	LMSCourseName  string
	LMSTerm        string
	EnrollmentType string
	LastSyncedAt   *time.Time
}

// LMSAssignment is a synced Canvas assignment row (school-scoped).
type LMSAssignment struct {
	ID              uuid.UUID
	SchoolID        uuid.UUID
	LMSAssignmentID string
	LMSCourseID     string
	Title           string
	DueAt           *time.Time
	PointsPossible  *float64
	AssignmentGroup string
	LastSyncedAt    time.Time
}

// LMSGradeEvent is a synced grade/submission snapshot for an assignment.
type LMSGradeEvent struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	SchoolID        uuid.UUID
	LMSAssignmentID string
	LMSCourseID     string
	Score           *float64
	PointsPossible  *float64
	SubmittedAt     *time.Time
	GradedAt        *time.Time
	GradeType       string
	SyncedAt        time.Time
}

type LMSAnnouncement struct {
	UserID            uuid.UUID
	SchoolID          uuid.UUID
	LMSCourseID       string
	LMSAnnouncementID string
	Title             string
	Message           string
	PostedAt          *time.Time
	HTMLURL           string
	LastSyncedAt      time.Time
}

type LMSDiscussionTopic struct {
	UserID      uuid.UUID
	SchoolID    uuid.UUID
	LMSCourseID string
	LMSTopicID  string
	Title       string
	Message     string
	PostedAt    *time.Time
	HTMLURL     string
	LastSyncedAt time.Time
}

type LMSSubmissionState struct {
	UserID          uuid.UUID
	SchoolID        uuid.UUID
	LMSAssignmentID string
	LMSCourseID     string
	WorkflowState   string
	Missing         bool
	Late            bool
	Excused         bool
	SubmittedAt     *time.Time
	GradedAt        *time.Time
	Score           *float64
	PointsPossible  *float64
	SyncedAt        time.Time
}

// LMSCourseResponse is returned by GET /v1/lms/courses.
type LMSCourseResponse struct {
	LMSCourseID  string     `json:"lms_course_id"`
	CourseName   string     `json:"course_name"`
	Term         string     `json:"term"`
	LastSyncedAt *time.Time `json:"last_synced_at"`
}

// LMSRepository abstracts persistence for LMS connection data.
type LMSRepository interface {
	UpsertConnection(ctx context.Context, conn *LMSConnection) error
	FindConnectionByUser(ctx context.Context, userID, schoolID uuid.UUID) (*LMSConnection, error)
	DeleteConnection(ctx context.Context, userID, schoolID uuid.UUID) error

	UpsertCourse(ctx context.Context, course *LMSCourse) error
	UpsertAssignment(ctx context.Context, assignment *LMSAssignment) error
	UpsertGradeEvent(ctx context.Context, event *LMSGradeEvent) error
	UpsertAnnouncement(ctx context.Context, item *LMSAnnouncement) error
	UpsertDiscussionTopic(ctx context.Context, item *LMSDiscussionTopic) error
	UpsertSubmissionState(ctx context.Context, item *LMSSubmissionState) error
	ListCoursesByUser(ctx context.Context, userID, schoolID uuid.UUID) ([]*LMSCourse, error)
	ListConnectionsForSync(ctx context.Context) ([]*LMSConnection, error)
	UpdateConnectionSyncStatus(ctx context.Context, userID, schoolID uuid.UUID, status string, lastSyncedAt *time.Time) error
	StartSyncRun(ctx context.Context, userID, schoolID uuid.UUID) (uuid.UUID, error)
	FinishSyncRun(ctx context.Context, runID uuid.UUID, status string, coursesSynced int, errMsg *string) error
}
