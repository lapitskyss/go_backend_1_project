package controller

import (
	"github.com/go-chi/render"

	"go.uber.org/zap"
	"net/http"

	"github.com/lapitskyss/go_backend_1_project/internal/model"
	"github.com/lapitskyss/go_backend_1_project/internal/repository"
	se "github.com/lapitskyss/go_backend_1_project/pkg/serverErrors"
)

type linkController struct {
	log *zap.SugaredLogger
	rep repository.Repository
}

func NewLinkController(log *zap.SugaredLogger, rep repository.Repository) *linkController {
	return &linkController{log, rep}
}

func (api *linkController) Add(w http.ResponseWriter, r *http.Request) {
	newLink := &model.Link{
		ShortLink:    "Short link",
		RedirectLink: "Redirect link",
		CreatedAt:    "created at",
	}
	createdLink, err := api.rep.Add(newLink)

	if err != nil {
		se.RenderBadRequestError(w, r, err)
		return
	}

	render.JSON(w, r, struct {
		Status bool        `json:"status"`
		Link   *model.Link `json:"link"`
	}{true, createdLink})
}

func (api *linkController) List(w http.ResponseWriter, r *http.Request) {
	links, err := api.rep.List(1, 100)
	if err != nil {
		se.RenderBadRequestError(w, r, err)
		return
	}

	render.JSON(w, r, links)
}
