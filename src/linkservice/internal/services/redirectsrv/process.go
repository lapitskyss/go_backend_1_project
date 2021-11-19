package redirectsrv

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrLinkNotFound = errors.New("not found")
	ErrInternal     = errors.New("internal error")
)

func (s *RedirectService) Process(ctx context.Context, hash string) (string, error) {
	if hash == "" || len(hash) > 20 {
		return "", ErrLinkNotFound
	}

	rl, err := s.rs.GetRedirectLink(ctx, hash)
	if err != nil {
		s.log.Error("Redirect service error to call GetRedirectLink", zap.Error(err))
		return "", ErrInternal
	}

	if rl == nil {
		return "", ErrLinkNotFound
	}

	go func(id uint64) {
		sErr := s.rs.SaveRedirectStatistics(context.Background(), Redirect{
			LinkID: id,
		})
		if sErr != nil {
			s.log.Error("Redirect service error to call SaveRedirectStatistics", zap.Error(err))
		}
	}(rl.ID)

	return rl.URL, nil
}
