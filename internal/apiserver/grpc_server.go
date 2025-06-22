package apiserver

import (
	"context"

	"github.com/loveRyujin/fast_blog/internal/pkg/server"
)

type GRPCServer struct{}

var _ server.Server = (*GRPCServer)(nil)

func (cfg *Config) NewGRPCServerOr() (*GRPCServer, error) {
	return &GRPCServer{}, nil
}

func (s *GRPCServer) Run() {
	select {}
}

func (s *GRPCServer) GracefulStop(ctx context.Context) {}
