package grpc

import (
	"context"
	"time"

	apiv1 "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
