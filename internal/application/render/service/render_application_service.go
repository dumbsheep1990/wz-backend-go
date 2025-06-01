package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
	"wz-backend-go/internal/application/render/dto"
	domainService "wz-backend-go/internal/domain/render/service"
	"wz-backend-go/internal/domain/render/valueobject"
)

// RenderApplicationService 渲染应用服务
type RenderApplicationService struct {
	renderDomainService *domainService.RenderDomainService
	// 其他依赖服务，如站点服务、页面服务等
	siteService         SiteService
	pageService         PageService
}

// SiteService 站点服务接口
type SiteService interface {
	GetSiteByID(ctx context.Context, siteID string) (interface{}, error)
	GetSiteByDomain(ctx context.Context, domain string) (interface{}, error)
	IsSitePublished(ctx context.Context, siteID string) bool
}

// PageService 页面服务接口
type PageService interface {
	GetPageByID(ctx context.Context, pageID string) (interface{}, error)
	GetPageBySlug(ctx context.Context, siteID string, slug string) (interface{}, error)
	GetHomePage(ctx context.Context, siteID string) (interface{}, error)
}

// NewRenderApplicationService 创建一个新的渲染应用服务
func NewRenderApplicationService(
	renderDomainService *domainService.RenderDomainService,
	siteService SiteService,
	pageService PageService,
) *RenderApplicationService {
	return &RenderApplicationService{
		renderDomainService: renderDomainService,
		siteService:         siteService,
		pageService:         pageService,
	}
}

// RenderTemplate 渲染模板
func (s *RenderApplicationService) RenderTemplate(
	ctx context.Context,
	request dto.RenderRequestDTO,
) (*dto.RenderResponseDTO, error) {
	// 验证请求
	if request.TemplateID == "" {
		return nil, errors.New("模板ID不能为空")
	}

	// 检查站点是否存在
	if request.SiteID != "" {
		_, err := s.siteService.GetSiteByID(ctx, request.SiteID)
		if err != nil {
			return nil, errors.New("站点不存在: " + err.Error())
		}
	}

	// 构造模板上下文
	templateContext := valueobject.NewTemplateContext(request.Data, request.Metadata)

	// 获取渲染格式
	format, err := valueobject.FormatFromString(request.Format)
	if err != nil {
		return nil, err
	}

	// 构造缓存选项
	cacheOptions := &domainService.CacheOptions{
		Enabled: request.CacheConfig.Enabled,
		TTL:     time.Duration(request.CacheConfig.TTL) * time.Second,
	}

	// 设置缓存级别
	if request.CacheConfig.Level != "" {
		cacheOptions.Level = valueobject.CacheLevel(request.CacheConfig.Level)
	} else {
		cacheOptions.Level = valueobject.CacheLevelMemory
	}

	// 调用领域服务进行渲染
	result, err := s.renderDomainService.RenderTemplate(
		ctx,
		request.TemplateID,
		templateContext,
		format,
		cacheOptions,
	)

	if err != nil {
		return nil, err
	}

	// 构造响应
	response := &dto.RenderResponseDTO{
		ID:          result.ID().String(),
		Content:     result.Content(),
		ContentType: result.Format().ContentType(),
		Format:      result.Format().Format(),
		CreatedAt:   result.CreatedAt(),
		CacheHit:    false, // 这里简化了缓存命中的判断
	}

	// 如果启用了缓存，设置过期时间
	if result.CacheStrategy().Enabled() {
		response.ExpiresAt = result.ExpiresAt()
	}

	return response, nil
}

// RenderPage 渲染页面
func (s *RenderApplicationService) RenderPage(
	ctx context.Context,
	request dto.PageRenderRequestDTO,
) (*dto.RenderResponseDTO, error) {
	var site, page interface{}
	var err error

	// 根据域名或站点ID获取站点信息
	if request.Domain != "" {
		site, err = s.siteService.GetSiteByDomain(ctx, request.Domain)
		if err != nil {
			return nil, errors.New("站点不存在: " + err.Error())
		}
	} else if request.SiteID != "" {
		// 检查站点是否已发布
		if !s.siteService.IsSitePublished(ctx, request.SiteID) {
			return nil, errors.New("站点不存在或未发布")
		}

		site, err = s.siteService.GetSiteByID(ctx, request.SiteID)
		if err != nil {
			return nil, errors.New("站点不存在: " + err.Error())
		}
	} else {
		return nil, errors.New("必须提供域名或站点ID")
	}

	// 获取页面信息
	siteID := extractSiteID(site) // 这个函数需要根据实际的Site结构实现
	if request.Slug != "" {
		page, err = s.pageService.GetPageBySlug(ctx, siteID, request.Slug)
		if err != nil {
			// 如果是首页相关的slug，尝试获取首页
			if request.Slug == "" || request.Slug == "index" || request.Slug == "home" {
				page, err = s.pageService.GetHomePage(ctx, siteID)
				if err != nil {
					return nil, errors.New("页面不存在: " + err.Error())
				}
			} else {
				return nil, errors.New("页面不存在: " + err.Error())
			}
		}
	} else {
		// 默认获取首页
		page, err = s.pageService.GetHomePage(ctx, siteID)
		if err != nil {
			return nil, errors.New("首页不存在: " + err.Error())
		}
	}

	// 获取页面模板ID
	templateID := extractTemplateID(page) // 这个函数需要根据实际的Page结构实现

	// 构造渲染请求
	renderRequest := dto.RenderRequestDTO{
		TemplateID: templateID,
		SiteID:     siteID,
		Format:     getFormatOrDefault(request.Format, "html"),
		Data: map[string]interface{}{
			"site": site,
			"page": page,
		},
		CacheConfig: dto.CacheConfigDTO{
			Enabled: true,
			TTL:     3600, // 默认缓存1小时
			Level:   "memory",
		},
	}

	// 调用渲染模板方法
	return s.RenderTemplate(ctx, renderRequest)
}

