package logic

import (
	"context"
	"fmt"

	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/internal/application/order/dto"
	"wz-backend-go/internal/delivery/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessPaymentLogic {
	return &ProcessPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 支付管理
func (l *ProcessPaymentLogic) ProcessPayment(in *trade.ProcessPaymentRequest) (*trade.ProcessPaymentResponse, error) {
	// 验证请求参数
	if in.OrderId == "" {
		l.Error("Invalid payment request: missing order ID")
		return &trade.ProcessPaymentResponse{
			Success: false,
			Error:   "Missing order ID",
		}, nil
	}

	if in.PaymentMethod == "" {
		l.Error("Invalid payment request: missing payment method")
		return &trade.ProcessPaymentResponse{
			Success: false,
			Error:   "Missing payment method",
		}, nil
	}

	// 获取订单信息
	orderQuery := dto.GetOrderQuery{
		OrderID: in.OrderId,
	}
	orderDTO, err := l.svcCtx.OrderApplicationService.GetOrder(l.ctx, orderQuery)
	if err != nil {
		l.Error("Failed to get order for payment", logx.Field("error", err), 
			logx.Field("orderId", in.OrderId))
		return &trade.ProcessPaymentResponse{
			Success: false,
			Error:   fmt.Sprintf("Order not found: %s", err.Error()),
		}, nil
	}

	// 检查订单状态是否允许支付
	if orderDTO.Status != "created" && orderDTO.Status != "submitted" {
		l.Error("Invalid order status for payment", 
			logx.Field("orderId", in.OrderId), 
			logx.Field("status", orderDTO.Status))
		return &trade.ProcessPaymentResponse{
			Success: false,
			Error:   fmt.Sprintf("Order status does not allow payment: %s", orderDTO.Status),
		}, nil
	}

	// 根据支付方式构建支付信息
	var paymentUrl string
	var paymentData map[string]string
	var externalRef string
	
	switch in.PaymentMethod {
	case "alipay":
		// 调用支付宝支付服务
		paymentUrl, paymentData, externalRef, err = l.processAlipayPayment(orderDTO)
	
	case "wxpay":
		// 调用微信支付服务
		paymentUrl, paymentData, externalRef, err = l.processWxPayment(orderDTO)
	
	case "bank_transfer":
		// 处理银行转账信息
		paymentUrl, paymentData, externalRef, err = l.processBankTransfer(orderDTO)
	
	default:
		err = fmt.Errorf("unsupported payment method: %s", in.PaymentMethod)
	}

	if err != nil {
		l.Error("Failed to process payment", 
			logx.Field("error", err), 
			logx.Field("orderId", in.OrderId), 
			logx.Field("method", in.PaymentMethod))
		return &trade.ProcessPaymentResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 返回支付信息
	return &trade.ProcessPaymentResponse{
		Success:     true,
		OrderId:     orderDTO.ID,
		OrderNumber: orderDTO.OrderNumber,
		Amount: &trade.Money{
			Amount:   orderDTO.TotalAmount.Amount,
			Currency: orderDTO.TotalAmount.Currency,
		},
		PaymentMethod: in.PaymentMethod,
		PaymentUrl:    paymentUrl,
		PaymentData:   paymentData,
		ExternalRef:   externalRef,
		ExpireAt:      in.ExpireAt,
	}, nil
}

// 处理支付宝支付
func (l *ProcessPaymentLogic) processAlipayPayment(orderDTO *dto.OrderDTO) (string, map[string]string, string, error) {
	// 这里应该集成实际的支付宝SDK处理逻辑
	// 目前返回模拟数据
	paymentUrl := fmt.Sprintf("https://payment.example.com/alipay?order=%s&amount=%s", 
		orderDTO.OrderNumber, orderDTO.TotalAmount.Amount)
	
	paymentData := map[string]string{
		"tradeNo":   fmt.Sprintf("ALI%s", orderDTO.OrderNumber),
		"qrCodeUrl": fmt.Sprintf("https://qr.alipay.com/%s", orderDTO.OrderNumber),
	}
	
	externalRef := fmt.Sprintf("ALI%s%s", orderDTO.OrderNumber, orderDTO.TotalAmount.Amount)
	
	return paymentUrl, paymentData, externalRef, nil
}

// 处理微信支付
func (l *ProcessPaymentLogic) processWxPayment(orderDTO *dto.OrderDTO) (string, map[string]string, string, error) {
	// 这里应该集成实际的微信支付SDK处理逻辑
	// 目前返回模拟数据
	paymentUrl := fmt.Sprintf("https://payment.example.com/wxpay?order=%s&amount=%s", 
		orderDTO.OrderNumber, orderDTO.TotalAmount.Amount)
	
	paymentData := map[string]string{
		"prepayId":  fmt.Sprintf("WX%s", orderDTO.OrderNumber),
		"qrCodeUrl": fmt.Sprintf("https://qr.wxpay.com/%s", orderDTO.OrderNumber),
	}
	
	externalRef := fmt.Sprintf("WX%s%s", orderDTO.OrderNumber, orderDTO.TotalAmount.Amount)
	
	return paymentUrl, paymentData, externalRef, nil
}

// 处理银行转账
func (l *ProcessPaymentLogic) processBankTransfer(orderDTO *dto.OrderDTO) (string, map[string]string, string, error) {
	// 返回银行转账信息
	paymentData := map[string]string{
		"bankName":    "万知银行",
		"accountName": "南京万知园信息工程有限公司",
		"accountNo":   "6225 8888 8888 8888",
		"reference":   orderDTO.OrderNumber,
	}
	
	externalRef := fmt.Sprintf("BT%s", orderDTO.OrderNumber)
	
	return "", paymentData, externalRef, nil
}
