package linksrv

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lapitskyss/go_backend_1_project/src/frontend/internal/server/handler"
	"github.com/lapitskyss/go_backend_1_project/src/frontend/proto"
)

type LinkSrv struct {
	linkSvcAddr string
	linkSvcConn *grpc.ClientConn
}

func InitLinkServiceClient(addr string) (handler.LinkService, error) {
	con, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "grpc: failed to connect %s", addr)
	}

	return &LinkSrv{
		linkSvcAddr: addr,
		linkSvcConn: con,
	}, nil
}

func (s *LinkSrv) GetRedirectUrl(ctx context.Context, hash string) (string, error) {
	link, err := proto.NewLinkServiceClient(s.linkSvcConn).GetLink(ctx, &proto.GetLinkRequest{
		Hash: hash,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return "", handler.ErrLinkNotFound
		}

		return "", err
	}

	return link.GetUrl(), nil
}
