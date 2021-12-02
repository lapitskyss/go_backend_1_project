package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/pkg/e"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/pkg/render"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/linksrv"
)

type LinkHandler struct {
	log *zap.Logger
	ls  *linksrv.LinkService
}

func NewLinkHandler(log *zap.Logger, linkService *linksrv.LinkService) *LinkHandler {
	return &LinkHandler{
		log: log,
		ls:  linkService,
	}
}

type createLinkRequest struct {
	URL string `json:"url"`
}

type singleLinkResponse struct {
	URL       string    `json:"url"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var linkRequest = &createLinkRequest{}
	err := json.NewDecoder(r.Body).Decode(linkRequest)
	if err != nil {
		render.BadRequestError(w, errors.New("incorrect request params"))
		return
	}

	l, err := h.ls.Create(r.Context(), linksrv.CreateLinkRequest{
		URL: linkRequest.URL,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	render.Success(w, singleLinkResponse{
		URL:       l.URL,
		Hash:      l.Hash,
		CreatedAt: l.CreatedAt,
	})
}

type listLinkParameters struct {
	Hashes string `schema:"ids"`
}

func (h *LinkHandler) List(w http.ResponseWriter, r *http.Request) {
	var params listLinkParameters
	err := schema.NewDecoder().Decode(&params, r.URL.Query())
	if err != nil {
		render.BadRequestError(w, errors.New("unexpected query parameter"))
		return
	}

	hashesSlice := strings.Split(params.Hashes, ",")
	hashesSliceLength := len(hashesSlice)

	links, errChan := h.ls.List(r.Context(), hashesSlice)

	respLinks := make([]singleLinkResponse, 0, hashesSliceLength)

	for l := range links {
		respLinks = append(respLinks, singleLinkResponse{
			URL:       l.URL,
			Hash:      l.Hash,
			CreatedAt: l.CreatedAt,
		})
	}

	select {
	case <-r.Context().Done():
		return
	case message, ok := <-errChan:
		if ok {
			h.log.Error("Link list handler error from link service", zap.Error(message))
			sendError(w, e.ErrInternal())
			return
		}
	}

	render.Success(w, respLinks)
}

type infoLinkResponse struct {
	URL       string    `json:"url"`
	Hash      string    `json:"hash"`
	Redirects uint64    `json:"redirects"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *LinkHandler) Info(w http.ResponseWriter, r *http.Request) {
	var hash = chi.URLParam(r, "hash")

	l, err := h.ls.Info(r.Context(), hash)
	if err != nil {
		sendError(w, err)
		return
	}

	render.Success(w, infoLinkResponse{
		URL:       l.URL,
		Hash:      l.Hash,
		Redirects: l.NumberOfRedirects,
		CreatedAt: l.CreatedAt,
	})
}

type searchLinkParameters struct {
	Page  *uint64 `schema:"page"`
	Limit *uint64 `schema:"limit"`
	Sort  *string `schema:"sort"`
	Order *string `schema:"order"`
	Query *string `schema:"query"`
}

type searchLinkResponse struct {
	Page  uint64               `json:"page"`
	Limit uint64               `json:"limit"`
	Pages uint64               `json:"pages"`
	Total uint64               `json:"total"`
	Links []singleLinkResponse `json:"links"`
}

func (h *LinkHandler) Search(w http.ResponseWriter, r *http.Request) {
	var params searchLinkParameters
	err := schema.NewDecoder().Decode(&params, r.URL.Query())
	if err != nil {
		render.BadRequestError(w, errors.New("unexpected query parameter"))
		return
	}

	links, err := h.ls.Search(r.Context(), linksrv.SearchLinkRequest{
		Page:  params.Page,
		Limit: params.Limit,
		Sort:  params.Sort,
		Order: params.Order,
		Query: params.Query,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	respLinks := make([]singleLinkResponse, 0, links.Limit)

	for l := range links.Links {
		respLinks = append(respLinks, singleLinkResponse{
			URL:       l.URL,
			Hash:      l.Hash,
			CreatedAt: l.CreatedAt,
		})
	}

	select {
	case <-r.Context().Done():
		return
	case message, ok := <-links.Err:
		if ok {
			h.log.Error("Link search handler error from link service", zap.Error(message))
			sendError(w, e.ErrInternal())
			return
		}
	default:
	}

	render.Success(w, &searchLinkResponse{
		Page:  links.Page,
		Limit: links.Limit,
		Pages: links.Pages,
		Total: links.Total,
		Links: respLinks,
	})
}

func sendError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *e.BadRequestError:
		render.BadRequestError(w, err)
	case *e.NotFoundError:
		render.NotFoundError(w)
	default:
		render.InternalServerError(w)
	}
}
