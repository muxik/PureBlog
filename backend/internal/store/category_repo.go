package store

import (
	"context"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"gorm.io/gorm"
)

// CategoryModel is the GORM persistence model for categories.
type CategoryModel struct {
	ID          int64  `gorm:"primaryKey"`
	ParentID    *int64 `gorm:"column:parent_id;index"`
	Name        string `gorm:"column:name"`
	Slug        string `gorm:"column:slug;uniqueIndex;size:200"`
	Description string `gorm:"column:description"`
	Sort        int    `gorm:"column:sort;index"`
}

// TableName pins the table created by the goose migration.
func (CategoryModel) TableName() string { return "categories" }

func toDomainCategory(m *CategoryModel) *domain.Category {
	return &domain.Category{
		ID:          m.ID,
		ParentID:    m.ParentID,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		Sort:        m.Sort,
	}
}

func fromDomainCategory(c *domain.Category) *CategoryModel {
	return &CategoryModel{
		ID:          c.ID,
		ParentID:    c.ParentID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Sort:        c.Sort,
	}
}

// CategoryRepo implements domain.CategoryRepository.
type CategoryRepo struct{ db *gorm.DB }

// NewCategoryRepo builds a CategoryRepo.
func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db: db} }

var _ domain.CategoryRepository = (*CategoryRepo)(nil)

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	m := fromDomainCategory(c)
	m.ID = 0
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapErr(err)
	}
	*c = *toDomainCategory(m)
	return nil
}

func (r *CategoryRepo) Update(ctx context.Context, c *domain.Category) error {
	res := r.db.WithContext(ctx).Model(&CategoryModel{}).Where("id = ?", c.ID).Updates(map[string]any{
		"parent_id":   c.ParentID,
		"name":        c.Name,
		"slug":        c.Slug,
		"description": c.Description,
		"sort":        c.Sort,
	})
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&CategoryModel{}, id)
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	var m CategoryModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, mapErr(err)
	}
	return toDomainCategory(&m), nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]*domain.Category, error) {
	var models []CategoryModel
	if err := r.db.WithContext(ctx).
		Order("sort ASC, id ASC").
		Find(&models).Error; err != nil {
		return nil, mapErr(err)
	}
	items := make([]*domain.Category, len(models))
	for i := range models {
		items[i] = toDomainCategory(&models[i])
	}
	return items, nil
}
