package repository

import (
	"wz-backend-go/internal/domain/order/entity"
	"wz-backend-go/internal/domain/shared/event"
	uservo "wz-backend-go/internal/domain/user/valueobject"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// Save 保存订单
	Save(order *entity.Order) error

	// FindByID 根据ID查找订单
	FindByID(id ordervo.OrderID) (*entity.Order, error)

	// FindByOrderNumber 根据订单号查找订单
	FindByOrderNumber(orderNumber ordervo.OrderNumber) (*entity.Order, error)

	// FindByCustomerID 查找客户的订单
	FindByCustomerID(customerID uservo.UserID, page, pageSize int) ([]*entity.Order, int64, error)

	// FindByStatus 根据状态查找订单
	FindByStatus(status ordervo.OrderStatus, page, pageSize int) ([]*entity.Order, int64, error)

	// FindAll 分页查询所有订单
	FindAll(page, pageSize int) ([]*entity.Order, int64, error)

	// FindActiveOrders 查询活跃订单（非终态）
	FindActiveOrders(customerID uservo.UserID, page, pageSize int) ([]*entity.Order, int64, error)

	// FindRecentOrders 查询最近的订单
	FindRecentOrders(limit int) ([]*entity.Order, error)

	// Search 搜索订单
	Search(keyword string, page, pageSize int) ([]*entity.Order, int64, error)

	// Delete 删除订单（逻辑删除）
	Delete(id ordervo.OrderID) error
}

// EventPublisher 事件发布接口
type EventPublisher interface {
	// Publish 发布事件
	Publish(event event.DomainEvent) error
}
