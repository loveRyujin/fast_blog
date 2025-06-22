package apiserver

import (
	"context"

	grpchandler "github.com/loveRyujin/fast_blog/internal/apiserver/handler/grpc"
	"github.com/loveRyujin/fast_blog/internal/pkg/server"
	apiv1 "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	srv  server.Server
	stop func(context.Context)
}

var _ server.Server = (*GRPCServer)(nil)

// NewGRPCServerOr 启动一个GRPC服务或者GRPC-GATEWAY服务
func (cfg *Config) NewGRPCServerOr() (*GRPCServer, error) {
	var grpcServerOptions []grpc.ServerOption
	grpcsrv, err := server.NewGRPCServerOr(
		cfg.GRPCOptions,
		grpcServerOptions,
		func(sr grpc.ServiceRegistrar) {
			apiv1.RegisterFastBlogServer(sr, grpchandler.NewHandler())
		},
	)
	if err != nil {
		return nil, err
	}

	if cfg.ServerMode == GRPCServerMode {
		return &GRPCServer{
			srv: grpcsrv,
			stop: func(ctx context.Context) {
				grpcsrv.GracefulStop(ctx)
			},
		}, nil
	}

	// 先启动GRPC服务，HTTP服务器依赖GRPC服务
	go grpcsrv.Run()

	httpsrv, err := server.NewGRPCGatewayServer(
		cfg.HTTPOptions,
		cfg.GRPCOptions,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &GRPCServer{
		srv: httpsrv,
		stop: func(ctx context.Context) {
			httpsrv.GracefulStop(ctx)
			grpcsrv.GracefulStop(ctx)
		},
	}, nil
}

func (s *GRPCServer) Run() {
	s.srv.Run()
}

func (s *GRPCServer) GracefulStop(ctx context.Context) {
	s.stop(ctx)
}
