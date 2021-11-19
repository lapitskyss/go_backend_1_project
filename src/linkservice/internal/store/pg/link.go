package pg

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/linksrv"
)

type LinkStore struct {
	db *pgxpool.Pool
}

func NewLinkStore(store *Store) linksrv.LinkStore {
	return &LinkStore{
		db: store.Connection,
	}
}

func (l LinkStore) Add(ctx context.Context, link *linksrv.Link) error {
	_, err := l.db.Exec(ctx, "INSERT INTO links (id, url, hash, created_at) values($1, $2, $3, $4)",
		link.ID, link.URL, link.Hash, link.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (l LinkStore) GetNextId(ctx context.Context) (uint64, error) {
	var id uint64
	err := l.db.QueryRow(ctx, "SELECT NEXTVAL('links_id_seq')").Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (l LinkStore) GetByURL(ctx context.Context, url string) (*linksrv.Link, error) {
	var ln = &linksrv.Link{}

	query := "SELECT id, url, hash, created_at FROM links WHERE url = $1"
	err := l.db.QueryRow(ctx, query, url).Scan(&ln.ID, &ln.URL, &ln.Hash, &ln.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return ln, nil
}

func (l LinkStore) GetByHash(ctx context.Context, hash string) (*linksrv.Link, error) {
	var ln = &linksrv.Link{}

	query := "SELECT id, url, hash, created_at FROM links WHERE hash = $1"
	err := l.db.QueryRow(ctx, query, hash).Scan(&ln.ID, &ln.URL, &ln.Hash, &ln.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return ln, nil
}

func (l LinkStore) GetByHashWithStatistics(ctx context.Context, hash string) (*linksrv.StatisticForLink, error) {
	var ln = &linksrv.StatisticForLink{}

	query := "SELECT id, url, hash, created_at, (SELECT count(*) FROM redirects WHERE redirects.link_id = links.id) as num FROM links WHERE hash = $1"
	err := l.db.QueryRow(ctx, query, hash).Scan(&ln.ID, &ln.URL, &ln.Hash, &ln.CreatedAt, &ln.NumberOfRedirects)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return ln, nil
}

func (l LinkStore) GetByHashes(ctx context.Context, hashes []string) (<-chan linksrv.Link, <-chan error) {
	linkChannel := make(chan linksrv.Link, len(hashes))
	errorChannel := make(chan error, 1)

	go func() {
		defer close(linkChannel)
		defer close(errorChannel)

		query := "SELECT id, url, hash, created_at FROM links WHERE hash = any($1)"
		hs := &pgtype.TextArray{}
		err := hs.Set(hashes)
		if err != nil {
			errorChannel <- err
			return
		}

		rows, err := l.db.Query(ctx, query, hs)
		if err != nil {
			errorChannel <- err
			return
		}

		defer rows.Close()

		for rows.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				var res linksrv.Link
				err = rows.Scan(&res.ID, &res.URL, &res.Hash, &res.CreatedAt)
				if err != nil {
					errorChannel <- err
					return
				}

				linkChannel <- res
			}
		}
	}()

	return linkChannel, errorChannel
}

func (l LinkStore) FindBy(ctx context.Context, params linksrv.FindByParameters) (<-chan linksrv.Link, <-chan error) {
	linkChannel := make(chan linksrv.Link, params.Limit)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(linkChannel)
		defer close(errorChannel)

		var rows pgx.Rows
		var err error

		if params.Query != "" {
			q := "SELECT id, url, hash, created_at FROM links WHERE url ILIKE $1 ORDER BY " +
				params.Sort + " " + params.Order + " LIMIT $2 OFFSET $3"
			query := "%" + params.Query + "%"
			rows, err = l.db.Query(ctx, q, query, params.Limit, params.Offset)
		} else {
			q := "SELECT id, url, hash, created_at FROM links ORDER BY " +
				params.Sort + " " + params.Order + " LIMIT $1 OFFSET $2"
			rows, err = l.db.Query(ctx, q, params.Limit, params.Offset)
		}
		if err != nil {
			errorChannel <- err
			return
		}

		defer rows.Close()

		for rows.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				var res linksrv.Link
				err = rows.Scan(&res.ID, &res.URL, &res.Hash, &res.CreatedAt)
				if err != nil {
					errorChannel <- err
					return
				}

				linkChannel <- res
			}
		}
	}()

	return linkChannel, errorChannel
}

func (l LinkStore) CountByQuery(ctx context.Context, query string) (uint64, error) {
	var numberOfLinks uint64
	var err error

	if query != "" {
		q := "SELECT count(*) as total FROM links WHERE url ILIKE $1"
		err = l.db.QueryRow(ctx, q, query).Scan(&numberOfLinks)
	} else {
		q := "SELECT count(*) as total FROM links"
		err = l.db.QueryRow(ctx, q).Scan(&numberOfLinks)
	}
	if err != nil {
		return 0, err
	}

	return numberOfLinks, err
}
