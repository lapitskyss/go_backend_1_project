package handler

import (
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/files"
)

type Handler struct {
	linkSrv LinkService
	tmp     *files.Templates
	log     *zap.Logger
}

func InitHandler(linkSrv LinkService, tmp *files.Templates, log *zap.Logger) *Handler {
	return &Handler{
		linkSrv: linkSrv,
		tmp:     tmp,
		log:     log,
	}
}
