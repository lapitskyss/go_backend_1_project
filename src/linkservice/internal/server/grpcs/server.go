package grpcs

import (
	"net"

	grpcm "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/grpcs/handler"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/server/grpcs/middleware"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/proto"
)

type GRPCServer struct {
	server *grpc.Server
	log    *zap.Logger
	errors chan error
}

func NewGRPCServer(log *zap.Logger, h *handler.RedirectHandler) *GRPCServer {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpcm.ChainUnaryServer(
			middleware.RecoverUnaryServerInterceptor(log),
		)),
	)

	proto.RegisterLinkServiceServer(srv, h)
	healthpb.RegisterHealthServer(srv, h)

	return &GRPCServer{
		server: srv,
		log:    log,
	}
}

func (s *GRPCServer) Start(port string) {
	s.log.Info("Linkservice GRPC server started on port " + port + ".")
	go func() {
		l, err := net.Listen("tcp", ":"+port)
		if err != nil {
			s.log.Error("Linkservice GRPC server return error", zap.Error(err))
			s.errors <- err
			return
		}
		err = s.server.Serve(l)
		if err != nil {
			s.log.Error("Linkservice GRPC server return error", zap.Error(err))
			s.errors <- err
		}
	}()
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) Notify() <-chan error {
	return s.errors
}
