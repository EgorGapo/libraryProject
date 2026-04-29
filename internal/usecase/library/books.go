package library

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
	"github.com/project/library/internal/usecase/repository"
	"go.uber.org/zap"
)

type LibraryImpl struct {
	storage repository.StorageInterface
	logger  *zap.Logger
}

func New(storage repository.StorageInterface, logger *zap.Logger) *LibraryImpl {
	return &LibraryImpl{
		storage: storage,
		logger:  logger,
	}
}

func (s *LibraryImpl) AddBook(ctx context.Context, bookName string, authorIds []string) (*entity.Book, error) {
	book := &entity.Book{
		ID:        uuid.New().String(),
		Name:      bookName,
		AuthorIDs: authorIds,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.storage.AddBook(book)
}

func (s *LibraryImpl) GetBook(ctx context.Context, ID string) (*entity.Book, error) {
	return s.storage.GetBook(ID)
}

func (s *LibraryImpl) UpdateBook(ctx context.Context, id, name string, authorIDs []string) (*entity.Book, error) {
	book := &entity.Book{
		ID:        id,
		Name:      name,
		AuthorIDs: authorIDs,
	}
	return s.storage.UpdateBook(book)
}
