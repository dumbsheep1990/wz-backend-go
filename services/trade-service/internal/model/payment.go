package model

import (
	"time"
)

// 支付状态
const (
	PaymentStatusPending  = "pending"  // 待支付
	PaymentStatusSuccess  = "success"  // 支付成功
	PaymentStatusFailed   = "failed"   // 支付失败
	PaymentStatusCanceled = "canceled" // 已取消
	PaymentStatusRefunded = "refunded" // 已退款
)

// 支付方式
const (
	PaymentMethodAliPay    = "alipay"    // 支付宝
	PaymentMethodWeChatPay = "wechatpay" // 微信支付
	PaymentMethodPayPal    = "paypal"    // PayPal
	PaymentMethodStripe    = "stripe"    // Stripe
)

// Payment 支付模型
type Payment struct {
	ID            string    `json:"id" db:"id"`
	OrderID       int64     `json:"order_id" db:"order_id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	Status        string    `json:"status" db:"status"`
	TransactionID string    `json:"transaction_id" db:"transaction_id"` // 支付平台交易流水号
	PaymentURL    string    `json:"payment_url" db:"payment_url"`
	QRCode        string    `json:"qr_code" db:"qr_code"`
	ClientIP      string    `json:"client_ip" db:"client_ip"`
	ReturnURL     string    `json:"return_url" db:"return_url"`
	PaidAt        time.Time `json:"paid_at" db:"paid_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Refund 退款模型
type Refund struct {
	ID            string    `json:"id" db:"id"`
	PaymentID     string    `json:"payment_id" db:"payment_id"`
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	Reason        string    `json:"reason" db:"reason"`
	Status        string    `json:"status" db:"status"`
	TransactionID string    `json:"transaction_id" db:"transaction_id"` // 支付平台交易流水号
	RefundedAt    time.Time `json:"refunded_at" db:"refunded_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// GeneratePaymentID 生成支付ID
func GeneratePaymentID() string {
	now := time.Now()
	paymentID := "P" + now.Format("20060102150405") + RandomNumeric(8)
	return paymentID
}

// GenerateRefundID 生成退款ID
func GenerateRefundID() string {
	now := time.Now()
	refundID := "R" + now.Format("20060102150405") + RandomNumeric(8)
	return refundID
}
