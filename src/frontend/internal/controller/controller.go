package controller

import (
	"context"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/rpc"
)

type controller struct {
	ctx context.Context
	log *zap.SugaredLogger
	fe  *rpc.FrontendServer
}

func NewController(ctx context.Context, log *zap.SugaredLogger, fe *rpc.FrontendServer) *controller {
	return &controller{ctx, log, fe}
}
