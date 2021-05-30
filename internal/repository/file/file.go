package file

import (
	"github.com/lapitskyss/go_backend_1_project/internal/model"
	"github.com/lapitskyss/go_backend_1_project/internal/repository"
)

type fileStore struct {
}

func New() repository.Repository {
	return &fileStore{}
}

func (f *fileStore) Add(data *model.Link) (*model.Link, error) {
	return data, nil
}

func (f *fileStore) List(page int, limit int) ([]model.Link, error) {
	// TODO: implement
	return []model.Link{}, nil
}
