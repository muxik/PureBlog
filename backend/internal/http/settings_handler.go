package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/domain"
)

// ---- DTOs ----

// settingsResponse mirrors domain.SiteSettings on the wire.
type settingsResponse struct {
	SiteName          string            `json:"siteName"`
	Description       string            `json:"description"`
	Author            string            `json:"author"`
	AboutMd           string            `json:"aboutMd"`
	Social            map[string]string `json:"social"`
	DefaultDateFormat string            `json:"defaultDateFormat"`
}

// updateSettingsRequest is the payload for replacing the site settings.
type updateSettingsRequest struct {
	SiteName          string            `json:"siteName"`
	Description       string            `json:"description"`
	Author            string            `json:"author"`
	AboutMd           string            `json:"aboutMd"`
	Social            map[string]string `json:"social"`
	DefaultDateFormat string            `json:"defaultDateFormat"`
}

func toSettingsResponse(s *domain.SiteSettings) settingsResponse {
	return settingsResponse{
		SiteName:          s.SiteName,
		Description:       s.Description,
		Author:            s.Author,
		AboutMd:           s.AboutMd,
		Social:            s.Social,
		DefaultDateFormat: s.DefaultDateFormat,
	}
}

// ---- handlers ----

// getSettings returns the public site settings.
//
//	@Summary	Get site settings
//	@Tags		settings
//	@Produce	json
//	@Success	200	{object}	settingsResponse
//	@Router		/settings [get]
func (s *Server) getSettings(c *gin.Context) {
	settings, err := s.settings.Get(c.Request.Context())
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toSettingsResponse(settings))
}

// updateSettings replaces the site settings (admin).
//
//	@Summary	Replace site settings
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	updateSettingsRequest	true	"settings"
//	@Success	200	{object}	settingsResponse
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/settings [put]
func (s *Server) updateSettings(c *gin.Context) {
	var req updateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	in := &domain.SiteSettings{
		SiteName:          req.SiteName,
		Description:       req.Description,
		Author:            req.Author,
		AboutMd:           req.AboutMd,
		Social:            req.Social,
		DefaultDateFormat: req.DefaultDateFormat,
	}
	if err := s.settings.Update(c.Request.Context(), in); err != nil {
		s.fail(c, err)
		return
	}
	settings, err := s.settings.Get(c.Request.Context())
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toSettingsResponse(settings))
}
