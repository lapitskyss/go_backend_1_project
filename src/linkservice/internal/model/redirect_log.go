package model

import "time"

type RedirectLog struct {
	ID        uint64
	LinkId    uint64
	UserAgent string
	CreatedAt time.Time
}
