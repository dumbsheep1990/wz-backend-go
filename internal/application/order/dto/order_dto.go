package dto

import (
	"time"
	ordervo "wz-backend-go/internal/domain/order/valueobject"
)

// 命令DTO

// CreateOrderCommand 创建订单命令
type CreateOrderCommand struct {
	CustomerID      string                 `json:"customer_id" binding:"required"`
	Items           []OrderItemDTO         `json:"items" binding:"required,min=1"`
	ShippingAddress AddressDTO             `json:"shipping_address" binding:"required"`
	BillingAddress  AddressDTO             `json:"billing_address"`
	ShippingMethod  string                 `json:"shipping_method" binding:"required"`
	PaymentMethod   string                 `json:"payment_method"`
	Note            string                 `json:"note"`
	Attributes      map[string]interface{} `json:"attributes"`
}

// UpdateOrderCommand 更新订单命令
type UpdateOrderCommand struct {
	OrderID         string     `json:"order_id" binding:"required"`
	ShippingAddress *AddressDTO `json:"shipping_address"`
	BillingAddress  *AddressDTO `json:"billing_address"`
	ShippingMethod  string     `json:"shipping_method"`
	Note            string     `json:"note"`
}

// AddOrderItemCommand 添加订单项命令
type AddOrderItemCommand struct {
	OrderID      string            `json:"order_id" binding:"required"`
	ProductID    string            `json:"product_id" binding:"required"`
	ProductName  string            `json:"product_name" binding:"required"`
	ProductSKU   string            `json:"product_sku" binding:"required"`
	Quantity     int32             `json:"quantity" binding:"required,min=1"`
	UnitPrice    MoneyDTO          `json:"unit_price" binding:"required"`
	Attributes   map[string]string `json:"attributes"`
}

// RemoveOrderItemCommand 移除订单项命令
type RemoveOrderItemCommand struct {
	OrderID string `json:"order_id" binding:"required"`
	ItemID  string `json:"item_id" binding:"required"`
}

// UpdateOrderItemQuantityCommand 更新订单项数量命令
type UpdateOrderItemQuantityCommand struct {
	OrderID  string `json:"order_id" binding:"required"`
	ItemID   string `json:"item_id" binding:"required"`
	Quantity int32  `json:"quantity" binding:"required,min=1"`
}

// AddOrderDiscountCommand 添加订单折扣命令
type AddOrderDiscountCommand struct {
	OrderID      string   `json:"order_id" binding:"required"`
	DiscountType string   `json:"discount_type" binding:"required"`
	DiscountName string   `json:"discount_name" binding:"required"`
	Amount       MoneyDTO `json:"amount" binding:"required"`
}

// RemoveOrderDiscountCommand 移除订单折扣命令
type RemoveOrderDiscountCommand struct {
	OrderID    string `json:"order_id" binding:"required"`
	DiscountID string `json:"discount_id" binding:"required"`
}

// SetOrderShippingFeeCommand 设置订单运费命令
type SetOrderShippingFeeCommand struct {
	OrderID string   `json:"order_id" binding:"required"`
	Fee     MoneyDTO `json:"fee" binding:"required"`
}

// SetOrderTaxCommand 设置订单税费命令
type SetOrderTaxCommand struct {
	OrderID string   `json:"order_id" binding:"required"`
	Tax     MoneyDTO `json:"tax" binding:"required"`
}

