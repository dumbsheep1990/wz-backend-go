.PHONY: all build gen-api gen-rpc clean run-user run-content test

# 项目变量定义
PROJ_ROOT := $(shell pwd)
BIN_DIR := $(PROJ_ROOT)/bin
USER_API_BIN := $(BIN_DIR)/user-api
CONTENT_API_BIN := $(BIN_DIR)/content-api
GO_BUILD := go build -ldflags="-s -w"

# 默认目标：构建所有服务
all: build

# 创建必要的目录
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# 构建所有服务
build: $(BIN_DIR) build-user build-content
	@echo "所有服务构建完成"

# 构建用户服务
build-user: $(BIN_DIR)
	@echo "构建用户服务..."
	$(GO_BUILD) -o $(USER_API_BIN) ./cmd/user-api
	@echo "用户服务构建完成: $(USER_API_BIN)"

# 构建内容服务
build-content: $(BIN_DIR)
	@echo "构建内容服务..."
	$(GO_BUILD) -o $(CONTENT_API_BIN) ./cmd/content-api
	@echo "内容服务构建完成: $(CONTENT_API_BIN)"

# 生成API代码
gen-api:
	@echo "生成HTTP API代码..."
	goctl api go -api api/http/user.api -dir ./internal/delivery/http
	goctl api go -api api/http/content.api -dir ./internal/delivery/http
	@echo "HTTP API代码生成完成"

# 生成RPC代码
gen-rpc:
	@echo "生成RPC代码..."
	goctl rpc protoc api/rpc/user.proto --go_out=./internal/delivery/rpc --go-grpc_out=./internal/delivery/rpc --zrpc_out=./internal/delivery/rpc
	goctl rpc protoc api/rpc/content.proto --go_out=./internal/delivery/rpc --go-grpc_out=./internal/delivery/rpc --zrpc_out=./internal/delivery/rpc
	@echo "RPC代码生成完成"

# 生成所有代码
gen: gen-api gen-rpc
	@echo "所有代码生成完成"

# 运行用户服务
run-user:
	@echo "启动用户服务..."
	go run ./cmd/user-api/main.go -f ./configs/user.yaml

# 运行内容服务
run-content:
	@echo "启动内容服务..."
	go run ./cmd/content-api/main.go -f ./configs/content.yaml

# 运行测试
test:
	go test -v ./...

# 清理编译产物
clean:
	@echo "清理编译产物..."
	rm -rf $(BIN_DIR)
	@echo "清理完成"

# 安装依赖
deps:
	go mod tidy

# 显示帮助信息
help:
	@echo "可用目标:"
	@echo "  all          - 构建所有服务 (默认)"
	@echo "  build-user   - 构建用户服务"
	@echo "  build-content - 构建内容服务"
	@echo "  gen-api      - 生成HTTP API代码"
	@echo "  gen-rpc      - 生成RPC代码" 
	@echo "  gen          - 生成所有代码"
	@echo "  run-user     - 运行用户服务"
	@echo "  run-content  - 运行内容服务"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理编译产物"
	@echo "  deps         - 安装依赖"
