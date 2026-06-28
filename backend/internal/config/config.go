// Package config loads runtime configuration from environment variables.
package config

import (
	"os"
	"strings"
	"time"
)

// Config holds all runtime configuration for the server.
type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	AdminUsername   string
	AdminPassword   string
	CORSOrigins     []string
	// UploadDir is the filesystem directory where uploaded images are stored
	// and served from under /uploads.
	UploadDir string
	// PublicBaseURL is the externally reachable origin of the site, used to
	// build absolute URLs for uploaded images so they resolve from both the
	// public site and the admin preview (which lives on a different host).
	PublicBaseURL string
}

// Load reads configuration from the environment, applying sane defaults for
// local development.
func Load() Config {
	return Config{
		Port:            env("PORT", "8080"),
		DatabaseURL:     databaseURL(),
		JWTSecret:       env("JWT_SECRET", "dev-insecure-change-me"),
		AccessTokenTTL:  envDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: envDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		AdminUsername:   env("ADMIN_USERNAME", "muxi"),
		AdminPassword:   env("ADMIN_PASSWORD", "muxi-change-me"),
		CORSOrigins:     splitAndTrim(env("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		UploadDir:       env("UPLOAD_DIR", "./uploads"),
		PublicBaseURL:   strings.TrimRight(env("PUBLIC_ORIGIN", "http://localhost:8080"), "/"),
	}
}

func databaseURL() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return strings.Join([]string{
		"host=" + env("POSTGRES_HOST", "localhost"),
		"port=" + env("POSTGRES_PORT", "5432"),
		"user=" + env("POSTGRES_USER", "postgres"),
		"password=" + env("POSTGRES_PASSWORD", "pureblog"),
		"dbname=" + env("POSTGRES_DB", "pureblog"),
		"sslmode=" + env("POSTGRES_SSLMODE", "disable"),
		"TimeZone=UTC",
	}, " ")
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
