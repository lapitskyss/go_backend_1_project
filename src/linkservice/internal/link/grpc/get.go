package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/lapitskyss/go_backend_1_project/src/linkservice/genproto"
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

	return &pb.Link{
		Url:  link.URL,
		Hash: link.Hash,
	}, nil
}
