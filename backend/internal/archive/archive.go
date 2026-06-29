// Package archive converts posts to and from portable Markdown files with YAML
// frontmatter, the on-disk format used by the import/export feature. It holds
// only pure functions (no database, filesystem, or HTTP), so the orchestration
// layer can unit-test the format and the HTTP handler can stay thin.
package archive

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// Frontmatter is the YAML header of an exported post. Field names are the
// project's own (not Hugo/Jekyll); view counts are runtime data and excluded.
type Frontmatter struct {
	Title   string     `yaml:"title"`
	Slug    string     `yaml:"slug"`
	Date    *time.Time `yaml:"date,omitempty"`
	Status  string     `yaml:"status"`
	Pinned  bool       `yaml:"pinned"`
	Tags    []string   `yaml:"tags,omitempty"`
	Summary string     `yaml:"summary,omitempty"`
	Cover   string     `yaml:"cover,omitempty"`
}

// uploadRef matches an uploaded-asset reference, with or without a host:
// "/uploads/<name>" or "http://host/uploads/<name>". The captured group is the
// bare filename, which is a random hex string plus extension.
var uploadRef = regexp.MustCompile(`(?:https?://[^/\s]+)?/uploads/([A-Za-z0-9._-]+)`)

// bundleRef matches a reference rewritten for the export bundle: "images/<name>".
var bundleRef = regexp.MustCompile(`images/([A-Za-z0-9._-]+)`)

// MarshalPost renders frontmatter + body into the bytes of a .md file.
func MarshalPost(fm Frontmatter, body string) ([]byte, error) {
	yml, err := yaml.Marshal(fm)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(yml)
	buf.WriteString("---\n\n")
	buf.WriteString(strings.TrimLeft(body, "\n"))
	if !strings.HasSuffix(body, "\n") {
		buf.WriteString("\n")
	}
	return buf.Bytes(), nil
}

// ParsePost splits a .md file into its frontmatter and body. A file without a
// leading "---" frontmatter block is treated as all body with empty frontmatter.
func ParsePost(raw []byte) (Frontmatter, string, error) {
	s := string(bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})) // strip a UTF-8 BOM if present
	if !strings.HasPrefix(s, "---\n") && !strings.HasPrefix(s, "---\r\n") {
		return Frontmatter{}, s, nil
	}
	// Find the closing fence on its own line.
	rest := s[strings.IndexByte(s, '\n')+1:]
	end := regexp.MustCompile(`(?m)^---[ \t]*\r?\n?`).FindStringIndex(rest)
	if end == nil {
		return Frontmatter{}, "", fmt.Errorf("frontmatter not closed by '---'")
	}
	header := rest[:end[0]]
	body := rest[end[1]:]

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(header), &fm); err != nil {
		return Frontmatter{}, "", fmt.Errorf("invalid frontmatter: %w", err)
	}
	return fm, strings.TrimLeft(body, "\n"), nil
}

// ReferencedImages returns the unique uploaded-asset filenames referenced across
// the given strings (post body, cover, …), e.g. "abc123.png".
func ReferencedImages(sources ...string) []string {
	seen := map[string]bool{}
	var names []string
	for _, src := range sources {
		for _, m := range uploadRef.FindAllStringSubmatch(src, -1) {
			if name := m[1]; !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// ToBundlePaths rewrites "/uploads/<name>" (and absolute variants) to the
// bundle-relative "images/<name>" used inside the export zip.
func ToBundlePaths(s string) string {
	return uploadRef.ReplaceAllString(s, "images/$1")
}

// FromBundlePaths rewrites bundle-relative "images/<name>" back to an absolute
// "<baseURL>/uploads/<name>" on import.
func FromBundlePaths(s, baseURL string) string {
	base := strings.TrimRight(baseURL, "/")
	return bundleRef.ReplaceAllString(s, base+"/uploads/$1")
}
