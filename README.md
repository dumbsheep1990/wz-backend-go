# wz-backend-go

万知微服务后端项目，基于go-zero框架实现。采用B2B2C架构，支持多租户模式，租户隔离，SaaS化部署。

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
- JWT: 认证与授权
- 多租户B2B2C架构: 支持子域名访问，租户隔离

## 主要功能

### 多租户平台架构 (2025-04-23更新)

万知文站采用B2B2C多租户架构，支持企业用户通过唯一子域名创建自定义站点，实现多租户隔离。主要特性包括：

- **多租户支持**: 每个企业用户(租户)拥有独立的数据和接口，通过租户ID和子域名实现隔离
- **B2B2C架构**: 企业用户可自定义站点，个人用户访问公共内容
- **总站/分站模式**: 总站提供全局导航和租户发现，分站为租户专属站点
- **租户类型划分**:
  - 企业租户: 主要功能为商品交易
  - 个人租户: 主要功能为博客作品集
  - 教育机构: 主要功能为课程、学术资料库等
- **鉴权控制**: 使用JWT令牌认证，确保租户数据安全访问
- **租户识别**: 通过子域名解析确定租户ID，或通过请求头(X-Tenant-ID)传递

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

### 公共接口服务 (2025-04-23更新)

公共接口服务提供总站和分站共用的接口，包括：

- **总站公共接口**: 提供全局功能，如租户列表、全局导航等
- **分站公共接口**: 在租户子域名上运行，返回租户专属数据
  - 导航相关: 获取租户导航分类
  - 搜索相关: 租户内搜索功能
  - 推荐相关: 获取租户推荐内容
  - 分类详情: 获取租户分类详细信息
  - 静态页面: 获取租户隐私政策等静态内容

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

# 导入租户服务表结构 (2025-04-23更新)
mysql -u用户名 -p密码 < internal/repository/sql/tenant_schema.sql
```

2. 修改配置文件
```bash
# 根据你的环境修改配置文件中的数据库连接信息
vi configs/user.yaml
vi configs/content.yaml
vi configs/search.yaml  # 2025-04-16更新
vi configs/trade.yaml   # 2025-04-16更新
vi configs/public.yaml  # 2025-04-23更新
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

# 启动公共接口服务 (2025-04-23更新)
./bin/public-api -f configs/public.yaml
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

### 多租户平台API (2025-04-23更新)

#### 总站公共接口

- GET /api/total/tenants - 获取租户列表
- GET /api/total/navigation - 获取全局导航分类列表

#### 分站公共接口

- GET /api/navigation - 获取当前租户的导航分类列表
- GET /api/search - 搜索当前租户的内容
- GET /api/recommendations - 获取当前租户的推荐内容
- GET /api/category/{id} - 获取当前租户的分类详情
- GET /api/static/{page} - 获取当前租户的静态页面内容

#### 租户管理接口

- POST /api/v1/users/register-enterprise - 企业注册并创建租户
- POST /api/v1/tenants - 创建新租户
- GET /api/v1/tenants/:id - 获取租户详情
- PUT /api/v1/tenants/:id - 更新租户信息
- GET /api/v1/tenants/:id/users - 获取租户用户列表
- POST /api/v1/tenants/:id/users - 添加用户到租户
- DELETE /api/v1/tenants/:id/users/:userId - 从租户移除用户

### grpc代码生成

```bash
# 根据API文件生成HTTP服务代码
goctl api go -api api/http/user.api -dir ./internal/delivery/http

# 根据proto文件生成RPC服务代码
goctl rpc protoc api/rpc/user.proto --go_out=./internal/delivery/rpc --go-grpc_out=./internal/delivery/rpc --zrpc_out=./internal/delivery/rpc
```

## 多租户架构实现 (2025-04-23更新)

### 认证与授权

- 使用JWT令牌进行认证，token有效期为24小时
- JWT令牌中包含用户ID、角色和租户ID，用于控制访问权限
- 私有接口(如评论提交、购买接口)需要携带有效的JWT令牌
- 公有接口(如导航、搜索、推荐)无需认证，但需通过子域名或租户ID访问特定租户的数据

### 角色定义

- 企业用户: 分为租户管理员(管理站点)和普通企业用户(访问站点功能)
- 个人用户: 访问公共接口和购买服务，无租户关联
- 平台管理员: 管理所有租户，访问总站数据

### 租户隔离实现

- 每个企业用户分配一个唯一租户ID和子域名
- 租户ID通过子域名解析确定，服务器根据请求主机名提取租户ID
- 可选通过请求头(X-Tenant-ID: {tenantId})传递租户ID，适用于无子域名场景
- 租户类型决定功能边界:
  - 企业租户: 商品交易功能
  - 个人租户: 博客作品集功能
  - 教育机构: 课程、学术资料库功能

### 数据库设计

主要租户相关表结构：

- `tenants` - 租户信息表
- `tenant_users` - 租户用户关联表
- `tenant_configs` - 租户配置表
- `tenant_categories` - 租户导航分类表
- `tenant_domains` - 租户域名绑定表
- `tenant_static_pages` - 租户静态页面表

