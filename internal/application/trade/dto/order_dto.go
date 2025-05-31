package dto

import (
	"time"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CustomerID      int64                      `json:"customer_id" binding:"required"`
	Items           []*CreateOrderItemRequest  `json:"items" binding:"required,min=1"`
	ShippingAddress *AddressRequest            `json:"shipping_address" binding:"required"`
	BillingAddress  *AddressRequest            `json:"billing_address" binding:"required"`
	ShippingMethod  int32                      `json:"shipping_method" binding:"required"`
	ShippingFee     int64                      `json:"shipping_fee"`
	Tax             int64                      `json:"tax"`
	Currency        string                     `json:"currency" binding:"required"`
	Note            string                     `json:"note"`
}

// CreateOrderItemRequest 创建订单项请求
type CreateOrderItemRequest struct {
	ProductID   int64             `json:"product_id" binding:"required"`
	ProductName string            `json:"product_name" binding:"required"`
	ProductSKU  string            `json:"product_sku" binding:"required"`
	Quantity    int32             `json:"quantity" binding:"required,min=1"`
	UnitPrice   int64             `json:"unit_price" binding:"required,min=0"`
	Attributes  map[string]string `json:"attributes"`
}

// AddressRequest 地址请求
type AddressRequest struct {
	Country       string `json:"country" binding:"required"`
	Province      string `json:"province" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district"`
	DetailAddress string `json:"detail_address" binding:"required"`
	PostalCode    string `json:"postal_code"`
	ContactName   string `json:"contact_name" binding:"required"`
	ContactPhone  string `json:"contact_phone" binding:"required"`
}

// GetOrderRequest 获取订单请求
type GetOrderRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  int64  `json:"user_id" binding:"required"`
}

// GetOrderByNumberRequest 根据订单号获取订单请求
type GetOrderByNumberRequest struct {
	OrderNumber string `json:"order_number" binding:"required"`
	UserID      int64  `json:"user_id" binding:"required"`
}

// ListUserOrdersRequest 获取用户订单列表请求
type ListUserOrdersRequest struct {
	CustomerID int64 `json:"customer_id" binding:"required"`
	Page       int   `json:"page" binding:"min=1"`
	PageSize   int   `json:"page_size" binding:"min=1,max=100"`
}

// PayOrderRequest 支付订单请求
type PayOrderRequest struct {
	OrderID       string `json:"order_id" binding:"required"`
	UserID        int64  `json:"user_id" binding:"required"`
	PaymentMethod int32  `json:"payment_method" binding:"required"`
}

// ShipOrderRequest 发货请求
type ShipOrderRequest struct {
	OrderID        string `json:"order_id" binding:"required"`
	TrackingNumber string `json:"tracking_number" binding:"required"`
}

// DeliverOrderRequest 确认送达请求
type DeliverOrderRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

// CompleteOrderRequest 完成订单请求
type CompleteOrderRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  int64  `json:"user_id" binding:"required"`
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  int64  `json:"user_id" binding:"required"`
	Reason  string `json:"reason"`
}

// RefundOrderRequest 退款请求
type RefundOrderRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  int64  `json:"user_id" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

// SearchOrdersRequest 搜索订单请求
type SearchOrdersRequest struct {
	Keyword  string `json:"keyword" binding:"required"`
	Page     int    `json:"page" binding:"min=1"`
	PageSize int    `json:"page_size" binding:"min=1,max=100"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID              string                    `json:"id"`
	OrderNumber     string                    `json:"order_number"`
	CustomerID      int64                     `json:"customer_id"`
	Status          int32                     `json:"status"`
	Items           []*OrderItemResponse      `json:"items"`
	Discounts       []*OrderDiscountResponse  `json:"discounts"`
	ShippingAddress *AddressResponse          `json:"shipping_address"`
	BillingAddress  *AddressResponse          `json:"billing_address"`
	PaymentMethod   int32                     `json:"payment_method"`
	ShippingMethod  int32                     `json:"shipping_method"`
	Subtotal        int64                     `json:"subtotal"`
	ShippingFee     int64                     `json:"shipping_fee"`
	Tax             int64                     `json:"tax"`
	DiscountAmount  int64                     `json:"discount_amount"`
	TotalAmount     int64                     `json:"total_amount"`
	Currency        string                    `json:"currency"`
	Note            string                    `json:"note"`
	TrackingNumber  string                    `json:"tracking_number"`
	PaidAt          *time.Time                `json:"paid_at"`
	ShippedAt       *time.Time                `json:"shipped_at"`
	DeliveredAt     *time.Time                `json:"delivered_at"`
	CompletedAt     *time.Time                `json:"completed_at"`
	CancelledAt     *time.Time                `json:"cancelled_at"`
	RefundedAt      *time.Time                `json:"refunded_at"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// OrderItemResponse 订单项响应
type OrderItemResponse struct {
	ID          string            `json:"id"`
	ProductID   int64             `json:"product_id"`
	ProductName string            `json:"product_name"`
	ProductSKU  string            `json:"product_sku"`
	Quantity    int32             `json:"quantity"`
	UnitPrice   int64             `json:"unit_price"`
	TotalPrice  int64             `json:"total_price"`
	Attributes  map[string]string `json:"attributes"`
}

// OrderDiscountResponse 订单折扣响应
type OrderDiscountResponse struct {
	ID             string `json:"id"`
	Type           int32  `json:"type"`
	Name           string `json:"name"`
	DiscountAmount int64  `json:"discount_amount"`
	Description    string `json:"description"`
}

// AddressResponse 地址响应
type AddressResponse struct {
	Country       string `json:"country"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detail_address"`
	PostalCode    string `json:"postal_code"`
	ContactName   string `json:"contact_name"`
	ContactPhone  string `json:"contact_phone"`
}

// ListOrdersResponse 订单列表响应
type ListOrdersResponse struct {
	Orders   []*OrderResponse `json:"orders"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// OrderStatisticsResponse 订单统计响应
type OrderStatisticsResponse struct {
	TotalOrders       int64 `json:"total_orders"`
	PendingOrders     int64 `json:"pending_orders"`
	PaidOrders        int64 `json:"paid_orders"`
	ShippedOrders     int64 `json:"shipped_orders"`
	CompletedOrders   int64 `json:"completed_orders"`
	CancelledOrders   int64 `json:"cancelled_orders"`
	RefundedOrders    int64 `json:"refunded_orders"`
	TotalRevenue      int64 `json:"total_revenue"`
	TodayOrders       int64 `json:"today_orders"`
	TodayRevenue      int64 `json:"today_revenue"`
}

// OrderFilterRequest 订单过滤请求
type OrderFilterRequest struct {
	Status        *int32     `json:"status"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	MinAmount     *int64     `json:"min_amount"`
	MaxAmount     *int64     `json:"max_amount"`
	PaymentMethod *int32     `json:"payment_method"`
	Page          int        `json:"page" binding:"min=1"`
	PageSize      int        `json:"page_size" binding:"min=1,max=100"`
} 