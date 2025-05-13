package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/fast_blog/internal/pkg/contextx"
	"github.com/loveRyujin/fast_blog/internal/pkg/core"
	"github.com/loveRyujin/fast_blog/internal/pkg/errorx"
	"github.com/loveRyujin/fast_blog/pkg/token"
)

// Authn 是认证中间件，用来从 gin.Context 中提取 token 并验证 token 是否合法，
// 如果合法则将 token 中的 sub 作为<用户名>存放在 gin.Context 的 XUsernameKey 键中.
func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 JWT Token
		userID, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errorx.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		// 将用户ID和用户名注入到上下文中
		ctx := contextx.WithUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)

		// 继续后续的操作
		c.Next()

	}
}
