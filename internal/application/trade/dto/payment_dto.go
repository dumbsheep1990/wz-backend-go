package dto

import (
	"time"
)

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID       string  `json:"orderId" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Currency      string  `json:"currency" validate:"required,oneof=CNY USD EUR JPY GBP HKD"`
	PaymentMethod string  `json:"paymentMethod" validate:"required,oneof=alipay wechatpay paypal stripe bankcard creditcard balance applepay googlepay"`
	ReturnURL     string  `json:"returnUrl,omitempty" validate:"omitempty,url"`
	NotifyURL     string  `json:"notifyUrl,omitempty" validate:"omitempty,url"`
	ClientIP      string  `json:"clientIp,omitempty"`
	Metadata      string  `json:"metadata,omitempty"`
}

// ProcessPaymentRequest 处理支付请求
type ProcessPaymentRequest struct {
	PaymentID string `json:"paymentId" validate:"required"`
}

// PaymentCallbackRequest 支付回调请求
type PaymentCallbackRequest struct {
	PaymentID     string                 `json:"paymentId" validate:"required"`
	TransactionID string                 `json:"transactionId" validate:"required"`
	Status        string                 `json:"status" validate:"required"`
	CallbackData  map[string]interface{} `json:"callbackData,omitempty"`
}

// CancelPaymentRequest 取消支付请求
type CancelPaymentRequest struct {
	PaymentID string `json:"paymentId" validate:"required"`
	Reason    string `json:"reason,omitempty"`
}

// RefundPaymentRequest 退款请求
type RefundPaymentRequest struct {
	PaymentID string `json:"paymentId" validate:"required"`
	Reason    string `json:"reason,omitempty"`
}

// PaymentDTO 支付数据传输对象
type PaymentDTO struct {
	ID            string     `json:"id"`
	OrderID       string     `json:"orderId"`
	UserID        string     `json:"userId"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	DisplayAmount string     `json:"displayAmount"`
	PaymentMethod string     `json:"paymentMethod"`
	Status        string     `json:"status"`
	TransactionID string     `json:"transactionId,omitempty"`
	PaymentURL    string     `json:"paymentUrl,omitempty"`
	QRCode        string     `json:"qrCode,omitempty"`
	ClientIP      string     `json:"clientIp,omitempty"`
	ReturnURL     string     `json:"returnUrl,omitempty"`
	NotifyURL     string     `json:"notifyUrl,omitempty"`
	FailureReason string     `json:"failureReason,omitempty"`
	Metadata      string     `json:"metadata,omitempty"`
	PaymentTime   *time.Time `json:"paymentTime,omitempty"`
	ExpiredAt     *time.Time `json:"expiredAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// PaymentListRequest 支付列表查询请求
type PaymentListRequest struct {
	OrderID       string     `json:"orderId" form:"orderId"`
	Status        string     `json:"status" form:"status"`
	PaymentMethod string     `json:"paymentMethod" form:"paymentMethod"`
	Currency      string     `json:"currency" form:"currency"`
	StartTime     *time.Time `json:"startTime" form:"startTime"`
	EndTime       *time.Time `json:"endTime" form:"endTime"`
	MinAmount     *float64   `json:"minAmount" form:"minAmount"`
	MaxAmount     *float64   `json:"maxAmount" form:"maxAmount"`
	Page          int        `json:"page" form:"page" validate:"min=1"`
	Size          int        `json:"size" form:"size" validate:"min=1,max=100"`
	SortBy        string     `json:"sortBy" form:"sortBy" validate:"omitempty,oneof=created_at updated_at amount"`
	SortOrder     string     `json:"sortOrder" form:"sortOrder" validate:"omitempty,oneof=asc desc"`
}

// PaymentListResponse 支付列表响应
type PaymentListResponse struct {
	Payments []PaymentDTO `json:"payments"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	Size     int          `json:"size"`
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	Payment    PaymentDTO `json:"payment"`
	PaymentURL string     `json:"paymentUrl,omitempty"`
	QRCode     string     `json:"qrCode,omitempty"`
	Message    string     `json:"message"`
}

// ProcessPaymentResponse 处理支付响应
type ProcessPaymentResponse struct {
	Payment PaymentDTO `json:"payment"`
	Message string     `json:"message"`
}

// PaymentCallbackResponse 支付回调响应
type PaymentCallbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CancelPaymentResponse 取消支付响应
type CancelPaymentResponse struct {
	Payment PaymentDTO `json:"payment"`
	Message string     `json:"message"`
}

