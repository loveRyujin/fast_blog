package apiserver

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpchandler "github.com/loveRyujin/fast_blog/internal/apiserver/handler/grpc"
	mw "github.com/loveRyujin/fast_blog/internal/pkg/middleware/grpc"
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
	grpcServerOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			// 请求 ID 拦截器
			mw.RequestIDInterceptor(),
		),
	}
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
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return apiv1.RegisterFastBlogHandler(context.Background(), mux, conn)
		},
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
