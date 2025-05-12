package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain"
)

// UserMessageRepository 用户消息仓储SQL实现
type UserMessageRepository struct {
	db *sqlx.DB
}

// NewUserMessageRepository 创建用户消息仓储实例
func NewUserMessageRepository(db *sqlx.DB) domain.UserMessageRepository {
	return &UserMessageRepository{
		db: db,
	}
}

// Create 创建用户消息
func (r *UserMessageRepository) Create(message *domain.UserMessage) (int64, error) {
	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now
	// 默认为未读
	if message.Status == 0 {
		message.Status = 0
	}

	query := `INSERT INTO user_messages (
        user_id, title, content, type, status, 
        is_important, related_id, related_type, tenant_id, 
        created_at, updated_at
    ) VALUES (
        :user_id, :title, :content, :type, :status, 
        :is_important, :related_id, :related_type, :tenant_id, 
        :created_at, :updated_at
    )`

	result, err := r.db.NamedExec(query, message)
	if err != nil {
		return 0, fmt.Errorf("创建用户消息失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}

	return id, nil
}

// GetByID 根据ID获取用户消息
func (r *UserMessageRepository) GetByID(id int64) (*domain.UserMessage, error) {
	var message domain.UserMessage
	query := `SELECT * FROM user_messages WHERE id = ?`

	err := r.db.Get(&message, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("获取用户消息失败: %w", err)
	}

	return &message, nil
}

// ListByUser 获取用户消息列表
func (r *UserMessageRepository) ListByUser(userID int64, page, pageSize int, query map[string]interface{}) ([]*domain.UserMessage, int64, error) {
	messages := []*domain.UserMessage{}
	var count int64

	// 构建查询条件
	conditions := []string{"user_id = ?"}
	args := []interface{}{userID}

	if status, ok := query["status"]; ok {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	if msgType, ok := query["type"]; ok {
		conditions = append(conditions, "type = ?")
		args = append(args, msgType)
	}

	if important, ok := query["is_important"]; ok {
		conditions = append(conditions, "is_important = ?")
		args = append(args, important)
	}

	if tenantID, ok := query["tenant_id"]; ok {
		conditions = append(conditions, "tenant_id = ?")
		args = append(args, tenantID)
	}

	if title, ok := query["title"]; ok && title != "" {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+title.(string)+"%")
	}

	// 构建WHERE子句
	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM user_messages %s", whereClause)
	err := r.db.Get(&count, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("统计用户消息数量失败: %w", err)
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf("SELECT * FROM user_messages %s ORDER BY is_important DESC, id DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, pageSize, offset)

	err = r.db.Select(&messages, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户消息列表失败: %w", err)
	}

	return messages, count, nil
}

// MarkAsRead 标记消息为已读
func (r *UserMessageRepository) MarkAsRead(id int64, userID int64) error {
	query := `UPDATE user_messages SET status = 1, updated_at = ? WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("标记消息已读失败: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("消息不存在或不属于该用户")
	}

	return nil
}

// MarkAllAsRead 标记所有消息为已读
func (r *UserMessageRepository) MarkAllAsRead(userID int64) error {
	query := `UPDATE user_messages SET status = 1, updated_at = ? WHERE user_id = ? AND status = 0`
	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("标记所有消息已读失败: %w", err)
	}

	return nil
}

// Delete 删除用户消息
func (r *UserMessageRepository) Delete(id int64, userID int64) error {
	query := `DELETE FROM user_messages WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("删除用户消息失败: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("消息不存在或不属于该用户")
	}

	return nil
}

// CountUnread 获取未读消息数量
func (r *UserMessageRepository) CountUnread(userID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM user_messages WHERE user_id = ? AND status = 0`
	err := r.db.Get(&count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("统计未读消息数量失败: %w", err)
	}

	return count, nil
}
