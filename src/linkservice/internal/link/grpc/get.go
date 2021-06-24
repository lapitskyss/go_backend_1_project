package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/lapitskyss/go_backend_1_project/src/linkservice/genproto"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
)

func (l *link) GetLink(ctx context.Context, req *pb.GetLinkRequest) (*pb.Link, error) {
	link, err := l.rep.Link.GetByHash(req.Hash)
	if err != nil {
		l.log.Error(err)
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if link == nil {
		return nil, status.New(codes.NotFound, "link not found").Err()
	}

	err = l.rep.RedirectLog.Add(&model.RedirectLog{
		LinkId:    link.ID,
		UserAgent: req.UserAgent,
		CreatedAt: time.Now(),
	})
	if err != nil {
		l.log.Error(err)
	}

	return &pb.Link{
		Url:  link.URL,
		Hash: link.Hash,
	}, nil
}
