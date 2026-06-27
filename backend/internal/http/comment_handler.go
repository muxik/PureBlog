package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/internal/service"
)

// ---- DTOs ----

type createCommentRequest struct {
	AuthorName  string `json:"authorName" binding:"required"`
	AuthorEmail string `json:"authorEmail"`
	Content     string `json:"content" binding:"required"`
	ParentID    *int64 `json:"parentId"`
}

type moderateCommentRequest struct {
	Status string `json:"status" binding:"required,oneof=approved pending"`
}

type commentResponse struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"postId"`
	ParentID   *int64    `json:"parentId"`
	AuthorName string    `json:"authorName"`
	Content    string    `json:"content"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

func toCommentResponse(c *domain.Comment) commentResponse {
	return commentResponse{
		ID:         c.ID,
		PostID:     c.PostID,
		ParentID:   c.ParentID,
		AuthorName: c.AuthorName,
		Content:    c.Content,
		Status:     c.Status,
		CreatedAt:  c.CreatedAt,
	}
}

type commentListResponse struct {
	Items []commentResponse `json:"items"`
}

func toCommentResponses(items []*domain.Comment) []commentResponse {
	out := make([]commentResponse, len(items))
	for i, c := range items {
		out[i] = toCommentResponse(c)
	}
	return out
}

// ---- handlers ----

// listPostComments lists approved comments for a published post.
//
//	@Summary	List approved comments for a post
//	@Tags		comments
//	@Produce	json
//	@Param		slug	path	string	true	"post slug"
//	@Success	200	{object}	commentListResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/posts/{slug}/comments [get]
func (s *Server) listPostComments(c *gin.Context) {
	items, err := s.comments.ListApprovedForSlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, commentListResponse{Items: toCommentResponses(items)})
}

// createPostComment submits a pending comment on a published post.
//
//	@Summary	Submit a comment on a post
//	@Tags		comments
//	@Accept		json
//	@Produce	json
//	@Param		slug	path	string					true	"post slug"
//	@Param		body	body	createCommentRequest	true	"comment"
//	@Success	201	{object}	commentResponse
//	@Failure	400	{object}	errorResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/posts/{slug}/comments [post]
func (s *Server) createPostComment(c *gin.Context) {
	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	in := service.CreateCommentInput{
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		Content:     req.Content,
		ParentID:    req.ParentID,
		IP:          c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}
	cm, err := s.comments.CreateForSlug(c.Request.Context(), c.Param("slug"), in)
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, toCommentResponse(cm))
}

// adminListComments lists comments, optionally filtered by status (admin).
//
//	@Summary	List comments (admin)
//	@Tags		admin
//	@Security	BearerAuth
//	@Produce	json
//	@Param		status	query	string	false	"pending|approved (default: all)"
//	@Success	200	{object}	commentListResponse
//	@Router		/admin/comments [get]
func (s *Server) adminListComments(c *gin.Context) {
	items, err := s.comments.List(c.Request.Context(), c.Query("status"))
	if err != nil {
		s.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, commentListResponse{Items: toCommentResponses(items)})
}

// moderateComment sets the moderation status of a comment (admin).
//
//	@Summary	Moderate a comment
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int						true	"comment id"
//	@Param		body	body	moderateCommentRequest	true	"status"
//	@Success	204	"no content"
//	@Failure	400	{object}	errorResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/comments/{id} [put]
func (s *Server) moderateComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	var req moderateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	if err := s.comments.SetStatus(c.Request.Context(), id, req.Status); err != nil {
		s.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// deleteComment deletes a comment (admin).
//
//	@Summary	Delete a comment
//	@Tags		admin
//	@Security	BearerAuth
//	@Param		id	path	int	true	"comment id"
//	@Success	204	"no content"
//	@Failure	404	{object}	errorResponse
//	@Router		/admin/comments/{id} [delete]
func (s *Server) deleteComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.badRequest(c, err)
		return
	}
	if err := s.comments.Delete(c.Request.Context(), id); err != nil {
		s.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
