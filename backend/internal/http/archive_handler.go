package http

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muxik/PureBlog/backend/internal/archive"
	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/internal/http/middleware"
	"github.com/muxik/PureBlog/backend/internal/service"
)

// maxImportBytes caps an uploaded import archive at 64 MiB.
const maxImportBytes = 64 << 20

// importSummary reports the outcome of an import to the client.
type importSummary struct {
	Created int           `json:"created"`
	Updated int           `json:"updated"`
	Images  int           `json:"images"`
	Errors  []importError `json:"errors"`
}

type importError struct {
	File   string `json:"file"`
	Reason string `json:"reason"`
}

// exportPosts streams a .zip of posts as Markdown files with YAML frontmatter,
// bundling every referenced uploaded image under images/ and rewriting the
// references to be bundle-relative.
//
//	@Summary	Export posts as a Markdown bundle
//	@Tags		admin
//	@Security	BearerAuth
//	@Produce	application/zip
//	@Param		slugs	query	string	false	"comma-separated slugs; default all"
//	@Success	200	{file}	binary
//	@Router		/admin/export [get]
func (s *Server) exportPosts(c *gin.Context) {
	ctx := c.Request.Context()
	// Page through every post (List caps page size at 100).
	var posts []*domain.Post
	for page := 1; ; page++ {
		batch, total, err := s.posts.List(ctx, domain.PostListFilter{Page: page, PageSize: 100})
		if err != nil {
			s.fail(c, err)
			return
		}
		posts = append(posts, batch...)
		if len(batch) == 0 || int64(len(posts)) >= total {
			break
		}
	}

	// Optional slug allow-list.
	if raw := strings.TrimSpace(c.Query("slugs")); raw != "" {
		want := map[string]bool{}
		for _, sl := range strings.Split(raw, ",") {
			if sl = strings.TrimSpace(sl); sl != "" {
				want[sl] = true
			}
		}
		filtered := posts[:0]
		for _, p := range posts {
			if want[p.Slug] {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="pureblog-export.zip"`)

	zw := zip.NewWriter(c.Writer)
	defer zw.Close()

	imageNames := map[string]bool{}
	for _, p := range posts {
		for _, name := range archive.ReferencedImages(p.ContentMD, p.CoverURL) {
			imageNames[name] = true
		}

		tagNames := make([]string, 0, len(p.Tags))
		for _, t := range p.Tags {
			tagNames = append(tagNames, t.Name)
		}
		fm := archive.Frontmatter{
			Title:   p.Title,
			Slug:    p.Slug,
			Date:    p.PublishedAt,
			Status:  string(p.Status),
			Pinned:  p.Pinned,
			Tags:    tagNames,
			Summary: p.Summary,
			Cover:   archive.ToBundlePaths(p.CoverURL),
		}
		md, err := archive.MarshalPost(fm, archive.ToBundlePaths(p.ContentMD))
		if err != nil {
			continue
		}
		w, err := zw.Create("posts/" + p.Slug + ".md")
		if err != nil {
			return
		}
		_, _ = w.Write(md)
	}

	// Bundle the referenced images (sorted for a deterministic archive).
	names := make([]string, 0, len(imageNames))
	for n := range imageNames {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(s.cfg.UploadDir, filepath.Base(name)))
		if err != nil {
			continue // a missing/external image is skipped, not fatal
		}
		w, err := zw.Create("images/" + name)
		if err != nil {
			return
		}
		_, _ = w.Write(data)
	}
}

// importPosts ingests a .zip produced by export (or hand-authored in the same
// layout): images/* are restored to the upload directory, posts/*.md are parsed
// and upserted by slug (existing slugs are overwritten).
//
//	@Summary	Import posts from a Markdown bundle
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		file	formData	file	true	"zip archive"
//	@Success	200	{object}	importSummary
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/import [post]
func (s *Server) importPosts(c *gin.Context) {
	ctx := c.Request.Context()
	authorID := middleware.UserID(c)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxImportBytes+1024)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		s.badRequest(c, fmt.Errorf("missing file field: %w", err))
		return
	}
	src, err := fileHeader.Open()
	if err != nil {
		s.fail(c, err)
		return
	}
	defer src.Close()
	raw, err := io.ReadAll(src)
	if err != nil {
		s.fail(c, err)
		return
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		s.badRequest(c, fmt.Errorf("not a valid zip archive: %w", err))
		return
	}

	summary := importSummary{Errors: []importError{}}

	// Pass 1: restore images so the body rewrite can point at real files.
	if err := os.MkdirAll(s.cfg.UploadDir, 0o755); err != nil {
		s.fail(c, err)
		return
	}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() || !strings.HasPrefix(f.Name, "images/") {
			continue
		}
		name := filepath.Base(path.Clean(f.Name))
		if name == "." || name == ".." || name == "" {
			continue
		}
		if err := writeZipEntry(f, filepath.Join(s.cfg.UploadDir, name)); err != nil {
			summary.Errors = append(summary.Errors, importError{File: f.Name, Reason: err.Error()})
			continue
		}
		summary.Images++
	}

	// Build the tag name→id cache once, up front, so the per-post resolution
	// doesn't re-query the whole tags table for every file (N+1).
	tagCache, err := s.loadTagCache(ctx)
	if err != nil {
		s.fail(c, err)
		return
	}

	// Pass 2: upsert posts.
	for _, f := range zr.File {
		if f.FileInfo().IsDir() || !strings.HasPrefix(f.Name, "posts/") || !strings.HasSuffix(f.Name, ".md") {
			continue
		}
		if err := s.importOnePost(ctx, f, authorID, tagCache, &summary); err != nil {
			summary.Errors = append(summary.Errors, importError{File: f.Name, Reason: err.Error()})
		}
	}

	c.JSON(http.StatusOK, summary)
}

// importOnePost parses and upserts a single posts/*.md entry.
func (s *Server) importOnePost(ctx context.Context, f *zip.File, authorID int64, tagCache map[string]int64, summary *importSummary) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	fm, body, err := archive.ParsePost(data)
	if err != nil {
		return err
	}

	slug := strings.TrimSpace(fm.Slug)
	if slug == "" {
		// Fall back to the filename without extension.
		slug = strings.TrimSuffix(filepath.Base(f.Name), ".md")
	}
	if strings.TrimSpace(fm.Title) == "" {
		return fmt.Errorf("missing title")
	}

	// Look up the existing post first. Only a genuine "not found" means we
	// should create; any other error (transient DB failure, cancellation) must
	// abort, otherwise we'd fall through to Create and produce a duplicate
	// "<slug>-2" instead of overwriting.
	existing, err := s.posts.FindBySlug(ctx, slug)
	switch {
	case err == nil:
		// found → update
	case errors.Is(err, domain.ErrNotFound):
		existing = nil
	default:
		return err
	}

	tagIDs, err := s.resolveTagIDsByName(ctx, tagCache, fm.Tags)
	if err != nil {
		return err
	}

	in := service.SavePostInput{
		Title:       fm.Title,
		Slug:        slug,
		Summary:     fm.Summary,
		ContentMD:   archive.FromBundlePaths(body, s.cfg.PublicBaseURL),
		CoverURL:    archive.FromBundlePaths(fm.Cover, s.cfg.PublicBaseURL),
		Status:      importStatus(fm.Status, existing),
		Pinned:      fm.Pinned,
		AuthorID:    authorID,
		TagIDs:      tagIDs,
		PublishedAt: fm.Date,
	}

	if existing != nil {
		if _, err := s.posts.Update(ctx, existing.ID, in); err != nil {
			return err
		}
		summary.Updated++
		return nil
	}
	if _, err := s.posts.Create(ctx, in); err != nil {
		return err
	}
	summary.Created++
	return nil
}

// importStatus resolves the status for an imported post. An explicit
// "published"/"draft" is honored; anything else (empty or unrecognised) does
// NOT silently unpublish — it preserves the existing post's status on update,
// or defaults to published for a new post (export always writes an explicit
// status, so an omitted one only occurs in hand-authored bundles).
func importStatus(raw string, existing *domain.Post) domain.PostStatus {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(domain.StatusPublished):
		return domain.StatusPublished
	case string(domain.StatusDraft):
		return domain.StatusDraft
	default:
		if existing != nil {
			return existing.Status
		}
		return domain.StatusPublished
	}
}

// loadTagCache reads all tags once into a name→id map for the import to reuse,
// avoiding a full tags-table query per imported post.
func (s *Server) loadTagCache(ctx context.Context) (map[string]int64, error) {
	existing, err := s.tags.List(ctx)
	if err != nil {
		return nil, err
	}
	byName := make(map[string]int64, len(existing))
	for _, t := range existing {
		byName[t.Name] = t.ID
	}
	return byName, nil
}

// resolveTagIDsByName maps tag names to ids using the shared cache, creating any
// that don't yet exist (and adding them to the cache so later posts reuse them).
func (s *Server) resolveTagIDsByName(ctx context.Context, cache map[string]int64, names []string) ([]int64, error) {
	if len(names) == 0 {
		return nil, nil
	}
	var ids []int64
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		if id, ok := cache[n]; ok {
			ids = append(ids, id)
			continue
		}
		created, err := s.tags.Create(ctx, service.SaveTagInput{Name: n})
		if err != nil {
			return nil, err
		}
		cache[n] = created.ID
		ids = append(ids, created.ID)
	}
	return ids, nil
}

// writeZipEntry copies a zip file entry to dst, capping the DECOMPRESSED size at
// maxUploadBytes. The request-body limit only bounds the compressed upload, so
// without this a small archive with a high-ratio deflate entry (a zip bomb)
// could expand without limit and fill the disk. A partially written file is
// removed on any failure.
func writeZipEntry(f *zip.File, dst string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	// LimitReader to one byte past the cap so we can detect an oversize entry.
	written, err := io.Copy(out, io.LimitReader(rc, maxUploadBytes+1))
	if err == nil && written > maxUploadBytes {
		err = fmt.Errorf("image exceeds %d MiB limit", maxUploadBytes>>20)
	}
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(dst)
	}
	return err
}
