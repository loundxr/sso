package app

import (
	"log/slog"
	grpc_app "sso/internal/app/grpc"
	service_auth "sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCServer *grpc_app.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := service_auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpc_app.New(log, authService, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
