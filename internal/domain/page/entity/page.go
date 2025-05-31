package entity

import (
	"errors"
	"strings"
	"time"

	"wz-backend-go/internal/domain/page/valueobject"
	"wz-backend-go/internal/domain/shared/event"
)

// Page 页面聚合根实体
type Page struct {
	id          valueobject.PageID
	siteID      string // 所属站点ID，应该也是值对象，但为简化暂用string
	name        string // 页面名称
	slug        valueobject.PageSlug
	title       valueobject.PageTitle
	seoMeta     valueobject.SEOMeta
	layout      string // 布局模板
	status      valueobject.PageStatus
	isHomepage  bool
	sortOrder   int32
	content     string // 页面内容（JSON格式的组件配置）
	publishedAt *time.Time
	createdAt   time.Time
	updatedAt   time.Time
	
	domainEvents []event.DomainEvent
}

// NewPage 创建新页面实体
func NewPage(
	id valueobject.PageID,
	siteID string,
	name string,
	slug valueobject.PageSlug,
	title valueobject.PageTitle,
	seoMeta valueobject.SEOMeta,
	layout string,
) (*Page, error) {
	if id.IsEmpty() {
		return nil, errors.New("页面ID不能为空")
	}
	if siteID == "" {
		return nil, errors.New("站点ID不能为空")
	}
	if name == "" {
		return nil, errors.New("页面名称不能为空")
	}
	if slug.IsEmpty() {
		return nil, errors.New("页面URL段不能为空")
	}
	if title.IsEmpty() {
		return nil, errors.New("页面标题不能为空")
	}
	if layout == "" {
		layout = "default" // 默认布局
	}

	now := time.Now()
	page := &Page{
		id:           id,
		siteID:       siteID,
		name:         name,
		slug:         slug,
		title:        title,
		seoMeta:      seoMeta,
		layout:       layout,
		status:       valueobject.NewDraftPageStatus(),
		isHomepage:   false,
		sortOrder:    0,
		content:      "",
		publishedAt:  nil,
		createdAt:    now,
		updatedAt:    now,
		domainEvents: make([]event.DomainEvent, 0),
	}

	page.addDomainEvent(NewPageCreatedEvent(page))
	return page, nil
}

// ReconstructPage 从存储中重建页面实体
func ReconstructPage(
	id valueobject.PageID,
	siteID string,
	name string,
	slug valueobject.PageSlug,
	title valueobject.PageTitle,
	seoMeta valueobject.SEOMeta,
	layout string,
	status valueobject.PageStatus,
	isHomepage bool,
	sortOrder int32,
	content string,
	publishedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Page {
	return &Page{
		id:           id,
		siteID:       siteID,
		name:         name,
		slug:         slug,
		title:        title,
		seoMeta:      seoMeta,
		layout:       layout,
		status:       status,
		isHomepage:   isHomepage,
		sortOrder:    sortOrder,
		content:      content,
		publishedAt:  publishedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		domainEvents: make([]event.DomainEvent, 0),
	}
}

// ID 返回页面ID
func (p *Page) ID() valueobject.PageID {
	return p.id
}

// SiteID 返回站点ID
func (p *Page) SiteID() string {
	return p.siteID
}

// Name 返回页面名称
func (p *Page) Name() string {
	return p.name
}

// Slug 返回页面URL段
func (p *Page) Slug() valueobject.PageSlug {
	return p.slug
}

// Title 返回页面标题
func (p *Page) Title() valueobject.PageTitle {
	return p.title
}

// SEOMeta 返回SEO元数据
func (p *Page) SEOMeta() valueobject.SEOMeta {
	return p.seoMeta
}

// Layout 返回布局模板
func (p *Page) Layout() string {
	return p.layout
}

// Status 返回页面状态
func (p *Page) Status() valueobject.PageStatus {
	return p.status
}

// IsHomepage 返回是否为首页
func (p *Page) IsHomepage() bool {
	return p.isHomepage
}

// SortOrder 返回排序顺序
func (p *Page) SortOrder() int32 {
	return p.sortOrder
}

// Content 返回页面内容
func (p *Page) Content() string {
	return p.content
}

// PublishedAt 返回发布时间
func (p *Page) PublishedAt() *time.Time {
	return p.publishedAt
}

// CreatedAt 返回创建时间
func (p *Page) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt 返回最后更新时间
func (p *Page) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdateName 更新页面名称
func (p *Page) UpdateName(name string) error {
	if name == "" {
		return errors.New("页面名称不能为空")
	}
	
	if p.name == name {
		return nil // 名称未变化
	}
	
	oldName := p.name
	p.name = name
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "name", oldName, name))
	
	return nil
}

