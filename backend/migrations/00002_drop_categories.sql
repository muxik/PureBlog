-- +goose Up
-- Categories were removed from the product (the design organises posts with tags
-- only). Drop the posts.category_id association first, then the table itself.
ALTER TABLE posts DROP COLUMN IF EXISTS category_id;
DROP TABLE IF EXISTS categories;

-- +goose Down
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    sort        INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_categories_parent ON categories(parent_id);
ALTER TABLE posts ADD COLUMN category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL;