// PreviewPage 预览页面
func (s *RenderApplicationService) PreviewPage(
	ctx context.Context,
	request dto.PreviewRequestDTO,
) (*dto.RenderResponseDTO, error) {
	// 获取站点信息
	site, err := s.siteService.GetSiteByID(ctx, request.SiteID)
	if err != nil {
		return nil, errors.New("站点不存在: " + err.Error())
	}

	var page interface{}
	// 获取页面信息
	if request.PageID != "" {
		page, err = s.pageService.GetPageByID(ctx, request.PageID)
		if err != nil {
			return nil, errors.New("页面不存在: " + err.Error())
		}
	} else {
		// 默认预览首页
		page, err = s.pageService.GetHomePage(ctx, request.SiteID)
		if err != nil {
			return nil, errors.New("首页不存在: " + err.Error())
		}
	}

	// 获取页面模板ID
	templateID := extractTemplateID(page) // 这个函数需要根据实际的Page结构实现

	// 合并请求数据和默认数据
	data := map[string]interface{}{
		"site":     site,
		"page":     page,
		"preview":  true,
		"deviceType": request.DeviceType,
	}

	// 添加请求中的自定义数据
	if request.Data != nil {
		for k, v := range request.Data {
			data[k] = v
		}
	}

	// 构造渲染请求，预览不启用缓存
	renderRequest := dto.RenderRequestDTO{
		TemplateID: templateID,
		SiteID:     request.SiteID,
		PageID:     request.PageID,
		Format:     "html",
		Data:       data,
		CacheConfig: dto.CacheConfigDTO{
			Enabled: false,
		},
	}

	// 调用渲染模板方法
	return s.RenderTemplate(ctx, renderRequest)
}

// CreateTemplate 创建模板
func (s *RenderApplicationService) CreateTemplate(
	ctx context.Context,
	request dto.CreateTemplateDTO,
) (*dto.TemplateDTO, error) {
	// 验证站点是否存在
	_, err := s.siteService.GetSiteByID(ctx, request.SiteID)
	if err != nil {
		return nil, errors.New("站点不存在: " + err.Error())
	}

	// 生成模板ID
	templateID := uuid.New().String()

	// 调用领域服务创建模板
	template, err := s.renderDomainService.CreateTemplate(
		ctx,
		templateID,
		request.Name,
		request.Description,
		request.Content,
		request.Type,
		request.Version,
		request.SiteID,
		request.Metadata,
	)

	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.TemplateDTO{
		ID:          template.ID(),
		Name:        template.Name(),
		Description: template.Description(),
		Content:     template.Content(),
		Type:        template.Type(),
		Version:     template.Version(),
		SiteID:      template.SiteID(),
		Metadata:    template.Metadata(),
		CreatedAt:   template.CreatedAt(),
		UpdatedAt:   template.UpdatedAt(),
	}, nil
}

// UpdateTemplate 更新模板
func (s *RenderApplicationService) UpdateTemplate(
	ctx context.Context,
	templateID string,
	request dto.UpdateTemplateDTO,
) (*dto.TemplateDTO, error) {
	// 调用领域服务更新模板
	template, err := s.renderDomainService.UpdateTemplate(
		ctx,
		templateID,
		request.Name,
		request.Description,
		request.Content,
		request.Type,
		request.Version,
		request.Metadata,
	)

	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.TemplateDTO{
		ID:          template.ID(),
		Name:        template.Name(),
		Description: template.Description(),
		Content:     template.Content(),
		Type:        template.Type(),
		Version:     template.Version(),
		SiteID:      template.SiteID(),
		Metadata:    template.Metadata(),
		CreatedAt:   template.CreatedAt(),
		UpdatedAt:   template.UpdatedAt(),
	}, nil
}

// DeleteTemplate 删除模板
func (s *RenderApplicationService) DeleteTemplate(
	ctx context.Context,
	templateID string,
) error {
	return s.renderDomainService.DeleteTemplate(ctx, templateID)
}

// GetTemplate 获取模板
func (s *RenderApplicationService) GetTemplate(
	ctx context.Context,
	templateID string,
) (*dto.TemplateDTO, error) {
	// 调用领域服务获取模板
	template, err := s.renderDomainService.GetTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &dto.TemplateDTO{
		ID:          template.ID(),
		Name:        template.Name(),
		Description: template.Description(),
		Content:     template.Content(),
		Type:        template.Type(),
		Version:     template.Version(),
		SiteID:      template.SiteID(),
		Metadata:    template.Metadata(),
		CreatedAt:   template.CreatedAt(),
		UpdatedAt:   template.UpdatedAt(),
	}, nil
}

// 工具函数：从站点对象中提取站点ID
func extractSiteID(site interface{}) string {
	// 这个函数需要根据实际的Site结构实现
	// 简化实现，实际应用中需要根据站点对象结构获取ID
	if site, ok := site.(map[string]interface{}); ok {
		if id, ok := site["id"].(string); ok {
			return id
		}
	}
	return ""
}

// 工具函数：从页面对象中提取模板ID
func extractTemplateID(page interface{}) string {
	// 这个函数需要根据实际的Page结构实现
	// 简化实现，实际应用中需要根据页面对象结构获取模板ID
	if page, ok := page.(map[string]interface{}); ok {
		if templateID, ok := page["templateId"].(string); ok {
			return templateID
		}
	}
	return ""
}

// 工具函数：获取格式或默认值
func getFormatOrDefault(format string, defaultFormat string) string {
	if format == "" {
		return defaultFormat
	}
	return format
}
