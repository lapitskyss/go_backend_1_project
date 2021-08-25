package server

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/lapitskyss/go_backend_1_project/src/linkservice/genproto"
	linkgrpc "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/controller/grpc/link"
	grpc_recovery "github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/middleware"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/postgres"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/pkg/server_port"
)

type GRPSServer struct {
	server *grpc.Server
	log    *zap.SugaredLogger
}

func NewGRPCServer(ctx context.Context, log *zap.SugaredLogger, rep *postgres.Store) *GRPSServer {
	svc := linkgrpc.New(ctx, log, rep)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(log),
		)),
	)
	pb.RegisterLinkServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)

	return &GRPSServer{srv, log}
}

func (srv *GRPSServer) StartServer() {
	srv.log.Info("Linkservice GRPC server started.")
	go func() {
		port := server_port.GetServerPortFromEnv("LINKSERVICE_GRPC_PORT", 3550)
		l, err := net.Listen("tcp", port)
		if err != nil {
			srv.log.Error("Linkservice GRPC server return error", zap.NamedError("sever_error", err))
		}
		err = srv.server.Serve(l)
		if err != nil {
			srv.log.Error("Linkservice GRPC server return error", zap.NamedError("sever_error", err))
		}
	}()
}

func (srv *GRPSServer) StopServer() {
	srv.server.GracefulStop()
}
