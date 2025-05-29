# 用户积分服务DDD重构记录

## 总体目标

将用户积分服务从传统的三层架构（service/handler/delivery）重构为领域驱动设计（DDD）架构（domain/application/infrastructure/interfaces），实现贫血模型到充血模型的转变。

## 重构步骤记录

### 1. 领域层（Domain Layer）实现 ✓

#### 1.1 值对象（Value Objects）✓

在 `internal/domain/user/valueobject/user_points_valueobjects.go` 中实现：

- ID - 唯一标识符值对象
- UserID - 用户标识符值对象
- Points - 积分值对象
- PointsType - 积分类型值对象（增加/减少）
- Source - 积分来源值对象（签到、购买、活动等）
- Description - 描述值对象
- RelatedType - 关联类型值对象
- TenantID - 租户ID值对象

这些值对象是不可变的，包含自己的验证逻辑，确保业务规则在领域层得到正确实施。

#### 1.2 实体（Entities）✓

##### 1.2.1 UserPoints 实体 ✓

在 `internal/domain/user/entity/user_points.go` 中实现，代表一条积分记录。包含：

- 属性: id, userID, points, totalPoints, pointsType, source 等
- 构造函数: NewUserPoints()
- 业务方法: Revoke() - 撤销积分记录
- 领域事件: 在实体创建和撤销时发布

##### 1.2.2 PointsRules 实体 ✓

在 `internal/domain/user/entity/points_rules.go` 中实现，代表积分规则。包含：

- 属性: signInPoints, commentPoints, sharePoints 等规则配置
- 构造函数: NewPointsRules()
- 业务方法: Update() - 更新规则
- 领域事件: 在规则创建和更新时发布

#### 1.3 领域事件（Domain Events）✓

在 `internal/domain/shared/event/domain_event.go` 中定义基础事件接口：

```go
type DomainEvent interface {
    EventID() string
    AggregateID() string
    EventType() string
    OccurredTime() time.Time
}
```

在实体文件中定义具体的领域事件：

- UserPointsCreatedEvent - 用户积分创建事件
- UserPointsRevokedEvent - 用户积分撤销事件
- PointsRulesCreatedEvent - 积分规则创建事件
- PointsRulesUpdatedEvent - 积分规则更新事件

#### 1.4 仓储接口（Repository Interfaces）✓

在 `internal/domain/user/repository/user_points_repository.go` 中定义：

- UserPointsRepository - 定义用户积分仓储接口
- PointsRulesRepository - 定义积分规则仓储接口

这些接口定义了领域对象的持久化方法，但不涉及具体实现，遵循依赖倒置原则。

### 2. 应用层（Application Layer）实现 ✓

#### 2.1 DTO（Data Transfer Objects）✓

在 `internal/application/user/dto/user_points_dto.go` 中定义：

- CreatePointsRequest - 创建积分请求DTO
- PointsDTO - 积分记录DTO
- ListPointsRequest/Response - 查询积分列表DTO
- PointsStatisticsResponse - 积分统计DTO
- PointsRulesRequest/Response - 积分规则DTO

这些DTO用于在应用层和接口层之间传输数据，避免领域对象暴露给外部。

#### 2.2 应用服务（Application Services）✓

在 `internal/application/user/service/user_points_application_service.go` 中实现：

- UserPointsApplicationService - 用户积分应用服务，协调领域对象和基础设施

包括以下业务用例：
- CreatePoints() - 创建积分记录
- GetPointByID() - 获取积分记录详情
- ListPointsByUserID() - 获取用户积分记录列表
- RevokePoint() - 撤销积分记录
- GetPointsStatistics() - 获取积分统计
- GetPointsRules() - 获取积分规则
- UpdatePointsRules() - 更新积分规则

### 3. 基础设施层（Infrastructure Layer）实现 ✓

#### 3.1 仓储实现 ✓

##### 3.1.1 用户积分仓储实现 ✓

在 `internal/infrastructure/persistence/sql/user_points_repository_impl.go` 中实现：

- SQLUserPointsRepository - 使用SQL实现的用户积分仓储

##### 3.1.2 积分规则仓储实现 ✓

在 `internal/infrastructure/persistence/sql/points_rules_repository_impl.go` 中实现：

- SQLPointsRulesRepository - 使用SQL实现的积分规则仓储

#### 3.2 事件总线实现 ✓

在 `internal/infrastructure/event/event_bus.go` 中实现：

- SimpleEventBus - 简单的事件总线实现，支持发布和订阅领域事件

#### 3.3 工作单元实现 ✓

在 `internal/infrastructure/persistence/database/unit_of_work_impl.go` 中实现：

- SQLUnitOfWorkImpl - 使用SQL实现的工作单元，支持事务管理

### 4. 接口层（Interfaces Layer）实现 ✓

#### 4.1 HTTP控制器 ✓

在 `internal/interfaces/http/controller/user_points_controller.go` 中实现：

- UserPointsController - 用户积分HTTP控制器，处理Web请求

路由包括：
- GET /points/users/:userId - 获取用户积分记录
- POST /points - 创建积分记录
- GET /points/:id - 获取积分记录详情
- PUT /points/:id/revoke - 撤销积分记录
- GET /points/statistics - 获取积分统计
- GET /points/rules - 获取积分规则
- PUT /points/rules - 更新积分规则

#### 4.2 中间件 ✓

在 `internal/interfaces/http/middleware/auth.go` 中实现：

- AdminAuth() - 管理员认证中间件
- UserAuth() - 用户认证中间件

### 5. 应用引导 ✓

在 `internal/app/bootstrap/user_points_service.go` 中实现：

- InitUserPointsService() - 初始化用户积分服务的依赖注入函数

## 重构效果对比

### 重构前（贫血模型）

1. 业务逻辑主要在Service层，实体只是数据容器
2. 缺乏值对象概念，基本类型随处使用
3. 没有领域事件，状态变化需要显式调用其他服务
4. 缺乏充血模型中的业务行为封装

### 重构后（充血模型）

1. 业务规则和不变性被封装在领域层（值对象和实体）
2. 使用值对象表示业务概念，确保一致性
3. 通过领域事件实现松耦合的状态变化通知
4. 实体包含业务行为，确保数据和行为一起封装
5. 应用服务变得更薄，主要负责协调和事务
6. 更清晰的架构分层和职责划分

## 总结

本次重构成功将用户积分服务从传统三层架构转换为DDD架构，实现了从贫血模型到充血模型的转变。领域模型更加丰富，业务规则得到更好的封装，系统各层之间的耦合度降低。重构后的代码更易于理解、测试和维护，也更易于扩展新功能。

完成度：100% 