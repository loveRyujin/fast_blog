package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/fast_blog/internal/pkg/core2"
	"github.com/loveRyujin/fast_blog/internal/pkg/log"
)

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	log.Infow("Login function called")

	core2.HandleJSONRequest(c, h.biz.UserV1().Login, h.validator.ValidateLoginRequest)
}

// RefreshToken 刷新token
func (h *Handler) RefreshToken(c *gin.Context) {
	log.Infow("Refresh token function called")

	core2.HandleJSONRequest(c, h.biz.UserV1().RefreshToken, h.validator.ValidateRefreshTokenRequest)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	log.Infow("Change user password function call")

	core2.HandleJSONRequest(c, h.biz.UserV1().ChangePassword, h.validator.ValidateChangePasswordRequest)
}

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
	log.Infow("Create user function call")

	core2.HandleJSONRequest(c, h.biz.UserV1().Create, h.validator.ValidateCreateUserRequest)
}

// UpdateUser 更新用户信息.
func (h *Handler) UpdateUser(c *gin.Context) {
	log.Infow("Update user function called")

	core2.HandleJSONRequest(c, h.biz.UserV1().Update, h.validator.ValidateUpdateUserRequest)
}

// DeleteUser 删除用户.
func (h *Handler) DeleteUser(c *gin.Context) {
	log.Infow("Delete user function called")

	core2.HandleURIRequest(c, h.biz.UserV1().Delete, h.validator.ValidateDeleteUserRequest)
}

// GetUser 获取用户信息.
func (h *Handler) GetUser(c *gin.Context) {
	log.Infow("Get user function called")

	core2.HandleURIRequest(c, h.biz.UserV1().Get, h.validator.ValidateGetUserRequest)
}

// ListUser 获取用户列表.
func (h *Handler) ListUser(c *gin.Context) {
	log.Infow("List user function called")

	core2.HandleQueryRequest(c, h.biz.UserV1().List, h.validator.ValidateListUserRequest)
}
