# 商品服务DDD设计文档

## 1. 概述

本文档描述了基于领域驱动设计(DDD)的商品服务重构设计。商品服务是WZ系统的核心服务之一，负责商品信息管理、库存管理和商品状态管理等功能。

## 2. 领域模型

### 2.1 值对象 (Value Objects)

值对象是没有唯一标识的不可变对象，用于表示领域中的概念。

- **ProductID**: 商品ID值对象，使用int64类型表示
- **ProductName**: 商品名称值对象，包含名称验证逻辑
- **ProductDescription**: 商品描述值对象，包含长度验证逻辑
- **ProductSKU**: 商品SKU值对象，包含SKU格式验证逻辑
- **Price**: 价格值对象，以分为单位，包含非负验证逻辑
- **ProductStatus**: 商品状态值对象，定义了商品的不同状态（草稿、上架、下架、售罄、已删除）

### 2.2 实体 (Entities)

实体是具有唯一标识的对象，用于表示领域中的核心概念。

- **Product**: 商品实体，代表系统中的商品，包含商品的基本信息和行为

### 2.3 领域事件 (Domain Events)

领域事件用于表示领域中发生的重要事件，有助于实现系统的解耦和扩展性。

- **ProductCreatedEvent**: 商品创建事件
- **ProductPublishedEvent**: 商品发布事件
- **ProductUnpublishedEvent**: 商品下架事件
- **ProductDeletedEvent**: 商品删除事件
- **ProductStockChangedEvent**: 商品库存变更事件

### 2.4 仓储接口 (Repository Interfaces)

仓储接口定义了领域对象的持久化操作。

- **ProductRepository**: 商品仓储接口，定义了商品实体的CRUD操作
- **EventPublisher**: 事件发布接口，用于发布领域事件

### 2.5 领域服务 (Domain Services)

领域服务封装了不适合放在实体或值对象中的领域逻辑。

- **ProductDomainService**: 商品领域服务，包含商品创建、发布、下架、库存管理等核心业务逻辑

## 3. 商品生命周期和状态管理

商品实体具有完整的生命周期管理，状态流转如下：

1. **草稿(Draft)**: 商品创建后的初始状态，此时商品信息可以自由编辑
2. **上架(Active)**: 商品发布后的状态，对外可见可购买（前提是有库存）
3. **下架(Inactive)**: 商品手动下架后的状态，对外不可见不可购买
4. **售罄(SoldOut)**: 商品库存为0时自动转换的状态，对外可见但不可购买
5. **删除(Deleted)**: 商品删除后的状态，对外不可见不可购买，且不可恢复正常状态

## 4. 库存管理

商品实体内部封装了库存管理的业务规则：

1. 库存减少时，自动检查是否有足够库存
2. 库存为0时，自动将状态更新为售罄
3. 售罄状态下，当库存增加时，自动恢复为上架状态
4. 库存变更时，自动发布库存变更事件，便于其他模块感知变化

## 5. 业务规则

1. 商品SKU全局唯一
2. 库存不能为负数
3. 价格不能为负数
4. 库存为0的商品不能上架
5. 已删除的商品不能更改状态
6. 商品发布前必须有库存

## 6. 目录结构

```
internal/
  ├── domain/
  │   └── product/
  │       ├── entity/           // 实体
  │       │   └── product.go
  │       ├── valueobject/      // 值对象
  │       │   ├── product_id.go
  │       │   └── product_valueobjects.go
  │       ├── event/            // 领域事件
  │       │   └── product_events.go
  │       ├── repository/       // 仓储接口
  │       │   └── product_repository.go
  │       └── service/          // 领域服务
  │           └── product_service.go
  ├── application/
  │   └── product/
  │       ├── dto/              // 数据传输对象
  │       │   └── product_dto.go
  │       └── service/          // 应用服务
  │           └── product_app_service.go
  ├── infrastructure/
  │   └── persistence/
  │       └── product/          // 仓储实现
  │           └── product_repository.go
  └── interfaces/
      └── http/
          └── product/          // HTTP处理器
              └── product_handler.go
```

## 7. 重构收益

1. **提高代码的可读性和可维护性**：通过领域驱动设计将业务逻辑清晰地组织在各个层次，使代码更易于理解和维护。

2. **增强业务逻辑的表达能力**：使用充血模型（实体和值对象）更准确地表达商品领域的业务概念和规则。

3. **提高系统的可测试性**：领域对象和应用服务都有明确的接口，便于单元测试和集成测试。

4. **增强系统的可扩展性**：通过领域事件实现系统的解耦，便于添加新功能和集成其他服务。

5. **规范库存管理**：在领域模型中明确定义库存管理规则，避免库存管理混乱导致的业务问题。

6. **状态管理更清晰**：商品状态转换逻辑被封装在实体内部，更符合实际业务规则，减少状态混乱的可能性。 