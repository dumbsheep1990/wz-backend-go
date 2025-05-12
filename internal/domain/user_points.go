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
	OperatorID  int64     `json:"operator_id,omitempty" db:"operator_id"`   // 操作员ID，管理员操作时使用
	Username    string    `json:"username,omitempty" db:"-"`                // 用户名称（非数据库字段）
	Operator    string    `json:"operator,omitempty" db:"-"`                // 操作员名称（非数据库字段）
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// PointsRules 积分规则结构体
type PointsRules struct {
	ID                int64     `json:"id,omitempty" db:"id"`
	SignInPoints      int       `json:"sign_in_points" db:"sign_in_points"`           // 签到积分
	CommentPoints     int       `json:"comment_points" db:"comment_points"`           // 评论积分
	SharePoints       int       `json:"share_points" db:"share_points"`               // 分享积分
	ArticlePoints     int       `json:"article_points" db:"article_points"`           // 发布文章积分
	InvitePoints      int       `json:"invite_points" db:"invite_points"`             // 邀请积分
	PurchaseRate      int       `json:"purchase_rate" db:"purchase_rate"`             // 购买积分比例
	MaxDailyPoints    int       `json:"max_daily_points" db:"max_daily_points"`       // 每日最大获取积分
	EnableExchange    bool      `json:"enable_exchange" db:"enable_exchange"`         // 是否可兑换商品
	ExchangeRate      int       `json:"exchange_rate" db:"exchange_rate"`             // 兑换比例
	MinExchangePoints int       `json:"min_exchange_points" db:"min_exchange_points"` // 最小兑换积分
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`                   // 更新时间
}

// SourceItem 来源数量统计项
type SourceItem struct {
	Source string `json:"source" db:"source"`
	Count  int64  `json:"count" db:"count"`
}

// UserPointsRepository 用户积分仓储接口
type UserPointsRepository interface {
	Create(points *UserPoints) (int64, error)
	GetByID(id int64) (*UserPoints, error)
	GetTotalPointsByUserID(userID int64) (int, error)
	ListByUserID(userID int64, offset, limit int64) ([]*UserPoints, error)
	CountByUserID(userID int64) (int64, error)
	ListWithConditions(conditions map[string]interface{}, offset, limit int64) ([]*UserPoints, error)
	CountWithConditions(conditions map[string]interface{}) (int64, error)
	MarkAsRevoked(id int64) error

	// 统计相关方法
	CountUsers() (int64, error)
	SumPoints() (int64, error)
	MaxPoints() (int64, error)
	SumPointsByConditions(conditions map[string]interface{}) (int64, error)
	GroupBySource() ([]*SourceItem, error)

	// 积分规则相关方法
	GetPointsRules() (*PointsRules, error)
	UpdatePointsRules(rules *PointsRules) error
}
