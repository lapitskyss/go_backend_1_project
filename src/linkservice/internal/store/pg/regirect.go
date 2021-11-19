package pg

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/redirectsrv"
)

type RedirectStore struct {
	db *pgxpool.Pool
}

func NewRedirectStore(store *Store) redirectsrv.RedirectStore {
	return &RedirectStore{
		db: store.Connection,
	}
}

func (r RedirectStore) GetRedirectLink(ctx context.Context, hash string) (*redirectsrv.RedirectLink, error) {
	var redirect = &redirectsrv.RedirectLink{}

	query := "SELECT id, url FROM links WHERE hash = $1"
	err := r.db.QueryRow(ctx, query, hash).Scan(&redirect.ID, &redirect.URL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return redirect, nil
}

func (r RedirectStore) SaveRedirectStatistics(ctx context.Context, redirect redirectsrv.Redirect) error {
	_, err := r.db.Exec(ctx, "INSERT INTO redirects (link_id) values($1)", redirect.LinkID)
	if err != nil {
		return err
	}

	return nil
}
