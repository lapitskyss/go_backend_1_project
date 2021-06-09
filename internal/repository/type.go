package repository

import "github.com/lapitskyss/go_backend_1_project/internal/model"

type Repository interface {
	AddLink(link *model.Link) (*model.Link, error)
	GetLink(hash string) (*model.Link, error)
	GetLinksBy(page uint64, limit uint64, sort *string, order *string, query *string) (*[]*model.Link, error)
	CountAllLinks() (*uint64, error)
}
