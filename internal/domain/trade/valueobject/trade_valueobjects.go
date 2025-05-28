package valueobject

import (
	"errors"
	"regexp"
	"time"
)

// OrderID 订单ID值对象
type OrderID string

// NewOrderID 创建一个新的订单ID
func NewOrderID(id string) OrderID {
	return OrderID(id)
}

// String 返回订单ID的字符串表示
func (id OrderID) String() string {
	return string(id)
}

// IsEmpty 检查订单ID是否为空
func (id OrderID) IsEmpty() bool {
	return id == ""
}

// OrderStatus 订单状态值对象
type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"    // 待付款
	StatusPaid       OrderStatus = "PAID"       // 已付款
	StatusShipping   OrderStatus = "SHIPPING"   // 配送中
	StatusDelivered  OrderStatus = "DELIVERED"  // 已送达
	StatusCompleted  OrderStatus = "COMPLETED"  // 已完成
	StatusCancelled  OrderStatus = "CANCELLED"  // 已取消
	StatusRefunding  OrderStatus = "REFUNDING"  // 退款中
	StatusRefunded   OrderStatus = "REFUNDED"   // 已退款
)

// PaymentID 支付ID值对象
type PaymentID string

// NewPaymentID 创建一个新的支付ID
func NewPaymentID(id string) PaymentID {
	return PaymentID(id)
}

// String 返回支付ID的字符串表示
func (id PaymentID) String() string {
	return string(id)
}

// IsEmpty 检查支付ID是否为空
func (id PaymentID) IsEmpty() bool {
	return id == ""
}

// PaymentStatus 支付状态值对象
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"   // 待支付
	PaymentStatusSuccess   PaymentStatus = "SUCCESS"   // 支付成功
	PaymentStatusFailed    PaymentStatus = "FAILED"    // 支付失败
	PaymentStatusRefunding PaymentStatus = "REFUNDING" // 退款中
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"  // 已退款
)

// PaymentMethod 支付方式值对象
type PaymentMethod string

const (
	PaymentMethodWechat   PaymentMethod = "WECHAT"   // 微信支付
	PaymentMethodAlipay   PaymentMethod = "ALIPAY"   // 支付宝
	PaymentMethodBankCard PaymentMethod = "BANKCARD" // 银行卡
)

// Money 金额值对象
type Money struct {
	amount   int64  // 金额，单位：分
	currency string // 货币类型，如CNY、USD
}

// NewMoney 创建一个新的Money值对象
func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("金额不能为负数")
	}
	
	if currency == "" {
		return Money{}, errors.New("货币类型不能为空")
	}
	
	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount 获取金额
func (m Money) Amount() int64 {
	return m.amount
}

// Currency 获取货币类型
func (m Money) Currency() string {
	return m.currency
}

// Add 金额相加
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("货币类型不匹配")
	}
	
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract 金额相减
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("货币类型不匹配")
	}
	
	result := m.amount - other.amount
	if result < 0 {
		return Money{}, errors.New("结果金额不能为负数")
	}
	
	return Money{
		amount:   result,
		currency: m.currency,
	}, nil
}

// Multiply 金额乘以一个因子
func (m Money) Multiply(factor int) (Money, error) {
	if factor < 0 {
		return Money{}, errors.New("乘数不能为负数")
	}
	
	return Money{
		amount:   m.amount * int64(factor),
		currency: m.currency,
	}, nil
}

// Equals 比较两个金额是否相等
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// CartID 购物车ID值对象
type CartID string

// NewCartID 创建一个新的购物车ID
func NewCartID(id string) CartID {
	return CartID(id)
}

// String 返回购物车ID的字符串表示
func (id CartID) String() string {
	return string(id)
}

// IsEmpty 检查购物车ID是否为空
func (id CartID) IsEmpty() bool {
	return id == ""
}

// UserID 用户ID值对象
type UserID string

// NewUserID 创建一个新的用户ID
func NewUserID(id string) UserID {
	return UserID(id)
}

// String 返回用户ID的字符串表示
func (id UserID) String() string {
	return string(id)
}

// IsEmpty 检查用户ID是否为空
func (id UserID) IsEmpty() bool {
	return id == ""
}

// ProductID 商品ID值对象
type ProductID string

// NewProductID 创建一个新的商品ID
func NewProductID(id string) ProductID {
	return ProductID(id)
}

// String 返回商品ID的字符串表示
func (id ProductID) String() string {
	return string(id)
}

// IsEmpty 检查商品ID是否为空
func (id ProductID) IsEmpty() bool {
	return id == ""
}

// Quantity 数量值对象
type Quantity int

// NewQuantity 创建一个新的数量值对象
func NewQuantity(q int) (Quantity, error) {
	if q <= 0 {
		return 0, errors.New("数量必须大于0")
	}
	return Quantity(q), nil
}

// Value 返回数量的整数值
func (q Quantity) Value() int {
	return int(q)
}

// Address 地址值对象
type Address struct {
	province    string // 省份
	city        string // 城市
	district    string // 区县
	street      string // 街道
	detail      string // 详细地址
	postalCode  string // 邮政编码
	receiver    string // 收件人
	phoneNumber string // 联系电话
}

// NewAddress 创建一个新的地址值对象
func NewAddress(
	province string,
	city string,
	district string,
	street string,
	detail string,
	postalCode string,
	receiver string,
	phoneNumber string,
) (Address, error) {
	if province == "" {
		return Address{}, errors.New("省份不能为空")
	}
	if city == "" {
		return Address{}, errors.New("城市不能为空")
	}
	if receiver == "" {
		return Address{}, errors.New("收件人不能为空")
	}
	
	// 验证手机号格式
	if !validatePhoneNumber(phoneNumber) {
		return Address{}, errors.New("手机号格式不正确")
	}
	
	return Address{
		province:    province,
		city:        city,
		district:    district,
		street:      street,
		detail:      detail,
		postalCode:  postalCode,
		receiver:    receiver,
		phoneNumber: phoneNumber,
	}, nil
}

// 验证手机号格式
func validatePhoneNumber(phoneNumber string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phoneNumber)
	return matched
}

// Province 获取省份
func (a Address) Province() string {
	return a.province
}

// City 获取城市
func (a Address) City() string {
	return a.city
}

// District 获取区县
func (a Address) District() string {
	return a.district
}

// Street 获取街道
func (a Address) Street() string {
	return a.street
}

// Detail 获取详细地址
func (a Address) Detail() string {
	return a.detail
}

// PostalCode 获取邮政编码
func (a Address) PostalCode() string {
	return a.postalCode
}

// Receiver 获取收件人
func (a Address) Receiver() string {
	return a.receiver
}

// PhoneNumber 获取联系电话
func (a Address) PhoneNumber() string {
	return a.phoneNumber
}

// FullAddress 获取完整地址
func (a Address) FullAddress() string {
	return a.province + a.city + a.district + a.street + a.detail
}
