package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/postgres"
)

type link struct {
	ctx context.Context
	log *zap.SugaredLogger
	rep *postgres.Store
}

func New(ctx context.Context, log *zap.SugaredLogger, rep *postgres.Store) *link {
	return &link{ctx, log, rep}
}

func (l *link) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (l *link) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}
