package link

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/render"
	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/internal/model"
	se "github.com/lapitskyss/go_backend_1_project/pkg/server_errors"
)

type createLinkParams struct {
	URL string `json:"url"`
}

func (api *linkController) Create(w http.ResponseWriter, r *http.Request) {
	params := &createLinkParams{}

	// Декодим тело запроса
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		se.BadRequestError(w, r, err)
		return
	}

	// Валидируем пришедшии параметры
	err = validateParams(params)
	if err != nil {
		se.BadRequestError(w, r, err)
		return
	}

	// Проверяем что короткая ссылка уже есть для URL
	isExist, existingLink, err := api.rep.GetExistingLink(params.URL)
	if err != nil {
		api.log.Error(zap.Error(err))
		se.BadRequestError(w, r, err)
		return
	}

	// Если короткая ссылка есть, то отдаем ее
	if *isExist {
		render.JSON(w, r, existingLink)
		return
	}

	// Получаем id короткой ссылки
	nextId, err := api.rep.GetNextLinkId()
	if err != nil {
		api.log.Error(zap.Error(err))
		se.BadRequestError(w, r, err)
		return
	}

	// Добавляем хэш и дату создания короткой ссылки
	link := initLink(params, nextId)

	// Создаем ссылку
	link, err = api.rep.AddLink(link)
	if err != nil {
		api.log.Error(zap.Error(err))
		se.BadRequestError(w, r, err)
		return
	}

	render.JSON(w, r, link)
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

func initLink(params *createLinkParams, id *uint64) *model.Link {
	return &model.Link{
		ID:        *id,
		URL:       params.URL,
		Hash:      getHash(id),
		CreatedAt: time.Now(),
	}
}

func getHash(id *uint64) string {
	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{int64(*id)})

	return e
}
