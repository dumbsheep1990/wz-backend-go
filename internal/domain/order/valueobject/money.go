package valueobject

import (
	"errors"
	"fmt"
)

// Money 金额值对象
type Money struct {
	amount   int64  // 金额，以分为单位
	currency string // 货币类型，默认CNY
}

// NewMoney 创建金额值对象
func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("金额不能为负数")
	}

	if currency == "" {
		currency = "CNY" // 默认人民币
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
		return Money{}, errors.New("不能对不同货币类型的金额进行加法运算")
	}

	return NewMoney(m.amount+other.amount, m.currency)
}

// Subtract 金额相减
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("不能对不同货币类型的金额进行减法运算")
	}

	if m.amount < other.amount {
		return Money{}, errors.New("金额不足")
	}

	return NewMoney(m.amount-other.amount, m.currency)
}

// Multiply 金额乘以数量
func (m Money) Multiply(quantity int) (Money, error) {
	if quantity < 0 {
		return Money{}, errors.New("数量不能为负数")
	}

	return NewMoney(m.amount*int64(quantity), m.currency)
}

// IsZero 判断金额是否为零
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive 判断金额是否为正数
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// Equals 判断两个金额是否相等
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String 金额的字符串表示
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.amount)/100.0, m.currency)
}

// Float 获取金额的浮点数表示（以元为单位）
func (m Money) Float() float64 {
	return float64(m.amount) / 100.0
}
