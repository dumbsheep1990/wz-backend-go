# wz-backend-go 更新日志 (2025-04-30)

## 本次更新内容：基于Nacos的服务注册与发现、健康检查和服务实例管理实现

### 实现概述

今日更新主要针对以下功能：

- 服务注册与发现机制
- 健康检查和自动恢复机制
- 服务实例管理

### 功能实现清单

1. **服务注册与发现模块**
   
   - 基于Nacos的服务注册实现
   - 服务发现与查询机制
   - 服务变更通知与订阅
   - 自动化的服务实例管理

2. **健康检查系统**
   
   - HTTP健康检查接口
   - 可扩展的健康检查机制
   - 自动化的健康状态更新
   - 支持存活检查(liveness)和就绪检查(readiness)

3. **服务实例管理**
   
   - 实例状态管理
   - 实例元数据管理
   - 自动心跳机制
   - 服务实例生命周期管理

4. **服务注册与发现API接口**
   
   - 服务列表查询接口：`GET /api/v1/registry/services`
   - 服务实例查询接口：`GET /api/v1/registry/services/:serviceName/instances`
   - 服务实例注册接口：`POST /api/v1/registry/services/:serviceName/instances`
   - 服务实例注销接口：`DELETE /api/v1/registry/services/:serviceName/instances/:instanceId`
   - 实例状态更新接口：`PATCH /api/v1/registry/services/:serviceName/instances/:instanceId/status`
   - 健康状态查询接口：`GET /api/v1/registry/health/status`
   - 手动健康检查接口：`POST /api/v1/registry/health/check`
   - 服务依赖关系接口：`GET /api/v1/registry/services/:serviceName/dependencies`

### 代码结构

```
internal/
├── registry/
│   ├── registry.go  # 服务注册与发现核心实现
│   ├── health.go    # 健康检查实现
│   ├── instance.go  # 服务实例管理
│   └── example.go   # 使用示例
└── delivery/
    └── http/
        └── internal/
            ├── types/
            │   └── registry_types.go    # API请求响应类型定义
            ├── handler/
            │   └── registryhandler.go   # API路由和处理器
            └── logic/
                ├── listserviceslogic.go             # 服务列表查询逻辑
                ├── getserviceinstanceslogic.go      # 获取服务实例逻辑
                ├── registerinstancelogic.go         # 注册服务实例逻辑
                ├── deregisterinstancelogic.go       # 注销服务实例逻辑
                ├── updateinstancestatuslogic.go     # 更新实例状态逻辑
                ├── gethealthstatuslogic.go          # 获取健康状态逻辑
                ├── triggerhealthchecklogic.go       # 触发健康检查逻辑
                ├── getservicedependencieslogic.go   # 获取服务依赖逻辑
                └── errors.go                        # 错误定义
```

### 模块详细说明

#### 1. 服务注册与发现模块 (registry.go)

- **NacosRegistry**: 封装了与Nacos服务器交互的核心逻辑
- **ServiceRegistry接口**: 定义了服务注册与发现的标准接口，支持扩展其他实现
- **服务订阅机制**: 支持实时获取服务列表变更通知
- **缓存优化**: 本地缓存服务实例信息，提高查询效率

关键接口：

```go
// ServiceRegistry 定义服务注册接口
type ServiceRegistry interface {
    Register(serviceName, ip string, port int, meta map[string]string) error
    Deregister(serviceName, ip string, port int) error
    GetService(serviceName string) ([]ServiceInstance, error)
    Subscribe(serviceName string, callback func(instances []ServiceInstance)) error
    Unsubscribe(serviceName string) error
    Close() error
}
```

#### 2. 健康检查模块 (health.go)

- **HealthCheckServer**: HTTP健康检查服务器，提供标准健康检查接口
- **可定制检查项**: 支持注册自定义健康检查函数
- **自动检查**: 周期性自动执行所有注册的健康检查项
- **标准接口**: 提供符合云原生规范的健康检查接口(/health, /health/live, /health/ready)

关键接口：

```go
// HealthChecker 定义健康检查器接口
type HealthChecker interface {
    Start(ctx context.Context) error
    Stop() error
    RegisterCheck(name string, check CheckFunc)
    DeregisterCheck(name string)
    GetStatus() map[string]CheckResult
}
```

#### 3. 服务实例管理模块 (instance.go)

- **NacosInstanceManager**: 管理服务实例的状态和元数据
- **实例状态机制**: 支持多种实例状态(UP, DOWN, STARTING等)
- **自动心跳**: 周期性发送心跳保持实例活跃
- **元数据管理**: 支持丰富的实例元数据，便于服务治理

关键接口：

```go
// InstanceManager 服务实例管理器接口
type InstanceManager interface {
    RegisterInstance(serviceName string, instanceInfo *InstanceInfo) error
    DeregisterInstance(serviceName string, instanceID string) error
    GetInstances(serviceName string) ([]*InstanceInfo, error)
    GetInstance(serviceName, instanceID string) (*InstanceInfo, error)
    UpdateInstanceStatus(serviceName, instanceID string, status InstanceStatus) error
    StartHeartbeat(serviceName string, instanceID string, interval time.Duration) error
    StopHeartbeat(serviceName string, instanceID string) error
}
```

### 使用示例 (example.go)

1. 服务实例注册
2. 健康检查配置
3. 心跳机制设置
4. 服务发现使用
5. 优雅关闭流程

### 配置说明

Nacos配置示例：

```go
config := &NacosConfig{
    ServerAddr: "127.0.0.1",          // Nacos服务器地址
    ServerPort: 8848,                 // Nacos服务器端口
    Namespace:  "public",             // 命名空间
    Group:      "DEFAULT_GROUP",      // 服务分组
    LogDir:     "logs/nacos",         // 日志目录
    CacheDir:   "cache/nacos",        // 缓存目录
    LogLevel:   "info",               // 日志级别
    Username:   "nacos",              // 用户名(可选)
    Password:   "nacos",              // 密码(可选)
}
```

### 最佳实践

1. **服务命名**
   
   - 使用有意义且一致的服务命名规则
   - 建议格式: `<组织>.<应用>.<服务>`，如 `wz.user.auth-service`

2. **实例元数据**
   
   - 合理利用元数据传递关键信息
   - 推荐元数据项: 版本、环境、区域、负责人等

3. **健康检查**
   
   - 实现细粒度的健康检查
   - 区分liveness(存活检查)和readiness(就绪检查)

4. **优雅关闭**
   
   - 先更新状态为下线中
   - 停止接收新请求
   - 等待处理中请求完成
   - 注销服务实例

### 后续计划

1. **监控集成**
   
   - 与Prometheus集成，实现服务指标监控
   - 添加更多运行时指标收集

2. **熔断机制**
   
   - 集成断路器模式
   - 实现服务降级策略

3. **数据一致性**
   
   - 优化服务实例缓存策略
   - 增强数据一致性保障

4. **自动恢复**
   
   - 增强实例自动恢复机制
   - 添加故障自愈能力
