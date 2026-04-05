package repository

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

// PlannerRepo implements domain.PlannerRepository using pgx.
type PlannerRepo struct {
	DB *pgxpool.Pool
}

func (r *PlannerRepo) ListStudySessions(ctx context.Context, userID, schoolID uuid.UUID, startDate, endDate string) ([]*domain.StudySession, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, title, course_id, topic_id, scheduled_date, start_time, duration_minutes, status, created_at, updated_at
		FROM study_sessions
		WHERE user_id = $1 AND school_id = $2
		  AND scheduled_date >= $3::date AND scheduled_date <= $4::date
		ORDER BY scheduled_date, start_time`,
		userID, schoolID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("ListStudySessions: %w", err)
	}
	defer rows.Close()

	var sessions []*domain.StudySession
	for rows.Next() {
		s := &domain.StudySession{UserID: userID, SchoolID: schoolID}
		var startTime time.Time
		if err := rows.Scan(&s.ID, &s.Title, &s.CourseID, &s.TopicID, &s.ScheduledDate, &startTime, &s.DurationMinutes, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("ListStudySessions scan: %w", err)
		}
		s.StartTime = startTime.Format("15:04")
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *PlannerRepo) CreateStudySession(ctx context.Context, session *domain.StudySession) (*domain.StudySession, error) {
	err := r.DB.QueryRow(ctx, `
		INSERT INTO study_sessions (user_id, school_id, title, course_id, topic_id, scheduled_date, start_time, duration_minutes, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7::time, $8, $9)
		RETURNING id, created_at, updated_at`,
		session.UserID, session.SchoolID, session.Title, session.CourseID, session.TopicID,
		session.ScheduledDate.Format("2006-01-02"), session.StartTime, session.DurationMinutes, session.Status,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateStudySession: %w", err)
	}
	return session, nil
}

func (r *PlannerRepo) UpdateStudySessionStatus(ctx context.Context, id, userID, schoolID uuid.UUID, status string) error {
	tag, err := r.DB.Exec(ctx, `
		UPDATE study_sessions SET status = $1, updated_at = now()
		WHERE id = $2 AND user_id = $3 AND school_id = $4`,
		status, id, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("UpdateStudySessionStatus: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PlannerRepo) DeleteStudySession(ctx context.Context, id, userID, schoolID uuid.UUID) error {
	tag, err := r.DB.Exec(ctx, `
		DELETE FROM study_sessions WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		id, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteStudySession: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PlannerRepo) MarkMissedSessions(ctx context.Context, userID, schoolID uuid.UUID, before string) (int, error) {
	tag, err := r.DB.Exec(ctx, `
		UPDATE study_sessions SET status = 'missed', updated_at = now()
		WHERE user_id = $1 AND school_id = $2
		  AND scheduled_date < $3::date
		  AND status = 'planned'`,
		userID, schoolID, before,
	)
	if err != nil {
		return 0, fmt.Errorf("MarkMissedSessions: %w", err)
	}
	return int(tag.RowsAffected()), nil
}

func (r *PlannerRepo) ListUpcomingDeadlines(ctx context.Context, userID, schoolID uuid.UUID, limit int) ([]*domain.UpcomingDeadline, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT 
			la.id::text AS id,
			la.title AS name,
			COALESCE(lc.lms_course_name, '') AS course_name,
			la.due_at AS due_date,
			'lms' AS source
		FROM lms_assignments la
		JOIN lms_courses lc 
		  ON lc.lms_course_id = la.lms_course_id 
		 AND lc.user_id = $1 
		 AND lc.school_id = $2
		WHERE la.school_id = $2
		  AND la.due_at IS NOT NULL
		  AND la.due_at >= now()
		ORDER BY due_date ASC
		LIMIT $3`,
		userID, schoolID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("ListUpcomingDeadlines: %w", err)
	}
	defer rows.Close()

	var deadlines []*domain.UpcomingDeadline
	for rows.Next() {
		d := &domain.UpcomingDeadline{}
		if err := rows.Scan(&d.ID, &d.Name, &d.CourseName, &d.DueDate, &d.Source); err != nil {
			return nil, fmt.Errorf("ListUpcomingDeadlines scan: %w", err)
		}
		d.DaysUntil = int(math.Ceil(time.Until(d.DueDate).Hours() / 24))
		if d.DaysUntil < 0 {
			d.DaysUntil = 0
		}
		switch {
		case d.DaysUntil <= 3:
			d.Urgency = "urgent"
		case d.DaysUntil <= 10:
			d.Urgency = "soon"
		default:
			d.Urgency = "safe"
		}
		deadlines = append(deadlines, d)
	}
	return deadlines, nil
}

func (r *PlannerRepo) CreateStudyDeadline(ctx context.Context, deadline *domain.StudyDeadline) (*domain.StudyDeadline, error) {
	err := r.DB.QueryRow(ctx, `
		INSERT INTO study_deadlines (user_id, school_id, name, course_name, due_date, source, lms_assignment_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`,
		deadline.UserID, deadline.SchoolID, deadline.Name, deadline.CourseName,
		deadline.DueDate.Format("2006-01-02"), deadline.Source, deadline.LMSAssignmentID,
	).Scan(&deadline.ID, &deadline.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateStudyDeadline: %w", err)
	}
	return deadline, nil
}

func (r *PlannerRepo) DeleteStudyDeadline(ctx context.Context, id, userID, schoolID uuid.UUID) error {
	tag, err := r.DB.Exec(ctx, `
		DELETE FROM study_deadlines WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		id, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteStudyDeadline: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
