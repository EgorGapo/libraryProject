package entity

import "errors"

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrAuthorNotFound    = errors.New("author not found")
	ErrInvalidAuthorName = errors.New("invalid author name")
)
