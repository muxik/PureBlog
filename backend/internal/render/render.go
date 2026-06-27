// Package render converts Markdown source into sanitised HTML. It is the single
// renderer of record: both the publish path and the /api/v1/render preview
// endpoint use it, so the admin preview matches what readers see.
package render

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// Renderer renders and sanitises Markdown.
type Renderer struct {
	md     goldmark.Markdown
	policy *bluemonday.Policy
}

// New builds a Renderer with GFM enabled and a UGC sanitisation policy.
func New() *Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(gmhtml.WithXHTML()),
	)
	return &Renderer{md: md, policy: bluemonday.UGCPolicy()}
}

// ToHTML converts Markdown to sanitised HTML.
func (r *Renderer) ToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return r.policy.Sanitize(buf.String()), nil
}
