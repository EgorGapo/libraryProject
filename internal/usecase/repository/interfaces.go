package repository

import "github.com/project/library/internal/entity"

type StorageInterface interface {
	AddBook(book *entity.Book) (*entity.Book, error)
	GetBook(ID string) (*entity.Book, error)
	UpdateBook(book *entity.Book) (*entity.Book, error)
	RegisterAuthor(author *entity.Author) (*entity.Author, error)
	GetAuthor(ID string) (*entity.Author, error)
	UpdateAuthor(author *entity.Author) (*entity.Author, error)
	GetAuthorBooks(ID string) ([]*entity.Book, error)
}
