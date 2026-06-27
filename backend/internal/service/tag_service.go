package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/muxik/PureBlog/backend/internal/domain"
)

// TagService implements the tag use-cases.
type TagService struct {
	repo domain.TagRepository
}

// NewTagService builds a TagService.
func NewTagService(repo domain.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// SaveTagInput is the payload for creating a tag.
type SaveTagInput struct {
	Name string
	Slug string
}

// Create derives a slug from the name when none is supplied and de-duplicates
// the slug by appending a numeric suffix on conflict.
func (s *TagService) Create(ctx context.Context, in SaveTagInput) (*domain.Tag, error) {
	t := &domain.Tag{
		Name: strings.TrimSpace(in.Name),
		Slug: strings.TrimSpace(in.Slug),
	}
	if t.Slug == "" {
		t.Slug = Slugify(t.Name)
	}

	base := t.Slug
	for attempt := 0; ; attempt++ {
		if attempt > 0 {
			t.Slug = fmt.Sprintf("%s-%d", base, attempt+1)
		}
		err := s.repo.Create(ctx, t)
		if err == nil {
			return t, nil
		}
		if errors.Is(err, domain.ErrConflict) && attempt < 20 {
			continue
		}
		return nil, err
	}
}

// Delete removes a tag.
func (s *TagService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// List returns all tags.
func (s *TagService) List(ctx context.Context) ([]*domain.Tag, error) {
	return s.repo.List(ctx)
}
