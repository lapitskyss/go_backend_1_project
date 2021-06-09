package link

import (
	"context"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/internal/repository/postgres"
)

type linkController struct {
	ctx context.Context
	log *zap.SugaredLogger
	rep *postgres.Store
}

func New(ctx context.Context, log *zap.SugaredLogger, rep *postgres.Store) *linkController {
	return &linkController{ctx, log, rep}
}
