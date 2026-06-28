-- +goose Up
-- Full-text-ish search via pg_trgm. Search previously hit posts.title and
-- posts.summary with ILIKE '%q%', which a B-tree index can't accelerate (the
-- leading wildcard). A GIN index with gin_trgm_ops makes those substring
-- (and fuzzy) matches index-backed, and lets us widen the search to the post
-- body (content_md) without a per-row sequential scan.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_posts_title_trgm   ON posts USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_summary_trgm ON posts USING gin (summary gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_content_trgm ON posts USING gin (content_md gin_trgm_ops);

-- +goose Down
DROP INDEX IF EXISTS idx_posts_content_trgm;
DROP INDEX IF EXISTS idx_posts_summary_trgm;
DROP INDEX IF EXISTS idx_posts_title_trgm;
-- Leave the pg_trgm extension in place: other objects may depend on it and
-- dropping an extension other migrations could rely on is needlessly risky.
