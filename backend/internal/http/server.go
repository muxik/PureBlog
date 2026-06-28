package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/auth"
	"github.com/muxik/PureBlog/backend/internal/config"
	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/internal/http/middleware"
	"github.com/muxik/PureBlog/backend/internal/service"
)

// Server wires the HTTP handlers to the application services.
type Server struct {
	cfg        config.Config
	posts    *service.PostService
	authS    *service.AuthService
	tokens   *auth.TokenManager
	tags     *service.TagService
	comments *service.CommentService
	settings *service.SettingsService
}

// NewServer builds a Server.
func NewServer(cfg config.Config, posts *service.PostService, authS *service.AuthService, tokens *auth.TokenManager, tags *service.TagService, comments *service.CommentService, settings *service.SettingsService) *Server {
	return &Server{
		cfg:      cfg,
		posts:    posts,
		authS:    authS,
		tokens:   tokens,
		tags:     tags,
		comments: comments,
		settings: settings,
	}
}

// Router builds the Gin engine with all routes registered.
func (s *Server) Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(s.cfg.CORSOrigins))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		// public
		v1.GET("/posts", s.listPosts)
		v1.GET("/posts/:slug", s.getPostBySlug)
		v1.GET("/tags", s.listTags)
		v1.GET("/settings", s.getSettings)
		v1.GET("/posts/:slug/comments", s.listPostComments)
		v1.POST("/posts/:slug/comments", s.createPostComment)

		// auth
		v1.POST("/auth/login", s.login)
		v1.POST("/auth/refresh", s.refresh)
		v1.POST("/auth/logout", s.logout)

		// admin (JWT-protected)
		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAuth(s.tokens))
		{
			admin.GET("/posts", s.adminListPosts)
			admin.GET("/posts/:id", s.adminGetPost)
			admin.POST("/posts", s.createPost)
			admin.PUT("/posts/:id", s.updatePost)
			admin.DELETE("/posts/:id", s.deletePost)
			admin.POST("/render", s.renderMarkdown)

			admin.POST("/tags", s.createTag)
			admin.DELETE("/tags/:id", s.deleteTag)

			admin.GET("/comments", s.adminListComments)
			admin.PUT("/comments/:id", s.moderateComment)
			admin.DELETE("/comments/:id", s.deleteComment)

			admin.PUT("/settings", s.updateSettings)
		}
	}

	return r
}

// fail maps domain errors onto HTTP status codes.
func (s *Server) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Error: "not found"})
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, errorResponse{Error: "conflict"})
	case errors.Is(err, domain.ErrInvalidReference):
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid reference"})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal error"})
	}
}

func (s *Server) badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
}
