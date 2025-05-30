package service

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"wz-backend-go/internal/application/site/dto"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/site/entity"
	"wz-backend-go/internal/domain/site/repository"
	siteService "wz-backend-go/internal/domain/site/service"
	"wz-backend-go/internal/domain/site/valueobject"
	"wz-backend-go/internal/infrastructure/database"
	"time"
)

// SiteApplicationService 站点应用服务
type SiteApplicationService struct {
	siteRepo      repository.SiteRepository
	domainService *siteService.SiteDomainService
	eventBus      event.EventBus
	validator     *validator.Validate
	unitOfWork    database.UnitOfWork
}

// NewSiteApplicationService 创建站点应用服务
func NewSiteApplicationService(
	siteRepo repository.SiteRepository,
	domainService *siteService.SiteDomainService,
	eventBus event.EventBus,
	validator *validator.Validate,
	unitOfWork database.UnitOfWork,
) *SiteApplicationService {
	return &SiteApplicationService{
		siteRepo:      siteRepo,
		domainService: domainService,
		eventBus:      eventBus,
		validator:     validator,
		unitOfWork:    unitOfWork,
	}
}

// CreateSite 创建站点
func (s *SiteApplicationService) CreateSite(ctx context.Context, req dto.CreateSiteRequest, tenantID string) (*dto.SiteDTO, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 创建值对象
	siteID, err := valueobject.NewSiteID(generateSiteID())
	if err != nil {
		return nil, fmt.Errorf("创建站点ID失败: %w", err)
	}
	
	siteName, err := valueobject.NewSiteName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("创建站点名称失败: %w", err)
	}
	
	domain, err := valueobject.NewDomain(req.Domain)
	if err != nil {
		return nil, fmt.Errorf("创建域名失败: %w", err)
	}
	
	// 验证域名唯一性
	if err := s.domainService.ValidateUniqueDomain(ctx, domain); err != nil {
		return nil, err
	}
	
	// 处理主题配置
	themeConfig := valueobject.NewDefaultThemeConfig()
	if req.Theme != (dto.ThemeDTO{}) {
		themeConfig, err = valueobject.NewThemeConfig(
			req.Theme.PrimaryColor,
			req.Theme.SecondaryColor,
			req.Theme.AccentColor,
			req.Theme.TextColor,
			req.Theme.BackgroundColor,
			req.Theme.FontFamily,
			req.Theme.HeaderStyle,
			req.Theme.BorderRadius,
			req.Theme.CustomCSS,
		)
		if err != nil {
			return nil, fmt.Errorf("创建主题配置失败: %w", err)
		}
	}
	
	// 创建站点聚合
	site, err := entity.NewSite(siteID, siteName, req.Description, domain, tenantID)
	if err != nil {
		return nil, fmt.Errorf("创建站点失败: %w", err)
	}
	
	// 设置可选属性
	if req.Logo != "" {
		site.UpdateLogo(req.Logo)
	}
	if req.Favicon != "" {
		site.UpdateFavicon(req.Favicon)
	}
	if req.Thumbnail != "" {
		site.UpdateThumbnail(req.Thumbnail)
	}
	site.UpdateTheme(themeConfig)
	
	// 在事务中保存站点
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.siteRepo.Save(ctx, site); err != nil {
			return fmt.Errorf("保存站点失败: %w", err)
		}
		
		// 发布领域事件
		events := site.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		site.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return s.siteToDTO(site), nil
}

