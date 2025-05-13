package model

import (
	"time"
)

// Order 订单模型
type Order struct {
	ID             int64       `json:"id" gorm:"primaryKey;column:id"`
	OrderNo        string      `json:"order_no" gorm:"column:order_no;type:varchar(64);not null;uniqueIndex"`
	UserID         int64       `json:"user_id" gorm:"column:user_id;not null;index"`
	TotalAmount    float64     `json:"total_amount" gorm:"column:total_amount;not null;type:decimal(10,2)"`
	PayAmount      float64     `json:"pay_amount" gorm:"column:pay_amount;not null;type:decimal(10,2)"`
	DiscountAmount float64     `json:"discount_amount" gorm:"column:discount_amount;default:0;type:decimal(10,2)"`
	Status         int         `json:"status" gorm:"column:status;not null;default:0"` // 0待支付,1已支付,2已发货,3已完成,4已取消,5已退款
	PayType        *int        `json:"pay_type" gorm:"column:pay_type"`                // 1支付宝,2微信
	PayTime        *time.Time  `json:"pay_time" gorm:"column:pay_time"`
	Consignee      string      `json:"consignee" gorm:"column:consignee;type:varchar(64)"`
	Phone          string      `json:"phone" gorm:"column:phone;type:varchar(20)"`
	Address        string      `json:"address" gorm:"column:address;type:varchar(255)"`
	Remark         string      `json:"remark" gorm:"column:remark;type:varchar(255)"`
	CreatedAt      time.Time   `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time   `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Items          []OrderItem `json:"items" gorm:"-"` // 订单项，不存储在数据库中
}

// TableName 返回表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	ID           int64     `json:"id" gorm:"primaryKey;column:id"`
	OrderID      int64     `json:"order_id" gorm:"column:order_id;not null;index"`
	ProductID    int64     `json:"product_id" gorm:"column:product_id;not null"`
	ProductName  string    `json:"product_name" gorm:"column:product_name;type:varchar(128);not null"`
	ProductImage string    `json:"product_image" gorm:"column:product_image;type:varchar(255)"`
	Price        float64   `json:"price" gorm:"column:price;not null;type:decimal(10,2)"`
	Quantity     int       `json:"quantity" gorm:"column:quantity;not null"`
	Subtotal     float64   `json:"subtotal" gorm:"column:subtotal;not null;type:decimal(10,2)"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 返回表名
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderQuery 订单查询条件
type OrderQuery struct {
	OrderNo   string `json:"order_no"`
	UserID    int64  `json:"user_id"`
	Status    *int   `json:"status"`
	PayType   *int   `json:"pay_type"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page" binding:"required,min=1"`
	PageSize  int    `json:"page_size" binding:"required,min=1,max=100"`
}

// OrderStatistics 订单统计
type OrderStatistics struct {
	TotalOrders     int     `json:"total_orders"`
	TotalAmount     float64 `json:"total_amount"`
	PaidOrders      int     `json:"paid_orders"`
	UnpaidOrders    int     `json:"unpaid_orders"`
	ShippedOrders   int     `json:"shipped_orders"`
	CompletedOrders int     `json:"completed_orders"`
	CanceledOrders  int     `json:"canceled_orders"`
	RefundedOrders  int     `json:"refunded_orders"`
}

// StatusCount 状态统计
type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// PaymentTypeCount 支付方式统计
type PaymentTypeCount struct {
	PayType string `json:"pay_type"`
	Count   int    `json:"count"`
}

// TrendData 趋势数据
type TrendData struct {
	TimeUnit string  `json:"time_unit"`
	Count    int     `json:"count"`
	Amount   float64 `json:"amount"`
}
