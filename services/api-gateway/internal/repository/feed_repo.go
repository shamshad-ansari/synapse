package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

type PostgresFeedRepo struct {
	pool *pgxpool.Pool
}

func NewFeedRepo(pool *pgxpool.Pool) *PostgresFeedRepo {
	return &PostgresFeedRepo{pool: pool}
}

func (r *PostgresFeedRepo) CreatePost(ctx context.Context, schoolID, userID uuid.UUID, title, body, postType string, courseID, topicID *uuid.UUID) (*domain.FeedPostView, error) {
	var out domain.FeedPostView
	err := r.pool.QueryRow(ctx,
		`WITH inserted AS (
		   INSERT INTO feed_posts (school_id, user_id, course_id, topic_id, post_type, title, body)
		   VALUES ($1, $2, $3, $4, $5, $6, $7)
		   RETURNING id, school_id, user_id, course_id, topic_id, post_type, title, body, upvotes, created_at, updated_at
		 )
		 SELECT
		   inserted.id, inserted.school_id, inserted.user_id, inserted.course_id, inserted.topic_id,
		   inserted.post_type, inserted.title, inserted.body, inserted.upvotes, inserted.created_at, inserted.updated_at,
		   u.name AS author_name,
		   FALSE AS upvoted,
		   0::int AS comment_count
		 FROM inserted
		 JOIN users u ON u.id = inserted.user_id`,
		schoolID, userID, courseID, topicID, postType, title, body,
	).Scan(
		&out.ID, &out.SchoolID, &out.UserID, &out.CourseID, &out.TopicID,
		&out.PostType, &out.Title, &out.Body, &out.Upvotes, &out.CreatedAt, &out.UpdatedAt,
		&out.AuthorName, &out.Upvoted, &out.CommentCount,
	)
	if err != nil {
		return nil, fmt.Errorf("CreatePost: %w", err)
	}
	return &out, nil
}

func (r *PostgresFeedRepo) ListPosts(ctx context.Context, schoolID, requestingUserID uuid.UUID, limit, offset int) ([]domain.FeedPostView, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT
		  fp.id, fp.school_id, fp.user_id, fp.course_id, fp.topic_id,
		  fp.post_type, fp.title, fp.body, fp.upvotes, fp.created_at, fp.updated_at,
		  u.name AS author_name,
		  EXISTS(SELECT 1 FROM feed_upvotes fu WHERE fu.post_id = fp.id AND fu.user_id = $2) AS upvoted,
		  COALESCE(fc.cnt, 0)::int AS comment_count
		FROM feed_posts fp
		JOIN users u ON u.id = fp.user_id
		LEFT JOIN (
		  SELECT post_id, COUNT(*)::bigint AS cnt
		  FROM feed_comments
		  WHERE school_id = $1
		  GROUP BY post_id
		) fc ON fc.post_id = fp.id
		WHERE fp.school_id = $1
		ORDER BY fp.upvotes DESC, fp.created_at DESC
		LIMIT $3 OFFSET $4`,
		schoolID, requestingUserID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("ListPosts: %w", err)
	}
	defer rows.Close()

	var list []domain.FeedPostView
	for rows.Next() {
		var v domain.FeedPostView
		if err := rows.Scan(
			&v.ID, &v.SchoolID, &v.UserID, &v.CourseID, &v.TopicID,
			&v.PostType, &v.Title, &v.Body, &v.Upvotes, &v.CreatedAt, &v.UpdatedAt,
			&v.AuthorName, &v.Upvoted, &v.CommentCount,
		); err != nil {
			return nil, fmt.Errorf("ListPosts: %w", err)
		}
		list = append(list, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListPosts: %w", err)
	}
	return list, nil
}

func (r *PostgresFeedRepo) ToggleUpvote(ctx context.Context, postID, userID, schoolID uuid.UUID) (int, bool, error) {
	// PostgreSQL aborts the whole transaction on a failed statement; do not rely on
	// INSERT-then-catch-23505. Branch on EXISTS instead.
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var already bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM feed_upvotes WHERE user_id = $1 AND post_id = $2)`,
		userID, postID,
	).Scan(&already)
	if err != nil {
		return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
	}

	if !already {
		_, err = tx.Exec(ctx,
			`INSERT INTO feed_upvotes (user_id, post_id) VALUES ($1, $2)`,
			userID, postID,
		)
		if err != nil {
			return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
		}
		var upvotes int
		err = tx.QueryRow(ctx,
			`UPDATE feed_posts SET upvotes = upvotes + 1, updated_at = now()
			 WHERE id = $1 AND school_id = $2
			 RETURNING upvotes`,
			postID, schoolID,
		).Scan(&upvotes)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				_, _ = tx.Exec(ctx, `DELETE FROM feed_upvotes WHERE user_id = $1 AND post_id = $2`, userID, postID)
				return 0, false, fmt.Errorf("ToggleUpvote: %w", domain.ErrNotFound)
			}
			return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
		}
		if err := tx.Commit(ctx); err != nil {
			return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
		}
		return upvotes, true, nil
	}

	tag, err := tx.Exec(ctx,
		`DELETE FROM feed_upvotes WHERE user_id = $1 AND post_id = $2`,
		userID, postID,
	)
	if err != nil {
		return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return 0, false, fmt.Errorf("ToggleUpvote: %w", domain.ErrNotFound)
	}
	var upvotes int
	err = tx.QueryRow(ctx,
		`UPDATE feed_posts SET upvotes = upvotes - 1, updated_at = now()
		 WHERE id = $1 AND school_id = $2
		 RETURNING upvotes`,
		postID, schoolID,
	).Scan(&upvotes)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false, fmt.Errorf("ToggleUpvote: %w", domain.ErrNotFound)
		}
		return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, false, fmt.Errorf("ToggleUpvote: %w", err)
	}
	return upvotes, false, nil
}

