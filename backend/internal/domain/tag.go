package domain

import "context"

// Tag is a lightweight label that can be attached to posts.
type Tag struct {
	ID   int64
	Name string
	Slug string
}

// TagRepository is the persistence port for tags, defined where it is consumed
// and implemented by the store layer.
type TagRepository interface {
	Create(ctx context.Context, t *Tag) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Tag, error)
	List(ctx context.Context) ([]*Tag, error)
}
