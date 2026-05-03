package library

import (
	"regexp"

	"github.com/project/library/internal/entity"
)

var authorNameRegexp = regexp.MustCompile(`^[\p{L}\d ]+$`)

func validateAuthorName(name string) error {
	if len(name) == 0 || len(name) >= 1024 {
		return entity.ErrInvalidAuthorName
	}
	if !authorNameRegexp.MatchString(name) {
		return entity.ErrInvalidAuthorName
	}
	return nil
}
