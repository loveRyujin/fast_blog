package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/onexstack_practice/fast_blog/internal/pkg/contextx"
	"github.com/onexstack_practice/fast_blog/internal/pkg/known"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取X-Request_ID， 如果没有就生成一个新的
		requestID := c.Request.Header.Get(known.XRequestID)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将请求ID存放到context中
		ctx := contextx.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		// 将请求ID存放到响应头中
		c.Writer.Header().Set(known.XRequestID, requestID)
		c.Next()
	}
}
