package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/pkg/api"
)

func (c *controller) Redirect(w http.ResponseWriter, r *http.Request) {
	var hash = chi.URLParam(r, "hash")
	if hash == "" {
		http.Error(w, "Hash not provided!", http.StatusBadRequest)
		return
	}

	link, err := c.client.Link.GetLinkByHash(c.ctx, hash)
	if err != nil {
		if err != api.ErrLinkNotFound {
			c.log.Error(zap.Error(err))
		}

		c.Home(w, r)
		return
	}

	if link == nil || link.URL == nil {
		c.log.Error(zap.Error(errors.New("link url not provided")))
		c.Home(w, r)
		return
	}

	// TODO: add redirect statistic

	http.Redirect(w, r, *link.URL, 301)
	return
}
