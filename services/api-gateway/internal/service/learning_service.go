package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/ai"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

// LMSCourseInput is one LMS course row for import into Synapse courses.
type LMSCourseInput struct {
	LMSCourseID string `json:"lms_course_id"`
	Name        string `json:"name"`
	Term        string `json:"term"`
}

// SubmitReviewResult is returned after logging a review and updating the scheduler.
type SubmitReviewResult struct {
	EaseFactor   float64
	IntervalDays int
	DueAt        time.Time
}

// LearningService defines course/topic/flashcard/review operations for the learning core.
type LearningService interface {
	CreateCourse(ctx context.Context, userID, schoolID uuid.UUID, name, term, color string) (*domain.Course, error)
	ListCourses(ctx context.Context, userID, schoolID uuid.UUID) ([]*domain.Course, error)
	GetCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) (*domain.Course, error)
	DeleteCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) error

	ImportCoursesFromLMS(ctx context.Context, userID, schoolID uuid.UUID, lmsCourses []LMSCourseInput) ([]*domain.Course, error)

	CreateTopic(ctx context.Context, userID, schoolID, courseID uuid.UUID, name string, examWeight *float64) (*domain.Topic, error)
	ListTopics(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.Topic, error)

	CreateFlashcard(ctx context.Context, userID, schoolID, courseID uuid.UUID, topicID *uuid.UUID, cardType, prompt, answer string) (*domain.Flashcard, error)
	ListFlashcards(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.Flashcard, error)
	DeleteFlashcard(ctx context.Context, id, userID, schoolID uuid.UUID) error

	GetDueCards(ctx context.Context, userID, schoolID uuid.UUID, limit int, courseID *uuid.UUID) ([]*domain.DueCard, error)
	GetConfusionInsights(ctx context.Context, userID, schoolID uuid.UUID, courseID *uuid.UUID, windowDays int) (*domain.ConfusionInsight, error)
	ListNoteMetrics(ctx context.Context, courseID, userID, schoolID uuid.UUID, windowDays int) ([]*domain.NoteMetric, error)

	SubmitReview(ctx context.Context, userID, schoolID, flashcardID, sessionID uuid.UUID, correct bool, confidence int, confused bool, responseTimeMs int) (*SubmitReviewResult, error)

	CreateNoteText(ctx context.Context, userID, schoolID, courseID uuid.UUID, title, content string, topicID *uuid.UUID) (*domain.NoteText, error)
	ListNoteTexts(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.NoteText, error)
	GetNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) (*domain.NoteText, error)
	UpdateNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID, title, content string, topicID *uuid.UUID) (*domain.NoteText, error)
	DeleteNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) error
}

type learningServiceImpl struct {
	repo domain.LearningRepository
}

func NewLearningService(repo domain.LearningRepository) LearningService {
	return &learningServiceImpl{repo: repo}
}

func (s *learningServiceImpl) CreateCourse(ctx context.Context, userID, schoolID uuid.UUID, name, term, color string) (*domain.Course, error) {
	if name == "" {
		return nil, &domain.ValidationError{Message: "name is required"}
	}
	c := &domain.Course{
		SchoolID: schoolID,
		UserID:   userID,
		Name:     name,
		Term:     term,
		Color:    color,
	}
	out, err := s.repo.CreateCourse(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("CreateCourse: %w", err)
	}
	return out, nil
}

func (s *learningServiceImpl) ListCourses(ctx context.Context, userID, schoolID uuid.UUID) ([]*domain.Course, error) {
	return s.repo.ListCourses(ctx, userID, schoolID)
}

func (s *learningServiceImpl) GetCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) (*domain.Course, error) {
	c, err := s.repo.GetCourse(ctx, courseID, userID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("GetCourse: %w", err)
	}
	return c, nil
}

func (s *learningServiceImpl) DeleteCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) error {
	if err := s.repo.DeleteCourse(ctx, courseID, userID, schoolID); err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}
	return nil
}

func (s *learningServiceImpl) ImportCoursesFromLMS(ctx context.Context, userID, schoolID uuid.UUID, lmsCourses []LMSCourseInput) ([]*domain.Course, error) {
	existing, err := s.repo.ListCourses(ctx, userID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("ImportCoursesFromLMS: %w", err)
	}
	seen := make(map[string]struct{})
	for _, c := range existing {
		if c.LMSCourseID != nil && *c.LMSCourseID != "" {
			seen[*c.LMSCourseID] = struct{}{}
		}
	}

	var created []*domain.Course
	for _, in := range lmsCourses {
		if in.LMSCourseID == "" {
			continue
		}
		if _, dup := seen[in.LMSCourseID]; dup {
			continue
		}
		if in.Name == "" {
			continue
		}
		lmsID := in.LMSCourseID
		c := &domain.Course{
			SchoolID:    schoolID,
			UserID:      userID,
			Name:        in.Name,
			Term:        in.Term,
			LMSCourseID: &lmsID,
		}
		out, err := s.repo.CreateCourse(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("ImportCoursesFromLMS: %w", err)
		}
		seen[in.LMSCourseID] = struct{}{}
		created = append(created, out)
	}
	return created, nil
}

