# 后端网关API文档

## 概述

本文档描述了微众后端网关(gateway)服务的API接口实现情况，包括已实现的接口和待完善的功能。

## 当前架构状态

微众后端系统采用微服务架构，通过API网关将请求路由到对应的微服务。系统由以下主要组件构成：

1. **API网关**：处理认证、路由、日志记录等横切关注点
2. **微服务集群**：用户服务、内容服务、交易服务等业务服务
3. **服务注册与发现**：基于Nacos实现的服务管理系统
4. **分布式追踪系统**：基于OpenTelemetry的链路追踪

## 已实现的网关接口

### 1. 服务注册与发现API

| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/registry/services` | GET | 获取服务列表 |
| `/api/v1/registry/services/:serviceName/instances` | GET | 获取服务实例列表 |
| `/api/v1/registry/services/:serviceName/instances` | POST | 注册服务实例 |
| `/api/v1/registry/services/:serviceName/instances/:instanceId` | DELETE | 注销服务实例 |
| `/api/v1/registry/services/:serviceName/instances/:instanceId/status` | PATCH | 更新实例状态 |
| `/api/v1/registry/health/status` | GET | 获取健康状态 |
| `/api/v1/registry/health/check` | POST | 手动健康检查 |
| `/api/v1/registry/services/:serviceName/dependencies` | GET | 获取服务依赖关系 |

### 2. 用户服务API

#### 认证类接口
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/auth/register` | POST | 用户注册 |
| `/api/v1/auth/login` | POST | 用户登录 |

#### 用户管理接口
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/users/info` | GET | 获取用户信息 |
| `/api/v1/users/info` | PUT | 更新用户信息 |
| `/api/v1/users/verify` | POST | 用户验证 |
| `/api/v1/users/verify-company` | POST | 企业验证 |
| `/api/v1/users/behavior` | GET | 获取用户行为 |

### 3. 内容服务API

#### 分类管理
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/categories` | POST | 创建分类 |
| `/api/v1/categories/:category_id` | PUT | 更新分类 |
| `/api/v1/categories/:category_id` | DELETE | 删除分类 |
| `/api/v1/categories/:category_id` | GET | 获取分类详情 |
| `/api/v1/categories` | GET | 获取分类列表 |

#### 帖子管理
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/posts` | POST | 创建帖子 |
| `/api/v1/posts/:post_id` | PUT | 更新帖子 |
| `/api/v1/posts/:post_id` | DELETE | 删除帖子 |
| `/api/v1/posts/:post_id` | GET | 获取帖子详情 |
| `/api/v1/posts` | GET | 获取帖子列表 |

#### 评论管理
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/reviews` | POST | 创建评论 |
| `/api/v1/reviews/:review_id` | PUT | 更新评论 |
| `/api/v1/reviews/:review_id` | DELETE | 删除评论 |
| `/api/v1/reviews/:review_id` | GET | 获取评论详情 |
| `/api/v1/reviews` | GET | 获取评论列表 |

#### 内容状态管理
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/content/status` | POST | 更新内容状态 |
| `/api/v1/content/status/:resource_type/:resource_id` | GET | 获取内容状态 |

#### 热门内容管理
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/content/hot` | GET | 获取热门内容 |
| `/api/v1/content/hot` | POST | 设置热门内容 |

### 4. 交互服务API

| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/interaction/like` | POST | 点赞 |
| `/api/v1/interaction/unlike` | POST | 取消点赞 |
| `/api/v1/interaction/comment` | POST | 发表评论 |
| `/api/v1/interaction/comment/:id` | DELETE | 删除评论 |
| `/api/v1/interaction/follow` | POST | 关注 |
| `/api/v1/interaction/unfollow` | POST | 取消关注 |
| `/api/v1/interaction/report` | POST | 举报 |

### 5. AI服务API

| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/ai/recommend` | POST | 获取推荐内容 |
| `/api/v1/ai/review` | POST | 内容审核 |
| `/api/v1/ai/chat` | POST | 客服对话 |

### 6. 通知服务API

| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/notification/ws` | GET | WebSocket连接 |
| `/api/v1/notification/email` | POST | 发送邮件通知 |
| `/api/v1/notification/sms` | POST | 发送短信通知 |
| `/api/v1/notification/list` | GET | 获取通知列表 |
| `/api/v1/notification/read/:id` | PUT | 标记通知为已读 |

### 7. 文件服务API

| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/file/upload` | POST | 上传文件 |
| `/api/v1/file/:fileId` | GET | 获取文件信息 |
| `/api/v1/file/:fileId` | DELETE | 删除文件 |
| `/api/v1/file/list` | GET | 获取文件列表 |

### 8. 搜索服务API

| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/search` | GET | 搜索 |
| `/api/v1/suggest` | GET | 获取搜索建议 |
| `/api/v1/hot` | GET | 获取热搜词 |
| `/api/v1/hot` | POST | 添加热搜词 |
| `/api/v1/hot/:id` | DELETE | 删除热搜词 |
| `/api/v1/search/logs` | GET | 获取搜索记录 |

### 9. 统计服务API

| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/statistics/user/behavior` | GET | 获取用户行为统计 |
| `/api/v1/statistics/content/popularity` | GET | 获取内容流行度统计 |
| `/api/v1/statistics/content/hot` | GET | 获取热门内容 |
| `/api/v1/statistics/user/profile` | GET | 获取用户画像 |

### 10. 交易服务API

#### 订单管理
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/orders` | POST | 创建订单 |
| `/api/v1/orders/:order_id` | GET | 获取订单详情 |
| `/api/v1/orders/:order_id/cancel` | POST | 取消订单 |
| `/api/v1/orders` | GET | 获取订单列表 |
| `/api/v1/orders/callback` | POST | 支付回调处理 |

#### 退款管理
| 接口路径 | 方法 | 描述 |
|---------|-----|------|
| `/api/v1/refunds` | POST | 创建退款 |
| `/api/v1/refunds/:refund_id` | GET | 获取退款详情 |
| `/api/v1/refunds` | GET | 获取退款列表 |
| `/api/v1/refunds/process` | POST | 处理退款 |

## 网关中间件功能

当前网关已实现以下中间件功能：

1. **请求日志记录**：记录请求路径、方法、状态码和响应时间
2. **错误处理**：统一错误响应格式
3. **请求ID生成**：为每个请求生成唯一ID，便于追踪
4. **认证拦截**：验证JWT令牌
5. **OpenTelemetry追踪**：集成分布式追踪功能

## 未完成/需要完善的功能

### 1. API网关完善

- [ ] **完善API网关与gRPC服务的整合**：实现HTTP请求到gRPC请求的转换
- [ ] **路由规则管理系统**：支持动态配置路由规则
- [ ] **增强认证机制**：支持多种认证方式（OAuth2.0、API Key等）
- [ ] **细粒度权限控制**：基于角色和资源的权限控制
- [ ] **API文档生成和管理**：自动生成OpenAPI文档

### 2. 服务注册与发现完善

- [ ] **健康检查机制**：增强健康检查策略
- [ ] **服务实例管理UI**：提供服务实例可视化管理界面

### 3. 限流和熔断

- [ ] **请求限流中间件**：支持基于IP和用户的限流策略
- [ ] **熔断器实现**：检测并防止服务调用失败导致的级联故障
- [ ] **降级策略**：服务不可用时提供降级响应

### 4. 监控与告警

- [ ] **完整的分布式追踪系统**：请求全链路追踪
- [ ] **监控告警机制**：服务异常指标监控和告警
- [ ] **性能指标收集**：收集和分析性能数据

### 5. 多租户数据隔离

- [ ] **完善租户隔离实现**：确保不同租户数据安全隔离
- [ ] **租户资源管理**：管理租户资源配额
- [ ] **租户自定义配置**：支持租户级别的自定义设置

### 6. 部署与运维

- [ ] **容器化配置**：提供Docker和Kubernetes配置
- [ ] **CI/CD流程**：自动化构建、测试和部署流程
- [ ] **环境配置管理**：管理不同环境的配置

## 技术实现细节

### 网关组件结构

网关服务基于以下组件构建：

1. **路由器**：将请求路由到对应的后端服务
2. **拦截器链**：处理横切关注点，如认证、日志等
3. **负载均衡器**：在多个实例之间分发请求
4. **服务发现客户端**：与Nacos交互获取服务实例信息
5. **追踪上下文传播**：传递追踪信息

### 请求处理流程

1. 接收HTTP请求
2. 生成请求ID并启动追踪
3. 记录请求日志
4. 执行认证和授权检查
5. 查询目标服务实例
6. 转发请求到目标服务
7. 接收服务响应
8. 处理响应并返回给客户端
9. 记录响应日志和追踪信息

## 后续开发计划

### 近期计划（2025-05-02至2025-05-12）

1. 完成API网关与gRPC服务的整合
2. 实现动态路由规则管理
3. 增强认证和授权机制
4. 实现API文档自动生成

### 中期计划（2025-05-13至2025-05-27）

1. 实现请求限流和熔断机制
2. 完善多租户数据隔离
3. 实现监控系统和告警机制
4. 完善服务实例管理界面

### 长期计划（2025-05-28至2025-06-15）

1. 优化部署和运维自动化
2. 实现数据库管理和优化策略
3. 完善测试和性能优化
4. 增强系统的安全性和可靠性