// UpdateSite 更新站点
func (s *SiteApplicationService) UpdateSite(ctx context.Context, siteID string, req dto.UpdateSiteRequest, tenantID string) (*dto.SiteDTO, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 创建站点ID
	id, err := valueobject.NewSiteID(siteID)
	if err != nil {
		return nil, fmt.Errorf("无效的站点ID: %w", err)
	}
	
	// 验证站点所有权和修改权限
	site, err := s.domainService.ValidateModification(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	
	// 更新站点信息
	if req.Name != nil {
		siteName, err := valueobject.NewSiteName(*req.Name)
		if err != nil {
			return nil, fmt.Errorf("创建站点名称失败: %w", err)
		}
		if err := site.UpdateName(siteName); err != nil {
			return nil, fmt.Errorf("更新站点名称失败: %w", err)
		}
	}
	
	if req.Description != nil {
		site.UpdateDescription(*req.Description)
	}
	
	if req.Domain != nil {
		domain, err := valueobject.NewDomain(*req.Domain)
		if err != nil {
			return nil, fmt.Errorf("创建域名失败: %w", err)
		}
		
		// 验证域名唯一性（排除当前站点）
		if err := s.domainService.ValidateUniqueDomainForUpdate(ctx, domain, id); err != nil {
			return nil, err
		}
		
		if err := site.UpdateDomain(domain); err != nil {
			return nil, fmt.Errorf("更新域名失败: %w", err)
		}
	}
	
	if req.Logo != nil {
		site.UpdateLogo(*req.Logo)
	}
	
	if req.Favicon != nil {
		site.UpdateFavicon(*req.Favicon)
	}
	
	if req.Thumbnail != nil {
		site.UpdateThumbnail(*req.Thumbnail)
	}
	
	if req.Theme != nil {
		themeConfig, err := valueobject.NewThemeConfig(
			req.Theme.PrimaryColor,
			req.Theme.SecondaryColor,
			req.Theme.AccentColor,
			req.Theme.TextColor,
			req.Theme.BackgroundColor,
			req.Theme.FontFamily,
			req.Theme.HeaderStyle,
			req.Theme.BorderRadius,
			req.Theme.CustomCSS,
		)
		if err != nil {
			return nil, fmt.Errorf("创建主题配置失败: %w", err)
		}
		site.UpdateTheme(themeConfig)
	}
	
	// 在事务中保存站点
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.siteRepo.Save(ctx, site); err != nil {
			return fmt.Errorf("保存站点失败: %w", err)
		}
		
		// 发布领域事件
		events := site.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		site.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return s.siteToDTO(site), nil
}

// GetSite 获取站点详情
func (s *SiteApplicationService) GetSite(ctx context.Context, siteID string, tenantID string) (*dto.SiteDTO, error) {
	id, err := valueobject.NewSiteID(siteID)
	if err != nil {
		return nil, fmt.Errorf("无效的站点ID: %w", err)
	}
	
	site, err := s.siteRepo.FindByIDAndTenant(ctx, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查找站点失败: %w", err)
	}
	if site == nil {
		return nil, fmt.Errorf("站点不存在")
	}
	
	return s.siteToDTO(site), nil
}

// ListSites 获取站点列表
func (s *SiteApplicationService) ListSites(ctx context.Context, req dto.SiteListRequest, tenantID string) (*dto.SiteListResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	
	filters := repository.SiteFilters{
		Status: req.Status,
		Search: req.Search,
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}
	
	sites, err := s.siteRepo.FindByTenant(ctx, tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("查询站点列表失败: %w", err)
	}
	
	total, err := s.siteRepo.Count(ctx, tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("统计站点数量失败: %w", err)
	}
	
	siteDTOs := make([]dto.SiteDTO, 0, len(sites))
	for _, site := range sites {
		siteDTOs = append(siteDTOs, *s.siteToDTO(site))
	}
	
	return &dto.SiteListResponse{
		Sites: siteDTOs,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// PublishSite 发布站点
func (s *SiteApplicationService) PublishSite(ctx context.Context, siteID string, tenantID string) (*dto.PublishSiteResponse, error) {
	id, err := valueobject.NewSiteID(siteID)
	if err != nil {
		return nil, fmt.Errorf("无效的站点ID: %w", err)
	}
	
	// 验证站点所有权
	site, err := s.domainService.ValidateOwnership(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	
	// 发布站点
	if err := site.Publish(); err != nil {
		return nil, fmt.Errorf("发布站点失败: %w", err)
	}
	
	// 在事务中保存站点
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.siteRepo.Save(ctx, site); err != nil {
			return fmt.Errorf("保存站点失败: %w", err)
		}
		
		// 发布领域事件
		events := site.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		site.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &dto.PublishSiteResponse{
		Site:        *s.siteToDTO(site),
		PublishedAt: *site.PublishedAt(),
	}, nil
}

// DeleteSite 删除站点
func (s *SiteApplicationService) DeleteSite(ctx context.Context, siteID string, tenantID string) error {
	id, err := valueobject.NewSiteID(siteID)
	if err != nil {
		return fmt.Errorf("无效的站点ID: %w", err)
	}
	
	// 验证是否可以删除
	if err := s.domainService.CanDeleteSite(ctx, id, tenantID); err != nil {
		return err
	}
	
	// 删除站点
	if err := s.siteRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除站点失败: %w", err)
	}
	
	return nil
}

// siteToDTO 将站点实体转换为DTO
func (s *SiteApplicationService) siteToDTO(site *entity.Site) *dto.SiteDTO {
	theme := site.Theme()
	return &dto.SiteDTO{
		ID:          site.ID().Value(),
		Name:        site.Name().Value(),
		Description: site.Description(),
		Domain:      site.Domain().Value(),
		Logo:        site.Logo(),
		Favicon:     site.Favicon(),
		TenantID:    site.TenantID(),
		Theme: dto.ThemeDTO{
			PrimaryColor:    theme.PrimaryColor(),
			SecondaryColor:  theme.SecondaryColor(),
			AccentColor:     theme.AccentColor(),
			TextColor:       theme.TextColor(),
			BackgroundColor: theme.BackgroundColor(),
			FontFamily:      theme.FontFamily(),
			HeaderStyle:     theme.HeaderStyle(),
			BorderRadius:    theme.BorderRadius(),
			CustomCSS:       theme.CustomCSS(),
		},
		Status:      site.Status().Value(),
		Thumbnail:   site.Thumbnail(),
		CreatedAt:   site.CreatedAt(),
		UpdatedAt:   site.UpdatedAt(),
		PublishedAt: site.PublishedAt(),
	}
}

// generateSiteID 生成站点ID (简单实现，实际应使用UUID)
func generateSiteID() string {
	// 这里应该使用更好的ID生成策略，比如UUID
	return fmt.Sprintf("site_%d", time.Now().UnixNano())
} 