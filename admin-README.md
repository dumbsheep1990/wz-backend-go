# wz-backend-go 后台管理API实现文档

## 本次更新内容：实现后台管理API框架与核心功能

### 实现概述

本次开发主要针对以下功能：

- 设计并实现统一的后台管理API接口
- 建立基于角色的权限控制机制
- 完善用户、租户、仪表盘等核心管理功能
- 构建可扩展的管理服务框架

### 功能实现清单

1. **后台管理API定义**
   
   - 用户管理API（查询、创建、更新、删除）
   - 租户管理API（查询、创建、更新、删除）
   - 系统仪表盘API（数据概览、系统设置）
   - 统一的请求与响应结构

2. **权限控制机制**
   
   - JWT认证中间件
   - 管理员角色检查中间件
   - 基于角色的接口访问控制
   - 分层的权限验证流程

3. **数据访问层实现**
   
   - 统一的仓库接口定义
   - SQL实现的数据存储
   - 动态查询与过滤条件构建
   - 安全的数据操作控制

4. **业务逻辑层设计**
   
   - 服务层与数据层分离
   - 业务逻辑集中处理
   - 微服务间的调用整合
   - 请求参数验证与转换

### 代码结构

```
wz-backend-go/
├── api/
│   └── http/
│       ├── admin_user.api        # 用户管理API定义
│       ├── admin_tenant.api      # 租户管理API定义
│       └── admin_dashboard.api   # 仪表盘API定义
├── cmd/
│   └── admin/
│       └── main.go               # 后台服务入口
├── configs/
│   └── admin.yaml               # 后台服务配置
├── internal/
│   ├── config/
│   │   └── admin_config.go      # 配置结构定义
│   ├── handler/
│   │   └── admin/
│   │       ├── admin.go         # 路由注册
│   │       ├── user/            # 用户管理处理器
│   │       ├── tenant/          # 租户管理处理器
│   │       └── dashboard/       # 仪表盘处理器
│   ├── logic/
│   │   └── admin/
│   │       ├── user/            # 用户业务逻辑
│   │       ├── tenant/          # 租户业务逻辑
│   │       └── dashboard/       # 仪表盘业务逻辑
│   ├── middleware/
│   │   ├── admin_check.go       # 管理员检查中间件
│   │   └── jwt_auth.go          # JWT认证中间件
│   ├── repository/
│   │   ├── repository.go        # 数据仓库接口
│   │   └── sql_user_repository.go # 用户仓库SQL实现
│   └── svc/
│       └── admin_service_context.go # 服务上下文
```

### 模块详细说明

#### 1. 后台管理API定义 (api/http/)

- **用户管理API**: 创建、查询、更新、删除用户
- **租户管理API**: 创建、查询、更新、删除租户
- **仪表盘API**: 数据概览、系统设置管理
- **请求响应结构**: 统一的数据交互格式

关键结构（用户管理API示例）：

```go
type (
    // 用户分页查询请求
    UserListReq {
        Page      int    `form:"page,default=1"`
        PageSize  int    `form:"pageSize,default=10"`
        Username  string `form:"username,optional"`
        Email     string `form:"email,optional"`
        Phone     string `form:"phone,optional"`
        Status    int    `form:"status,optional"`
        Role      string `form:"role,optional"`
        StartTime string `form:"startTime,optional"`
        EndTime   string `form:"endTime,optional"`
    }

    // 用户信息详情
    UserDetail {
        ID               int64  `json:"id"`
        Username         string `json:"username"`
        Email            string `json:"email"`
        Phone            string `json:"phone"`
        Role             string `json:"role"`
        Status           int32  `json:"status"`
        IsVerified       bool   `json:"is_verified"`
        IsCompanyVerified bool  `json:"is_company_verified"`
        DefaultTenantID  int64  `json:"default_tenant_id"`
        CreatedAt        string `json:"created_at"`
        UpdatedAt        string `json:"updated_at"`
    }

    // 用户列表响应
    UserListResp {
        Total int64        `json:"total"`
        List  []UserDetail `json:"list"`
    }
)
```

#### 2. 权限控制机制 (internal/middleware/)

- **JWT认证**: 基于JWT令牌验证用户身份
- **管理员检查**: 验证用户是否具有管理员权限
- **角色信息传递**: 在上下文中传递用户角色信息
- **权限继承关系**: 设置管理员角色的层次结构

关键代码：

```go
// AdminCheck 检查用户是否有管理员权限
func AdminCheck(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 从JWT中获取用户角色信息
        role := r.Context().Value("role")
        if role == nil {
            httpx.Error(w, errors.New("未找到用户角色信息"), http.StatusForbidden)
            return
        }

        // 检查角色是否是管理员
        roleStr, ok := role.(string)
        if !ok || (roleStr != "platform_admin" && !strings.HasPrefix(roleStr, "admin")) {
            httpx.Error(w, errors.New("没有管理员权限"), http.StatusForbidden)
            return
        }

        next(w, r)
    }
}
```