// UpdateSlug 更新页面URL段
func (p *Page) UpdateSlug(slug valueobject.PageSlug) error {
	if slug.IsEmpty() {
		return errors.New("页面URL段不能为空")
	}
	
	if p.slug.IsEquals(slug) {
		return nil // URL段未变化
	}
	
	oldSlug := p.slug
	p.slug = slug
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "slug", oldSlug.Value(), slug.Value()))
	
	return nil
}

// UpdateTitle 更新页面标题
func (p *Page) UpdateTitle(title valueobject.PageTitle) error {
	if title.IsEmpty() {
		return errors.New("页面标题不能为空")
	}
	
	if p.title.IsEquals(title) {
		return nil // 标题未变化
	}
	
	oldTitle := p.title
	p.title = title
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "title", oldTitle.Value(), title.Value()))
	
	return nil
}

// UpdateSEOMeta 更新SEO元数据
func (p *Page) UpdateSEOMeta(seoMeta valueobject.SEOMeta) {
	p.seoMeta = seoMeta
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "seo_meta", "", ""))
}

// UpdateLayout 更新布局模板
func (p *Page) UpdateLayout(layout string) error {
	if layout == "" {
		return errors.New("布局模板不能为空")
	}
	
	if p.layout == layout {
		return nil // 布局未变化
	}
	
	oldLayout := p.layout
	p.layout = layout
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "layout", oldLayout, layout))
	
	return nil
}

// UpdateContent 更新页面内容
func (p *Page) UpdateContent(content string) {
	if p.content == content {
		return // 内容未变化
	}
	
	p.content = content
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "content", "", ""))
}

