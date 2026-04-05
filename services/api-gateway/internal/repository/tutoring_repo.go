package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

type PostgresTutoringRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresTutoringRepo(pool *pgxpool.Pool) *PostgresTutoringRepo {
	return &PostgresTutoringRepo{pool: pool}
}

func (r *PostgresTutoringRepo) CreateRequest(ctx context.Context, schoolID, requesterID, tutorID uuid.UUID, topicName, message string, topicID *uuid.UUID) (*domain.TutorRequest, error) {
	var out domain.TutorRequest
	err := r.pool.QueryRow(ctx,
		`INSERT INTO tutor_requests (school_id, requester_id, tutor_id, topic_id, topic_name, message)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, school_id, requester_id, tutor_id, topic_id, topic_name, status, message, created_at, updated_at`,
		schoolID, requesterID, tutorID, topicID, topicName, message,
	).Scan(
		&out.ID, &out.SchoolID, &out.RequesterID, &out.TutorID, &out.TopicID, &out.TopicName, &out.Status, &out.Message, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateRequest: %w", err)
	}
	return &out, nil
}

func normalizeTutorRequestStatusFilter(status string) string {
	s := strings.TrimSpace(strings.ToLower(status))
	allowed := map[string]bool{
		"pending": true, "accepted": true, "declined": true, "completed": true, "cancelled": true, "all": true,
	}
	if s == "" {
		return "pending"
	}
	if !allowed[s] {
		return "pending"
	}
	return s
}

func normalizeOutgoingRequestStatusFilter(status string) string {
	s := strings.TrimSpace(strings.ToLower(status))
	allowed := map[string]bool{
		"pending": true, "accepted": true, "declined": true, "completed": true, "cancelled": true, "all": true,
	}
	if s == "" {
		return "all"
	}
	if !allowed[s] {
		return "all"
	}
	return s
}

