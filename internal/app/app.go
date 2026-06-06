package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/project/library/config"
	"github.com/project/library/db"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Postgres
	poolCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN())
	if err != nil {
		logger.Fatal("db parse config", zap.Error(err))
	}
	poolCfg.MaxConns = cfg.Postgres.MaxConn

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		logger.Fatal("db connect", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("db ping", zap.Error(err))
	}

	db.SetupPostgres(pool, logger)

	// 2. Зависимости
	storage := repository.NewPostgresRepository(pool, logger)
	usecase := library.New(storage, logger)
	ctrl := controller.New(logger, usecase)

	grpcServer := newGrpcServer(ctrl)
	httpSrv := newHTTPServer(ctx, cfg)

	// 3. Запуск серверов с каналом ошибок
	errCh := make(chan error, 2)

	go func() {
		errCh <- runGrpc(grpcServer, cfg)
	}()

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// 4. Ждём сигнал ИЛИ падение одного из серверов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-quit:
		logger.Info("shutdown signal", zap.String("signal", sig.String()))
	case err := <-errCh:
		logger.Error("server failed", zap.Error(err))
	}

	// 5. Graceful shutdown с таймаутом
	shutdownCtx, sCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer sCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown failed", zap.Error(err))
	}

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		logger.Info("grpc stopped gracefully")
	case <-shutdownCtx.Done():
		logger.Warn("grpc graceful stop timeout, forcing")
		grpcServer.Stop()
	}
}

func newHTTPServer(ctx context.Context, cfg *config.Config) *http.Server {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := generated.RegisterLibraryHandlerFromEndpoint(ctx, mux, "127.0.0.1:"+cfg.GRPC.Port, opts); err != nil {
		panic(err)
	}
	return &http.Server{
		Addr:    ":" + cfg.GRPC.GatewayPort,
		Handler: corsMiddleware(mux),
	}
}

func runGrpc(server *grpc.Server, cfg *config.Config) error {
	lsn, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		return err
	}
	return server.Serve(lsn)
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
