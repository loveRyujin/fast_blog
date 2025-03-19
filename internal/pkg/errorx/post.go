package errorx

import "net/http"

// ErrPostNotFound 表示文章未找到
var ErrPostNotFound = New(http.StatusNotFound, "NotFound.PostNotFound", "Post not found")
