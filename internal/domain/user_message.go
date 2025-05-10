package domain

import "time"

// UserMessage 用户消息结构体
type UserMessage struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	UserID      int64     `json:"user_id,omitempty" db:"user_id"`           // 用户ID
	Title       string    `json:"title,omitempty" db:"title"`               // 消息标题
	Content     string    `json:"content,omitempty" db:"content"`           // 消息内容
	Type        int       `json:"type,omitempty" db:"type"`                 // 消息类型（1系统通知，2业务通知，3私信）
	Status      int       `json:"status,omitempty" db:"status"`             // 状态（0未读，1已读）
	IsImportant int       `json:"is_important,omitempty" db:"is_important"` // 是否重要
	RelatedID   int64     `json:"related_id,omitempty" db:"related_id"`     // 关联ID
	RelatedType string    `json:"related_type,omitempty" db:"related_type"` // 关联类型
	TenantID    int64     `json:"tenant_id,omitempty" db:"tenant_id"`       // 租户ID，多租户支持
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UserMessageRepository 用户消息仓储接口
type UserMessageRepository interface {
	Create(message *UserMessage) (int64, error)
	GetByID(id int64) (*UserMessage, error)
	ListByUser(userID int64, page, pageSize int, query map[string]interface{}) ([]*UserMessage, int64, error)
	MarkAsRead(id int64, userID int64) error
	MarkAllAsRead(userID int64) error
	Delete(id int64, userID int64) error
	CountUnread(userID int64) (int64, error)
}