# 后端开发工作总结

## 今日完成工作

### 1. 友情链接管理模块

完成了友情链接管理相关功能的开发，主要涉及以下内容：

- `domain/links.go`: 定义友情链接实体及仓储接口
- `handler/admin/links/links.go`: 实现友情链接相关API处理器
- `types/request.go`: 定义友情链接相关请求结构体
- `service/service.go`: 定义友情链接服务接口

主要功能包括：
- 友情链接列表查询（分页、条件筛选）
- 友情链接详情获取
- 友情链接创建
- 友情链接更新
- 友情链接删除

### 2. 站点配置管理模块

完成了站点配置管理相关功能的开发，主要涉及以下内容：

- `domain/site_config.go`: 定义站点配置实体及仓储接口
- `handler/admin/site_config/site_config.go`: 实现站点配置相关API处理器
- `types/request.go`: 定义站点配置相关请求结构体
- `service/service.go`: 定义站点配置服务接口

主要功能包括：
- 站点配置获取
- 站点配置更新

### 3. 主题/模板管理模块

完成了主题/模板管理相关功能的开发，主要涉及以下内容：

- `domain/theme.go`: 定义主题/模板实体及仓储接口
- `types/request.go`: 定义主题/模板相关请求结构体
- `service/service.go`: 定义主题服务接口

主要功能包括：
- 主题/模板列表查询
- 主题/模板详情获取
- 主题/模板创建
- 主题/模板更新
- 主题/模板删除
- 设置默认主题

### 4. 用户消息系统模块

完成了用户消息系统相关功能的开发，主要涉及以下内容：

- `domain/user_message.go`: 定义用户消息实体及仓储接口
- `types/request.go`: 定义用户消息相关请求结构体
- `service/service.go`: 定义用户消息服务接口

主要功能包括：
- 用户消息创建
- 用户消息列表查询（分页、条件筛选）
- 用户消息标记已读
- 用户消息全部标记已读
- 用户消息删除
- 未读消息统计

### 5. 用户积分系统模块

完成了用户积分系统相关功能的开发，主要涉及以下内容：

- `domain/user_points.go`: 定义用户积分实体及仓储接口
- `types/request.go`: 定义用户积分相关请求结构体
- `service/service.go`: 定义用户积分服务接口

主要功能包括：
- 用户积分记录创建
- 用户积分记录列表查询（分页）
- 用户总积分获取

### 6. 用户收藏功能模块

完成了用户收藏功能相关功能的开发，主要涉及以下内容：

- `domain/user_favorite.go`: 定义用户收藏实体及仓储接口
- `types/request.go`: 定义用户收藏相关请求结构体
- `service/service.go`: 定义用户收藏服务接口

主要功能包括：
- 用户收藏创建
- 用户收藏列表查询（分页、条件筛选）
- 用户收藏删除
- 检查是否已收藏

### 7. 站点统计分析模块

完成了站点统计分析相关功能的开发，主要涉及以下内容：

- `domain/statistics.go`: 定义统计数据实体及仓储接口
- `types/request.go`: 定义统计数据相关请求结构体
- `service/service.go`: 定义统计服务接口

主要功能包括：
- 统计数据记录
- 站点概览数据获取
- 特定类型统计数据获取（按时间段）
- 内容排行榜获取

### 8. 新增文件结构目录

```
wz-backend-go/
├── internal/
│   ├── domain/
│   │   ├── links.go                # 友情链接实体及仓储接口
│   │   ├── site_config.go          # 站点配置实体及仓储接口
│   │   ├── theme.go                # 主题/模板实体及仓储接口
│   │   ├── user_message.go         # 用户消息实体及仓储接口
│   │   ├── user_points.go          # 用户积分实体及仓储接口
│   │   ├── user_favorite.go        # 用户收藏实体及仓储接口
│   │   └── statistics.go           # 统计数据实体及仓储接口
│   ├── handler/
│   │   └── admin/
│   │       ├── links/              # 友情链接处理器
│   │       ├── site_config/        # 站点配置处理器
│   │       ├── theme/              # 主题处理器
│   │       ├── user/               # 用户相关处理器
│   │       └── statistics/         # 统计数据处理器
│   ├── service/
│   │   ├── service.go              # 服务接口定义
│   │   ├── link_service.go         # 友情链接服务实现
│   │   ├── site_config_service.go  # 站点配置服务实现
│   │   ├── theme_service.go        # 主题服务实现
│   │   ├── user_message_service.go # 用户消息服务实现
│   │   ├── user_points_service.go  # 用户积分服务实现
│   │   ├── user_favorite_service.go # 用户收藏服务实现
│   │   └── statistics_service.go   # 统计服务实现
│   └── types/
│       ├── request.go              # 请求参数结构体定义
│       └── response.go             # 响应结构体定义
└── tests/
    └── integration/
        ├── repositories/           # 仓储接口测试
        └── services/               # 服务接口测试
```

## 后续计划

1. 实现各模块的数据访问层（SQL仓储实现）
2. 完善服务层实现
3. 对接前端接口
4. 编写单元测试与集成测试
5. 性能优化与安全加固
