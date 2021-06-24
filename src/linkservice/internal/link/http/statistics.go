package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	se "github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/server_errors"
)

type statisticsInfo struct {
	Redirects uint64 `json:"redirects"`
}

func (lc *linkController) Statistics(w http.ResponseWriter, r *http.Request) {
	var hash = chi.URLParam(r, "hash")
	if hash == "" {
		se.NotFoundError(w, r)
		return
	}

	// Находим ссылку по хэш параметру
	link, err := lc.rep.Link.GetByHash(hash)
	if err != nil {
		lc.log.Error(err)
		se.NotFoundError(w, r)
		return
	}
	if link == nil {
		se.NotFoundError(w, r)
		return
	}

	total, err := lc.rep.RedirectLog.CountRedirects(link.ID)
	if err != nil {
		lc.log.Error(err)
		se.InternalServerError(w, r)
		return
	}

	render.JSON(w, r, statisticsInfo{Redirects: *total})
}
