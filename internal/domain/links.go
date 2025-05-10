package domain

import "time"

// Links 友情链接结构体
type Link struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	Name        string    `json:"name,omitempty" db:"name"`               // 链接名称
	URL         string    `json:"url,omitempty" db:"url"`                 // 链接URL
	Logo        string    `json:"logo,omitempty" db:"logo"`               // 链接Logo
	Sort        int       `json:"sort,omitempty" db:"sort"`               // 排序
	Status      int       `json:"status,omitempty" db:"status"`           // 状态（0禁用，1启用）
	Description string    `json:"description,omitempty" db:"description"` // 链接描述
	TenantID    int64     `json:"tenant_id,omitempty" db:"tenant_id"`     // 租户ID，多租户支持
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// LinkRepository 链接仓储接口，定义友情链接相关的数据库操作
type LinkRepository interface {
	Create(link *Link) (int64, error)
	GetByID(id int64) (*Link, error)
	List(page, pageSize int, query map[string]interface{}) ([]*Link, int64, error)
	Update(link *Link) error
	Delete(id int64) error
}
