package grpc

import (
	apiv1 "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1"
)

type Handler struct {
	apiv1.UnimplementedFastBlogServer
}

func NewHandler() *Handler {
	return &Handler{}
}
