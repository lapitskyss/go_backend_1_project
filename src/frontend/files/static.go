package files

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed static
var res embed.FS

func FileServerForStatic(r chi.Router) {
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(res)).ServeHTTP(w, r)
	})
}
