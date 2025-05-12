package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain"
)

// SiteConfigRepository SQL实现
type SiteConfigRepository struct {
	db *sqlx.DB
}

// NewSiteConfigRepository 创建站点配置仓储实例
func NewSiteConfigRepository(db *sqlx.DB) domain.SiteConfigRepository {
	return &SiteConfigRepository{
		db: db,
	}
}

// GetSiteConfig 获取站点配置
func (r *SiteConfigRepository) GetSiteConfig(tenantID int64) (*domain.SiteConfig, error) {
	var config domain.SiteConfig

	query := `SELECT * FROM site_configs WHERE tenant_id = ? LIMIT 1`
	err := r.db.Get(&config, query, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有配置，返回默认配置
			return &domain.SiteConfig{
				SiteName:    "默认站点",
				SeoTitle:    "默认站点",
				SeoKeywords: "默认站点,默认关键词",
				TenantID:    tenantID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("获取站点配置失败: %w", err)
	}

	return &config, nil
}

// UpdateSiteConfig 更新站点配置
func (r *SiteConfigRepository) UpdateSiteConfig(config *domain.SiteConfig) error {
	var exists bool
	err := r.db.Get(&exists, "SELECT 1 FROM site_configs WHERE tenant_id = ? LIMIT 1", config.TenantID)

	config.UpdatedAt = time.Now()

	if err == sql.ErrNoRows {
		// 如果配置不存在，则创建
		config.CreatedAt = time.Now()
		query := `INSERT INTO site_configs (
            site_name, site_logo, seo_title, seo_keywords, seo_description, 
            icp_number, copyright, theme_id, contact_email, contact_phone, 
            address, tenant_id, created_at, updated_at
        ) VALUES (
            :site_name, :site_logo, :seo_title, :seo_keywords, :seo_description, 
            :icp_number, :copyright, :theme_id, :contact_email, :contact_phone, 
            :address, :tenant_id, :created_at, :updated_at
        )`
		_, err := r.db.NamedExec(query, config)
		if err != nil {
			return fmt.Errorf("创建站点配置失败: %w", err)
		}
	} else {
		// 如果配置存在，则更新
		query := `UPDATE site_configs SET 
            site_name = :site_name,
            site_logo = :site_logo,
            seo_title = :seo_title,
            seo_keywords = :seo_keywords,
            seo_description = :seo_description,
            icp_number = :icp_number,
            copyright = :copyright,
            theme_id = :theme_id,
            contact_email = :contact_email,
            contact_phone = :contact_phone,
            address = :address,
            updated_at = :updated_at
        WHERE tenant_id = :tenant_id`

		_, err := r.db.NamedExec(query, config)
		if err != nil {
			return fmt.Errorf("更新站点配置失败: %w", err)
		}
	}

	return nil
}
