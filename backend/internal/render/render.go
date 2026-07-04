// Package render converts Markdown source into sanitised HTML. It is the single
// renderer of record: both the publish path and the /api/v1/render preview
// endpoint use it, so the admin preview matches what readers see.
package render

import (
	"bytes"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

var (
	footnoteClassPattern    = regexp.MustCompile(`^(footnotes|footnote-ref|footnote-backref)$`)
	footnoteRolePattern     = regexp.MustCompile(`^doc-(endnotes|noteref|backlink)$`)
	taskCheckboxTypePattern = regexp.MustCompile(`(?i)^checkbox$`)
	checkedAttrPattern      = regexp.MustCompile(`(?i)^(|checked)$`)
	disabledAttrPattern     = regexp.MustCompile(`(?i)^(|disabled)$`)
)

// Renderer renders and sanitises Markdown.
type Renderer struct {
	md     goldmark.Markdown
	policy *bluemonday.Policy
}

// New builds a Renderer with GFM and footnotes enabled, using a UGC sanitisation policy.
func New() *Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote),
		goldmark.WithRendererOptions(gmhtml.WithXHTML(), gmhtml.WithHardWraps()),
	)
	return &Renderer{md: md, policy: newPolicy()}
}

// ToHTML converts Markdown to sanitised HTML.
func (r *Renderer) ToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return r.policy.Sanitize(buf.String()), nil
}

func newPolicy() *bluemonday.Policy {
	policy := bluemonday.UGCPolicy()

	// UGCPolicy allows footnote ids and fragment links, but not class or ARIA role.
	policy.AllowAttrs("class").Matching(footnoteClassPattern).OnElements("div", "a")
	policy.AllowAttrs("role").Matching(footnoteRolePattern).OnElements("div", "a")

	// GFM task lists render disabled checkbox inputs; keep only that inert shape.
	policy.AllowAttrs("type").Matching(taskCheckboxTypePattern).OnElements("input")
	policy.AllowAttrs("checked").Matching(checkedAttrPattern).OnElements("input")
	policy.AllowAttrs("disabled").Matching(disabledAttrPattern).OnElements("input")

	return policy
}
