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
)
