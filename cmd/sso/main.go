package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.Setup(cfg.Env)
	log.Info("starting application", slog.Any("cfg", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCServer.MustRun()

	// TODO: start gRPC app-server

	// graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sig := <-stop
	log.Info("stopping application", slog.String("signal", sig.String()))

	application.GRPCServer.Stop()

	log.Info("application stopped")
}