// SetOrderPaymentMethodCommand 设置订单支付方式命令
type SetOrderPaymentMethodCommand struct {
	OrderID       string `json:"order_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

// SetOrderShippingMethodCommand 设置订单配送方式命令
type SetOrderShippingMethodCommand struct {
	OrderID        string `json:"order_id" binding:"required"`
	ShippingMethod string `json:"shipping_method" binding:"required"`
}

// SetOrderNoteCommand 设置订单备注命令
type SetOrderNoteCommand struct {
	OrderID string `json:"order_id" binding:"required"`
	Note    string `json:"note"`
}

// SetOrderTrackingNumberCommand 设置订单物流单号命令
type SetOrderTrackingNumberCommand struct {
	OrderID        string `json:"order_id" binding:"required"`
	TrackingNumber string `json:"tracking_number" binding:"required"`
}

// SubmitOrderCommand 提交订单命令
type SubmitOrderCommand struct {
	OrderID string `json:"order_id" binding:"required"`
}

// PayOrderCommand 支付订单命令
type PayOrderCommand struct {
	OrderID       string `json:"order_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

// ShipOrderCommand 发货命令
type ShipOrderCommand struct {
	OrderID        string `json:"order_id" binding:"required"`
	TrackingNumber string `json:"tracking_number"`
}

// DeliverOrderCommand 标记订单为已送达命令
type DeliverOrderCommand struct {
	OrderID string `json:"order_id" binding:"required"`
}

// CompleteOrderCommand 完成订单命令
type CompleteOrderCommand struct {
	OrderID string `json:"order_id" binding:"required"`
}

// CancelOrderCommand 取消订单命令
type CancelOrderCommand struct {
	OrderID string `json:"order_id" binding:"required"`
}

// RequestRefundCommand 申请退款命令
type RequestRefundCommand struct {
	OrderID      string `json:"order_id" binding:"required"`
	RefundReason string `json:"refund_reason"`
}

// RefundOrderCommand 退款命令
type RefundOrderCommand struct {
	OrderID string `json:"order_id" binding:"required"`
}

// 查询DTO

// GetOrderQuery 获取订单查询
type GetOrderQuery struct {
	OrderID string `json:"order_id" binding:"required"`
}

// ListOrdersQuery 获取订单列表查询
type ListOrdersQuery struct {
	CustomerID string    `form:"customer_id"`
	Status     []string  `form:"status"`
	StartDate  time.Time `form:"start_date"`
	EndDate    time.Time `form:"end_date"`
	Page       int       `form:"page,default=1" binding:"min=1"`
	PageSize   int       `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy     string    `form:"sort_by,default=created_at"`
	SortOrder  string    `form:"sort_order,default=desc" binding:"oneof=asc desc"`
}

// SearchOrdersQuery 搜索订单查询
type SearchOrdersQuery struct {
	Keyword    string    `form:"keyword" binding:"required"`
	Status     []string  `form:"status"`
	StartDate  time.Time `form:"start_date"`
	EndDate    time.Time `form:"end_date"`
	Page       int       `form:"page,default=1" binding:"min=1"`
	PageSize   int       `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy     string    `form:"sort_by,default=created_at"`
	SortOrder  string    `form:"sort_order,default=desc" binding:"oneof=asc desc"`
}

// 响应DTO

// OrderDTO 订单DTO
type OrderDTO struct {
	ID              string           `json:"id"`
	OrderNumber     string           `json:"order_number"`
	CustomerID      string           `json:"customer_id"`
	Status          string           `json:"status"`
	StatusCode      int32            `json:"status_code"`
	Items           []OrderItemDTO   `json:"items"`
	Discounts       []OrderDiscountDTO `json:"discounts"`
	ShippingAddress AddressDTO       `json:"shipping_address"`
	BillingAddress  AddressDTO       `json:"billing_address"`
	PaymentMethod   string           `json:"payment_method"`
	ShippingMethod  string           `json:"shipping_method"`
	Subtotal        MoneyDTO         `json:"subtotal"`
	ShippingFee     MoneyDTO         `json:"shipping_fee"`
	Tax             MoneyDTO         `json:"tax"`
	DiscountAmount  MoneyDTO         `json:"discount_amount"`
	TotalAmount     MoneyDTO         `json:"total_amount"`
	Note            string           `json:"note"`
	TrackingNumber  string           `json:"tracking_number"`
	PaidAt          *time.Time       `json:"paid_at"`
	ShippedAt       *time.Time       `json:"shipped_at"`
	DeliveredAt     *time.Time       `json:"delivered_at"`
	CompletedAt     *time.Time       `json:"completed_at"`
	CancelledAt     *time.Time       `json:"cancelled_at"`
	RefundedAt      *time.Time       `json:"refunded_at"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// OrderItemDTO 订单项DTO
type OrderItemDTO struct {
	ID          string            `json:"id"`
	ProductID   string            `json:"product_id"`
	ProductName string            `json:"product_name"`
	ProductSKU  string            `json:"product_sku"`
	Quantity    int32             `json:"quantity"`
	UnitPrice   MoneyDTO          `json:"unit_price"`
	TotalPrice  MoneyDTO          `json:"total_price"`
	Attributes  map[string]string `json:"attributes"`
}

// OrderDiscountDTO 订单折扣DTO
type OrderDiscountDTO struct {
	ID           string   `json:"id"`
	DiscountType string   `json:"discount_type"`
	DiscountName string   `json:"discount_name"`
	Amount       MoneyDTO `json:"amount"`
}

// MoneyDTO 金额DTO
type MoneyDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// AddressDTO 地址DTO
type AddressDTO struct {
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	District     string `json:"district" binding:"required"`
	DetailAddress string `json:"detail_address" binding:"required"`
	PostalCode   string `json:"postal_code"`
}

// OrderListDTO 订单列表DTO
type OrderListDTO struct {
	Total  int64      `json:"total"`
	Orders []OrderDTO `json:"orders"`
}
