package repository

import "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"

type Repository interface {
	AddLink(link *model.Link) (*model.Link, error)
	GetNextLinkId() (*uint64, error)
	GetExistingLink(url string) (*bool, *model.Link, error)
	GetLinkByHash(hash string) (*model.Link, error)
	FindLinks(page uint64, limit uint64, sort *string, order *string, query *string) (*[]*model.Link, error)
	CountLinksByQuery(query *string) (*uint64, error)
}
