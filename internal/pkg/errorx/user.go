package errorx

import "net/http"

var (
	// ErrUsernameInvalid 表示用户名无效
	ErrUsernameInvalid = New(http.StatusBadRequest, "InvalidArgument.InvalidUsername", "Invalid username: Username must consist of letters, digits, and underscores only, and its length must be between 3 and 20 characters.")
	// ErrPasswordInvalid 表示密码无效
	ErrPasswordInvalid = New(http.StatusBadRequest, "InvalidArgument.InvalidPassword", "Password is incorrect")
	// ErrUserAlreadyExists 表示用户已存在
	ErrUserAlreadyExists = New(http.StatusBadRequest, "AlreadyExists.UserAlreadyExists", "User already exists")
	// ErrUserNotFound 表示用户未找到
	ErrUserNotFound = New(http.StatusNotFound, "NotFound.UserNotFound", "User not found")
)
