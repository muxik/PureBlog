// Package service holds the application use-cases. It depends on the domain
// ports and the renderer, never on Gin or GORM.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/internal/render"
)

// PostService implements the article use-cases.
type PostService struct {
	repo domain.PostRepository
	rdr  *render.Renderer
}

// NewPostService builds a PostService.
func NewPostService(repo domain.PostRepository, rdr *render.Renderer) *PostService {
	return &PostService{repo: repo, rdr: rdr}
}

// SavePostInput is the payload for creating or updating a post.
type SavePostInput struct {
	Title     string
	Slug      string
	Summary   string
	ContentMD string
	CoverURL  string
	Status    domain.PostStatus
	Pinned    bool
	AuthorID  int64
}

// Create renders the Markdown, derives a unique slug, and persists the post.
func (s *PostService) Create(ctx context.Context, in SavePostInput) (*domain.Post, error) {
	html, err := s.rdr.ToHTML(in.ContentMD)
	if err != nil {
		return nil, err
	}
	p := &domain.Post{
		Slug:        strings.TrimSpace(in.Slug),
		Title:       in.Title,
		Summary:     in.Summary,
		ContentMD:   in.ContentMD,
		ContentHTML: html,
		CoverURL:    in.CoverURL,
		Status:      normalizeStatus(in.Status),
		Pinned:      in.Pinned,
		AuthorID:    in.AuthorID,
	}
	if p.Slug == "" {
		p.Slug = Slugify(in.Title)
	}
	if p.Status == domain.StatusPublished {
		now := time.Now()
		p.PublishedAt = &now
	}

	base := p.Slug
	for attempt := 0; ; attempt++ {
		if attempt > 0 {
			p.Slug = fmt.Sprintf("%s-%d", base, attempt+1)
		}
		err := s.repo.Create(ctx, p)
		if err == nil {
			return p, nil
		}
		if errors.Is(err, domain.ErrConflict) && attempt < 20 {
			continue
		}
		return nil, err
	}
}

// Update mutates an existing post, re-rendering Markdown and stamping the
// publish time on the draft → published transition.
func (s *PostService) Update(ctx context.Context, id int64, in SavePostInput) (*domain.Post, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	html, err := s.rdr.ToHTML(in.ContentMD)
	if err != nil {
		return nil, err
	}

	existing.Title = in.Title
	existing.Summary = in.Summary
	existing.ContentMD = in.ContentMD
	existing.ContentHTML = html
	existing.CoverURL = in.CoverURL
	existing.Pinned = in.Pinned
	if s := strings.TrimSpace(in.Slug); s != "" {
		existing.Slug = s
	}

	newStatus := normalizeStatus(in.Status)
	if existing.Status != domain.StatusPublished && newStatus == domain.StatusPublished {
		now := time.Now()
		existing.PublishedAt = &now
	}
	existing.Status = newStatus

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete removes a post.
func (s *PostService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// List returns posts matching the filter together with the total count.
func (s *PostService) List(ctx context.Context, f domain.PostListFilter) ([]*domain.Post, int64, error) {
	return s.repo.List(ctx, f)
}

// GetByID returns any post by id (admin use).
func (s *PostService) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

// GetPublishedBySlug returns a published post by slug and bumps its view count.
func (s *PostService) GetPublishedBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	p, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if p.Status != domain.StatusPublished {
		return nil, domain.ErrNotFound
	}
	_ = s.repo.IncrementViews(ctx, p.ID)
	p.ViewCount++
	return p, nil
}

// Render converts Markdown to sanitised HTML for the admin live preview, using
// the same renderer as the publish path so preview == published output.
func (s *PostService) Render(md string) (string, error) {
	return s.rdr.ToHTML(md)
}

func normalizeStatus(st domain.PostStatus) domain.PostStatus {
	if st == domain.StatusPublished {
		return domain.StatusPublished
	}
	return domain.StatusDraft
}

// Slugify produces a URL-safe slug from a title, falling back to "post" when
// the title contains no ASCII alphanumerics (e.g. a purely Chinese title) — the
// caller de-duplicates by appending a numeric suffix.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '-' || r == '_':
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "post"
	}
	return out
}
