package repository

import (
	"context"

	"wz-backend-go/trade-service/core/model"
)

// CartRepository 购物车仓储接口
type CartRepository interface {
	// 添加购物车
	Add(ctx context.Context, cart *model.Cart) error

	// 更新购物车
	Update(ctx context.Context, cart *model.Cart) error

	// 删除购物车项
	Delete(ctx context.Context, id int64) error

	// 根据用户ID获取购物车
	GetByUserID(ctx context.Context, userID int64) ([]*model.Cart, error)

	// 清空用户购物车
	Clear(ctx context.Context, userID int64) error

	// 获取选中的购物车项
	GetSelected(ctx context.Context, userID int64) ([]*model.Cart, error)

	// 更新选中状态
	UpdateSelected(ctx context.Context, userID int64, productID int64, selected bool) error

	// 批量更新选中状态
	UpdateAllSelected(ctx context.Context, userID int64, selected bool) error

	// 根据用户ID和商品ID获取购物车项
	GetByUserIDAndProductID(ctx context.Context, userID int64, productID int64) (*model.Cart, error)

	// 计算购物车总金额
	GetCartTotal(ctx context.Context, userID int64) (float64, error)
}
