package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

// PlannerService defines the planner business operations.
type PlannerService interface {
	ListStudySessions(ctx context.Context, userID, schoolID uuid.UUID, startDate, endDate string) ([]*domain.StudySession, error)
	CreateStudySession(ctx context.Context, userID, schoolID uuid.UUID, title, scheduledDate, startTime string, durationMinutes int) (*domain.StudySession, error)
	UpdateStudySessionStatus(ctx context.Context, id, userID, schoolID uuid.UUID, status string) error
	DeleteStudySession(ctx context.Context, id, userID, schoolID uuid.UUID) error
	MarkMissedYesterday(ctx context.Context, userID, schoolID uuid.UUID) (int, error)
	ListUpcomingDeadlines(ctx context.Context, userID, schoolID uuid.UUID, limit int) ([]*domain.UpcomingDeadline, error)
	CreateStudyDeadline(ctx context.Context, userID, schoolID uuid.UUID, name, courseName, dueDate string) (*domain.StudyDeadline, error)
	DeleteStudyDeadline(ctx context.Context, id, userID, schoolID uuid.UUID) error
	RegeneratePlan(ctx context.Context, userID, schoolID uuid.UUID) ([]*domain.StudySession, error)
}

type plannerServiceImpl struct {
	repo domain.PlannerRepository
}

// NewPlannerService returns a PlannerService backed by the given repository.
func NewPlannerService(repo domain.PlannerRepository) PlannerService {
	return &plannerServiceImpl{repo: repo}
}

func (s *plannerServiceImpl) ListStudySessions(ctx context.Context, userID, schoolID uuid.UUID, startDate, endDate string) ([]*domain.StudySession, error) {
	return s.repo.ListStudySessions(ctx, userID, schoolID, startDate, endDate)
}

func (s *plannerServiceImpl) CreateStudySession(ctx context.Context, userID, schoolID uuid.UUID, title, scheduledDate, startTime string, durationMinutes int) (*domain.StudySession, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if durationMinutes < 5 {
		durationMinutes = 30
	}
	if startTime == "" {
		startTime = "09:00"
	}
	sd, err := time.Parse("2006-01-02", scheduledDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	session := &domain.StudySession{
		UserID:          userID,
		SchoolID:        schoolID,
		Title:           title,
		ScheduledDate:   sd,
		StartTime:       startTime,
		DurationMinutes: durationMinutes,
		Status:          "planned",
	}
	return s.repo.CreateStudySession(ctx, session)
}

func (s *plannerServiceImpl) UpdateStudySessionStatus(ctx context.Context, id, userID, schoolID uuid.UUID, status string) error {
	switch status {
	case "done", "missed", "planned":
		// valid
	default:
		return fmt.Errorf("invalid status: %s", status)
	}
	return s.repo.UpdateStudySessionStatus(ctx, id, userID, schoolID, status)
}

func (s *plannerServiceImpl) DeleteStudySession(ctx context.Context, id, userID, schoolID uuid.UUID) error {
	return s.repo.DeleteStudySession(ctx, id, userID, schoolID)
}

func (s *plannerServiceImpl) MarkMissedYesterday(ctx context.Context, userID, schoolID uuid.UUID) (int, error) {
	today := time.Now().Format("2006-01-02")
	return s.repo.MarkMissedSessions(ctx, userID, schoolID, today)
}

func (s *plannerServiceImpl) ListUpcomingDeadlines(ctx context.Context, userID, schoolID uuid.UUID, limit int) ([]*domain.UpcomingDeadline, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.ListUpcomingDeadlines(ctx, userID, schoolID, limit)
}

func (s *plannerServiceImpl) CreateStudyDeadline(ctx context.Context, userID, schoolID uuid.UUID, name, courseName, dueDate string) (*domain.StudyDeadline, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	dd, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	deadline := &domain.StudyDeadline{
		UserID:     userID,
		SchoolID:   schoolID,
		Name:       name,
		CourseName: courseName,
		DueDate:    dd,
		Source:     "manual",
	}
	return s.repo.CreateStudyDeadline(ctx, deadline)
}

func (s *plannerServiceImpl) DeleteStudyDeadline(ctx context.Context, id, userID, schoolID uuid.UUID) error {
	return s.repo.DeleteStudyDeadline(ctx, id, userID, schoolID)
}

// RegeneratePlan creates a simple algorithmic study plan.
// It looks at upcoming deadlines and distributes study sessions across the days before each one.
func (s *plannerServiceImpl) RegeneratePlan(ctx context.Context, userID, schoolID uuid.UUID) ([]*domain.StudySession, error) {
	deadlines, err := s.repo.ListUpcomingDeadlines(ctx, userID, schoolID, 20)
	if err != nil {
		return nil, fmt.Errorf("RegeneratePlan: list deadlines: %w", err)
	}
	if len(deadlines) == 0 {
		return nil, nil
	}

	// Delete existing planned sessions from today onward
	today := time.Now().UTC().Truncate(24 * time.Hour)
	endDate := today.Add(28 * 24 * time.Hour) // 4 weeks ahead
	existing, _ := s.repo.ListStudySessions(ctx, userID, schoolID, today.Format("2006-01-02"), endDate.Format("2006-01-02"))
	for _, sess := range existing {
		if sess.Status == "planned" {
			_ = s.repo.DeleteStudySession(ctx, sess.ID, userID, schoolID)
		}
	}

	// Available time slots: 9am and 2pm daily
	type slot struct {
		hour int
		used map[string]bool
	}
	slots := []slot{
		{hour: 9, used: make(map[string]bool)},
		{hour: 14, used: make(map[string]bool)},
	}

	var created []*domain.StudySession

	for _, dl := range deadlines {
		daysUntil := dl.DaysUntil
		if daysUntil < 1 {
			continue
		}
		// Spread 2-4 sessions before each deadline
		numSessions := 2
		if daysUntil >= 7 {
			numSessions = 3
		}
		if daysUntil >= 14 {
			numSessions = 4
		}

		spacing := daysUntil / (numSessions + 1)
		if spacing < 1 {
			spacing = 1
		}

		title := dl.Name
		if dl.CourseName != "" {
			title = dl.CourseName + ": " + dl.Name
		}

		for i := 1; i <= numSessions; i++ {
			dayOffset := i * spacing
			if dayOffset >= daysUntil {
				dayOffset = daysUntil - 1
			}
			if dayOffset < 0 {
				dayOffset = 0
			}
			sDate := today.Add(time.Duration(dayOffset) * 24 * time.Hour)
			dateStr := sDate.Format("2006-01-02")

			// Find available slot
			placed := false
			for si := range slots {
				if !slots[si].used[dateStr] {
					slots[si].used[dateStr] = true
					sess := &domain.StudySession{
						UserID:          userID,
						SchoolID:        schoolID,
						Title:           title,
						ScheduledDate:   sDate,
						StartTime:       fmt.Sprintf("%02d:00", slots[si].hour),
						DurationMinutes: 45,
						Status:          "planned",
					}
					out, err := s.repo.CreateStudySession(ctx, sess)
					if err == nil {
						created = append(created, out)
					}
					placed = true
					break
				}
			}
			if !placed {
				// Both slots used, skip this day
				continue
			}
		}
	}

	return created, nil
}
