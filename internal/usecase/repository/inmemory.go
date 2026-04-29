package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/project/library/internal/entity"
	"go.uber.org/zap"
)

type Storage struct {
	books  map[string]*entity.Book
	logger *zap.Logger
	mu     sync.RWMutex
}

func New(logger *zap.Logger) *Storage {
	return &Storage{
		books:  make(map[string]*entity.Book),
		logger: logger,
	}
}
func (s *Storage) AddBook(book *entity.Book) (*entity.Book, error) {
	s.mu.Lock()
	s.books[book.ID] = book
	s.mu.Unlock()
	s.logger.Info("book added", zap.String("id", book.ID))
	return book, nil
}

func (s *Storage) GetBook(ID string) (*entity.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.books[ID]
	if !ok {
		return nil, ErrBookNotFound
	}
	return val, nil
}

func (s *Storage) UpdateBook(book *entity.Book) (*entity.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.books[book.ID]
	if !ok {
		return nil, ErrBookNotFound
	}
	val.Name = book.Name
	val.AuthorIDs = book.AuthorIDs
	val.UpdatedAt = time.Now()
	s.logger.Info("book updated", zap.String("id", book.ID))
	return val, nil
}

var ErrBookNotFound = errors.New("book not found")
