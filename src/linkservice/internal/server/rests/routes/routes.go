package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/rests/handler"
)

func Routes(r *chi.Mux, linkHandler *handler.LinkHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/links", linkHandler.Create)
		r.Get("/links", linkHandler.List)
		r.Get("/links/{hash}", linkHandler.Info)
		r.Get("/links/search", linkHandler.Search)
	})
}
