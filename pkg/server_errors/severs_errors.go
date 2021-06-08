package server_errors

import (
	"net/http"

	"github.com/go-chi/render"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}{false, "Not found"})
}

func RenderBadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, struct {
		Status bool  `json:"status"`
		Error  error `json:"error"`
	}{false, err})
}
