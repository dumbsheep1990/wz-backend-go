# wz-backend-go 更新日志 (2025-05-06)

## 本次更新内容：完善多租户数据隔离机制及多平台支持

### 实现概述

今日更新主要针对以下功能：

- 完善基于租户的数据隔离机制
- 增强对多种客户端平台的支持（Web、移动端、UniApp）
- 统一租户上下文的管理与传递
- 优化跨平台API响应适配

### 功能实现清单

1. **租户上下文管理**
   
   - 集中式租户上下文定义与处理
   - 多平台类型识别与支持
   - 租户信息在各层级间的传递
   - 系统内部操作的绕过机制

2. **增强的租户中间件**
   
   - 多来源租户识别（子域名、请求头、URL参数）
   - 平台自动识别与分类
   - 请求路径级别的租户策略控制
   - 日志记录与追踪支持

3. **数据隔离策略优化**
   
   - 多种隔离模式支持（逻辑隔离、物理隔离、共享表）
   - SQL查询自动注入租户条件
   - JOIN查询的租户条件处理
   - 防止跨租户数据操作的安全保障

4. **平台适配响应处理**
   
   - 平台感知的响应适配器架构
   - 不同客户端的数据格式转换
   - 可扩展的适配器注册机制
   - 自定义字段与属性处理

### 代码结构

```
internal/
├── pkg/
│   └── tenantctx/
│       ├── tenant_context.go    # 租户上下文定义
│       ├── tenant_service.go    # 租户服务接口
│       └── response_adapter.go  # 响应适配器
├── gateway/
│   └── middleware/
│       ├── tenant.go            # 基础租户中间件
│       ├── tenant_context.go    # 租户上下文中间件
│       └── enhanced_tenant.go   # 增强版租户中间件
├── repository/
│   ├── tenant_scope.go          # 基础租户数据隔离
│   ├── enhanced_tenant_scope.go # 增强版租户数据隔离
│   └── tenantrepo.go            # 租户仓库操作
├── service/
│   ├── tenant_config_service.go # 租户配置服务
│   └── tenantservice.go         # 租户服务实现
└── delivery/
    └── adapter/
        └── platform_adapters.go # 平台特定适配器
```

### 模块详细说明

#### 1. 租户上下文管理 (internal/pkg/tenantctx/)

- **上下文键定义**: 统一管理租户相关上下文键
- **平台类型识别**: 定义并处理多种客户端平台
- **错误处理**: 标准化租户操作错误
- **服务接口**: 定义租户感知服务的共同接口

关键结构：

```go
// AppPlatform 表示请求来源的平台
type AppPlatform string

const (
    // PlatformWeb 表示来自网页浏览器的请求
    PlatformWeb AppPlatform = "web"
    
    // PlatformMobile 表示来自原生移动应用的请求
    PlatformMobile AppPlatform = "mobile"
    
    // PlatformUniApp 表示来自uniapp应用的请求
    PlatformUniApp AppPlatform = "uniapp"
)

// TenantAwareService 定义需要租户隔离的服务接口
type TenantAwareService interface {
    // WithContext 将包含租户信息的上下文应用到服务操作中
    WithContext(ctx context.Context) TenantAwareService
    
    // GetCurrentTenant 返回服务当前使用的租户ID
    GetCurrentTenant() (string, error)
    
    // IsTenantValid 验证租户ID是否存在且活跃
    IsTenantValid(tenantID string) (bool, error)
}
```

#### 2. 增强的租户中间件 (internal/gateway/middleware/)

- **多来源租户识别**: 从多种来源提取租户信息
- **平台自动检测**: 根据请求特征识别客户端平台
- **上下文注入**: 将租户和平台信息注入请求上下文
- **公共路由处理**: 对公共路由的特殊处理

关键函数：

```go
// EnhancedTenantMiddleware 增强版租户中间件，支持多种平台类型识别
// 同时支持子域名和请求头方式传递租户信息
func EnhancedTenantMiddleware() gin.HandlerFunc {
    // 检测平台类型
    // 提取租户ID
    // 创建和注入上下文
    // 处理公共路由
}

// detectPlatform 根据请求信息检测平台类型
func detectPlatform(c *gin.Context) tenantctx.AppPlatform {
    // 从请求头识别平台
    // 从User-Agent识别设备类型
    // 返回对应平台枚举
}
```

