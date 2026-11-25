package author

import (
	"context"
	"errors"
)

var (
	ErrNotFound  = errors.New("Author not found")
	ErrDuplicate = errors.New("Author already exists")
)

type Repository interface {
	ExistByName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, author Author) (Author, error)
	ViewAuthor(ctx context.Context, id string) (Author, error)
	SearchAuthor(ctx context.Context, name string, limit, offset int) ([]Author, error)
	CountAuthors(ctx context.Context, name string) (int, error)
	HasBooksAuthors(ctx context.Context, id string) (bool, error)
	DeleteAuthor(ctx context.Context, id string) (Author, error)
	HasName(ctx context.Context, id string, name string) (bool, error)
	UpdateAuthor(ctx context.Context, author Author) (Author, error)
}