func (r *PostgresTutoringRepo) ListRequestsForTutor(ctx context.Context, tutorID, schoolID uuid.UUID, status string) ([]domain.TutorRequestView, error) {
	status = normalizeTutorRequestStatusFilter(status)
	base := `SELECT tr.id, tr.school_id, tr.requester_id, tr.tutor_id, tr.topic_id,
			tr.topic_name, tr.status, tr.message, tr.created_at, tr.updated_at,
			req.name AS requester_name, tut.name AS tutor_name
		 FROM tutor_requests tr
		 JOIN users req ON req.id = tr.requester_id
		 JOIN users tut ON tut.id = tr.tutor_id
		 WHERE tr.tutor_id = $1 AND tr.school_id = $2`
	var (
		query string
		args  []any
	)
	if status == "all" {
		query = base + ` ORDER BY tr.created_at DESC`
		args = []any{tutorID, schoolID}
	} else {
		query = base + ` AND tr.status = $3 ORDER BY tr.created_at DESC`
		args = []any{tutorID, schoolID, status}
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListRequestsForTutor: %w", err)
	}
	defer rows.Close()

	var list []domain.TutorRequestView
	for rows.Next() {
		var v domain.TutorRequestView
		if err := rows.Scan(
			&v.ID, &v.SchoolID, &v.RequesterID, &v.TutorID, &v.TopicID,
			&v.TopicName, &v.Status, &v.Message, &v.CreatedAt, &v.UpdatedAt,
			&v.RequesterName, &v.TutorName,
		); err != nil {
			return nil, fmt.Errorf("ListRequestsForTutor row scan: %w", err)
		}
		list = append(list, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListRequestsForTutor: %w", err)
	}
	return list, nil
}

func (r *PostgresTutoringRepo) ListRequestsByRequester(ctx context.Context, requesterID, schoolID uuid.UUID, status string) ([]domain.TutorRequestView, error) {
	status = normalizeOutgoingRequestStatusFilter(status)
	base := `SELECT tr.id, tr.school_id, tr.requester_id, tr.tutor_id, tr.topic_id,
			tr.topic_name, tr.status, tr.message, tr.created_at, tr.updated_at,
			req.name AS requester_name, tut.name AS tutor_name
		 FROM tutor_requests tr
		 JOIN users req ON req.id = tr.requester_id
		 JOIN users tut ON tut.id = tr.tutor_id
		 WHERE tr.requester_id = $1 AND tr.school_id = $2`
	var (
		query string
		args  []any
	)
	if status == "all" {
		query = base + ` ORDER BY tr.created_at DESC`
		args = []any{requesterID, schoolID}
	} else {
		query = base + ` AND tr.status = $3 ORDER BY tr.created_at DESC`
		args = []any{requesterID, schoolID, status}
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListRequestsByRequester: %w", err)
	}
	defer rows.Close()

	var list []domain.TutorRequestView
	for rows.Next() {
		var v domain.TutorRequestView
		if err := rows.Scan(
			&v.ID, &v.SchoolID, &v.RequesterID, &v.TutorID, &v.TopicID,
			&v.TopicName, &v.Status, &v.Message, &v.CreatedAt, &v.UpdatedAt,
			&v.RequesterName, &v.TutorName,
		); err != nil {
			return nil, fmt.Errorf("ListRequestsByRequester row scan: %w", err)
		}
		list = append(list, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListRequestsByRequester: %w", err)
	}
	return list, nil
}

func (r *PostgresTutoringRepo) UpdateRequestStatus(ctx context.Context, requestID, userID, schoolID uuid.UUID, status string) (*domain.TutorRequest, error) {
	var out domain.TutorRequest
	err := r.pool.QueryRow(ctx,
		`UPDATE tutor_requests
		 SET status = $4, updated_at = now()
		 WHERE id = $1
		   AND school_id = $3
		   AND (requester_id = $2 OR tutor_id = $2)
		 RETURNING id, school_id, requester_id, tutor_id, topic_id, topic_name, status, message, created_at, updated_at`,
		requestID, userID, schoolID, status,
	).Scan(
		&out.ID, &out.SchoolID, &out.RequesterID, &out.TutorID, &out.TopicID, &out.TopicName, &out.Status, &out.Message, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UpdateRequestStatus: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("UpdateRequestStatus: %w", err)
	}
	return &out, nil
}

func (r *PostgresTutoringRepo) FindTutorMatches(ctx context.Context, schoolID, requesterID uuid.UUID, topicName string, limit int) ([]domain.TutorMatch, error) {
	topicName = strings.TrimSpace(topicName)
	if topicName == "" {
		return []domain.TutorMatch{}, nil
	}
	if limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx,
		`SELECT sub.user_id, sub.user_name, sub.topic_name, sub.mastery_score, sub.review_count
		 FROM (
		   SELECT DISTINCT ON (u.id)
		     u.id AS user_id,
		     u.name AS user_name,
		     t.name AS topic_name,
		     tm.mastery_score,
		     tm.review_count
		   FROM users u
		   INNER JOIN topic_mastery tm ON tm.user_id = u.id AND tm.school_id = u.school_id
		   INNER JOIN topics t ON t.id = tm.topic_id AND t.school_id = tm.school_id
		   WHERE u.school_id = $1
		     AND u.id <> $2
		     AND tm.mastery_score >= 0.75
		     AND t.name ILIKE '%' || $3 || '%'
		   ORDER BY u.id, tm.mastery_score DESC
		 ) sub
		 ORDER BY sub.mastery_score DESC
		 LIMIT $4`,
		schoolID, requesterID, topicName, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("FindTutorMatches: %w", err)
	}
	defer rows.Close()

	var list []domain.TutorMatch
	for rows.Next() {
		var (
			uid     uuid.UUID
			name    string
			topic   string
			mastery float64
			reviews int
		)
		if err := rows.Scan(&uid, &name, &topic, &mastery, &reviews); err != nil {
			return nil, fmt.Errorf("FindTutorMatches row: %w", err)
		}
		rep := int(math.Min(float64(reviews)*2.5, 99))
		if rep < 1 {
			rep = 1
		}
		list = append(list, domain.TutorMatch{
			UserID:     uid,
			Name:       name,
			TopicName:  topic,
			Mastery:    mastery,
			Sessions:   reviews,
			Reputation: rep,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindTutorMatches: %w", err)
	}
	return list, nil
}

func (r *PostgresTutoringRepo) ListTeachingTopics(ctx context.Context, userID, schoolID uuid.UUID) ([]domain.TeachingTopic, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.name, tm.mastery_score
		 FROM topic_mastery tm
		 INNER JOIN topics t ON t.id = tm.topic_id AND t.school_id = tm.school_id
		 WHERE tm.user_id = $1 AND tm.school_id = $2 AND tm.mastery_score >= 0.75
		 ORDER BY tm.mastery_score DESC`,
		userID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListTeachingTopics: %w", err)
	}
	defer rows.Close()

	var list []domain.TeachingTopic
	for rows.Next() {
		var row domain.TeachingTopic
		if err := rows.Scan(&row.TopicID, &row.TopicName, &row.Mastery); err != nil {
			return nil, fmt.Errorf("ListTeachingTopics row: %w", err)
		}
		list = append(list, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListTeachingTopics: %w", err)
	}
	return list, nil
}
