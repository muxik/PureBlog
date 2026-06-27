package store

import (
	"context"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"gorm.io/gorm"
)

// TagModel is the GORM persistence model for tags.
type TagModel struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
	Slug string `gorm:"column:slug;uniqueIndex;size:200"`
}

// TableName pins the table created by the goose migration.
func (TagModel) TableName() string { return "tags" }

func toDomainTag(m *TagModel) *domain.Tag {
	return &domain.Tag{
		ID:   m.ID,
		Name: m.Name,
		Slug: m.Slug,
	}
}

func fromDomainTag(t *domain.Tag) *TagModel {
	return &TagModel{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
}

// TagRepo implements domain.TagRepository.
type TagRepo struct{ db *gorm.DB }

// NewTagRepo builds a TagRepo.
func NewTagRepo(db *gorm.DB) *TagRepo { return &TagRepo{db: db} }

var _ domain.TagRepository = (*TagRepo)(nil)

func (r *TagRepo) Create(ctx context.Context, t *domain.Tag) error {
	m := fromDomainTag(t)
	m.ID = 0
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapErr(err)
	}
	*t = *toDomainTag(m)
	return nil
}

func (r *TagRepo) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&TagModel{}, id)
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TagRepo) GetByID(ctx context.Context, id int64) (*domain.Tag, error) {
	var m TagModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, mapErr(err)
	}
	return toDomainTag(&m), nil
}

func (r *TagRepo) List(ctx context.Context) ([]*domain.Tag, error) {
	var models []TagModel
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&models).Error; err != nil {
		return nil, mapErr(err)
	}
	items := make([]*domain.Tag, len(models))
	for i := range models {
		items[i] = toDomainTag(&models[i])
	}
	return items, nil
}
