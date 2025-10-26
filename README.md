# FastBlog

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/loveRyujin/fast_blog)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 📖 项目简介

FastBlog 是一个基于 Go 语言开发的高性能博客 API 服务，采用现代化的微服务架构设计。项目支持多种服务模式（HTTP、gRPC、gRPC-Gateway），使用 Protocol Buffers 定义 API，提供完整的用户管理和文章管理功能。

## ✨ 核心特性

- 🚀 **多协议支持**：支持 HTTP、gRPC、gRPC-Gateway 三种服务模式，灵活切换
- 🔐 **JWT 认证**：完善的身份认证机制，支持 token 刷新
- 📝 **博客管理**：完整的文章 CRUD 操作，支持标题搜索和分页
- 👤 **用户系统**：用户注册、登录、信息更新、密码修改等功能
- 🏗️ **分层架构**：清晰的分层设计（Handler -> Biz -> Store），易于维护和扩展
- 📊 **性能优化**：使用 errgroup 并发处理，提升接口响应速度
- 🔍 **日志系统**：基于 zap 的结构化日志，支持请求追踪
- 🛡️ **安全加密**：密码加密存储，防止明文泄露
- 📈 **性能分析**：集成 pprof，方便性能调优

## 🛠️ 技术栈

### 核心框架
- **Go 1.24**：核心开发语言
- **Gin**：HTTP Web 框架
- **gRPC**：高性能 RPC 框架
- **Protocol Buffers**：API 定义和序列化

### 数据存储
- **MySQL**：关系型数据库
- **GORM**：ORM 框架

### 工具库
- **Viper**：配置管理
- **Cobra**：命令行工具
- **Zap**：结构化日志
- **JWT**：身份认证
- **UUID**：唯一标识生成

## 📦 快速开始

### 前置要求

- Go 1.24+
- MySQL 5.7+
- Protocol Buffers 编译器（如需修改 proto 文件）

### 安装部署

```bash
# 1. 克隆项目
$ mkdir -p $HOME/golang/src/github.com/loveRyujin/
$ cd $HOME/golang/src/github.com/loveRyujin/
$ git clone https://github.com/loveRyujin/fast_blog
$ cd fast_blog/

# 2. 配置数据库
# 编辑 configs/fb-apiserver.yaml，修改数据库连接信息

# 3. 构建项目
$ make build
# 或者使用
$ ./build.sh

# 4. 运行服务
$ _output/fb-apiserver -c configs/fb-apiserver.yaml
```

### 配置说明

编辑 `configs/fb-apiserver.yaml`：

```yaml
# MySQL 数据库配置
mysql:
  addr: 127.0.0.1:3306
  username: your_username
  password: your_password
  database: fastgo
  max-idle-connections: 100
  max-open-connections: 100
  max-connection-life-time: 10s

# 日志配置
log:
  caller-enabled: true
  stacktrace-enabled: true
  level: debug
  format: json
  output: 
    - stdout

# HTTP 服务配置
http:
  addr: 127.0.0.1:8080

# gRPC 服务配置
grpc:
  addr: 127.0.0.1:6666

# 服务模式：http、grpc、grpc-gateway
server-mode: grpc-gateway

# JWT 配置
jwt-key: your_secret_key
expiration: 1000h
```

## 📁 项目结构

```
fast_blog/
├── cmd/                          # 应用程序入口
│   └── fb-apiserver/            # API 服务器主程序
│       ├── app/                 # 应用初始化和配置
│       └── main.go              # 程序入口点
│
├── configs/                      # 配置文件
│   └── fb-apiserver.yaml        # 服务器配置（数据库、日志、JWT等）
│
├── internal/                     # 私有应用和库代码
│   ├── apiserver/               # API 服务器核心实现
│   │   ├── biz/                 # 业务逻辑层（Business Logic）
│   │   │   └── v1/              # v1 版本业务逻辑
│   │   │       ├── user/        # 用户业务逻辑
│   │   │       └── post/        # 文章业务逻辑
│   │   ├── handler/             # 请求处理层（HTTP/gRPC Handler）
│   │   │   ├── user.go          # 用户接口实现
│   │   │   └── post.go          # 文章接口实现
│   │   ├── model/               # 数据模型（Database Models）
│   │   │   ├── user.gen.go      # 用户模型
│   │   │   └── post.gen.go      # 文章模型
│   │   ├── store/               # 数据访问层（Data Access Layer）
│   │   │   ├── store.go         # Store 接口定义
│   │   │   ├── user.go          # 用户数据访问
│   │   │   └── post.go          # 文章数据访问
│   │   ├── pkg/                 # API 服务器内部工具包
│   │   │   ├── conversion/      # 数据转换工具
│   │   │   └── validation/      # 参数验证工具
│   │   ├── http_server.go       # HTTP 服务器实现
│   │   ├── grpc_server.go       # gRPC 服务器实现
│   │   └── server.go            # 服务器统一接口
│   │
│   └── pkg/                     # 内部共享库
│       ├── contextx/            # Context 扩展工具
│       ├── core/                # 核心响应封装
│       ├── errorx/              # 统一错误码定义
│       ├── known/               # 常量定义
│       ├── log/                 # 日志封装
│       ├── middleware/          # 中间件（认证、请求ID等）
│       └── rid/                 # 唯一ID生成器
│
├── pkg/                         # 公共库代码（可被外部引用）
│   ├── api/                     # API 定义
│   │   └── apiserver/v1/        # v1 版本 API
│   │       ├── *.proto          # Protocol Buffers 定义
│   │       ├── *.pb.go          # 生成的 protobuf 代码
│   │       └── *.pb.gw.go       # 生成的 gRPC-Gateway 代码
│   ├── auth/                    # 认证工具（密码加密等）
│   ├── token/                   # JWT Token 工具
│   ├── options/                 # 配置选项
│   └── version/                 # 版本信息
│
├── api/                         # 外部 API 定义
│   └── openapi/                 # OpenAPI/Swagger 文档
│
├── docs/                        # 项目文档
│   ├── agents.md                # 架构设计文档
│   └── images/                  # 文档图片
│
├── third_party/                 # 第三方依赖（Protocol Buffers）
│   └── protobuf/                # Google API proto 文件
│
├── scripts/                     # 脚本文件
│   └── test.sh                  # 测试脚本
│
├── _output/                     # 构建输出目录
│   └── fb-apiserver             # 编译后的二进制文件
│
├── go.mod                       # Go 模块依赖
├── go.sum                       # 依赖版本锁定
├── Makefile                     # 构建脚本
├── build.sh                     # 快速构建脚本
└── README.md                    # 项目说明文档
```

