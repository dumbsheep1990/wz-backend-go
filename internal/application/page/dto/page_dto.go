package dto

import (
	"time"
)

// Commands (命令) - 写操作

// CreatePageCommand 创建页面命令
type CreatePageCommand struct {
	SiteID      string            `json:"site_id" validate:"required"`
	Name        string            `json:"name" validate:"required,min=1,max=100"`
	Slug        string            `json:"slug" validate:"required,min=1,max=100"`
	Title       string            `json:"title" validate:"required,min=1,max=200"`
	Description string            `json:"description" validate:"max=300"`
	Keywords    []string          `json:"keywords" validate:"max=10,dive,min=1,max=30"`
	Layout      string            `json:"layout" validate:"required"`
	Content     string            `json:"content"`
	Meta        map[string]string `json:"meta"`
}

// UpdatePageCommand 更新页面命令
type UpdatePageCommand struct {
	PageID      string            `json:"page_id" validate:"required"`
	Name        *string           `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Slug        *string           `json:"slug,omitempty" validate:"omitempty,min=1,max=100"`
	Title       *string           `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=300"`
	Keywords    []string          `json:"keywords,omitempty" validate:"omitempty,max=10,dive,min=1,max=30"`
	Layout      *string           `json:"layout,omitempty" validate:"omitempty,min=1"`
	Content     *string           `json:"content,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
}

// UpdateSEOCommand 更新SEO信息命令
type UpdateSEOCommand struct {
	PageID      string   `json:"page_id" validate:"required"`
	Description string   `json:"description" validate:"max=300"`
	Keywords    []string `json:"keywords" validate:"max=10,dive,min=1,max=30"`
}

// PublishPageCommand 发布页面命令
type PublishPageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// UnpublishPageCommand 取消发布页面命令
type UnpublishPageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// SetPrivatePageCommand 设置页面为私有命令
type SetPrivatePageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// ArchivePageCommand 归档页面命令
type ArchivePageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// DeletePageCommand 删除页面命令
type DeletePageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// RestorePageCommand 恢复页面命令
type RestorePageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// SetHomepageCommand 设置首页命令
type SetHomepageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// UnsetHomepageCommand 取消首页设置命令
type UnsetHomepageCommand struct {
	PageID string `json:"page_id" validate:"required"`
}

// UpdateSortOrderCommand 更新排序命令
type UpdateSortOrderCommand struct {
	PageID    string `json:"page_id" validate:"required"`
	SortOrder int32  `json:"sort_order" validate:"min=0"`
}

// Queries (查询) - 读操作

// GetPageQuery 获取单个页面查询
type GetPageQuery struct {
	PageID string `json:"page_id" validate:"required"`
}

// GetPageBySlugQuery 根据URL段获取页面查询
type GetPageBySlugQuery struct {
	SiteID string `json:"site_id" validate:"required"`
	Slug   string `json:"slug" validate:"required"`
}

// GetHomepageQuery 获取首页查询
type GetHomepageQuery struct {
	SiteID string `json:"site_id" validate:"required"`
}

// ListPagesQuery 页面列表查询
type ListPagesQuery struct {
	SiteID   string            `json:"site_id" validate:"required"`
	Status   *int32            `json:"status,omitempty"`
	Keyword  string            `json:"keyword,omitempty"`
	Layout   string            `json:"layout,omitempty"`
	Page     int               `json:"page" validate:"min=1"`
	PageSize int               `json:"page_size" validate:"min=1,max=100"`
	OrderBy  string            `json:"order_by,omitempty"` // created_at, updated_at, sort_order, name
	OrderDir string            `json:"order_dir,omitempty"` // asc, desc
	Filters  map[string]string `json:"filters,omitempty"`
}

// SearchPagesQuery 搜索页面查询
type SearchPagesQuery struct {
	SiteID   string `json:"site_id" validate:"required"`
	Keyword  string `json:"keyword" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// GetPageSEOAnalysisQuery 获取页面SEO分析查询
type GetPageSEOAnalysisQuery struct {
	PageID string `json:"page_id" validate:"required"`
}

// GetSitePageStatsQuery 获取站点页面统计查询
type GetSitePageStatsQuery struct {
	SiteID string `json:"site_id" validate:"required"`
}

// Responses (响应)

// PageResponse 页面响应
type PageResponse struct {
	ID           string                 `json:"id"`
	SiteID       string                 `json:"site_id"`
	Name         string                 `json:"name"`
	Slug         string                 `json:"slug"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Keywords     []string               `json:"keywords"`
	Layout       string                 `json:"layout"`
	Status       int32                  `json:"status"`
	StatusText   string                 `json:"status_text"`
	IsHomepage   bool                   `json:"is_homepage"`
	IsVisible    bool                   `json:"is_visible"`
	IsEditable   bool                   `json:"is_editable"`
	SortOrder    int32                  `json:"sort_order"`
	Content      string                 `json:"content"`
	URL          string                 `json:"url"`
	SEOScore     int                    `json:"seo_score"`
	PublishedAt  *time.Time             `json:"published_at"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Meta         map[string]interface{} `json:"meta"`
}

// PageListResponse 页面列表响应
type PageListResponse struct {
	Pages      []PageResponse `json:"pages"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// PageSearchResponse 页面搜索响应
type PageSearchResponse struct {
	Pages    []PageResponse `json:"pages"`
	Total    int64          `json:"total"`
	Keyword  string         `json:"keyword"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// PageSEOAnalysisResponse 页面SEO分析响应
type PageSEOAnalysisResponse struct {
	PageID      string                 `json:"page_id"`
	SEOScore    int                    `json:"seo_score"`
	Title       SEOItemAnalysis        `json:"title"`
	Description SEOItemAnalysis        `json:"description"`
	Keywords    SEOItemAnalysis        `json:"keywords"`
	URL         SEOItemAnalysis        `json:"url"`
	Content     SEOItemAnalysis        `json:"content"`
	Suggestions []string               `json:"suggestions"`
	MetaTags    map[string]string      `json:"meta_tags"`
}

// SEOItemAnalysis SEO项目分析
type SEOItemAnalysis struct {
	Score       int                    `json:"score"`
	MaxScore    int                    `json:"max_score"`
	Status      string                 `json:"status"` // optimal, good, poor, missing
	Suggestion  string                 `json:"suggestion"`
	Details     map[string]interface{} `json:"details"`
}

// SitePageStatsResponse 站点页面统计响应
type SitePageStatsResponse struct {
	SiteID         string               `json:"site_id"`
	TotalPages     int64                `json:"total_pages"`
	PublishedPages int64                `json:"published_pages"`
	DraftPages     int64                `json:"draft_pages"`
	PrivatePages   int64                `json:"private_pages"`
	ArchivedPages  int64                `json:"archived_pages"`
	DeletedPages   int64                `json:"deleted_pages"`
	StatusStats    map[string]int64     `json:"status_stats"`
	LayoutStats    map[string]int64     `json:"layout_stats"`
	SEOStats       SEOStatsResponse     `json:"seo_stats"`
	RecentPages    []PageResponse       `json:"recent_pages"`
}

// SEOStatsResponse SEO统计响应
type SEOStatsResponse struct {
	AverageSEOScore    float64 `json:"average_seo_score"`
	PagesWithGoodSEO   int64   `json:"pages_with_good_seo"`   // score >= 80
	PagesWithPoorSEO   int64   `json:"pages_with_poor_seo"`   // score < 50
	PagesNeedingUpdate int64   `json:"pages_needing_update"`
}

// PageCreatedResponse 页面创建响应
type PageCreatedResponse struct {
	PageID    string `json:"page_id"`
	Message   string `json:"message"`
	URL       string `json:"url"`
	SEOScore  int    `json:"seo_score"`
}

// PageUpdatedResponse 页面更新响应
type PageUpdatedResponse struct {
	PageID   string `json:"page_id"`
	Message  string `json:"message"`
	SEOScore int    `json:"seo_score"`
}

// PageDeletedResponse 页面删除响应
type PageDeletedResponse struct {
	PageID  string `json:"page_id"`
	Message string `json:"message"`
}

// PageStatusChangedResponse 页面状态变更响应
type PageStatusChangedResponse struct {
	PageID    string `json:"page_id"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	Message   string `json:"message"`
	URL       string `json:"url,omitempty"`
}

// CommandResult 命令结果
type CommandResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

// BatchCommandResult 批量命令结果
type BatchCommandResult struct {
	Success      bool     `json:"success"`
	ProcessedIDs []string `json:"processed_ids"`
	FailedIDs    []string `json:"failed_ids"`
	Message      string   `json:"message"`
} 