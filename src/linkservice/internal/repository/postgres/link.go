package postgres

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/pointer"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type LinkRepository repository

func (lr *LinkRepository) Add(link *model.Link) (*model.Link, error) {
	sql, args, err := psql.Insert("links").
		Columns("id", "url", "hash", "created_at").
		Values(link.ID, link.URL, link.Hash, link.CreatedAt).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = lr.store.client.QueryRow(lr.store.ctx, sql, args...).
		Scan(&link.URL, &link.Hash, &link.CreatedAt)

	if err == pgx.ErrNoRows {
		return link, nil
	}
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (lr *LinkRepository) GetNextId() (*uint64, error) {
	var id uint64

	err := lr.store.client.QueryRow(lr.store.ctx, "SELECT NEXTVAL('links_id_seq')").
		Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (lr *LinkRepository) GetByURL(url string) (*bool, *model.Link, error) {
	q, args, err := psql.Select("url, hash, created_at").
		From("links").
		Where(sq.Eq{"url": url}).
		ToSql()

	if err != nil {
		return nil, nil, err
	}

	var link model.Link

	err = lr.store.client.QueryRow(lr.store.ctx, q, args...).
		Scan(&link.URL, &link.Hash, &link.CreatedAt)

	if err == pgx.ErrNoRows {
		return pointer.Bool(false), nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	return pointer.Bool(true), &link, nil
}

func (lr *LinkRepository) GetByHash(hash string) (*model.Link, error) {
	q, args, err := psql.Select("url, hash, created_at").
		From("links").
		Where(sq.Eq{"hash": hash}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var link model.Link

	err = lr.store.client.QueryRow(lr.store.ctx, q, args...).
		Scan(&link.URL, &link.Hash, &link.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &link, nil
}

func (lr *LinkRepository) GetByHashes(hashes *[]string) (*[]*model.Link, error) {
	q, args, err := psql.Select("url, hash, created_at").
		From("links").
		Where(sq.Eq{"hash": *hashes}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := lr.store.client.Query(lr.store.ctx, q, args...)
	if err != nil {
		return nil, err
	}

	var links []*model.Link

	for rows.Next() {
		var link model.Link

		err = rows.Scan(&link.URL, &link.Hash, &link.CreatedAt)
		if err != nil {
			return nil, err
		}

		links = append(links, &link)
	}

	resultedLinks := []*model.Link{} // nolint errcheck
	for _, hash := range *hashes {
		for _, link := range links {
			if link.Hash == hash {
				resultedLinks = append(resultedLinks, link)
			}
		}
	}

	return &resultedLinks, nil
}

type FindByParameters struct {
	Page  uint64
	Limit uint64
	Sort  *string // url, hash, created_at
	Order *string // asc, desc
	Query *string
}

func (lr *LinkRepository) FindBy(params *FindByParameters) (*[]*model.Link, error) {
	sb := psql.Select("url, hash, created_at").
		From("links").
		Limit(params.Limit)

	if params.Page > 1 {
		sb = sb.Offset((params.Page - 1) * params.Limit)
	}
	if params.Sort != nil && params.Order != nil {
		sb = sb.OrderBy(*params.Sort + " " + strings.ToUpper(*params.Order))
	}
	if params.Sort != nil && params.Order == nil {
		sb = sb.OrderBy(*params.Sort + " ASC")
	}
	if params.Sort == nil && params.Order != nil {
		sb = sb.OrderBy("id " + strings.ToUpper(*params.Order))
	}
	if params.Query != nil {
		sb = sb.Where(sq.Like{"url": "%" + *params.Query + "%"})
	}

	q, args, err := sb.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := lr.store.client.Query(lr.store.ctx, q, args...)
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

func (lr *LinkRepository) CountByQuery(query *string) (*uint64, error) {
	sb := psql.Select("count(*) as n_links").From("links")
	if query != nil {
		sb = sb.Where(sq.Like{"url": "%" + *query + "%"})
	}
	q, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var total uint64
	err = lr.store.client.QueryRow(lr.store.ctx, q, args...).
		Scan(&total)
	if err != nil {
		return nil, err
	}

	return &total, err
}
