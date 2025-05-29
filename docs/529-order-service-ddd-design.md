# 订单服务DDD设计文档

## 1. 概述

本文档描述了基于领域驱动设计(DDD)的订单服务设计。订单服务是WZ系统的核心服务之一，负责订单创建、支付、发货、取消和退款等业务流程管理。

## 2. 领域模型

### 2.1 值对象 (Value Objects)

- **OrderID**: 订单ID值对象，使用string类型表示
- **OrderNumber**: 订单编号值对象，具有特定的格式和唯一性
- **OrderStatus**: 订单状态值对象，包含订单的不同状态
- **Money**: 金额值对象，包含金额和货币类型
- **Address**: 地址值对象，包含收货地址信息
- **PaymentMethod**: 支付方式值对象
- **ShippingMethod**: 配送方式值对象
- **OrderItemID**: 订单项ID值对象

### 2.2 实体 (Entities)

- **Order**: 订单实体，代表一个客户的订单
- **OrderItem**: 订单项实体，代表订单中的一个商品项
- **OrderDiscount**: 订单折扣实体，代表应用于订单的折扣

### 2.3 领域事件 (Domain Events)

- **OrderCreatedEvent**: 订单创建事件
- **OrderPaidEvent**: 订单支付事件
- **OrderShippedEvent**: 订单发货事件
- **OrderDeliveredEvent**: 订单送达事件
- **OrderCancelledEvent**: 订单取消事件
- **OrderRefundedEvent**: 订单退款事件
- **OrderItemAddedEvent**: 订单项添加事件
- **OrderItemRemovedEvent**: 订单项移除事件

### 2.4 聚合 (Aggregates)

订单聚合是以Order实体为根的聚合，包含以下成员：

- Order (聚合根)
- OrderItem
- OrderDiscount

聚合确保订单数据的一致性和完整性，如订单总额必须等于所有订单项金额加上运费减去折扣。

### 2.5 仓储接口 (Repository Interfaces)

- **OrderRepository**: 订单仓储接口，定义了订单聚合的CRUD操作
- **EventPublisher**: 事件发布接口，用于发布领域事件

### 2.6 领域服务 (Domain Services)

- **OrderDomainService**: 订单领域服务，包含订单创建、支付、发货等核心业务逻辑
- **PricingService**: 定价服务，计算订单价格、折扣和税费
- **InventoryService**: 库存服务，处理库存检查和扣减

## 3. 订单生命周期和状态管理

订单实体具有完整的生命周期管理，状态流转如下：

1. **已创建(Created)**: 订单初始创建状态，购物车转为订单
2. **待支付(Pending)**: 订单已提交，等待支付
3. **已支付(Paid)**: 订单已完成支付
4. **已发货(Shipped)**: 订单已发货
5. **已送达(Delivered)**: 订单已送达
6. **已完成(Completed)**: 订单流程完成
7. **已取消(Cancelled)**: 订单被取消
8. **退款中(Refunding)**: 订单正在处理退款
9. **已退款(Refunded)**: 订单已完成退款

订单状态转换必须遵循特定规则，例如：
- 只有在"待支付"状态的订单才能被支付
- 只有在"已支付"状态的订单才能发货
- 已取消或已完成的订单不能再修改

## 4. 业务规则

1. 订单创建时必须检查商品库存是否充足
2. 订单总额必须等于所有订单项金额之和加上运费减去折扣
3. 订单支付后必须生成相应的支付记录
4. 订单取消时，如果已支付，必须进行退款
5. 订单支付后，必须扣减相应的商品库存
6. 订单号必须全局唯一且符合特定格式
7. 已发货的订单不能取消，只能申请退款

## 5. 应用服务设计

应用服务层将提供以下功能：

1. 创建订单
2. 添加/移除订单项
3. 应用折扣
4. 提交订单
5. 支付订单
6. 取消订单
7. 申请退款
8. 确认发货
9. 确认收货
10. 查询订单

## 6. 目录结构

```
internal/
  ├── domain/
  │   └── order/
  │       ├── entity/           // 实体
  │       │   ├── order.go
  │       │   ├── order_item.go
  │       │   └── order_discount.go
  │       ├── valueobject/      // 值对象
  │       │   ├── order_id.go
  │       │   ├── order_number.go
  │       │   ├── order_status.go
  │       │   ├── money.go
  │       │   ├── address.go
  │       │   └── payment_method.go
  │       ├── event/            // 领域事件
  │       │   └── order_events.go
  │       ├── repository/       // 仓储接口
  │       │   └── order_repository.go
  │       └── service/          // 领域服务
  │           ├── order_domain_service.go
  │           └── pricing_service.go
  ├── application/
  │   └── order/
  │       ├── dto/              // 数据传输对象
  │       │   └── order_dto.go
  │       └── service/          // 应用服务
  │           └── order_app_service.go
  ├── infrastructure/
  │   └── persistence/
  │       └── order/            // 仓储实现
  │           └── order_repository.go
  └── interfaces/
      └── http/
          └── order/            // HTTP处理器
              └── order_handler.go
```

## 7. 重构收益

1. **提高业务模型的准确性**：通过领域驱动设计，订单模型更准确地反映业务规则和流程。

2. **提高代码的可维护性**：清晰的领域边界和业务规则封装，使代码更易于理解和维护。

3. **增强业务逻辑的封装**：订单状态转换、价格计算等关键业务逻辑封装在领域层，避免逻辑泄漏。

4. **提高系统的可扩展性**：通过领域事件实现系统的解耦，便于添加新功能和与其他服务集成。

5. **减少业务错误**：状态机确保订单状态转换的正确性，减少状态混乱导致的业务错误。

6. **提高测试覆盖率**：领域模型和业务规则可以独立测试，提高测试覆盖率和系统稳定性。

7. **支持复杂业务场景**：如部分发货、部分退款等复杂场景可以更容易实现。

8. **提高团队协作效率**：统一的业务语言和模型，减少沟通成本，提高团队协作效率。

## 8. 关键技术点

1. **订单状态机实现**：使用状态模式处理订单状态转换
2. **领域事件发布与订阅**：基于事件的异步处理机制
3. **聚合边界设计**：保证数据一致性的聚合设计
4. **值对象不可变性**：确保值对象的不可变特性
5. **仓储模式实现**：处理复杂聚合的持久化
6. **领域服务与应用服务分离**：明确责任边界 