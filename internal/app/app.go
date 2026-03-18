package app

import (
	"context"
	"log/slog"
	grpcapp "sso/internal/app/grpcapp"

	//"sso/internal/grpc/auth"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировать хранилище (storage)
	cxt := context.Background()
	storage, err := postgres.New(cxt, storagePath)
	if err != nil {
		panic(err)
	}

	// TODO: init auth service (auth)
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{GRPCSrv: grpcApp}

}
