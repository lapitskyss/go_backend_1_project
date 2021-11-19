package middleware

import (
	"net/http"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/pkg/render"
)

func NotFound(w http.ResponseWriter, _ *http.Request) {
	render.NotFoundError(w)
}
