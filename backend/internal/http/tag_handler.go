package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/internal/service"
)

// ---- DTOs ----

type saveTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

type tagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type tagListResponse struct {
	Items []tagResponse `json:"items"`
}

func toTagResponse(t *domain.Tag) tagResponse {
	return tagResponse{ID: t.ID, Name: t.Name, Slug: t.Slug}
}

// ---- handlers ----

// listTags lists all tags.
//
//	@Summary	List tags
//	@Tags		tags
//	@Produce	json
//	@Success	200	{object}	tagListResponse
//	@Router		/tags [get]
func (s *Server) listTags(c *gin.Context) {
	items, err := s.tags.List(c.Request.Context())
	if err != nil {
		s.fail(c, err)
		return
	}
	out := make([]tagResponse, len(items))
	for i, t := range items {
		out[i] = toTagResponse(t)
	}
	c.JSON(http.StatusOK, tagListResponse{Items: out})
}

// createTag creates a tag (admin).
//
//	@Summary	Create a tag
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	saveTagRequest	true	"tag"
//	@Success	201	{object}	tagResponse
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/tags [post]
func (s *Server) createTag(c *gin.Context) {
	var req saveTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	t, err := s.tags.Create(c.Request.Context(), service.SaveTagInput{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, toTagResponse(t))
}

// deleteTag deletes a tag (admin).
//
//	@Summary	Delete a tag
//	@Tags		admin
//	@Security	BearerAuth
//	@Param		id	path	int	true	"tag id"
//	@Success	204	"no content"
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/tags/{id} [delete]
func (s *Server) deleteTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	if err := s.tags.Delete(c.Request.Context(), id); err != nil {
		s.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
