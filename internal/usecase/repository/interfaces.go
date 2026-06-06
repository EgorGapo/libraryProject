package repository

import (
	"context"

	"github.com/project/library/internal/entity"
)

type StorageInterface interface {
	AddBook(ctx context.Context, book *entity.Book) (*entity.Book, error)
	GetBook(ctx context.Context, ID string) (*entity.Book, error)
	UpdateBook(ctx context.Context, book *entity.Book) (*entity.Book, error)
	RegisterAuthor(ctx context.Context, author *entity.Author) (*entity.Author, error)
	GetAuthor(ctx context.Context, ID string) (*entity.Author, error)
	UpdateAuthor(ctx context.Context, author *entity.Author) (*entity.Author, error)
	GetAuthorBooks(ctx context.Context, ID string) ([]*entity.Book, error)
}
