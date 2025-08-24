package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/fast_blog/internal/pkg/core"
	"github.com/loveRyujin/fast_blog/internal/pkg/errorx"
	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	apiv1 "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1"
)

// CreatePost 创建文章
func (h *Handler) CreatePost(c *gin.Context) {
	log.Infow("create post function call")

	var rq apiv1.CreatePostRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrBind)
		return
	}

	if err := h.validator.ValidateCreatePostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrInvalidArugment.WithMessage(err.Error()))
		return
	}

	resp, err := h.biz.PostV1().Create(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// UpdatePost 更新文章
func (h *Handler) UpdatePost(c *gin.Context) {
	log.Infow("update post function call")

	var rq apiv1.UpdatePostRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrBind)
		return
	}
	rq.PostID = c.Param("postID")

	if err := h.validator.ValidateUpdatePostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrInvalidArugment.WithMessage(err.Error()))
		return
	}

	resp, err := h.biz.PostV1().Update(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// DeletePost 删除文章
func (h *Handler) DeletePost(c *gin.Context) {
	log.Infow("delete post function call")

	var rq apiv1.DeletePostRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrBind)
		return
	}

	if err := h.validator.ValidateDeletePostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrInvalidArugment.WithMessage(err.Error()))
		return
	}

	resp, err := h.biz.PostV1().Delete(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// GetPost 获取文章
func (h *Handler) GetPost(c *gin.Context) {
	log.Infow("get post function call")

	var rq apiv1.GetPostRequest
	if err := c.ShouldBindUri(&rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrBind)
		return
	}

	if err := h.validator.ValidateGetPostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrInvalidArugment.WithMessage(err.Error()))
		return
	}

	resp, err := h.biz.PostV1().Get(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}

// ListPost 获取文章列表
func (h *Handler) ListPost(c *gin.Context) {
	log.Infow("list post function call")

	var rq apiv1.ListPostRequest
	if err := c.ShouldBindQuery(&rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrBind)
		return
	}

	if err := h.validator.ValidateListPostRequest(c.Request.Context(), &rq); err != nil {
		core.WriteResponse(c, nil, errorx.ErrInvalidArugment.WithMessage(err.Error()))
		return
	}

	resp, err := h.biz.PostV1().List(c.Request.Context(), &rq)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, resp, nil)
}
