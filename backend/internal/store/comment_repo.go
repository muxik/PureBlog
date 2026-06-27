package store

import (
	"context"
	"time"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"gorm.io/gorm"
)

// CommentModel is the GORM persistence model for comments.
type CommentModel struct {
	ID          int64     `gorm:"primaryKey"`
	PostID      int64     `gorm:"column:post_id;index"`
	ParentID    *int64    `gorm:"column:parent_id;index"`
	AuthorName  string    `gorm:"column:author_name"`
	AuthorEmail string    `gorm:"column:author_email"`
	Content     string    `gorm:"column:content"`
	Status      string    `gorm:"column:status;index;size:20"`
	IP          string    `gorm:"column:ip"`
	UserAgent   string    `gorm:"column:user_agent"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

// TableName pins the table created by the goose migration.
func (CommentModel) TableName() string { return "comments" }

func toDomainComment(m *CommentModel) *domain.Comment {
	return &domain.Comment{
		ID:          m.ID,
		PostID:      m.PostID,
		ParentID:    m.ParentID,
		AuthorName:  m.AuthorName,
		AuthorEmail: m.AuthorEmail,
		Content:     m.Content,
		Status:      m.Status,
		IP:          m.IP,
		UserAgent:   m.UserAgent,
		CreatedAt:   m.CreatedAt,
	}
}

func fromDomainComment(c *domain.Comment) *CommentModel {
	return &CommentModel{
		ID:          c.ID,
		PostID:      c.PostID,
		ParentID:    c.ParentID,
		AuthorName:  c.AuthorName,
		AuthorEmail: c.AuthorEmail,
		Content:     c.Content,
		Status:      c.Status,
		IP:          c.IP,
		UserAgent:   c.UserAgent,
		CreatedAt:   c.CreatedAt,
	}
}

// CommentRepo implements domain.CommentRepository.
type CommentRepo struct{ db *gorm.DB }

// NewCommentRepo builds a CommentRepo.
func NewCommentRepo(db *gorm.DB) *CommentRepo { return &CommentRepo{db: db} }

var _ domain.CommentRepository = (*CommentRepo)(nil)

func (r *CommentRepo) Create(ctx context.Context, c *domain.Comment) error {
	m := fromDomainComment(c)
	m.ID = 0
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapErr(err)
	}
	*c = *toDomainComment(m)
	return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int64) (*domain.Comment, error) {
	var m CommentModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, mapErr(err)
	}
	return toDomainComment(&m), nil
}

func (r *CommentRepo) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&CommentModel{}, id)
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CommentRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	res := r.db.WithContext(ctx).Model(&CommentModel{}).
		Where("id = ?", id).
		Update("status", status)
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CommentRepo) ListByPost(ctx context.Context, postID int64, onlyApproved bool) ([]*domain.Comment, error) {
	q := r.db.WithContext(ctx).Model(&CommentModel{}).Where("post_id = ?", postID)
	if onlyApproved {
		q = q.Where("status = ?", string(domain.CommentApproved))
	}

	var models []CommentModel
	if err := q.Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, mapErr(err)
	}

	items := make([]*domain.Comment, len(models))
	for i := range models {
		items[i] = toDomainComment(&models[i])
	}
	return items, nil
}

func (r *CommentRepo) ListByStatus(ctx context.Context, status string) ([]*domain.Comment, error) {
	q := r.db.WithContext(ctx).Model(&CommentModel{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var models []CommentModel
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, mapErr(err)
	}

	items := make([]*domain.Comment, len(models))
	for i := range models {
		items[i] = toDomainComment(&models[i])
	}
	return items, nil
}
