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
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UserFavoriteRepository 用户收藏仓储接口
type UserFavoriteRepository interface {
	Create(favorite *UserFavorite) (int64, error)
	GetByID(id int64) (*UserFavorite, error)
	ListByUser(userID int64, page, pageSize int, itemType string) ([]*UserFavorite, int64, error)
	Delete(id int64, userID int64) error
	CheckFavorite(userID int64, itemID int64, itemType string) (bool, error)
}
