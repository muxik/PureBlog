package service

import (
	"context"

	"github.com/muxik/PureBlog/backend/internal/domain"
)

// CommentService implements the comment use-cases. It depends on the comment
// repository and on the post repository to resolve and validate the target
// post (only published posts accept and expose comments).
type CommentService struct {
	repo  domain.CommentRepository
	posts domain.PostRepository
}

// NewCommentService builds a CommentService.
func NewCommentService(repo domain.CommentRepository, posts domain.PostRepository) *CommentService {
	return &CommentService{repo: repo, posts: posts}
}

// CreateCommentInput is the payload for submitting a comment.
type CreateCommentInput struct {
	AuthorName  string
	AuthorEmail string
	Content     string
	ParentID    *int64
	IP          string
	UserAgent   string
}

// CreateForSlug resolves the published post by slug and creates a pending
// comment for it. A missing or unpublished post yields domain.ErrNotFound.
func (s *CommentService) CreateForSlug(ctx context.Context, slug string, in CreateCommentInput) (*domain.Comment, error) {
	p, err := s.posts.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if p.Status != domain.StatusPublished {
		return nil, domain.ErrNotFound
	}

	// Validate a reply target: the parent must exist and belong to THIS post,
	// otherwise untrusted input could trigger a raw FK violation (500) or graft
	// a reply onto another post's thread.
	if in.ParentID != nil {
		parent, err := s.repo.GetByID(ctx, *in.ParentID)
		if err != nil {
			return nil, err // ErrNotFound when the parent is missing
		}
		if parent.PostID != p.ID {
			return nil, domain.ErrNotFound
		}
	}

	c := &domain.Comment{
		PostID:      p.ID,
		ParentID:    in.ParentID,
		AuthorName:  in.AuthorName,
		AuthorEmail: in.AuthorEmail,
		Content:     in.Content,
		Status:      string(domain.CommentPending),
		IP:          in.IP,
		UserAgent:   in.UserAgent,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// ListApprovedForSlug resolves the published post by slug and returns its
// approved comments oldest-first.
func (s *CommentService) ListApprovedForSlug(ctx context.Context, slug string) ([]*domain.Comment, error) {
	p, err := s.posts.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if p.Status != domain.StatusPublished {
		return nil, domain.ErrNotFound
	}
	return s.repo.ListByPost(ctx, p.ID, true)
}

// List returns comments filtered by status (empty = all), newest-first.
func (s *CommentService) List(ctx context.Context, status string) ([]*domain.Comment, error) {
	return s.repo.ListByStatus(ctx, status)
}

// SetStatus moderates a comment.
func (s *CommentService) SetStatus(ctx context.Context, id int64, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

// Delete removes a comment.
func (s *CommentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
