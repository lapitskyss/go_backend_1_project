package serverErrors

import (
	"github.com/go-chi/render"
	"net/http"
)

func RenderBadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, struct {
		Status bool  `json:"status"`
		Error  error `json:"error"`
	}{false, err})
}
