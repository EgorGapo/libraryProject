package library

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidAuthorName = errors.New("invalid author name")
	authorNameRegexp     = regexp.MustCompile(`^[\p{L}\d ]+$`)
)

func validateAuthorName(name string) error {
	if len(name) == 0 || len(name) >= 1024 {
		return ErrInvalidAuthorName
	}
	if !authorNameRegexp.MatchString(name) {
		return ErrInvalidAuthorName
	}
	return nil
}
