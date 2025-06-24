package server

import (
	"context"
	"net/http"
)

type Server interface {
	Run()
	GracefulStop(ctx context.Context)
}

// protocolName 根据http.Server返回协议名称
func protocolName(server *http.Server) string {
	if server.TLSConfig != nil {
		return "https"
	}

	return "http"
}
