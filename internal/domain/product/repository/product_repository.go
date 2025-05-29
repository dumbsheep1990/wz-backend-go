package repository

import (
	"wz-backend-go/internal/domain/product/entity"
	"wz-backend-go/internal/domain/product/valueobject"
)

// ProductRepository 商品仓储接口
type ProductRepository interface {
	// 保存商品
	Save(product *entity.Product) error

	// 根据ID查找商品
	FindByID(id valueobject.ProductID) (*entity.Product, error)

	// 根据SKU查找商品
	FindBySKU(sku valueobject.ProductSKU) (*entity.Product, error)

	// 查找分类下的商品
	FindByCategoryID(categoryID int64, page, pageSize int) ([]*entity.Product, int64, error)

	// 分页查询商品列表
	FindAll(page, pageSize int) ([]*entity.Product, int64, error)

	// 查询创建者的商品
	FindByCreatorID(creatorID int64, page, pageSize int) ([]*entity.Product, int64, error)

	// 搜索商品
	Search(keyword string, page, pageSize int) ([]*entity.Product, int64, error)

	// 查询上架商品
	FindActiveProducts(page, pageSize int) ([]*entity.Product, int64, error)
}

// EventPublisher 事件发布接口
type EventPublisher interface {
	// 发布事件
	Publish(event interface{}) error
}