#### 3. 数据隔离策略优化 (internal/repository/)

- **租户Schema定义**: 多种隔离策略支持
- **SQL条件注入**: 自动添加租户过滤条件
- **表名映射**: 物理隔离模式的表名映射
- **安全保障**: 防止跨租户操作的安全机制

关键结构与函数：

```go
// TenantSchema 表示需要租户隔离的数据库架构类型
type TenantSchema int

const (
    // SchemaDefault 表示默认的租户隔离策略，使用tenant_id列隔离
    SchemaDefault TenantSchema = iota
    
    // SchemaShared 表示共享表，不需要租户隔离
    SchemaShared
    
    // SchemaPhysical 表示物理隔离，使用不同的数据库或表
    SchemaPhysical
)

// EnhancedTenantScope 增强版的租户数据隔离
// 支持多种隔离策略和平台特定处理
func EnhancedTenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
    // 检查系统内部操作
    // 获取租户ID
    // 根据隔离策略应用过滤条件
    // 处理特殊操作（删除、更新、JOIN查询）
}
```

#### 4. 平台适配响应处理 (internal/pkg/tenantctx/ & internal/delivery/adapter/)

- **适配器接口**: 定义响应适配器接口
- **注册机制**: 支持多种平台适配器注册
- **上下文感知**: 基于上下文进行响应适配
- **平台特定实现**: 为不同平台提供特定适配

关键接口与实现：

```go
// ResponseAdapter 定义基于平台适配的API响应方法
type ResponseAdapter interface {
    // AdaptResponse 根据平台需求修改响应
    // 如果响应被修改返回true，否则返回false
    AdaptResponse(ctx context.Context, response interface{}) (interface{}, bool)
    
    // SupportsPlatform 检查适配器是否支持给定平台
    SupportsPlatform(platform AppPlatform) bool
}

// MobileResponseAdapter 适配移动平台的响应
type MobileResponseAdapter struct {
    *tenantctx.BaseResponseAdapter
}

// UniAppResponseAdapter 适配UniApp平台的响应
type UniAppResponseAdapter struct {
    *tenantctx.BaseResponseAdapter
}
```

### 优化效果

1. **数据安全性提升**:
   
   - 强化租户间数据隔离
   - 防止跨租户数据泄露
   - 多层次的安全校验机制
   - 集中式的租户身份管理

2. **多平台支持增强**:
   
   - 自动识别Web、移动端和UniApp请求
   - 平台特定的响应格式优化
   - 统一的租户身份验证机制
   - 平台感知的配置与特性控制

3. **开发效率与可维护性提高**:
   
   - 统一的租户上下文管理
   - 可扩展的适配器架构
   - 声明式的隔离策略配置
   - 代码复用与关注点分离

4. **用户体验优化**:
   
   - 根据客户端平台优化响应格式
   - 减少不必要的数据传输
   - 提供平台特定的API接口
   - 支持多种租户识别方式

### 租户配置示例

租户配置示例，支持多平台特定设置：

```go
config := &TenantConfig{
    ID:     "tenant123",
    Name:   "测试租户",
    Status: 1, // 正常
    Features: map[string]bool{
        "content_management": true,
        "e_commerce":        true,
        "analytics":         true,
    },
    PlatformConfigs: map[string]PlatformConfig{
        // Web平台配置
        "web": {
            Enabled:    true,
            Theme:      "default",
            CustomSettings: map[string]string{
                "max_upload_size": "10MB",
            },
        },
        // 移动端配置
        "mobile": {
            Enabled:    true,
            Theme:      "mobile_dark",
            MinVersion: "1.0.0",
            CustomSettings: map[string]string{
                "max_upload_size": "5MB",
                "enable_push":     "true",
            },
        },
        // UniApp配置
        "uniapp": {
            Enabled:    true,
            Theme:      "uniapp_light",
            CustomSettings: map[string]string{
                "max_upload_size": "8MB",
            },
        },
    },
}
```

### 后续计划

- 实现租户资源配额管理
- 添加租户使用量统计与计费
- 完善租户数据迁移与备份机制
- 实现跨租户数据共享与授权
- 优化大规模多租户环境下的性能
