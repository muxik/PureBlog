package service

import (
	"context"

	"github.com/muxik/PureBlog/backend/internal/domain"
)

// SettingsService implements the site-settings use-cases.
type SettingsService struct {
	repo domain.SettingsRepository
}

// NewSettingsService builds a SettingsService.
func NewSettingsService(repo domain.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

// Get returns the site settings, defaulting the date format to "numeric" when
// it has not been configured yet.
func (s *SettingsService) Get(ctx context.Context) (*domain.SiteSettings, error) {
	settings, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	if settings.DefaultDateFormat == "" {
		settings.DefaultDateFormat = "numeric"
	}
	return settings, nil
}

// Update replaces the site settings.
func (s *SettingsService) Update(ctx context.Context, in *domain.SiteSettings) error {
	return s.repo.Update(ctx, in)
}
