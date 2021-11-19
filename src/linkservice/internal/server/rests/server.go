package rests

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/rests/handler"
	mw "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/rests/middleware"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/rests/routes"
)

type RESTServer struct {
	server http.Server
	log    *zap.Logger
	errors chan error
}

func NewRESTServer(port string, log *zap.Logger, linkHandler *handler.LinkHandler) *RESTServer {
	r := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(mw.Recoverer(log))
	r.Use(corsHandler.Handler)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.NotFound(mw.NotFound)

	routes.Routes(r, linkHandler)

	return &RESTServer{
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

func (s *RESTServer) Start() {
	s.log.Info("Linkservice HTTP server started on port " + s.server.Addr + ".")
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.log.Error("Linkservice HTTP server return error", zap.Error(err))
			s.errors <- err
		}
	}()
}

func (s *RESTServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *RESTServer) Notify() <-chan error {
	return s.errors
}
