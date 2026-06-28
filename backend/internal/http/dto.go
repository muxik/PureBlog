// Package http is the Gin transport layer: routing, request/response DTOs, and
// middleware. It depends on the service layer, never the other way around.
package http

import (
	"strconv"
	"time"

	"github.com/muxik/PureBlog/backend/internal/domain"
)

// ---- requests ----

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type renderRequest struct {
	Markdown string `json:"markdown"`
}

type savePostRequest struct {
	Title     string  `json:"title" binding:"required"`
	Slug      string  `json:"slug"`
	Summary   string  `json:"summary"`
	ContentMD string  `json:"contentMd"`
	CoverURL  string  `json:"coverUrl"`
	Status    string  `json:"status"`
	Pinned    bool    `json:"pinned"`
	TagIDs    []int64 `json:"tagIds"`
}

// ---- responses ----

type renderResponse struct {
	HTML string `json:"html"`
}

type userResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
}

type tokenResponse struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	User         *userResponse `json:"user,omitempty"`
}

type postResponse struct {
	ID          int64         `json:"id"`
	Slug        string        `json:"slug"`
	Title       string        `json:"title"`
	Summary     string        `json:"summary"`
	ContentMD   string        `json:"contentMd,omitempty"`
	ContentHTML string        `json:"contentHtml,omitempty"`
	CoverURL    string        `json:"coverUrl"`
	Status      string        `json:"status"`
	Pinned      bool          `json:"pinned"`
	ViewCount   int64         `json:"viewCount"`
	Tags        []tagResponse `json:"tags"`
	PublishedAt *time.Time    `json:"publishedAt"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type listResponse struct {
	Items    []postResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// ---- mappers ----

func toUserResponse(u *domain.User) *userResponse {
	return &userResponse{ID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role}
}

func toPostResponse(p *domain.Post, includeContent bool) postResponse {
	resp := postResponse{
		ID:          p.ID,
		Slug:        p.Slug,
		Title:       p.Title,
		Summary:     p.Summary,
		CoverURL:    p.CoverURL,
		Status:      string(p.Status),
		Pinned:      p.Pinned,
		ViewCount:   p.ViewCount,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	resp.Tags = make([]tagResponse, len(p.Tags))
	for i := range p.Tags {
		resp.Tags[i] = toTagResponse(&p.Tags[i])
	}
	if includeContent {
		resp.ContentMD = p.ContentMD
		resp.ContentHTML = p.ContentHTML
	}
	return resp
}

func toListResponse(items []*domain.Post, total int64, page, size int, includeContent bool) listResponse {
	out := make([]postResponse, len(items))
	for i, p := range items {
		out[i] = toPostResponse(p, includeContent)
	}
	return listResponse{Items: out, Total: total, Page: page, PageSize: size}
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
