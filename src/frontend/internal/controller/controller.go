package controller

import (
	"context"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/files"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/rpc"
)

type controller struct {
	ctx context.Context
	log *zap.Logger
	fe  *rpc.FrontendServer
	tmp *files.Templates
}

func NewController(ctx context.Context, log *zap.Logger, fe *rpc.FrontendServer, tmp *files.Templates) *controller {
	return &controller{
		ctx: ctx,
		log: log,
		fe:  fe,
		tmp: tmp,
	}
}
