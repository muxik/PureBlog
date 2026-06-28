package http

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// maxUploadBytes caps an uploaded image at 10 MiB.
const maxUploadBytes = 10 << 20

// imageExtByType maps the content types we accept onto the extension we store
// the file with. Anything not in this allowlist is rejected.
var imageExtByType = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// uploadResponse is returned after a successful image upload.
type uploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// mediaItem describes one stored upload for the media library.
type mediaItem struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	// Modified is an RFC3339 timestamp of the file's last modification.
	Modified string `json:"modified"`
}

// uploadImage accepts a multipart "file" field, validates that it is an image
// within the size limit, stores it under the upload directory with a random
// name, and returns its public URL.
//
//	@Summary	Upload an image
//	@Tags		admin
//	@Security	BearerAuth
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		file	formData	file	true	"image file"
//	@Success	201	{object}	uploadResponse
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/uploads [post]
func (s *Server) uploadImage(c *gin.Context) {
	// Reject oversize bodies before reading them fully into memory/disk.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes+1024)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		s.badRequest(c, fmt.Errorf("missing file field: %w", err))
		return
	}
	if fileHeader.Size > maxUploadBytes {
		s.badRequest(c, fmt.Errorf("file exceeds %d MiB limit", maxUploadBytes>>20))
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		s.fail(c, err)
		return
	}
	defer src.Close()

	// Sniff the real content type from the first 512 bytes rather than trusting
	// the client-supplied Content-Type or the filename extension.
	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	contentType := http.DetectContentType(head[:n])
	ext, ok := imageExtByType[contentType]
	if !ok {
		s.badRequest(c, fmt.Errorf("unsupported image type %q (allowed: png, jpeg, gif, webp)", contentType))
		return
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		s.fail(c, err)
		return
	}

	if err := os.MkdirAll(s.cfg.UploadDir, 0o755); err != nil {
		s.fail(c, err)
		return
	}
	name, err := randomName(ext)
	if err != nil {
		s.fail(c, err)
		return
	}

	dst, err := os.Create(filepath.Join(s.cfg.UploadDir, name))
	if err != nil {
		s.fail(c, err)
		return
	}
	written, err := io.Copy(dst, src)
	if cerr := dst.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		// Clean up a partially written file on failure.
		_ = os.Remove(filepath.Join(s.cfg.UploadDir, name))
		s.fail(c, err)
		return
	}

	c.JSON(http.StatusCreated, uploadResponse{
		URL:      s.cfg.PublicBaseURL + "/uploads/" + name,
		Filename: name,
		Size:     written,
	})
}

// listUploads returns the stored uploads, newest first, for the media library.
//
//	@Summary	List uploaded images
//	@Tags		admin
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{array}	mediaItem
//	@Router		/admin/uploads [get]
func (s *Server) listUploads(c *gin.Context) {
	entries, err := os.ReadDir(s.cfg.UploadDir)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, []mediaItem{})
			return
		}
		s.fail(c, err)
		return
	}

	items := make([]mediaItem, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		items = append(items, mediaItem{
			URL:      s.cfg.PublicBaseURL + "/uploads/" + e.Name(),
			Filename: e.Name(),
			Size:     info.Size(),
			Modified: info.ModTime().UTC().Format("2006-01-02T15:04:05Z"),
		})
	}
	// Newest first.
	sort.Slice(items, func(i, j int) bool { return items[i].Modified > items[j].Modified })

	c.JSON(http.StatusOK, items)
}

// deleteUpload removes a stored upload by filename.
//
//	@Summary	Delete an uploaded image
//	@Tags		admin
//	@Security	BearerAuth
//	@Param		name	path	string	true	"stored filename"
//	@Success	204	"No Content"
//	@Failure	400	{object}	errorResponse
//	@Router		/admin/uploads/{name} [delete]
func (s *Server) deleteUpload(c *gin.Context) {
	name := c.Param("name")
	// Reject path traversal: the name must be a bare filename, not a path or a
	// relative directory reference.
	if name == "" || name == "." || name == ".." ||
		name != filepath.Base(name) || strings.ContainsAny(name, `/\`) {
		s.badRequest(c, fmt.Errorf("invalid filename"))
		return
	}
	if err := os.Remove(filepath.Join(s.cfg.UploadDir, name)); err != nil {
		if os.IsNotExist(err) {
			c.Status(http.StatusNoContent)
			return
		}
		s.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// randomName builds a collision-resistant filename: 16 random bytes hex-encoded
// plus the given extension.
func randomName(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + ext, nil
}
