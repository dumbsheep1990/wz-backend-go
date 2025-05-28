package repository

import (
	"context"

	"github.com/yourusername/wz-backend-go/internal/domain/trade/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// 保存订单（创建或更新）
	Save(ctx context.Context, order *entity.Order) error
	
	// 根据ID查找订单
	FindByID(ctx context.Context, id valueobject.OrderID) (*entity.Order, error)
	
	// 根据用户ID查找订单列表
	FindByUserID(ctx context.Context, userID valueobject.UserID, page, pageSize int) ([]*entity.Order, int64, error)
	
	// 根据状态查找订单列表
	FindByStatus(ctx context.Context, status valueobject.OrderStatus, page, pageSize int) ([]*entity.Order, int64, error)
	
	// 根据用户ID和状态查找订单列表
	FindByUserIDAndStatus(ctx context.Context, userID valueobject.UserID, status valueobject.OrderStatus, page, pageSize int) ([]*entity.Order, int64, error)
	
	// 删除订单
	Delete(ctx context.Context, id valueobject.OrderID) error
}

// CartRepository 购物车仓储接口
type CartRepository interface {
	// 保存购物车（创建或更新）
	Save(ctx context.Context, cart *entity.Cart) error
	
	// 根据ID查找购物车
	FindByID(ctx context.Context, id valueobject.CartID) (*entity.Cart, error)
	
	// 根据用户ID查找购物车
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*entity.Cart, error)
	
	// 删除购物车
	Delete(ctx context.Context, id valueobject.CartID) error
}

// PaymentRepository 支付仓储接口
type PaymentRepository interface {
	// 保存支付（创建或更新）
	Save(ctx context.Context, payment *entity.Payment) error
	
	// 根据ID查找支付
	FindByID(ctx context.Context, id valueobject.PaymentID) (*entity.Payment, error)
	
	// 根据订单ID查找支付
	FindByOrderID(ctx context.Context, orderID valueobject.OrderID) (*entity.Payment, error)
	
	// 根据用户ID查找支付列表
	FindByUserID(ctx context.Context, userID valueobject.UserID, page, pageSize int) ([]*entity.Payment, int64, error)
	
	// 根据状态查找支付列表
	FindByStatus(ctx context.Context, status valueobject.PaymentStatus, page, pageSize int) ([]*entity.Payment, int64, error)
	
	// 删除支付
	Delete(ctx context.Context, id valueobject.PaymentID) error
}
