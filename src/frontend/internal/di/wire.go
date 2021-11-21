//go:generate wire
//go:build wireinject
// +build wireinject

package di

import (
	"strconv"

	"github.com/caarlos0/env/v6"
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/files"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/config"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/linksrv"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/server"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/server/handler"
)

type Frontend struct {
	Log *zap.Logger
	Srv *server.Frontend
}

var FrontendSet = wire.NewSet(
	InitFrontend,
	InitConfig,
	InitLogger,
	InitTemplates,
	InitServer,
)

func InitFrontend(log *zap.Logger, f *server.Frontend) (*Frontend, error) {
	return &Frontend{
		Log: log,
		Srv: f,
	}, nil
}

func InitializeFrontend() (*Frontend, func(), error) {
	panic(wire.Build(FrontendSet))
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

func InitTemplates() (*files.Templates, error) {
	return files.InitTemplates()
}

func InitServer(config *config.Config, tmp *files.Templates, log *zap.Logger) (*server.Frontend, func(), error) {
	linkService, err := linksrv.InitLinkServiceClient(config.LinkServiceAddr)
	if err != nil {
		return nil, nil, err
	}

	h := handler.InitHandler(linkService, tmp, log)
	srv := server.NewFrontendServer(strconv.Itoa(config.ServerPort), h, log)

	cleanup := func() {
		_ = srv.Stop()
	}

	srv.Start()

	return srv, cleanup, nil
}
