package controller

import (
	"net/http"

	"go.uber.org/zap"
)

func (c *controller) Home(w http.ResponseWriter, r *http.Request) {
	err := c.tmp.HomeTemplate.Execute(w, nil)
	if err != nil {
		c.log.Error("Frontend error to show home page", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
