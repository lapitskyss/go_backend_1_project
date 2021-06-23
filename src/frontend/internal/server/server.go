package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/controller"
	mw "github.com/lapitskyss/go_backend_1_project/src/frontend/internal/middleware"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/rpc"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/server_port"
)

type Frontend struct {
	server http.Server
	log    *zap.SugaredLogger
}

func NewFrontendServer(ctx context.Context, log *zap.SugaredLogger, fe *rpc.FrontendServer) *Frontend {
	r := chi.NewRouter()

	r.Use(mw.Recoverer(log))
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	cnt := controller.NewController(ctx, log, fe)

	FileServerForStatic(r)
	r.Get("/{hash}", cnt.Redirect)
	r.Get("/*", cnt.Home)

	port := server_port.GetServerPortFromEnv("FRONTEND_HTTP_PORT", 3001)

	return &Frontend{
		server: http.Server{
			Addr:    port,
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
			f.log.Error("Frontend server return error", zap.NamedError("sever_error", err))
		}
	}()
}

func (f *Frontend) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return f.server.Shutdown(ctx)
}
