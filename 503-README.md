# wz-backend-go 更新日志 (2025-05-03)

## 本次更新内容：API网关功能增强与扩展

### 实现概述

今日更新主要针对以下功能：

- API网关与gRPC服务的整合
- 动态路由规则管理系统
- 增强的认证和授权机制
- API文档自动生成功能

### 功能实现清单

1. **API网关与gRPC服务整合**
   
   - gRPC客户端连接池管理
   - HTTP到gRPC请求自动转换
   - gRPC反射API支持
   - 健康检查与连接状态管理

2. **动态路由规则管理**
   
   - 运行时路由规则更新
   - RESTful路由管理API
   - 服务注册与注销机制
   - 路由生命周期管理

3. **增强认证和授权机制**
   
   - 多种认证方式支持（JWT、API Key、OAuth2、Basic Auth）
   - 基于角色的访问控制(RBAC)
   - 细粒度资源权限管理
   - 认证提供者可扩展架构

4. **API文档自动生成**
   
   - 基于OpenAPI/Swagger的文档生成
   - 从路由配置自动提取API信息
   - 交互式Swagger UI集成
   - 文档自动更新机制

### 代码结构

```
internal/
└── gateway/
    ├── auth/                # 认证与授权模块
    │   ├── auth_manager.go  # 认证管理器
    │   └── provider.go      # 认证提供者实现
    ├── config/
    │   └── config.go        # 网关配置结构
    ├── docs/
    │   └── generator.go     # API文档生成器
    ├── grpc/
    │   ├── client_manager.go # gRPC客户端管理器
    │   └── converter.go     # HTTP到gRPC转换器
    ├── handler/
    │   ├── grpc_handler.go  # gRPC请求处理器
    │   ├── http_routes.go   # HTTP路由处理器
    │   └── swagger.go       # Swagger文档处理器
    ├── middleware/
    │   ├── auth.go          # 基础认证中间件
    │   └── enhanced_auth.go # 增强认证中间件
    ├── router/
    │   ├── api.go           # 路由管理API
    │   └── dynamic.go       # 动态路由管理器
    └── server/
        └── server.go        # 网关服务器实现
```

### 模块详细说明

#### 1. gRPC集成模块 (grpc/)

- **ClientManager**: 管理gRPC客户端连接池和健康状态
- **Converter**: 处理HTTP到gRPC的转换和元数据传递
- **反射支持**: 利用gRPC反射API动态获取服务描述
- **连接健康检查**: 周期性检查连接状态并自动恢复

关键结构：

```go
// ClientManager 管理gRPC客户端连接
type ClientManager struct {
    clients     map[string]*Client
    config      config.GRPCConfig
    mu          sync.RWMutex
    healthCheck time.Duration
    stopCh      chan struct{}
}

// Client 封装gRPC客户端连接及其元数据
type Client struct {
    Conn         *grpc.ClientConn
    ServiceName  string
    Methods      map[string]MethodInfo
    LastRefresh  time.Time
    Healthy      bool
    ReflectionCl grpc_reflection_v1alpha.ServerReflectionClient
}
```

#### 2. 动态路由模块 (router/)

- **DynamicRouter**: 实现运行时路由规则的管理
- **RouteAPI**: 提供REST API用于管理路由规则
- **路由组管理**: 支持路由前缀和路由组管理
- **路由生命周期**: 处理路由注册、更新和注销

关键接口：

```go
// DynamicRouter 动态路由管理器
type DynamicRouter struct {
    engine      *gin.Engine
    routeGroups map[string]*gin.RouterGroup
    services    map[string]config.ServiceConfig
    handlers    map[string]gin.HandlerFunc
    mu          sync.RWMutex
}

// RegisterService 注册或更新服务
func (r *DynamicRouter) RegisterService(service config.ServiceConfig, 
    registerFunc func(*gin.RouterGroup, config.ServiceConfig))
```

#### 3. 认证与授权模块 (auth/)

- **AuthManager**: 管理多种认证提供者
- **Provider接口**: 定义认证提供者的统一接口
- **多种实现**: 支持JWT、API Key、OAuth2和基本认证
- **RBAC支持**: 基于角色的权限控制机制

关键结构与接口：

```go
// AuthProvider 认证提供者接口
type AuthProvider interface {
    // Authenticate 从请求中提取凭据并验证
    Authenticate(c *gin.Context) (*AuthUser, error)
    
    // GenerateCredentials 生成新的凭据
    GenerateCredentials(user *AuthUser) (map[string]interface{}, error)
    
    // Name 获取提供者名称
    Name() string
}

// AuthManager 认证管理器
type AuthManager struct {
    providers map[string]AuthProvider
    config    config.SecurityConfig
    mu        sync.RWMutex
}
```

#### 4. API文档生成模块 (docs/)

- **APIDocGenerator**: 从服务配置自动生成API文档
- **OpenAPI规范**: 遵循OpenAPI 3.0.3规范
- **认证整合**: 自动添加认证方案信息
- **自动更新**: 支持文档的定期自动更新

关键结构：

```go
// APIDocGenerator API文档生成器
type APIDocGenerator struct {
    swagger    *openapi3.T
    config     config.Config
    outputPath string
    mu         sync.RWMutex
}

// GenerateFromServices 从服务配置生成API文档
func (g *APIDocGenerator) GenerateFromServices() error
```

#### 5. 服务器模块 (server/)

- **Server**: 网关服务器实现，整合所有模块
- **优雅关闭**: 支持优雅启动和关闭
- **信号处理**: 捕获系统信号进行安全退出
- **统一配置**: 集中管理各模块配置

关键结构：

```go
// Server API网关服务器
type Server struct {
    engine         *gin.Engine
    config         config.Config
    httpServer     *http.Server
    authManager    *auth.AuthManager
    dynamicRouter  *router.DynamicRouter
    grpcHandler    *handler.GrpcHandler
    swaggerHandler *handler.SwaggerHandler
}

// Start 启动服务器
func (s *Server) Start() error
```

### 后续优化方向

1. **依赖引入**: 需要添加`github.com/getkin/kin-openapi/openapi3`依赖
2. **单元测试**: 为各模块添加完整的单元测试和集成测试
3. **性能优化**: 进一步优化gRPC连接池和路由匹配性能
4. **服务发现集成**: 将动态路由与Nacos服务发现更深度整合