func (s *learningServiceImpl) CreateTopic(ctx context.Context, userID, schoolID, courseID uuid.UUID, name string, examWeight *float64) (*domain.Topic, error) {
	if name == "" {
		return nil, &domain.ValidationError{Message: "name is required"}
	}
	if _, err := s.repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("CreateTopic: %w", err)
	}
	t := &domain.Topic{
		SchoolID:   schoolID,
		CourseID:   courseID,
		UserID:     userID,
		Name:       name,
		ExamWeight: examWeight,
		Source:     "manual",
	}
	out, err := s.repo.CreateTopic(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("CreateTopic: %w", err)
	}
	return out, nil
}

func (s *learningServiceImpl) ListTopics(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.Topic, error) {
	if _, err := s.repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("ListTopics: %w", err)
	}
	return s.repo.ListTopics(ctx, courseID, userID, schoolID)
}

func (s *learningServiceImpl) CreateFlashcard(ctx context.Context, userID, schoolID, courseID uuid.UUID, topicID *uuid.UUID, cardType, prompt, answer string) (*domain.Flashcard, error) {
	if prompt == "" || answer == "" {
		return nil, &domain.ValidationError{Message: "prompt and answer are required"}
	}
	if cardType == "" {
		cardType = "qa"
	}
	if _, err := s.repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("CreateFlashcard: %w", err)
	}
	if topicID != nil {
		topics, err := s.repo.ListTopics(ctx, courseID, userID, schoolID)
		if err != nil {
			return nil, fmt.Errorf("CreateFlashcard: %w", err)
		}
		ok := false
		for _, t := range topics {
			if t.ID == *topicID {
				ok = true
				break
			}
		}
		if !ok {
			return nil, &domain.ValidationError{Message: "topic not found in this course"}
		}
	}

	fc := &domain.Flashcard{
		SchoolID:   schoolID,
		CourseID:   courseID,
		UserID:     userID,
		TopicID:    topicID,
		CardType:   cardType,
		Prompt:     prompt,
		Answer:     answer,
		CreatedBy:  "user",
		Visibility: "private",
	}
	out, err := s.repo.CreateFlashcard(ctx, fc)
	if err != nil {
		return nil, fmt.Errorf("CreateFlashcard: %w", err)
	}

	now := time.Now()
	state := &domain.SchedulerState{
		FlashcardID:  out.ID,
		UserID:       userID,
		SchoolID:     schoolID,
		EaseFactor:   2.5,
		IntervalDays: 1,
		DueAt:        now,
		LapseCount:   0,
		LastReviewAt: nil,
	}
	if err := s.repo.UpsertSchedulerState(ctx, state); err != nil {
		return nil, fmt.Errorf("CreateFlashcard: scheduler: %w", err)
	}
	return out, nil
}

func (s *learningServiceImpl) ListFlashcards(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.Flashcard, error) {
	if _, err := s.repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("ListFlashcards: %w", err)
	}
	return s.repo.ListFlashcards(ctx, courseID, userID, schoolID)
}

func (s *learningServiceImpl) DeleteFlashcard(ctx context.Context, id, userID, schoolID uuid.UUID) error {
	if err := s.repo.DeleteFlashcard(ctx, id, userID, schoolID); err != nil {
		return fmt.Errorf("DeleteFlashcard: %w", err)
	}
	return nil
}

func (s *learningServiceImpl) GetDueCards(ctx context.Context, userID, schoolID uuid.UUID, limit int, courseID *uuid.UUID) ([]*domain.DueCard, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	return s.repo.GetDueCards(ctx, userID, schoolID, limit, courseID)
}

func (s *learningServiceImpl) GetConfusionInsights(ctx context.Context, userID, schoolID uuid.UUID, courseID *uuid.UUID, windowDays int) (*domain.ConfusionInsight, error) {
	return s.repo.GetConfusionInsights(ctx, userID, schoolID, courseID, windowDays)
}

func (s *learningServiceImpl) ListNoteMetrics(ctx context.Context, courseID, userID, schoolID uuid.UUID, windowDays int) ([]*domain.NoteMetric, error) {
	if _, err := s.repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("ListNoteMetrics: %w", err)
	}
	return s.repo.ListNoteMetrics(ctx, courseID, userID, schoolID, windowDays)
}

