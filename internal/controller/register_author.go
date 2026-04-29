package controller

import (
	"context"

	"github.com/project/library/generated/api/library"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) RegisterAuthor(ctx context.Context, req *library.RegisterAuthorRequest) (*library.RegisterAuthorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	s.logger.Info("RegisterAuthor called", zap.String("Name", req.Name))

	author, err := s.useCase.RegisterAuthor(ctx, req.Name)
	if err != nil {
		s.logger.Error("RegisterAuthor failed", zap.String("Name", req.Name), zap.Error(err))
		return nil, toGRPCError(err)
	}
	return &library.RegisterAuthorResponse{Id: author.ID}, nil
}
