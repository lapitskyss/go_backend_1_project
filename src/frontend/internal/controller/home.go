package controller

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func (c *controller) Home(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		c.log.Error(zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templatePath := filepath.Join(wd, "./templates/index.html")

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		c.log.Error(zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		c.log.Error(zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
