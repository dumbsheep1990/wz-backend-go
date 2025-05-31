package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"wz-backend-go/internal/application/page/dto"
	"wz-backend-go/internal/domain/page/entity"
	"wz-backend-go/internal/domain/page/valueobject"
	"wz-backend-go/internal/domain/shared/event"
)

// PageRepository 页面仓储接口
type PageRepository interface {
	// 写操作
	Save(ctx context.Context, page *entity.Page) error
	Update(ctx context.Context, page *entity.Page) error
	Delete(ctx context.Context, pageID valueobject.PageID) error

	// 读操作
	GetByID(ctx context.Context, pageID valueobject.PageID) (*entity.Page, error)
	GetBySlug(ctx context.Context, siteID string, slug valueobject.PageSlug) (*entity.Page, error)
	GetHomepage(ctx context.Context, siteID string) (*entity.Page, error)
	List(ctx context.Context, query dto.ListPagesQuery) ([]entity.Page, int64, error)
	Search(ctx context.Context, query dto.SearchPagesQuery) ([]entity.Page, int64, error)
	GetStats(ctx context.Context, siteID string) (*dto.SitePageStatsResponse, error)

	// 业务查询
	ExistsBySlug(ctx context.Context, siteID string, slug valueobject.PageSlug, excludePageID *valueobject.PageID) (bool, error)
	GetNextSortOrder(ctx context.Context, siteID string) (int32, error)
}

// EventPublisher 事件发布器接口
type EventPublisher interface {
	Publish(ctx context.Context, events []event.DomainEvent) error
}

// SlugGenerator URL段生成器接口
type SlugGenerator interface {
	GenerateFromTitle(title string) (valueobject.PageSlug, error)
	EnsureUnique(ctx context.Context, siteID string, baseSlug valueobject.PageSlug, excludePageID *valueobject.PageID) (valueobject.PageSlug, error)
}

// PageApplicationService 页面应用服务
type PageApplicationService struct {
	pageRepo      PageRepository
	eventPublisher EventPublisher
	slugGenerator SlugGenerator
}

// NewPageApplicationService 创建页面应用服务实例
func NewPageApplicationService(
	pageRepo PageRepository,
	eventPublisher EventPublisher,
	slugGenerator SlugGenerator,
) *PageApplicationService {
	return &PageApplicationService{
		pageRepo:      pageRepo,
		eventPublisher: eventPublisher,
		slugGenerator: slugGenerator,
	}
}

// Command Handlers (命令处理器)

