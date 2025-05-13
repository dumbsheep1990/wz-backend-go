package wechat

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// WechatPayClient 微信支付客户端
type WechatPayClient struct {
	AppID     string
	MchID     string
	APIKey    string
	NotifyURL string
	TradeType string
	CertFile  string
	KeyFile   string
}

// WXPayResponse 微信支付响应
type WXPayResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
}

// NewWechatPayClient 创建微信支付客户端
func NewWechatPayClient(appID, mchID, apiKey, notifyURL string) *WechatPayClient {
	return &WechatPayClient{
		AppID:     appID,
		MchID:     mchID,
		APIKey:    apiKey,
		NotifyURL: notifyURL,
		TradeType: "NATIVE", // 二维码支付
	}
}

// CreatePayment 创建微信支付订单
func (c *WechatPayClient) CreatePayment(orderNo, body string, totalFee int) (*WXPayResponse, error) {
	params := make(map[string]string)
	params["appid"] = c.AppID
	params["mch_id"] = c.MchID
	params["nonce_str"] = c.generateNonceStr()
	params["body"] = body
	params["out_trade_no"] = orderNo
	params["total_fee"] = fmt.Sprintf("%d", totalFee) // 微信支付金额单位为分
	params["spbill_create_ip"] = "127.0.0.1"          // 实际使用中应该是用户的真实IP
	params["notify_url"] = c.NotifyURL
	params["trade_type"] = c.TradeType

	// 签名
	// sign := c.Sign(params)
	// params["sign"] = sign
	params["sign"] = "mock_signature" // 模拟签名，实际使用时应该使用真正的签名算法

	// 生成请求XML
	xmlStr, err := c.generateXML(params)
	if err != nil {
		return nil, err
	}

	// 实际应该发送HTTP请求，这里简化返回模拟响应
	return &WXPayResponse{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
		ResultCode: "SUCCESS",
		PrepayID:   fmt.Sprintf("wx%s%d", orderNo, time.Now().Unix()),
		CodeURL:    fmt.Sprintf("weixin://wxpay/bizpayurl?pr=%s", orderNo),
	}, nil
}

// VerifyNotify 验证微信支付回调
func (c *WechatPayClient) VerifyNotify(notifyData string) (bool, string, error) {
	// 实际使用时应该解析XML并验证签名
	// notifyMap, err := c.parseXML(notifyData)
	// if err != nil {
	//     return false, "", err
	// }
	//
	// sign := notifyMap["sign"]
	// delete(notifyMap, "sign")
	//
	// if c.Sign(notifyMap) != sign {
	//     return false, "", fmt.Errorf("签名验证失败")
	// }
	//
	// return true, notifyMap["out_trade_no"], nil

	// 模拟验证成功，返回订单号
	return true, "mock_order_no", nil
}

// RefundPayment 退款
func (c *WechatPayClient) RefundPayment(orderNo, refundNo string, totalFee, refundFee int) (map[string]string, error) {
	params := make(map[string]string)
	params["appid"] = c.AppID
	params["mch_id"] = c.MchID
	params["nonce_str"] = c.generateNonceStr()
	params["out_trade_no"] = orderNo
	params["out_refund_no"] = refundNo
	params["total_fee"] = fmt.Sprintf("%d", totalFee)
	params["refund_fee"] = fmt.Sprintf("%d", refundFee)

	// 签名
	// sign := c.Sign(params)
	// params["sign"] = sign
	params["sign"] = "mock_signature" // 模拟签名，实际使用时应该使用真正的签名算法

	// 实际应该发送HTTP请求，这里简化返回模拟响应
	result := make(map[string]string)
	result["return_code"] = "SUCCESS"
	result["result_code"] = "SUCCESS"
	result["refund_id"] = fmt.Sprintf("refund_%s_%d", orderNo, time.Now().Unix())

	return result, nil
}

// QueryPayment 查询支付状态
func (c *WechatPayClient) QueryPayment(orderNo string) (map[string]string, error) {
	params := make(map[string]string)
	params["appid"] = c.AppID
	params["mch_id"] = c.MchID
	params["out_trade_no"] = orderNo
	params["nonce_str"] = c.generateNonceStr()

	// 签名
	// sign := c.Sign(params)
	// params["sign"] = sign
	params["sign"] = "mock_signature" // 模拟签名，实际使用时应该使用真正的签名算法

	// 实际应该发送HTTP请求，这里简化返回模拟响应
	result := make(map[string]string)
	result["return_code"] = "SUCCESS"
	result["result_code"] = "SUCCESS"
	result["trade_state"] = "SUCCESS" // 模拟支付成功
	result["transaction_id"] = fmt.Sprintf("wxpay%s", orderNo)
	result["out_trade_no"] = orderNo
	result["total_fee"] = "100" // 假设100分，即1元
	result["time_end"] = time.Now().Format("20060102150405")

	return result, nil
}

// 生成随机字符串
func (c *WechatPayClient) generateNonceStr() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 32
	result := make([]byte, length)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}

// 生成XML
func (c *WechatPayClient) generateXML(params map[string]string) (string, error) {
	var buf strings.Builder
	buf.WriteString("<xml>")
	for k, v := range params {
		buf.WriteString(fmt.Sprintf("<%s>%s</%s>", k, v, k))
	}
	buf.WriteString("</xml>")
	return buf.String(), nil
}

// 解析XML
func (c *WechatPayClient) parseXML(xmlStr string) (map[string]string, error) {
	result := make(map[string]string)

	// 实际使用时应该实现真正的XML解析
	// decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	// for {
	//     token, err := decoder.Token()
	//     if err != nil {
	//         break
	//     }
	//     if el, ok := token.(xml.StartElement); ok {
	//         var data string
	//         decoder.DecodeElement(&data, &el)
	//         result[el.Name.Local] = data
	//     }
	// }

	return result, nil
}

// Sign 签名
func (c *WechatPayClient) Sign(params map[string]string) string {
	var keys []string
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	var buf strings.Builder
	for _, k := range keys {
		if params[k] != "" {
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(params[k])
			buf.WriteString("&")
		}
	}
	buf.WriteString("key=")
	buf.WriteString(c.APIKey)

	// 计算MD5
	h := md5.New()
	h.Write([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
