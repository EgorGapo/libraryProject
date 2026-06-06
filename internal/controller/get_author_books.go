package controller

import (
	"github.com/project/library/generated/api/library"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Implementation) GetAuthorBooks(req *library.GetAuthorBooksRequest, stream grpc.ServerStreamingServer[library.Book]) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if err := validateUUID(req.AuthorId); err != nil {
		return err
	}

	s.logger.Info("GetAuthorBooks called", zap.String("author_id", req.AuthorId))

	books, err := s.useCase.GetAuthorBooks(stream.Context(), req.AuthorId)
	if err != nil {
		s.logger.Error("GetAuthorBooks failed", zap.String("author_id", req.AuthorId), zap.Error(err))
		return toGRPCError(err)
	}

	for _, book := range books {
		if err := stream.Send(&library.Book{
			Id:        book.ID,
			Name:      book.Name,
			AuthorId:  book.AuthorIDs,
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
		}); err != nil {
			return err
		}
	}
	return nil
}
