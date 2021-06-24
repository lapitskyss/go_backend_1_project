package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	linkhttp "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/link/http"
	mw "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/middleware"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/postgres"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/server_errors"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/server_port"
)

type HTTPServer struct {
	server http.Server
	log    *zap.SugaredLogger
}

func NewHTTPServer(ctx context.Context, log *zap.SugaredLogger, rep *postgres.Store) *HTTPServer {
	r := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(corsHandler.Handler)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(mw.Recoverer(log))
	r.Use(middleware.Timeout(60 * time.Second))
	r.NotFound(server_errors.NotFoundError)

	linkController := linkhttp.New(ctx, log, rep)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/links", linkController.Create)
		r.Get("/links", linkController.List)
		r.Get("/links/{hash}/statistics", linkController.Statistics)
	})

	port := server_port.GetServerPortFromEnv("LINKSERVICE_HTTP_PORT", 3000)

	return &HTTPServer{
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

func (srv *HTTPServer) Start() {
	srv.log.Info("Linkservice HTTP server started.")
	go func() {
		err := srv.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			srv.log.Error("Linkservice HTTP server return error", zap.NamedError("sever_error", err))
		}
	}()
}

func (srv *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.server.Shutdown(ctx)
}
