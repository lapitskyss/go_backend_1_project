package link

import (
	"errors"
	"math"
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/pointer"
	"go.uber.org/zap"

	se "github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/server_errors"
)

type listParameters struct {
	Page  *uint64 `schema:"page"`
	Limit *uint64 `schema:"limit"`
	Sort  *string `schema:"sort"`  // url, hash, created_at
	Order *string `schema:"order"` // asc, desc
	Query *string `schema:"query"`
}

type linksList struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
	Pages uint64 `json:"pages"`
	Total uint64 `json:"total"`

	Link *[]*model.Link `json:"links"`
}

func (api *linkController) List(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры с запроса
	params, err := getQueryParams(r)
	if err != nil {
		se.BadRequestError(w, r, err)
		return
	}

	// Валидируем параметры
	err = validateQueryParams(params)
	if err != nil {
		se.BadRequestError(w, r, err)
		return
	}

	// Получаем слайс ссылок по заданным параметрам
	links, err := api.rep.FindLinks(*params.Page, *params.Limit, params.Sort, params.Order, params.Query)
	if err != nil {
		api.log.Error(zap.Error(err))
		se.NotFoundError(w, r)
		return
	}

	// Определяем количество ссылок в базе
	totalLinks, err := api.rep.CountLinksByQuery(params.Query)
	if err != nil {
		api.log.Error(zap.Error(err))
		se.NotFoundError(w, r)
		return
	}

	render.JSON(w, r, linksList{
		Page:  *params.Page,
		Limit: *params.Limit,
		Pages: getPages(*totalLinks, *params.Limit),
		Total: *totalLinks,
		Link:  links,
	})
}

func getQueryParams(r *http.Request) (*listParameters, error) {
	var params listParameters
	decoder := schema.NewDecoder()
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		return nil, err
	}

	if params.Page == nil {
		params.Page = pointer.Uint64(1)
	}
	if params.Limit == nil {
		params.Limit = pointer.Uint64(10)
	}

	return &params, nil
}

func validateQueryParams(params *listParameters) error {
	if params.Limit != nil && *params.Limit > 100 {
		return errors.New("maximum limit is 100")
	}
	if params.Limit != nil && *params.Limit <= 0 {
		return errors.New("limit can not be less or 0")
	}
	if params.Page != nil && *params.Page <= 0 {
		return errors.New("page can not be less or 0")
	}
	if params.Sort != nil && !isValidSort(*params.Sort) {
		return errors.New("invalid sort value, available values: url, hash, created_at")
	}
	if params.Order != nil && !isValidOrder(*params.Order) {
		return errors.New("invalid order value, available values: asc, desc")
	}
	if params.Query != nil && len(*params.Query) > 1000 {
		return errors.New("query is to long")
	}

	return nil
}

func isValidSort(sort string) bool {
	switch sort {
	case
		"url",
		"hash",
		"created_at":
		return true
	}
	return false
}

func isValidOrder(order string) bool {
	switch order {
	case
		"asc",
		"ASC",
		"desc",
		"DESC":
		return true
	}
	return false
}

func getPages(total, limit uint64) uint64 {
	if total == 0 {
		return 1
	}
	d := float64(total) / float64(limit)
	return uint64(math.Ceil(d))
}
