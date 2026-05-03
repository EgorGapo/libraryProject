package controller

import (
	"errors"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, entity.ErrBookNotFound):
		return status.Error(codes.NotFound, "book not found")
	case errors.Is(err, entity.ErrAuthorNotFound):
		return status.Error(codes.NotFound, "author not found")
	case errors.Is(err, entity.ErrInvalidAuthorName):
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
