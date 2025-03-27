package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/onexstack_practice/fast_blog/internal/pkg/core"
	"github.com/onexstack_practice/fast_blog/internal/pkg/errorx"
	v1 "github.com/onexstack_practice/fast_blog/pkg/api/apiserver/v1"
)

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	slog.Info("Login function called")

	var request v1.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateLoginRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().Login(c.Request.Context(), &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// RefreshToken 刷新token
func (h *Handler) RefreshToken(c *gin.Context) {
	slog.Info("Refresh token function called")

	var request v1.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateRefreshTokenRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().RefreshToken(c.Request.Context(), &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	slog.Info("Change user password function call")

	var request v1.ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateChangePasswordRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().ChangePassword(c.Request.Context(), &request)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
	slog.Info("Create user function call")

	var request v1.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateCreateUserRequest(c.Request.Context(), &request); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.UserV1().Create(c.Request.Context(), &request)
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
	if err := c.ShouldBindUri(&request); err != nil {
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
	if err := c.ShouldBindUri(&request); err != nil {
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
	if err := c.ShouldBindQuery(&request); err != nil {
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
