package controller

import (
	generated "github.com/project/library/generated/api/library"
	"github.com/project/library/internal/usecase/library"
	"go.uber.org/zap"
)

type Implementation struct {
	logger  *zap.Logger
	useCase library.LibraryUseCase
	generated.UnimplementedLibraryServer
}

func New(
	logger *zap.Logger,
	useCase library.LibraryUseCase,
) *Implementation {
	return &Implementation{
		logger:  logger,
		useCase: useCase,
	}
}
