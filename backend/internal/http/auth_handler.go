package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/service"
)

// login authenticates a user and returns an access + refresh token pair.
//
//	@Summary	Log in
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	loginRequest	true	"credentials"
//	@Success	200	{object}	tokenResponse
//	@Failure	401	{object}	errorResponse
//	@Router		/auth/login [post]
func (s *Server) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	pair, user, err := s.authS.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid credentials"})
			return
		}
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		User:         toUserResponse(user),
	})
}

// refresh rotates a refresh token for a new token pair.
//
//	@Summary	Refresh tokens
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	refreshRequest	true	"refresh token"
//	@Success	200	{object}	tokenResponse
//	@Failure	401	{object}	errorResponse
//	@Router		/auth/refresh [post]
func (s *Server) refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	pair, err := s.authS.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

// logout is a no-op for stateless JWTs; the client discards its tokens.
//
//	@Summary	Log out
//	@Tags		auth
//	@Success	204	"no content"
//	@Router		/auth/logout [post]
func (s *Server) logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