// Publish 发布页面
func (p *Page) Publish() error {
	if !p.status.CanBePublished() {
		return errors.New("当前状态不允许发布")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewPublishedPageStatus()
	now := time.Now()
	p.publishedAt = &now
	p.updatedAt = now
	
	// 添加状态变更事件
	p.addDomainEvent(NewPageStatusChangedEvent(p, oldStatus, p.status, "发布页面"))
	
	return nil
}

// Unpublish 取消发布页面
func (p *Page) Unpublish() error {
	if !p.status.IsPublished() {
		return errors.New("页面未发布")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewDraftPageStatus()
	p.publishedAt = nil
	p.updatedAt = time.Now()
	
	// 添加状态变更事件
	p.addDomainEvent(NewPageStatusChangedEvent(p, oldStatus, p.status, "取消发布"))
	
	return nil
}

// SetAsPrivate 设置为私有页面
func (p *Page) SetAsPrivate() error {
	if !p.status.CanTransitionTo(valueobject.NewPrivatePageStatus()) {
		return errors.New("当前状态不允许设为私有")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewPrivatePageStatus()
	p.updatedAt = time.Now()
	
	// 添加状态变更事件
	p.addDomainEvent(NewPageStatusChangedEvent(p, oldStatus, p.status, "设为私有"))
	
	return nil
}

// Archive 归档页面
func (p *Page) Archive() error {
	if !p.status.CanBeArchived() {
		return errors.New("当前状态不允许归档")
	}
	
	oldStatus := p.status
	p.status = valueobject.PageStatusArchived
	p.updatedAt = time.Now()
	
	// 添加状态变更事件
	p.addDomainEvent(NewPageStatusChangedEvent(p, oldStatus, p.status, "归档页面"))
	
	return nil
}

// Delete 删除页面（软删除）
func (p *Page) Delete() error {
	if !p.status.CanBeDeleted() {
		return errors.New("当前状态不允许删除")
	}
	
	// 首页不能删除
	if p.isHomepage {
		return errors.New("首页不能删除")
	}
	
	oldStatus := p.status
	p.status = valueobject.PageStatusDeleted
	p.updatedAt = time.Now()
	
	// 添加删除事件
	p.addDomainEvent(NewPageDeletedEvent(p))
	
	return nil
}

// Restore 恢复页面
func (p *Page) Restore() error {
	if !p.status.CanBeRestored() {
		return errors.New("当前状态不允许恢复")
	}
	
	oldStatus := p.status
	p.status = valueobject.NewDraftPageStatus()
	p.updatedAt = time.Now()
	
	// 添加状态变更事件
	p.addDomainEvent(NewPageStatusChangedEvent(p, oldStatus, p.status, "恢复页面"))
	
	return nil
}

// SetAsHomepage 设置为首页
func (p *Page) SetAsHomepage() error {
	if p.isHomepage {
		return nil // 已经是首页
	}
	
	// 只有已发布的页面才能设为首页
	if !p.status.IsPublished() {
		return errors.New("只有已发布的页面才能设为首页")
	}
	
	p.isHomepage = true
	p.updatedAt = time.Now()
	
	// 添加首页设置事件
	p.addDomainEvent(NewPageHomepageSetEvent(p))
	
	return nil
}

// UnsetAsHomepage 取消首页设置
func (p *Page) UnsetAsHomepage() {
	if !p.isHomepage {
		return // 本来就不是首页
	}
	
	p.isHomepage = false
	p.updatedAt = time.Now()
	
	// 添加首页取消事件
	p.addDomainEvent(NewPageHomepageUnsetEvent(p))
}

// UpdateSortOrder 更新排序顺序
func (p *Page) UpdateSortOrder(sortOrder int32) {
	if p.sortOrder == sortOrder {
		return // 排序未变化
	}
	
	oldSortOrder := p.sortOrder
	p.sortOrder = sortOrder
	p.updatedAt = time.Now()
	
	// 添加排序更新事件
	p.addDomainEvent(NewPageSortOrderChangedEvent(p, oldSortOrder, sortOrder))
}

// IsVisible 检查页面是否对外可见
func (p *Page) IsVisible() bool {
	return p.status.IsVisible()
}

// IsEditable 检查页面是否可编辑
func (p *Page) IsEditable() bool {
	return p.status.IsEditable()
}

// CanChangeStatus 检查是否可以改变状态
func (p *Page) CanChangeStatus(targetStatus valueobject.PageStatus) bool {
	return p.status.CanTransitionTo(targetStatus)
}

// GetURL 获取页面完整URL
func (p *Page) GetURL() string {
	if p.isHomepage {
		return "/"
	}
	return p.slug.ToURL()
}

// GetSEOScore 获取页面SEO评分
func (p *Page) GetSEOScore() int {
	score := 0
	
	// 标题评分（30分）
	if p.title.IsSEOOptimized() {
		score += 30
	} else if !p.title.IsEmpty() {
		score += 15
	}
	
	// SEO元数据评分（40分）
	score += int(float64(p.seoMeta.GetSEOScore()) * 0.4)
	
	// URL段评分（20分）
	if p.slug.IsSEOFriendly() {
		score += 20
	} else if !p.slug.IsEmpty() {
		score += 10
	}
	
	// 内容评分（10分）
	if len(p.content) > 100 {
		score += 10
	} else if len(p.content) > 0 {
		score += 5
	}
	
	return score
}

// GetSEOSuggestions 获取SEO优化建议
func (p *Page) GetSEOSuggestions() []string {
	var suggestions []string
	
	// 标题建议
	if !p.title.IsSEOOptimized() {
		titleInfo := p.title.GetSuggestedLength()
		if suggestion, ok := titleInfo["suggestion"].(string); ok {
			suggestions = append(suggestions, suggestion)
		}
	}
	
	// SEO元数据建议
	seoSuggestions := p.seoMeta.GetSEOSuggestions()
	suggestions = append(suggestions, seoSuggestions...)
	
	// URL段建议
	if !p.slug.IsSEOFriendly() {
		suggestions = append(suggestions, "建议优化URL段，使其更符合SEO规范")
	}
	
	// 内容建议
	if len(p.content) < 100 {
		suggestions = append(suggestions, "建议增加页面内容，提高页面价值")
	}
	
	return suggestions
}

// GetDisplayInfo 获取显示信息
func (p *Page) GetDisplayInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":           p.id.Value(),
		"site_id":      p.siteID,
		"name":         p.name,
		"slug":         p.slug.Value(),
		"title":        p.title.GetUnescaped(),
		"status":       p.status.String(),
		"is_homepage":  p.isHomepage,
		"is_visible":   p.IsVisible(),
		"is_editable":  p.IsEditable(),
		"url":          p.GetURL(),
		"seo_score":    p.GetSEOScore(),
		"sort_order":   p.sortOrder,
		"published_at": p.publishedAt,
		"created_at":   p.createdAt,
		"updated_at":   p.updatedAt,
	}
}

// MatchesSearch 检查是否匹配搜索条件
func (p *Page) MatchesSearch(keyword string) bool {
	if keyword == "" {
		return true
	}
	
	keyword = strings.ToLower(keyword)
	
	// 检查名称
	if strings.Contains(strings.ToLower(p.name), keyword) {
		return true
	}
	
	// 检查标题
	if p.title.ContainsKeyword(keyword) {
		return true
	}
	
	// 检查URL段
	if strings.Contains(strings.ToLower(p.slug.Value()), keyword) {
		return true
	}
	
	// 检查SEO关键词
	if p.seoMeta.ContainsKeyword(keyword) {
		return true
	}
	
	return false
}

// GetDomainEvents 返回页面关联的所有领域事件
func (p *Page) GetDomainEvents() []event.DomainEvent {
	return p.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (p *Page) ClearDomainEvents() {
	p.domainEvents = p.domainEvents[:0]
}

// addDomainEvent 添加领域事件到页面
func (p *Page) addDomainEvent(event event.DomainEvent) {
	p.domainEvents = append(p.domainEvents, event)
} 