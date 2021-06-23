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

	common repository

	Link *LinkRepository
}

type repository struct {
	store *Store
}

func NewStore(ctx context.Context) (*Store, error) {
	config := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	client, err := pgxpool.Connect(ctx, config)
	if err != nil {
		return nil, err
	}

	c := &Store{ctx: ctx, client: client}
	c.common.store = c
	c.Link = (*LinkRepository)(&c.common)

	return c, nil
}

func (store *Store) GetConn() interface{} {
	return store.client
}

func (store *Store) Close() {
	store.client.Close()
}
