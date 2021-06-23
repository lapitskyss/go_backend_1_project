//go:generate wire
//+build wireinject

package di

import (
	"context"
	"os"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/server"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/api"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/rpc"
)

type FrontendService struct {
	Log *zap.SugaredLogger
}

var FrontendSet = wire.NewSet(
	InitFrontendService,
	InitContext,
	InitLogger,
	InitClient,
	NewGRPCClient,
	InitServer,
)

func InitFrontendService(log *zap.SugaredLogger, f *server.Frontend) (*FrontendService, error) {
	return &FrontendService{
		Log: log,
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

func InitLogger() (*zap.SugaredLogger, func(), error) {
	logger, _ := zap.NewProduction()

	cleanup := func() {
		logger.Sync()
	}

	sugar := logger.Sugar()

	return sugar, cleanup, nil
}

func InitClient() (*api.Client, error) {
	client, err := api.NewClient(os.Getenv("LINKSERVICE_BASE_URL"), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewGRPCClient(ctx context.Context) (*rpc.FrontendServer, error) {
	client, err := rpc.NewGRPCClient(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func InitServer(ctx context.Context, log *zap.SugaredLogger, fe *rpc.FrontendServer) (*server.Frontend, func(), error) {
	server := server.NewFrontendServer(ctx, log, fe)

	cleanup := func() {
		server.Stop()
	}

	server.Start()

	return server, cleanup, nil
}
