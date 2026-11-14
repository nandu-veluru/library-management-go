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
}