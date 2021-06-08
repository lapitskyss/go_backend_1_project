// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"context"
	"github.com/google/wire"
	"github.com/lapitskyss/go_backend_1_project/internal/repository/postgres"
	"github.com/lapitskyss/go_backend_1_project/internal/server"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func InitializeAPIService() (*ApiService, func(), error) {
	context, cleanup, err := InitContext()
	if err != nil {
		return nil, nil, err
	}
	sugaredLogger, cleanup2, err := InitLogger()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	store, cleanup3, err := InitPostgresqlStore(context)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	api, cleanup4, err := InitServer(sugaredLogger, store)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	apiService, err := InitApiService(context, sugaredLogger, api)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return apiService, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

type ApiService struct {
	Log *zap.SugaredLogger
}

var APISet = wire.NewSet(
	InitApiService,
	InitContext,
	InitLogger,
	InitServer,
	InitPostgresqlStore,
)

func InitApiService(ctx context.Context, log *zap.SugaredLogger, api *server.Api) (*ApiService, error) {
	return &ApiService{
		Log: log,
	}, nil
}

func InitContext() (context.Context, func(), error) {
	ctx := context.Background()

	cb := func() {
		ctx.Done()
	}

	return ctx, cb, nil
}

func InitLogger() (*zap.SugaredLogger, func(), error) {
	logger, _ := zap.NewProduction()

	cleanup := func() {
		logger.Sync()
	}

	sugar := logger.Sugar()

	return sugar, cleanup, nil
}

func InitServer(log *zap.SugaredLogger, rep *postgres.Store) (*server.Api, func(), error) {
	server2 := server.New(log, rep)

	cleanup := func() {
		server2.
			Stop()
	}
	server2.
		Start()

	return server2, cleanup, nil
}

func InitPostgresqlStore(ctx context.Context) (*postgres.Store, func(), error) {
	store := &postgres.Store{}
	err := store.Init(ctx)

	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		store.Close()
	}

	return store, cleanup, nil
}
