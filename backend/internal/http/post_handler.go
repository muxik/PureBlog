package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/internal/http/middleware"
	"github.com/muxik/PureBlog/backend/internal/service"
)

// listPosts lists published posts.
//
//	@Summary	List published posts
//	@Tags		posts
//	@Produce	json
//	@Param		page		query	int		false	"page (1-based)"
//	@Param		pageSize	query	int		false	"page size"
//	@Param		q			query	string	false	"search query"
//	@Success	200	{object}	listResponse
//	@Router		/posts [get]
func (s *Server) listPosts(c *gin.Context) {
	f := domain.PostListFilter{
		Status:   domain.StatusPublished,
		Query:    c.Query("q"),
		Page:     atoiDefault(c.Query("page"), 1),
		PageSize: atoiDefault(c.Query("pageSize"), 10),
	}
	items, total, err := s.posts.List(c.Request.Context(), f)
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toListResponse(items, total, f.Page, f.PageSize, false))
}

// getPostBySlug returns a single published post by slug.
//
//	@Summary	Get a published post by slug
//	@Tags		posts
//	@Produce	json
//	@Param		slug	path	string	true	"post slug"
//	@Success	200	{object}	postResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/posts/{slug} [get]
func (s *Server) getPostBySlug(c *gin.Context) {
	p, err := s.posts.GetPublishedBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toPostResponse(p, true))
}

// adminListPosts lists posts of any status (admin).
//
//	@Summary	List posts (admin)
//	@Tags		admin
//	@Security	BearerAuth
//	@Produce	json
//	@Param		status		query	string	false	"draft|published (default: any)"
//	@Param		page		query	int		false	"page (1-based)"
//	@Param		pageSize	query	int		false	"page size"
//	@Param		q			query	string	false	"search query"
//	@Success	200	{object}	listResponse
//	@Router		/admin/posts [get]
func (s *Server) adminListPosts(c *gin.Context) {
	f := domain.PostListFilter{
		Status:   domain.PostStatus(c.Query("status")),
		Query:    c.Query("q"),
		Page:     atoiDefault(c.Query("page"), 1),
		PageSize: atoiDefault(c.Query("pageSize"), 20),
	}
	items, total, err := s.posts.List(c.Request.Context(), f)
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toListResponse(items, total, f.Page, f.PageSize, false))
}

// adminGetPost returns any post by id (admin).
//
//	@Summary	Get a post by id (admin)
//	@Tags		admin
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path	int	true	"post id"
//	@Success	200	{object}	postResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/posts/{id} [get]
func (s *Server) adminGetPost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	p, err := s.posts.GetByID(c.Request.Context(), id)
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toPostResponse(p, true))
}

// createPost creates a post (admin).
//
//	@Summary	Create a post
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	savePostRequest	true	"post"
//	@Success	201	{object}	postResponse
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/posts [post]
func (s *Server) createPost(c *gin.Context) {
	var req savePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	p, err := s.posts.Create(c.Request.Context(), s.toInput(req, middleware.UserID(c)))
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, toPostResponse(p, true))
}

// updatePost updates a post (admin).
//
//	@Summary	Update a post
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int				true	"post id"
//	@Param		body	body	savePostRequest	true	"post"
//	@Success	200	{object}	postResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/posts/{id} [put]
func (s *Server) updatePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	var req savePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	p, err := s.posts.Update(c.Request.Context(), id, s.toInput(req, middleware.UserID(c)))
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toPostResponse(p, true))
}

// deletePost deletes a post (admin).
//
//	@Summary	Delete a post
//	@Tags		admin
//	@Security	BearerAuth
//	@Param		id	path	int	true	"post id"
//	@Success	204	"no content"
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/posts/{id} [delete]
func (s *Server) deletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	if err := s.posts.Delete(c.Request.Context(), id); err != nil {
		s.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// renderMarkdown renders Markdown to sanitised HTML for the admin live preview.
//
//	@Summary	Render Markdown preview
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	renderRequest	true	"markdown"
//	@Success	200	{object}	renderResponse
//	@Router		/admin/render [post]
func (s *Server) renderMarkdown(c *gin.Context) {
	var req renderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	html, err := s.posts.Render(req.Markdown)
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, renderResponse{HTML: html})
}

func (s *Server) toInput(req savePostRequest, authorID int64) service.SavePostInput {
	return service.SavePostInput{
		Title:     req.Title,
		Slug:      req.Slug,
		Summary:   req.Summary,
		ContentMD: req.ContentMD,
		CoverURL:  req.CoverURL,
		Status:    domain.PostStatus(req.Status),
		Pinned:    req.Pinned,
		AuthorID:  authorID,
	}
}
