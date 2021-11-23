package handler

import (
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/files"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/config"
)

type Handler struct {
	linkSrv LinkService
	config  *config.Config
	tmp     *files.Templates
	log     *zap.Logger
}

func InitHandler(linkSrv LinkService, config *config.Config, tmp *files.Templates, log *zap.Logger) *Handler {
	return &Handler{
		linkSrv: linkSrv,
		config:  config,
		tmp:     tmp,
		log:     log,
	}
}
