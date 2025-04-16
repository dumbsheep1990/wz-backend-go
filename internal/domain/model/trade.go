package model

import (
	"time"
)

// Order 订单模型
type Order struct {
	ID           int64      `json:"id" db:"id"`
	OrderID      string     `json:"order_id" db:"order_id"`           // 订单ID，业务唯一标识
	UserID       int64      `json:"user_id" db:"user_id"`             // 用户ID
	ProductID    int64      `json:"product_id" db:"product_id"`       // 产品ID
	ProductType  string     `json:"product_type" db:"product_type"`   // 产品类型
	Quantity     int        `json:"quantity" db:"quantity"`           // 数量
	Amount       float64    `json:"amount" db:"amount"`               // 金额
	Currency     string     `json:"currency" db:"currency"`           // 货币类型
	Status       string     `json:"status" db:"status"`               // 订单状态
	PaymentID    string     `json:"payment_id" db:"payment_id"`       // 支付ID
	PaymentType  string     `json:"payment_type" db:"payment_type"`   // 支付类型
	PaymentTime  *time.Time `json:"payment_time" db:"payment_time"`   // 支付时间
	Description  string     `json:"description" db:"description"`     // 描述
	Metadata     string     `json:"metadata" db:"metadata"`           // 元数据，JSON格式
	ClientIP     string     `json:"client_ip" db:"client_ip"`         // 客户端IP
	DeviceID     string     `json:"device_id" db:"device_id"`         // 设备ID
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`       // 创建时间
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`       // 更新时间
	ExpireTime   *time.Time `json:"expire_time" db:"expire_time"`     // 过期时间
	OrderItems   []OrderItem `json:"order_items,omitempty" db:"-"`    // 订单项列表
}

// OrderItem 订单项模型
type OrderItem struct {
	ID          int64      `json:"id" db:"id"`
	OrderID     string     `json:"order_id" db:"order_id"`             // 订单ID
	ProductID   int64      `json:"product_id" db:"product_id"`         // 产品ID
	ProductType string     `json:"product_type" db:"product_type"`     // 产品类型
	ProductName string     `json:"product_name" db:"product_name"`     // 产品名称
	Quantity    int        `json:"quantity" db:"quantity"`             // 数量
	UnitPrice   float64    `json:"unit_price" db:"unit_price"`         // 单价
	TotalPrice  float64    `json:"total_price" db:"total_price"`       // 总价
	Discount    float64    `json:"discount" db:"discount"`             // 折扣
	Metadata    string     `json:"metadata" db:"metadata"`             // 元数据，JSON格式
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`         // 创建时间
}

// Payment 支付模型
type Payment struct {
	ID              int64      `json:"id" db:"id"`
	PaymentID       string     `json:"payment_id" db:"payment_id"`               // 支付ID，业务唯一标识
	OrderID         string     `json:"order_id" db:"order_id"`                   // 订单ID
	UserID          int64      `json:"user_id" db:"user_id"`                     // 用户ID
	Amount          float64    `json:"amount" db:"amount"`                       // 支付金额
	Currency        string     `json:"currency" db:"currency"`                   // 货币类型
	PaymentType     string     `json:"payment_type" db:"payment_type"`           // 支付类型
	Status          string     `json:"status" db:"status"`                       // 支付状态
	TransactionID   string     `json:"transaction_id" db:"transaction_id"`       // 第三方交易ID
	PaymentTime     *time.Time `json:"payment_time" db:"payment_time"`           // 支付时间
	CallbackTime    *time.Time `json:"callback_time" db:"callback_time"`         // 回调时间
	CallbackData    string     `json:"callback_data" db:"callback_data"`         // 回调原始数据
	ClientIP        string     `json:"client_ip" db:"client_ip"`                 // 客户端IP
	Metadata        string     `json:"metadata" db:"metadata"`                   // 元数据，JSON格式
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`               // 创建时间
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`               // 更新时间
}

// Refund 退款模型
type Refund struct {
	ID                  int64      `json:"id" db:"id"`
	RefundID            string     `json:"refund_id" db:"refund_id"`                           // 退款ID，业务唯一标识
	OrderID             string     `json:"order_id" db:"order_id"`                             // 订单ID
	PaymentID           string     `json:"payment_id" db:"payment_id"`                         // 支付ID
	UserID              int64      `json:"user_id" db:"user_id"`                               // 用户ID
	Amount              float64    `json:"amount" db:"amount"`                                 // 退款金额
	Currency            string     `json:"currency" db:"currency"`                             // 货币类型
	Status              string     `json:"status" db:"status"`                                 // 退款状态
	Reason              string     `json:"reason" db:"reason"`                                 // 退款原因
	Description         string     `json:"description" db:"description"`                       // 描述
	ProcessedBy         string     `json:"processed_by" db:"processed_by"`                     // 处理人
	ProcessTime         *time.Time `json:"process_time" db:"process_time"`                     // 处理时间
	RefundTransactionID string     `json:"refund_transaction_id" db:"refund_transaction_id"`   // 退款交易ID
	Metadata            string     `json:"metadata" db:"metadata"`                             // 元数据，JSON格式
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`                         // 创建时间
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`                         // 更新时间
}

// AccountBalance 账户余额模型
type AccountBalance struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`           // 用户ID
	Currency   string    `json:"currency" db:"currency"`         // 货币类型
	Available  float64   `json:"available" db:"available"`       // 可用余额
	Pending    float64   `json:"pending" db:"pending"`           // 待结算余额
	Frozen     float64   `json:"frozen" db:"frozen"`             // 冻结余额
	Total      float64   `json:"total" db:"total"`               // 总余额
	CreatedAt  time.Time `json:"created_at" db:"created_at"`     // 创建时间
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`     // 更新时间
}

// Transaction 交易记录模型
type Transaction struct {
	ID             int64     `json:"id" db:"id"`
	TransactionID  string    `json:"transaction_id" db:"transaction_id"`    // 交易ID，业务唯一标识
	UserID         int64     `json:"user_id" db:"user_id"`                  // 用户ID
	RelatedID      string    `json:"related_id" db:"related_id"`            // 关联ID（订单ID或退款ID）
	RelatedType    string    `json:"related_type" db:"related_type"`        // 关联类型
	Type           string    `json:"type" db:"type"`                        // 交易类型
	Amount         float64   `json:"amount" db:"amount"`                    // 金额
	Currency       string    `json:"currency" db:"currency"`                // 货币类型
	BalanceBefore  float64   `json:"balance_before" db:"balance_before"`    // 交易前余额
	BalanceAfter   float64   `json:"balance_after" db:"balance_after"`      // 交易后余额
	Status         string    `json:"status" db:"status"`                    // 交易状态
	Description    string    `json:"description" db:"description"`          // 描述
	Metadata       string    `json:"metadata" db:"metadata"`                // 元数据，JSON格式
	OperatorID     string    `json:"operator_id" db:"operator_id"`          // 操作员ID
	ClientIP       string    `json:"client_ip" db:"client_ip"`              // 客户端IP
	CreatedAt      time.Time `json:"created_at" db:"created_at"`            // 创建时间
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`            // 更新时间
}

