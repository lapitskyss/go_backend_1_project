package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/files"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/server/handler"
	mw "github.com/lapitskyss/go_backend_1_project/src/frontend/internal/server/middleware"
)

type Frontend struct {
	server http.Server
	log    *zap.Logger
	errors chan error
}

func NewFrontendServer(port string, h *handler.Handler, log *zap.Logger) *Frontend {
	r := chi.NewRouter()

	r.Use(mw.Recoverer(log))

	files.FileServerForStatic(r)
	r.Get("/{hash}", h.Redirect)
	r.Get("/*", h.Home)

	return &Frontend{
		server: http.Server{
			Addr:    ":" + port,
			Handler: r,

			ReadTimeout:       1 * time.Second,
			WriteTimeout:      90 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
		log: log,
	}
}

func (f *Frontend) Start() {
	f.log.Info("Frontend server started on port " + f.server.Addr + ".")
	go func() {
		err := f.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			f.log.Error("Frontend server return error", zap.Error(err))
			f.errors <- err
		}
	}()
}

func (f *Frontend) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.server.Shutdown(ctx)
}

func (f *Frontend) Notify() <-chan error {
	return f.errors
}