// RefundPaymentResponse 退款响应
type RefundPaymentResponse struct {
	Payment PaymentDTO `json:"payment"`
	Message string     `json:"message"`
}

// PaymentStatisticsRequest 支付统计请求
type PaymentStatisticsRequest struct {
	StartTime     *time.Time `json:"startTime" form:"startTime"`
	EndTime       *time.Time `json:"endTime" form:"endTime"`
	Currency      string     `json:"currency" form:"currency"`
	PaymentMethod string     `json:"paymentMethod" form:"paymentMethod"`
	GroupBy       string     `json:"groupBy" form:"groupBy" validate:"omitempty,oneof=day week month method status"`
}

// PaymentStatisticsResponse 支付统计响应
type PaymentStatisticsResponse struct {
	TotalCount       int64                               `json:"totalCount"`
	TotalAmount      float64                             `json:"totalAmount"`
	Currency         string                              `json:"currency"`
	SuccessCount     int64                               `json:"successCount"`
	SuccessAmount    float64                             `json:"successAmount"`
	SuccessRate      float64                             `json:"successRate"`
	FailedCount      int64                               `json:"failedCount"`
	PendingCount     int64                               `json:"pendingCount"`
	RefundedCount    int64                               `json:"refundedCount"`
	RefundedAmount   float64                             `json:"refundedAmount"`
	MethodStats      []PaymentMethodStatisticsDTO        `json:"methodStats"`
	DailyStats       []DailyPaymentStatisticsDTO         `json:"dailyStats,omitempty"`
	TimeRange        PaymentStatisticsTimeRangeDTO       `json:"timeRange"`
}

// PaymentMethodStatisticsDTO 支付方式统计DTO
type PaymentMethodStatisticsDTO struct {
	Method        string  `json:"method"`
	Count         int64   `json:"count"`
	Amount        float64 `json:"amount"`
	SuccessCount  int64   `json:"successCount"`
	SuccessAmount float64 `json:"successAmount"`
	SuccessRate   float64 `json:"successRate"`
}

// DailyPaymentStatisticsDTO 日支付统计DTO
type DailyPaymentStatisticsDTO struct {
	Date          time.Time `json:"date"`
	Count         int64     `json:"count"`
	Amount        float64   `json:"amount"`
	SuccessCount  int64     `json:"successCount"`
	SuccessAmount float64   `json:"successAmount"`
	SuccessRate   float64   `json:"successRate"`
}

// PaymentStatisticsTimeRangeDTO 支付统计时间范围DTO
type PaymentStatisticsTimeRangeDTO struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// PaymentMethodsResponse 支付方式列表响应
type PaymentMethodsResponse struct {
	Methods []PaymentMethodDTO `json:"methods"`
}

// PaymentMethodDTO 支付方式DTO
type PaymentMethodDTO struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	IsEnabled   bool    `json:"isEnabled"`
	MinAmount   float64 `json:"minAmount"`
	MaxAmount   float64 `json:"maxAmount"`
	Description string  `json:"description"`
}

// PaymentStatusRequest 支付状态查询请求
type PaymentStatusRequest struct {
	PaymentID string `json:"paymentId" validate:"required"`
}

// PaymentStatusResponse 支付状态查询响应
type PaymentStatusResponse struct {
	PaymentID     string     `json:"paymentId"`
	Status        string     `json:"status"`
	StatusText    string     `json:"statusText"`
	PaymentTime   *time.Time `json:"paymentTime,omitempty"`
	FailureReason string     `json:"failureReason,omitempty"`
	CanRetry      bool       `json:"canRetry"`
	CanCancel     bool       `json:"canCancel"`
	CanRefund     bool       `json:"canRefund"`
}

// RetryPaymentRequest 重试支付请求
type RetryPaymentRequest struct {
	PaymentID     string `json:"paymentId" validate:"required"`
	PaymentMethod string `json:"paymentMethod,omitempty"`
	ReturnURL     string `json:"returnUrl,omitempty" validate:"omitempty,url"`
	NotifyURL     string `json:"notifyUrl,omitempty" validate:"omitempty,url"`
}

// RetryPaymentResponse 重试支付响应
type RetryPaymentResponse struct {
	NewPayment PaymentDTO `json:"newPayment"`
	PaymentURL string     `json:"paymentUrl,omitempty"`
	QRCode     string     `json:"qrCode,omitempty"`
	Message    string     `json:"message"`
} 