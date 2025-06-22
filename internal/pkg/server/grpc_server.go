package server

import (
	"context"
	"net"

	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	"github.com/loveRyujin/fast_blog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	srv *grpc.Server
	lis net.Listener
}

func NewGRPCServerOr(grpcOptions *options.GRPCOptions, serverOptions []grpc.ServerOption, registerServer func(grpc.ServiceRegistrar)) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", grpcOptions.Addr)
	if err != nil {
		return nil, err
	}

	grpcsrv := grpc.NewServer(serverOptions...)

	registerHealthCheckServer(grpcsrv)
	registerServer(grpcsrv)
	reflection.Register(grpcsrv)

	return &GRPCServer{
		srv: grpcsrv,
		lis: lis,
	}, nil
}

func (s *GRPCServer) Run() {
	log.Infow("Start to listening the incoming requests", "protocol", "grpc", "addr", s.lis.Addr().String())
	if err := s.srv.Serve(s.lis); err != nil {
		log.Fatalw("Failed to serve grpc server", "err", err)
	}
}

func (s *GRPCServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop grpc server")
	s.srv.GracefulStop()
}

// registerHealthCheckServer 注册健康检查服务
func registerHealthCheckServer(grpcsrv *grpc.Server) {
	// 创建健康检查服务实例
	healthServer := health.NewServer()

	// 设定服务的健康状态
	healthServer.SetServingStatus("fast_blog", grpc_health_v1.HealthCheckResponse_SERVING)

	// 注册健康检查服务
	grpc_health_v1.RegisterHealthServer(grpcsrv, healthServer)
}
