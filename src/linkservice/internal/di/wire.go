//go:generate wire
//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"strconv"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/config"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/grpcs"
	grpch "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/grpcs/handler"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/rests"
	resth "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/rests/handler"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/linksrv"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/redirectsrv"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/store/pg"
)

type REST struct {
	Log        *zap.Logger
	RESTServer *rests.RESTServer
}

type GRPC struct {
	Log        *zap.Logger
	GRPCServer *grpcs.GRPCServer
}

var RestSet = wire.NewSet(
	InitRESTService,
	InitContext,
	InitConfig,
	InitLogger,
	InitPgStore,
	InitRESTServer,
)

var GRPCSet = wire.NewSet(
	InitGRPCService,
	InitContext,
	InitConfig,
	InitLogger,
	InitPgStore,
	InitGRPCServer,
)

func InitializeREST() (*REST, func(), error) {
	panic(wire.Build(RestSet))
}

func InitializeGRPC() (*GRPC, func(), error) {
	panic(wire.Build(GRPCSet))
}

func InitRESTService(log *zap.Logger, restServer *rests.RESTServer) (*REST, error) {
	return &REST{
		Log:        log,
		RESTServer: restServer,
	}, nil
}

func InitGRPCService(log *zap.Logger, grpcServer *grpcs.GRPCServer) (*GRPC, error) {
	return &GRPC{
		Log:        log,
		GRPCServer: grpcServer,
	}, nil
}

func InitContext() (context.Context, func(), error) {
	ctx := context.Background()

	cb := func() {
		ctx.Done()
	}

	return ctx, cb, nil
}

func InitConfig() (*config.Config, error) {
	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func InitLogger() (*zap.Logger, func(), error) {
	logger, _ := zap.NewProduction()

	cleanup := func() {
		_ = logger.Sync()
	}

	return logger, cleanup, nil
}

func InitPgStore(ctx context.Context, cfg *config.Config, log *zap.Logger) (*pg.Store, func(), error) {
	store, err := pg.Connect(ctx, cfg.DatabaseURL, log)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		store.Close()
	}

	return store, cleanup, nil
}

func InitRESTServer(log *zap.Logger, store *pg.Store, cfg *config.Config) (*rests.RESTServer, func(), error) {
	linkStore := pg.NewLinkStore(store)
	linkService := linksrv.NewLinkService(log, linkStore)
	linkHandler := resth.NewLinkHandler(log, linkService)

	srv := rests.NewRESTServer(strconv.Itoa(cfg.RESTPort), log, linkHandler)

	cleanup := func() {
		_ = srv.Stop()
	}

	srv.Start()

	return srv, cleanup, nil
}

func InitGRPCServer(log *zap.Logger, store *pg.Store, cfg *config.Config) (*grpcs.GRPCServer, func(), error) {
	redirectStore := pg.NewRedirectStore(store)
	redirectService := redirectsrv.NewRedirectService(log, redirectStore)
	redirectHandler := grpch.NewRedirectHandler(log, redirectService)

	srv := grpcs.NewGRPCServer(log, redirectHandler)

	cleanup := func() {
		srv.Stop()
	}

	srv.Start(strconv.Itoa(cfg.GRPCPort))

	return srv, cleanup, nil
}
