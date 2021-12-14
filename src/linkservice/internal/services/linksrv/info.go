package linksrv

import (
	"context"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/pkg/response"
)

func (s *LinkService) Info(ctx context.Context, hash string) (*StatisticForLink, error) {
	if hash == "" || len(hash) > 20 {
		return nil, response.ErrNotFound()
	}

	return s.ls.GetByHashWithStatistics(ctx, hash)
}