func (s *learningServiceImpl) SubmitReview(ctx context.Context, userID, schoolID, flashcardID, sessionID uuid.UUID, correct bool, confidence int, confused bool, responseTimeMs int) (*SubmitReviewResult, error) {
	if confidence < 1 || confidence > 4 {
		return nil, &domain.ValidationError{Message: "confidence must be between 1 and 4"}
	}

	if _, err := s.repo.GetFlashcard(ctx, flashcardID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("SubmitReview: %w", err)
	}

	state, err := s.repo.GetSchedulerState(ctx, flashcardID, userID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("SubmitReview: %w", err)
	}

	easeBefore := state.EaseFactor
	intervalBefore := state.IntervalDays
	update := domain.ComputeSM2(domain.ReviewInput{
		Correct:         correct,
		Confidence:      confidence,
		Confused:        confused,
		ResponseTimeMs:  responseTimeMs,
		CurrentEase:     state.EaseFactor,
		CurrentInterval: state.IntervalDays,
		LapseCount:      state.LapseCount,
	})

	now := time.Now()
	lastAt := now
	newState := &domain.SchedulerState{
		FlashcardID:  flashcardID,
		UserID:       userID,
		SchoolID:     schoolID,
		EaseFactor:   update.NewEase,
		IntervalDays: update.NewInterval,
		DueAt:        update.NewDueAt,
		LapseCount:   update.NewLapses,
		LastReviewAt: &lastAt,
	}
	if err := s.repo.UpsertSchedulerState(ctx, newState); err != nil {
		return nil, fmt.Errorf("SubmitReview: %w", err)
	}

	eb := easeBefore
	ib := intervalBefore
	ev := &domain.ReviewEvent{
		ID:             uuid.New(),
		SchoolID:       schoolID,
		UserID:         userID,
		FlashcardID:    flashcardID,
		SessionID:      sessionID,
		Ts:             now,
		Correct:        correct,
		Confidence:     confidence,
		Confused:       confused,
		ResponseTimeMs: responseTimeMs,
		EaseBefore:     &eb,
		IntervalBefore: &ib,
	}
	if err := s.repo.CreateReviewEvent(ctx, ev); err != nil {
		return nil, fmt.Errorf("SubmitReview: %w", err)
	}

	return &SubmitReviewResult{
		EaseFactor:   update.NewEase,
		IntervalDays: update.NewInterval,
		DueAt:        update.NewDueAt,
	}, nil
}

func (s *learningServiceImpl) CreateNoteText(ctx context.Context, userID, schoolID, courseID uuid.UUID, title, content string, topicID *uuid.UUID) (*domain.NoteText, error) {
	if content == "" {
		return nil, &domain.ValidationError{Message: "content is required"}
	}
	if _, err := s.repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("CreateNoteText: %w", err)
	}
	if topicID != nil {
		topics, err := s.repo.ListTopics(ctx, courseID, userID, schoolID)
		if err != nil {
			return nil, fmt.Errorf("CreateNoteText: %w", err)
		}
		ok := false
		for _, t := range topics {
			if t.ID == *topicID {
				ok = true
				break
			}
		}
		if !ok {
			return nil, &domain.ValidationError{Message: "topic not found in this course"}
		}
	}
	n := &domain.NoteText{
		SchoolID: schoolID,
		CourseID: courseID,
		UserID:   userID,
		TopicID:  topicID,
		Title:    title,
		Content:  content,
	}
	out, err := s.repo.CreateNoteText(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("CreateNoteText: %w", err)
	}
	return out, nil
}

func (s *learningServiceImpl) ListNoteTexts(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.NoteText, error) {
	if _, err := s.repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("ListNoteTexts: %w", err)
	}
	return s.repo.ListNoteTexts(ctx, courseID, userID, schoolID)
}

func (s *learningServiceImpl) GetNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) (*domain.NoteText, error) {
	return s.repo.GetNoteText(ctx, noteID, userID, schoolID)
}

func (s *learningServiceImpl) UpdateNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID, title, content string, topicID *uuid.UUID) (*domain.NoteText, error) {
	note, err := s.repo.GetNoteText(ctx, noteID, userID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("UpdateNoteText: %w", err)
	}
	if topicID != nil {
		topics, err := s.repo.ListTopics(ctx, note.CourseID, userID, schoolID)
		if err != nil {
			return nil, fmt.Errorf("UpdateNoteText: %w", err)
		}
		ok := false
		for _, t := range topics {
			if t.ID == *topicID {
				ok = true
				break
			}
		}
		if !ok {
			return nil, &domain.ValidationError{Message: "topic not found in this course"}
		}
	}
	return s.repo.UpdateNoteText(ctx, noteID, userID, schoolID, title, content, topicID)
}

func (s *learningServiceImpl) DeleteNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) error {
	return s.repo.DeleteNoteText(ctx, noteID, userID, schoolID)
}

// MapLearningError maps domain errors to HTTP-friendly handling in handlers.
func MapLearningError(err error) (status int, msg string) {
	var upstream *ai.UpstreamError
	if errors.As(err, &upstream) {
		clientStatus := 502
		if upstream.StatusCode == 429 {
			clientStatus = 429
		}
		// Use provider-aware messaging while avoiding raw upstream payload leakage.
		msg := upstream.ClientMessage()
		if strings.TrimSpace(msg) == "" {
			msg = "AI provider request failed"
		}
		return clientStatus, msg
	}

	var v *domain.ValidationError
	if errors.As(err, &v) {
		return 422, v.Message
	}
	if errors.Is(err, domain.ErrNotFound) {
		return 404, "not found"
	}
	return 500, "internal error"
}
