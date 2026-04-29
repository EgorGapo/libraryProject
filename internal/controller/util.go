package controller

import (
	"errors"

	"github.com/google/uuid"
	librarysvc "github.com/project/library/internal/usecase/library"
	"github.com/project/library/internal/usecase/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, repository.ErrBookNotFound):
		return status.Error(codes.NotFound, "book not found")
	case errors.Is(err, repository.ErrAuthorNotFound):
		return status.Error(codes.NotFound, "author not found")
	case errors.Is(err, librarysvc.ErrInvalidAuthorName):
		return status.Error(codes.InvalidArgument, "invalid author name")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

func validateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return status.Error(codes.InvalidArgument, "invalid id format")
	}
	return nil
}
