package server_errors

import (
	"net/http"

	"github.com/go-chi/render"
)

func NotFoundError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}{false, "Not found."})
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}{false, err.Error()})
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}{false, "Unexpected error occurred. Please try request later."})
}
