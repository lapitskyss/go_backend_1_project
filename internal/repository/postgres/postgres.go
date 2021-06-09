package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Store struct {
	ctx    context.Context
	client *pgxpool.Pool
}

func (store *Store) Init(ctx context.Context) error {
	var err error

	store.ctx = ctx

	config := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	store.client, err = pgxpool.Connect(ctx, config)
	if err != nil {
		return err
	}

	return nil
}

func (store *Store) GetConn() interface{} {
	return store.client
}

func (store *Store) Close() {
	store.client.Close()
}
