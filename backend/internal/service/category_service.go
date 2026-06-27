package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/muxik/PureBlog/backend/internal/domain"
)

// CategoryService implements the category use-cases.
type CategoryService struct {
	repo domain.CategoryRepository
}

// NewCategoryService builds a CategoryService.
func NewCategoryService(repo domain.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// SaveCategoryInput is the payload for creating or updating a category.
type SaveCategoryInput struct {
	ParentID    *int64
	Name        string
	Slug        string
	Description string
	Sort        int
}

// Create derives a unique slug (falling back to the name) and persists the
// category, de-duplicating on conflict by appending -2, -3, ... exactly like
// PostService.Create.
func (s *CategoryService) Create(ctx context.Context, in SaveCategoryInput) (*domain.Category, error) {
	c := &domain.Category{
		ParentID:    in.ParentID,
		Name:        in.Name,
		Slug:        strings.TrimSpace(in.Slug),
		Description: in.Description,
		Sort:        in.Sort,
	}
	if c.Slug == "" {
		c.Slug = Slugify(in.Name)
	}

	base := c.Slug
	for attempt := 0; ; attempt++ {
		if attempt > 0 {
			c.Slug = fmt.Sprintf("%s-%d", base, attempt+1)
		}
		err := s.repo.Create(ctx, c)
		if err == nil {
			return c, nil
		}
		if errors.Is(err, domain.ErrConflict) && attempt < 20 {
			continue
		}
		return nil, err
	}
}

// Update mutates an existing category, re-deriving the slug from the name when
// it is left blank and de-duplicating on conflict.
func (s *CategoryService) Update(ctx context.Context, id int64, in SaveCategoryInput) (*domain.Category, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.ParentID = in.ParentID
	existing.Name = in.Name
	existing.Description = in.Description
	existing.Sort = in.Sort
	if slug := strings.TrimSpace(in.Slug); slug != "" {
		existing.Slug = slug
	} else {
		existing.Slug = Slugify(in.Name)
	}

	base := existing.Slug
	for attempt := 0; ; attempt++ {
		if attempt > 0 {
			existing.Slug = fmt.Sprintf("%s-%d", base, attempt+1)
		}
		err := s.repo.Update(ctx, existing)
		if err == nil {
			return existing, nil
		}
		if errors.Is(err, domain.ErrConflict) && attempt < 20 {
			continue
		}
		return nil, err
	}
}

// Delete removes a category.
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// List returns all categories ordered by sort then id.
func (s *CategoryService) List(ctx context.Context) ([]*domain.Category, error) {
	return s.repo.List(ctx)
}
