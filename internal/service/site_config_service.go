package service

import (
	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// SiteConfigServiceImpl 站点配置服务实现
type SiteConfigServiceImpl struct {
	repo domain.SiteConfigRepository
}

// NewSiteConfigService 创建站点配置服务
func NewSiteConfigService(repo domain.SiteConfigRepository) SiteConfigService {
	return &SiteConfigServiceImpl{
		repo: repo,
	}
}

// 通过依赖注入SQL连接创建服务
func NewSiteConfigServiceWithConn(conn interface{}) SiteConfigService {
	return &SiteConfigServiceImpl{
		repo: sql.NewSiteConfigRepository(conn),
	}
}

// GetSiteConfig 获取站点配置
func (s *SiteConfigServiceImpl) GetSiteConfig(tenantID int64) (*domain.SiteConfig, error) {
	logx.Infof("获取站点配置: tenantID=%d", tenantID)

	// 参数验证
	if tenantID <= 0 {
		logx.Error("无效的租户ID")
		return nil, domain.ErrInvalidParam
	}

	// 调用仓储层获取站点配置
	config, err := s.repo.GetSiteConfig(tenantID)
	if err != nil {
		logx.Errorf("获取站点配置失败: %v, tenantID: %d", err, tenantID)
		return nil, err
	}

	return config, nil
}

// UpdateSiteConfig 更新站点配置
func (s *SiteConfigServiceImpl) UpdateSiteConfig(config *domain.SiteConfig) error {
	logx.Infof("更新站点配置: %+v", config)

	// 参数验证
	if config.TenantID <= 0 {
		logx.Error("无效的租户ID")
		return domain.ErrInvalidParam
	}

	// 验证配置数据
	if err := s.validateSiteConfig(config); err != nil {
		return err
	}

	// 调用仓储层更新站点配置
	err := s.repo.UpdateSiteConfig(config)
	if err != nil {
		logx.Errorf("更新站点配置失败: %v", err)
		return err
	}

	return nil
}

// validateSiteConfig 验证站点配置
func (s *SiteConfigServiceImpl) validateSiteConfig(config *domain.SiteConfig) error {
	// 这里可以添加业务规则验证
	// 例如：站点名称不能为空、邮箱格式验证等
	return nil
}
