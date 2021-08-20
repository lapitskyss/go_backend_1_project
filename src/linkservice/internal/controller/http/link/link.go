package link

import (
	"context"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/repository"
)

type linkController struct {
	ctx context.Context
	log *zap.SugaredLogger
	rep repository.Store
}

func New(ctx context.Context, log *zap.SugaredLogger, rep repository.Store) *linkController {
	return &linkController{ctx, log, rep}
}
