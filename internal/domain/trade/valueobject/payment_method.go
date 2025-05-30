package valueobject

import (
	"errors"
)

// PaymentMethod 支付方式值对象
type PaymentMethod struct {
	value string
}

// 支付方式常量
const (
	MethodAliPay      = "alipay"      // 支付宝
	MethodWeChatPay   = "wechatpay"   // 微信支付
	MethodPayPal      = "paypal"      // PayPal
	MethodStripe      = "stripe"      // Stripe
	MethodBankCard    = "bankcard"    // 银行卡
	MethodCreditCard  = "creditcard"  // 信用卡
	MethodBalance     = "balance"     // 余额支付
	MethodApplePay    = "applepay"    // Apple Pay
	MethodGooglePay   = "googlepay"   // Google Pay
)

// NewPaymentMethod 创建支付方式
func NewPaymentMethod(value string) (PaymentMethod, error) {
	if err := validatePaymentMethod(value); err != nil {
		return PaymentMethod{}, err
	}
	return PaymentMethod{value: value}, nil
}

// NewAliPayMethod 创建支付宝支付方式
func NewAliPayMethod() PaymentMethod {
	return PaymentMethod{value: MethodAliPay}
}

// NewWeChatPayMethod 创建微信支付方式
func NewWeChatPayMethod() PaymentMethod {
	return PaymentMethod{value: MethodWeChatPay}
}

// NewPayPalMethod 创建PayPal支付方式
func NewPayPalMethod() PaymentMethod {
	return PaymentMethod{value: MethodPayPal}
}

// NewStripeMethod 创建Stripe支付方式
func NewStripeMethod() PaymentMethod {
	return PaymentMethod{value: MethodStripe}
}

// NewBalanceMethod 创建余额支付方式
func NewBalanceMethod() PaymentMethod {
	return PaymentMethod{value: MethodBalance}
}

// Value 获取支付方式值
func (m PaymentMethod) Value() string {
	return m.value
}

// IsAliPay 是否为支付宝
func (m PaymentMethod) IsAliPay() bool {
	return m.value == MethodAliPay
}

// IsWeChatPay 是否为微信支付
func (m PaymentMethod) IsWeChatPay() bool {
	return m.value == MethodWeChatPay
}

// IsPayPal 是否为PayPal
func (m PaymentMethod) IsPayPal() bool {
	return m.value == MethodPayPal
}

// IsStripe 是否为Stripe
func (m PaymentMethod) IsStripe() bool {
	return m.value == MethodStripe
}

// IsBalance 是否为余额支付
func (m PaymentMethod) IsBalance() bool {
	return m.value == MethodBalance
}

// IsThirdParty 是否为第三方支付
func (m PaymentMethod) IsThirdParty() bool {
	return m.value == MethodAliPay || m.value == MethodWeChatPay || 
		   m.value == MethodPayPal || m.value == MethodStripe ||
		   m.value == MethodApplePay || m.value == MethodGooglePay
}

// RequiresCallback 是否需要回调处理
func (m PaymentMethod) RequiresCallback() bool {
	return m.IsThirdParty()
}

// Equals 判断两个支付方式是否相等
func (m PaymentMethod) Equals(other PaymentMethod) bool {
	return m.value == other.value
}

// String 字符串表示
func (m PaymentMethod) String() string {
	return m.value
}

// validatePaymentMethod 验证支付方式
func validatePaymentMethod(value string) error {
	validMethods := []string{
		MethodAliPay, MethodWeChatPay, MethodPayPal, MethodStripe,
		MethodBankCard, MethodCreditCard, MethodBalance,
		MethodApplePay, MethodGooglePay,
	}
	
	for _, validMethod := range validMethods {
		if value == validMethod {
			return nil
		}
	}
	
	return errors.New("无效的支付方式")
} 