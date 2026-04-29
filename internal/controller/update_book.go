package controller

import (
	"context"

	"github.com/project/library/generated/api/library"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) UpdateBook(ctx context.Context, req *library.UpdateBookRequest) (*library.UpdateBookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	s.logger.Info("UpdateBook called", zap.String("id", req.Id))
	_, err := s.bookUseCase.UpdateBook(ctx, req.Id, req.Name, req.AuthorIds)
	if err != nil {
		s.logger.Error("UpdateBook failed", zap.String("id", req.Id), zap.Error(err))
		return nil, toGRPCError(err)
	}
	return &library.UpdateBookResponse{}, nil
}
