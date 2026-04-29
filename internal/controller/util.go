package controller

import (
	"errors"

	"github.com/project/library/internal/usecase/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, repository.ErrBookNotFound):
		return status.Error(codes.NotFound, "book not found")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
