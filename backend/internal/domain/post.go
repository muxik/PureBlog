// Package domain holds the pure business entities and the repository
// interfaces they depend on. It imports neither GORM nor Gin — the dependency
// arrow always points inward to this package.
package domain

import (
	"context"
	"errors"
	"time"
)

// Sentinel errors translated by the store layer from driver-specific errors.
var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
	// ErrInvalidReference is for foreign-key violations from untrusted input
	// (e.g. a comment whose parentId does not exist).
	ErrInvalidReference = errors.New("invalid reference")
)

// PostStatus enumerates the lifecycle states of a post.
type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
)

// Post is the core article entity.
type Post struct {
	ID          int64
	Slug        string
	Title       string
	Summary     string
	ContentMD   string
	ContentHTML string
	CoverURL    string
	CategoryID  *int64
	Status      PostStatus
	Pinned      bool
	AuthorID    int64
	ViewCount   int64
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Associations
	TagIDs []int64 // write-side: tag ids to attach on create/update
	Tags   []Tag   // read-side: the attached tags
}

// PostListFilter parameterises a paginated post listing.
type PostListFilter struct {
	Status       PostStatus // empty = any status
	Query        string     // free-text search over title/summary
	CategorySlug string     // filter to posts whose category has this slug
	TagSlug      string     // filter to posts that have a tag with this slug
	Page         int        // 1-based
	PageSize     int
}

// PostRepository is the persistence port for posts, defined where it is
// consumed and implemented by the store layer.
type PostRepository interface {
	Create(ctx context.Context, p *Post) error
	Update(ctx context.Context, p *Post) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Post, error)
	GetBySlug(ctx context.Context, slug string) (*Post, error)
	List(ctx context.Context, f PostListFilter) (items []*Post, total int64, err error)
	IncrementViews(ctx context.Context, id int64) error
}
