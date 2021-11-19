package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/files"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/controller"
	mw "github.com/lapitskyss/go_backend_1_project/src/frontend/internal/middleware"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/rpc"
)

type Frontend struct {
	server http.Server
	errors chan error
	log    *zap.Logger
}

func NewFrontendServer(ctx context.Context, log *zap.Logger, fe *rpc.FrontendServer, tmp *files.Templates) *Frontend {
	r := chi.NewRouter()

	r.Use(mw.Recoverer(log))

	cnt := controller.NewController(ctx, log, fe, tmp)

	files.FileServerForStatic(r)
	r.Get("/{hash}", cnt.Redirect)
	r.Get("/*", cnt.Home)

	// TODO get server port from env

	return &Frontend{
		server: http.Server{
			Addr:    ":3001",
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
	f.log.Info("Frontend server started.")
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
