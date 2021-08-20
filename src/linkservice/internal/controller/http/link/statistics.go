package link

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/render"
)

type statisticsInfo struct {
	Redirects uint64 `json:"redirects"`
}

func (lc *linkController) Statistics(w http.ResponseWriter, r *http.Request) {
	var hash = chi.URLParam(r, "hash")
	if hash == "" {
		render.NotFoundError(w)
		return
	}

	// Находим ссылку по хэш параметру
	link, err := lc.rep.Link().GetByHash(hash)
	if err != nil {
		lc.log.Error(err)
		render.InternalServerError(w)
		return
	}
	if link == nil {
		render.NotFoundError(w)
		return
	}

	total, err := lc.rep.RedirectLog().CountRedirects(link.ID)
	if err != nil {
		lc.log.Error(err)
		render.InternalServerError(w)
		return
	}

	render.Success(w, statisticsInfo{Redirects: *total})
}
