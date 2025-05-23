package trade

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// Order 订单模型
type Order struct {
	common.BaseIDModel
	common.BaseTimeModel
	OrderID      string    `json:"orderId" db:"order_id" gorm:"uniqueIndex;not null;comment:订单ID，业务唯一标识"`
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	ProductID    int64     `json:"productId" db:"product_id" gorm:"index;not null;comment:产品ID"`
	ProductType  string    `json:"productType" db:"product_type" gorm:"not null;comment:产品类型"`
	Quantity     int       `json:"quantity" db:"quantity" gorm:"not null;comment:数量"`
	Amount       float64   `json:"amount" db:"amount" gorm:"not null;comment:金额"`
	Currency     string    `json:"currency" db:"currency" gorm:"default:CNY;comment:货币类型"`
	Status       string    `json:"status" db:"status" gorm:"index;not null;comment:订单状态"`
	PaymentID    string    `json:"paymentId" db:"payment_id" gorm:"index;comment:支付ID"`
	PaymentType  string    `json:"paymentType" db:"payment_type" gorm:"comment:支付类型"`
	PaymentTime  *time.Time `json:"paymentTime" db:"payment_time" gorm:"comment:支付时间"`
	Description  string    `json:"description" db:"description" gorm:"type:text;comment:描述"`
	Metadata     string    `json:"metadata" db:"metadata" gorm:"type:json;comment:元数据，JSON格式"`
	ClientIP     string    `json:"clientIp" db:"client_ip" gorm:"comment:客户端IP"`
	DeviceID     string    `json:"deviceId" db:"device_id" gorm:"comment:设备ID"`
	ExpireTime   *time.Time `json:"expireTime" db:"expire_time" gorm:"comment:过期时间"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// OrderItem 订单项
type OrderItem struct {
	common.BaseIDModel
	OrderID      string    `json:"orderId" db:"order_id" gorm:"index;not null;comment:订单ID"`
	ProductID    int64     `json:"productId" db:"product_id" gorm:"index;not null;comment:产品ID"`
	ProductType  string    `json:"productType" db:"product_type" gorm:"not null;comment:产品类型"`
	ProductName  string    `json:"productName" db:"product_name" gorm:"not null;comment:产品名称"`
	Quantity     int       `json:"quantity" db:"quantity" gorm:"not null;comment:数量"`
	UnitPrice    float64   `json:"unitPrice" db:"unit_price" gorm:"not null;comment:单价"`
	TotalPrice   float64   `json:"totalPrice" db:"total_price" gorm:"not null;comment:总价"`
	Discount     float64   `json:"discount" db:"discount" gorm:"default:0.00;comment:折扣"`
	Metadata     string    `json:"metadata" db:"metadata" gorm:"type:json;comment:元数据，JSON格式"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Payment 支付
type Payment struct {
	common.BaseIDModel
	common.BaseTimeModel
	PaymentID    string    `json:"paymentId" db:"payment_id" gorm:"uniqueIndex;not null;comment:支付ID，业务唯一标识"`
	OrderID      string    `json:"orderId" db:"order_id" gorm:"index;not null;comment:订单ID"`
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	Amount       float64   `json:"amount" db:"amount" gorm:"not null;comment:支付金额"`
	Currency     string    `json:"currency" db:"currency" gorm:"default:CNY;comment:货币类型"`
	PaymentType  string    `json:"paymentType" db:"payment_type" gorm:"not null;comment:支付类型"`
	Status       string    `json:"status" db:"status" gorm:"index;not null;comment:支付状态"`
	TransactionID string   `json:"transactionId" db:"transaction_id" gorm:"index;comment:第三方交易ID"`
	PaymentTime  *time.Time `json:"paymentTime" db:"payment_time" gorm:"comment:支付时间"`
	CallbackTime *time.Time `json:"callbackTime" db:"callback_time" gorm:"comment:回调时间"`
	CallbackData string    `json:"callbackData" db:"callback_data" gorm:"type:text;comment:回调原始数据"`
	ClientIP     string    `json:"clientIp" db:"client_ip" gorm:"comment:客户端IP"`
	Metadata     string    `json:"metadata" db:"metadata" gorm:"type:json;comment:元数据，JSON格式"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Refund 退款
type Refund struct {
	common.BaseIDModel
	common.BaseTimeModel
	RefundID     string    `json:"refundId" db:"refund_id" gorm:"uniqueIndex;not null;comment:退款ID，业务唯一标识"`
	OrderID      string    `json:"orderId" db:"order_id" gorm:"index;not null;comment:订单ID"`
	PaymentID    string    `json:"paymentId" db:"payment_id" gorm:"index;not null;comment:支付ID"`
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	Amount       float64   `json:"amount" db:"amount" gorm:"not null;comment:退款金额"`
	Currency     string    `json:"currency" db:"currency" gorm:"default:CNY;comment:货币类型"`
	Status       string    `json:"status" db:"status" gorm:"index;not null;comment:退款状态"`
	Reason       string    `json:"reason" db:"reason" gorm:"comment:退款原因"`
	Description  string    `json:"description" db:"description" gorm:"type:text;comment:描述"`
	ProcessedBy  string    `json:"processedBy" db:"processed_by" gorm:"comment:处理人"`
	ProcessTime  *time.Time `json:"processTime" db:"process_time" gorm:"comment:处理时间"`
	RefundTransactionID string `json:"refundTransactionId" db:"refund_transaction_id" gorm:"index;comment:退款交易ID"`
	Metadata     string    `json:"metadata" db:"metadata" gorm:"type:json;comment:元数据，JSON格式"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// AccountBalance 账户余额
type AccountBalance struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID       int64     `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_user_currency;index;not null;comment:用户ID"`
	Currency     string    `json:"currency" db:"currency" gorm:"uniqueIndex:idx_user_currency;default:CNY;comment:货币类型"`
	Available    float64   `json:"available" db:"available" gorm:"default:0.00;comment:可用余额"`
	Pending      float64   `json:"pending" db:"pending" gorm:"default:0.00;comment:待结算余额"`
	Frozen       float64   `json:"frozen" db:"frozen" gorm:"default:0.00;comment:冻结余额"`
	Total        float64   `json:"total" db:"total" gorm:"default:0.00;comment:总余额"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// Transaction 交易记录
type Transaction struct {
	common.BaseIDModel
	common.BaseTimeModel
	TransactionID string   `json:"transactionId" db:"transaction_id" gorm:"uniqueIndex;not null;comment:交易ID，业务唯一标识"`
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	RelatedID    string    `json:"relatedId" db:"related_id" gorm:"index;not null;comment:关联ID（订单ID或退款ID）"`
	RelatedType  string    `json:"relatedType" db:"related_type" gorm:"index;not null;comment:关联类型"`
	Type         string    `json:"type" db:"type" gorm:"index;not null;comment:交易类型"`
	Amount       float64   `json:"amount" db:"amount" gorm:"not null;comment:金额"`
	Currency     string    `json:"currency" db:"currency" gorm:"default:CNY;comment:货币类型"`
	BalanceBefore float64  `json:"balanceBefore" db:"balance_before" gorm:"not null;comment:交易前余额"`
	BalanceAfter float64   `json:"balanceAfter" db:"balance_after" gorm:"not null;comment:交易后余额"`
	Status       string    `json:"status" db:"status" gorm:"index;not null;comment:交易状态"`
	Description  string    `json:"description" db:"description" gorm:"type:text;comment:描述"`
	Metadata     string    `json:"metadata" db:"metadata" gorm:"type:json;comment:元数据，JSON格式"`
	OperatorID   string    `json:"operatorId" db:"operator_id" gorm:"comment:操作员ID"`
	ClientIP     string    `json:"clientIp" db:"client_ip" gorm:"comment:客户端IP"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// PaymentMethod 支付方式
type PaymentMethod struct {
	common.BaseIDModel
	common.BaseTimeModel
	MethodCode   string    `json:"methodCode" db:"method_code" gorm:"uniqueIndex;not null;comment:方式代码"`
	MethodName   string    `json:"methodName" db:"method_name" gorm:"not null;comment:方式名称"`
	MethodType   string    `json:"methodType" db:"method_type" gorm:"not null;comment:方式类型"`
	Config       string    `json:"config" db:"config" gorm:"type:json;comment:配置，JSON格式"`
	IsEnabled    bool      `json:"isEnabled" db:"is_enabled" gorm:"default:true;comment:是否启用"`
	SortOrder    int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ShoppingCart 购物车
type ShoppingCart struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID       int64     `json:"userId" db:"user_id" gorm:"uniqueIndex:idx_user_product_spec;index;not null;comment:用户ID"`
	ProductID    int64     `json:"productId" db:"product_id" gorm:"uniqueIndex:idx_user_product_spec;index;not null;comment:产品ID"`
	SpecID       int64     `json:"specId" db:"spec_id" gorm:"uniqueIndex:idx_user_product_spec;index;comment:规格ID"`
	Quantity     int       `json:"quantity" db:"quantity" gorm:"default:1;comment:数量"`
	Selected     bool      `json:"selected" db:"selected" gorm:"default:true;comment:是否选中"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ShippingAddress 收货地址
type ShippingAddress struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID        int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	ReceiverName  string    `json:"receiverName" db:"receiver_name" gorm:"not null;comment:收货人姓名"`
	ReceiverPhone string    `json:"receiverPhone" db:"receiver_phone" gorm:"not null;comment:收货人电话"`
	Province      string    `json:"province" db:"province" gorm:"not null;comment:省份"`
	City          string    `json:"city" db:"city" gorm:"not null;comment:城市"`
	District      string    `json:"district" db:"district" gorm:"not null;comment:区县"`
	DetailAddress string    `json:"detailAddress" db:"detail_address" gorm:"not null;comment:详细地址"`
	PostalCode    string    `json:"postalCode" db:"postal_code" gorm:"comment:邮政编码"`
	IsDefault     bool      `json:"isDefault" db:"is_default" gorm:"default:false;comment:是否默认"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
