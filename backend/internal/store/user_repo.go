package store

import (
	"context"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"gorm.io/gorm"
)

// UserRepo implements domain.UserRepository.
type UserRepo struct{ db *gorm.DB }

// NewUserRepo builds a UserRepo.
func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

var _ domain.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	m := fromDomainUser(u)
	m.ID = 0
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return mapErr(err)
	}
	*u = *toDomainUser(m)
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, mapErr(err)
	}
	return toDomainUser(&m), nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		return nil, mapErr(err)
	}
	return toDomainUser(&m), nil
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&n).Error; err != nil {
		return 0, mapErr(err)
	}
	return n, nil
}