// FinancialDailyReport 财务日报表模型
type FinancialDailyReport struct {
	ID           int64     `json:"id" db:"id"`
	ReportDate   time.Time `json:"report_date" db:"report_date"`         // 报表日期
	Currency     string    `json:"currency" db:"currency"`               // 货币类型
	Income       float64   `json:"income" db:"income"`                   // 收入
	Refund       float64   `json:"refund" db:"refund"`                   // 退款
	Net          float64   `json:"net" db:"net"`                         // 净收入
	OrderCount   int       `json:"order_count" db:"order_count"`         // 订单数
	PaymentCount int       `json:"payment_count" db:"payment_count"`     // 支付数
	RefundCount  int       `json:"refund_count" db:"refund_count"`       // 退款数
	UserCount    int       `json:"user_count" db:"user_count"`           // 用户数
	CreatedAt    time.Time `json:"created_at" db:"created_at"`           // 创建时间
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`           // 更新时间
}

// FinancialMonthlyReport 财务月报表模型
type FinancialMonthlyReport struct {
	ID           int64     `json:"id" db:"id"`
	ReportYear   int       `json:"report_year" db:"report_year"`         // 报表年份
	ReportMonth  int       `json:"report_month" db:"report_month"`       // 报表月份
	Currency     string    `json:"currency" db:"currency"`               // 货币类型
	Income       float64   `json:"income" db:"income"`                   // 收入
	Refund       float64   `json:"refund" db:"refund"`                   // 退款
	Net          float64   `json:"net" db:"net"`                         // 净收入
	OrderCount   int       `json:"order_count" db:"order_count"`         // 订单数
	PaymentCount int       `json:"payment_count" db:"payment_count"`     // 支付数
	RefundCount  int       `json:"refund_count" db:"refund_count"`       // 退款数
	UserCount    int       `json:"user_count" db:"user_count"`           // 用户数
	CreatedAt    time.Time `json:"created_at" db:"created_at"`           // 创建时间
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`           // 更新时间
}

