package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/lms-service/internal/domain"
)

type PostgresLMSRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresLMSRepo(pool *pgxpool.Pool) *PostgresLMSRepo {
	return &PostgresLMSRepo{pool: pool}
}

func (r *PostgresLMSRepo) UpsertConnection(ctx context.Context, conn *domain.LMSConnection) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_connections (user_id, school_id, lms_type, institution_url, access_token, refresh_token, token_expires_at, sync_status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (user_id, lms_type) DO UPDATE SET
		   institution_url  = EXCLUDED.institution_url,
		   access_token     = EXCLUDED.access_token,
		   refresh_token    = EXCLUDED.refresh_token,
		   token_expires_at = EXCLUDED.token_expires_at,
		   sync_status      = EXCLUDED.sync_status`,
		conn.UserID, conn.SchoolID, conn.LMSType, conn.InstitutionURL,
		conn.AccessToken, conn.RefreshToken, conn.TokenExpiresAt, conn.SyncStatus,
	)
	if err != nil {
		return fmt.Errorf("UpsertConnection: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) FindConnectionByUser(ctx context.Context, userID, schoolID uuid.UUID) (*domain.LMSConnection, error) {
	var c domain.LMSConnection
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, school_id, lms_type, institution_url,
		        access_token, refresh_token, token_expires_at,
		        last_synced_at, sync_status, created_at
		 FROM lms_connections
		 WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	).Scan(
		&c.ID, &c.UserID, &c.SchoolID, &c.LMSType, &c.InstitutionURL,
		&c.AccessToken, &c.RefreshToken, &c.TokenExpiresAt,
		&c.LastSyncedAt, &c.SyncStatus, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("FindConnectionByUser: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("FindConnectionByUser: %w", err)
	}
	return &c, nil
}

func (r *PostgresLMSRepo) DeleteConnection(ctx context.Context, userID, schoolID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM lms_connections WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteConnection: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) UpsertCourse(ctx context.Context, course *domain.LMSCourse) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_courses (user_id, school_id, lms_course_id, lms_course_name, lms_term, enrollment_type, last_synced_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (user_id, lms_course_id) DO UPDATE SET
		   lms_course_name = EXCLUDED.lms_course_name,
		   lms_term = EXCLUDED.lms_term,
		   last_synced_at = EXCLUDED.last_synced_at`,
		course.UserID, course.SchoolID, course.LMSCourseID, course.LMSCourseName,
		course.LMSTerm, course.EnrollmentType, course.LastSyncedAt,
	)
	if err != nil {
		return fmt.Errorf("UpsertCourse: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) UpsertAssignment(ctx context.Context, assignment *domain.LMSAssignment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_assignments
		 (school_id, lms_assignment_id, lms_course_id, title, due_at, points_possible, assignment_group, last_synced_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (lms_assignment_id, school_id) DO UPDATE SET
		   title = EXCLUDED.title,
		   due_at = EXCLUDED.due_at,
		   points_possible = EXCLUDED.points_possible,
		   last_synced_at = EXCLUDED.last_synced_at`,
		assignment.SchoolID, assignment.LMSAssignmentID, assignment.LMSCourseID, assignment.Title,
		assignment.DueAt, assignment.PointsPossible, assignment.AssignmentGroup, assignment.LastSyncedAt,
	)
	if err != nil {
		return fmt.Errorf("UpsertAssignment: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) UpsertGradeEvent(ctx context.Context, event *domain.LMSGradeEvent) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_grade_events
		 (user_id, school_id, lms_assignment_id, lms_course_id, score, points_possible,
		  submitted_at, graded_at, grade_type, synced_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (user_id, lms_assignment_id) DO UPDATE SET
		   score = EXCLUDED.score,
		   points_possible = EXCLUDED.points_possible,
		   graded_at = EXCLUDED.graded_at,
		   synced_at = EXCLUDED.synced_at`,
		event.UserID, event.SchoolID, event.LMSAssignmentID, event.LMSCourseID,
		event.Score, event.PointsPossible, event.SubmittedAt, event.GradedAt,
		event.GradeType, event.SyncedAt,
	)
	if err != nil {
		return fmt.Errorf("UpsertGradeEvent: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) UpsertAnnouncement(ctx context.Context, item *domain.LMSAnnouncement) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_announcements
		 (user_id, school_id, lms_course_id, lms_announcement_id, title, message, posted_at, html_url, last_synced_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 ON CONFLICT (user_id, lms_announcement_id) DO UPDATE SET
		   title = EXCLUDED.title,
		   message = EXCLUDED.message,
		   posted_at = EXCLUDED.posted_at,
		   html_url = EXCLUDED.html_url,
		   last_synced_at = EXCLUDED.last_synced_at`,
		item.UserID, item.SchoolID, item.LMSCourseID, item.LMSAnnouncementID,
		item.Title, item.Message, item.PostedAt, item.HTMLURL, item.LastSyncedAt,
	)
	if err != nil {
		return fmt.Errorf("UpsertAnnouncement: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) UpsertDiscussionTopic(ctx context.Context, item *domain.LMSDiscussionTopic) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_discussion_topics
		 (user_id, school_id, lms_course_id, lms_topic_id, title, message, posted_at, html_url, last_synced_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 ON CONFLICT (user_id, lms_topic_id) DO UPDATE SET
		   title = EXCLUDED.title,
		   message = EXCLUDED.message,
		   posted_at = EXCLUDED.posted_at,
		   html_url = EXCLUDED.html_url,
		   last_synced_at = EXCLUDED.last_synced_at`,
		item.UserID, item.SchoolID, item.LMSCourseID, item.LMSTopicID,
		item.Title, item.Message, item.PostedAt, item.HTMLURL, item.LastSyncedAt,
	)
	if err != nil {
		return fmt.Errorf("UpsertDiscussionTopic: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) UpsertSubmissionState(ctx context.Context, item *domain.LMSSubmissionState) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO lms_submission_states
		 (user_id, school_id, lms_assignment_id, lms_course_id, workflow_state, missing, late, excused,
		  submitted_at, graded_at, score, points_possible, synced_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 ON CONFLICT (user_id, lms_assignment_id) DO UPDATE SET
		   workflow_state = EXCLUDED.workflow_state,
		   missing = EXCLUDED.missing,
		   late = EXCLUDED.late,
		   excused = EXCLUDED.excused,
		   submitted_at = EXCLUDED.submitted_at,
		   graded_at = EXCLUDED.graded_at,
		   score = EXCLUDED.score,
		   points_possible = EXCLUDED.points_possible,
		   synced_at = EXCLUDED.synced_at`,
		item.UserID, item.SchoolID, item.LMSAssignmentID, item.LMSCourseID,
		item.WorkflowState, item.Missing, item.Late, item.Excused,
		item.SubmittedAt, item.GradedAt, item.Score, item.PointsPossible, item.SyncedAt,
	)
	if err != nil {
		return fmt.Errorf("UpsertSubmissionState: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) ListCoursesByUser(ctx context.Context, userID, schoolID uuid.UUID) ([]*domain.LMSCourse, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT lms_course_id, lms_course_name, lms_term, last_synced_at
		 FROM lms_courses WHERE user_id = $1 AND school_id = $2
		 ORDER BY lms_course_name ASC`,
		userID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListCoursesByUser: %w", err)
	}
	defer rows.Close()

	var out []*domain.LMSCourse
	for rows.Next() {
		var c domain.LMSCourse
		c.UserID = userID
		c.SchoolID = schoolID
		var term *string
		if err := rows.Scan(&c.LMSCourseID, &c.LMSCourseName, &term, &c.LastSyncedAt); err != nil {
			return nil, fmt.Errorf("ListCoursesByUser: scan: %w", err)
		}
		if term != nil {
			c.LMSTerm = *term
		}
		out = append(out, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListCoursesByUser: %w", err)
	}
	return out, nil
}

func (r *PostgresLMSRepo) ListConnectionsForSync(ctx context.Context) ([]*domain.LMSConnection, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, school_id, lms_type, institution_url,
		        access_token, refresh_token, token_expires_at,
		        last_synced_at, sync_status, created_at
		 FROM lms_connections
		 WHERE lms_type = 'canvas'
		 ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListConnectionsForSync: %w", err)
	}
	defer rows.Close()

	var out []*domain.LMSConnection
	for rows.Next() {
		var c domain.LMSConnection
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.SchoolID, &c.LMSType, &c.InstitutionURL,
			&c.AccessToken, &c.RefreshToken, &c.TokenExpiresAt,
			&c.LastSyncedAt, &c.SyncStatus, &c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListConnectionsForSync: scan: %w", err)
		}
		out = append(out, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListConnectionsForSync: %w", err)
	}
	return out, nil
}

func (r *PostgresLMSRepo) UpdateConnectionSyncStatus(ctx context.Context, userID, schoolID uuid.UUID, status string, lastSyncedAt *time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE lms_connections
		 SET sync_status = $3, last_synced_at = $4
		 WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID, status, lastSyncedAt,
	)
	if err != nil {
		return fmt.Errorf("UpdateConnectionSyncStatus: %w", err)
	}
	return nil
}

func (r *PostgresLMSRepo) StartSyncRun(ctx context.Context, userID, schoolID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.pool.QueryRow(ctx,
		`INSERT INTO lms_sync_runs (user_id, school_id, status)
		 VALUES ($1, $2, 'running')
		 RETURNING id`,
		userID, schoolID,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("StartSyncRun: %w", err)
	}
	return id, nil
}

func (r *PostgresLMSRepo) FinishSyncRun(ctx context.Context, runID uuid.UUID, status string, coursesSynced int, errMsg *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE lms_sync_runs
		 SET status = $2, courses_synced = $3, error_message = $4, finished_at = now()
		 WHERE id = $1`,
		runID, status, coursesSynced, errMsg,
	)
	if err != nil {
		return fmt.Errorf("FinishSyncRun: %w", err)
	}
	return nil
}
