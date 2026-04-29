package repository

import "github.com/project/library/internal/entity"

type StorageInterface interface {
	AddBook(book *entity.Book) (*entity.Book, error)
	GetBook(ID string) (*entity.Book, error)
	UpdateBook(book *entity.Book) (*entity.Book, error)
}
