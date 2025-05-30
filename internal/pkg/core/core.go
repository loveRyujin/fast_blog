package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/fast_blog/internal/pkg/errorx"
)

// ErrorResponse 是通用的错误响应结构体
// 用于API请求发生错误时返回统一的格式化错误信息
type ErrorResponse struct {
	// 错误原因，表示错误类型
	Reason string `json:"reason,omitempty"`
	// 错误详情的描述信息
	Message string `json:"message,omitempty"`
}

// WriteResponse 是通用的响应函数
// 根据是否有错误进行不同的响应表现
func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		errx := errorx.FromError(err)
		c.JSON(errx.Code, ErrorResponse{
			Reason:  errx.Reason,
			Message: errx.Message,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
