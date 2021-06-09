package link

import (
	"encoding/json"
	"errors"
	"hash/crc64"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/render"
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
	existingLink, err := api.rep.GetLinkByURL(params.URL)
	if err != nil {
		api.log.Error(zap.Error(err))
		se.BadRequestError(w, r, err)
		return
	}

	// Если короткая ссылка есть, то отдаем ее
	if existingLink != nil {
		render.JSON(w, r, existingLink)
		return
	}

	// Добавляем хэш и дату создания короткой ссылки
	link := initLink(params)

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

func initLink(params *createLinkParams) *model.Link {
	return &model.Link{
		URL:       params.URL,
		Hash:      getHash(params.URL),
		CreatedAt: time.Now(),
	}
}

// TODO: change function
func getHash(url string) string {
	crc64Table := crc64.MakeTable(0xC96C5795D7870F42)
	crc64Int := crc64.Checksum([]byte(url), crc64Table)
	return strconv.FormatUint(crc64Int, 16)
}
