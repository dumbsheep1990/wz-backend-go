package domain

import "time"

// UserPoints 用户积分结构体
type UserPoints struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	UserID      int64     `json:"user_id,omitempty" db:"user_id"`           // 用户ID
	Points      int       `json:"points,omitempty" db:"points"`             // 积分变动值
	TotalPoints int       `json:"total_points,omitempty" db:"total_points"` // 总积分（冗余字段）
	Type        int       `json:"type,omitempty" db:"type"`                 // 类型（1增加，2减少）
	Source      string    `json:"source,omitempty" db:"source"`             // 来源（如签到、消费、活动）
	Description string    `json:"description,omitempty" db:"description"`   // 描述
	RelatedID   int64     `json:"related_id,omitempty" db:"related_id"`     // 关联ID
	RelatedType string    `json:"related_type,omitempty" db:"related_type"` // 关联类型
	TenantID    int64     `json:"tenant_id,omitempty" db:"tenant_id"`       // 租户ID，多租户支持
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UserPointsRepository 用户积分仓储接口
type UserPointsRepository interface {
	Create(points *UserPoints) (int64, error)
	GetByID(id int64) (*UserPoints, error)
	ListByUser(userID int64, page, pageSize int) ([]*UserPoints, int64, error)
	GetUserTotalPoints(userID int64) (int, error)
}
