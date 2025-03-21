package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/onexstack_practice/fast_blog/internal/pkg/core"
	"github.com/onexstack_practice/fast_blog/internal/pkg/errorx"
	v1 "github.com/onexstack_practice/fast_blog/pkg/api/apiserver/v1"
)

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
	slog.Info("Create user function call")

	var request v1.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateCreateUserRequest(c, &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().Create(c, &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// UpdateUser 更新用户信息.
func (h *Handler) UpdateUser(c *gin.Context) {
	slog.Info("Update user function called")

	var request v1.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateUpdateUserRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().Update(c.Request.Context(), &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// DeleteUser 删除用户.
func (h *Handler) DeleteUser(c *gin.Context) {
	slog.Info("Delete user function called")

	var request v1.DeleteUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateDeleteUserRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().Delete(c.Request.Context(), &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// GetUser 获取用户信息.
func (h *Handler) GetUser(c *gin.Context) {
	slog.Info("Get user function called")

	var request v1.GetUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateGetUserRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().Get(c.Request.Context(), &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// ListUser 获取用户列表.
func (h *Handler) ListUser(c *gin.Context) {
	slog.Info("List user function called")

	var request v1.ListUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateListUserRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().List(c.Request.Context(), &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
