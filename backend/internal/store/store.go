// Package store implements the domain repository ports on top of GORM +
// PostgreSQL. GORM types live here only; they never leak into the domain.
package store

import (
	"errors"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to PostgreSQL via GORM. TranslateError lets GORM surface
// portable sentinel errors (ErrRecordNotFound, ErrDuplicatedKey).
func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Warn),
		TranslateError: true,
	})
}

// mapErr translates GORM errors into domain sentinel errors.
func mapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return domain.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return domain.ErrConflict
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return domain.ErrInvalidReference
	default:
		return err
	}
}
