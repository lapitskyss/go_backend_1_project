package controller

import (
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/internal/repository/postgres"
)

type linkController struct {
	log *zap.SugaredLogger
	rep *postgres.Store
}

func NewLinkController(log *zap.SugaredLogger, rep *postgres.Store) *linkController {
	return &linkController{log, rep}
}

func (api *linkController) Add(w http.ResponseWriter, r *http.Request) {
	//newLink := &model.Link{
	//	ShortLink:    "Short link",
	//	RedirectLink: "Redirect link",
	//	CreatedAt:    "created at",
	//}
	//createdLink, err := api.rep.Add(newLink)
	//
	//if err != nil {
	//	se.RenderBadRequestError(w, r, err)
	//	return
	//}
	//
	//render.JSON(w, r, struct {
	//	Status bool        `json:"status"`
	//	Link   *model.Link `json:"link"`
	//}{true, createdLink})
}

func (api *linkController) List(w http.ResponseWriter, r *http.Request) {
	//links, err := api.rep.List(1, 100)
	//if err != nil {
	//	se.RenderBadRequestError(w, r, err)
	//	return
	//}

	render.JSON(w, r, struct {
		Status bool `json:"status"`
	}{true})
}
