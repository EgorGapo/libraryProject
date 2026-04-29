package app

import (
	"net"

	"github.com/project/library/config"
	generated "github.com/project/library/generated/api/library"
	"github.com/project/library/internal/controller"
	"github.com/project/library/internal/usecase/library"
	"github.com/project/library/internal/usecase/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(logger *zap.Logger, cfg *config.Config) {
	storage := repository.New(logger)
	usecase := library.New(storage, logger)
	controller := controller.New(logger, usecase)
	runGrpc(controller, cfg)

}

func runRest() {}

func runGrpc(controller *controller.Implementation, cfg *config.Config) {
	port := ":" + cfg.GRPC.Port
	lsn, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	reflection.Register(server)
	generated.RegisterLibraryServer(server, controller)
	if err := server.Serve(lsn); err != nil {
		panic(err)
	}
}