// PaymentMethod 支付方式模型
type PaymentMethod struct {
	ID         int64     `json:"id" db:"id"`
	MethodCode string    `json:"method_code" db:"method_code"`     // 方式代码
	MethodName string    `json:"method_name" db:"method_name"`     // 方式名称
	MethodType string    `json:"method_type" db:"method_type"`     // 方式类型
	Config     string    `json:"config" db:"config"`               // 配置，JSON格式
	IsEnabled  bool      `json:"is_enabled" db:"is_enabled"`       // 是否启用
	SortOrder  int       `json:"sort_order" db:"sort_order"`       // 排序
	CreatedAt  time.Time `json:"created_at" db:"created_at"`       // 创建时间
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`       // 更新时间
}

// 订单状态常量
const (
	OrderStatusPending  = "pending"  // 待支付
	OrderStatusPaid     = "paid"     // 已支付
	OrderStatusCanceled = "canceled" // 已取消
	OrderStatusRefunded = "refunded" // 已退款
	OrderStatusExpired  = "expired"  // 已过期
)

// 支付状态常量
const (
	PaymentStatusPending = "pending" // 处理中
	PaymentStatusSuccess = "success" // 成功
	PaymentStatusFailed  = "failed"  // 失败
)

// 退款状态常量
const (
	RefundStatusPending   = "pending"    // 待处理
	RefundStatusApproved  = "approved"   // 已批准
	RefundStatusProcessing = "processing" // 处理中
	RefundStatusSuccess    = "success"    // 成功
	RefundStatusRejected   = "rejected"   // 已拒绝
	RefundStatusFailed     = "failed"     // 失败
)

// 交易类型常量
const (
	TransactionTypePayment    = "payment"    // 支付
	TransactionTypeRefund     = "refund"     // 退款
	TransactionTypeWithdraw   = "withdraw"   // 提现
	TransactionTypeAdjustment = "adjustment" // 调整
)

// TradeService 交易服务接口
type TradeService interface {
	// 订单相关
	CreateOrder(order *Order) (*Order, error)
	GetOrder(orderID string, userID int64) (*Order, error)
	ListOrders(userID int64, status, startTime, endTime, productType string, page, pageSize int) ([]*Order, int, error)
	CancelOrder(orderID string, userID int64, reason string) error
	UpdateOrderStatus(orderID, status, operatorID, reason string) error
	
	// 支付相关
	ProcessPayment(orderID, paymentType string, amount float64, currency, returnURL, notifyURL, clientIP, metadata string) (*Payment, error)
	HandlePaymentCallback(callbackData map[string]interface{}) error
	
	// 退款相关
	CreateRefund(orderID string, userID int64, amount float64, reason, description string) (*Refund, error)
	GetRefund(refundID string, userID int64) (*Refund, error)
	ListRefunds(userID int64, orderID, status, startTime, endTime string, page, pageSize int) ([]*Refund, int, error)
	ProcessRefund(refundID, action, comment, processedBy string) error
	
	// 账户相关
	GetBalance(userID int64) ([]*AccountBalance, error)
	GetTransactions(userID int64, transactionType, status, startTime, endTime string, page, pageSize int) ([]*Transaction, int, error)
	
	// 报表相关
	GetFinancialReport(startTime, endTime, reportType, currency string) (interface{}, error)
}

// TradeRepository 交易仓储接口
type TradeRepository interface {
	// 订单相关
	SaveOrder(order *Order) error
	GetOrder(orderID string) (*Order, error)
	GetOrderByPaymentID(paymentID string) (*Order, error)
	UpdateOrder(order *Order) error
	ListOrders(userID int64, status, startTime, endTime, productType string, page, pageSize int) ([]*Order, int, error)
	
	// 订单项相关
	SaveOrderItems(items []*OrderItem) error
	GetOrderItems(orderID string) ([]*OrderItem, error)
	
	// 支付相关
	SavePayment(payment *Payment) error
	GetPayment(paymentID string) (*Payment, error)
	GetPaymentByTransactionID(transactionID string) (*Payment, error)
	UpdatePayment(payment *Payment) error
	ListPayments(userID int64, orderID, status string, page, pageSize int) ([]*Payment, int, error)
	
	// 退款相关
	SaveRefund(refund *Refund) error
	GetRefund(refundID string) (*Refund, error)
	UpdateRefund(refund *Refund) error
	ListRefunds(userID int64, orderID, status, startTime, endTime string, page, pageSize int) ([]*Refund, int, error)
	
	// 账户相关
	GetAccountBalance(userID int64, currency string) (*AccountBalance, error)
	UpdateAccountBalance(balance *AccountBalance) error
	GetAllUserBalances(userID int64) ([]*AccountBalance, error)
	
	// 交易记录相关
	SaveTransaction(transaction *Transaction) error
	GetTransaction(transactionID string) (*Transaction, error)
	ListTransactions(userID int64, transactionType, status, startTime, endTime string, page, pageSize int) ([]*Transaction, int, error)
	
	// 报表相关
	GetDailyReport(date time.Time, currency string) (*FinancialDailyReport, error)
	SaveDailyReport(report *FinancialDailyReport) error
	UpdateDailyReport(report *FinancialDailyReport) error
	GetDailyReports(startDate, endDate time.Time, currency string) ([]*FinancialDailyReport, error)
	
	GetMonthlyReport(year, month int, currency string) (*FinancialMonthlyReport, error)
	SaveMonthlyReport(report *FinancialMonthlyReport) error
	UpdateMonthlyReport(report *FinancialMonthlyReport) error
	GetMonthlyReports(startYear, startMonth, endYear, endMonth int, currency string) ([]*FinancialMonthlyReport, error)
	
	// 支付方式相关
	GetPaymentMethods(enabled bool) ([]*PaymentMethod, error)
	GetPaymentMethod(methodCode string) (*PaymentMethod, error)
}
