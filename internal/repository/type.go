package repository

import "github.com/lapitskyss/go_backend_1_project/internal/model"

type Repository interface {
	Add(data *model.Link) (*model.Link, error)
	List(page int, limit int) ([]model.Link, error)
}
