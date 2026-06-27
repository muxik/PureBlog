package domain

import "context"

// SiteSettings holds the site-wide configuration exposed publicly and edited by
// the admin. It is persisted as a single JSON document.
type SiteSettings struct {
	SiteName    string
	Description string
	Author      string
	AboutMd     string
	Social      map[string]string
	// DefaultDateFormat is either "numeric" or "lunar".
	DefaultDateFormat string
}

// SettingsRepository is the persistence port for the single site-settings row,
// defined where it is consumed and implemented by the store layer.
type SettingsRepository interface {
	Get(ctx context.Context) (*SiteSettings, error)
	Update(ctx context.Context, s *SiteSettings) error
}
