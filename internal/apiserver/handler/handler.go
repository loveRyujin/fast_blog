package handler

import (
	"github.com/onexstack_practice/fast_blog/internal/apiserver/biz"
	"github.com/onexstack_practice/fast_blog/internal/apiserver/pkg/validation"
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
