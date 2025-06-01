package service

import (
	"context"
	"errors"
	"time"
	"wz-backend-go/internal/domain/render/entity"
	"wz-backend-go/internal/domain/render/repository"
	"wz-backend-go/internal/domain/render/valueobject"
)

// RenderDomainService 渲染领域服务
type RenderDomainService struct {
	renderRepo    repository.RenderResultRepository
	templateRepo  repository.TemplateRepository
	eventPublisher EventPublisher
}

// EventPublisher 事件发布接口
type EventPublisher interface {
	Publish(event entity.RenderEvent)
}

// NewRenderDomainService 创建一个新的渲染领域服务
func NewRenderDomainService(
	renderRepo repository.RenderResultRepository,
	templateRepo repository.TemplateRepository,
	eventPublisher EventPublisher,
) *RenderDomainService {
	return &RenderDomainService{
		renderRepo:    renderRepo,
		templateRepo:  templateRepo,
		eventPublisher: eventPublisher,
	}
}

// RenderTemplate 渲染模板
func (s *RenderDomainService) RenderTemplate(
	ctx context.Context,
	templateID string,
	context valueobject.TemplateContext,
	format valueobject.RenderFormat,
	cacheOptions *CacheOptions,
) (*entity.RenderResult, error) {
	// 查找模板
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, errors.New("模板不存在: " + err.Error())
	}

	// 检查缓存
	if cacheOptions != nil && cacheOptions.Enabled {
		cacheKey := s.buildCacheKey(templateID, format.Format(), context)
		cachedResult, err := s.renderRepo.FindByCacheKey(ctx, cacheKey)
		
		if err == nil && !cachedResult.IsExpired() {
			// 发布缓存命中事件
			s.eventPublisher.Publish(entity.NewCacheHitEvent(
				templateID,
				cacheKey,
				valueobject.CacheLevelMemory,
				0, // 响应时间，实际实现需要计算
			))
			
			return cachedResult, nil
		}
	}

	// 开始渲染计时
	startTime := time.Now()

	// 渲染模板
	result, err := template.Render(context, format)
	if err != nil {
		// 发布渲染失败事件
		s.eventPublisher.Publish(entity.NewRenderFailedEvent(
			templateID,
			err.Error(),
			templateID,
			template.Type(),
			time.Since(startTime).Milliseconds(),
		))
		
		return nil, errors.New("渲染失败: " + err.Error())
	}

	// 配置缓存策略
	if cacheOptions != nil {
		cacheKey := s.buildCacheKey(templateID, format.Format(), context)
		cacheGroups := []string{template.SiteID(), template.Type()}
		
		cacheStrategy, _ := valueobject.NewCacheStrategy(
			cacheOptions.Enabled,
			cacheOptions.TTL,
			cacheOptions.Level,
			cacheKey,
			cacheGroups,
		)
		
		result.WithCacheStrategy(cacheStrategy)
	}

	// 保存渲染结果到仓储
	if result.CacheStrategy().Enabled() {
		if err := s.renderRepo.Save(ctx, result); err != nil {
			// 记录错误但不中断渲染流程
			// 实际应用中可能需要日志记录
		}
	}

	// 发布渲染完成事件
	s.eventPublisher.Publish(entity.NewRenderCompletedEvent(
		templateID,
		result.ID(),
		format,
		result.CacheStrategy().Enabled(),
		time.Since(startTime).Milliseconds(),
		len(result.Content()),
	))

	return result, nil
}

// CacheOptions 缓存选项
type CacheOptions struct {
	Enabled bool
	TTL     time.Duration
	Level   valueobject.CacheLevel
}

// 构建缓存键
func (s *RenderDomainService) buildCacheKey(templateID string, format string, context valueobject.TemplateContext) string {
	// 在实际实现中，这里应该基于模板ID、格式和关键上下文数据生成唯一的缓存键
	// 简化实现，实际应用中可能需要哈希计算
	return templateID + ":" + format
}

