package render

import (
	"strings"
	"testing"
)

func TestGFMTableAndTaskListRender(t *testing.T) {
	html := renderMarkdown(t, `| Done | Item |
| --- | --- |
| yes | table |

- [x] checked
- [ ] unchecked
`)

	assertContains(t, html, "<table>")
	assertContains(t, html, "<thead>")
	assertContains(t, html, "<tbody>")
	assertContains(t, html, "<th>Done</th>")
	assertContains(t, html, "<td>table</td>")

	if got := strings.Count(html, "<input"); got != 2 {
		t.Fatalf("expected two task checkboxes, got %d in:\n%s", got, html)
	}
	assertContains(t, html, `type="checkbox"`)
	assertContains(t, html, `disabled=""`)
	assertContains(t, html, `checked=""`)
}

func TestFootnotesRenderUsableHTML(t *testing.T) {
	html := renderMarkdown(t, `Here is a footnote.[^note]

[^note]: Footnote body.
`)

	assertContains(t, html, `id="fnref:1"`)
	assertContains(t, html, `href="#fn:1"`)
	assertContains(t, html, `class="footnote-ref"`)
	assertContains(t, html, `role="doc-noteref"`)
	assertContains(t, html, `class="footnotes"`)
	assertContains(t, html, `role="doc-endnotes"`)
	assertContains(t, html, `id="fn:1"`)
	assertContains(t, html, `href="#fnref:1"`)
	assertContains(t, html, `class="footnote-backref"`)
	assertContains(t, html, `role="doc-backlink"`)
	assertContains(t, html, "Footnote body.")
}

func TestSoftLineBreaksRenderAsVisibleBreaks(t *testing.T) {
	html := renderMarkdown(t, `你好
你好
`)

	assertContains(t, html, "你好<br")
	assertContains(t, html, "<br")
}

func TestRawHTMLIsNotRenderedAsDangerousHTML(t *testing.T) {
	html := renderMarkdown(t, `Before

<script>alert(1)</script>
<img src="x" onerror="alert(1)">
<span onclick="alert(1)">click</span>

After
`)

	for _, forbidden := range []string{
		"<script",
		"<img",
		"onerror",
		"onclick",
		"alert(1)",
	} {
		assertNotContains(t, html, forbidden)
	}
	assertContains(t, html, "Before")
	assertContains(t, html, "After")
}

func renderMarkdown(t *testing.T, markdown string) string {
	t.Helper()

	html, err := New().ToHTML(markdown)
	if err != nil {
		t.Fatalf("render markdown: %v", err)
	}
	return html
}

func assertContains(t *testing.T, html, want string) {
	t.Helper()
	if !strings.Contains(html, want) {
		t.Fatalf("expected HTML to contain %q, got:\n%s", want, html)
	}
}

func assertNotContains(t *testing.T, html, want string) {
	t.Helper()
	if strings.Contains(html, want) {
		t.Fatalf("expected HTML not to contain %q, got:\n%s", want, html)
	}
}
