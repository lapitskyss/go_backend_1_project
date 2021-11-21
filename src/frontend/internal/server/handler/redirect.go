package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var ErrLinkNotFound = errors.New("link not found")

type LinkService interface {
	GetRedirectUrl(ctx context.Context, hash string) (string, error)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	var hash = chi.URLParam(r, "hash")
	if hash == "" {
		http.Error(w, "Hash not provided!", http.StatusBadRequest)
		return
	}

	link, err := h.linkSrv.GetRedirectUrl(r.Context(), hash)
	if err != nil {
		if err == ErrLinkNotFound {
			h.Home(w, r)
			return
		}

		h.log.Error("Frontend error to call GetRedirectUrl", zap.Error(err))
		http.Error(w, "Unexpected error occurred. Please try again later.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, link, 301)
	return
}
