package domain

import "time"

// SiteConfig 站点配置结构体
type SiteConfig struct {
	ID             int64     `json:"id,omitempty" db:"id"`
	SiteName       string    `json:"site_name,omitempty" db:"site_name"`             // 站点名称
	SiteLogo       string    `json:"site_logo,omitempty" db:"site_logo"`             // 站点LOGO
	SeoTitle       string    `json:"seo_title,omitempty" db:"seo_title"`             // SEO标题
	SeoKeywords    string    `json:"seo_keywords,omitempty" db:"seo_keywords"`       // SEO关键词
	SeoDescription string    `json:"seo_description,omitempty" db:"seo_description"` // SEO描述
	IcpNumber      string    `json:"icp_number,omitempty" db:"icp_number"`           // ICP备案号
	Copyright      string    `json:"copyright,omitempty" db:"copyright"`             // 版权信息
	ThemeID        int64     `json:"theme_id,omitempty" db:"theme_id"`               // 当前主题ID
	ContactEmail   string    `json:"contact_email,omitempty" db:"contact_email"`     // 联系邮箱
	ContactPhone   string    `json:"contact_phone,omitempty" db:"contact_phone"`     // 联系电话
	Address        string    `json:"address,omitempty" db:"address"`                 // 联系地址
	TenantID       int64     `json:"tenant_id,omitempty" db:"tenant_id"`             // 租户ID，多租户支持
	CreatedAt      time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// SiteConfigRepository 站点配置仓储接口
type SiteConfigRepository interface {
	GetSiteConfig(tenantID int64) (*SiteConfig, error)
	UpdateSiteConfig(config *SiteConfig) error
}
