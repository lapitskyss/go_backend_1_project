//go:generate wire
//+build wireinject

package di

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/postgres"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server"
)

type LinkService struct {
	Log *zap.SugaredLogger
}

var LinkServiceSet = wire.NewSet(
	InitLinkService,
	InitContext,
	InitLogger,
	InitHttpServer,
	InitGRPCServer,
	InitPostgresqlStore,
)

func InitLinkService(log *zap.SugaredLogger, hs *server.HTTPServer, gs *server.GRPSServer) (*LinkService, error) {
	return &LinkService{
		Log: log,
	}, nil
}

func InitializeLinkService() (*LinkService, func(), error) {
	panic(wire.Build(LinkServiceSet))
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

func InitHttpServer(ctx context.Context, log *zap.SugaredLogger, rep *postgres.Store) (*server.HTTPServer, func(), error) {
	server := server.NewHTTPServer(ctx, log, rep)

	cleanup := func() {
		server.Stop()
	}

	server.Start()

	return server, cleanup, nil
}

func InitGRPCServer(ctx context.Context, log *zap.SugaredLogger, rep *postgres.Store) (*server.GRPSServer, func(), error) {
	server := server.NewGRPCServer(ctx, log, rep)

	cleanup := func() {
		server.StopServer()
	}

	server.StartServer()

	return server, cleanup, nil
}

func InitPostgresqlStore(ctx context.Context) (*postgres.Store, func(), error) {
	store, err := postgres.NewStore(ctx)

	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		store.Close()
	}

	return store, cleanup, nil
}
