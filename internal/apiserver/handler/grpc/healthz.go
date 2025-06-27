package grpc

import (
	"context"
	"time"

	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	apiv1 "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	log.With(ctx).Infow("Received health check request", "timestamp", time.Now().Format(time.DateTime))
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
