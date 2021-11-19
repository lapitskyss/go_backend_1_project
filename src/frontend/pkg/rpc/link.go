package rpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/lapitskyss/go_backend_1_project/src/frontend/proto"
)

var ErrLinkNotFound = errors.New("link not found")

func (fe *FrontendServer) GetLink(ctx context.Context, hash string) (*pb.Link, error) {
	link, err := pb.NewLinkServiceClient(fe.linkSvcConn).GetLink(ctx, &pb.GetLinkRequest{
		Hash: hash,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, ErrLinkNotFound
		}

		return nil, err
	}

	return link, nil
}
