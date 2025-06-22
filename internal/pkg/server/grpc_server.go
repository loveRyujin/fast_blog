package server

import "context"

type GRPCServer struct{}

func (s *GRPCServer) Run() {
	select {}
}

func (s *GRPCServer) GracefulStop(ctx context.Context) {

}
