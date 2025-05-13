package handler

import (
	"github.com/loveRyujin/fast_blog/internal/apiserver/biz"
	"github.com/loveRyujin/fast_blog/internal/apiserver/pkg/validation"
)

type Handler struct {
	biz       biz.IBiz
	validator *validation.Validator
}

func NewHandler(biz biz.IBiz, validator *validation.Validator) *Handler {
	return &Handler{
		biz:       biz,
		validator: validator,
	}
}
