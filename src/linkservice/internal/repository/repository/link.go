package repository

import "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"

type LinkInterface interface {
	Add(link *model.Link) (*model.Link, error)
	GetNextId() (*uint64, error)
	GetByURL(url string) (*model.Link, error)
	GetByHash(hash string) (*model.Link, error)
	GetByHashes(hashes *[]string) ([]*model.Link, error)
	FindBy(params *FindByParameters) ([]*model.Link, error)
	CountByQuery(query *string) (*uint64, error)
}

type FindByParameters struct {
	Page  uint64
	Limit uint64
	Sort  *string // url, hash, created_at
	Order *string // asc, desc
	Query *string
}
