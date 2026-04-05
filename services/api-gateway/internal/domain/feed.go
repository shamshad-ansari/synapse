package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FeedPost struct {
	ID        uuid.UUID  `json:"id"`
	SchoolID  uuid.UUID  `json:"school_id"`
	UserID    uuid.UUID  `json:"user_id"`
	CourseID  *uuid.UUID `json:"course_id"`
	TopicID   *uuid.UUID `json:"topic_id"`
	PostType  string     `json:"post_type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Upvotes   int        `json:"upvotes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// FeedPostView is returned in list responses — includes author info and upvoted flag.
type FeedPostView struct {
	FeedPost
	AuthorName   string `json:"author_name"`
	Upvoted      bool   `json:"upvoted"`
	CommentCount int    `json:"comment_count"`
}

type FeedComment struct {
	ID        uuid.UUID  `json:"id"`
	SchoolID  uuid.UUID  `json:"school_id"`
	PostID    uuid.UUID  `json:"post_id"`
	UserID    uuid.UUID  `json:"user_id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// FeedCommentView includes the author name for list responses.
type FeedCommentView struct {
	FeedComment
	AuthorName string `json:"author_name"`
}

type FeedRepository interface {
	CreatePost(ctx context.Context, schoolID, userID uuid.UUID, title, body, postType string, courseID, topicID *uuid.UUID) (*FeedPostView, error)
	ListPosts(ctx context.Context, schoolID, requestingUserID uuid.UUID, limit, offset int) ([]FeedPostView, error)
	ToggleUpvote(ctx context.Context, postID, userID, schoolID uuid.UUID) (upvotes int, upvoted bool, err error)

	// Comments
	CreateComment(ctx context.Context, schoolID, postID, userID uuid.UUID, body string, parentID *uuid.UUID) (*FeedCommentView, error)
	ListComments(ctx context.Context, schoolID, postID uuid.UUID) ([]FeedCommentView, error)

	// Delete post (must validate that userID == post.user_id)
	DeletePost(ctx context.Context, postID, userID, schoolID uuid.UUID) error
}

type FeedService interface {
	CreatePost(ctx context.Context, schoolID, userID uuid.UUID, title, body, postType string, courseID, topicID *uuid.UUID) (*FeedPostView, error)
	ListPosts(ctx context.Context, schoolID, requestingUserID uuid.UUID, limit, offset int) ([]FeedPostView, error)
	ToggleUpvote(ctx context.Context, postID, userID, schoolID uuid.UUID) (upvotes int, upvoted bool, err error)

	// Comments
	CreateComment(ctx context.Context, schoolID, postID, userID uuid.UUID, body string, parentID *uuid.UUID) (*FeedCommentView, error)
	ListComments(ctx context.Context, schoolID, postID uuid.UUID) ([]FeedCommentView, error)

	// Delete post (must validate that userID == post.user_id)
	DeletePost(ctx context.Context, postID, userID, schoolID uuid.UUID) error
}
