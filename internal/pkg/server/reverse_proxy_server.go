package server

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	"github.com/loveRyujin/fast_blog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type GRPCGatewayServer struct {
	srv *http.Server
}

// NewGRPCGatewayServer 创建一个新的 GRPC 网关服务器实例.
func NewGRPCGatewayServer(
	httpOptions *options.HTTPOptions,
	grpcOptions *options.GRPCOptions,
	registerHandler func(mux *runtime.ServeMux, conn *grpc.ClientConn) error,
) (*GRPCGatewayServer, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 10 * time.Second, // 最小连接超时时间
		}),
	}
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(grpcOptions.Addr, dialOptions...)
	if err != nil {
		log.Errorw("Failed to dial context", "err", err)
		return nil, err
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			// 设置序列化 protobuf 数据时，枚举类型的字段以数字格式输出.
			// 否则，默认会以字符串格式输出，跟枚举类型定义不一致，带来理解成本.
			UseEnumNumbers: true,
		},
	}))
	if err := registerHandler(gwmux, conn); err != nil {
		log.Errorw("Failed to register handler", "err", err)
		return nil, err
	}

	return &GRPCGatewayServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: gwmux,
		},
	}, nil
}

func (s *GRPCGatewayServer) Run() {
	log.Infow("Start to listen the incoming requests", "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalw("Failed to serve HTTP(s) server", "err", err)
	}
}

func (s *GRPCGatewayServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
