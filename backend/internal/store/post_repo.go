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
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		return setPostTags(tx, m.ID, p.TagIDs)
	})
	if err != nil {
		return mapErr(err)
	}
	*p = *toDomainPost(m)
	return r.attachTags(ctx, p)
}

func (r *PostRepo) Update(ctx context.Context, p *domain.Post) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&PostModel{}).Where("id = ?", p.ID).Updates(map[string]any{
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
			return res.Error
		}
		if res.RowsAffected == 0 {
			return domain.ErrNotFound
		}
		return setPostTags(tx, p.ID, p.TagIDs)
	})
	if err != nil {
		return mapErr(err)
	}
	return r.attachTags(ctx, p)
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
	p := toDomainPost(&m)
	if err := r.attachTags(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepo) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	var m PostModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error; err != nil {
		return nil, mapErr(err)
	}
	p := toDomainPost(&m)
	if err := r.attachTags(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepo) List(ctx context.Context, f domain.PostListFilter) ([]*domain.Post, int64, error) {
	q := r.db.WithContext(ctx).Model(&PostModel{})
	if f.Status != "" {
		q = q.Where("posts.status = ?", string(f.Status))
	}
	if s := strings.TrimSpace(f.Query); s != "" {
		// Substring match across title, summary, and body. The leading-wildcard
		// ILIKE is index-backed by the pg_trgm GIN indexes (see migration 00003),
		// so widening to content_md does not force a sequential scan.
		like := "%" + s + "%"
		q = q.Where(
			"posts.title ILIKE ? OR posts.summary ILIKE ? OR posts.content_md ILIKE ?",
			like, like, like,
		)
	}
	if ts := strings.TrimSpace(f.TagSlug); ts != "" {
		// Use EXISTS to avoid duplicate rows from the many-to-many join table.
		q = q.Where(
			"EXISTS (SELECT 1 FROM post_tags JOIN tags ON tags.id = post_tags.tag_id WHERE post_tags.post_id = posts.id AND tags.slug = ?)",
			ts,
		)
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
		Order("posts.pinned DESC, COALESCE(posts.published_at, posts.created_at) DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&models).Error; err != nil {
		return nil, 0, mapErr(err)
	}

	items := make([]*domain.Post, len(models))
	for i := range models {
		items[i] = toDomainPost(&models[i])
	}
	if err := r.attachTags(ctx, items...); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *PostRepo) IncrementViews(ctx context.Context, id int64) error {
	return mapErr(r.db.WithContext(ctx).Model(&PostModel{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error)
}

// setPostTags replaces a post's tag associations inside the given transaction.
// A non-existent tag id triggers a foreign-key violation, which mapErr turns
// into domain.ErrInvalidReference (→ HTTP 400).
func setPostTags(tx *gorm.DB, postID int64, tagIDs []int64) error {
	if err := tx.Where("post_id = ?", postID).Delete(&PostTagModel{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}
	seen := make(map[int64]bool, len(tagIDs))
	rows := make([]PostTagModel, 0, len(tagIDs))
	for _, id := range tagIDs {
		if !seen[id] {
			seen[id] = true
			rows = append(rows, PostTagModel{PostID: postID, TagID: id})
		}
	}
	return tx.Create(&rows).Error
}

// attachTags loads the tags for the given posts in a single query and fills
// each post's Tags and TagIDs (avoiding N+1).
func (r *PostRepo) attachTags(ctx context.Context, posts ...*domain.Post) error {
	if len(posts) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(posts))
	byID := make(map[int64]*domain.Post, len(posts))
	for _, p := range posts {
		p.Tags = []domain.Tag{}
		p.TagIDs = []int64{}
		ids = append(ids, p.ID)
		byID[p.ID] = p
	}

	var rows []struct {
		PostID int64
		ID     int64
		Name   string
		Slug   string
	}
	if err := r.db.WithContext(ctx).
		Table("post_tags").
		Select("post_tags.post_id, tags.id, tags.name, tags.slug").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id IN ?", ids).
		Order("tags.name ASC").
		Scan(&rows).Error; err != nil {
		return mapErr(err)
	}

	for _, row := range rows {
		if p, ok := byID[row.PostID]; ok {
			p.Tags = append(p.Tags, domain.Tag{ID: row.ID, Name: row.Name, Slug: row.Slug})
			p.TagIDs = append(p.TagIDs, row.ID)
		}
	}
	return nil
}
