package errorx

import "net/http"

var (
	// OK 表示请求成功
	OK = New(http.StatusOK, "", "Success")
	// ErrInternal 表示所有未知的服务器端错误
	ErrInternal = New(http.StatusInternalServerError, "InternalError", "Internal server error")
	// ErrNotFound 表示资源未找到
	ErrNotFound = New(http.StatusNotFound, "NotFound", "Resource not found")
)
