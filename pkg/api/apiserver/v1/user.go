package v1

import (
	"time"
)

// User 表示用户信息
type User struct {
	// userID 表示用户 ID
	UserID string `json:"userID"`
	// username 表示用户名称
	Username string `json:"username"`
	// nickname 表示用户昵称
	Nickname string `json:"nickname"`
	// email 表示用户电子邮箱
	Email string `json:"email"`
	// phone 表示用户手机号
	Phone string `json:"phone"`
	// postCount 表示用户拥有的博客数量
	PostCount int64 `json:"postCount"`
	// createdAt 表示用户注册时间
	CreatedAt time.Time `json:"createdAt"`
	// updatedAt 表示用户最后更新时间
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginRequest 表示登录请求
type LoginRequest struct {
	// username 表示用户名称
	Username string `json:"username"`
	// password 表示用户密码
	Password string `json:"password"`
}

// LoginResponse 表示登录响应
type LoginResponse struct {
	// token 表示登录成功后返回的 token
	Token string `json:"token"`
	// expireAt 表示 token 过期时间
	ExpireAt time.Time `json:"expireAt"`
}

// RefreshTokenRequest 表示刷新 token 请求
type RefreshTokenRequest struct {
}

// RefreshTokenResponse 表示刷新 token 响应
type RefreshTokenResponse struct {
	// token 表示刷新后的 token
	Token string `json:"token"`
	// expireAt 表示 刷新过后的token 过期时间
	ExpireAt time.Time `json:"expireAt"`
}

// CreateUserRequest 表示创建用户请求
type CreateUserRequest struct {
	// username 表示用户名称
	Username string `json:"username"`
	// password 表示用户密码
	Password string `json:"password"`
	// nickname 表示用户昵称
	Nickname *string `json:"nickname"`
	// email 表示用户电子邮箱
	Email string `json:"email"`
	// phone 表示用户手机号
	Phone string `json:"phone"`
}

// CreateUserResponse 表示创建用户响应
type CreateUserResponse struct {
	// userID 表示新创建的用户 ID
	UserID string `json:"userID"`
}

// UpdateUserRequest 表示更新用户请求
type UpdateUserRequest struct {
	// username 表示可选的用户名称
	Username *string `json:"username"`
	// nickname 表示可选的用户昵称
	Nickname *string `json:"nickname"`
	// email 表示可选的用户电子邮箱
	Email *string `json:"email"`
	// phone 表示可选的用户手机号
	Phone *string `json:"phone"`
}

// UpdateUserResponse 表示更新用户响应
type UpdateUserResponse struct {
}

// DeleteUserRequest 表示删除用户请求
type DeleteUserRequest struct {
}

// DeleteUserResponse 表示删除用户响应
type DeleteUserResponse struct {
}

// GetUserRequest 表示获取用户请求
type GetUserRequest struct {
}

// GetUserResponse 表示获取用户响应
type GetUserResponse struct {
	// user 表示返回的用户信息
	User *User `json:"user"`
}

// ListUserRequest 表示用户列表请求
type ListUserRequest struct {
	// offset 表示偏移量
	Offset int64 `json:"offset"`
	// limit 表示每页数量
	Limit int64 `json:"limit"`
}

// ListUserResponse 表示用户列表响应
type ListUserResponse struct {
	// totalCount 表示总用户数
	TotalCount int64 `json:"totalCount"`
	// users 表示用户列表
	Users []*User `json:"users"`
}
