package repository

import "github.com/lapitskyss/go_backend_1_project/internal/model"

type Repository interface {
	AddLink(link *model.Link) (*model.Link, error)
	GetNextLinkId() (*uint64, error)
	GetExistingLink(url string) (*bool, *model.Link, error)
	GetLinkByHash(hash string) (*model.Link, error)
	FindLinks(page uint64, limit uint64, sort *string, order *string, query *string) (*[]*model.Link, error)
	CountAllLinks() (*uint64, error)
}
