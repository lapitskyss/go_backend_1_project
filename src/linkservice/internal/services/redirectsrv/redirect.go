package redirectsrv

import (
	"context"
	"time"
)

type Redirect struct {
	ID        uint64
	LinkID    uint64
	CreatedAt time.Time
}

type RedirectLink struct {
	ID  uint64
	URL string
}

type RedirectStore interface {
	GetRedirectLink(ctx context.Context, hash string) (*RedirectLink, error)
	SaveRedirectStatistics(ctx context.Context, redirect Redirect) error
}
