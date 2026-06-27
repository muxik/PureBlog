package domain

import (
	"context"
	"time"
)

// User is an author / administrator of the blog.
type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	DisplayName  string
	Bio          string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRepository is the persistence port for users.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Count(ctx context.Context) (int64, error)
}
