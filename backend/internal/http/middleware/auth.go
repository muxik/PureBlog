package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/auth"
)

const (
	ctxUserID   = "userID"
	ctxUsername = "username"
)

// RequireAuth rejects requests without a valid access token and stashes the
// authenticated user id/username in the Gin context.
func RequireAuth(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		claims, err := tokens.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil || claims.Type != auth.TokenAccess {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(ctxUserID, claims.UserID())
		c.Set(ctxUsername, claims.Username)
		c.Next()
	}
}

// UserID returns the authenticated user id stored by RequireAuth (0 if absent).
func UserID(c *gin.Context) int64 {
	if v, ok := c.Get(ctxUserID); ok {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}
