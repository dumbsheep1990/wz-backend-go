package domain

import "time"

// Theme 主题/模板结构体
type Theme struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	Name        string    `json:"name,omitempty" db:"name"`               // 主题名称
	Code        string    `json:"code,omitempty" db:"code"`               // 主题代码
	Preview     string    `json:"preview,omitempty" db:"preview"`         // 预览图
	Description string    `json:"description,omitempty" db:"description"` // 主题描述
	Status      int       `json:"status,omitempty" db:"status"`           // 状态（0禁用，1启用）
	IsDefault   int       `json:"is_default,omitempty" db:"is_default"`   // 是否默认主题
	Config      string    `json:"config,omitempty" db:"config"`           // 主题配置（JSON格式）
	TenantID    int64     `json:"tenant_id,omitempty" db:"tenant_id"`     // 租户ID，多租户支持
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// ThemeRepository 主题仓储接口
type ThemeRepository interface {
	Create(theme *Theme) (int64, error)
	GetByID(id int64) (*Theme, error)
	List(page, pageSize int, query map[string]interface{}) ([]*Theme, int64, error)
	Update(theme *Theme) error
	Delete(id int64) error
	SetDefault(id int64, tenantID int64) error
	GetDefault(tenantID int64) (*Theme, error)
}
