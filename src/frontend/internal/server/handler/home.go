package handler

import (
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) Home(w http.ResponseWriter, _ *http.Request) {
	err := h.tmp.HomeTemplate.Execute(w, nil)
	if err != nil {
		h.log.Error("Frontend error to show home page", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
