package store

import (
	"time"

	"github.com/muxik/PureBlog/backend/internal/domain"
)

// PostModel is the GORM persistence model for posts.
type PostModel struct {
	ID          int64      `gorm:"primaryKey"`
	Slug        string     `gorm:"column:slug;uniqueIndex;size:200"`
	Title       string     `gorm:"column:title"`
	Summary     string     `gorm:"column:summary"`
	ContentMD   string     `gorm:"column:content_md"`
	ContentHTML string     `gorm:"column:content_html"`
	CoverURL    string     `gorm:"column:cover_url"`
	Status      string     `gorm:"column:status;index;size:20"`
	Pinned      bool       `gorm:"column:pinned"`
	AuthorID    int64      `gorm:"column:author_id;index"`
	ViewCount   int64      `gorm:"column:view_count"`
	PublishedAt *time.Time `gorm:"column:published_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

// TableName pins the table created by the goose migration.
func (PostModel) TableName() string { return "posts" }

func toDomainPost(m *PostModel) *domain.Post {
	return &domain.Post{
		ID:          m.ID,
		Slug:        m.Slug,
		Title:       m.Title,
		Summary:     m.Summary,
		ContentMD:   m.ContentMD,
		ContentHTML: m.ContentHTML,
		CoverURL:    m.CoverURL,
		Status:      domain.PostStatus(m.Status),
		Pinned:      m.Pinned,
		AuthorID:    m.AuthorID,
		ViewCount:   m.ViewCount,
		PublishedAt: m.PublishedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func fromDomainPost(p *domain.Post) *PostModel {
	return &PostModel{
		ID:          p.ID,
		Slug:        p.Slug,
		Title:       p.Title,
		Summary:     p.Summary,
		ContentMD:   p.ContentMD,
		ContentHTML: p.ContentHTML,
		CoverURL:    p.CoverURL,
		Status:      string(p.Status),
		Pinned:      p.Pinned,
		AuthorID:    p.AuthorID,
		ViewCount:   p.ViewCount,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// UserModel is the GORM persistence model for users.
type UserModel struct {
	ID           int64     `gorm:"primaryKey"`
	Username     string    `gorm:"column:username;uniqueIndex;size:80"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	DisplayName  string    `gorm:"column:display_name"`
	Bio          string    `gorm:"column:bio"`
	Role         string    `gorm:"column:role;size:20"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

// TableName pins the table created by the goose migration.
func (UserModel) TableName() string { return "users" }

func toDomainUser(m *UserModel) *domain.User {
	return &domain.User{
		ID:           m.ID,
		Username:     m.Username,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		DisplayName:  m.DisplayName,
		Bio:          m.Bio,
		Role:         m.Role,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func fromDomainUser(u *domain.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
		Bio:          u.Bio,
		Role:         u.Role,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
