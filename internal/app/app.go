package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/project/library/config"
	generated "github.com/project/library/generated/api/library"
	"github.com/project/library/internal/controller"
	"github.com/project/library/internal/usecase/library"
	"github.com/project/library/internal/usecase/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(logger *zap.Logger, cfg *config.Config) {
	storage := repository.New(logger)
	usecase := library.New(storage, logger)
	ctrl := controller.New(logger, usecase)

	grpcServer := newGrpcServer(ctrl)

	go runGrpc(grpcServer, cfg)
	go runRest(cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	grpcServer.GracefulStop()
}

func runRest(cfg *config.Config) {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := generated.RegisterLibraryHandlerFromEndpoint(ctx, mux, "127.0.0.1:"+cfg.GRPC.Port, opts); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":"+cfg.GRPC.GatewayPort, corsMiddleware(mux)); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func newGrpcServer(ctrl *controller.Implementation) *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)
	generated.RegisterLibraryServer(server, ctrl)
	return server
}

func runGrpc(server *grpc.Server, cfg *config.Config) {
	port := ":" + cfg.GRPC.Port
	lsn, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	if err := server.Serve(lsn); err != nil {
		panic(err)
	}
}
