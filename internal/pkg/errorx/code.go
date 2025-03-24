package errorx

import "net/http"

var (
	// OK 表示请求成功
	OK = New(http.StatusOK, "", "Success")
	// ErrInternal 表示所有未知的服务器端错误
	ErrInternal = New(http.StatusInternalServerError, "InternalError", "Internal server error")
	// ErrNotFound 表示资源未找到
	ErrNotFound = New(http.StatusNotFound, "NotFound", "Resource not found")
	// ErrDBWrite 表示数据库写入错误
	ErrDBWrite = New(http.StatusInternalServerError, "InternalError.DBWrite", "Database write error")
	// ErrDBRead 表示数据库读取错误
	ErrDBRead = New(http.StatusInternalServerError, "InternalError.DBRead", "Database read error")
	// ErrBind 表示请求参数绑定错误
	ErrBind = New(http.StatusBadRequest, "BindError", "Request parameter binding error")
	// ErrInvalidArugment 表示参数验证失败
	ErrInvalidArugment = New(http.StatusBadRequest, "InvalidArgument", "Argument valification failed")
	// ErrSignToken 表示签名令牌失败
	ErrSignToken = New(http.StatusUnauthorized, "Unauthenticated.SignToken", "Failed to sign token")
	// ErrTokenInvalid 表示令牌无效
	ErrTokenInvalid = New(http.StatusUnauthorized, "Unauthenticated.TokenInvalid", "Invalid token")
)
