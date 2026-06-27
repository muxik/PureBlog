package domain

import (
	"context"
	"time"
)

// CommentStatus enumerates the moderation states of a comment.
type CommentStatus string

const (
	CommentPending  CommentStatus = "pending"
	CommentApproved CommentStatus = "approved"
)

// Comment is a reader-submitted comment attached to a post. Replies reference
// their parent via ParentID (nil for a top-level comment).
type Comment struct {
	ID          int64
	PostID      int64
	ParentID    *int64
	AuthorName  string
	AuthorEmail string
	Content     string
	Status      string
	IP          string
	UserAgent   string
	CreatedAt   time.Time
}

// CommentRepository is the persistence port for comments, defined where it is
// consumed and implemented by the store layer.
type CommentRepository interface {
	Create(ctx context.Context, c *Comment) error
	GetByID(ctx context.Context, id int64) (*Comment, error)
	Delete(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	// ListByPost returns comments for a post oldest-first; when onlyApproved is
	// true only approved comments are returned.
	ListByPost(ctx context.Context, postID int64, onlyApproved bool) ([]*Comment, error)
	// ListByStatus returns comments newest-first; an empty status means all.
	ListByStatus(ctx context.Context, status string) ([]*Comment, error)
}
