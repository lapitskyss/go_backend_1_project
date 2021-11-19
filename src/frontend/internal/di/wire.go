//go:generate wire
//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/files"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/server"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/rpc"
)

type FrontendService struct {
	Log *zap.Logger
	Srv *server.Frontend
}

var FrontendSet = wire.NewSet(
	InitFrontendService,
	InitContext,
	InitLogger,
	InitTemplates,
	InitGRPCClient,
	InitServer,
)

func InitFrontendService(log *zap.Logger, f *server.Frontend) (*FrontendService, error) {
	return &FrontendService{
		Log: log,
		Srv: f,
	}, nil
}

func InitializeFrontendService() (*FrontendService, func(), error) {
	panic(wire.Build(FrontendSet))
}

func InitContext() (context.Context, func(), error) {
	ctx := context.Background()

	cb := func() {
		ctx.Done()
	}

	return ctx, cb, nil
}

func InitLogger() (*zap.Logger, func(), error) {
	logger, _ := zap.NewProduction()

	cleanup := func() {
		_ = logger.Sync()
	}

	return logger, cleanup, nil
}

func InitTemplates() (*files.Templates, error) {
	return files.InitTemplates()
}

func InitGRPCClient(ctx context.Context) (*rpc.FrontendServer, error) {
	client, err := rpc.NewGRPCClient(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func InitServer(ctx context.Context, log *zap.Logger, fe *rpc.FrontendServer, tmp *files.Templates) (*server.Frontend, func(), error) {
	srv := server.NewFrontendServer(ctx, log, fe, tmp)

	cleanup := func() {
		_ = srv.Stop()
	}

	srv.Start()

	return srv, cleanup, nil
}
