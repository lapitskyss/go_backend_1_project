package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/rpc"
)

func (c *controller) Redirect(w http.ResponseWriter, r *http.Request) {
	var hash = chi.URLParam(r, "hash")
	if hash == "" {
		http.Error(w, "Hash not provided!", http.StatusBadRequest)
		return
	}

	link, err := c.fe.GetLink(c.ctx, hash)
	if err != nil {
		if err != rpc.ErrLinkNotFound {
			c.log.Error(err)
		}

		c.Home(w, r)
		return
	}

	if link == nil {
		c.log.Error("link url not provided")
		c.Home(w, r)
		return
	}

	// TODO: add redirect statistic

	http.Redirect(w, r, link.Url, 301)
	return
}
