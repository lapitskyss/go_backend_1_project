package linksrv

import (
	"context"
	"errors"
	"math"

	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/pkg/e"
)

type SearchLinkRequest struct {
	Page  *uint64 // min 1
	Limit *uint64 // max 100; min 1
	Sort  *string // url, hash, created_at
	Order *string // asc, desc
	Query *string
}

type SearchLinkResponse struct {
	Page  uint64
	Limit uint64
	Pages uint64
	Total uint64
	Links <-chan Link
	Err   <-chan error
}

func (s *LinkService) Search(ctx context.Context, r SearchLinkRequest) (*SearchLinkResponse, error) {
	err := r.validate()
	if err != nil {
		return nil, e.ErrBadRequest(err.Error())
	}

	page := r.getPage()
	limit := r.getLimit()
	sort := r.getSort()
	order := r.getOrder()
	query := r.getQuery()
	offset := (page - 1) * limit

	totalLinks, err := s.ls.CountByQuery(ctx, query)
	if err != nil {
		s.log.Error("Link service error to call CountByQuery", zap.Error(err))
		return nil, e.ErrInternal()
	}

	numberOfPages := getNumberOfPages(totalLinks, limit)

	linkChannel, errChannel := s.ls.FindBy(ctx, FindByParameters{
		Query:  query,
		Sort:   sort,
		Order:  order,
		Limit:  limit,
		Offset: offset,
	})

	return &SearchLinkResponse{
		Page:  page,
		Limit: limit,
		Pages: numberOfPages,
		Total: totalLinks,
		Links: linkChannel,
		Err:   errChannel,
	}, nil
}

func (r *SearchLinkRequest) validate() error {
	if r.Limit != nil {
		if *r.Limit <= 0 {
			return errors.New("limit can not be less or equal 0")
		}
		if *r.Limit > 100 {
			return errors.New("maximum limit is 100")
		}
	}
	if r.Page != nil && *r.Page <= 0 {
		return errors.New("page can not be less or equal 0")
	}
	if r.Sort != nil && !isValidSort(*r.Sort) {
		return errors.New("invalid sort value, available values: url, hash, created_at")
	}
	if r.Order != nil && !isValidOrder(*r.Order) {
		return errors.New("invalid order value, available values: asc, desc")
	}
	if r.Query != nil && len(*r.Query) > 1000 {
		return errors.New("query is to long")
	}

	return nil
}

func (r *SearchLinkRequest) getPage() uint64 {
	if r.Page != nil {
		return *r.Page
	}

	return 1
}

func (r *SearchLinkRequest) getLimit() uint64 {
	if r.Limit != nil {
		return *r.Limit
	}

	return 10
}

func (r *SearchLinkRequest) getSort() string {
	if r.Sort != nil {
		return *r.Sort
	}

	return "created_at"
}

func (r *SearchLinkRequest) getOrder() string {
	if r.Order != nil {
		return *r.Order
	}

	return "asc"
}

func (r *SearchLinkRequest) getQuery() string {
	if r.Query != nil {
		return *r.Query
	}

	return ""
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

func getNumberOfPages(total, limit uint64) uint64 {
	if total == 0 {
		return 1
	}
	d := float64(total) / float64(limit)
	return uint64(math.Ceil(d))
}
