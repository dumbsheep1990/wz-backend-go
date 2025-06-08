# WanZhi平台 - DDD微服务开发指南

## 一、前言

本文档旨在规范化万智平台基于领域驱动设计（Domain-Driven Design，简称DDD）的微服务开发流程，为团队成员提供清晰的开发指导和架构规范。通过遵循本指南，可确保各微服务实现结构一致、代码风格统一、架构设计合理，便于维护和扩展。

## 二、DDD架构概述

### 1. 分层架构

万智平台采用经典四层DDD架构，各层职责明确分离：

```
wz-backend-go/
├── cmd/                          # 应用程序入口
│   ├── <microservice>/           # 各微服务的入口程序
│   │   └── main.go               # 主函数
├── internal/                     # 内部包
│   ├── domain/                   # 领域层
│   │   ├── <microservice>/       # 微服务领域模型
│   │   │   ├── entity/           # 领域实体
│   │   │   ├── repository/       # 仓储接口
│   │   │   ├── service/          # 领域服务
│   │   │   ├── valueobject/      # 值对象
│   │   │   └── dto/              # 领域层数据传输对象
│   ├── application/              # 应用层
│   │   ├── <microservice>/       # 微服务应用服务
│   │   │   ├── <service>.go      # 应用服务实现
│   │   │   └── dto/              # 应用层数据传输对象
│   ├── infrastructure/           # 基础设施层
│   │   ├── repository/           # 仓储实现
│   │   │   ├── mysql/            # MySQL实现
│   │   │   ├── mongodb/          # MongoDB实现
│   │   │   └── redis/            # Redis实现
│   │   ├── client/               # 微服务客户端
│   │   └── middleware/           # 中间件
│   └── delivery/                 # 交付层
│       ├── http/                 # HTTP交付
│       │   └── internal/
│       │       ├── handler/      # HTTP处理器
│       │       │   └── <microservice>/
│       │       └── router/       # 路由配置
│       └── rpc/                  # RPC交付
└── pkg/                          # 公共包
```

### 2. 各层职责

**领域层（Domain Layer）**：
- 定义业务核心概念和规则
- 包含实体、值对象、聚合根、领域事件和领域服务
- 保持独立性，不依赖其他层
- 领域逻辑和业务规则的唯一可信来源

**应用层（Application Layer）**：
- 协调领域对象完成特定用例
- 编排业务流程
- 进行事务管理
- 将领域对象转换为DTO
- 不包含业务规则，只有业务流程

**基础设施层（Infrastructure Layer）**：
- 实现领域层定义的接口
- 提供技术能力支持
- 处理数据持久化
- 管理消息、缓存等基础设施
- 实现与外部系统的集成

**交付层（Delivery Layer）**：
- 处理用户请求和响应
- 进行参数验证和格式转换
- 调用应用服务
- 定义API接口规范
- 处理权限和认证
