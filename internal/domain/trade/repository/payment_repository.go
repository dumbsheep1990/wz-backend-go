package repository

import (
	"context"
	"time"
	"wz-backend-go/internal/domain/trade/entity"
	"wz-backend-go/internal/domain/trade/valueobject"
)

// PaymentRepository 支付仓储接口
type PaymentRepository interface {
	// Save 保存支付
	Save(ctx context.Context, payment *entity.Payment) error
	
	// FindByID 根据ID查找支付
	FindByID(ctx context.Context, id valueobject.PaymentID) (*entity.Payment, error)
	
	// FindByOrderID 根据订单ID查找支付
	FindByOrderID(ctx context.Context, orderID string) ([]*entity.Payment, error)
	
	// FindByUserID 根据用户ID查找支付列表
	FindByUserID(ctx context.Context, userID string, filters PaymentFilters) ([]*entity.Payment, error)
	
	// FindByTransactionID 根据交易ID查找支付
	FindByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error)
	
	// FindExpiredPayments 查找过期的支付
	FindExpiredPayments(ctx context.Context, before time.Time) ([]*entity.Payment, error)
	
	// FindByStatus 根据状态查找支付
	FindByStatus(ctx context.Context, status valueobject.PaymentStatus, filters PaymentFilters) ([]*entity.Payment, error)
	
	// FindByMethod 根据支付方式查找支付
	FindByMethod(ctx context.Context, method valueobject.PaymentMethod, filters PaymentFilters) ([]*entity.Payment, error)
	
	// Count 统计支付数量
	Count(ctx context.Context, filters PaymentFilters) (int64, error)
	
	// CountByStatus 按状态统计支付数量
	CountByStatus(ctx context.Context, status valueobject.PaymentStatus, filters PaymentFilters) (int64, error)
	
	// CountByUser 统计用户支付数量
	CountByUser(ctx context.Context, userID string, filters PaymentFilters) (int64, error)
	
	// Delete 删除支付（软删除）
	Delete(ctx context.Context, id valueobject.PaymentID) error
	
	// ExistsByOrderID 检查订单是否已有支付
	ExistsByOrderID(ctx context.Context, orderID string) (bool, error)
}

// PaymentFilters 支付查询过滤器
type PaymentFilters struct {
	UserID        string                      // 用户ID过滤
	OrderID       string                      // 订单ID过滤
	Status        *valueobject.PaymentStatus  // 状态过滤
	Method        *valueobject.PaymentMethod  // 支付方式过滤
	StartTime     *time.Time                  // 开始时间
	EndTime       *time.Time                  // 结束时间
	MinAmount     *valueobject.Money          // 最小金额
	MaxAmount     *valueobject.Money          // 最大金额
	Currency      string                      // 货币类型
	TransactionID string                      // 交易ID
	ClientIP      string                      // 客户端IP
	Limit         int                         // 限制数量
	Offset        int                         // 偏移量
	SortBy        string                      // 排序字段 (created_at, updated_at, amount)
	SortOrder     string                      // 排序顺序 (asc, desc)
}

// PaymentStatistics 支付统计
type PaymentStatistics struct {
	TotalCount       int64                              // 总数量
	TotalAmount      valueobject.Money                  // 总金额
	SuccessCount     int64                              // 成功数量
	SuccessAmount    valueobject.Money                  // 成功金额
	FailedCount      int64                              // 失败数量
	PendingCount     int64                              // 待支付数量
	RefundedCount    int64                              // 退款数量
	RefundedAmount   valueobject.Money                  // 退款金额
	MethodStats      map[string]PaymentMethodStatistics // 按支付方式统计
	DailyStats       []DailyPaymentStatistics           // 按日期统计
}

// PaymentMethodStatistics 支付方式统计
type PaymentMethodStatistics struct {
	Method        valueobject.PaymentMethod `json:"method"`
	Count         int64                     `json:"count"`
	Amount        valueobject.Money         `json:"amount"`
	SuccessCount  int64                     `json:"success_count"`
	SuccessAmount valueobject.Money         `json:"success_amount"`
	SuccessRate   float64                   `json:"success_rate"`
}

// DailyPaymentStatistics 日支付统计
type DailyPaymentStatistics struct {
	Date          time.Time         `json:"date"`
	Count         int64             `json:"count"`
	Amount        valueobject.Money `json:"amount"`
	SuccessCount  int64             `json:"success_count"`
	SuccessAmount valueobject.Money `json:"success_amount"`
	SuccessRate   float64           `json:"success_rate"`
} 