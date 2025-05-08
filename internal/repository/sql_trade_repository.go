package repository

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Order 订单实体
type Order struct {
	OrderId     string
	UserId      int64
	TenantId    int64
	ProductId   int64
	ProductType string
	ProductName string
	Quantity    int32
	Amount      float64
	Currency    string
	Status      string
	PaymentId   string
	PaymentType string
	PaymentTime time.Time
	Description string
	Metadata    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpireTime  time.Time
}

// Refund 退款实体
type Refund struct {
	RefundId    string
	OrderId     string
	UserId      int64
	TenantId    int64
	Amount      float64
	Currency    string
	Status      string
	Reason      string
	Description string
	ProcessedBy string
	ProcessTime time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Transaction 交易记录实体
type Transaction struct {
	TransactionId string
	UserId        int64
	TenantId      int64
	RelatedId     string
	Type          string
	Amount        float64
	Currency      string
	Status        string
	Description   string
	Metadata      string
	CreatedAt     time.Time
}

// ReportItem 报表项
type ReportItem struct {
	Date         string
	OrderCount   int64
	OrderAmount  float64
	RefundCount  int64
	RefundAmount float64
	NetAmount    float64
	CurrencyUnit string
}

// TradeRepository 交易数据仓库接口
type TradeRepository interface {
	// 订单相关方法
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetOrderList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Order, int64, error)
	UpdateOrderStatus(ctx context.Context, id string, status, reason string, operatorId int64) error

	// 退款相关方法
	GetRefundById(ctx context.Context, id string) (*Refund, error)
	GetRefundList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Refund, int64, error)
	ProcessRefund(ctx context.Context, id, action, comment string, operatorId int64) error

	// 交易记录相关方法
	GetTransactionList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Transaction, int64, error)

	// 财务报表相关方法
	GetFinancialReport(ctx context.Context, startTime, endTime string, tenantId int64, groupBy string) ([]*ReportItem, float64, error)
}

// SqlTradeRepository SQL交易仓库实现
type SqlTradeRepository struct {
	conn sqlx.SqlConn
}

// NewSqlTradeRepository 创建交易仓库实例
func NewSqlTradeRepository(conn sqlx.SqlConn) TradeRepository {
	return &SqlTradeRepository{
		conn: conn,
	}
}

