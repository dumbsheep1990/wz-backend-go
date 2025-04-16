# wz-backend-go

万知微服务后端项目，基于go-zero框架实现。

## 项目结构

```
wz-backend-go/
├── api/                # API定义目录
│   ├── http/           # HTTP API定义
│   │   ├── user.api    # 用户服务API定义
│   │   ├── content.api # 内容服务API定义
│   │   ├── search.api  # 搜索服务API定义
│   │   └── trade.api   # 交易服务API定义
│   └── rpc/            # RPC服务定义
│       ├── user.proto  # 用户服务proto文件
│       ├── content.proto # 内容服务proto文件
│       ├── search.proto # 搜索服务proto文件
│       └── trade.proto # 交易服务proto文件
├── cmd/                # 程序入口
├── configs/            # 配置文件
│   ├── user.yaml       # 用户服务配置
│   ├── content.yaml    # 内容服务配置
│   ├── search.yaml     # 搜索服务配置
│   └── trade.yaml      # 交易服务配置
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

### 搜索服务 (2025-04-16更新)

搜索服务负责全站搜索相关功能，包括：

- 全文检索
- 搜索建议
- 热搜管理
- 搜索日志记录
- 搜索统计分析
- 同义词管理
- 停用词管理
- 搜索排序规则

### 交易服务 (2025-04-16更新)

交易服务负责交易相关功能，包括：

- 订单管理
- 支付处理
- 退款管理
- 账户余额
- 交易记录
- 财务报表

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

# 构建搜索服务 (2025-04-16更新)
make build-search

# 构建交易服务 (2025-04-16更新)
make build-trade

# 构建全部服务
make build
```

### 运行

1. 首先导入数据库结构
```bash
# 导入基础表结构
mysql -u用户名 -p密码 < internal/repository/sql/schema.sql

# 导入搜索服务表结构 (2025-04-16更新)
mysql -u用户名 -p密码 < internal/repository/sql/search_schema.sql

# 导入交易服务表结构 (2025-04-16更新)
mysql -u用户名 -p密码 < internal/repository/sql/trade_schema.sql
```

2. 修改配置文件
```bash
# 根据你的环境修改配置文件中的数据库连接信息
vi configs/user.yaml
vi configs/content.yaml
vi configs/search.yaml  # 2025-04-16更新
vi configs/trade.yaml   # 2025-04-16更新
```

3. 启动服务
```bash
# 启动用户服务
./bin/user-api -f configs/user.yaml

# 启动内容服务
./bin/content-api -f configs/content.yaml

# 启动搜索服务 (2025-04-16更新)
./bin/search-api -f configs/search.yaml

# 启动交易服务 (2025-04-16更新)
./bin/trade-api -f configs/trade.yaml
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
### 用户及文章服务（2025-04-15更新）

- [分类相关] - CRUD操作 /api/v1/categories
- [帖子相关] - CRUD操作 /api/v1/posts 
- [评论相关] - CRUD操作 /api/v1/reviews
- [内容状态管理] - /api/v1/content/status
- [热门内容管理] - /api/v1/content/hot

### 搜索服务 (2025-04-16更新)

搜索服务提供的主要API有：

- GET /api/v1/search - 全文搜索
- GET /api/v1/search/suggestions - 获取搜索建议
- GET /api/v1/search/hot - 获取热搜列表
- POST /api/v1/search/hot - 添加热搜
- PUT /api/v1/search/hot/:id - 更新热搜
- DELETE /api/v1/search/hot/:id - 删除热搜
- GET /api/v1/search/statistics - 获取搜索统计数据
- GET /api/v1/search/synonyms - 获取同义词
- PUT /api/v1/search/synonyms - 更新同义词
- GET /api/v1/search/config - 获取搜索配置
- PUT /api/v1/search/config - 更新搜索配置

### 交易服务 (2025-04-16更新)

交易服务提供的主要API有：

- POST /api/v1/orders - 创建订单
- GET /api/v1/orders/:id - 获取订单详情
- GET /api/v1/orders - 获取订单列表
- POST /api/v1/orders/:id/cancel - 取消订单
- POST /api/v1/payments/process - 处理支付
- POST /api/v1/payments/callback - 支付回调
- POST /api/v1/refunds - 申请退款
- GET /api/v1/refunds/:id - 获取退款详情
- GET /api/v1/refunds - 获取退款列表
- POST /api/v1/refunds/:id/process - 处理退款
- GET /api/v1/balance - 获取账户余额
- GET /api/v1/transactions - 获取交易记录
- GET /api/v1/reports/financial - 获取财务报表

### grpc代码生成

```bash
# 根据API文件生成HTTP服务代码
goctl api go -api api/http/user.api -dir ./internal/delivery/http

# 根据proto文件生成RPC服务代码
goctl rpc protoc api/rpc/user.proto --go_out=./internal/delivery/rpc --go-grpc_out=./internal/delivery/rpc --zrpc_out=./internal/delivery/rpc
```
