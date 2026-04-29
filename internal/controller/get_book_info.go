package controller

import (
	"context"

	"github.com/project/library/generated/api/library"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetBookInfo(ctx context.Context, req *library.GetBookInfoRequest) (*library.GetBookInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	s.logger.Info("GetBookInfo called", zap.String("id", req.Id))

	book, err := s.bookUseCase.GetBook(ctx, req.Id)
	if err != nil {
		s.logger.Error("GetBookInfo failed", zap.String("id", req.Id), zap.Error(err))
		return nil, toGRPCError(err)
	}
	return &library.GetBookInfoResponse{Book: &library.Book{
		Id:       book.ID,
		Name:     book.Name,
		AuthorId: book.AuthorIDs,
	}}, nil
}
