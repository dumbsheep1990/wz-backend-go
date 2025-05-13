package model

import (
	"time"
)

// Cart 购物车模型
type Cart struct {
	ID           int64     `json:"id" gorm:"primaryKey;column:id"`
	UserID       int64     `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_product"`
	ProductID    int64     `json:"product_id" gorm:"column:product_id;not null;uniqueIndex:idx_user_product"`
	ProductName  string    `json:"product_name" gorm:"column:product_name;type:varchar(128);not null"`
	ProductImage string    `json:"product_image" gorm:"column:product_image;type:varchar(255)"`
	Price        float64   `json:"price" gorm:"column:price;not null;type:decimal(10,2)"`
	Quantity     int       `json:"quantity" gorm:"column:quantity;not null"`
	Selected     bool      `json:"selected" gorm:"column:selected;not null;default:1"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 返回表名
func (Cart) TableName() string {
	return "carts"
}

// CartQuery 购物车查询条件
type CartQuery struct {
	UserID   int64 `json:"user_id" binding:"required"`
	Selected *bool `json:"selected"`
}

// CartItemRequest 购物车添加/更新请求
type CartItemRequest struct {
	ProductID    int64   `json:"product_id" binding:"required"`
	ProductName  string  `json:"product_name" binding:"required"`
	ProductImage string  `json:"product_image"`
	Price        float64 `json:"price" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required,min=1"`
	Selected     bool    `json:"selected"`
}

// CartResponse 购物车响应
type CartResponse struct {
	CartItems []Cart  `json:"cart_items"`
	Total     float64 `json:"total"`
	Count     int     `json:"count"`
}
