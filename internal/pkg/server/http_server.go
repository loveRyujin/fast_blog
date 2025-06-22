package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	"github.com/loveRyujin/fast_blog/pkg/options"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(httpOptions *options.HTTPOptions, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: handler,
		},
	}
}

func (s *HTTPServer) Run() {
	log.Infow("Start to listen the incoming request on http address", "addr", s.srv.Addr)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
