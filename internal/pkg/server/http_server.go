package server

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	"github.com/loveRyujin/fast_blog/pkg/options"
)

type HTTPServer struct {
	srv *http.Server
}

// NewHTTPServer 创建一个新的HTTP服务器实例
func NewHTTPServer(httpOptions *options.HTTPOptions, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: handler,
		},
	}
}

func (s *HTTPServer) Run() {
	log.Infow("Start to listen the incoming request on http address", "protocol", protocolName(s.srv), "addr", s.srv.Addr)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Errorw(err.Error())
		os.Exit(1)
	}
}

func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
