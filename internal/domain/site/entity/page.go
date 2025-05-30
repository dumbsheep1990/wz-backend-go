package entity

import (
	"errors"
	"time"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/site/valueobject"
)

// Page 页面实体
type Page struct {
	id          valueobject.PageID
	siteID      valueobject.SiteID
	title       valueobject.PageTitle
	slug        valueobject.PageSlug
	content     string
	metaTitle   string
	metaDesc    string
	isHomepage  bool
	sortOrder   int
	isPublished bool
	createdAt   time.Time
	updatedAt   time.Time
	
	// 领域事件
	domainEvents []event.DomainEvent
}

// NewPage 创建新页面
func NewPage(
	id valueobject.PageID,
	siteID valueobject.SiteID,
	title valueobject.PageTitle,
	slug valueobject.PageSlug,
	content string,
) (*Page, error) {
	if id.IsEmpty() {
		return nil, errors.New("页面ID不能为空")
	}
	if siteID.IsEmpty() {
		return nil, errors.New("站点ID不能为空")
	}
	if title.IsEmpty() {
		return nil, errors.New("页面标题不能为空")
	}
	if slug.IsEmpty() {
		return nil, errors.New("页面路径不能为空")
	}
	
	now := time.Now()
	page := &Page{
		id:           id,
		siteID:       siteID,
		title:        title,
		slug:         slug,
		content:      content,
		isHomepage:   false,
		sortOrder:    0,
		isPublished:  false,
		createdAt:    now,
		updatedAt:    now,
		domainEvents: make([]event.DomainEvent, 0),
	}
	
	// 添加页面创建事件
	page.addDomainEvent(NewPageCreatedEvent(page))
	
	return page, nil
}

// Getters
func (p *Page) ID() valueobject.PageID {
	return p.id
}

func (p *Page) SiteID() valueobject.SiteID {
	return p.siteID
}

func (p *Page) Title() valueobject.PageTitle {
	return p.title
}

func (p *Page) Slug() valueobject.PageSlug {
	return p.slug
}

func (p *Page) Content() string {
	return p.content
}

func (p *Page) MetaTitle() string {
	return p.metaTitle
}

func (p *Page) MetaDesc() string {
	return p.metaDesc
}

func (p *Page) IsHomepage() bool {
	return p.isHomepage
}

func (p *Page) SortOrder() int {
	return p.sortOrder
}

func (p *Page) IsPublished() bool {
	return p.isPublished
}

func (p *Page) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Page) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdateTitle 更新页面标题
func (p *Page) UpdateTitle(title valueobject.PageTitle) error {
	if title.IsEmpty() {
		return errors.New("页面标题不能为空")
	}
	
	if p.title.Equals(title) {
		return nil // 标题未变化
	}
	
	oldTitle := p.title
	p.title = title
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "title", oldTitle.Value(), title.Value()))
	
	return nil
}

// UpdateSlug 更新页面路径
func (p *Page) UpdateSlug(slug valueobject.PageSlug) error {
	if slug.IsEmpty() {
		return errors.New("页面路径不能为空")
	}
	
	if p.slug.Equals(slug) {
		return nil // 路径未变化
	}
	
	oldSlug := p.slug
	p.slug = slug
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "slug", oldSlug.Value(), slug.Value()))
	
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

// UpdateMeta 更新页面元信息
func (p *Page) UpdateMeta(metaTitle, metaDesc string) {
	changed := false
	
	if p.metaTitle != metaTitle {
		p.metaTitle = metaTitle
		changed = true
	}
	
	if p.metaDesc != metaDesc {
		p.metaDesc = metaDesc
		changed = true
	}
	
	if changed {
		p.updatedAt = time.Now()
		p.addDomainEvent(NewPageUpdatedEvent(p, "meta", "", ""))
	}
}

// SetAsHomepage 设置为首页
func (p *Page) SetAsHomepage() {
	if p.isHomepage {
		return // 已经是首页
	}
	
	p.isHomepage = true
	p.updatedAt = time.Now()
	
	// 添加首页设置事件
	p.addDomainEvent(NewPageSetAsHomepageEvent(p))
}

// UnsetAsHomepage 取消首页设置
func (p *Page) UnsetAsHomepage() {
	if !p.isHomepage {
		return // 本来就不是首页
	}
	
	p.isHomepage = false
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "isHomepage", "true", "false"))
}

// UpdateSortOrder 更新排序
func (p *Page) UpdateSortOrder(sortOrder int) {
	if p.sortOrder == sortOrder {
		return // 排序未变化
	}
	
	oldOrder := p.sortOrder
	p.sortOrder = sortOrder
	p.updatedAt = time.Now()
	
	// 添加页面更新事件
	p.addDomainEvent(NewPageUpdatedEvent(p, "sortOrder", string(rune(oldOrder)), string(rune(sortOrder))))
}

// Publish 发布页面
func (p *Page) Publish() {
	if p.isPublished {
		return // 已经发布
	}
	
	p.isPublished = true
	p.updatedAt = time.Now()
	
	// 添加页面发布事件
	p.addDomainEvent(NewPagePublishedEvent(p))
}

// Unpublish 取消发布页面
func (p *Page) Unpublish() {
	if !p.isPublished {
		return // 本来就未发布
	}
	
	p.isPublished = false
	p.updatedAt = time.Now()
	
	// 添加页面取消发布事件
	p.addDomainEvent(NewPageUnpublishedEvent(p))
}

// CanBeDeleted 是否可以删除
func (p *Page) CanBeDeleted() bool {
	// 首页不能删除
	return !p.isHomepage
}

// GetDomainEvents 获取领域事件
func (p *Page) GetDomainEvents() []event.DomainEvent {
	return p.domainEvents
}

// ClearDomainEvents 清除领域事件
func (p *Page) ClearDomainEvents() {
	p.domainEvents = make([]event.DomainEvent, 0)
}

// addDomainEvent 添加领域事件
func (p *Page) addDomainEvent(domainEvent event.DomainEvent) {
	p.domainEvents = append(p.domainEvents, domainEvent)
} 