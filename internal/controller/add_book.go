package controller

import (
	"context"

	"github.com/project/library/generated/api/library"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Implementation) AddBook(ctx context.Context, req *library.AddBookRequest) (*library.AddBookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	s.logger.Info("AddBook called", zap.String("name", req.Name))

	book, err := s.useCase.AddBook(ctx, req.Name, req.AuthorIds)
	if err != nil {
		s.logger.Error("AddBook failed", zap.Error(err))
		return nil, toGRPCError(err)
	}
	return &library.AddBookResponse{Book: &library.Book{
		Id:        book.ID,
		Name:      book.Name,
		AuthorId:  book.AuthorIDs,
		CreatedAt: timestamppb.New(book.CreatedAt),
		UpdatedAt: timestamppb.New(book.UpdatedAt),
	}}, nil
}