func (r *PostgresFeedRepo) CreateComment(
	ctx context.Context,
	schoolID, postID, userID uuid.UUID,
	body string,
	parentID *uuid.UUID,
) (*domain.FeedCommentView, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateComment: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var targetPostID uuid.UUID
	err = tx.QueryRow(ctx,
		`SELECT id
		 FROM feed_posts
		 WHERE id = $1 AND school_id = $2`,
		postID, schoolID,
	).Scan(&targetPostID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("CreateComment: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("CreateComment: %w", err)
	}

	if parentID != nil {
		var parentPostID uuid.UUID
		err = tx.QueryRow(ctx,
			`SELECT post_id
			 FROM feed_comments
			 WHERE id = $1 AND school_id = $2`,
			*parentID, schoolID,
		).Scan(&parentPostID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("CreateComment: %w", &domain.ValidationError{Message: "parent comment not found"})
			}
			return nil, fmt.Errorf("CreateComment: %w", err)
		}
		if parentPostID != postID {
			return nil, fmt.Errorf("CreateComment: %w", &domain.ValidationError{Message: "parent comment must belong to this post"})
		}
	}

	var out domain.FeedCommentView
	err = tx.QueryRow(ctx,
		`WITH inserted AS (
		   INSERT INTO feed_comments (school_id, post_id, user_id, parent_id, body)
		   VALUES ($1, $2, $3, $4, $5)
		   RETURNING id, school_id, post_id, user_id, parent_id, body, created_at, updated_at
		 )
		 SELECT
		   inserted.id, inserted.school_id, inserted.post_id, inserted.user_id, inserted.parent_id,
		   inserted.body, inserted.created_at, inserted.updated_at, u.name AS author_name
		 FROM inserted
		 JOIN users u ON u.id = inserted.user_id`,
		schoolID, postID, userID, parentID, body,
	).Scan(&out.ID, &out.SchoolID, &out.PostID, &out.UserID, &out.ParentID,
		&out.Body, &out.CreatedAt, &out.UpdatedAt, &out.AuthorName)
	if err != nil {
		return nil, fmt.Errorf("CreateComment: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("CreateComment: %w", err)
	}
	return &out, nil
}

func (r *PostgresFeedRepo) ListComments(
	ctx context.Context,
	schoolID, postID uuid.UUID,
) ([]domain.FeedCommentView, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT fc.id, fc.school_id, fc.post_id, fc.user_id, fc.parent_id,
		        fc.body, fc.created_at, fc.updated_at, u.name AS author_name
		 FROM feed_comments fc
		 JOIN users u ON u.id = fc.user_id
		 WHERE fc.post_id = $1 AND fc.school_id = $2
		 ORDER BY fc.created_at ASC`,
		postID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListComments: %w", err)
	}
	defer rows.Close()
	var list []domain.FeedCommentView
	for rows.Next() {
		var v domain.FeedCommentView
		if err := rows.Scan(
			&v.ID, &v.SchoolID, &v.PostID, &v.UserID, &v.ParentID,
			&v.Body, &v.CreatedAt, &v.UpdatedAt, &v.AuthorName,
		); err != nil {
			return nil, fmt.Errorf("ListComments row: %w", err)
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *PostgresFeedRepo) DeletePost(
	ctx context.Context,
	postID, userID, schoolID uuid.UUID,
) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM feed_posts
		 WHERE id = $1 AND user_id = $2 AND school_id = $3`,
		postID, userID, schoolID,
	)
	if err != nil {
		return fmt.Errorf("DeletePost: %w", err)
	}
	if tag.RowsAffected() == 0 {
		// Either not found or not the author — return a not-found error
		// so the handler returns 404, which leaks no information about
		// whether the post exists under a different author.
		return fmt.Errorf("DeletePost: %w", domain.ErrNotFound)
	}
	return nil
}
