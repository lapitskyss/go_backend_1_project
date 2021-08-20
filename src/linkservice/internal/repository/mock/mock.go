//go:generate mockgen -destination=mock_link.go -package=mock github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/repository LinkInterface
//go:generate mockgen -destination=mock_redirect_log.go -package=mock github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/repository RedirectLogInterface

package mock

import (
	"context"

	r "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/repository"
)

type Store struct {
	link        r.LinkInterface
	redirectLog r.RedirectLogInterface
}

func (store *Store) Connect(_ context.Context) error {
	return nil
}

func (store *Store) GetConnection() interface{} {
	return nil
}

func (store *Store) CloseConnection() {
	return
}

func (store *Store) Link() r.LinkInterface {
	return store.link
}

func (store *Store) SetLink(linkInterface r.LinkInterface) {
	store.link = linkInterface
}

func (store *Store) RedirectLog() r.RedirectLogInterface {
	return store.redirectLog
}

func (store *Store) SetRedirectLog(redirectLogInterface r.RedirectLogInterface) {
	store.redirectLog = redirectLogInterface
}
