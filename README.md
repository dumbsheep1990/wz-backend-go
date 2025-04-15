# wz-backend-go

万知微服务后端项目，基于go-zero框架实现。

## 项目结构

```
wz-backend-go/
├── api/                # API定义目录
│   ├── http/           # HTTP API定义
│   │   ├── user.api    # 用户服务API定义
│   │   └── content.api # 内容服务API定义
│   └── rpc/            # RPC服务定义
│       ├── user.proto  # 用户服务proto文件
│       └── content.proto # 内容服务proto文件
├── cmd/                # 程序入口
├── configs/            # 配置文件
│   ├── user.yaml       # 用户服务配置
│   └── content.yaml    # 内容服务配置
├── internal/           # 内部代码，不对外暴露
│   ├── delivery/       # 传输层
│   ├── domain/         # 领域模型
│   │   └── model/      # 数据模型
│   ├── repository/     # 数据仓库
│   │   └── sql/        # SQL相关定义
│   └── service/        # 业务逻辑
└── pkg/                # 公共包
```

## 技术栈

- Go 1.22
- [go-zero](https://github.com/zeromicro/go-zero): 微服务框架
- gRPC: 服务间通信
- REST API: 对外接口
- MySQL: 数据存储
- JWT: 认证

## 主要功能

### 用户服务

用户服务负责用户相关功能，包括：

- 用户注册/登录
- 用户信息管理
- 实名认证
- 企业认证
- 用户行为分析

### 内容服务

内容服务负责内容相关功能，包括：

- 分类管理
- 帖子管理
- 评论管理
- 内容状态管理
- 热门内容管理

## 快速开始

### 环境要求

- Go 1.22 或以上
- MySQL 5.7 或以上
- etcd (用于服务发现)

### 安装依赖

```bash
go mod tidy
```

### 构建

```bash
# 构建用户服务
make build-user

# 构建内容服务
make build-content

# 构建全部服务
make build
```

### 运行

1. 首先导入数据库结构
```bash
mysql -u用户名 -p密码 < internal/repository/sql/schema.sql
```

2. 修改配置文件
```bash
# 根据你的环境修改配置文件中的数据库连接信息
vi configs/user.yaml
vi configs/content.yaml
```

3. 启动服务
```bash
# 启动用户服务
./bin/user-api -f configs/user.yaml

# 启动内容服务
./bin/content-api -f configs/content.yaml
```

## 当前实现的API

### 用户服务

用户服务提供的主要API有：

- POST /api/v1/auth/register - 用户注册
- POST /api/v1/auth/login - 用户登录
- GET /api/v1/users/info - 获取用户信息
- PUT /api/v1/users/info - 更新用户信息
- POST /api/v1/users/verify - 用户实名认证
- POST /api/v1/users/verify-company - 企业认证
- GET /api/v1/users/behavior - 获取用户行为数据

### 内容服务

内容服务提供的主要API有：

- [分类相关] - CRUD操作 /api/v1/categories
- [帖子相关] - CRUD操作 /api/v1/posts 
- [评论相关] - CRUD操作 /api/v1/reviews
- [内容状态管理] - /api/v1/content/status
- [热门内容管理] - /api/v1/content/hot


### grpc代码生成

```bash
# 根据API文件生成HTTP服务代码
goctl api go -api api/http/user.api -dir ./internal/delivery/http

# 根据proto文件生成RPC服务代码
goctl rpc protoc api/rpc/user.proto --go_out=./internal/delivery/rpc --go-grpc_out=./internal/delivery/rpc --zrpc_out=./internal/delivery/rpc
```

