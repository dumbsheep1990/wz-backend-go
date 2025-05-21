# 万知相同社区和群组微服务

这是"万知相同"平台的社区与群组功能的后端微服务实现。此服务提供了创建和管理社区、群组、帖子及评论的完整功能。

## 功能特点

- 社区(Community)管理：创建、查询、更新和删除社区
- 群组(Group)管理：在社区内创建群组，用户可加入/退出群组
- 帖子(Post)管理：用户可以在社区或群组内发布帖子
- 评论(Comment)系统：支持帖子评论及嵌套回复
- 点赞及查看统计：跟踪帖子的点赞数和查看次数
- RESTful API 接口和 gRPC 服务

## 系统架构

本服务采用微服务架构设计，提供两种访问方式：

1. **gRPC服务**：用于内部微服务间通信，端口50054
2. **REST API**：用于外部客户端访问，端口8084

## 技术栈

- **语言**: Go
- **Web框架**: Gin
- **RPC框架**: gRPC
- **数据模型**: Protobuf
- **存储**: 内存存储（可扩展至数据库存储）

## API接口说明

### 社区管理

- `GET /api/v1/communities` - 获取社区列表
- `GET /api/v1/communities/:id` - 获取社区详情
- `POST /api/v1/communities` - 创建社区
- `PUT /api/v1/communities/:id` - 更新社区
- `DELETE /api/v1/communities/:id` - 删除社区

### 群组管理

- `GET /api/v1/groups` - 获取群组列表，支持按社区过滤
- `GET /api/v1/groups/:id` - 获取群组详情
- `POST /api/v1/groups` - 创建群组
- `PUT /api/v1/groups/:id` - 更新群组
- `DELETE /api/v1/groups/:id` - 删除群组
- `POST /api/v1/groups/:id/join` - 加入群组
- `POST /api/v1/groups/:id/leave` - 退出群组

### 帖子管理

- `GET /api/v1/posts` - 获取帖子列表，支持按社区或群组过滤
- `GET /api/v1/posts/:id` - 获取帖子详情
- `POST /api/v1/posts` - 创建帖子
- `PUT /api/v1/posts/:id` - 更新帖子
- `DELETE /api/v1/posts/:id` - 删除帖子
- `POST /api/v1/posts/:id/like` - 点赞或取消点赞
- `POST /api/v1/posts/:id/view` - 记录帖子被查看

### 评论管理

- `GET /api/v1/comments` - 获取评论列表，支持按帖子和父评论过滤
- `POST /api/v1/comments` - 创建评论或回复
- `DELETE /api/v1/comments/:id` - 删除评论

## 数据模型

### 社区 (Community)

```json
{
  "id": "社区ID",
  "name": "社区名称",
  "description": "社区描述",
  "owner_id": "创建者ID",
  "owner_name": "创建者名称",
  "create_time": "创建时间",
  "status": "状态(ACTIVE/INACTIVE/DELETED)",
  "tags": ["标签1", "标签2"],
  "location": "地区",
  "group_count": 群组数量,
  "member_count": 成员数量,
  "post_count": 帖子数量
}
```

### 群组 (Group)

```json
{
  "id": "群组ID",
  "name": "群组名称",
  "description": "群组描述",
  "community_id": "所属社区ID",
  "owner_id": "创建者ID",
  "owner_name": "创建者名称",
  "create_time": "创建时间",
  "status": "状态(ACTIVE/INACTIVE/DELETED)",
  "members": ["成员ID1", "成员ID2"],
  "tags": ["标签1", "标签2"],
  "member_count": 成员数量,
  "post_count": 帖子数量
}
```

### 帖子 (Post)

```json
{
  "id": "帖子ID",
  "title": "标题",
  "content": "内容",
  "author_id": "作者ID",
  "author_name": "作者名称",
  "community_id": "所属社区ID",
  "group_id": "所属群组ID",
  "create_time": "创建时间",
  "update_time": "更新时间",
  "status": "状态(ACTIVE/INACTIVE/DELETED)",
  "like_count": 点赞数,
  "view_count": 查看数,
  "comment_count": 评论数,
  "tags": ["标签1", "标签2"],
  "images": ["图片URL1", "图片URL2"]
}
```

### 评论 (Comment)

```json
{
  "id": "评论ID",
  "content": "评论内容",
  "author_id": "作者ID",
  "author_name": "作者名称",
  "post_id": "帖子ID",
  "parent_id": "父评论ID",
  "create_time": "创建时间",
  "status": "状态(ACTIVE/DELETED)",
  "like_count": 点赞数
}
```

## 运行服务

启动服务：

```bash
go run main.go
```

服务启动后：
- gRPC服务将在 `:50054` 端口运行
- REST API将在 `:8084` 端口运行

## 使用示例

### 创建社区

```bash
curl -X POST http://localhost:8084/api/v1/communities \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -H "X-User-Name: 张三" \
  -d '{
    "name": "南京市技术交流",
    "description": "南京市技术交流社区，讨论各种IT技术话题",
    "tags": ["技术", "IT", "编程"],
    "location": "江苏省-南京市"
  }'
```

### 在社区内创建群组

```bash
curl -X POST http://localhost:8084/api/v1/groups \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -H "X-User-Name: 张三" \
  -d '{
    "name": "Java开发者",
    "description": "讨论Java相关技术和最佳实践",
    "community_id": "comm-123456789",
    "tags": ["Java", "Spring", "编程"]
  }'
```

### 发布帖子

```bash
curl -X POST http://localhost:8084/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -H "X-User-Name: 张三" \
  -d '{
    "title": "Spring Boot 3.0新特性介绍",
    "content": "Spring Boot 3.0带来了许多令人兴奋的新特性...",
    "community_id": "comm-123456789",
    "group_id": "group-123456789",
    "tags": ["Spring", "Spring Boot", "Java"],
    "images": ["https://example.com/image1.jpg"]
  }'
```

## 后续优化方向

1. **数据持久化**：接入MySQL或PostgreSQL数据库
2. **用户认证**：集成OAuth2.0或JWT完整认证
3. **搜索功能**：添加全文搜索能力
4. **通知系统**：实现消息和通知功能
5. **内容审核**：敏感内容过滤和审核机制
6. **性能优化**：添加缓存层
7. **测试**：完善单元测试和集成测试
