package entity

import (
	"errors"
	"time"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/site/valueobject"
)

// Site 站点聚合根
type Site struct {
	id          valueobject.SiteID
	name        valueobject.SiteName
	description string
	domain      valueobject.Domain
	logo        string
	favicon     string
	tenantID    string // 租户ID - 应该也是值对象，但为简化暂用string
	theme       valueobject.ThemeConfig
	status      valueobject.SiteStatus
	thumbnail   string
	createdAt   time.Time
	updatedAt   time.Time
	publishedAt *time.Time
	
	// 领域事件
	domainEvents []event.DomainEvent
}

// NewSite 创建新站点
func NewSite(
	id valueobject.SiteID,
	name valueobject.SiteName,
	description string,
	domain valueobject.Domain,
	tenantID string,
) (*Site, error) {
	if id.IsEmpty() {
		return nil, errors.New("站点ID不能为空")
	}
	if name.IsEmpty() {
		return nil, errors.New("站点名称不能为空")
	}
	if tenantID == "" {
		return nil, errors.New("租户ID不能为空")
	}
	
	now := time.Now()
	site := &Site{
		id:           id,
		name:         name,
		description:  description,
		domain:       domain,
		tenantID:     tenantID,
		theme:        valueobject.NewDefaultThemeConfig(),
		status:       valueobject.NewDraftStatus(),
		createdAt:    now,
		updatedAt:    now,
		domainEvents: make([]event.DomainEvent, 0),
	}
	
	// 添加站点创建事件
	site.addDomainEvent(NewSiteCreatedEvent(site))
	
	return site, nil
}

// Getters
func (s *Site) ID() valueobject.SiteID {
	return s.id
}

func (s *Site) Name() valueobject.SiteName {
	return s.name
}

func (s *Site) Description() string {
	return s.description
}

func (s *Site) Domain() valueobject.Domain {
	return s.domain
}

func (s *Site) Logo() string {
	return s.logo
}

func (s *Site) Favicon() string {
	return s.favicon
}

func (s *Site) TenantID() string {
	return s.tenantID
}

func (s *Site) Theme() valueobject.ThemeConfig {
	return s.theme
}

func (s *Site) Status() valueobject.SiteStatus {
	return s.status
}

func (s *Site) Thumbnail() string {
	return s.thumbnail
}

func (s *Site) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Site) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Site) PublishedAt() *time.Time {
	return s.publishedAt
}

// UpdateName 更新站点名称
func (s *Site) UpdateName(name valueobject.SiteName) error {
	if name.IsEmpty() {
		return errors.New("站点名称不能为空")
	}
	
	if s.name.Equals(name) {
		return nil // 名称未变化
	}
	
	oldName := s.name
	s.name = name
	s.updatedAt = time.Now()
	
	// 添加站点更新事件
	s.addDomainEvent(NewSiteUpdatedEvent(s, "name", oldName.Value(), name.Value()))
	
	return nil
}

// UpdateDescription 更新描述
func (s *Site) UpdateDescription(description string) {
	if s.description == description {
		return // 描述未变化
	}
	
	oldDescription := s.description
	s.description = description
	s.updatedAt = time.Now()
	
	// 添加站点更新事件
	s.addDomainEvent(NewSiteUpdatedEvent(s, "description", oldDescription, description))
}

// UpdateDomain 更新域名
func (s *Site) UpdateDomain(domain valueobject.Domain) error {
	if s.domain.Equals(domain) {
		return nil // 域名未变化
	}
	
	oldDomain := s.domain
	s.domain = domain
	s.updatedAt = time.Now()
	
	// 添加站点更新事件
	s.addDomainEvent(NewSiteUpdatedEvent(s, "domain", oldDomain.Value(), domain.Value()))
	
	return nil
}

// UpdateTheme 更新主题配置
func (s *Site) UpdateTheme(theme valueobject.ThemeConfig) {
	if s.theme.Equals(theme) {
		return // 主题未变化
	}
	
	s.theme = theme
	s.updatedAt = time.Now()
	
	// 添加站点更新事件
	s.addDomainEvent(NewSiteUpdatedEvent(s, "theme", "", ""))
}

// UpdateLogo 更新Logo
func (s *Site) UpdateLogo(logo string) {
	if s.logo == logo {
		return // Logo未变化
	}
	
	oldLogo := s.logo
	s.logo = logo
	s.updatedAt = time.Now()
	
	// 添加站点更新事件
	s.addDomainEvent(NewSiteUpdatedEvent(s, "logo", oldLogo, logo))
}

// UpdateFavicon 更新Favicon
func (s *Site) UpdateFavicon(favicon string) {
	if s.favicon == favicon {
		return // Favicon未变化
	}
	
	oldFavicon := s.favicon
	s.favicon = favicon
	s.updatedAt = time.Now()
	
	// 添加站点更新事件
	s.addDomainEvent(NewSiteUpdatedEvent(s, "favicon", oldFavicon, favicon))
}

// UpdateThumbnail 更新缩略图
func (s *Site) UpdateThumbnail(thumbnail string) {
	if s.thumbnail == thumbnail {
		return // 缩略图未变化
	}
	
	oldThumbnail := s.thumbnail
	s.thumbnail = thumbnail
	s.updatedAt = time.Now()
	
	// 添加站点更新事件
	s.addDomainEvent(NewSiteUpdatedEvent(s, "thumbnail", oldThumbnail, thumbnail))
}

// Publish 发布站点
func (s *Site) Publish() error {
	if !s.status.CanPublish() {
		return errors.New("当前状态下不能发布站点")
	}
	
	now := time.Now()
	s.status = valueobject.NewPublishedStatus()
	s.publishedAt = &now
	s.updatedAt = now
	
	// 添加站点发布事件
	s.addDomainEvent(NewSitePublishedEvent(s))
	
	return nil
}

// Archive 归档站点
func (s *Site) Archive() error {
	if !s.status.CanArchive() {
		return errors.New("当前状态下不能归档站点")
	}
	
	s.status = valueobject.NewArchivedStatus()
	s.updatedAt = time.Now()
	
	// 添加站点归档事件
	s.addDomainEvent(NewSiteArchivedEvent(s))
	
	return nil
}

// IsOwnedBy 检查是否属于指定租户
func (s *Site) IsOwnedBy(tenantID string) bool {
	return s.tenantID == tenantID
}

// CanBeModified 是否可以修改
func (s *Site) CanBeModified() bool {
	return !s.status.IsArchived()
}

// CanBeDeleted 是否可以删除
func (s *Site) CanBeDeleted() bool {
	return s.status.IsDraft()
}

// GetDomainEvents 获取领域事件
func (s *Site) GetDomainEvents() []event.DomainEvent {
	return s.domainEvents
}

// ClearDomainEvents 清除领域事件
func (s *Site) ClearDomainEvents() {
	s.domainEvents = make([]event.DomainEvent, 0)
}

// addDomainEvent 添加领域事件
func (s *Site) addDomainEvent(domainEvent event.DomainEvent) {
	s.domainEvents = append(s.domainEvents, domainEvent)
} 