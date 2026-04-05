package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

// formatVectorLiteral builds a PostgreSQL vector literal for casting to vector.
func formatVectorLiteral(v []float32) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, x := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatFloat(float64(x), 'f', -1, 32))
	}
	b.WriteByte(']')
	return b.String()
}

type PostgresLearningRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresLearningRepo(pool *pgxpool.Pool) *PostgresLearningRepo {
	return &PostgresLearningRepo{pool: pool}
}

func (r *PostgresLearningRepo) CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	var out domain.Course
	err := r.pool.QueryRow(ctx,
		`INSERT INTO courses (school_id, user_id, name, term, color, lms_course_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, school_id, user_id, name, term, color, lms_course_id, created_at`,
		course.SchoolID, course.UserID, course.Name, course.Term, course.Color, course.LMSCourseID,
	).Scan(
		&out.ID, &out.SchoolID, &out.UserID, &out.Name, &out.Term, &out.Color, &out.LMSCourseID, &out.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateCourse: %w", err)
	}
	return &out, nil
}

// UpsertCourseFromLMS inserts a new course or updates an existing one matched by
// (user_id, lms_course_id). This is idempotent — safe to call on every sync.
func (r *PostgresLearningRepo) UpsertCourseFromLMS(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	var out domain.Course
	err := r.pool.QueryRow(ctx,
		`INSERT INTO courses (school_id, user_id, name, term, color, lms_course_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (user_id, lms_course_id) WHERE lms_course_id IS NOT NULL
		 DO UPDATE SET name = EXCLUDED.name, term = EXCLUDED.term
		 RETURNING id, school_id, user_id, name, term, color, lms_course_id, created_at`,
		course.SchoolID, course.UserID, course.Name, course.Term, course.Color, course.LMSCourseID,
	).Scan(
		&out.ID, &out.SchoolID, &out.UserID, &out.Name, &out.Term, &out.Color, &out.LMSCourseID, &out.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("UpsertCourseFromLMS: %w", err)
	}
	return &out, nil
}

func (r *PostgresLearningRepo) ListCourses(ctx context.Context, userID, schoolID uuid.UUID) ([]*domain.Course, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, school_id, user_id, name, term, color, lms_course_id, created_at
		 FROM courses
		 WHERE user_id = $1 AND school_id = $2
		 ORDER BY created_at DESC`,
		userID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListCourses: %w", err)
	}
	defer rows.Close()

	var list []*domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(&c.ID, &c.SchoolID, &c.UserID, &c.Name, &c.Term, &c.Color, &c.LMSCourseID, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("ListCourses: %w", err)
		}
		list = append(list, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListCourses: %w", err)
	}
	return list, nil
}

func (r *PostgresLearningRepo) GetCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) (*domain.Course, error) {
	var c domain.Course
	err := r.pool.QueryRow(ctx,
		`SELECT id, school_id, user_id, name, term, color, lms_course_id, created_at
		 FROM courses
		 WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		courseID, userID, schoolID,
	).Scan(&c.ID, &c.SchoolID, &c.UserID, &c.Name, &c.Term, &c.Color, &c.LMSCourseID, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetCourse: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("GetCourse: %w", err)
	}
	return &c, nil
}

