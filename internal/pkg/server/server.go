package server

import "context"

type Server interface {
	Run()
	GracefulStop(ctx context.Context)
}
