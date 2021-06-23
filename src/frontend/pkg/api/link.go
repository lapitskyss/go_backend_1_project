package api

import (
	"context"
	"errors"
)

type LinkService service

type Link struct {
	URL       *string    `json:"url,omitempty"`
	Hash      *string    `json:"hash,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

var ErrLinkNotFound = errors.New("link not found")

func (s *LinkService) GetLinkByHash(ctx context.Context, hash string) (*Link, error) {
	u := "link/" + hash
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	link := &Link{}
	resp, err := s.client.Do(ctx, req, link)
	if resp != nil {
		if c := resp.StatusCode; 404 == c {
			return nil, ErrLinkNotFound
		}
	}

	if err != nil {
		return nil, err
	}

	return link, nil
}
