package model

import "time"

type Link struct {
	ID        uint64    `json:"-"`
	URL       string    `json:"url"`
	Hash      string    `json:"hash,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
