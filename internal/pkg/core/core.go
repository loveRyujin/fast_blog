package core

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/fast_blog/internal/pkg/errorx"
	"github.com/onexstack/onexstack/pkg/errorsx"
)

// Validator 是验证函数的类型，用于对绑定的数据结构进行验证.
type Validator[T any] func(context.Context, *T) error

// Binder 定义绑定函数的类型，用于绑定请求数据到相应结构体.
type Binder func(any) error

// Handler 是处理函数的类型，用于处理已经绑定和验证的数据.
type Handler[T any, R any] func(ctx context.Context, req *T) (R, error)

// ErrorResponse 定义了错误响应的结构，
// 用于 API 请求中发生错误时返回统一的格式化错误信息.
type ErrorResponse struct {
	// 错误原因，标识错误类型
	Reason string `json:"reason,omitempty"`
	// 错误详情的描述信息
	Message string `json:"message,omitempty"`
	// 附带的元数据信息
	Metadata map[string]string `json:"metadata,omitempty"`
}

func HandleJSONRequest[T any, R any](c *gin.Context, handler Handler[T, R], validator ...Validator[T]) {
	HandleRequest(c, c.ShouldBindJSON, handler, validator...)
}

func HandleQueryRequest[T any, R any](c *gin.Context, handler Handler[T, R], validator ...Validator[T]) {
	HandleRequest(c, c.ShouldBindQuery, handler, validator...)
}

func HandleURIRequest[T any, R any](c *gin.Context, handler Handler[T, R], validator ...Validator[T]) {
	HandleRequest(c, c.ShouldBindUri, handler, validator...)
}

func HandleRequest[T any, R any](c *gin.Context, binder Binder, handler Handler[T, R], validator ...Validator[T]) {
	var req T

	if err := ReadRequest(c, binder, &req, validator...); err != nil {
		WriteResponse(c, nil, err)
		return
	}

	resp, err := handler(c.Request.Context(), &req)
	WriteResponse(c, resp, err)
}

func ReadRequest[T any](c *gin.Context, binder Binder, request *T, validator ...Validator[T]) error {
	if err := binder(request); err != nil {
		return errorx.ErrBind.WithMessage(err.Error())
	}

	for _, validate := range validator {
		if validate == nil {
			continue
		}
		if err := validate(c.Request.Context(), request); err != nil {
			return err
		}
	}

	return nil
}

// WriteResponse 是通用的响应函数.
// 它会根据是否发生错误，生成成功响应或标准化的错误响应.
func WriteResponse(c *gin.Context, data any, err error) {
	if err != nil {
		// 如果发生错误，生成错误响应
		errx := errorsx.FromError(err) // 提取错误详细信息
		c.JSON(errx.Code, ErrorResponse{
			Reason:   errx.Reason,
			Message:  errx.Message,
			Metadata: errx.Metadata,
		})
		return
	}

	// 如果没有错误，返回成功响应
	c.JSON(http.StatusOK, data)
}
