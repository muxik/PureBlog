package archive

import (
	"strings"
	"testing"
	"time"
)

func TestMarshalParseRoundTrip(t *testing.T) {
	when := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	fm := Frontmatter{
		Title:   "农历这件小事",
		Slug:    "lunar-calendar",
		Date:    &when,
		Status:  "published",
		Pinned:  true,
		Tags:    []string{"随笔", "排版"},
		Summary: "一句话摘要",
		Cover:   "images/abc.png",
	}
	body := "正文第一段。\n\n第二段。"

	raw, err := MarshalPost(fm, body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.HasPrefix(string(raw), "---\n") {
		t.Fatalf("expected leading frontmatter fence, got:\n%s", raw)
	}

	got, gotBody, err := ParsePost(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.Title != fm.Title || got.Slug != fm.Slug || got.Status != fm.Status || got.Pinned != fm.Pinned {
		t.Fatalf("scalar fields differ: %+v", got)
	}
	if got.Date == nil || !got.Date.Equal(when) {
		t.Fatalf("date not round-tripped: %v", got.Date)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "随笔" || got.Tags[1] != "排版" {
		t.Fatalf("tags not round-tripped: %v", got.Tags)
	}
	if strings.TrimSpace(gotBody) != body {
		t.Fatalf("body not round-tripped: %q", gotBody)
	}
}

func TestParseNoFrontmatter(t *testing.T) {
	fm, body, err := ParsePost([]byte("# Just markdown\n\nno header."))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if fm.Title != "" || !strings.HasPrefix(body, "# Just markdown") {
		t.Fatalf("expected empty frontmatter and full body, got fm=%+v body=%q", fm, body)
	}
}

func TestParseUnclosedFrontmatter(t *testing.T) {
	if _, _, err := ParsePost([]byte("---\ntitle: x\nno closing fence")); err == nil {
		t.Fatal("expected an error for unclosed frontmatter")
	}
}

func TestReferencedImages(t *testing.T) {
	body := "![a](/uploads/aaa.png) text ![b](http://localhost/uploads/bbb.jpg) and again ![a2](/uploads/aaa.png)"
	cover := "http://localhost/uploads/ccc.webp"
	got := ReferencedImages(body, cover)
	want := map[string]bool{"aaa.png": true, "bbb.jpg": true, "ccc.webp": true}
	if len(got) != 3 {
		t.Fatalf("expected 3 unique images, got %v", got)
	}
	for _, n := range got {
		if !want[n] {
			t.Fatalf("unexpected image %q in %v", n, got)
		}
	}
}

func TestBundlePathRewriteRoundTrip(t *testing.T) {
	orig := "![x](/uploads/aaa.png) and ![y](http://localhost/uploads/bbb.jpg)"
	bundled := ToBundlePaths(orig)
	if strings.Contains(bundled, "/uploads/") || !strings.Contains(bundled, "images/aaa.png") || !strings.Contains(bundled, "images/bbb.jpg") {
		t.Fatalf("ToBundlePaths did not rewrite: %q", bundled)
	}
	restored := FromBundlePaths(bundled, "https://blog.example.com/")
	if !strings.Contains(restored, "https://blog.example.com/uploads/aaa.png") ||
		!strings.Contains(restored, "https://blog.example.com/uploads/bbb.jpg") {
		t.Fatalf("FromBundlePaths did not restore absolute URLs: %q", restored)
	}
}

func TestFromBundlePathsDoesNotClobberForeignRefs(t *testing.T) {
	base := "https://blog.example.com"
	cases := []struct {
		name string
		in   string
		want string // unchanged
	}{
		{"external image url", "![cat](https://cdn.other.com/images/cat.png)", "![cat](https://cdn.other.com/images/cat.png)"},
		{"prose mention", "see the images/readme file for details", "see the images/readme file for details"},
		{"foreign path", "path is /var/images/logo.svg here", "path is /var/images/logo.svg here"},
	}
	for _, c := range cases {
		if got := FromBundlePaths(c.in, base); got != c.want {
			t.Errorf("%s: FromBundlePaths clobbered a non-bundle reference\n in:   %q\n got:  %q\n want: %q", c.name, c.in, got, c.want)
		}
	}

	// A genuine bundle reference (Markdown image, and a bare cover) IS rewritten.
	if got := FromBundlePaths("![x](images/aaa.png)", base); got != "![x]("+base+"/uploads/aaa.png)" {
		t.Errorf("markdown bundle ref not rewritten: %q", got)
	}
	if got := FromBundlePaths("images/cover.png", base); got != base+"/uploads/cover.png" {
		t.Errorf("bare cover bundle ref not rewritten: %q", got)
	}
}
