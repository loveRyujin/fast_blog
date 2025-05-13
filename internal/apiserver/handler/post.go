package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/fast_blog/internal/pkg/core"
	"github.com/loveRyujin/fast_blog/internal/pkg/errorx"
	apiv1 "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1"
)

// CreatePost 创建文章
func (h *Handler) CreatePost(c *gin.Context) {
	slog.Info("create post function call")

	var rq apiv1.CreatePostRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateCreatePostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.PostV1().Create(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// UpdatePost 更新文章
func (h *Handler) UpdatePost(c *gin.Context) {
	slog.Info("update post function call")

	var rq apiv1.UpdatePostRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}
	rq.PostID = c.Param("postID")

	if err := h.validator.ValidateUpdatePostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.PostV1().Update(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// DeletePost 删除文章
func (h *Handler) DeletePost(c *gin.Context) {
	slog.Info("delete post function call")

	var rq apiv1.DeletePostRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateDeletePostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.PostV1().Delete(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// GetPost 获取文章
func (h *Handler) GetPost(c *gin.Context) {
	slog.Info("get post function call")

	var rq apiv1.GetPostRequest
	if err := c.ShouldBindUri(&rq); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateGetPostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.PostV1().Get(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}

// ListPost 获取文章列表
func (h *Handler) ListPost(c *gin.Context) {
	slog.Info("list post function call")

	var rq apiv1.ListPostRequest
	if err := c.ShouldBindQuery(&rq); err != nil {
		core.WriteResponse(c, errorx.ErrBind, nil)
		return
	}

	if err := h.validator.ValidateListPostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, errorx.ErrInvalidArugment.WithMessage(err.Error()), nil)
		return
	}

	resp, err := h.biz.PostV1().List(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
