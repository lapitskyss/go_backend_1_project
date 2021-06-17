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

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/link"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/postgres"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/server_errors"
)

type Api struct {
	server http.Server
	log    *zap.SugaredLogger
}

func New(ctx context.Context, log *zap.SugaredLogger, rep *postgres.Store) *Api {
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
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(60 * time.Second))
	r.NotFound(server_errors.NotFoundError)

	linkController := link.New(ctx, log, rep)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/links", linkController.Create)
		r.Get("/links", linkController.List)
		r.Get("/link/{hash}", linkController.Get)
	})

	return &Api{
		server: http.Server{
			Addr:    ":3000",
			Handler: r,

			ReadTimeout:       1 * time.Second,
			WriteTimeout:      90 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
		log: log,
	}
}

func (api *Api) Start() {
	api.log.Info("Server started.")
	go func() {
		err := api.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			api.log.Error("Server return error", zap.NamedError("sever_error", err))
		}
	}()
}

func (api *Api) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return api.server.Shutdown(ctx)
}
