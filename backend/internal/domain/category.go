package domain

import "context"

// Category is a (possibly nested) taxonomy node used to organise posts.
type Category struct {
	ID          int64
	ParentID    *int64 // nil for a top-level category
	Name        string
	Slug        string
	Description string
	Sort        int
}

// CategoryRepository is the persistence port for categories, defined where it is
// consumed and implemented by the store layer.
type CategoryRepository interface {
	Create(ctx context.Context, c *Category) error
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
}
