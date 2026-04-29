package library

import (
	"context"

	"github.com/project/library/internal/entity"
)

type LibraryUseCase interface {
	AddBook(ctx context.Context, bookName string, authorIds []string) (*entity.Book, error)
	GetBook(ctx context.Context, ID string) (*entity.Book, error)
	UpdateBook(ctx context.Context, id, name string, authorIDs []string) (*entity.Book, error)
	RegisterAuthor(ctx context.Context, name string) (*entity.Author, error)
	GetAuthor(ctx context.Context, ID string) (*entity.Author, error)
	UpdateAuthor(ctx context.Context, id, name string) (*entity.Author, error)
	GetAuthorBooks(ctx context.Context, ID string) ([]*entity.Book, error)
}
