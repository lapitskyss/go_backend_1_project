package linksrv

import (
	"context"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/pkg/e"
)

func (s *LinkService) Info(ctx context.Context, hash string) (*StatisticForLink, error) {
	if hash == "" || len(hash) > 20 {
		return nil, e.ErrNotFound()
	}

	return s.ls.GetByHashWithStatistics(ctx, hash)
}
