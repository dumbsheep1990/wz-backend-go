package alipay

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// AliPayClient 支付宝客户端
type AliPayClient struct {
	AppID      string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Gateway    string
	NotifyURL  string
	ReturnURL  string
}

// NewAliPayClient 创建支付宝客户端
func NewAliPayClient(appID, privateKey, publicKey, notifyURL, returnURL string) (*AliPayClient, error) {
	// 这里应该解析私钥和公钥，但为了简化示例，暂时略过
	// parsedPrivateKey, err := parsePrivateKey(privateKey)
	// parsedPublicKey, err := parsePublicKey(publicKey)

	return &AliPayClient{
		AppID: appID,
		// PrivateKey: parsedPrivateKey,
		// PublicKey:  parsedPublicKey,
		Gateway:   "https://openapi.alipay.com/gateway.do",
		NotifyURL: notifyURL,
		ReturnURL: returnURL,
	}, nil
}

// CreatePayment 创建支付订单
func (c *AliPayClient) CreatePayment(orderNo, subject string, totalAmount float64) (string, error) {
	params := make(map[string]string)
	params["app_id"] = c.AppID
	params["method"] = "alipay.trade.page.pay"
	params["charset"] = "utf-8"
	params["sign_type"] = "RSA2"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	params["version"] = "1.0"
	params["notify_url"] = c.NotifyURL
	params["return_url"] = c.ReturnURL

	bizContent := make(map[string]interface{})
	bizContent["out_trade_no"] = orderNo
	bizContent["product_code"] = "FAST_INSTANT_TRADE_PAY"
	bizContent["total_amount"] = totalAmount
	bizContent["subject"] = subject

	bizContentJSON, err := json.Marshal(bizContent)
	if err != nil {
		return "", err
	}

	params["biz_content"] = string(bizContentJSON)

	// 签名
	// sign, err := c.sign(params)
	// if err != nil {
	//     return "", err
	// }
	// params["sign"] = sign
	params["sign"] = "mock_signature" // 模拟签名，实际使用时应该使用真正的签名算法

	// 构建请求URL
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return c.Gateway + "?" + values.Encode(), nil
}

// VerifyNotify 验证支付宝异步通知
func (c *AliPayClient) VerifyNotify(notifyParams map[string]string) (bool, error) {
	// 验证签名，实际使用时应该实现真正的签名验证
	// sign := notifyParams["sign"]
	// signType := notifyParams["sign_type"]
	// delete(notifyParams, "sign")
	// delete(notifyParams, "sign_type")

	// return c.verifySign(notifyParams, sign, signType)
	return true, nil // 模拟验证成功，实际使用时应该验证签名
}

// RefundPayment 退款
func (c *AliPayClient) RefundPayment(orderNo string, refundAmount float64, refundReason string) (string, error) {
	params := make(map[string]string)
	params["app_id"] = c.AppID
	params["method"] = "alipay.trade.refund"
	params["charset"] = "utf-8"
	params["sign_type"] = "RSA2"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	params["version"] = "1.0"

	bizContent := make(map[string]interface{})
	bizContent["out_trade_no"] = orderNo
	bizContent["refund_amount"] = refundAmount
	bizContent["refund_reason"] = refundReason

	bizContentJSON, err := json.Marshal(bizContent)
	if err != nil {
		return "", err
	}

	params["biz_content"] = string(bizContentJSON)

	// 签名
	// sign, err := c.sign(params)
	// if err != nil {
	//     return "", err
	// }
	// params["sign"] = sign
	params["sign"] = "mock_signature" // 模拟签名，实际使用时应该使用真正的签名算法

	// 实际应该发送HTTP请求，这里简化返回退款流水号
	return fmt.Sprintf("refund_%s_%d", orderNo, time.Now().Unix()), nil
}

// QueryPayment 查询支付状态
func (c *AliPayClient) QueryPayment(orderNo string) (map[string]interface{}, error) {
	params := make(map[string]string)
	params["app_id"] = c.AppID
	params["method"] = "alipay.trade.query"
	params["charset"] = "utf-8"
	params["sign_type"] = "RSA2"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	params["version"] = "1.0"

	bizContent := make(map[string]interface{})
	bizContent["out_trade_no"] = orderNo

	bizContentJSON, err := json.Marshal(bizContent)
	if err != nil {
		return nil, err
	}

	params["biz_content"] = string(bizContentJSON)

	// 签名
	// sign, err := c.sign(params)
	// if err != nil {
	//     return "", err
	// }
	// params["sign"] = sign
	params["sign"] = "mock_signature" // 模拟签名，实际使用时应该使用真正的签名算法

	// 实际应该发送HTTP请求，这里简化返回支付状态
	result := make(map[string]interface{})
	result["trade_status"] = "TRADE_SUCCESS" // 模拟支付成功
	result["trade_no"] = fmt.Sprintf("2021%s", orderNo)
	result["out_trade_no"] = orderNo
	result["total_amount"] = 100.00
	result["gmt_payment"] = time.Now().Format("2006-01-02 15:04:05")

	return result, nil
}
