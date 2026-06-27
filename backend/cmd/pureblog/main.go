// Command pureblog is the PureBlog v3 API server entrypoint. It wires
// config → db → migrations → admin seed → services → HTTP server.
package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/muxik/PureBlog/backend/internal/auth"
	"github.com/muxik/PureBlog/backend/internal/config"
	httpapi "github.com/muxik/PureBlog/backend/internal/http"
	"github.com/muxik/PureBlog/backend/internal/render"
	"github.com/muxik/PureBlog/backend/internal/service"
	"github.com/muxik/PureBlog/backend/internal/store"
	"github.com/muxik/PureBlog/backend/migrations"
	"github.com/pressly/goose/v3"
)

// @title			PureBlog API
// @version		3.0
// @description	PureBlog v3 backend API.
// @BasePath		/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in				header
// @name			Authorization
func main() {
	cfg := config.Load()

	db, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		fatal("open database", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		fatal("acquire *sql.DB", err)
	}
	if err := runMigrations(sqlDB); err != nil {
		fatal("run migrations", err)
	}

	rdr := render.New()
	tokens := auth.NewTokenManager(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	postRepo := store.NewPostRepo(db)
	postSvc := service.NewPostService(postRepo, rdr)
	authSvc := service.NewAuthService(store.NewUserRepo(db), tokens)
	categorySvc := service.NewCategoryService(store.NewCategoryRepo(db))
	tagSvc := service.NewTagService(store.NewTagRepo(db))
	commentSvc := service.NewCommentService(store.NewCommentRepo(db), postRepo)
	settingsSvc := service.NewSettingsService(store.NewSettingsRepo(db))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := authSvc.EnsureAdmin(ctx, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		fatal("seed admin", err)
	}

	srv := httpapi.NewServer(cfg, postSvc, authSvc, tokens, categorySvc, tagSvc, commentSvc, settingsSvc)
	slog.Info("PureBlog API listening", "addr", ":"+cfg.Port)
	if err := srv.Router().Run(":" + cfg.Port); err != nil {
		fatal("serve", err)
	}
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}

func fatal(msg string, err error) {
	slog.Error(msg, "err", err)
	os.Exit(1)
}
