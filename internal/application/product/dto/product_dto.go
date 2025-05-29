package dto

import (
	"time"

	"wz-backend-go/internal/domain/product/entity"
)

// ProductDTO 商品数据传输对象
type ProductDTO struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	SKU          string    `json:"sku"`
	Price        int64     `json:"price"`
	PriceDisplay string    `json:"priceDisplay"` // 价格显示字符串，例如 "99.99"
	Stock        int32     `json:"stock"`
	Status       int32     `json:"status"`
	StatusDesc   string    `json:"statusDesc"`
	CategoryID   int64     `json:"categoryId"`
	CreatorID    int64     `json:"creatorId"`
	IsActive     bool      `json:"isActive"`
	IsAvailable  bool      `json:"isAvailable"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ProductCreateRequest 创建商品请求
type ProductCreateRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=2000"`
	SKU         string `json:"sku" binding:"required,min=6,max=20"`
	Price       int64  `json:"price" binding:"min=0"`
	Stock       int32  `json:"stock" binding:"min=0"`
	CategoryID  int64  `json:"categoryId" binding:"required,min=1"`
}

// ProductUpdateRequest 更新商品请求
type ProductUpdateRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=2000"`
	Price       int64  `json:"price" binding:"min=0"`
	CategoryID  int64  `json:"categoryId" binding:"required,min=1"`
}

// StockUpdateRequest 更新库存请求
type StockUpdateRequest struct {
	Stock  int32  `json:"stock" binding:"min=0"`
	Reason string `json:"reason" binding:"required"`
}

// StockDecrementRequest 减少库存请求
type StockDecrementRequest struct {
	Quantity int32  `json:"quantity" binding:"required,min=1"`
	Reason   string `json:"reason" binding:"required"`
}

// StockIncrementRequest 增加库存请求
type StockIncrementRequest struct {
	Quantity int32  `json:"quantity" binding:"required,min=1"`
	Reason   string `json:"reason" binding:"required"`
}

// ProductResponse 商品响应
type ProductResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *ProductDTO `json:"data,omitempty"`
}

// ProductsResponse 商品列表响应
type ProductsResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *ProductsData `json:"data,omitempty"`
}

// ProductsData 商品列表数据
type ProductsData struct {
	Total int64         `json:"total"`
	List  []*ProductDTO `json:"list"`
}

// 将领域实体转换为DTO
func ToProductDTO(product *entity.Product) *ProductDTO {
	return &ProductDTO{
		ID:           product.ID().Value(),
		Name:         product.Name().Value(),
		Description:  product.Description().Value(),
		SKU:          product.SKU().Value(),
		Price:        product.Price().Value(),
		PriceDisplay: product.Price().String(),
		Stock:        product.Stock(),
		Status:       product.Status().Value(),
		StatusDesc:   product.Status().String(),
		CategoryID:   product.CategoryID(),
		CreatorID:    product.CreatorID().Value(),
		IsActive:     product.IsActive(),
		IsAvailable:  product.IsAvailable(),
		CreatedAt:    product.CreatedAt(),
		UpdatedAt:    product.UpdatedAt(),
	}
}

// 将商品DTO列表转换为响应
func ToProductsResponse(products []*entity.Product, total int64) *ProductsResponse {
	productDTOs := make([]*ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = ToProductDTO(product)
	}

	return &ProductsResponse{
		Code:    0,
		Message: "success",
		Data: &ProductsData{
			Total: total,
			List:  productDTOs,
		},
	}
}
