package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/internal/service"
)

// ---- DTOs ----

type saveCategoryRequest struct {
	ParentID    *int64 `json:"parentId"`
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

type categoryResponse struct {
	ID          int64  `json:"id"`
	ParentID    *int64 `json:"parentId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

type categoryListResponse struct {
	Items []categoryResponse `json:"items"`
}

func toCategoryResponse(c *domain.Category) categoryResponse {
	return categoryResponse{
		ID:          c.ID,
		ParentID:    c.ParentID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Sort:        c.Sort,
	}
}

func toCategoryInput(req saveCategoryRequest) service.SaveCategoryInput {
	return service.SaveCategoryInput{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Sort:        req.Sort,
	}
}

// ---- handlers ----

// listCategories lists all categories (flat, ordered by sort then id).
//
//	@Summary	List categories
//	@Tags		categories
//	@Produce	json
//	@Success	200	{object}	categoryListResponse
//	@Router		/categories [get]
func (s *Server) listCategories(c *gin.Context) {
	items, err := s.categories.List(c.Request.Context())
	if err != nil {
		s.fail(c, err)
		return
	}
	out := make([]categoryResponse, len(items))
	for i, cat := range items {
		out[i] = toCategoryResponse(cat)
	}
	c.JSON(http.StatusOK, categoryListResponse{Items: out})
}

// createCategory creates a category (admin).
//
//	@Summary	Create a category
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	saveCategoryRequest	true	"category"
//	@Success	201	{object}	categoryResponse
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/categories [post]
func (s *Server) createCategory(c *gin.Context) {
	var req saveCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	cat, err := s.categories.Create(c.Request.Context(), toCategoryInput(req))
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, toCategoryResponse(cat))
}

// updateCategory updates a category (admin).
//
//	@Summary	Update a category
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int					true	"category id"
//	@Param		body	body	saveCategoryRequest	true	"category"
//	@Success	200	{object}	categoryResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/categories/{id} [put]
func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	var req saveCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	cat, err := s.categories.Update(c.Request.Context(), id, toCategoryInput(req))
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toCategoryResponse(cat))
}

// deleteCategory deletes a category (admin).
//
//	@Summary	Delete a category
//	@Tags		admin
//	@Security	BearerAuth
//	@Param		id	path	int	true	"category id"
//	@Success	204	"no content"
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/categories/{id} [delete]
func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	if err := s.categories.Delete(c.Request.Context(), id); err != nil {
		s.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
