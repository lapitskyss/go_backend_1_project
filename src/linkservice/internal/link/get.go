package link

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	se "github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/server_errors"
)

func (api *linkController) Get(w http.ResponseWriter, r *http.Request) {
	var hash = chi.URLParam(r, "hash")
	if hash == "" {
		se.NotFoundError(w, r)
		return
	}

	// Находим ссылку по хэш параметру
	link, err := api.rep.GetLinkByHash(hash)
	if err != nil {
		api.log.Error(zap.Error(err))
		se.NotFoundError(w, r)
		return
	}

	render.JSON(w, r, link)
}
