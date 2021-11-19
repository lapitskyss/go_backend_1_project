package linksrv

import (
	"context"
	"net/url"
	"time"

	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/pkg/e"
)

type CreateLinkRequest struct {
	URL string
}

// Create new link
func (s *LinkService) Create(ctx context.Context, r CreateLinkRequest) (*Link, error) {
	// validate request data
	err := r.validate()
	if err != nil {
		return nil, e.ErrBadRequest(err.Error())
	}

	// check is link already exist in database
	existingLink, err := s.ls.GetByURL(ctx, r.URL)
	if err != nil {
		s.log.Error("Link service error to call GetByURL", zap.Error(err))
		return nil, e.ErrInternal()
	}

	// return existing link
	if existingLink != nil {
		return existingLink, nil
	}

	// get next id for link
	linkId, err := s.ls.GetNextId(ctx)
	if err != nil {
		s.log.Error("Link service error to call GetNextId", zap.Error(err))
		return nil, e.ErrInternal()
	}

	// generate hash for link
	hash, err := getHashForLink(linkId)
	if err != nil {
		s.log.Error("Link service error to call getHashForLink", zap.Error(err))
		return nil, e.ErrInternal()
	}

	newLink := &Link{
		ID:        linkId,
		URL:       r.URL,
		Hash:      hash,
		CreatedAt: time.Now(),
	}

	// save link to store
	err = s.ls.Add(ctx, newLink)
	if err != nil {
		s.log.Error("Link service error to call Add", zap.Error(err))
		return nil, e.ErrInternal()
	}

	return newLink, nil
}

func (r *CreateLinkRequest) validate() error {
	result, err := url.ParseRequestURI(r.URL)

	if err != nil || result.Scheme == "" {
		return e.ErrBadRequest("incorrect link URL")
	}

	if len(r.URL) > 2048 {
		return e.ErrBadRequest("link URL is to long")
	}

	return nil
}

func getHashForLink(id uint64) (string, error) {
	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 5
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	return h.EncodeInt64([]int64{int64(id)})
}
