package controller

import (
	"context"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/api"
)

type controller struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	client *api.Client
}

func NewController(ctx context.Context, log *zap.SugaredLogger, client *api.Client) *controller {
	return &controller{ctx, log, client}
}
