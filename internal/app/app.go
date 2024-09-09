package app

import (
	"log/slog"
	"sso/internal/app/grpcapp"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokentTTL time.Duration,
) *App {
	// TODO: Инициализировать хранилище
	storage, err := sqlite.New(storagePath)
	storage.InitApp()
	if err != nil {
		panic(err)
	}

	// TODO: init auth service (auth)
	authService := auth.New(log, storage, storage, storage, tokentTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
