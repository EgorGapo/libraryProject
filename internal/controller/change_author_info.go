package controller

import (
	"context"

	"github.com/project/library/generated/api/library"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) ChangeAuthorInfo(ctx context.Context, req *library.ChangeAuthorInfoRequest) (*library.ChangeAuthorInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	if err := validateUUID(req.Id); err != nil {
		return nil, err
	}
	s.logger.Info("ChangeAuthorInfo called", zap.String("id", req.Id))

	_, err := s.useCase.UpdateAuthor(ctx, req.Id, req.Name)
	if err != nil {
		s.logger.Error("ChangeAuthorInfo failed", zap.String("id", req.Id), zap.Error(err))
		return nil, toGRPCError(err)
	}
	return &library.ChangeAuthorInfoResponse{}, nil
}
