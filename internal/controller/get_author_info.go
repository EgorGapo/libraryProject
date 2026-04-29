package controller

import (
	"context"

	"github.com/project/library/generated/api/library"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetAuthorInfo(ctx context.Context, req *library.GetAuthorInfoRequest) (*library.GetAuthorInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	if err := validateUUID(req.Id); err != nil {
		return nil, err
	}
	s.logger.Info("GetAuthorInfo called", zap.String("id", req.Id))

	author, err := s.useCase.GetAuthor(ctx, req.Id)
	if err != nil {
		s.logger.Error("GetAuthorInfo failed", zap.String("id", req.Id), zap.Error(err))
		return nil, toGRPCError(err)
	}
	return &library.GetAuthorInfoResponse{Id: author.ID, Name: author.Name}, nil
}
