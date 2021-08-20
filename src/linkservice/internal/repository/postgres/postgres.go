package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	r "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/repository"
)

type Store struct {
	ctx    context.Context
	client *pgxpool.Pool

	common repository

	link        r.LinkInterface
	redirectLog r.RedirectLogInterface
}

type repository struct {
	store *Store
}

func (store *Store) Connect(ctx context.Context) error {
	config := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	client, err := pgxpool.Connect(ctx, config)
	if err != nil {
		return err
	}

	store.ctx = ctx
	store.client = client

	store.common.store = store
	store.link = (*LinkRepository)(&store.common)
	store.redirectLog = (*RedirectLogRepository)(&store.common)

	return nil
}

func (store *Store) GetConnection() interface{} {
	return store.client
}

func (store *Store) CloseConnection() {
	store.client.Close()
}

func (store *Store) Link() r.LinkInterface {
	return store.link
}

func (store *Store) RedirectLog() r.RedirectLogInterface {
	return store.redirectLog
}
