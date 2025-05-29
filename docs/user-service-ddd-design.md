# 用户服务DDD设计文档

## 1. 概述

本文档描述了基于领域驱动设计(DDD)的用户服务重构设计。用户服务是WZ系统的核心服务之一，负责用户账号管理、认证和授权等功能。

## 2. 领域模型

### 2.1 值对象 (Value Objects)

值对象是没有唯一标识的不可变对象，用于表示领域中的概念。

- **UserID**: 用户ID值对象，使用int64类型表示
- **Username**: 用户名值对象，包含用户名验证逻辑
- **Email**: 邮箱值对象，包含邮箱格式验证逻辑
- **Phone**: 手机号值对象，包含手机号格式验证逻辑
- **UserStatus**: 用户状态值对象，定义了用户的不同状态（未激活、活跃、锁定、已删除）

### 2.2 实体 (Entities)

实体是具有唯一标识的对象，用于表示领域中的核心概念。

- **User**: 用户实体，代表系统中的用户，包含用户的基本信息和行为
- **UserBehavior**: 用户行为实体，记录用户的操作行为

### 2.3 领域事件 (Domain Events)

领域事件用于表示领域中发生的重要事件，有助于实现系统的解耦和扩展性。

- **UserCreatedEvent**: 用户创建事件
- **UserVerifiedEvent**: 用户验证事件
- **UserCompanyVerifiedEvent**: 用户企业认证事件
- **UserPasswordChangedEvent**: 用户密码修改事件
- **UserLoggedInEvent**: 用户登录事件
- **UserBehaviorRecordedEvent**: 用户行为记录事件

### 2.4 仓储接口 (Repository Interfaces)

仓储接口定义了领域对象的持久化操作。

- **UserRepository**: 用户仓储接口，定义了用户实体的CRUD操作
- **EventPublisher**: 事件发布接口，用于发布领域事件

### 2.5 领域服务 (Domain Services)

领域服务封装了不适合放在实体或值对象中的领域逻辑。

- **UserDomainService**: 用户领域服务，包含用户注册、登录、验证等核心业务逻辑

## 3. 应用层设计

应用层负责编排领域对象，实现用户的用例。

### 3.1 数据传输对象 (DTOs)

- **UserDTO**: 用户数据传输对象，用于向外部展示用户信息
- **UserBehaviorDTO**: 用户行为数据传输对象
- **UserRegisterRequest**: 用户注册请求
- **UserLoginRequest**: 用户登录请求
- **ChangePasswordRequest**: 修改密码请求
- **UserBehaviorRequest**: 记录用户行为请求

### 3.2 应用服务 (Application Services)

- **UserApplicationService**: 用户应用服务，协调领域对象完成用户的用例

## 4. 基础设施层设计

基础设施层提供技术实现，支持领域层和应用层。

### 4.1 持久化实现

- **UserRepository**: 用户仓储的SQL实现
- **UnitOfWork**: 工作单元实现，用于事务管理

### 4.2 事件总线实现

- **RabbitMQEventPublisher**: RabbitMQ事件发布器实现
- **LogEventPublisher**: 日志事件发布器实现

## 5. 接口层设计

接口层负责处理外部请求，并将其转换为应用层的调用。

### 5.1 HTTP控制器

- **UserHandler**: 用户HTTP处理器，提供RESTful API

### 5.2 API路由

- `POST /api/v1/users`: 用户注册
- `POST /api/v1/users/login`: 用户登录
- `GET /api/v1/users`: 获取用户列表
- `GET /api/v1/users/:id`: 获取用户详情
- `PUT /api/v1/users/:id/verify`: 验证用户
- `PUT /api/v1/users/:id/verify-company`: 完成企业认证
- `PUT /api/v1/users/:id/password`: 修改密码
- `POST /api/v1/users/:id/behaviors`: 记录用户行为
- `GET /api/v1/users/:id/behaviors`: 获取用户行为列表

## 6. 目录结构

```
internal/
  ├── domain/
  │   └── user/
  │       ├── entity/           // 实体
  │       │   └── user.go
  │       ├── valueobject/      // 值对象
  │       │   └── user_valueobjects.go
  │       ├── event/            // 领域事件
  │       │   └── user_events.go
  │       ├── repository/       // 仓储接口
  │       │   └── user_repository.go
  │       └── service/          // 领域服务
  │           └── user_service.go
  ├── application/
  │   └── user/
  │       ├── dto/              // 数据传输对象
  │       │   └── user_dto.go
  │       └── service/          // 应用服务
  │           └── user_app_service.go
  ├── infrastructure/
  │   ├── persistence/
  │   │   └── sql/              // SQL实现
  │   │       └── user_repository.go
  │   ├── eventbus/             // 事件总线
  │   │   └── event_publisher.go
  │   └── transaction/          // 事务管理
  │       └── unit_of_work.go
  └── interfaces/
      └── http/
          └── handler/          // HTTP处理器
              └── user_handler.go
```

## 7. 重构收益

1. **提高代码的可读性和可维护性**：通过领域驱动设计将业务逻辑清晰地组织在各个层次，使代码更易于理解和维护。

2. **增强业务逻辑的表达能力**：使用充血模型（Entities和Value Objects）更准确地表达业务概念和规则。

3. **提高系统的可测试性**：领域对象和应用服务都有明确的接口，便于单元测试和集成测试。

4. **增强系统的可扩展性**：通过领域事件实现系统的解耦，便于添加新功能和集成其他服务。

5. **规范代码结构**：使团队成员能够快速理解系统架构，提高开发效率。 