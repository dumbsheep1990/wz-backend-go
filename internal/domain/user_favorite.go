package domain

import "time"

// UserFavorite 用户收藏结构体
type UserFavorite struct {
	ID        int64     `json:"id,omitempty" db:"id"`
	UserID    int64     `json:"user_id,omitempty" db:"user_id"`     // 用户ID
	ItemID    int64     `json:"item_id,omitempty" db:"item_id"`     // 收藏项目ID
	ItemType  string    `json:"item_type,omitempty" db:"item_type"` // 收藏项目类型（如文章、商品）
	Title     string    `json:"title,omitempty" db:"title"`         // 收藏项目标题
	Cover     string    `json:"cover,omitempty" db:"cover"`         // 收藏项目封面
	Summary   string    `json:"summary,omitempty" db:"summary"`     // 收藏项目摘要
	URL       string    `json:"url,omitempty" db:"url"`             // 收藏项目URL
	Remark    string    `json:"remark,omitempty" db:"remark"`       // 用户备注
	TenantID  int64     `json:"tenant_id,omitempty" db:"tenant_id"` // 租户ID，多租户支持
	Username  string    `json:"username,omitempty" db:"-"`          // 用户名称（非数据库字段）
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// TypeDistributionItem 收藏类型分布项
type TypeDistributionItem struct {
	Type  string `json:"type" db:"type"`
	Count int64  `json:"count" db:"count"`
}

// HotContentItem 热门收藏内容项
type HotContentItem struct {
	ItemID     int64  `json:"item_id" db:"item_id"`
	ItemType   string `json:"item_type" db:"item_type"`
	Title      string `json:"title" db:"title"`
	Cover      string `json:"cover" db:"cover"`
	Count      int64  `json:"count" db:"count"`
	CreateDate string `json:"create_date" db:"create_date"`
}

// TrendItem 趋势数据项
type TrendItem struct {
	Date  string `json:"date" db:"date"`
	Count int64  `json:"count" db:"count"`
}

// UserFavoriteRepository 用户收藏仓储接口
type UserFavoriteRepository interface {
	Create(favorite *UserFavorite) (int64, error)
	GetByID(id int64) (*UserFavorite, error)
	ListByUserID(userID int64, offset, limit int64, itemType string) ([]*UserFavorite, error)
	CountByUserID(userID int64, itemType string) (int64, error)
	ListWithConditions(conditions map[string]interface{}, offset, limit int64) ([]*UserFavorite, error)
	CountWithConditions(conditions map[string]interface{}) (int64, error)
	DeleteByID(id int64) error
	BatchDelete(ids []int64) error
	CheckFavorite(userID int64, itemID int64, itemType string) (bool, error)

	// 统计相关方法
	CountUsers() (int64, error)
	CountFavorites() (int64, error)
	CountTodayFavorites() (int64, error)
	CountMonthFavorites() (int64, error)
	GroupByType() ([]*TypeDistributionItem, error)
	GetHotContent(limit int) ([]*HotContentItem, error)
	GetTrend(period string) ([]*TrendItem, error)
}
