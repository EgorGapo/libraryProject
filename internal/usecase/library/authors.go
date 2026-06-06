package library

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
)

func (s *LibraryImpl) RegisterAuthor(ctx context.Context, name string) (*entity.Author, error) {
	if err := validateAuthorName(name); err != nil {
		return nil, err
	}
	now := time.Now()
	author := &entity.Author{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.storage.RegisterAuthor(ctx, author)
}

func (s *LibraryImpl) GetAuthor(ctx context.Context, ID string) (*entity.Author, error) {
	return s.storage.GetAuthor(ctx, ID)
}

func (s *LibraryImpl) UpdateAuthor(ctx context.Context, id, name string) (*entity.Author, error) {
	if err := validateAuthorName(name); err != nil {
		return nil, err
	}
	author := &entity.Author{
		ID:   id,
		Name: name,
	}
	return s.storage.UpdateAuthor(ctx, author)
}

func (s *LibraryImpl) GetAuthorBooks(ctx context.Context, ID string) ([]*entity.Book, error) {
	return s.storage.GetAuthorBooks(ctx, ID)
}
