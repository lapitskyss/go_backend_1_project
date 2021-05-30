//go:generate wire
//+build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/internal/repository"
	fileStore "github.com/lapitskyss/go_backend_1_project/internal/repository/file"
	"github.com/lapitskyss/go_backend_1_project/internal/server"
)

type ApiService struct {
	Log *zap.SugaredLogger
}

var APISet = wire.NewSet(
	InitApiService,
	InitLogger,
	InitServer,
	InitFileStore,
)

func InitApiService(log *zap.SugaredLogger, api *server.Api) (*ApiService, error) {
	return &ApiService{
		Log: log,
	}, nil
}

func InitializeAPIService() (*ApiService, func(), error) {
	panic(wire.Build(APISet))
}

func InitLogger() (*zap.SugaredLogger, func(), error) {
	logger, _ := zap.NewProduction()

	cleanup := func() {
		logger.Sync()
	}

	sugar := logger.Sugar()

	return sugar, cleanup, nil
}

func InitServer(log *zap.SugaredLogger, rep repository.Repository) (*server.Api, func(), error) {
	server := server.New(log, rep)

	cleanup := func() {
		server.Stop()
	}

	server.Start()

	return server, cleanup, nil
}

func InitFileStore() (repository.Repository, func(), error) {
	fileStore := fileStore.New()

	cleanup := func() {

	}

	return fileStore, cleanup, nil
}
