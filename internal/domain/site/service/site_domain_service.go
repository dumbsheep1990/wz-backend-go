package service

import (
	"context"
	"errors"
	"fmt"
	"wz-backend-go/internal/domain/site/entity"
	"wz-backend-go/internal/domain/site/repository"
	"wz-backend-go/internal/domain/site/valueobject"
)

// SiteDomainService 站点领域服务
type SiteDomainService struct {
	siteRepo repository.SiteRepository
}

// NewSiteDomainService 创建站点领域服务
func NewSiteDomainService(siteRepo repository.SiteRepository) *SiteDomainService {
	return &SiteDomainService{
		siteRepo: siteRepo,
	}
}

// ValidateUniqueDomain 验证域名唯一性
func (s *SiteDomainService) ValidateUniqueDomain(ctx context.Context, domain valueobject.Domain) error {
	exists, err := s.siteRepo.ExistsByDomain(ctx, domain)
	if err != nil {
		return fmt.Errorf("检查域名唯一性失败: %w", err)
	}
	if exists {
		return errors.New("域名已被使用")
	}
	return nil
}

// ValidateUniqueDomainForUpdate 验证更新时的域名唯一性
func (s *SiteDomainService) ValidateUniqueDomainForUpdate(ctx context.Context, domain valueobject.Domain, siteID valueobject.SiteID) error {
	exists, err := s.siteRepo.ExistsByDomainExcludeID(ctx, domain, siteID)
	if err != nil {
		return fmt.Errorf("检查域名唯一性失败: %w", err)
	}
	if exists {
		return errors.New("域名已被使用")
	}
	return nil
}

// CanDeleteSite 检查是否可以删除站点
func (s *SiteDomainService) CanDeleteSite(ctx context.Context, siteID valueobject.SiteID, tenantID string) error {
	site, err := s.siteRepo.FindByIDAndTenant(ctx, siteID, tenantID)
	if err != nil {
		return fmt.Errorf("查找站点失败: %w", err)
	}
	if site == nil {
		return errors.New("站点不存在或无权限")
	}
	
	if !site.CanBeDeleted() {
		return errors.New("只有草稿状态的站点才能被删除")
	}
	
	return nil
}

// ValidateOwnership 验证站点所有权
func (s *SiteDomainService) ValidateOwnership(ctx context.Context, siteID valueobject.SiteID, tenantID string) (*entity.Site, error) {
	site, err := s.siteRepo.FindByIDAndTenant(ctx, siteID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查找站点失败: %w", err)
	}
	if site == nil {
		return nil, errors.New("站点不存在或无权限")
	}
	
	if !site.IsOwnedBy(tenantID) {
		return nil, errors.New("无权限操作此站点")
	}
	
	return site, nil
}

// ValidateModification 验证是否可以修改站点
func (s *SiteDomainService) ValidateModification(ctx context.Context, siteID valueobject.SiteID, tenantID string) (*entity.Site, error) {
	site, err := s.ValidateOwnership(ctx, siteID, tenantID)
	if err != nil {
		return nil, err
	}
	
	if !site.CanBeModified() {
		return nil, errors.New("归档状态的站点不能被修改")
	}
	
	return site, nil
}

// GenerateUniqueDomain 生成唯一域名
func (s *SiteDomainService) GenerateUniqueDomain(ctx context.Context, baseDomain string) (valueobject.Domain, error) {
	// 尝试原始域名
	domain, err := valueobject.NewDomain(baseDomain)
	if err != nil {
		return valueobject.Domain{}, err
	}
	
	exists, err := s.siteRepo.ExistsByDomain(ctx, domain)
	if err != nil {
		return valueobject.Domain{}, err
	}
	if !exists {
		return domain, nil
	}
	
	// 尝试添加数字后缀
	for i := 1; i <= 100; i++ {
		candidateDomain := fmt.Sprintf("%s-%d", baseDomain, i)
		domain, err := valueobject.NewDomain(candidateDomain)
		if err != nil {
			continue
		}
		
		exists, err := s.siteRepo.ExistsByDomain(ctx, domain)
		if err != nil {
			return valueobject.Domain{}, err
		}
		if !exists {
			return domain, nil
		}
	}
	
	return valueobject.Domain{}, errors.New("无法生成唯一域名")
} 