// GetOrderById 根据ID获取订单
func (r *SqlTradeRepository) GetOrderById(ctx context.Context, id string) (*Order, error) {
	var order Order
	query := `SELECT order_id, user_id, tenant_id, product_id, product_type, product_name, 
	          quantity, amount, currency, status, payment_id, payment_type, payment_time, 
	          description, metadata, created_at, updated_at, expire_time 
	          FROM orders WHERE order_id = ?`
	
	err := r.conn.QueryRowCtx(ctx, &order, query, id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetOrderList 获取订单列表
func (r *SqlTradeRepository) GetOrderList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if filters != nil {
		if orderId, ok := filters["orderId"].(string); ok && orderId != "" {
			whereClause += " AND order_id = ?"
			args = append(args, orderId)
		}
		if userId, ok := filters["userId"].(int64); ok && userId > 0 {
			whereClause += " AND user_id = ?"
			args = append(args, userId)
		}
		if status, ok := filters["status"].(string); ok && status != "" {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if productType, ok := filters["productType"].(string); ok && productType != "" {
			whereClause += " AND product_type = ?"
			args = append(args, productType)
		}
		if startTime, ok := filters["startTime"].(string); ok && startTime != "" {
			whereClause += " AND created_at >= ?"
			args = append(args, startTime)
		}
		if endTime, ok := filters["endTime"].(string); ok && endTime != "" {
			whereClause += " AND created_at <= ?"
			args = append(args, endTime)
		}
		if tenantId, ok := filters["tenantId"].(int64); ok && tenantId > 0 {
			whereClause += " AND tenant_id = ?"
			args = append(args, tenantId)
		}
	}

	// 执行查询
	query := fmt.Sprintf(`SELECT order_id, user_id, tenant_id, product_id, product_type, product_name, 
	                     quantity, amount, currency, status, payment_id, payment_type, payment_time, 
	                     description, metadata, created_at, updated_at, expire_time 
	                     FROM orders %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM orders %s`, whereClause)
	
	// 添加分页参数
	queryArgs := append(args, pageSize, offset)
	
	var orders []*Order
	err := r.conn.QueryRowsCtx(ctx, &orders, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	var count int64
	err = r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

// UpdateOrderStatus 更新订单状态
func (r *SqlTradeRepository) UpdateOrderStatus(ctx context.Context, id string, status, reason string, operatorId int64) error {
	// 使用事务确保更新订单状态和记录状态变更历史的一致性
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 更新订单状态
	updateQuery := `UPDATE orders SET status = ?, updated_at = ? WHERE order_id = ?`
	_, err = tx.ExecCtx(ctx, updateQuery, status, time.Now(), id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 记录状态变更历史
	if operatorId > 0 {
		historyQuery := `INSERT INTO order_status_history (order_id, status, reason, operator_id, created_at) 
		               VALUES (?, ?, ?, ?, ?)`
		_, err = tx.ExecCtx(ctx, historyQuery, id, status, reason, operatorId, time.Now())
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetRefundById 根据ID获取退款信息
func (r *SqlTradeRepository) GetRefundById(ctx context.Context, id string) (*Refund, error) {
	var refund Refund
	query := `SELECT refund_id, order_id, user_id, tenant_id, amount, currency, status, 
	          reason, description, processed_by, process_time, created_at, updated_at 
	          FROM refunds WHERE refund_id = ?`
	
	err := r.conn.QueryRowCtx(ctx, &refund, query, id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &refund, nil
}

// GetRefundList 获取退款列表
func (r *SqlTradeRepository) GetRefundList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Refund, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if filters != nil {
		if refundId, ok := filters["refundId"].(string); ok && refundId != "" {
			whereClause += " AND refund_id = ?"
			args = append(args, refundId)
		}
		if orderId, ok := filters["orderId"].(string); ok && orderId != "" {
			whereClause += " AND order_id = ?"
			args = append(args, orderId)
		}
		if userId, ok := filters["userId"].(int64); ok && userId > 0 {
			whereClause += " AND user_id = ?"
			args = append(args, userId)
		}
		if status, ok := filters["status"].(string); ok && status != "" {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if startTime, ok := filters["startTime"].(string); ok && startTime != "" {
			whereClause += " AND created_at >= ?"
			args = append(args, startTime)
		}
		if endTime, ok := filters["endTime"].(string); ok && endTime != "" {
			whereClause += " AND created_at <= ?"
			args = append(args, endTime)
		}
		if tenantId, ok := filters["tenantId"].(int64); ok && tenantId > 0 {
			whereClause += " AND tenant_id = ?"
			args = append(args, tenantId)
		}
	}

	// 执行查询
	query := fmt.Sprintf(`SELECT refund_id, order_id, user_id, tenant_id, amount, currency, status, 
	                     reason, description, processed_by, process_time, created_at, updated_at 
	                     FROM refunds %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM refunds %s`, whereClause)
	
	// 添加分页参数
	queryArgs := append(args, pageSize, offset)
	
	var refunds []*Refund
	err := r.conn.QueryRowsCtx(ctx, &refunds, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	var count int64
	err = r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return refunds, count, nil
}

// ProcessRefund 处理退款请求
func (r *SqlTradeRepository) ProcessRefund(ctx context.Context, id, action, comment string, operatorId int64) error {
	// 使用事务确保处理的一致性
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 首先检查退款状态
	var currentStatus string
	statusQuery := "SELECT status FROM refunds WHERE refund_id = ?"
	err = tx.QueryRowCtx(ctx, &currentStatus, statusQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 根据当前状态和操作类型确定新状态
	var newStatus string
	if action == "approve" {
		newStatus = "approved"
	} else if action == "reject" {
		newStatus = "rejected"
	} else {
		tx.Rollback()
		return fmt.Errorf("无效的操作类型: %s", action)
	}

	// 检查状态转换是否合法
	if currentStatus != "pending" {
		tx.Rollback()
		return fmt.Errorf("只能处理待处理状态的退款，当前状态: %s", currentStatus)
	}

	// 更新退款状态
	now := time.Now()
	updateQuery := `UPDATE refunds SET status = ?, processed_by = ?, process_time = ?, updated_at = ? WHERE refund_id = ?`
	_, err = tx.ExecCtx(ctx, updateQuery, newStatus, fmt.Sprintf("%d", operatorId), now, now, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 记录处理历史
	historyQuery := `INSERT INTO refund_process_history (refund_id, action, status, comment, operator_id, created_at) 
	               VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecCtx(ctx, historyQuery, id, action, newStatus, comment, operatorId, now)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 如果是批准退款，更新相关订单和交易记录
	if newStatus == "approved" {
		// 获取退款信息
		var refund Refund
		refundQuery := `SELECT order_id, amount, currency, user_id, tenant_id FROM refunds WHERE refund_id = ?`
		err = tx.QueryRowCtx(ctx, &refund, refundQuery, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		
		// 记录交易
		transactionQuery := `INSERT INTO transactions (transaction_id, user_id, tenant_id, related_id, type, amount, currency, status, description, created_at) 
		                   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		transactionId := fmt.Sprintf("TR%s", generateRandomString(12))
		description := fmt.Sprintf("退款 %s 的交易记录", id)
		_, err = tx.ExecCtx(ctx, transactionQuery, transactionId, refund.UserId, refund.TenantId, id, "refund", refund.Amount, refund.Currency, "completed", description, now)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetTransactionList 获取交易记录列表
func (r *SqlTradeRepository) GetTransactionList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if filters != nil {
		if userId, ok := filters["userId"].(int64); ok && userId > 0 {
			whereClause += " AND user_id = ?"
			args = append(args, userId)
		}
		if transactionType, ok := filters["type"].(string); ok && transactionType != "" {
			whereClause += " AND type = ?"
			args = append(args, transactionType)
		}
		if status, ok := filters["status"].(string); ok && status != "" {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if startTime, ok := filters["startTime"].(string); ok && startTime != "" {
			whereClause += " AND created_at >= ?"
			args = append(args, startTime)
		}
		if endTime, ok := filters["endTime"].(string); ok && endTime != "" {
			whereClause += " AND created_at <= ?"
			args = append(args, endTime)
		}
		if tenantId, ok := filters["tenantId"].(int64); ok && tenantId > 0 {
			whereClause += " AND tenant_id = ?"
			args = append(args, tenantId)
		}
	}

	// 执行查询
	query := fmt.Sprintf(`SELECT transaction_id, user_id, tenant_id, related_id, type, amount, currency, 
	                     status, description, metadata, created_at 
	                     FROM transactions %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM transactions %s`, whereClause)
	
	// 添加分页参数
	queryArgs := append(args, pageSize, offset)
	
	var transactions []*Transaction
	err := r.conn.QueryRowsCtx(ctx, &transactions, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	var count int64
	err = r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

// GetFinancialReport 获取财务报表
func (r *SqlTradeRepository) GetFinancialReport(ctx context.Context, startTime, endTime string, tenantId int64, groupBy string) ([]*ReportItem, float64, error) {
	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	// 添加时间范围条件
	if startTime != "" {
		whereClause += " AND created_at >= ?"
		args = append(args, startTime)
	}
	if endTime != "" {
		whereClause += " AND created_at <= ?"
		args = append(args, endTime)
	}

	// 添加租户条件
	if tenantId > 0 {
		whereClause += " AND tenant_id = ?"
		args = append(args, tenantId)
	}

	// 根据分组方式确定日期格式化
	var dateFormat string
	switch groupBy {
	case "day":
		dateFormat = "%Y-%m-%d"
	case "week":
		dateFormat = "%Y-%u" // ISO周格式：年份-周数
	case "month":
		dateFormat = "%Y-%m"
	default:
		dateFormat = "%Y-%m-%d"
	}

	// 订单统计查询
	orderQuery := fmt.Sprintf(`
		SELECT 
			DATE_FORMAT(created_at, '%s') AS date,
			COUNT(*) AS order_count,
			SUM(amount) AS order_amount,
			currency AS currency_unit
		FROM orders
		%s
		AND status = 'paid'
		GROUP BY date, currency_unit
	`, dateFormat, whereClause)

	// 退款统计查询
	refundQuery := fmt.Sprintf(`
		SELECT 
			DATE_FORMAT(created_at, '%s') AS date,
			COUNT(*) AS refund_count,
			SUM(amount) AS refund_amount,
			currency AS currency_unit
		FROM refunds
		%s
		AND status = 'approved'
		GROUP BY date, currency_unit
	`, dateFormat, whereClause)

	// 执行查询
	type OrderStats struct {
		Date        string
		OrderCount  int64
		OrderAmount float64
		CurrencyUnit string
	}

	type RefundStats struct {
		Date         string
		RefundCount  int64
		RefundAmount float64
		CurrencyUnit string
	}

	var orderStats []OrderStats
	var refundStats []RefundStats

	err := r.conn.QueryRowsCtx(ctx, &orderStats, orderQuery, args...)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, 0, err
	}

	err = r.conn.QueryRowsCtx(ctx, &refundStats, refundQuery, args...)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, 0, err
	}

	// 合并结果
	resultMap := make(map[string]*ReportItem)
	
	// 处理订单数据
	for _, stat := range orderStats {
		key := stat.Date + "_" + stat.CurrencyUnit
		if item, exists := resultMap[key]; exists {
			item.OrderCount = stat.OrderCount
			item.OrderAmount = stat.OrderAmount
			item.NetAmount = item.OrderAmount - item.RefundAmount
		} else {
			resultMap[key] = &ReportItem{
				Date:        stat.Date,
				OrderCount:  stat.OrderCount,
				OrderAmount: stat.OrderAmount,
				CurrencyUnit: stat.CurrencyUnit,
				NetAmount:   stat.OrderAmount,
			}
		}
	}

	// 处理退款数据
	for _, stat := range refundStats {
		key := stat.Date + "_" + stat.CurrencyUnit
		if item, exists := resultMap[key]; exists {
			item.RefundCount = stat.RefundCount
			item.RefundAmount = stat.RefundAmount
			item.NetAmount = item.OrderAmount - stat.RefundAmount
		} else {
			resultMap[key] = &ReportItem{
				Date:         stat.Date,
				RefundCount:  stat.RefundCount,
				RefundAmount: stat.RefundAmount,
				CurrencyUnit:  stat.CurrencyUnit,
				NetAmount:    -stat.RefundAmount,
			}
		}
	}

	// 转换为切片并计算总净额
	var result []*ReportItem
	var totalAmount float64

	for _, item := range resultMap {
		result = append(result, item)
		totalAmount += item.NetAmount
	}

	// 按日期排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result, totalAmount, nil
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
} 