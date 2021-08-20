package link

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/speps/go-hashids/v2"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/render"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/util"
)

type createLinkParams struct {
	URL string `json:"url"`
}

func (lc *linkController) Create(w http.ResponseWriter, r *http.Request) {
	// Декодим тело запроса
	params := &createLinkParams{}
	err := util.DecodeJSON(r, params)
	if err != nil {
		render.BadRequestError(w, err)
		return
	}

	// Валидируем пришедшии параметры
	err = validateParams(params)
	if err != nil {
		render.BadRequestError(w, err)
		return
	}

	// Проверяем что короткая ссылка уже есть для URL
	existingLink, err := lc.rep.Link().GetByURL(params.URL)
	if err != nil {
		lc.log.Error(err)
		render.InternalServerError(w)
		return
	}

	// Если короткая ссылка есть, то отдаем ее
	if existingLink != nil {
		render.Success(w, existingLink)
		return
	}

	// Получаем id короткой ссылки
	nextId, err := lc.rep.Link().GetNextId()
	if err != nil {
		lc.log.Error(err)
		render.InternalServerError(w)
		return
	}

	// Создаем ссылку
	link, err := lc.rep.Link().Add(&model.Link{
		ID:        *nextId,
		URL:       params.URL,
		Hash:      getHash(nextId),
		CreatedAt: time.Now(),
	})
	if err != nil {
		lc.log.Error(err)
		render.InternalServerError(w)
		return
	}

	render.Success(w, link)
}

func validateParams(params *createLinkParams) error {
	result, err := url.ParseRequestURI(params.URL)

	if err != nil || result.Scheme == "" {
		return errors.New("incorrect URL")
	}

	if len(params.URL) > 10000 {
		return errors.New("URL is to long")
	}

	return nil
}

func getHash(id *uint64) string {
	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{int64(*id)})

	return e
}
