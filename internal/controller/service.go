package controller

import (
	generated "github.com/project/library/generated/api/library"
	"github.com/project/library/internal/usecase/library"
	"go.uber.org/zap"
)

type Implementation struct {
	logger      *zap.Logger
	bookUseCase library.BookUseCase
	generated.UnimplementedLibraryServer
}

func New(
	logger *zap.Logger,
	booksUseCase library.BookUseCase,
) *Implementation {
	return &Implementation{
		logger:      logger,
		bookUseCase: booksUseCase,
	}
}
