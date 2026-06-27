package store

import (
	"context"
	"strings"
	"time"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"gorm.io/gorm"
)

// PostRepo implements domain.PostRepository.
type PostRepo struct{ db *gorm.DB }

// NewPostRepo builds a PostRepo.
func NewPostRepo(db *gorm.DB) *PostRepo { return &PostRepo{db: db} }

var _ domain.PostRepository = (*PostRepo)(nil)

func (r *PostRepo) Create(ctx context.Context, p *domain.Post) error {
	m := fromDomainPost(p)
	m.ID = 0
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapErr(err)
	}
	*p = *toDomainPost(m)
	return nil
}

func (r *PostRepo) Update(ctx context.Context, p *domain.Post) error {
	res := r.db.WithContext(ctx).Model(&PostModel{}).Where("id = ?", p.ID).Updates(map[string]any{
		"slug":         p.Slug,
		"title":        p.Title,
		"summary":      p.Summary,
		"content_md":   p.ContentMD,
		"content_html": p.ContentHTML,
		"cover_url":    p.CoverURL,
		"status":       string(p.Status),
		"pinned":       p.Pinned,
		"published_at": p.PublishedAt,
		"updated_at":   time.Now(),
	})
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostRepo) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&PostModel{}, id)
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostRepo) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	var m PostModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, mapErr(err)
	}
	return toDomainPost(&m), nil
}

func (r *PostRepo) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	var m PostModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error; err != nil {
		return nil, mapErr(err)
	}
	return toDomainPost(&m), nil
}

func (r *PostRepo) List(ctx context.Context, f domain.PostListFilter) ([]*domain.Post, int64, error) {
	q := r.db.WithContext(ctx).Model(&PostModel{})
	if f.Status != "" {
		q = q.Where("status = ?", string(f.Status))
	}
	if s := strings.TrimSpace(f.Query); s != "" {
		like := "%" + s + "%"
		q = q.Where("title ILIKE ? OR summary ILIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, mapErr(err)
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	size := f.PageSize
	if size < 1 || size > 100 {
		size = 10
	}

	var models []PostModel
	if err := q.
		Order("pinned DESC, COALESCE(published_at, created_at) DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&models).Error; err != nil {
		return nil, 0, mapErr(err)
	}

	items := make([]*domain.Post, len(models))
	for i := range models {
		items[i] = toDomainPost(&models[i])
	}
	return items, total, nil
}

func (r *PostRepo) IncrementViews(ctx context.Context, id int64) error {
	return mapErr(r.db.WithContext(ctx).Model(&PostModel{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error)
}
