package postgres

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/lapitskyss/go_backend_1_project/internal/model"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

func (store *Store) AddLink(link *model.Link) (*model.Link, error) {
	sql, args, err := psql.Insert("links").
		Columns("url", "hash", "created_at").
		Values(link.URL, link.Hash, link.CreatedAt).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = store.client.QueryRow(store.ctx, sql, args...).
		Scan(&link.URL, &link.Hash, &link.CreatedAt)

	if err == pgx.ErrNoRows {
		return link, nil
	}
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (store *Store) GetLinkByURL(url string) (*model.Link, error) {
	q, args, err := psql.Select("url, hash, created_at").
		From("links").
		Where(sq.Eq{"url": url}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var link model.Link

	err = store.client.QueryRow(store.ctx, q, args...).
		Scan(&link.URL, &link.Hash, &link.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (store *Store) GetLinkByHash(hash string) (*model.Link, error) {
	q, args, err := psql.Select("url, hash, created_at").
		From("links").
		Where(sq.Eq{"hash": hash}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var link model.Link

	err = store.client.QueryRow(store.ctx, q, args...).
		Scan(&link.URL, &link.Hash, &link.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (store *Store) GetLinksBy(page uint64, limit uint64, sort *string, order *string, query *string) (*[]*model.Link, error) {
	sb := psql.Select("url, hash, created_at").
		From("links").
		Limit(limit)

	if page > 1 {
		sb = sb.Offset((page - 1) * limit)
	}
	if sort != nil && order != nil {
		sb = sb.OrderBy(*sort + " " + strings.ToUpper(*order))
	}
	if sort != nil && order == nil {
		sb = sb.OrderBy(*sort + " ASC")
	}
	if sort == nil && order != nil {
		sb = sb.OrderBy("id " + strings.ToUpper(*order))
	}
	if query != nil {
		sb = sb.Where(sq.Like{"url": "%" + *query + "%"})
	}

	q, args, err := sb.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := store.client.Query(store.ctx, q, args...)
	if err != nil {
		return nil, err
	}

	links := []*model.Link{} // nolint errcheck

	for rows.Next() {
		var link model.Link

		err = rows.Scan(&link.URL, &link.Hash, &link.CreatedAt)
		if err != nil {
			return nil, err
		}

		links = append(links, &link)
	}

	return &links, nil
}

func (store *Store) CountAllLinks() (*uint64, error) {
	q, args, err := psql.Select("count(*) as n_links").From("links").ToSql()
	if err != nil {
		return nil, err
	}

	var total uint64
	err = store.client.QueryRow(store.ctx, q, args...).
		Scan(&total)
	if err != nil {
		return nil, err
	}

	return &total, err
}