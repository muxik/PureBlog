package store

import (
	"context"
	"encoding/json"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"gorm.io/gorm"
)

// SettingsModel is the GORM persistence model for the single site-settings row
// (id = 1) whose payload lives in a JSONB column.
type SettingsModel struct {
	ID   int    `gorm:"primaryKey"`
	Data []byte `gorm:"column:data;type:jsonb"`
}

// TableName pins the table created by the goose migration.
func (SettingsModel) TableName() string { return "settings" }

// SettingsRepo implements domain.SettingsRepository.
type SettingsRepo struct{ db *gorm.DB }

// NewSettingsRepo builds a SettingsRepo.
func NewSettingsRepo(db *gorm.DB) *SettingsRepo { return &SettingsRepo{db: db} }

var _ domain.SettingsRepository = (*SettingsRepo)(nil)

// settingsJSON mirrors domain.SiteSettings for JSON (de)serialisation in the
// JSONB column, keeping the wire keys stable independent of the domain struct.
type settingsJSON struct {
	SiteName          string            `json:"siteName"`
	Description       string            `json:"description"`
	Author            string            `json:"author"`
	AboutMd           string            `json:"aboutMd"`
	Social            map[string]string `json:"social"`
	DefaultDateFormat string            `json:"defaultDateFormat"`
}

// Get loads the single settings row. An empty or "{}" payload maps to a zero
// SiteSettings value without error.
func (r *SettingsRepo) Get(ctx context.Context) (*domain.SiteSettings, error) {
	var m SettingsModel
	if err := r.db.WithContext(ctx).First(&m, 1).Error; err != nil {
		return nil, mapErr(err)
	}
	var s domain.SiteSettings
	if len(m.Data) > 0 {
		var raw settingsJSON
		if err := json.Unmarshal(m.Data, &raw); err != nil {
			return nil, err
		}
		s = domain.SiteSettings{
			SiteName:          raw.SiteName,
			Description:       raw.Description,
			Author:            raw.Author,
			AboutMd:           raw.AboutMd,
			Social:            raw.Social,
			DefaultDateFormat: raw.DefaultDateFormat,
		}
	}
	return &s, nil
}

// Update marshals the settings and writes them back to the single row.
func (r *SettingsRepo) Update(ctx context.Context, s *domain.SiteSettings) error {
	raw := settingsJSON{
		SiteName:          s.SiteName,
		Description:       s.Description,
		Author:            s.Author,
		AboutMd:           s.AboutMd,
		Social:            s.Social,
		DefaultDateFormat: s.DefaultDateFormat,
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return mapErr(r.db.WithContext(ctx).Model(&SettingsModel{}).
		Where("id = ?", 1).
		Update("data", bytes).Error)
}