#### 3. 数据访问层 (internal/repository/)

- **统一接口**: 定义标准数据操作接口
- **SQL实现**: 基于SQL的仓库实现
- **动态查询**: 支持灵活的过滤和分页
- **事务支持**: 确保数据一致性

关键接口与实现：

```go
// UserRepository 用户数据仓库接口
type UserRepository interface {
    GetUserById(ctx context.Context, id int64) (*User, error)
    GetUserList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*User, int64, error)
    CreateUser(ctx context.Context, user *User) (int64, error)
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id int64) error
}

// SqlUserRepository SQL用户仓库实现
type SqlUserRepository struct {
    conn sqlx.SqlConn
}

// GetUserList 获取用户列表
func (r *SqlUserRepository) GetUserList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*User, int64, error) {
    // 构建动态查询条件
    whereClause := "1=1"
    args := []interface{}{}

    // 动态添加查询条件
    if filters != nil {
        for key, value := range filters {
            // 根据不同字段添加不同的查询条件
            // ...
        }
    }

    // 执行查询
    // ...

    return users, count, nil
}
```

#### 4. 业务逻辑层 (internal/logic/admin/)

- **业务逻辑集中**: 在逻辑层处理业务规则
- **参数验证**: 验证输入参数的合法性
- **错误处理**: 统一的错误处理机制
- **服务协调**: 协调多个数据操作

关键逻辑示例：

```go
// GetUserList 获取用户列表
func (l *GetUserListLogic) GetUserList(req *UserListReq) (*UserListResp, error) {
    // 构建查询过滤条件
    filters := make(map[string]interface{})
    if req.Username != "" {
        filters["username"] = req.Username
    }
    if req.Email != "" {
        filters["email"] = req.Email
    }
    // ...其他条件

    // 调用仓库层查询数据
    users, total, err := l.svcCtx.UserRepo.GetUserList(l.ctx, req.Page, req.PageSize, filters)
    if err != nil {
        return nil, fmt.Errorf("查询用户列表失败: %v", err)
    }

    // 转换数据格式
    var list []UserDetail
    for _, user := range users {
        // 数据转换逻辑
        // ...
    }

    return &UserListResp{
        Total: total,
        List:  list,
    }, nil
}
```

### 服务集成方式

1. **配置生成**:
   
   项目使用标准配置文件 `configs/admin.yaml` 管理服务配置。包括：
   
   - 服务监听地址和端口
   - 数据库连接信息
   - JWT认证密钥
   - 微服务调用配置

2. **启动服务**:
   
   ```bash
   # 从主目录启动
   go build -o cmd/admin/admin cmd/admin/main.go
   ./cmd/admin/admin -f configs/admin.yaml
   ```

3. **API访问方式**:
   
   - 所有API需要JWT认证，通过Authorization头传递
   - 请求格式：`Authorization: Bearer {token}`
   - 管理API统一前缀：`/api/v1/admin/`

### 后台管理功能列表

1. **用户管理**:
   
   - 获取用户列表: `GET /api/v1/admin/users`
   - 获取用户详情: `GET /api/v1/admin/users/:id`
   - 创建用户: `POST /api/v1/admin/users`
   - 更新用户: `PUT /api/v1/admin/users/:id`
   - 删除用户: `DELETE /api/v1/admin/users/:id`

2. **租户管理**:
   
   - 获取租户列表: `GET /api/v1/admin/tenants`
   - 获取租户详情: `GET /api/v1/admin/tenants/:id`
   - 创建租户: `POST /api/v1/admin/tenants`
   - 更新租户: `PUT /api/v1/admin/tenants/:id`
   - 删除租户: `DELETE /api/v1/admin/tenants/:id`

3. **仪表盘功能**:
   
   - 获取概览数据: `GET /api/v1/admin/dashboard/overview`
   - 获取系统设置: `GET /api/v1/admin/settings`
   - 更新系统设置: `PUT /api/v1/admin/settings`

### 后续开发计划

1. **功能扩展**:
   
   - 完善租户管理功能的实现
   - 增加内容管理后台接口
   - 实现交易数据管理功能
   - 添加系统日志与审计功能

2. **性能优化**:
   
   - 添加缓存层减轻数据库负担
   - 优化查询以提高响应速度
   - 增加数据预加载与批处理功能

3. **安全增强**:
   
   - 完善基于角色的权限控制
   - 实现操作日志记录与追踪
   - 加强数据验证与防护措施

4. **集成优化**:
   
   - 与前端后台管理系统深度集成
   - 优化API响应格式与前端需求匹配
   - 提供更丰富的数据分析与展示功能 