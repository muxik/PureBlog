package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/domain"
)

// batchRequest is the payload for a bulk action over posts.
type batchRequest struct {
	// Action is one of: "delete", "publish", "draft".
	Action string  `json:"action"`
	IDs    []int64 `json:"ids"`
}

// batchResponse reports how many posts the action affected and which ids failed.
type batchResponse struct {
	Affected int          `json:"affected"`
	Failed   []batchError `json:"failed"`
}

type batchError struct {
	ID     int64  `json:"id"`
	Reason string `json:"reason"`
}

// batchPosts applies a bulk action (delete / publish / draft) to the given post
// ids. It is non-transactional, so rather than aborting on the first error and
// discarding what already succeeded, it applies every id and reports the exact
// affected count and per-id failures — the client can reconcile instead of
// blindly retrying a partially-applied operation.
//
//	@Summary	Bulk action over posts
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	batchRequest	true	"action and ids"
//	@Success	200	{object}	batchResponse
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/posts/batch [post]
func (s *Server) batchPosts(c *gin.Context) {
	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.badRequest(c, err)
		return
	}
	ctx := c.Request.Context()

	var apply func(id int64) error
	switch req.Action {
	case "delete":
		apply = func(id int64) error { return s.posts.Delete(ctx, id) }
	case "publish":
		apply = func(id int64) error { return s.posts.SetStatus(ctx, id, domain.StatusPublished) }
	case "draft":
		apply = func(id int64) error { return s.posts.SetStatus(ctx, id, domain.StatusDraft) }
	default:
		s.badRequest(c, fmt.Errorf("unknown action %q (want delete|publish|draft)", req.Action))
		return
	}

	resp := batchResponse{Failed: []batchError{}}
	for _, id := range req.IDs {
		if err := apply(id); err != nil {
			resp.Failed = append(resp.Failed, batchError{ID: id, Reason: err.Error()})
			continue
		}
		resp.Affected++
	}
	c.JSON(http.StatusOK, resp)
}
