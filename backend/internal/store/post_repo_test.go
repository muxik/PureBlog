package store

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/muxik/PureBlog/backend/internal/domain"
	"github.com/muxik/PureBlog/backend/migrations"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// These are integration tests: they run against a real PostgreSQL instance
// pointed to by TEST_DATABASE_URL. When that variable is unset the tests skip,
// so `go test ./...` stays green on machines without a database. CI provides a
// Postgres service container and sets TEST_DATABASE_URL.

// testDB opens the integration database, applies all migrations, and truncates
// every table so each test starts from a clean, deterministic slate.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping Postgres integration test")
	}
	db, err := Open(dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("acquire *sql.DB: %v", err)
	}
	if err := migrateUp(sqlDB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := db.Exec("TRUNCATE post_tags, posts, tags, users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return db
}

func migrateUp(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}

// seedUser inserts a user and returns its id, for use as a post's author.
func seedUser(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	u := &domain.User{Username: "tester", PasswordHash: "x", Role: "admin"}
	if err := NewUserRepo(db).Create(context.Background(), u); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u.ID
}

// seedTag inserts a tag and returns its id.
func seedTag(t *testing.T, db *gorm.DB, name, slug string) int64 {
	t.Helper()
	tag := &domain.Tag{Name: name, Slug: slug}
	if err := NewTagRepo(db).Create(context.Background(), tag); err != nil {
		t.Fatalf("seed tag: %v", err)
	}
	return tag.ID
}

func TestPostRepo_CreateGetUpdateDelete(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	repo := NewPostRepo(db)
	author := seedUser(t, db)
	tagID := seedTag(t, db, "Go", "go")

	now := time.Now()
	p := &domain.Post{
		Slug:        "hello-world",
		Title:       "Hello World",
		Summary:     "first post",
		ContentMD:   "# hi",
		ContentHTML: "<h1>hi</h1>",
		Status:      domain.StatusPublished,
		AuthorID:    author,
		PublishedAt: &now,
		TagIDs:      []int64{tagID},
	}
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.ID == 0 {
		t.Fatal("create did not assign an id")
	}
	if len(p.Tags) != 1 || p.Tags[0].Slug != "go" {
		t.Fatalf("create did not attach tags: %+v", p.Tags)
	}

	// GetByID round-trips the persisted fields and tags.
	got, err := repo.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Title != "Hello World" || got.Slug != "hello-world" {
		t.Fatalf("unexpected post: %+v", got)
	}
	if len(got.Tags) != 1 || got.Tags[0].Name != "Go" {
		t.Fatalf("tags not loaded: %+v", got.Tags)
	}

	// GetBySlug finds the same row.
	bySlug, err := repo.GetBySlug(ctx, "hello-world")
	if err != nil || bySlug.ID != p.ID {
		t.Fatalf("get by slug: %v (%+v)", err, bySlug)
	}

	// Update changes fields and clears the tag association.
	got.Title = "Hello Again"
	got.TagIDs = []int64{}
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	reloaded, err := repo.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if reloaded.Title != "Hello Again" {
		t.Fatalf("update did not persist title: %q", reloaded.Title)
	}
	if len(reloaded.Tags) != 0 {
		t.Fatalf("update did not clear tags: %+v", reloaded.Tags)
	}

	// Delete removes it; a second delete reports ErrNotFound.
	if err := repo.Delete(ctx, p.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := repo.Delete(ctx, p.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound on second delete, got %v", err)
	}
	if _, err := repo.GetByID(ctx, p.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestPostRepo_CreateDuplicateSlugConflicts(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	repo := NewPostRepo(db)
	author := seedUser(t, db)

	mk := func() *domain.Post {
		return &domain.Post{Slug: "dup", Title: "t", Status: domain.StatusDraft, AuthorID: author}
	}
	if err := repo.Create(ctx, mk()); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if err := repo.Create(ctx, mk()); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected ErrConflict on duplicate slug, got %v", err)
	}
}

func TestPostRepo_CreateUnknownTagIsInvalidReference(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	repo := NewPostRepo(db)
	author := seedUser(t, db)

	p := &domain.Post{Slug: "p", Title: "t", Status: domain.StatusDraft, AuthorID: author, TagIDs: []int64{999}}
	if err := repo.Create(ctx, p); !errors.Is(err, domain.ErrInvalidReference) {
		t.Fatalf("expected ErrInvalidReference for unknown tag, got %v", err)
	}
}

func TestPostRepo_ListFilters(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	repo := NewPostRepo(db)
	author := seedUser(t, db)
	goTag := seedTag(t, db, "Go", "go")
	seedTag(t, db, "Vue", "vue")

	mustCreate := func(slug, title, summary, contentMD string, status domain.PostStatus, pinned bool, tags []int64) {
		t.Helper()
		ts := time.Now()
		p := &domain.Post{
			Slug: slug, Title: title, Summary: summary, ContentMD: contentMD, Status: status,
			Pinned: pinned, AuthorID: author, PublishedAt: &ts, TagIDs: tags,
		}
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("create %s: %v", slug, err)
		}
	}
	mustCreate("a-published", "Alpha", "about go", "", domain.StatusPublished, false, []int64{goTag})
	mustCreate("b-pinned", "Bravo", "pinned one", "", domain.StatusPublished, true, nil)
	mustCreate("c-draft", "Charlie", "secret", "", domain.StatusDraft, false, nil)
	// Body-only term: appears in neither title nor summary, only content_md.
	mustCreate("d-body", "Delta", "no keyword here", "deep in the **xyzzy** paragraph", domain.StatusPublished, false, nil)

	// Status filter: the three published posts (a, b, d); c is a draft.
	items, total, err := repo.List(ctx, domain.PostListFilter{Status: domain.StatusPublished, Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list published: %v", err)
	}
	if total != 3 || len(items) != 3 {
		t.Fatalf("expected 3 published, got total=%d len=%d", total, len(items))
	}
	// Pinned post sorts first.
	if items[0].Slug != "b-pinned" {
		t.Fatalf("pinned post should sort first, got %q", items[0].Slug)
	}

	// Query filter matches title or summary, case-insensitively.
	items, total, err = repo.List(ctx, domain.PostListFilter{Query: "go", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list query: %v", err)
	}
	if total != 1 || items[0].Slug != "a-published" {
		t.Fatalf("query 'go' should match one post, got total=%d items=%+v", total, items)
	}

	// Query filter also reaches the post body (content_md), not just title/summary.
	items, total, err = repo.List(ctx, domain.PostListFilter{Query: "xyzzy", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list body query: %v", err)
	}
	if total != 1 || items[0].Slug != "d-body" {
		t.Fatalf("query 'xyzzy' should match the body-only post, got total=%d items=%+v", total, items)
	}

	// Tag filter via the EXISTS subquery.
	items, total, err = repo.List(ctx, domain.PostListFilter{TagSlug: "go", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list by tag: %v", err)
	}
	if total != 1 || items[0].Slug != "a-published" {
		t.Fatalf("tag 'go' should match one post, got total=%d items=%+v", total, items)
	}

	// No filter returns all four.
	_, total, err = repo.List(ctx, domain.PostListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if total != 4 {
		t.Fatalf("expected 4 total, got %d", total)
	}
}

func TestPostRepo_IncrementViews(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	repo := NewPostRepo(db)
	author := seedUser(t, db)

	p := &domain.Post{Slug: "viewed", Title: "t", Status: domain.StatusPublished, AuthorID: author}
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("create: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := repo.IncrementViews(ctx, p.ID); err != nil {
			t.Fatalf("increment views: %v", err)
		}
	}
	got, err := repo.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ViewCount != 3 {
		t.Fatalf("expected view_count 3, got %d", got.ViewCount)
	}
}
