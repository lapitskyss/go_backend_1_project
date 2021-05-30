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

	"github.com/lapitskyss/go_backend_1_project/internal/repository"
	"github.com/lapitskyss/go_backend_1_project/internal/server/controller"
)

type Api struct {
	server http.Server
}

func New(log *zap.SugaredLogger, rep repository.Repository) *Api {
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
	r.NotFound(NotFoundHandler)

	linkController := controller.NewLinkController(log, rep)

	r.Post("/api/link", linkController.Add)
	r.Get("/api/links", linkController.List)

	return &Api{
		server: http.Server{
			Addr:    ":3000",
			Handler: r,

			ReadTimeout:       1 * time.Second,
			WriteTimeout:      90 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
	}
}

func (api *Api) Start() {
	go func() {
		_ = api.server.ListenAndServe()
	}()
}

func (api *Api) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return api.server.Shutdown(ctx)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}{false, "Not found"})
}
