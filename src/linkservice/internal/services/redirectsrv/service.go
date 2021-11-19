package redirectsrv

import "go.uber.org/zap"

type RedirectService struct {
	log *zap.Logger
	rs  RedirectStore
}

func NewRedirectService(log *zap.Logger, redirectStore RedirectStore) *RedirectService {
	return &RedirectService{
		log: log,
		rs:  redirectStore,
	}
}
