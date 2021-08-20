package repository

import "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"

type RedirectLogInterface interface {
	Add(redirectLog *model.RedirectLog) error
	CountRedirects(linkId uint64) (*uint64, error)
}
