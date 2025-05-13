package model

import (
	"time"
)

// Payment 支付记录模型
type Payment struct {
	ID          int64      `json:"id" gorm:"primaryKey;column:id"`
	OrderID     int64      `json:"order_id" gorm:"column:order_id;not null;index"`
	OrderNo     string     `json:"order_no" gorm:"column:order_no;type:varchar(64);not null;index"`
	UserID      int64      `json:"user_id" gorm:"column:user_id;not null;index"`
	PayType     int        `json:"pay_type" gorm:"column:pay_type;not null"` // 1支付宝,2微信
	TradeNo     string     `json:"trade_no" gorm:"column:trade_no;type:varchar(64)"`
	TotalAmount float64    `json:"total_amount" gorm:"column:total_amount;not null;type:decimal(10,2)"`
	Status      int        `json:"status" gorm:"column:status;not null;default:0"` // 0未支付,1支付成功,2支付失败,3已退款
	PayTime     *time.Time `json:"pay_time" gorm:"column:pay_time"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 返回表名
func (Payment) TableName() string {
	return "payments"
}

// PaymentConfig 支付配置模型
type PaymentConfig struct {
	ID         int64     `json:"id" gorm:"primaryKey;column:id"`
	PayType    int       `json:"pay_type" gorm:"column:pay_type;not null;uniqueIndex"` // 1支付宝,2微信
	PayName    string    `json:"pay_name" gorm:"column:pay_name;type:varchar(32);not null"`
	AppID      string    `json:"app_id" gorm:"column:app_id;type:varchar(64);not null"`
	MerchantID string    `json:"merchant_id" gorm:"column:merchant_id;type:varchar(64)"`
	PrivateKey string    `json:"private_key" gorm:"column:private_key;type:text"`
	PublicKey  string    `json:"public_key" gorm:"column:public_key;type:text"`
	NotifyURL  string    `json:"notify_url" gorm:"column:notify_url;type:varchar(255)"`
	ReturnURL  string    `json:"return_url" gorm:"column:return_url;type:varchar(255)"`
	Status     bool      `json:"status" gorm:"column:status;not null;default:1"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 返回表名
func (PaymentConfig) TableName() string {
	return "payment_configs"
}

// PaymentNotify 支付回调通知
type PaymentNotify struct {
	OrderNo     string  `json:"order_no"`
	TradeNo     string  `json:"trade_no"`
	TotalAmount float64 `json:"total_amount"`
	PayTime     string  `json:"pay_time"`
	Status      int     `json:"status"`
}

// PaymentRequest 支付请求
type PaymentRequest struct {
	OrderNo     string  `json:"order_no" binding:"required"`
	PayType     int     `json:"pay_type" binding:"required,oneof=1 2"`
	Subject     string  `json:"subject" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required,gt=0"`
	ReturnURL   string  `json:"return_url"`
	ClientIP    string  `json:"client_ip"`
}

// PaymentResponse 支付响应
type PaymentResponse struct {
	OrderNo     string `json:"order_no"`
	PayURL      string `json:"pay_url"`
	QRCodeURL   string `json:"qr_code_url"`
	PaymentInfo string `json:"payment_info"`
}

// RefundRequest 退款请求
type RefundRequest struct {
	OrderNo      string  `json:"order_no" binding:"required"`
	RefundAmount float64 `json:"refund_amount" binding:"required,gt=0"`
	RefundReason string  `json:"refund_reason"`
}

// RefundResponse 退款响应
type RefundResponse struct {
	OrderNo      string `json:"order_no"`
	RefundNo     string `json:"refund_no"`
	RefundStatus int    `json:"refund_status"`
}
