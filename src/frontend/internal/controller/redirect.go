package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

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
			c.log.Error("Frontend link not found", zap.Error(err))
		}

		// TODO show home page with error

		c.Home(w, r)
		return
	}

	if link == nil {
		c.log.Error("link url not provided")
		c.Home(w, r)
		return
	}

	http.Redirect(w, r, link.GetUrl(), 301)
	return
}
