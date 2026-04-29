package repository

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/project/library/internal/entity"
	"go.uber.org/zap"
)

type Storage struct {
	books       map[string]*entity.Book
	authors     map[string]*entity.Author
	authorBooks map[string][]string
	logger      *zap.Logger
	mu          sync.RWMutex
}

func New(logger *zap.Logger) *Storage {
	return &Storage{
		books:       make(map[string]*entity.Book),
		authors:     make(map[string]*entity.Author),
		authorBooks: make(map[string][]string),
		logger:      logger,
	}
}

func (s *Storage) AddBook(book *entity.Book) (*entity.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, authorID := range book.AuthorIDs {
		if _, ok := s.authors[authorID]; !ok {
			return nil, ErrAuthorNotFound
		}
	}

	s.books[book.ID] = book

	for _, authorID := range book.AuthorIDs {
		s.authorBooks[authorID] = append(s.authorBooks[authorID], book.ID)
	}

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

	for _, authorID := range book.AuthorIDs {
		if _, ok := s.authors[authorID]; !ok {
			return nil, ErrAuthorNotFound
		}
	}

	oldIDs := make(map[string]bool, len(val.AuthorIDs))
	for _, id := range val.AuthorIDs {
		oldIDs[id] = true
	}

	newIDs := make(map[string]bool, len(book.AuthorIDs))
	for _, id := range book.AuthorIDs {
		newIDs[id] = true
	}

	for _, oldAuthorID := range val.AuthorIDs {
		if !newIDs[oldAuthorID] {
			ids := s.authorBooks[oldAuthorID]
			for i, id := range ids {
				if id == book.ID {
					s.authorBooks[oldAuthorID] = slices.Delete(ids, i, i+1)
					break
				}
			}
		}
	}

	for _, newAuthorID := range book.AuthorIDs {
		if !oldIDs[newAuthorID] {
			s.authorBooks[newAuthorID] = append(s.authorBooks[newAuthorID], book.ID)
		}
	}

	val.Name = book.Name
	val.AuthorIDs = book.AuthorIDs
	val.UpdatedAt = time.Now()
	s.logger.Info("book updated", zap.String("id", book.ID))
	return val, nil
}

func (s *Storage) RegisterAuthor(author *entity.Author) (*entity.Author, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authors[author.ID] = author
	s.logger.Info("author added", zap.String("id", author.ID))
	return author, nil
}

func (s *Storage) GetAuthor(ID string) (*entity.Author, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.authors[ID]
	if !ok {
		return nil, ErrAuthorNotFound
	}
	return val, nil
}

func (s *Storage) UpdateAuthor(author *entity.Author) (*entity.Author, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.authors[author.ID]
	if !ok {
		return nil, ErrAuthorNotFound
	}
	val.Name = author.Name
	val.UpdatedAt = time.Now()
	s.logger.Info("author updated", zap.String("id", author.ID))
	return val, nil
}

func (s *Storage) GetAuthorBooks(ID string) ([]*entity.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.authors[ID]; !ok {
		return nil, ErrAuthorNotFound
	}
	bookIDs := s.authorBooks[ID]
	result := make([]*entity.Book, 0, len(bookIDs))
	for _, bookID := range bookIDs {
		result = append(result, s.books[bookID])
	}
	return result, nil
}

var ErrBookNotFound = errors.New("book not found")
var ErrAuthorNotFound = errors.New("author not found")
