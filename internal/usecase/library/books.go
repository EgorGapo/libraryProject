package library

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
	"github.com/project/library/internal/usecase/repository"
	"go.uber.org/zap"
)

type BookUseCaseImp struct {
	storage repository.StorageInterface
	logger  *zap.Logger
}

func New(storage repository.StorageInterface, logger *zap.Logger) *BookUseCaseImp {
	return &BookUseCaseImp{
		storage: storage,
		logger:  logger,
	}
}

func (s *BookUseCaseImp) AddBook(ctx context.Context, bookName string, aithorIds []string) (*entity.Book, error) {
	book := &entity.Book{
		ID:        uuid.New().String(),
		Name:      bookName,
		AuthorIDs: aithorIds,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.storage.AddBook(book)
}

func (s *BookUseCaseImp) GetBook(ctx context.Context, ID string) (*entity.Book, error) {
	return s.storage.GetBook(ID)
}

func (s *BookUseCaseImp) UpdateBook(ctx context.Context, id, name string, authorIDs []string) (*entity.Book, error) {
	book := &entity.Book{
		ID:        id,
		Name:      name,
		AuthorIDs: authorIDs,
	}
	return s.storage.UpdateBook(book)
}
