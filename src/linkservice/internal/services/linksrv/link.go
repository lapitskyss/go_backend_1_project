package linksrv

import (
	"context"
	"time"
)

type Link struct {
	ID        uint64
	URL       string
	Hash      string
	CreatedAt time.Time
}

type FindByParameters struct {
	Query  string
	Sort   string // url, hash, created_at
	Order  string // asc, desc
	Limit  uint64
	Offset uint64
}

type StatisticForLink struct {
	Link
	NumberOfRedirects uint64
}

type LinkStore interface {
	Add(ctx context.Context, link *Link) error
	GetNextId(ctx context.Context) (uint64, error)
	GetByURL(ctx context.Context, url string) (*Link, error)
	GetByHash(ctx context.Context, hash string) (*Link, error)
	GetByHashWithStatistics(ctx context.Context, hash string) (*StatisticForLink, error)
	GetByHashes(ctx context.Context, hashes []string) (<-chan Link, <-chan error)
	FindBy(ctx context.Context, params FindByParameters) (<-chan Link, <-chan error)
	CountByQuery(ctx context.Context, query string) (uint64, error)
}
