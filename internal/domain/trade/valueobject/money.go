package valueobject

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Money 金额值对象
type Money struct {
	amount   int64  // 以分为单位存储，避免浮点数精度问题
	currency string // 货币类型
}

// 支持的货币类型
const (
	CurrencyCNY = "CNY" // 人民币
	CurrencyUSD = "USD" // 美元
	CurrencyEUR = "EUR" // 欧元
	CurrencyJPY = "JPY" // 日元
	CurrencyGBP = "GBP" // 英镑
	CurrencyHKD = "HKD" // 港币
)

// NewMoney 创建金额（以元为单位输入）
func NewMoney(amount float64, currency string) (Money, error) {
	if err := validateCurrency(currency); err != nil {
		return Money{}, err
	}
	
	if amount < 0 {
		return Money{}, errors.New("金额不能为负数")
	}
	
	// 转换为分（处理精度问题）
	amountInCents := int64(math.Round(amount * 100))
	
	return Money{
		amount:   amountInCents,
		currency: strings.ToUpper(currency),
	}, nil
}

// NewMoneyFromCents 从分创建金额
func NewMoneyFromCents(amountInCents int64, currency string) (Money, error) {
	if err := validateCurrency(currency); err != nil {
		return Money{}, err
	}
	
	if amountInCents < 0 {
		return Money{}, errors.New("金额不能为负数")
	}
	
	return Money{
		amount:   amountInCents,
		currency: strings.ToUpper(currency),
	}, nil
}

// NewZeroMoney 创建零金额
func NewZeroMoney(currency string) (Money, error) {
	return NewMoney(0, currency)
}

// Amount 获取金额（以元为单位）
func (m Money) Amount() float64 {
	return float64(m.amount) / 100
}

// AmountInCents 获取金额（以分为单位）
func (m Money) AmountInCents() int64 {
	return m.amount
}

// Currency 获取货币类型
func (m Money) Currency() string {
	return m.currency
}

// IsZero 是否为零金额
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive 是否为正数
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// Add 加法（同币种）
func (m Money) Add(other Money) (Money, error) {
	if !m.SameCurrency(other) {
		return Money{}, errors.New("不同货币不能相加")
	}
	
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract 减法（同币种）
func (m Money) Subtract(other Money) (Money, error) {
	if !m.SameCurrency(other) {
		return Money{}, errors.New("不同货币不能相减")
	}
	
	if m.amount < other.amount {
		return Money{}, errors.New("金额不足")
	}
	
	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Multiply 乘法
func (m Money) Multiply(factor float64) (Money, error) {
	if factor < 0 {
		return Money{}, errors.New("乘数不能为负数")
	}
	
	newAmount := int64(math.Round(float64(m.amount) * factor))
	
	return Money{
		amount:   newAmount,
		currency: m.currency,
	}, nil
}

// Divide 除法
func (m Money) Divide(divisor float64) (Money, error) {
	if divisor <= 0 {
		return Money{}, errors.New("除数必须大于零")
	}
	
	newAmount := int64(math.Round(float64(m.amount) / divisor))
	
	return Money{
		amount:   newAmount,
		currency: m.currency,
	}, nil
}

// GreaterThan 大于比较
func (m Money) GreaterThan(other Money) (bool, error) {
	if !m.SameCurrency(other) {
		return false, errors.New("不同货币不能比较")
	}
	return m.amount > other.amount, nil
}

// LessThan 小于比较
func (m Money) LessThan(other Money) (bool, error) {
	if !m.SameCurrency(other) {
		return false, errors.New("不同货币不能比较")
	}
	return m.amount < other.amount, nil
}

// Equals 等于比较
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// SameCurrency 是否为同一货币
func (m Money) SameCurrency(other Money) bool {
	return m.currency == other.currency
}

// GreaterThanOrEqual 大于等于比较
func (m Money) GreaterThanOrEqual(other Money) (bool, error) {
	if !m.SameCurrency(other) {
		return false, errors.New("不同货币不能比较")
	}
	return m.amount >= other.amount, nil
}

// String 字符串表示
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount(), m.currency)
}

// DisplayString 显示字符串（带货币符号）
func (m Money) DisplayString() string {
	symbol := getCurrencySymbol(m.currency)
	return fmt.Sprintf("%s%.2f", symbol, m.Amount())
}

// validateCurrency 验证货币类型
func validateCurrency(currency string) error {
	validCurrencies := []string{
		CurrencyCNY, CurrencyUSD, CurrencyEUR, 
		CurrencyJPY, CurrencyGBP, CurrencyHKD,
	}
	
	upperCurrency := strings.ToUpper(currency)
	for _, validCurrency := range validCurrencies {
		if upperCurrency == validCurrency {
			return nil
		}
	}
	
	return errors.New("不支持的货币类型")
}

// getCurrencySymbol 获取货币符号
func getCurrencySymbol(currency string) string {
	symbols := map[string]string{
		CurrencyCNY: "¥",
		CurrencyUSD: "$",
		CurrencyEUR: "€",
		CurrencyJPY: "¥",
		CurrencyGBP: "£",
		CurrencyHKD: "HK$",
	}
	
	if symbol, exists := symbols[currency]; exists {
		return symbol
	}
	return currency + " "
} 