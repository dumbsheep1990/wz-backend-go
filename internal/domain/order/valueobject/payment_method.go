package valueobject

import (
	"errors"
)

// PaymentMethod 支付方式值对象
type PaymentMethod int32

// 支付方式常量
const (
	PaymentMethodUnknown        PaymentMethod = 0 // 未知
	PaymentMethodAlipay         PaymentMethod = 1 // 支付宝
	PaymentMethodWechatPay      PaymentMethod = 2 // 微信支付
	PaymentMethodCreditCard     PaymentMethod = 3 // 信用卡
	PaymentMethodBankTransfer   PaymentMethod = 4 // 银行转账
	PaymentMethodCashOnDelivery PaymentMethod = 5 // 货到付款
	PaymentMethodBalance        PaymentMethod = 6 // 余额支付
)

// NewPaymentMethod 创建支付方式值对象
func NewPaymentMethod(method int32) (PaymentMethod, error) {
	switch method {
	case int32(PaymentMethodUnknown), int32(PaymentMethodAlipay), int32(PaymentMethodWechatPay),
		int32(PaymentMethodCreditCard), int32(PaymentMethodBankTransfer), int32(PaymentMethodCashOnDelivery),
		int32(PaymentMethodBalance):
		return PaymentMethod(method), nil
	default:
		return 0, errors.New("无效的支付方式")
	}
}

// Value 获取支付方式值
func (p PaymentMethod) Value() int32 {
	return int32(p)
}

// String 支付方式的字符串表示
func (p PaymentMethod) String() string {
	switch p {
	case PaymentMethodAlipay:
		return "支付宝"
	case PaymentMethodWechatPay:
		return "微信支付"
	case PaymentMethodCreditCard:
		return "信用卡"
	case PaymentMethodBankTransfer:
		return "银行转账"
	case PaymentMethodCashOnDelivery:
		return "货到付款"
	case PaymentMethodBalance:
		return "余额支付"
	default:
		return "未知支付方式"
	}
}

// IsOnlinePayment 判断是否为在线支付方式
func (p PaymentMethod) IsOnlinePayment() bool {
	return p == PaymentMethodAlipay || p == PaymentMethodWechatPay ||
		p == PaymentMethodCreditCard || p == PaymentMethodBalance
}

// RequiresExternalProcessor 判断是否需要外部支付处理器
func (p PaymentMethod) RequiresExternalProcessor() bool {
	return p == PaymentMethodAlipay || p == PaymentMethodWechatPay || p == PaymentMethodCreditCard
}

// SupportRefund 判断是否支持退款
func (p PaymentMethod) SupportRefund() bool {
	return p != PaymentMethodCashOnDelivery && p != PaymentMethodUnknown
}
