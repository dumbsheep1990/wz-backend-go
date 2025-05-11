package sql

import (
	"time"

	"wz-backend-go/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// SiteConfigRepository 站点配置SQL仓储实现
type SiteConfigRepository struct {
	conn sqlx.SqlConn
}

// NewSiteConfigRepository 创建站点配置仓储实现
func NewSiteConfigRepository(conn sqlx.SqlConn) *SiteConfigRepository {
	return &SiteConfigRepository{
		conn: conn,
	}
}

// GetSiteConfig 获取站点配置
func (r *SiteConfigRepository) GetSiteConfig(tenantID int64) (*domain.SiteConfig, error) {
	var config domain.SiteConfig
	query := `
		SELECT 
			id, site_name, site_logo, seo_title, seo_keywords, 
			seo_description, icp_number, copyright, theme_id, 
			contact_email, contact_phone, address, tenant_id, 
			created_at, updated_at
		FROM site_configs 
		WHERE tenant_id = ?
		LIMIT 1
	`

	err := r.conn.QueryRow(&config, query, tenantID)
	if err != nil {
		// 如果没有找到记录，创建一个默认配置
		if err == sqlx.ErrNotFound {
			logx.Infof("未找到租户(%d)的站点配置，将创建默认配置", tenantID)
			defaultConfig := &domain.SiteConfig{
				SiteName:       "默认站点名称",
				SeoTitle:       "默认SEO标题",
				SeoKeywords:    "默认关键词",
				SeoDescription: "默认网站描述",
				Copyright:      "版权所有 © " + time.Now().Format("2006"),
				TenantID:       tenantID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// 插入默认配置
			id, insertErr := r.createSiteConfig(defaultConfig)
			if insertErr != nil {
				logx.Errorf("创建默认站点配置失败: %v", insertErr)
				return nil, insertErr
			}

			defaultConfig.ID = id
			return defaultConfig, nil
		}

		logx.Errorf("获取站点配置失败: %v, tenantID: %d", err, tenantID)
		return nil, err
	}

	return &config, nil
}

// createSiteConfig 创建站点配置（内部方法）
func (r *SiteConfigRepository) createSiteConfig(config *domain.SiteConfig) (int64, error) {
	query := `
		INSERT INTO site_configs (
			site_name, site_logo, seo_title, seo_keywords, seo_description,
			icp_number, copyright, theme_id, contact_email, contact_phone,
			address, tenant_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	result, err := r.conn.Exec(query,
		config.SiteName, config.SiteLogo, config.SeoTitle, config.SeoKeywords,
		config.SeoDescription, config.IcpNumber, config.Copyright, config.ThemeID,
		config.ContactEmail, config.ContactPhone, config.Address, config.TenantID,
		config.CreatedAt, config.UpdatedAt,
	)
	if err != nil {
		logx.Errorf("创建站点配置失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取新创建站点配置ID失败: %v", err)
		return 0, err
	}

	return id, nil
}

// UpdateSiteConfig 更新站点配置
func (r *SiteConfigRepository) UpdateSiteConfig(config *domain.SiteConfig) error {
	query := `
		UPDATE site_configs SET
			site_name = ?, site_logo = ?, seo_title = ?, seo_keywords = ?,
			seo_description = ?, icp_number = ?, copyright = ?, theme_id = ?,
			contact_email = ?, contact_phone = ?, address = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`
	now := time.Now()
	config.UpdatedAt = now

	_, err := r.conn.Exec(query,
		config.SiteName, config.SiteLogo, config.SeoTitle, config.SeoKeywords,
		config.SeoDescription, config.IcpNumber, config.Copyright, config.ThemeID,
		config.ContactEmail, config.ContactPhone, config.Address, now,
		config.ID, config.TenantID,
	)
	if err != nil {
		logx.Errorf("更新站点配置失败: %v, id: %d", err, config.ID)
		return err
	}

	return nil
}
