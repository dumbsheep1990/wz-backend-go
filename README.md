# 万知市场站点建设器 - 后端服务

这是"万知市场"平台自定义站点建设器功能的后端微服务实现。

## 功能特点

- 微服务架构设计，包含多个独立服务
- RESTful API接口
- 支持站点、页面、区块和组件的完整CRUD操作
- 提供站点模板管理
- 支持站点预览和渲染
- 支持多种设备的响应式预览
- 集成API网关

## 系统架构

系统采用微服务架构，包含以下服务：

1. **站点服务(site-service)**: 负责站点基本信息和模板管理
2. **页面服务(page-service)**: 管理站点的页面和区块
3. **组件服务(component-service)**: 提供组件库和管理页面组件
4. **渲染服务(render-service)**: 负责站点和页面的预览和渲染
5. **API网关**: 提供统一的API入口，负责请求路由和认证

## 技术栈

- **语言**: Go
- **Web框架**: Gin
- **认证**: JWT (简化实现)
- **API文档**: Swagger (计划中)
- **数据库**: 目前使用内存数据，计划迁移到MySQL/PostgreSQL

## 目录结构

```
wz-backend-go/
├── api-gateway/                # API网关
├── middleware/                 # 中间件
│   ├── auth.go                 # 认证中间件
│   └── cors.go                 # CORS中间件
├── models/                     # 数据模型
│   ├── site.go                 # 站点模型
│   └── page.go                 # 页面和区块模型
├── services/                   # 微服务
│   ├── site-service/           # 站点服务
│   ├── page-service/           # 页面服务
│   ├── component-service/      # 组件服务
│   └── render-service/         # 渲染服务
├── scripts/                    # 脚本
│   ├── start.bat               # Windows启动脚本
│   └── start.sh                # Linux/Mac启动脚本
├── bin/                        # 编译后的二进制文件
└── README.md                   # 项目说明
```

## 安装与运行

### 前提条件

- Go 1.16+
- Git

### 安装步骤

1. 克隆仓库:

```bash
git clone https://github.com/yourusername/wz-backend-go.git
cd wz-backend-go
```

2. 安装依赖:

```bash
go mod tidy
```

3. 运行服务:

**Windows**:

```bash
scripts\start.bat
```

**Linux/Mac**:

```bash
chmod +x scripts/start.sh
scripts/start.sh
```

所有服务启动后，API网关将在 http://localhost:8080 上运行。

## API接口

### 站点管理

- `GET /api/v1/sites` - 获取站点列表
- `GET /api/v1/sites/:id` - 获取站点详情
- `POST /api/v1/sites` - 创建站点
- `PUT /api/v1/sites/:id` - 更新站点
- `DELETE /api/v1/sites/:id` - 删除站点
- `PUT /api/v1/sites/:id/publish` - 发布站点

### 页面管理

- `GET /api/v1/sites/:siteId/pages` - 获取页面列表
- `GET /api/v1/sites/:siteId/pages/:id` - 获取页面详情
- `POST /api/v1/sites/:siteId/pages` - 创建页面
- `PUT /api/v1/sites/:siteId/pages/:id` - 更新页面
- `DELETE /api/v1/sites/:siteId/pages/:id` - 删除页面
- `PUT /api/v1/sites/:siteId/pages/:id/homepage` - 设置为首页

### 区块和组件

- `GET /api/v1/sites/:siteId/pages/:pageId/sections` - 获取区块列表
- `POST /api/v1/sites/:siteId/pages/:pageId/sections` - 添加区块
- `PUT /api/v1/sites/:siteId/pages/:pageId/sections/:id` - 更新区块
- `DELETE /api/v1/sites/:siteId/pages/:pageId/sections/:id` - 删除区块
- `POST /api/v1/sites/:siteId/pages/:pageId/sections/:sectionId/components` - 添加组件
- `PUT /api/v1/sites/:siteId/pages/:pageId/sections/:sectionId/components/:id` - 更新组件
- `DELETE /api/v1/sites/:siteId/pages/:pageId/sections/:sectionId/components/:id` - 删除组件

### 模板管理

- `GET /api/v1/site-templates` - 获取模板列表
- `GET /api/v1/site-templates/:id` - 获取模板详情

### 预览和渲染

- `GET /api/v1/preview/sites/:siteId` - 预览整个站点
- `GET /api/v1/preview/sites/:siteId/pages/:pageId` - 预览特定页面
- `GET /render/sites/:siteId/:slug` - 渲染站点页面

## 开发计划

- [ ] 集成实际数据库存储
- [ ] 添加用户认证与授权系统
- [ ] 实现文件上传和媒体管理
- [ ] 添加缓存层提高性能
- [ ] 完善错误处理和日志系统
- [ ] 添加单元测试和集成测试
- [ ] 实现CI/CD流程

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交变更 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请查看 LICENSE 文件


