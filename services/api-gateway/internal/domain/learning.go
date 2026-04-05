package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Course mirrors the courses table.
type Course struct {
	ID          uuid.UUID
	SchoolID    uuid.UUID
	UserID      uuid.UUID
	Name        string
	Term        string
	Color       string
	LMSCourseID *string
	CreatedAt   time.Time
}

// Topic mirrors the topics table.
type Topic struct {
	ID            uuid.UUID
	SchoolID      uuid.UUID
	CourseID      uuid.UUID
	UserID        uuid.UUID
	Name          string
	ParentTopicID *uuid.UUID
	ExamWeight    *float64
	Source        string
	CreatedAt     time.Time
}

// Flashcard mirrors the flashcards table (embedding omitted from domain).
type Flashcard struct {
	ID         uuid.UUID
	SchoolID   uuid.UUID
	CourseID   uuid.UUID
	UserID     uuid.UUID
	TopicID    *uuid.UUID
	CardType   string
	Prompt     string
	Answer     string
	CreatedBy  string
	Visibility string
	CreatedAt  time.Time
}

// SchedulerState mirrors scheduler_states.
type SchedulerState struct {
	FlashcardID  uuid.UUID
	UserID       uuid.UUID
	SchoolID     uuid.UUID
	EaseFactor   float64
	IntervalDays int
	DueAt        time.Time
	LapseCount   int
	LastReviewAt *time.Time
}

// ReviewEvent mirrors review_events.
type ReviewEvent struct {
	ID             uuid.UUID
	SchoolID       uuid.UUID
	UserID         uuid.UUID
	FlashcardID    uuid.UUID
	SessionID      uuid.UUID
	Ts             time.Time
	Correct        bool
	Confidence     int
	Confused       bool
	ResponseTimeMs int
	EaseBefore     *float64
	IntervalBefore *int
}

// NoteText mirrors note_texts (embedding stored in DB; slice optional in domain).
type NoteText struct {
	ID         uuid.UUID
	SchoolID   uuid.UUID
	CourseID   uuid.UUID
	UserID     uuid.UUID
	TopicID    *uuid.UUID
	Title      string
	Content    string
	Embedding  []float32
	EmbeddedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// DueCard is a row returned by the due-cards scheduler query.
type DueCard struct {
	FlashcardID  uuid.UUID
	Prompt       string
	Answer       string
	CardType     string
	TopicName    string
	EaseFactor   float64
	IntervalDays int
	DueAt        time.Time
}

// ConfusionInsight summarizes recent confusion hotspots for banner widgets.
type ConfusionInsight struct {
	CourseID       *uuid.UUID
	CourseName     string
	HotspotCount   int
	TopTopicName   string
	ConfusedEvents int
	WindowDays     int
}

// NoteMetric summarizes dynamic readiness/confusion indicators for one note.
type NoteMetric struct {
	NoteID         uuid.UUID  `json:"note_id"`
	ReadinessPct   int        `json:"readiness_pct"`
	ConnectedCards int        `json:"connected_cards"`
	ReviewCount    int        `json:"review_count"`
	ConfusedCount  int        `json:"confused_count"`
	ConfusionFlag  bool       `json:"confusion_flag"`
	LastReviewAt   *time.Time `json:"last_review_at,omitempty"`
}

// LearningRepository abstracts persistence for the learning core.
type LearningRepository interface {
	CreateCourse(ctx context.Context, course *Course) (*Course, error)
	UpsertCourseFromLMS(ctx context.Context, course *Course) (*Course, error)
	ListCourses(ctx context.Context, userID, schoolID uuid.UUID) ([]*Course, error)
	GetCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) (*Course, error)
	DeleteCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) error

	CreateTopic(ctx context.Context, topic *Topic) (*Topic, error)
	ListTopics(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*Topic, error)

	CreateFlashcard(ctx context.Context, fc *Flashcard) (*Flashcard, error)
	ListFlashcards(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*Flashcard, error)
	GetFlashcard(ctx context.Context, id, userID, schoolID uuid.UUID) (*Flashcard, error)
	DeleteFlashcard(ctx context.Context, id, userID, schoolID uuid.UUID) error
	BatchCreateFlashcards(ctx context.Context, cards []*Flashcard) ([]*Flashcard, error)

	GetDueCards(ctx context.Context, userID, schoolID uuid.UUID, limit int, courseID *uuid.UUID) ([]*DueCard, error)
	UpsertSchedulerState(ctx context.Context, state *SchedulerState) error
	GetSchedulerState(ctx context.Context, flashcardID, userID, schoolID uuid.UUID) (*SchedulerState, error)

	CreateReviewEvent(ctx context.Context, event *ReviewEvent) error

	CreateNoteText(ctx context.Context, note *NoteText) (*NoteText, error)
	UpdateNoteTextEmbedding(ctx context.Context, noteID, schoolID uuid.UUID, embedding []float32) error
	ListNoteTexts(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*NoteText, error)
	GetNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) (*NoteText, error)
	UpdateNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID, title, content string, topicID *uuid.UUID) (*NoteText, error)
	DeleteNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) error

	SearchSimilarFlashcards(ctx context.Context, userID, schoolID uuid.UUID, queryVec []float32, limit int) ([]*Flashcard, error)
	UpdateFlashcardEmbedding(ctx context.Context, flashcardID, userID, schoolID uuid.UUID, embedding []float32) error
	GetConfusionInsights(ctx context.Context, userID, schoolID uuid.UUID, courseID *uuid.UUID, windowDays int) (*ConfusionInsight, error)
	ListNoteMetrics(ctx context.Context, courseID, userID, schoolID uuid.UUID, windowDays int) ([]*NoteMetric, error)
}
