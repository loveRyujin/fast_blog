# FastBlog 架构设计文档

## 目录

- [1. 架构概览](#1-架构概览)
- [2. 分层设计](#2-分层设计)
- [3. 核心组件](#3-核心组件)
- [4. 数据流转](#4-数据流转)
- [5. 设计模式](#5-设计模式)
- [6. 扩展指南](#6-扩展指南)

## 1. 架构概览

FastBlog 采用经典的三层架构设计，结合了 DDD（领域驱动设计）和 Clean Architecture 的思想，实现了高内聚、低耦合的系统架构。

### 1.1 架构原则

- **单一职责**：每个组件只负责一个明确的功能
- **依赖倒置**：上层依赖抽象接口，不依赖具体实现
- **开闭原则**：对扩展开放，对修改关闭
- **接口隔离**：提供细粒度的接口定义

### 1.2 技术架构图

```
┌──────────────────────────────────────────────────────────┐
│                     客户端层                              │
│   (HTTP Client / gRPC Client / gRPC-Gateway)             │
└──────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────┐
│                   接入层 (Handler)                        │
│  ┌─────────────────────────────────────────────────┐    │
│  │  • 请求参数验证                                  │    │
│  │  • 路由分发                                      │    │
│  │  • 中间件处理（认证、日志、请求ID）             │    │
│  │  • 响应封装                                      │    │
│  └─────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────┐
│                  业务逻辑层 (Biz)                         │
│  ┌─────────────────────────────────────────────────┐    │
│  │  • 业务规则实现                                  │    │
│  │  • 业务流程编排                                  │    │
│  │  • 数据转换与聚合                                │    │
│  │  • 并发控制                                      │    │
│  └─────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────┐
│                 数据访问层 (Store)                        │
│  ┌─────────────────────────────────────────────────┐    │
│  │  • CRUD 操作                                     │    │
│  │  • 查询条件构建                                  │    │
│  │  • 事务管理                                      │    │
│  │  • 数据模型映射                                  │    │
│  └─────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────┐
│                    数据库层 (MySQL)                       │
└──────────────────────────────────────────────────────────┘
```

## 2. 分层设计

### 2.1 Handler 层（接入层）

**位置**：`internal/apiserver/handler/`

**职责**：
- 处理 HTTP/gRPC 请求
- 参数验证和校验
- 调用 Biz 层完成业务处理
- 构建响应数据
- 错误处理和统一响应

**核心文件**：
```
handler/
├── handler.go      # Handler 结构体定义
├── user.go         # 用户相关接口实现
└── post.go         # 文章相关接口实现
```

**示例代码**：
```go
// Handler 持有业务逻辑层和验证器
type Handler struct {
    biz       biz.IBiz
    validator *validation.Validator
}

// 处理请求的典型流程
func (h *Handler) CreateUser(ctx *gin.Context) {
    // 1. 参数绑定
    var req apiv1.CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        core.WriteResponse(ctx, errorx.ErrInvalidParameter, nil)
        return
    }
    
    // 2. 参数验证
    if err := h.validator.Validate(&req); err != nil {
        core.WriteResponse(ctx, err, nil)
        return
    }
    
    // 3. 调用业务逻辑
    resp, err := h.biz.User().Create(ctx, &req)
    if err != nil {
        core.WriteResponse(ctx, err, nil)
        return
    }
    
    // 4. 返回响应
    core.WriteResponse(ctx, nil, resp)
}
```

**设计要点**：
- Handler 不包含业务逻辑，仅做请求转发
- 统一的错误处理和响应格式
- 使用验证器进行参数校验
- 依赖 Biz 层接口，不直接访问数据库

### 2.2 Biz 层（业务逻辑层）

**位置**：`internal/apiserver/biz/`

**职责**：
- 实现核心业务逻辑
- 业务流程编排
- 调用 Store 层进行数据操作
- 数据转换和聚合
- 并发控制和性能优化

**核心文件**：
```
biz/
├── biz.go              # Biz 接口定义
└── v1/
    ├── user/
    │   └── user.go     # 用户业务逻辑
    └── post/
        └── post.go     # 文章业务逻辑
```

**接口定义**：
```go
// IBiz 定义了业务逻辑层的顶层接口
type IBiz interface {
    User() user.UserBiz
    Post() post.PostBiz
}

// UserBiz 定义用户相关的业务方法
type UserBiz interface {
    Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
    Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
    Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
    Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
    List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)
    
    UserExpansion
}

// UserExpansion 定义扩展方法
type UserExpansion interface {
    Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error)
    RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error)
    ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error)
}
```

**示例实现**：
```go
// userBiz 是 UserBiz 接口的实现
type userBiz struct {
    store store.IStore
}

// Create 实现用户创建逻辑
func (b *userBiz) Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
    // 1. 数据转换
    var userM model.User
    _ = copier.Copy(&userM, rq)
    
    // 2. 业务逻辑处理（密码加密）
    encryptedPassword, err := auth.Encrypt(userM.Password)
    if err != nil {
        return nil, err
    }
    userM.Password = encryptedPassword
    
    // 3. 调用 Store 层保存数据
    if err := b.store.User().Create(ctx, &userM); err != nil {
        return nil, err
    }
    
    // 4. 返回结果
    return &apiv1.CreateUserResponse{UserID: userM.UserID}, nil
}

// List 展示了并发优化的使用
func (b *userBiz) List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
    // 1. 查询用户列表
    whr := where.P(int(rq.Offset), int(rq.Limit))
    count, userList, err := b.store.User().List(ctx, whr)
    if err != nil {
        return nil, err
    }
    
    // 2. 使用 errgroup 并发查询每个用户的文章数量
    var m sync.Map
    eg, ctx := errgroup.WithContext(ctx)
    eg.SetLimit(known.MaxErrGroupConcurrency)
    
    for _, user := range userList {
        eg.Go(func() error {
            // 并发查询文章数量
            count, _, err := b.store.Post().List(ctx, where.F("userID", user.UserID))
            if err != nil {
                return err
            }
            
            // 数据转换并存储
            converted := conversion.UserodelToUserV1(user)
            converted.PostCount = count
            m.Store(user.ID, converted)
            
            return nil
        })
    }
    
    // 3. 等待所有并发任务完成
    if err := eg.Wait(); err != nil {
        return nil, err
    }
    
    // 4. 组装结果
    users := make([]*apiv1.User, 0, len(userList))
    for _, item := range userList {
        user, _ := m.Load(item.ID)
        users = append(users, user.(*apiv1.User))
    }
    
    return &apiv1.ListUserResponse{TotalCount: count, Users: users}, nil
}
```

**设计要点**：
- 接口和实现分离，便于测试和扩展
- 使用 Expansion 接口扩展非标准 CRUD 方法
- 合理使用并发提升性能
- 业务逻辑与数据访问解耦

### 2.3 Store 层（数据访问层）

**位置**：`internal/apiserver/store/`

**职责**：
- 封装数据库 CRUD 操作
- 查询条件构建
- 事务管理
- 数据模型定义

**核心文件**：
```
store/
├── store.go        # Store 接口定义
├── user.go         # 用户数据访问实现
└── post.go         # 文章数据访问实现
```

**接口定义**：
```go
// IStore 定义数据访问层的顶层接口
type IStore interface {
    User() UserStore
    Post() PostStore
    TX(ctx context.Context, fn func(context.Context) error) error
    Close() error
}

// UserStore 定义用户数据访问接口
type UserStore interface {
    Create(ctx context.Context, user *model.User) error
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, whr *where.Where) error
    Get(ctx context.Context, whr *where.Where) (*model.User, error)
    List(ctx context.Context, whr *where.Where) (count int64, users []*model.User, err error)
}
```

**示例实现**：
```go
type userStore struct {
    db *gorm.DB
}

// Create 创建用户
func (s *userStore) Create(ctx context.Context, user *model.User) error {
    return s.db.WithContext(ctx).Create(user).Error
}

// Get 根据条件获取用户
func (s *userStore) Get(ctx context.Context, whr *where.Where) (*model.User, error) {
    var user model.User
    if err := s.db.WithContext(ctx).Where(whr.Conditions, whr.Args...).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// List 查询用户列表
func (s *userStore) List(ctx context.Context, whr *where.Where) (int64, []*model.User, error) {
    var count int64
    var users []*model.User
    
    db := s.db.WithContext(ctx).Model(&model.User{})
    
    // 应用查询条件
    if whr != nil && whr.Conditions != "" {
        db = db.Where(whr.Conditions, whr.Args...)
    }
    
    // 获取总数
    if err := db.Count(&count).Error; err != nil {
        return 0, nil, err
    }
    
    // 分页查询
    if whr != nil {
        db = db.Offset(whr.Offset).Limit(whr.Limit)
    }
    
    if err := db.Find(&users).Error; err != nil {
        return 0, nil, err
    }
    
    return count, users, nil
}
```

**设计要点**：
- 使用 Where 对象统一管理查询条件
- 所有操作都传递 context，支持超时控制
- 返回明确的错误信息
- 支持事务操作

### 2.4 Model 层（数据模型）

**位置**：`internal/apiserver/model/`

**职责**：
- 定义数据库表结构
- 实现 GORM 钩子函数
- 数据库模型与 API 模型分离

**核心文件**：
```
model/
├── hook.go         # GORM 钩子实现
├── user.gen.go     # 用户模型
└── post.gen.go     # 文章模型
```

**示例**：
```go
// User 用户数据模型
type User struct {
    ID        uint64         `gorm:"primarykey;column:id"`
    UserID    string         `gorm:"column:userID;type:varchar(128);uniqueIndex;not null"`
    Username  string         `gorm:"column:username;type:varchar(255);uniqueIndex;not null"`
    Password  string         `gorm:"column:password;type:varchar(255);not null"`
    Nickname  string         `gorm:"column:nickname;type:varchar(255)"`
    Email     string         `gorm:"column:email;type:varchar(255);not null"`
    Phone     string         `gorm:"column:phone;type:varchar(32);not null"`
    CreatedAt time.Time      `gorm:"column:createdAt"`
    UpdatedAt time.Time      `gorm:"column:updatedAt"`
}

// TableName 指定表名
func (u *User) TableName() string {
    return "user"
}

// BeforeCreate GORM 创建前钩子，自动生成 UserID
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.UserID == "" {
        u.UserID = "user-" + rid.MustGenString()
    }
    return nil
}
```

## 3. 核心组件

### 3.1 中间件系统

**位置**：`internal/pkg/middleware/`

#### 3.1.1 认证中间件 (auth.go)
```go
// Authn 认证中间件，验证 JWT Token
func Authn() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 提取 Token
        token := c.GetHeader("Authorization")
        token = strings.TrimPrefix(token, "Bearer ")
        
        // 2. 解析和验证 Token
        claims, err := token.Parse(token)
        if err != nil {
            core.WriteResponse(c, errorx.ErrUnauthorized, nil)
            c.Abort()
            return
        }
        
        // 3. 将用户信息存入 Context
        c.Set(known.XUserIDKey, claims.UserID)
        c.Next()
    }
}
```

#### 3.1.2 请求 ID 中间件 (requestid.go)
```go
// RequestID 为每个请求生成唯一 ID
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从 Header 获取或生成新的 RequestID
        requestID := c.GetHeader(known.XRequestIDKey)
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // 设置到 Context 和响应 Header
        c.Set(known.XRequestIDKey, requestID)
        c.Header(known.XRequestIDKey, requestID)
        c.Next()
    }
}
```

### 3.2 错误处理系统

**位置**：`internal/pkg/errorx/`

**设计思路**：
- 定义统一的错误码
- 支持错误链
- 国际化支持
- 错误详情附加

**核心结构**：
```go
// ErrorCode 错误码定义
type ErrorCode struct {
    Code    int    // 错误码
    Message string // 错误消息
    HTTPStatus int // HTTP 状态码
}

// WithMessage 附加错误详情
func (e *ErrorCode) WithMessage(msg string) error {
    return fmt.Errorf("%s: %s", e.Message, msg)
}
```

**常用错误码**：
```go
var (
    ErrSuccess           = &ErrorCode{0, "OK", 200}
    ErrInvalidParameter  = &ErrorCode{10001, "Invalid parameter", 400}
    ErrUnauthorized      = &ErrorCode{10002, "Unauthorized", 401}
    ErrUserNotFound      = &ErrorCode{20001, "User not found", 404}
    ErrPasswordInvalid   = &ErrorCode{20002, "Password invalid", 401}
    ErrPostNotFound      = &ErrorCode{30001, "Post not found", 404}
)
```

### 3.3 验证系统

**位置**：`internal/apiserver/pkg/validation/`

**功能**：
- 参数格式验证
- 业务规则验证
- 自定义验证规则

**示例**：
```go
// Validator 验证器
type Validator struct {
    validate *validator.Validate
}

// Validate 执行验证
func (v *Validator) Validate(obj interface{}) error {
    if err := v.validate.Struct(obj); err != nil {
        return errorx.ErrInvalidParameter.WithMessage(err.Error())
    }
    return nil
}

// 自定义验证规则
func ValidateUsername(username string) error {
    if len(username) < 3 || len(username) > 20 {
        return errors.New("username length must be between 3 and 20")
    }
    if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
        return errors.New("username can only contain letters, numbers, and underscores")
    }
    return nil
}
```

### 3.4 数据转换系统

**位置**：`internal/apiserver/pkg/conversion/`

**职责**：
- Model 到 API 对象的转换
- API 对象到 Model 的转换

**示例**：
```go
// UserodelToUserV1 将数据库模型转换为 API 模型
func UserodelToUserV1(user *model.User) *apiv1.User {
    return &apiv1.User{
        UserID:    user.UserID,
        Username:  user.Username,
        Nickname:  user.Nickname,
        Email:     user.Email,
        Phone:     user.Phone,
        CreatedAt: timestamppb.New(user.CreatedAt),
        UpdatedAt: timestamppb.New(user.UpdatedAt),
    }
}
```

### 3.5 日志系统

**位置**：`internal/pkg/log/`

**特性**：
- 基于 Zap 的高性能日志
- 结构化日志输出
- 支持日志级别配置
- 自动添加 RequestID 和 UserID

**使用示例**：
```go
// 带上下文的日志
log.With(ctx).Infow("User created", "userID", userID)
log.With(ctx).Errorw("Failed to create user", "err", err)

// 调试日志
log.Debugw("Processing request", "params", params)
```

## 4. 数据流转

### 4.1 请求处理流程

```
1. 客户端发起请求
   ↓
2. 中间件处理链
   - RequestID 生成
   - 请求日志记录
   - JWT 认证（需要时）
   ↓
3. 路由匹配到 Handler
   ↓
4. Handler 层
   - 参数绑定
   - 参数验证
   - 调用 Biz 层
   ↓
5. Biz 层
   - 业务逻辑处理
   - 调用 Store 层
   - 数据转换
   ↓
6. Store 层
   - 构建查询条件
   - 执行数据库操作
   - 返回结果
   ↓
7. 响应返回
   - 统一响应格式
   - 错误码转换
   - 日志记录
```

### 4.2 用户登录流程示例

```
[客户端] → POST /v1/login
    ↓
[RequestID 中间件] → 生成请求 ID
    ↓
[Handler] → 接收登录请求
    ↓ 验证参数
    ↓
[Biz.Login]
    ↓ 1. 根据用户名查询用户
    ↓
[Store.User.Get] → 查询数据库
    ↓
[Biz.Login]
    ↓ 2. 验证密码
    ↓ 3. 生成 JWT Token
    ↓
[Handler] → 返回 Token
    ↓
[客户端] ← 收到 Token
```

### 4.3 创建文章流程示例

```
[客户端] → POST /v1/posts + JWT Token
    ↓
[RequestID 中间件]
    ↓
[Authn 中间件] → 验证 Token，提取 UserID
    ↓
[Handler] → 接收创建文章请求
    ↓ 验证参数
    ↓
[Biz.CreatePost]
    ↓ 1. 从 Context 获取 UserID
    ↓ 2. 构建 Post 对象
    ↓
[Store.Post.Create] → 保存到数据库
    ↓ GORM BeforeCreate Hook → 生成 PostID
    ↓
[Biz.CreatePost] ← 返回 PostID
    ↓
[Handler] → 返回创建结果
    ↓
[客户端] ← 收到 PostID
```

## 5. 设计模式

### 5.1 工厂模式

用于创建 Biz 和 Store 实例：

```go
// Biz 工厂
func NewBiz(store store.IStore) IBiz {
    return &biz{
        userBiz: user.New(store),
        postBiz: post.New(store),
    }
}

// Store 工厂
func NewStore(db *gorm.DB) IStore {
    return &datastore{
        db: db,
    }
}
```

### 5.2 依赖注入

通过接口注入依赖，实现解耦：

```go
// Handler 依赖 Biz 接口
type Handler struct {
    biz biz.IBiz
}

// Biz 依赖 Store 接口
type userBiz struct {
    store store.IStore
}
```

### 5.3 适配器模式

支持多种服务模式（HTTP、gRPC、gRPC-Gateway）：

```go
switch cfg.ServerMode {
case HTTPServerMode:
    srv, err = cfg.NewHTTPServer()
case GRPCServerMode:
    srv, err = cfg.NewGRPCServer()
case GRPCGatewayServerMode:
    srv, err = cfg.NewGRPCGatewayServer()
}
```

### 5.4 装饰器模式

通过中间件链装饰 HTTP 处理器：

```go
router.Use(
    middleware.RequestID(),
    middleware.Header(),
    gin.Logger(),
    gin.Recovery(),
)

// 需要认证的路由
authRouter := router.Group("/v1")
authRouter.Use(middleware.Authn())
```

### 5.5 策略模式

Where 对象用于构建不同的查询策略：

```go
// 简单查询
whr := where.F("userID", userID)

// 分页查询
whr := where.P(offset, limit)

// 复杂查询
whr := where.F("userID", userID).P(0, 10).Q("title like ?", "%keyword%")
```

## 6. 扩展指南

### 6.1 添加新的业务模块

假设要添加"评论"功能：

#### 步骤 1：定义 Proto 文件

创建 `pkg/api/apiserver/v1/comment.proto`：

```protobuf
syntax = "proto3";

package v1;

option go_package = "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1";

message Comment {
    string commentID = 1;
    string postID = 2;
    string userID = 3;
    string content = 4;
    google.protobuf.Timestamp createdAt = 5;
}

message CreateCommentRequest {
    string postID = 1;
    string content = 2;
}

message CreateCommentResponse {
    string commentID = 1;
}

// ... 其他消息定义
```

#### 步骤 2：生成代码

```bash
make protoc
```

#### 步骤 3：定义数据模型

创建 `internal/apiserver/model/comment.gen.go`：

```go
type Comment struct {
    ID        uint64    `gorm:"primarykey"`
    CommentID string    `gorm:"column:commentID;uniqueIndex"`
    PostID    string    `gorm:"column:postID;index"`
    UserID    string    `gorm:"column:userID;index"`
    Content   string    `gorm:"column:content;type:text"`
    CreatedAt time.Time `gorm:"column:createdAt"`
    UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
    if c.CommentID == "" {
        c.CommentID = "comment-" + rid.MustGenString()
    }
    return nil
}
```

#### 步骤 4：实现 Store 层

创建 `internal/apiserver/store/comment.go`：

```go
type CommentStore interface {
    Create(ctx context.Context, comment *model.Comment) error
    Get(ctx context.Context, whr *where.Where) (*model.Comment, error)
    List(ctx context.Context, whr *where.Where) (int64, []*model.Comment, error)
    Delete(ctx context.Context, whr *where.Where) error
}

type commentStore struct {
    db *gorm.DB
}

func newCommentStore(db *gorm.DB) *commentStore {
    return &commentStore{db: db}
}

// 实现接口方法...
```

在 `store.go` 中添加：

```go
type IStore interface {
    User() UserStore
    Post() PostStore
    Comment() CommentStore  // 新增
    // ...
}

func (ds *datastore) Comment() CommentStore {
    return newCommentStore(ds.db)
}
```

#### 步骤 5：实现 Biz 层

创建 `internal/apiserver/biz/v1/comment/comment.go`：

```go
type CommentBiz interface {
    Create(ctx context.Context, rq *apiv1.CreateCommentRequest) (*apiv1.CreateCommentResponse, error)
    // ... 其他方法
}

type commentBiz struct {
    store store.IStore
}

func New(store store.IStore) *commentBiz {
    return &commentBiz{store: store}
}

func (b *commentBiz) Create(ctx context.Context, rq *apiv1.CreateCommentRequest) (*apiv1.CreateCommentResponse, error) {
    // 实现创建评论逻辑
    var commentM model.Comment
    copier.Copy(&commentM, rq)
    commentM.UserID = contextx.UserID(ctx)
    
    if err := b.store.Comment().Create(ctx, &commentM); err != nil {
        return nil, err
    }
    
    return &apiv1.CreateCommentResponse{CommentID: commentM.CommentID}, nil
}
```

在 `biz.go` 中添加：

```go
type IBiz interface {
    User() user.UserBiz
    Post() post.PostBiz
    Comment() comment.CommentBiz  // 新增
}

func (b *biz) Comment() comment.CommentBiz {
    return b.commentBiz
}

func NewBiz(store store.IStore) IBiz {
    return &biz{
        userBiz:    user.New(store),
        postBiz:    post.New(store),
        commentBiz: comment.New(store),  // 新增
    }
}
```

#### 步骤 6：实现 Handler 层

在 `handler/comment.go` 中：

```go
func (h *Handler) CreateComment(ctx *gin.Context) {
    var req apiv1.CreateCommentRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        core.WriteResponse(ctx, errorx.ErrInvalidParameter, nil)
        return
    }
    
    resp, err := h.biz.Comment().Create(ctx, &req)
    if err != nil {
        core.WriteResponse(ctx, err, nil)
        return
    }
    
    core.WriteResponse(ctx, nil, resp)
}
```

#### 步骤 7：注册路由

在相应的服务器文件中注册路由：

```go
authRouter.POST("/comments", handler.CreateComment)
authRouter.GET("/comments/:commentID", handler.GetComment)
```

### 6.2 添加中间件

创建 `internal/pkg/middleware/ratelimit.go`：

```go
func RateLimit(limit int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(limit), limit)
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            core.WriteResponse(c, errorx.ErrTooManyRequests, nil)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

注册中间件：

```go
router.Use(middleware.RateLimit(100))
```

### 6.3 添加自定义错误码

在 `internal/pkg/errorx/comment.go` 中：

```go
var (
    ErrCommentNotFound = &ErrorCode{40001, "Comment not found", 404}
    ErrCommentTooLong  = &ErrorCode{40002, "Comment content too long", 400}
)
```

## 7. 最佳实践

### 7.1 错误处理

**DO**：
```go
// 使用统一的错误码
if user == nil {
    return nil, errorx.ErrUserNotFound
}

// 附加错误详情
if err != nil {
    return nil, errorx.ErrInternalError.WithMessage(err.Error())
}
```

**DON'T**：
```go
// 不要直接返回底层错误
return nil, err

// 不要使用模糊的错误信息
return nil, errors.New("error")
```

### 7.2 日志记录

**DO**：
```go
// 使用结构化日志
log.With(ctx).Infow("User created", "userID", userID, "username", username)

// 记录关键操作
log.With(ctx).Infow("User login successful", "userID", userID)
```

**DON'T**：
```go
// 不要使用 fmt.Println
fmt.Println("User created:", userID)

// 不要在日志中输出敏感信息
log.Infow("User login", "password", password)  // 错误！
```

### 7.3 性能优化

**DO**：
```go
// 使用并发提升性能
eg, ctx := errgroup.WithContext(ctx)
eg.SetLimit(known.MaxErrGroupConcurrency)
for _, item := range items {
    eg.Go(func() error {
        // 处理逻辑
        return nil
    })
}
eg.Wait()

// 使用连接池
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(100)

// 添加索引
`gorm:"index"`
```

**DON'T**：
```go
// 不要在循环中执行数据库查询（N+1 问题）
for _, userID := range userIDs {
    user, _ := store.User().Get(ctx, where.F("userID", userID))
    // ...
}
```

### 7.4 安全实践

**DO**：
```go
// 密码加密
encryptedPassword, _ := auth.Encrypt(password)

// 使用参数化查询
db.Where("username = ?", username)

// 验证用户输入
if err := validator.Validate(req); err != nil {
    return err
}
```

**DON'T**：
```go
// 不要明文存储密码
user.Password = password

// 不要拼接 SQL
query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)

// 不要跳过验证
// 总是验证用户输入！
```

### 7.5 Context 使用

**DO**：
```go
// 传递 context
func (b *userBiz) Create(ctx context.Context, req *Request) error {
    return b.store.User().Create(ctx, user)
}

// 从 context 获取用户信息
userID := contextx.UserID(ctx)
```

**DON'T**：
```go
// 不要忽略 context
func (b *userBiz) Create(req *Request) error {
    return b.store.User().Create(context.Background(), user)
}

// 不要创建新的空 context
ctx := context.Background()  // 应该使用传入的 ctx
```

## 8. 总结

FastBlog 项目通过清晰的分层架构、合理的设计模式和最佳实践，实现了一个可扩展、可维护、高性能的博客 API 服务。主要特点：

1. **清晰的职责分离**：Handler、Biz、Store 各司其职
2. **面向接口编程**：依赖抽象而非具体实现
3. **统一的错误处理**：标准化的错误码和响应格式
4. **完善的中间件**：认证、日志、请求追踪等
5. **性能优化**：并发处理、连接池、索引优化
6. **安全保障**：密码加密、JWT 认证、参数验证
7. **易于扩展**：遵循开闭原则，便于添加新功能

希望这份文档能帮助你理解项目架构，并指导后续的开发工作！