func (r *PostgresLearningRepo) DeleteCourse(ctx context.Context, courseID, userID, schoolID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	_, err = tx.Exec(ctx,
		`DELETE FROM review_events
		 WHERE school_id = $1 AND user_id = $2
		   AND flashcard_id IN (
		     SELECT id FROM flashcards WHERE course_id = $3 AND school_id = $1 AND user_id = $2
		   )`,
		schoolID, userID, courseID,
	)
	if err != nil {
		return fmt.Errorf("DeleteCourse: review_events: %w", err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM scheduler_states
		 WHERE school_id = $1 AND user_id = $2
		   AND flashcard_id IN (
		     SELECT id FROM flashcards WHERE course_id = $3 AND school_id = $1 AND user_id = $2
		   )`,
		schoolID, userID, courseID,
	)
	if err != nil {
		return fmt.Errorf("DeleteCourse: scheduler_states: %w", err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM flashcards
		 WHERE course_id = $1 AND user_id = $2 AND school_id = $3`,
		courseID, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteCourse: flashcards: %w", err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM topic_mastery
		 WHERE school_id = $1 AND user_id = $2
		   AND topic_id IN (
		     SELECT id FROM topics WHERE course_id = $3 AND school_id = $1 AND user_id = $2
		   )`,
		schoolID, userID, courseID,
	)
	if err != nil {
		return fmt.Errorf("DeleteCourse: topic_mastery: %w", err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM topics
		 WHERE course_id = $1 AND user_id = $2 AND school_id = $3`,
		courseID, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteCourse: topics: %w", err)
	}

	tag, err := tx.Exec(ctx,
		`DELETE FROM courses
		 WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		courseID, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteCourse: courses: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteCourse: %w", domain.ErrNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}
	return nil
}

func (r *PostgresLearningRepo) CreateTopic(ctx context.Context, topic *domain.Topic) (*domain.Topic, error) {
	var out domain.Topic
	err := r.pool.QueryRow(ctx,
		`INSERT INTO topics (school_id, course_id, user_id, name, parent_topic_id, exam_weight, source)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, school_id, course_id, user_id, name, parent_topic_id, exam_weight, source, created_at`,
		topic.SchoolID, topic.CourseID, topic.UserID, topic.Name, topic.ParentTopicID, topic.ExamWeight, topic.Source,
	).Scan(
		&out.ID, &out.SchoolID, &out.CourseID, &out.UserID, &out.Name,
		&out.ParentTopicID, &out.ExamWeight, &out.Source, &out.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateTopic: %w", err)
	}
	return &out, nil
}

func (r *PostgresLearningRepo) ListTopics(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.Topic, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, school_id, course_id, user_id, name, parent_topic_id, exam_weight, source, created_at
		 FROM topics
		 WHERE course_id = $1 AND user_id = $2 AND school_id = $3
		 ORDER BY created_at ASC`,
		courseID, userID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListTopics: %w", err)
	}
	defer rows.Close()

	var list []*domain.Topic
	for rows.Next() {
		var t domain.Topic
		if err := rows.Scan(
			&t.ID, &t.SchoolID, &t.CourseID, &t.UserID, &t.Name,
			&t.ParentTopicID, &t.ExamWeight, &t.Source, &t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListTopics: %w", err)
		}
		list = append(list, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListTopics: %w", err)
	}
	return list, nil
}

func (r *PostgresLearningRepo) CreateFlashcard(ctx context.Context, fc *domain.Flashcard) (*domain.Flashcard, error) {
	var out domain.Flashcard
	err := r.pool.QueryRow(ctx,
		`INSERT INTO flashcards (school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by, visibility)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by, visibility, created_at`,
		fc.SchoolID, fc.CourseID, fc.UserID, fc.TopicID, fc.CardType, fc.Prompt, fc.Answer, fc.CreatedBy, fc.Visibility,
	).Scan(
		&out.ID, &out.SchoolID, &out.CourseID, &out.UserID, &out.TopicID,
		&out.CardType, &out.Prompt, &out.Answer, &out.CreatedBy, &out.Visibility, &out.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateFlashcard: %w", err)
	}
	return &out, nil
}

func (r *PostgresLearningRepo) ListFlashcards(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.Flashcard, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by, visibility, created_at
		 FROM flashcards
		 WHERE course_id = $1 AND user_id = $2 AND school_id = $3
		 ORDER BY created_at DESC`,
		courseID, userID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListFlashcards: %w", err)
	}
	defer rows.Close()

	var list []*domain.Flashcard
	for rows.Next() {
		var f domain.Flashcard
		if err := rows.Scan(
			&f.ID, &f.SchoolID, &f.CourseID, &f.UserID, &f.TopicID,
			&f.CardType, &f.Prompt, &f.Answer, &f.CreatedBy, &f.Visibility, &f.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListFlashcards: %w", err)
		}
		list = append(list, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListFlashcards: %w", err)
	}
	return list, nil
}

func (r *PostgresLearningRepo) GetFlashcard(ctx context.Context, id, userID, schoolID uuid.UUID) (*domain.Flashcard, error) {
	var f domain.Flashcard
	err := r.pool.QueryRow(ctx,
		`SELECT id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by, visibility, created_at
		 FROM flashcards
		 WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		id, userID, schoolID,
	).Scan(
		&f.ID, &f.SchoolID, &f.CourseID, &f.UserID, &f.TopicID,
		&f.CardType, &f.Prompt, &f.Answer, &f.CreatedBy, &f.Visibility, &f.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetFlashcard: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("GetFlashcard: %w", err)
	}
	return &f, nil
}

func (r *PostgresLearningRepo) DeleteFlashcard(ctx context.Context, id, userID, schoolID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("DeleteFlashcard: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	_, err = tx.Exec(ctx,
		`DELETE FROM review_events WHERE flashcard_id = $1 AND user_id = $2 AND school_id = $3`,
		id, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteFlashcard: review_events: %w", err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM scheduler_states WHERE flashcard_id = $1 AND user_id = $2 AND school_id = $3`,
		id, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteFlashcard: scheduler_states: %w", err)
	}

	tag, err := tx.Exec(ctx,
		`DELETE FROM flashcards WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		id, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteFlashcard: flashcards: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteFlashcard: %w", domain.ErrNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("DeleteFlashcard: %w", err)
	}
	return nil
}

// BatchCreateFlashcards inserts cards one-by-one (MVP). TODO: bulk insert / COPY for performance.
func (r *PostgresLearningRepo) BatchCreateFlashcards(ctx context.Context, cards []*domain.Flashcard) ([]*domain.Flashcard, error) {
	out := make([]*domain.Flashcard, 0, len(cards))
	for _, c := range cards {
		created, err := r.CreateFlashcard(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("BatchCreateFlashcards: %w", err)
		}
		out = append(out, created)
	}
	return out, nil
}

func (r *PostgresLearningRepo) GetDueCards(ctx context.Context, userID, schoolID uuid.UUID, limit int, courseID *uuid.UUID) ([]*domain.DueCard, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if courseID != nil {
		rows, err = r.pool.Query(ctx,
			`SELECT f.id, f.prompt, f.answer, f.card_type, COALESCE(t.name, '') AS topic_name,
			        s.ease_factor, s.interval_days, s.due_at
			 FROM flashcards f
			 LEFT JOIN topics t ON t.id = f.topic_id AND t.school_id = f.school_id
			 JOIN scheduler_states s ON s.flashcard_id = f.id AND s.user_id = $1 AND s.school_id = $2
			 WHERE s.user_id = $1 AND f.school_id = $2 AND s.school_id = $2 AND s.due_at <= now()
			   AND f.course_id = $3
			 ORDER BY s.due_at ASC
			 LIMIT $4`,
			userID, schoolID, *courseID, limit,
		)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT f.id, f.prompt, f.answer, f.card_type, COALESCE(t.name, '') AS topic_name,
			        s.ease_factor, s.interval_days, s.due_at
			 FROM flashcards f
			 LEFT JOIN topics t ON t.id = f.topic_id AND t.school_id = f.school_id
			 JOIN scheduler_states s ON s.flashcard_id = f.id AND s.user_id = $1 AND s.school_id = $2
			 WHERE s.user_id = $1 AND f.school_id = $2 AND s.school_id = $2 AND s.due_at <= now()
			 ORDER BY s.due_at ASC
			 LIMIT $3`,
			userID, schoolID, limit,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("GetDueCards: %w", err)
	}
	defer rows.Close()

	var list []*domain.DueCard
	for rows.Next() {
		var d domain.DueCard
		if err := rows.Scan(
			&d.FlashcardID, &d.Prompt, &d.Answer, &d.CardType, &d.TopicName,
			&d.EaseFactor, &d.IntervalDays, &d.DueAt,
		); err != nil {
			return nil, fmt.Errorf("GetDueCards: %w", err)
		}
		list = append(list, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetDueCards: %w", err)
	}
	return list, nil
}

func (r *PostgresLearningRepo) UpsertSchedulerState(ctx context.Context, state *domain.SchedulerState) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO scheduler_states (flashcard_id, user_id, school_id, ease_factor, interval_days, due_at, lapse_count, last_review_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (flashcard_id, user_id) DO UPDATE SET
		   ease_factor = EXCLUDED.ease_factor,
		   interval_days = EXCLUDED.interval_days,
		   due_at = EXCLUDED.due_at,
		   lapse_count = EXCLUDED.lapse_count,
		   last_review_at = EXCLUDED.last_review_at`,
		state.FlashcardID, state.UserID, state.SchoolID, state.EaseFactor, state.IntervalDays,
		state.DueAt, state.LapseCount, state.LastReviewAt,
	)
	if err != nil {
		return fmt.Errorf("UpsertSchedulerState: %w", err)
	}
	return nil
}

func (r *PostgresLearningRepo) GetSchedulerState(ctx context.Context, flashcardID, userID, schoolID uuid.UUID) (*domain.SchedulerState, error) {
	var s domain.SchedulerState
	err := r.pool.QueryRow(ctx,
		`SELECT flashcard_id, user_id, school_id, ease_factor, interval_days, due_at, lapse_count, last_review_at
		 FROM scheduler_states
		 WHERE flashcard_id = $1 AND user_id = $2 AND school_id = $3`,
		flashcardID, userID, schoolID,
	).Scan(
		&s.FlashcardID, &s.UserID, &s.SchoolID, &s.EaseFactor, &s.IntervalDays,
		&s.DueAt, &s.LapseCount, &s.LastReviewAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetSchedulerState: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("GetSchedulerState: %w", err)
	}
	return &s, nil
}

func (r *PostgresLearningRepo) CreateReviewEvent(ctx context.Context, event *domain.ReviewEvent) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO review_events (id, school_id, user_id, flashcard_id, session_id, ts, correct, confidence, confused, response_time_ms, ease_before, interval_before)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		event.ID, event.SchoolID, event.UserID, event.FlashcardID, event.SessionID, event.Ts,
		event.Correct, event.Confidence, event.Confused, event.ResponseTimeMs, event.EaseBefore, event.IntervalBefore,
	)
	if err != nil {
		return fmt.Errorf("CreateReviewEvent: %w", err)
	}
	return nil
}

func (r *PostgresLearningRepo) CreateNoteText(ctx context.Context, note *domain.NoteText) (*domain.NoteText, error) {
	var out domain.NoteText
	err := r.pool.QueryRow(ctx,
		`INSERT INTO note_texts (school_id, course_id, user_id, topic_id, title, content)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at, updated_at`,
		note.SchoolID, note.CourseID, note.UserID, note.TopicID, note.Title, note.Content,
	).Scan(&out.ID, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateNoteText: %w", err)
	}
	out.SchoolID = note.SchoolID
	out.CourseID = note.CourseID
	out.UserID = note.UserID
	out.TopicID = note.TopicID
	out.Title = note.Title
	out.Content = note.Content
	return &out, nil
}

func (r *PostgresLearningRepo) UpdateNoteTextEmbedding(ctx context.Context, noteID, schoolID uuid.UUID, embedding []float32) error {
	vec := formatVectorLiteral(embedding)
	tag, err := r.pool.Exec(ctx,
		`UPDATE note_texts SET embedding = $1::vector, embedded_at = now(), updated_at = now()
		 WHERE id = $2 AND school_id = $3`,
		vec, noteID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("UpdateNoteTextEmbedding: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateNoteTextEmbedding: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *PostgresLearningRepo) ListNoteTexts(ctx context.Context, courseID, userID, schoolID uuid.UUID) ([]*domain.NoteText, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, school_id, course_id, user_id, topic_id, title, content, embedded_at, created_at, updated_at
		 FROM note_texts
		 WHERE course_id = $1 AND user_id = $2 AND school_id = $3
		 ORDER BY updated_at DESC`,
		courseID, userID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListNoteTexts: %w", err)
	}
	defer rows.Close()

	var list []*domain.NoteText
	for rows.Next() {
		var n domain.NoteText
		if err := rows.Scan(
			&n.ID, &n.SchoolID, &n.CourseID, &n.UserID, &n.TopicID,
			&n.Title, &n.Content, &n.EmbeddedAt, &n.CreatedAt, &n.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListNoteTexts: %w", err)
		}
		list = append(list, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListNoteTexts: %w", err)
	}
	return list, nil
}

func (r *PostgresLearningRepo) GetNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) (*domain.NoteText, error) {
	var n domain.NoteText
	err := r.pool.QueryRow(ctx,
		`SELECT id, school_id, course_id, user_id, topic_id, title, content, created_at, updated_at
		 FROM note_texts
		 WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		noteID, userID, schoolID,
	).Scan(
		&n.ID, &n.SchoolID, &n.CourseID, &n.UserID, &n.TopicID,
		&n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetNoteText: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("GetNoteText: %w", err)
	}
	return &n, nil
}

func (r *PostgresLearningRepo) UpdateNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID, title, content string, topicID *uuid.UUID) (*domain.NoteText, error) {
	var n domain.NoteText
	err := r.pool.QueryRow(ctx,
		`UPDATE note_texts
		 SET title = $4, content = $5, topic_id = $6, updated_at = now()
		 WHERE id = $1 AND user_id = $2 AND school_id = $3
		 RETURNING id, school_id, course_id, user_id, topic_id, title, content, created_at, updated_at`,
		noteID, userID, schoolID, title, content, topicID,
	).Scan(
		&n.ID, &n.SchoolID, &n.CourseID, &n.UserID, &n.TopicID,
		&n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UpdateNoteText: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("UpdateNoteText: %w", err)
	}
	return &n, nil
}

func (r *PostgresLearningRepo) DeleteNoteText(ctx context.Context, noteID, userID, schoolID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM note_texts WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		noteID, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeleteNoteText: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteNoteText: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *PostgresLearningRepo) SearchSimilarFlashcards(ctx context.Context, userID, schoolID uuid.UUID, queryVec []float32, limit int) ([]*domain.Flashcard, error) {
	if limit <= 0 {
		limit = 5
	}
	vec := formatVectorLiteral(queryVec)
	rows, err := r.pool.Query(ctx,
		`SELECT id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by, visibility, created_at
		 FROM flashcards
		 WHERE user_id = $1 AND school_id = $2
		   AND embedding IS NOT NULL
		 ORDER BY embedding <=> $3::vector
		 LIMIT $4`,
		userID, schoolID, vec, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("SearchSimilarFlashcards: %w", err)
	}
	defer rows.Close()

	var list []*domain.Flashcard
	for rows.Next() {
		var f domain.Flashcard
		if err := rows.Scan(
			&f.ID, &f.SchoolID, &f.CourseID, &f.UserID, &f.TopicID,
			&f.CardType, &f.Prompt, &f.Answer, &f.CreatedBy, &f.Visibility, &f.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("SearchSimilarFlashcards: %w", err)
		}
		list = append(list, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("SearchSimilarFlashcards: %w", err)
	}
	return list, nil
}

func (r *PostgresLearningRepo) UpdateFlashcardEmbedding(ctx context.Context, flashcardID, userID, schoolID uuid.UUID, embedding []float32) error {
	vec := formatVectorLiteral(embedding)
	tag, err := r.pool.Exec(ctx,
		`UPDATE flashcards SET embedding = $1::vector WHERE id = $2 AND user_id = $3 AND school_id = $4`,
		vec, flashcardID, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("UpdateFlashcardEmbedding: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateFlashcardEmbedding: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *PostgresLearningRepo) GetConfusionInsights(ctx context.Context, userID, schoolID uuid.UUID, courseID *uuid.UUID, windowDays int) (*domain.ConfusionInsight, error) {
	if windowDays <= 0 {
		windowDays = 14
	}
	if windowDays > 90 {
		windowDays = 90
	}

	out := &domain.ConfusionInsight{WindowDays: windowDays}

	whereCourse := ""
	args := []any{userID, schoolID, windowDays}
	if courseID != nil {
		whereCourse = " AND f.course_id = $4"
		args = append(args, *courseID)
		out.CourseID = courseID
	}

	querySummary := fmt.Sprintf(`
SELECT
  COALESCE(COUNT(*), 0) AS confused_events,
  COALESCE(COUNT(DISTINCT COALESCE(f.topic_id::text, f.id::text)), 0) AS hotspot_count,
  COALESCE(MAX(c.name), '') AS course_name
FROM review_events re
JOIN flashcards f ON f.id = re.flashcard_id AND f.user_id = $1 AND f.school_id = $2
JOIN courses c ON c.id = f.course_id AND c.user_id = $1 AND c.school_id = $2
WHERE re.user_id = $1
  AND re.school_id = $2
  AND re.confused = true
  AND re.ts >= now() - make_interval(days => $3)
  %s
`, whereCourse)
	if err := r.pool.QueryRow(ctx, querySummary, args...).Scan(&out.ConfusedEvents, &out.HotspotCount, &out.CourseName); err != nil {
		return nil, fmt.Errorf("GetConfusionInsights summary: %w", err)
	}

	queryTopTopic := fmt.Sprintf(`
SELECT COALESCE(t.name, 'Unscoped')
FROM review_events re
JOIN flashcards f ON f.id = re.flashcard_id AND f.user_id = $1 AND f.school_id = $2
LEFT JOIN topics t ON t.id = f.topic_id AND t.school_id = $2
WHERE re.user_id = $1
  AND re.school_id = $2
  AND re.confused = true
  AND re.ts >= now() - make_interval(days => $3)
  %s
GROUP BY COALESCE(t.name, 'Unscoped')
ORDER BY COUNT(*) DESC, COALESCE(t.name, 'Unscoped') ASC
LIMIT 1
`, whereCourse)
	if err := r.pool.QueryRow(ctx, queryTopTopic, args...).Scan(&out.TopTopicName); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetConfusionInsights top topic: %w", err)
		}
		out.TopTopicName = ""
	}

	return out, nil
}

func (r *PostgresLearningRepo) ListNoteMetrics(ctx context.Context, courseID, userID, schoolID uuid.UUID, windowDays int) ([]*domain.NoteMetric, error) {
	if windowDays <= 0 {
		windowDays = 30
	}
	if windowDays > 120 {
		windowDays = 120
	}

	rows, err := r.pool.Query(ctx, `
WITH topic_rollup AS (
  SELECT
    f.topic_id AS topic_id,
    COUNT(f.id) AS connected_cards,
    COUNT(re.id) AS review_count,
    COUNT(re.id) FILTER (WHERE re.confused = true) AS confused_count,
    COUNT(re.id) FILTER (WHERE re.correct = true) AS correct_count,
    MAX(re.ts) AS last_review_at
  FROM flashcards f
  LEFT JOIN review_events re
    ON re.flashcard_id = f.id
   AND re.user_id = $2
   AND re.school_id = $3
   AND re.ts >= now() - make_interval(days => $4)
  WHERE f.course_id = $1
    AND f.user_id = $2
    AND f.school_id = $3
  GROUP BY f.topic_id
)
SELECT
  n.id AS note_id,
  COALESCE(
    ROUND(
      CASE
        WHEN tr.review_count = 0 THEN 0
        ELSE
          (
            ((tr.correct_count::numeric / tr.review_count::numeric) * 100)
            - ((tr.confused_count::numeric / tr.review_count::numeric) * 40)
          )
      END
    ),
    0
  )::int AS readiness_pct,
  COALESCE(tr.connected_cards, 0) AS connected_cards,
  COALESCE(tr.review_count, 0) AS review_count,
  COALESCE(tr.confused_count, 0) AS confused_count,
  CASE
    WHEN COALESCE(tr.review_count, 0) = 0 THEN false
    ELSE (tr.confused_count::numeric / tr.review_count::numeric) >= 0.25
  END AS confusion_flag,
  tr.last_review_at
FROM note_texts n
LEFT JOIN topic_rollup tr
  ON n.topic_id IS NOT DISTINCT FROM tr.topic_id
WHERE n.course_id = $1
  AND n.user_id = $2
  AND n.school_id = $3
ORDER BY n.updated_at DESC
`, courseID, userID, schoolID, windowDays)
	if err != nil {
		return nil, fmt.Errorf("ListNoteMetrics: %w", err)
	}
	defer rows.Close()

	var out []*domain.NoteMetric
	for rows.Next() {
		var m domain.NoteMetric
		var lastReview sql.NullTime
		if err := rows.Scan(
			&m.NoteID,
			&m.ReadinessPct,
			&m.ConnectedCards,
			&m.ReviewCount,
			&m.ConfusedCount,
			&m.ConfusionFlag,
			&lastReview,
		); err != nil {
			return nil, fmt.Errorf("ListNoteMetrics: %w", err)
		}
		if lastReview.Valid {
			t := lastReview.Time.UTC()
			m.LastReviewAt = &t
		}
		out = append(out, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListNoteMetrics: %w", err)
	}
	return out, nil
}
