package service

import (
	"context"

	"wz-backend-go/trade-service/core/model"
)

// CartService 购物车服务接口
type CartService interface {
	// 添加商品到购物车
	AddToCart(ctx context.Context, userID int64, item model.CartItemRequest) error

	// 更新购物车商品数量
	UpdateQuantity(ctx context.Context, userID int64, productID int64, quantity int) error

	// 删除购物车商品
	RemoveFromCart(ctx context.Context, userID int64, productID int64) error

	// 获取用户购物车
	GetUserCart(ctx context.Context, userID int64) (*model.CartResponse, error)

	// 清空购物车
	ClearCart(ctx context.Context, userID int64) error

	// 更新选中状态
	UpdateSelected(ctx context.Context, userID int64, productID int64, selected bool) error

	// 全选/全不选
	SelectAll(ctx context.Context, userID int64, selected bool) error

	// 获取选中的购物车项
	GetSelectedItems(ctx context.Context, userID int64) ([]*model.Cart, error)

	// 获取购物车商品数量
	GetCartItemCount(ctx context.Context, userID int64) (int, error)

	// 获取购物车总金额
	GetCartTotal(ctx context.Context, userID int64) (float64, error)

	// 购物车结算，创建订单
	Checkout(ctx context.Context, userID int64, address, consignee, phone, remark string) (*model.Order, error)
}
