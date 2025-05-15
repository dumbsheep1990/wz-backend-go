package model

import (
	"time"
)

// 订单状态
const (
	OrderStatusPending    = "pending"    // 待支付
	OrderStatusPaid       = "paid"       // 已支付
	OrderStatusProcessing = "processing" // 处理中
	OrderStatusShipped    = "shipped"    // 已发货
	OrderStatusDelivered  = "delivered"  // 已送达
	OrderStatusCompleted  = "completed"  // 已完成
	OrderStatusCancelled  = "cancelled"  // 已取消
	OrderStatusRefunded   = "refunded"   // 已退款
)

// Order 订单模型
type Order struct {
	ID           int64       `json:"id" db:"id"`
	OrderNumber  string      `json:"order_number" db:"order_number"`
	UserID       int64       `json:"user_id" db:"user_id"`
	TotalAmount  float64     `json:"total_amount" db:"total_amount"`
	Status       string      `json:"status" db:"status"`
	Address      string      `json:"address" db:"address"`
	ContactName  string      `json:"contact_name" db:"contact_name"`
	ContactPhone string      `json:"contact_phone" db:"contact_phone"`
	Note         string      `json:"note" db:"note"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	Items        []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Payments     []Payment   `json:"payments" gorm:"foreignKey:OrderID"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID         int64     `json:"id" db:"id"`
	OrderID    int64     `json:"order_id" db:"order_id"`
	ProductID  int64     `json:"product_id" db:"product_id"`
	Quantity   int32     `json:"quantity" db:"quantity"`
	Price      float64   `json:"price" db:"price"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// OrderHistory 订单历史记录
type OrderHistory struct {
	ID        int64     `json:"id" db:"id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	Status    string    `json:"status" db:"status"`
	Remark    string    `json:"remark" db:"remark"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// GenerateOrderNumber 生成订单号
func GenerateOrderNumber() string {
	now := time.Now()
	orderNumber := now.Format("20060102150405") + RandomNumeric(6)
	return orderNumber
}

// RandomNumeric 生成指定长度的随机数字字符串
func RandomNumeric(length int) string {
	const numeric = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = numeric[time.Now().UnixNano()%int64(len(numeric))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(result)
}
