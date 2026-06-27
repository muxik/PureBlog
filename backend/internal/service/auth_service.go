package service

import (
	"context"
	"errors"

	"github.com/muxik/PureBlog/backend/internal/auth"
	"github.com/muxik/PureBlog/backend/internal/domain"
)

// ErrInvalidCredentials is returned for bad logins and unusable refresh tokens.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthService implements authentication use-cases.
type AuthService struct {
	users  domain.UserRepository
	tokens *auth.TokenManager
}

// NewAuthService builds an AuthService.
func NewAuthService(users domain.UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

// TokenPair bundles an access + refresh token.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Login verifies credentials and issues a fresh token pair.
func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenPair, *domain.User, error) {
	u, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}
	ok, err := auth.VerifyPassword(password, u.PasswordHash)
	if err != nil || !ok {
		return nil, nil, ErrInvalidCredentials
	}
	pair, err := s.issue(u)
	if err != nil {
		return nil, nil, err
	}
	return pair, u, nil
}

// Refresh validates a refresh token and rotates it for a new pair.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.tokens.Parse(refreshToken)
	if err != nil || claims.Type != auth.TokenRefresh {
		return nil, ErrInvalidCredentials
	}
	u, err := s.users.GetByID(ctx, claims.UserID())
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issue(u)
}

func (s *AuthService) issue(u *domain.User) (*TokenPair, error) {
	access, err := s.tokens.Access(u.ID, u.Username)
	if err != nil {
		return nil, err
	}
	refresh, err := s.tokens.Refresh(u.ID, u.Username)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

// EnsureAdmin seeds the first administrator when the users table is empty.
func (s *AuthService) EnsureAdmin(ctx context.Context, username, password string) error {
	n, err := s.users.Count(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	return s.users.Create(ctx, &domain.User{
		Username:    username,
		PasswordHash: hash,
		DisplayName: username,
		Role:        "admin",
	})
}
