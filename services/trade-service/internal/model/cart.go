package model

import (
	"time"
)

// Cart 购物车模型
type Cart struct {
	ID        int64      `json:"id" db:"id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

// CartItem 购物车项模型
type CartItem struct {
	ID        int64     `json:"id" db:"id"`
	CartID    int64     `json:"cart_id" db:"cart_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int32     `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// 产品信息，不存储在数据库中
	ProductName string `json:"product_name" db:"-"`
	ImageURL    string `json:"image_url" db:"-"`
}

// CartResult 购物车查询结果
type CartResult struct {
	Items         []CartItem `json:"items"`
	TotalAmount   float64    `json:"total_amount"`
	TotalQuantity int32      `json:"total_quantity"`
}

// CalculateTotal 计算购物车总价和总数量
func (c *Cart) CalculateTotal() (float64, int32) {
	var totalAmount float64
	var totalQuantity int32

	for _, item := range c.Items {
		totalAmount += item.Price * float64(item.Quantity)
		totalQuantity += item.Quantity
	}

	return totalAmount, totalQuantity
}