// CreatePage 创建页面
func (s *PageApplicationService) CreatePage(ctx context.Context, cmd dto.CreatePageCommand) (*dto.PageCreatedResponse, error) {
	// 生成页面ID
	pageID := valueobject.GeneratePageID()

	// 创建页面标题
	title, err := valueobject.NewPageTitle(cmd.Title)
	if err != nil {
		return nil, fmt.Errorf("创建页面标题失败: %w", err)
	}

	// 处理URL段
	var slug valueobject.PageSlug
	if cmd.Slug != "" {
		slug, err = valueobject.NewPageSlug(cmd.Slug)
		if err != nil {
			return nil, fmt.Errorf("创建页面URL段失败: %w", err)
		}
	} else {
		// 从标题生成URL段
		slug, err = s.slugGenerator.GenerateFromTitle(cmd.Title)
		if err != nil {
			return nil, fmt.Errorf("生成页面URL段失败: %w", err)
		}
	}

	// 确保URL段唯一性
	slug, err = s.slugGenerator.EnsureUnique(ctx, cmd.SiteID, slug, nil)
	if err != nil {
		return nil, fmt.Errorf("确保URL段唯一性失败: %w", err)
	}

	// 创建SEO元数据
	seoMeta, err := valueobject.NewSEOMeta(cmd.Description, cmd.Keywords)
	if err != nil {
		return nil, fmt.Errorf("创建SEO元数据失败: %w", err)
	}

	// 创建页面实体
	page, err := entity.NewPage(
		pageID,
		cmd.SiteID,
		cmd.Name,
		slug,
		title,
		seoMeta,
		cmd.Layout,
	)
	if err != nil {
		return nil, fmt.Errorf("创建页面实体失败: %w", err)
	}

	// 设置内容
	if cmd.Content != "" {
		page.UpdateContent(cmd.Content)
	}

	// 设置排序顺序
	sortOrder, err := s.pageRepo.GetNextSortOrder(ctx, cmd.SiteID)
	if err != nil {
		return nil, fmt.Errorf("获取排序顺序失败: %w", err)
	}
	page.UpdateSortOrder(sortOrder)

	// 保存页面
	if err := s.pageRepo.Save(ctx, page); err != nil {
		return nil, fmt.Errorf("保存页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
		// log.Error("发布领域事件失败", err)
	}
	page.ClearDomainEvents()

	return &dto.PageCreatedResponse{
		PageID:   pageID.Value(),
		Message:  "页面创建成功",
		URL:      page.GetURL(),
		SEOScore: page.GetSEOScore(),
	}, nil
}

// UpdatePage 更新页面
func (s *PageApplicationService) UpdatePage(ctx context.Context, cmd dto.UpdatePageCommand) (*dto.PageUpdatedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	// 更新页面名称
	if cmd.Name != nil && *cmd.Name != page.Name() {
		if err := page.UpdateName(*cmd.Name); err != nil {
			return nil, fmt.Errorf("更新页面名称失败: %w", err)
		}
	}

	// 更新页面标题
	if cmd.Title != nil {
		title, err := valueobject.NewPageTitle(*cmd.Title)
		if err != nil {
			return nil, fmt.Errorf("创建页面标题失败: %w", err)
		}
		if err := page.UpdateTitle(title); err != nil {
			return nil, fmt.Errorf("更新页面标题失败: %w", err)
		}
	}

	// 更新URL段
	if cmd.Slug != nil {
		slug, err := valueobject.NewPageSlug(*cmd.Slug)
		if err != nil {
			return nil, fmt.Errorf("创建页面URL段失败: %w", err)
		}

		// 检查URL段唯一性
		slug, err = s.slugGenerator.EnsureUnique(ctx, page.SiteID(), slug, &pageID)
		if err != nil {
			return nil, fmt.Errorf("确保URL段唯一性失败: %w", err)
		}

		if err := page.UpdateSlug(slug); err != nil {
			return nil, fmt.Errorf("更新页面URL段失败: %w", err)
		}
	}

	// 更新SEO元数据
	if cmd.Description != nil || len(cmd.Keywords) > 0 {
		description := ""
		keywords := page.SEOMeta().Keywords()

		if cmd.Description != nil {
			description = *cmd.Description
		} else {
			description = page.SEOMeta().Description()
		}

		if len(cmd.Keywords) > 0 {
			keywords = cmd.Keywords
		}

		seoMeta, err := valueobject.NewSEOMeta(description, keywords)
		if err != nil {
			return nil, fmt.Errorf("创建SEO元数据失败: %w", err)
		}
		page.UpdateSEOMeta(seoMeta)
	}

	// 更新布局
	if cmd.Layout != nil {
		if err := page.UpdateLayout(*cmd.Layout); err != nil {
			return nil, fmt.Errorf("更新布局失败: %w", err)
		}
	}

	// 更新内容
	if cmd.Content != nil {
		page.UpdateContent(*cmd.Content)
	}

	// 更新页面
	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageUpdatedResponse{
		PageID:   pageID.Value(),
		Message:  "页面更新成功",
		SEOScore: page.GetSEOScore(),
	}, nil
}

// UpdatePageSEO 更新页面SEO信息
func (s *PageApplicationService) UpdatePageSEO(ctx context.Context, cmd dto.UpdateSEOCommand) (*dto.PageUpdatedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	seoMeta, err := valueobject.NewSEOMeta(cmd.Description, cmd.Keywords)
	if err != nil {
		return nil, fmt.Errorf("创建SEO元数据失败: %w", err)
	}

	page.UpdateSEOMeta(seoMeta)

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageUpdatedResponse{
		PageID:   pageID.Value(),
		Message:  "SEO信息更新成功",
		SEOScore: page.GetSEOScore(),
	}, nil
}

// PublishPage 发布页面
func (s *PageApplicationService) PublishPage(ctx context.Context, cmd dto.PublishPageCommand) (*dto.PageStatusChangedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	oldStatus := page.Status()
	if err := page.Publish(); err != nil {
		return nil, fmt.Errorf("发布页面失败: %w", err)
	}

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageStatusChangedResponse{
		PageID:    pageID.Value(),
		OldStatus: oldStatus.String(),
		NewStatus: page.Status().String(),
		Message:   "页面发布成功",
		URL:       page.GetURL(),
	}, nil
}

// UnpublishPage 取消发布页面
func (s *PageApplicationService) UnpublishPage(ctx context.Context, cmd dto.UnpublishPageCommand) (*dto.PageStatusChangedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	oldStatus := page.Status()
	if err := page.Unpublish(); err != nil {
		return nil, fmt.Errorf("取消发布失败: %w", err)
	}

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageStatusChangedResponse{
		PageID:    pageID.Value(),
		OldStatus: oldStatus.String(),
		NewStatus: page.Status().String(),
		Message:   "页面已取消发布",
	}, nil
}

// SetPagePrivate 设置页面为私有
func (s *PageApplicationService) SetPagePrivate(ctx context.Context, cmd dto.SetPrivatePageCommand) (*dto.PageStatusChangedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	oldStatus := page.Status()
	if err := page.SetAsPrivate(); err != nil {
		return nil, fmt.Errorf("设置私有失败: %w", err)
	}

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageStatusChangedResponse{
		PageID:    pageID.Value(),
		OldStatus: oldStatus.String(),
		NewStatus: page.Status().String(),
		Message:   "页面已设为私有",
	}, nil
}

// ArchivePage 归档页面
func (s *PageApplicationService) ArchivePage(ctx context.Context, cmd dto.ArchivePageCommand) (*dto.PageStatusChangedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	oldStatus := page.Status()
	if err := page.Archive(); err != nil {
		return nil, fmt.Errorf("归档页面失败: %w", err)
	}

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageStatusChangedResponse{
		PageID:    pageID.Value(),
		OldStatus: oldStatus.String(),
		NewStatus: page.Status().String(),
		Message:   "页面已归档",
	}, nil
}

// DeletePage 删除页面
func (s *PageApplicationService) DeletePage(ctx context.Context, cmd dto.DeletePageCommand) (*dto.PageDeletedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	if err := page.Delete(); err != nil {
		return nil, fmt.Errorf("删除页面失败: %w", err)
	}

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageDeletedResponse{
		PageID:  pageID.Value(),
		Message: "页面已删除",
	}, nil
}

// RestorePage 恢复页面
func (s *PageApplicationService) RestorePage(ctx context.Context, cmd dto.RestorePageCommand) (*dto.PageStatusChangedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	oldStatus := page.Status()
	if err := page.Restore(); err != nil {
		return nil, fmt.Errorf("恢复页面失败: %w", err)
	}

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageStatusChangedResponse{
		PageID:    pageID.Value(),
		OldStatus: oldStatus.String(),
		NewStatus: page.Status().String(),
		Message:   "页面已恢复",
	}, nil
}

// SetHomepage 设置首页
func (s *PageApplicationService) SetHomepage(ctx context.Context, cmd dto.SetHomepageCommand) (*dto.PageStatusChangedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	// 先取消当前首页设置
	currentHomepage, err := s.pageRepo.GetHomepage(ctx, page.SiteID())
	if err == nil && currentHomepage != nil && !currentHomepage.ID().IsEquals(pageID) {
		currentHomepage.UnsetAsHomepage()
		if err := s.pageRepo.Update(ctx, currentHomepage); err != nil {
			return nil, fmt.Errorf("取消当前首页设置失败: %w", err)
		}
	}

	// 设置新首页
	if err := page.SetAsHomepage(); err != nil {
		return nil, fmt.Errorf("设置首页失败: %w", err)
	}

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageStatusChangedResponse{
		PageID:  pageID.Value(),
		Message: "首页设置成功",
		URL:     page.GetURL(),
	}, nil
}

// UnsetHomepage 取消首页设置
func (s *PageApplicationService) UnsetHomepage(ctx context.Context, cmd dto.UnsetHomepageCommand) (*dto.PageStatusChangedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	if !page.IsHomepage() {
		return nil, errors.New("该页面不是首页")
	}

	page.UnsetAsHomepage()

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageStatusChangedResponse{
		PageID:  pageID.Value(),
		Message: "首页设置已取消",
	}, nil
}

// UpdateSortOrder 更新排序顺序
func (s *PageApplicationService) UpdateSortOrder(ctx context.Context, cmd dto.UpdateSortOrderCommand) (*dto.PageUpdatedResponse, error) {
	pageID, err := valueobject.NewPageID(cmd.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	page.UpdateSortOrder(cmd.SortOrder)

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("更新页面失败: %w", err)
	}

	// 发布领域事件
	if err := s.eventPublisher.Publish(ctx, page.GetDomainEvents()); err != nil {
		// 记录日志但不影响主流程
	}
	page.ClearDomainEvents()

	return &dto.PageUpdatedResponse{
		PageID:  pageID.Value(),
		Message: "排序顺序更新成功",
	}, nil
}

// Query Handlers (查询处理器)

// GetPage 获取单个页面
func (s *PageApplicationService) GetPage(ctx context.Context, query dto.GetPageQuery) (*dto.PageResponse, error) {
	pageID, err := valueobject.NewPageID(query.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	return s.toPageResponse(page), nil
}

// GetPageBySlug 根据URL段获取页面
func (s *PageApplicationService) GetPageBySlug(ctx context.Context, query dto.GetPageBySlugQuery) (*dto.PageResponse, error) {
	slug, err := valueobject.NewPageSlug(query.Slug)
	if err != nil {
		return nil, fmt.Errorf("无效的URL段: %w", err)
	}

	page, err := s.pageRepo.GetBySlug(ctx, query.SiteID, slug)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	return s.toPageResponse(page), nil
}

// GetHomepage 获取首页
func (s *PageApplicationService) GetHomepage(ctx context.Context, query dto.GetHomepageQuery) (*dto.PageResponse, error) {
	page, err := s.pageRepo.GetHomepage(ctx, query.SiteID)
	if err != nil {
		return nil, fmt.Errorf("获取首页失败: %w", err)
	}

	return s.toPageResponse(page), nil
}

// ListPages 获取页面列表
func (s *PageApplicationService) ListPages(ctx context.Context, query dto.ListPagesQuery) (*dto.PageListResponse, error) {
	pages, total, err := s.pageRepo.List(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("获取页面列表失败: %w", err)
	}

	pageResponses := make([]dto.PageResponse, len(pages))
	for i, page := range pages {
		pageResponses[i] = *s.toPageResponse(&page)
	}

	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

	return &dto.PageListResponse{
		Pages:      pageResponses,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchPages 搜索页面
func (s *PageApplicationService) SearchPages(ctx context.Context, query dto.SearchPagesQuery) (*dto.PageSearchResponse, error) {
	pages, total, err := s.pageRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("搜索页面失败: %w", err)
	}

	pageResponses := make([]dto.PageResponse, len(pages))
	for i, page := range pages {
		pageResponses[i] = *s.toPageResponse(&page)
	}

	return &dto.PageSearchResponse{
		Pages:    pageResponses,
		Total:    total,
		Keyword:  query.Keyword,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// GetPageSEOAnalysis 获取页面SEO分析
func (s *PageApplicationService) GetPageSEOAnalysis(ctx context.Context, query dto.GetPageSEOAnalysisQuery) (*dto.PageSEOAnalysisResponse, error) {
	pageID, err := valueobject.NewPageID(query.PageID)
	if err != nil {
		return nil, fmt.Errorf("无效的页面ID: %w", err)
	}

	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("获取页面失败: %w", err)
	}

	return s.buildSEOAnalysis(page), nil
}

// GetSitePageStats 获取站点页面统计
func (s *PageApplicationService) GetSitePageStats(ctx context.Context, query dto.GetSitePageStatsQuery) (*dto.SitePageStatsResponse, error) {
	stats, err := s.pageRepo.GetStats(ctx, query.SiteID)
	if err != nil {
		return nil, fmt.Errorf("获取页面统计失败: %w", err)
	}

	return stats, nil
}

// Helper Methods (辅助方法)

// toPageResponse 将页面实体转换为响应DTO
func (s *PageApplicationService) toPageResponse(page *entity.Page) *dto.PageResponse {
	return &dto.PageResponse{
		ID:          page.ID().Value(),
		SiteID:      page.SiteID(),
		Name:        page.Name(),
		Slug:        page.Slug().Value(),
		Title:       page.Title().GetUnescaped(),
		Description: page.SEOMeta().GetUnescapedDescription(),
		Keywords:    page.SEOMeta().Keywords(),
		Layout:      page.Layout(),
		Status:      page.Status().Value(),
		StatusText:  page.Status().String(),
		IsHomepage:  page.IsHomepage(),
		IsVisible:   page.IsVisible(),
		IsEditable:  page.IsEditable(),
		SortOrder:   page.SortOrder(),
		Content:     page.Content(),
		URL:         page.GetURL(),
		SEOScore:    page.GetSEOScore(),
		PublishedAt: page.PublishedAt(),
		CreatedAt:   page.CreatedAt(),
		UpdatedAt:   page.UpdatedAt(),
		Meta:        page.GetDisplayInfo(),
	}
}

// buildSEOAnalysis 构建SEO分析
func (s *PageApplicationService) buildSEOAnalysis(page *entity.Page) *dto.PageSEOAnalysisResponse {
	// 标题分析
	titleLength := page.Title().GetSuggestedLength()
	titleAnalysis := dto.SEOItemAnalysis{
		Score:      0,
		MaxScore:   30,
		Status:     "poor",
		Suggestion: "",
		Details:    titleLength,
	}

	if page.Title().IsSEOOptimized() {
		titleAnalysis.Score = 30
		titleAnalysis.Status = "optimal"
		titleAnalysis.Suggestion = "标题SEO优化良好"
	} else if !page.Title().IsEmpty() {
		titleAnalysis.Score = 15
		titleAnalysis.Status = "good"
		if suggestion, ok := titleLength["suggestion"].(string); ok {
			titleAnalysis.Suggestion = suggestion
		}
	} else {
		titleAnalysis.Suggestion = "请添加页面标题"
	}

	// 描述分析
	descAnalysis := page.SEOMeta().GetDescriptionAnalysis()
	descriptionAnalysis := dto.SEOItemAnalysis{
		Score:      0,
		MaxScore:   25,
		Status:     "missing",
		Suggestion: "请添加页面描述",
		Details:    descAnalysis,
	}

	if page.SEOMeta().IsDescriptionOptimal() {
		descriptionAnalysis.Score = 25
		descriptionAnalysis.Status = "optimal"
		descriptionAnalysis.Suggestion = "描述SEO优化良好"
	} else if page.SEOMeta().Description() != "" {
		descriptionAnalysis.Score = 15
		descriptionAnalysis.Status = "good"
		if suggestion, ok := descAnalysis["suggestion"].(string); ok {
			descriptionAnalysis.Suggestion = suggestion
		}
	}

	// 关键词分析
	keywordAnalysis := page.SEOMeta().GetKeywordsAnalysis()
	keywordsAnalysis := dto.SEOItemAnalysis{
		Score:      0,
		MaxScore:   15,
		Status:     "missing",
		Suggestion: "请添加3-8个相关关键词",
		Details:    keywordAnalysis,
	}

	keywordCount := len(page.SEOMeta().Keywords())
	if keywordCount >= 3 && keywordCount <= 8 {
		keywordsAnalysis.Score = 15
		keywordsAnalysis.Status = "optimal"
		keywordsAnalysis.Suggestion = "关键词数量最佳"
	} else if keywordCount > 0 {
		keywordsAnalysis.Score = 8
		keywordsAnalysis.Status = "good"
		if suggestion, ok := keywordAnalysis["suggestion"].(string); ok {
			keywordsAnalysis.Suggestion = suggestion
		}
	}

	// URL分析
	urlAnalysis := dto.SEOItemAnalysis{
		Score:      0,
		MaxScore:   20,
		Status:     "poor",
		Suggestion: "请优化URL段",
		Details:    map[string]interface{}{},
	}

	if page.Slug().IsSEOFriendly() {
		urlAnalysis.Score = 20
		urlAnalysis.Status = "optimal"
		urlAnalysis.Suggestion = "URL SEO优化良好"
	} else if !page.Slug().IsEmpty() {
		urlAnalysis.Score = 10
		urlAnalysis.Status = "good"
		urlAnalysis.Suggestion = "建议使用更符合SEO规范的URL格式"
	}

	// 内容分析
	contentAnalysis := dto.SEOItemAnalysis{
		Score:      0,
		MaxScore:   10,
		Status:     "missing",
		Suggestion: "请添加页面内容",
		Details: map[string]interface{}{
			"length": len(page.Content()),
		},
	}

	contentLength := len(page.Content())
	if contentLength > 100 {
		contentAnalysis.Score = 10
		contentAnalysis.Status = "optimal"
		contentAnalysis.Suggestion = "内容丰富度良好"
	} else if contentLength > 0 {
		contentAnalysis.Score = 5
		contentAnalysis.Status = "good"
		contentAnalysis.Suggestion = "建议增加更多内容"
	}

	return &dto.PageSEOAnalysisResponse{
		PageID:      page.ID().Value(),
		SEOScore:    page.GetSEOScore(),
		Title:       titleAnalysis,
		Description: descriptionAnalysis,
		Keywords:    keywordsAnalysis,
		URL:         urlAnalysis,
		Content:     contentAnalysis,
		Suggestions: page.GetSEOSuggestions(),
		MetaTags:    page.SEOMeta().GetHTMLMetaTags(),
	}
} 