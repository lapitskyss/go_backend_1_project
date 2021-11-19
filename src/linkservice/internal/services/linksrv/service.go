package linksrv

import "go.uber.org/zap"

type LinkService struct {
	log *zap.Logger
	ls  LinkStore
}

func NewLinkService(log *zap.Logger, linksStore LinkStore) *LinkService {
	return &LinkService{
		log: log,
		ls:  linksStore,
	}
}
