package sql

import (
	"fmt"
	"time"

	"wz-backend-go/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserMessageRepository 用户消息SQL仓储实现
type UserMessageRepository struct {
	conn sqlx.SqlConn
}

// NewUserMessageRepository 创建用户消息仓储实现
func NewUserMessageRepository(conn sqlx.SqlConn) *UserMessageRepository {
	return &UserMessageRepository{
		conn: conn,
	}
}

// Create 创建用户消息
func (r *UserMessageRepository) Create(message *domain.UserMessage) (int64, error) {
	query := `
		INSERT INTO user_messages (
			user_id, title, content, type, status, is_important,
			related_id, related_type, tenant_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	result, err := r.conn.Exec(query,
		message.UserID, message.Title, message.Content, message.Type,
		message.Status, message.IsImportant, message.RelatedID,
		message.RelatedType, message.TenantID, message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		logx.Errorf("创建用户消息失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取新创建用户消息ID失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetByID 根据ID获取用户消息
func (r *UserMessageRepository) GetByID(id int64) (*domain.UserMessage, error) {
	var message domain.UserMessage
	query := `SELECT 
		id, user_id, title, content, type, status, is_important, 
		related_id, related_type, tenant_id, created_at, updated_at
	FROM user_messages 
	WHERE id = ?`

	err := r.conn.QueryRow(&message, query, id)
	if err != nil {
		logx.Errorf("根据ID查询用户消息失败: %v, id: %d", err, id)
		return nil, err
	}

	return &message, nil
}

// ListByUser 获取用户消息列表
func (r *UserMessageRepository) ListByUser(userID int64, page, pageSize int, query map[string]interface{}) ([]*domain.UserMessage, int64, error) {
	// 构建查询条件
	whereClause := "user_id = ?"
	args := []interface{}{userID}

	// 动态添加查询条件
	if query != nil {
		if status, ok := query["status"].(int); ok {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if messageType, ok := query["type"].(int); ok {
			whereClause += " AND type = ?"
			args = append(args, messageType)
		}
		if isImportant, ok := query["is_important"].(int); ok {
			whereClause += " AND is_important = ?"
			args = append(args, isImportant)
		}
		if title, ok := query["title"].(string); ok && title != "" {
			whereClause += " AND title LIKE ?"
			args = append(args, fmt.Sprintf("%%%s%%", title))
		}
		if relatedID, ok := query["related_id"].(int64); ok && relatedID > 0 {
			whereClause += " AND related_id = ?"
			args = append(args, relatedID)
		}
		if relatedType, ok := query["related_type"].(string); ok && relatedType != "" {
			whereClause += " AND related_type = ?"
			args = append(args, relatedType)
		}
		if tenantID, ok := query["tenant_id"].(int64); ok && tenantID > 0 {
			whereClause += " AND tenant_id = ?"
			args = append(args, tenantID)
		}
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM user_messages WHERE %s", whereClause)
	var count int64
	err := r.conn.QueryRow(&count, countQuery, args...)
	if err != nil {
		logx.Errorf("查询用户消息总数失败: %v", err)
		return nil, 0, err
	}

	// 计算分页
	offset := (page - 1) * pageSize

	// 查询列表
	listQuery := fmt.Sprintf(`
		SELECT 
			id, user_id, title, content, type, status, is_important, 
			related_id, related_type, tenant_id, created_at, updated_at
		FROM user_messages 
		WHERE %s 
		ORDER BY is_important DESC, created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)
	var messages []*domain.UserMessage
	err = r.conn.QueryRows(&messages, listQuery, args...)
	if err != nil {
		logx.Errorf("查询用户消息列表失败: %v", err)
		return nil, 0, err
	}

	return messages, count, nil
}

// MarkAsRead 标记消息为已读
func (r *UserMessageRepository) MarkAsRead(id int64, userID int64) error {
	query := "UPDATE user_messages SET status = 1, updated_at = ? WHERE id = ? AND user_id = ?"
	now := time.Now()
	_, err := r.conn.Exec(query, now, id, userID)
	if err != nil {
		logx.Errorf("标记消息为已读失败: %v, id: %d, userID: %d", err, id, userID)
		return err
	}

	return nil
}

// MarkAllAsRead 标记用户所有消息为已读
func (r *UserMessageRepository) MarkAllAsRead(userID int64) error {
	query := "UPDATE user_messages SET status = 1, updated_at = ? WHERE user_id = ? AND status = 0"
	now := time.Now()
	_, err := r.conn.Exec(query, now, userID)
	if err != nil {
		logx.Errorf("标记所有消息为已读失败: %v, userID: %d", err, userID)
		return err
	}

	return nil
}

// Delete 删除用户消息
func (r *UserMessageRepository) Delete(id int64, userID int64) error {
	query := "DELETE FROM user_messages WHERE id = ? AND user_id = ?"
	_, err := r.conn.Exec(query, id, userID)
	if err != nil {
		logx.Errorf("删除用户消息失败: %v, id: %d, userID: %d", err, id, userID)
		return err
	}

	return nil
}

// CountUnread 统计用户未读消息数量
func (r *UserMessageRepository) CountUnread(userID int64) (int64, error) {
	query := "SELECT COUNT(*) FROM user_messages WHERE user_id = ? AND status = 0"
	var count int64
	err := r.conn.QueryRow(&count, query, userID)
	if err != nil {
		logx.Errorf("统计用户未读消息数量失败: %v, userID: %d", err, userID)
		return 0, err
	}

	return count, nil
}