// GetTemplate 获取模板
func (s *RenderDomainService) GetTemplate(ctx context.Context, templateID string) (*entity.Template, error) {
	return s.templateRepo.FindByID(ctx, templateID)
}

// CreateTemplate 创建模板
func (s *RenderDomainService) CreateTemplate(
	ctx context.Context,
	id string,
	name string,
	description string,
	content string,
	type_ string,
	version string,
	siteID string,
	metadata map[string]string,
) (*entity.Template, error) {
	// 检查名称是否已存在
	existingTemplate, err := s.templateRepo.FindByName(ctx, name, siteID)
	if err == nil && existingTemplate != nil {
		return nil, errors.New("同名模板已存在")
	}

	// 创建新模板
	template, err := entity.NewTemplate(
		id,
		name,
		description,
		content,
		type_,
		version,
		siteID,
		metadata,
	)
	
	if err != nil {
		return nil, err
	}

	// 保存到仓储
	if err := s.templateRepo.Save(ctx, template); err != nil {
		return nil, err
	}

	// 发布模板创建事件
	s.eventPublisher.Publish(entity.NewTemplateUpdatedEvent(
		template.ID(),
		template.Name(),
		template.Type(),
		template.SiteID(),
	))

	return template, nil
}

// UpdateTemplate 更新模板
func (s *RenderDomainService) UpdateTemplate(
	ctx context.Context,
	templateID string,
	name string,
	description string,
	content string,
	type_ string,
	version string,
	metadata map[string]string,
) (*entity.Template, error) {
	// 查找模板
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, errors.New("模板不存在: " + err.Error())
	}

	// 如果名称变更，检查是否与其他模板冲突
	if name != template.Name() {
		existingTemplate, err := s.templateRepo.FindByName(ctx, name, template.SiteID())
		if err == nil && existingTemplate != nil && existingTemplate.ID() != templateID {
			return nil, errors.New("同名模板已存在")
		}

		if err := template.SetName(name); err != nil {
			return nil, err
		}
	}

	// 更新其他字段
	template.SetDescription(description)
	if err := template.SetContent(content); err != nil {
		return nil, err
	}
	
	if err := template.SetType(type_); err != nil {
		return nil, err
	}
	
	template.SetVersion(version)

	// 更新元数据
	for key, value := range metadata {
		template.SetMetadata(key, value)
	}

	// 保存到仓储
	if err := s.templateRepo.Save(ctx, template); err != nil {
		return nil, err
	}

	// 发布模板更新事件
	s.eventPublisher.Publish(entity.NewTemplateUpdatedEvent(
		template.ID(),
		template.Name(),
		template.Type(),
		template.SiteID(),
	))

	// 删除相关的缓存
	s.renderRepo.DeleteByGroup(ctx, templateID)

	return template, nil
}

// DeleteTemplate 删除模板
func (s *RenderDomainService) DeleteTemplate(ctx context.Context, templateID string) error {
	// 查找模板
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return errors.New("模板不存在: " + err.Error())
	}

	// 删除模板
	if err := s.templateRepo.Delete(ctx, templateID); err != nil {
		return err
	}

	// 删除相关的缓存
	s.renderRepo.DeleteByGroup(ctx, templateID)

	// 发布模板删除事件
	s.eventPublisher.Publish(&entity.BaseRenderEvent{
		EventType:  entity.EventTypeTemplateDeleted,
		OccurredAt: time.Now(),
		EntityID:   templateID,
	})

	return nil
}

// CleanExpiredCache 清理过期缓存
func (s *RenderDomainService) CleanExpiredCache(ctx context.Context) error {
	return s.renderRepo.DeleteExpired(ctx)
}

// InvalidateCacheByGroup 按组失效缓存
func (s *RenderDomainService) InvalidateCacheByGroup(ctx context.Context, group string) error {
	return s.renderRepo.DeleteByGroup(ctx, group)
}

// InvalidateCacheBySite 按站点失效缓存
func (s *RenderDomainService) InvalidateCacheBySite(ctx context.Context, siteID string) error {
	return s.renderRepo.DeleteBySiteID(ctx, siteID)
}
