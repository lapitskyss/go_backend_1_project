package rpc

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type FrontendServer struct {
	linkSvcAddr string
	linkSvcConn *grpc.ClientConn
}

func NewGRPCClient(ctx context.Context) (*FrontendServer, error) {
	var err error
	var svc = &FrontendServer{}

	err = mustMapEnv(&svc.linkSvcAddr, "LINKSERVICE_SERVICE_ADDR")
	if err != nil {
		return nil, err
	}

	err = mustConnGRPC(ctx, &svc.linkSvcConn, svc.linkSvcAddr)
	if err != nil {
		return nil, err
	}

	return svc, err
}

func mustMapEnv(target *string, envKey string) error {
	v := os.Getenv(envKey)
	if v == "" {
		return errors.New(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v

	return nil
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) error {
	var err error
	*conn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrapf(err, "grpc: failed to connect %s", addr)
	}

	return nil
}
