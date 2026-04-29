package library

import (
	"context"

	"github.com/project/library/internal/entity"
)

type BookUseCase interface {
	AddBook(ctx context.Context, bookName string, aithorIds []string) (*entity.Book, error)
	GetBook(ctx context.Context, ID string) (*entity.Book, error)
	UpdateBook(ctx context.Context, id, name string, authorIDs []string) (*entity.Book, error)
}
