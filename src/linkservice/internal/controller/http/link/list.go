package link

import (
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/openlyinc/pointy"

	"github.com/gorilla/schema"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/repository"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/render"
)

type listParameters struct {
	Page        *uint64   `schema:"page"`
	Limit       *uint64   `schema:"limit"`
	Sort        *string   `schema:"sort"`  // url, hash, created_at
	Order       *string   `schema:"order"` // asc, desc
	Query       *string   `schema:"query"`
	Hashes      *string   `schema:"ids"`
	HashesSlice *[]string `schema:"-"`
}

type linksList struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
	Pages uint64 `json:"pages"`
	Total uint64 `json:"total"`

	Links []*model.Link `json:"links"`
}

func (lc *linkController) List(w http.ResponseWriter, r *http.Request) {
	// Получаем провалидированные параметры с запроса
	params, err := getQueryParams(r)
	if err != nil {
		render.BadRequestError(w, err)
		return
	}

	// Находим список ссылок по слайсу хэшов
	if params.HashesSlice != nil {
		links := []*model.Link{}
		links, err = lc.rep.Link().GetByHashes(params.HashesSlice)
		if err != nil {
			lc.log.Error(err)
			render.InternalServerError(w)
			return
		}

		render.Success(w, links)
		return
	}

	// Получаем слайс ссылок по заданным параметрам
	links, err := lc.rep.Link().FindBy(&repository.FindByParameters{
		Page:  *params.Page,
		Limit: *params.Limit,
		Sort:  params.Sort,
		Order: params.Order,
		Query: params.Query,
	})
	if err != nil {
		lc.log.Error(err)
		render.InternalServerError(w)
		return
	}

	// Определяем количество ссылок в базе
	totalLinks, err := lc.rep.Link().CountByQuery(params.Query)
	if err != nil {
		lc.log.Error(err)
		render.InternalServerError(w)
		return
	}

	render.Success(w, linksList{
		Page:  *params.Page,
		Limit: *params.Limit,
		Pages: getPages(*totalLinks, *params.Limit),
		Total: *totalLinks,
		Links: links,
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
		params.Page = pointy.Uint64(1)
	}
	if params.Limit == nil {
		params.Limit = pointy.Uint64(10)
	}

	err = validateQueryParams(&params)
	if err != nil {
		return nil, err
	}

	if params.Hashes != nil {
		err := validateQueryHashes(&params)
		if err != nil {
			return nil, err
		}
	}

	return &params, nil
}

func validateQueryHashes(params *listParameters) error {
	if len(*params.Hashes) > 2000 {
		return errors.New("maximum limit for URLs is 100")
	}
	hashesSlice := strings.Split(*params.Hashes, ",")
	hashesSliceLength := len(hashesSlice)

	if hashesSliceLength > 100 {
		return errors.New("maximum limit for URLs is 100")
	}
	params.HashesSlice = &hashesSlice

	return nil
}

func validateQueryParams(params *listParameters) error {
	if params.Limit != nil {
		if *params.Limit <= 0 {
			return errors.New("limit can not be less or equal 0")
		}
		if *params.Limit > 100 {
			return errors.New("maximum limit is 100")
		}
	}
	if params.Page != nil && *params.Page <= 0 {
		return errors.New("page can not be less or equal 0")
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
