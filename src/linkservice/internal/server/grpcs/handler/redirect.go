package handler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/redirectsrv"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/proto"
)

type RedirectHandler struct {
	proto.UnimplementedLinkServiceServer
	log *zap.Logger
	rs  *redirectsrv.RedirectService
}

func NewRedirectHandler(log *zap.Logger, redirectService *redirectsrv.RedirectService) *RedirectHandler {
	return &RedirectHandler{
		log: log,
		rs:  redirectService,
	}
}

func (r RedirectHandler) GetLink(ctx context.Context, request *proto.GetLinkRequest) (*proto.Link, error) {
	url, err := r.rs.URL(ctx, request.GetHash())

	if err != nil {
		if errors.Is(err, redirectsrv.ErrLinkNotFound) {
			return nil, status.Error(codes.NotFound, "link not found")
		} else {
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &proto.Link{
		Url: url,
	}, nil
}

func (r RedirectHandler) Check(_ context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (r RedirectHandler) Watch(_ *healthpb.HealthCheckRequest, _ healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}
