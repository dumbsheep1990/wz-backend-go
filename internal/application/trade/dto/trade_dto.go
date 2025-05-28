package dto

import (
	"time"

	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
)

// ---------- 订单相关 ----------

// OrderItemDTO 订单项DTO
type OrderItemDTO struct {
	ProductID  string  `json:"productId"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Currency   string  `json:"currency"`
	Quantity   int     `json:"quantity"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID          string         `json:"userId" validate:"required"`
	OrderItems      []OrderItemDTO `json:"orderItems" validate:"required,min=1"`
	ShippingAddress AddressDTO     `json:"shippingAddress" validate:"required"`
}

// CreateOrderFromCartRequest 从购物车创建订单请求
type CreateOrderFromCartRequest struct {
	UserID          string     `json:"userId" validate:"required"`
	ShippingAddress AddressDTO `json:"shippingAddress" validate:"required"`
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	OrderID string `json:"orderId" validate:"required"`
	Status  string `json:"status" validate:"required,oneof=待支付 已支付 已发货 已送达 已完成 已取消 退款中 已退款"`
}

// OrderDTO 订单DTO
type OrderDTO struct {
	ID              string         `json:"id"`
	UserID          string         `json:"userId"`
	OrderItems      []OrderItemDTO `json:"orderItems"`
	TotalAmount     float64        `json:"totalAmount"`
	Currency        string         `json:"currency"`
	Status          string         `json:"status"`
	ShippingAddress AddressDTO     `json:"shippingAddress"`
	PaymentID       string         `json:"paymentId,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

// OrderQueryRequest 订单查询请求
type OrderQueryRequest struct {
	UserID  string `json:"userId,omitempty"`
	Status  string `json:"status,omitempty"`
	Page    int    `json:"page" validate:"min=1"`
	PerPage int    `json:"perPage" validate:"min=1,max=100"`
}

// OrderQueryResponse 订单查询响应
type OrderQueryResponse struct {
	Orders     []OrderDTO `json:"orders"`
	TotalCount int64      `json:"totalCount"`
	Page       int        `json:"page"`
	PerPage    int        `json:"perPage"`
}

// ---------- 购物车相关 ----------

// CartItemDTO 购物车项DTO
type CartItemDTO struct {
	ProductID  string  `json:"productId"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Currency   string  `json:"currency"`
	Quantity   int     `json:"quantity"`
	AddedAt    time.Time `json:"addedAt"`
	Subtotal   float64 `json:"subtotal"`
}

// AddCartItemRequest 添加购物车项请求
type AddCartItemRequest struct {
	UserID    string  `json:"userId" validate:"required"`
	ProductID string  `json:"productId" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Price     float64 `json:"price" validate:"required,min=0"`
	Currency  string  `json:"currency" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
}

// UpdateCartItemQuantityRequest 更新购物车项数量请求
type UpdateCartItemQuantityRequest struct {
	UserID    string `json:"userId" validate:"required"`
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// RemoveCartItemRequest 移除购物车项请求
type RemoveCartItemRequest struct {
	UserID    string `json:"userId" validate:"required"`
	ProductID string `json:"productId" validate:"required"`
}

// ClearCartRequest 清空购物车请求
type ClearCartRequest struct {
	UserID string `json:"userId" validate:"required"`
}

// CartDTO 购物车DTO
type CartDTO struct {
	ID           string       `json:"id"`
	UserID       string       `json:"userId"`
	Items        []CartItemDTO `json:"items"`
	TotalItems   int          `json:"totalItems"`
	TotalQuantity int         `json:"totalQuantity"`
	TotalAmount  float64      `json:"totalAmount"`
	Currency     string       `json:"currency,omitempty"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

// ---------- 支付相关 ----------

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID string  `json:"orderId" validate:"required"`
	UserID  string  `json:"userId" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,min=0.01"`
	Currency string `json:"currency" validate:"required"`
	Method  string  `json:"method" validate:"required,oneof=微信支付 支付宝 银行卡 余额支付"`
}

// CompletePaymentRequest 完成支付请求
type CompletePaymentRequest struct {
	PaymentID     string `json:"paymentId" validate:"required"`
	TransactionID string `json:"transactionId" validate:"required"`
}

// FailPaymentRequest 标记支付失败请求
type FailPaymentRequest struct {
	PaymentID string `json:"paymentId" validate:"required"`
	Reason    string `json:"reason" validate:"required"`
}

// RequestRefundRequest 申请退款请求
type RequestRefundRequest struct {
	PaymentID string `json:"paymentId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
}

// CompleteRefundRequest 完成退款请求
type CompleteRefundRequest struct {
	PaymentID string `json:"paymentId" validate:"required"`
}

// PaymentDTO 支付DTO
type PaymentDTO struct {
	ID            string     `json:"id"`
	OrderID       string     `json:"orderId"`
	UserID        string     `json:"userId"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	Method        string     `json:"method"`
	Status        string     `json:"status"`
	TransactionID string     `json:"transactionId,omitempty"`
	PaymentTime   *time.Time `json:"paymentTime,omitempty"`
	RefundTime    *time.Time `json:"refundTime,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// PaymentQueryRequest 支付查询请求
type PaymentQueryRequest struct {
	UserID  string `json:"userId,omitempty"`
	OrderID string `json:"orderId,omitempty"`
	Status  string `json:"status,omitempty"`
	Page    int    `json:"page" validate:"min=1"`
	PerPage int    `json:"perPage" validate:"min=1,max=100"`
}

// PaymentQueryResponse 支付查询响应
type PaymentQueryResponse struct {
	Payments   []PaymentDTO `json:"payments"`
	TotalCount int64        `json:"totalCount"`
	Page       int          `json:"page"`
	PerPage    int          `json:"perPage"`
}

// ---------- 公共DTO ----------

// AddressDTO 地址DTO
type AddressDTO struct {
	Province    string `json:"province" validate:"required"`
	City        string `json:"city" validate:"required"`
	District    string `json:"district" validate:"required"`
	Detail      string `json:"detail" validate:"required"`
	Receiver    string `json:"receiver" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}
