package handler

import (
	"net/http"

	"go.uber.org/zap"
)

type homePage struct {
	API string
}

func (h *Handler) Home(w http.ResponseWriter, _ *http.Request) {
	err := h.tmp.HomeTemplate.Execute(w, homePage{
		API: h.config.API,
	})
	if err != nil {
		h.log.Error("Frontend error to show home page", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
