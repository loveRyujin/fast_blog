package apiserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	"github.com/loveRyujin/fast_blog/internal/pkg/server"
	genericclioptions "github.com/loveRyujin/fast_blog/pkg/options"
)

const (
	GRPCServerMode        = "grpc"
	HTTPServerMode        = "http"
	GRPCGatewayServerMode = "grpc-gateway"
)

// Config存储应用配置
type Config struct {
	ServerMode   string
	MysqlOptions *genericclioptions.MysqlOptions
	HTTPOptions  *genericclioptions.HTTPOptions
	GRPCOptions  *genericclioptions.GRPCOptions
	JWTKey       string
	Expiration   time.Duration
}

// UnionServer是一个服务器结构体类型
type UnionServer struct {
	srv server.Server
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	log.Infow("Initializing UnionServer", "server-mode", cfg.ServerMode)

	var srv server.Server
	var err error
	switch cfg.ServerMode {
	case HTTPServerMode:
		srv, err = cfg.NewHTTPServer()
	default:
		srv, err = cfg.NewGRPCServerOr()
	}
	if err != nil {
		return nil, err
	}

	return &UnionServer{
		srv: srv,
	}, nil
}

func (s *UnionServer) Run() error {
	go s.srv.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infow("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	// 10s内关闭所有服务，超过10s就超时退出
	s.srv.GracefulStop(ctx)

	log.Infow("Server exited")

	return nil
}