### 核心目录说明

| 目录 | 说明 | 关键文件 |
|------|------|----------|
| `cmd/` | 应用程序入口，每个子目录是一个可执行程序 | `main.go` |
| `internal/apiserver/` | API 服务器核心实现，采用三层架构 | `handler/`, `biz/`, `store/` |
| `internal/pkg/` | 内部共享工具库 | `middleware/`, `errorx/` |
| `pkg/` | 可被外部引用的公共库 | `api/`, `auth/`, `token/` |
| `configs/` | 配置文件存放目录 | `fb-apiserver.yaml` |
| `docs/` | 项目文档 | `agents.md` |

## 🏗️ 项目架构

![架构图](./docs/images/architecture.png)

### 分层设计

项目采用经典的三层架构：

```
┌─────────────┐
│   Handler   │ ← HTTP/gRPC 请求入口，参数验证
└─────────────┘
      ↓
┌─────────────┐
│     Biz     │ ← 业务逻辑层，核心业务处理
└─────────────┘
      ↓
┌─────────────┐
│    Store    │ ← 数据访问层，数据库操作
└─────────────┘
```

## 📚 API 文档

### 用户相关接口

#### 1. 用户注册
```bash
POST /v1/users
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "nickname": "测试用户",
  "email": "test@example.com",
  "phone": "13800138000"
}
```

#### 2. 用户登录
```bash
POST /v1/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

# 响应
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expireAt": "2025-12-31T23:59:59Z"
}
```

#### 3. 刷新 Token
```bash
POST /v1/refresh-token
Authorization: Bearer <your-token>
```

#### 4. 修改密码
```bash
PUT /v1/change-password
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "userID": "user-id",
  "oldPassword": "oldpass123",
  "newPassword": "newpass456"
}
```

#### 5. 获取用户信息
```bash
GET /v1/users/{userID}
Authorization: Bearer <your-token>
```

#### 6. 更新用户信息
```bash
PUT /v1/users/{userID}
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "email": "newemail@example.com"
}
```

#### 7. 删除用户
```bash
DELETE /v1/users/{userID}
Authorization: Bearer <your-token>
```

#### 8. 用户列表
```bash
GET /v1/users?offset=0&limit=10
Authorization: Bearer <your-token>
```

### 文章相关接口

#### 1. 创建文章
```bash
POST /v1/posts
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "title": "我的第一篇博客",
  "content": "这是博客内容..."
}
```

#### 2. 获取文章详情
```bash
GET /v1/posts/{postID}
Authorization: Bearer <your-token>
```

#### 3. 更新文章
```bash
PUT /v1/posts/{postID}
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容..."
}
```

#### 4. 删除文章
```bash
DELETE /v1/posts
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "postIDs": ["post-id-1", "post-id-2"]
}
```

#### 5. 文章列表（支持搜索和分页）
```bash
GET /v1/posts?offset=0&limit=10&title=搜索关键词
Authorization: Bearer <your-token>
```

## 🔧 开发指南

### 编译命令

```bash
# 格式化代码
make format

# 整理依赖
make tidy

# 构建二进制
make build

# 生成 Protocol Buffers 代码
make protoc

# 清理构建产物
make clean

# 完整构建流程（格式化 + 整理依赖 + 编译）
make all
```

### 添加新的 API

1. **定义 Proto 文件**：在 `pkg/api/apiserver/v1/` 目录添加 `.proto` 文件
2. **生成代码**：运行 `make protoc` 生成 Go 代码
3. **实现 Store 层**：在 `internal/apiserver/store/` 实现数据访问
4. **实现 Biz 层**：在 `internal/apiserver/biz/` 实现业务逻辑
5. **实现 Handler 层**：在 `internal/apiserver/handler/` 实现请求处理
6. **注册路由**：在相应的服务器文件中注册路由

### 项目规范

- **代码风格**：遵循 Go 官方代码规范
- **错误处理**：使用统一的错误码系统（见 `internal/pkg/errorx/`）
- **日志记录**：使用结构化日志，携带 context 信息
- **命名规范**：
  - 接口以 `I` 开头（如 `IStore`）
  - 私有实现使用小写开头（如 `userBiz`）
  - 公开类型使用大写开头（如 `UserBiz`）

## 📊 性能特性

- **并发处理**：使用 `errgroup` 进行并发查询，提升列表接口性能
- **连接池**：合理配置数据库连接池，提高数据库访问效率
- **上下文传递**：全链路传递 context，支持请求取消和超时控制
- **优雅关闭**：支持服务优雅关闭，不丢失正在处理的请求

## 🔐 安全特性

- **密码加密**：使用 bcrypt 算法加密存储用户密码
- **JWT 认证**：基于 JWT 的无状态身份认证
- **请求 ID**：每个请求分配唯一 ID，便于追踪和排查问题
- **参数验证**：严格的参数校验，防止非法输入
