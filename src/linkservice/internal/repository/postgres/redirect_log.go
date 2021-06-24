package postgres

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
)

type RedirectLogRepository repository

func (rl *RedirectLogRepository) Add(redirectLog *model.RedirectLog) error {
	sql, args, err := psql.Insert("redirects_logs").
		Columns("link_id", "user_agent", "created_at").
		Values(redirectLog.LinkId, redirectLog.UserAgent, redirectLog.CreatedAt).
		ToSql()

	if err != nil {
		return err
	}

	_, err = rl.store.client.Exec(rl.store.ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (rl *RedirectLogRepository) CountRedirects(linkId uint64) (*uint64, error) {
	q, args, err := psql.Select("count(*) as t0").
		From("redirects_logs").
		Where(sq.Eq{"link_id": linkId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var total uint64
	err = rl.store.client.QueryRow(rl.store.ctx, q, args...).
		Scan(&total)
	if err != nil {
		return nil, err
	}

	return &total, nil
}
