package entity

import (
	"errors"
	"time"
	ordervo "wz-backend-go/internal/domain/order/valueobject"
)

// DiscountType 折扣类型
type DiscountType int32

const (
	DiscountTypeAmount    DiscountType = 1 // 金额折扣
	DiscountTypePercent   DiscountType = 2 // 百分比折扣
	DiscountTypeCoupon    DiscountType = 3 // 优惠券
	DiscountTypePromotion DiscountType = 4 // 促销活动
)

// OrderDiscount 订单折扣实体
type OrderDiscount struct {
	id            string
	name          string
	discountType  DiscountType
	value         float64       // 折扣值：金额折扣为具体金额，百分比折扣为0-100之间的值
	amount        ordervo.Money // 折扣金额
	code          string        // 折扣码或优惠券码
	description   string
	applyToItems  []string      // 适用的订单项ID，为空表示适用于整个订单
	minOrderValue ordervo.Money // 最低订单金额要求
	startTime     time.Time     // 生效开始时间
	endTime       time.Time     // 生效结束时间
	isActive      bool          // 是否有效
}

// NewOrderDiscount 创建订单折扣
func NewOrderDiscount(
	id string,
	name string,
	discountType DiscountType,
	value float64,
	code string,
	description string,
	minOrderValue ordervo.Money,
	startTime time.Time,
	endTime time.Time,
) (*OrderDiscount, error) {
	if name == "" {
		return nil, errors.New("折扣名称不能为空")
	}

	if discountType == DiscountTypePercent && (value < 0 || value > 100) {
		return nil, errors.New("百分比折扣必须在0-100之间")
	}

	if discountType == DiscountTypeAmount && value < 0 {
		return nil, errors.New("金额折扣不能为负数")
	}

	// 创建零值金额作为初始折扣金额
	zeroMoney, _ := ordervo.NewMoney(0, "CNY")

	return &OrderDiscount{
		id:            id,
		name:          name,
		discountType:  discountType,
		value:         value,
		amount:        zeroMoney,
		code:          code,
		description:   description,
		applyToItems:  []string{},
		minOrderValue: minOrderValue,
		startTime:     startTime,
		endTime:       endTime,
		isActive:      true,
	}, nil
}

// ID 获取折扣ID
func (d *OrderDiscount) ID() string {
	return d.id
}

// Name 获取折扣名称
func (d *OrderDiscount) Name() string {
	return d.name
}

// DiscountType 获取折扣类型
func (d *OrderDiscount) DiscountType() DiscountType {
	return d.discountType
}

// Value 获取折扣值
func (d *OrderDiscount) Value() float64 {
	return d.value
}

// Amount 获取折扣金额
func (d *OrderDiscount) Amount() ordervo.Money {
	return d.amount
}

// Code 获取折扣码
func (d *OrderDiscount) Code() string {
	return d.code
}

// Description 获取折扣描述
func (d *OrderDiscount) Description() string {
	return d.description
}

// ApplyToItems 获取适用的订单项ID
func (d *OrderDiscount) ApplyToItems() []string {
	return d.applyToItems
}

// MinOrderValue 获取最低订单金额要求
func (d *OrderDiscount) MinOrderValue() ordervo.Money {
	return d.minOrderValue
}

// StartTime 获取生效开始时间
func (d *OrderDiscount) StartTime() time.Time {
	return d.startTime
}

// EndTime 获取生效结束时间
func (d *OrderDiscount) EndTime() time.Time {
	return d.endTime
}

// IsActive 获取是否有效
func (d *OrderDiscount) IsActive() bool {
	return d.isActive
}

// Deactivate 停用折扣
func (d *OrderDiscount) Deactivate() {
	d.isActive = false
}

// Activate 启用折扣
func (d *OrderDiscount) Activate() {
	d.isActive = true
}

// IsValid 检查折扣是否有效
func (d *OrderDiscount) IsValid(orderTotal ordervo.Money) bool {
	now := time.Now()

	// 检查折扣是否处于有效期内
	if !d.isActive || now.Before(d.startTime) || now.After(d.endTime) {
		return false
	}

	// 检查订单金额是否满足最低要求
	if !d.minOrderValue.IsZero() && orderTotal.Amount() < d.minOrderValue.Amount() {
		return false
	}

	return true
}

// ApplyTo 设置适用的订单项
func (d *OrderDiscount) ApplyTo(itemIDs []string) {
	d.applyToItems = itemIDs
}

// ApplyToOrder 应用于整个订单
func (d *OrderDiscount) ApplyToOrder() {
	d.applyToItems = []string{}
}

// IsAppliedToEntireOrder 检查是否应用于整个订单
func (d *OrderDiscount) IsAppliedToEntireOrder() bool {
	return len(d.applyToItems) == 0
}

// CalculateDiscount 计算折扣金额
func (d *OrderDiscount) CalculateDiscount(orderTotal ordervo.Money) (ordervo.Money, error) {
	if !d.IsValid(orderTotal) {
		zeroMoney, _ := ordervo.NewMoney(0, orderTotal.Currency())
		return zeroMoney, errors.New("折扣无效或不适用")
	}

	var discountAmount int64

	switch d.discountType {
	case DiscountTypeAmount:
		// 金额折扣直接使用设定的值
		discountAmount = int64(d.value * 100) // 转换为分
	case DiscountTypePercent:
		// 百分比折扣需要计算
		discountAmount = int64(float64(orderTotal.Amount()) * d.value / 100)
	default:
		// 其他类型折扣逻辑
		discountAmount = 0
	}

	// 确保折扣金额不超过订单总额
	if discountAmount > orderTotal.Amount() {
		discountAmount = orderTotal.Amount()
	}

	// 创建折扣金额值对象
	amount, err := ordervo.NewMoney(discountAmount, orderTotal.Currency())
	if err != nil {
		return ordervo.Money{}, err
	}

	// 更新折扣金额
	d.amount = amount

	return amount, nil
}
