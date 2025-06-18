# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
PROJ_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output
# 构建二进制文件名
EXECUTABLE := fb-apiserver
# ProtoBuf 文件存放路径
APIROOT_DIR := $(PROJ_ROOT_DIR)/pkg/api

# ==============================================================================
# 定义版本相关变量

## 指定应用使用的 version 包，会通过 `-ldflags -X` 向该包中指定的变量注入值
VERSION_PACKAGE=github.com/loveRyujin/fast_blog/pkg/version
## 定义 VERSION 语义化版本号
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

## 检查代码仓库是否是 dirty（默认dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
    GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
    -X $(VERSION_PACKAGE).gitVersion=$(VERSION) \
    -X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) \
    -X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) \
    -X $(VERSION_PACKAGE).buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# ==============================================================================
# 定义默认目标为 all
.DEFAULT_GOAL := all

# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标
.PHONY: all
all: tidy format build

# ==============================================================================
# 定义其他需要的伪目标

.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/$(EXECUTABLE) $(PROJ_ROOT_DIR)/cmd/$(EXECUTABLE)/main.go

.PHONY: format
format: # 格式化 Go 源码.
	gofmt -s -w ./

.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	go mod tidy

.PHONY: clean
clean: # 清理构建产物、临时文件等.
	-rm -vrf $(OUTPUT_DIR)

.PHONY: protoc
protoc: # 生成 Protobuf 相关代码.
	@echo "===========> Generating protobuf code"
	protoc \
		--proto_path=$(APIROOT_DIR) \
		--proto_path=$(PROJ_ROOT_DIR)/third_party/protobuf \
		--go_out=$(APIROOT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(APIROOT_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(shell find $(APIROOT_DIR) -name '*.proto' -print0 | xargs -0)