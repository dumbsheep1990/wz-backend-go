package entity

import (
	"errors"
	"time"
	"wz-backend-go/internal/domain/component/valueobject"
	"wz-backend-go/internal/domain/shared/event"
)

// Component 组件聚合根
type Component struct {
	id          valueobject.ComponentID
	name        string
	description string
	componentType valueobject.ComponentType
	template    string  // 组件模板内容
	config      string  // 组件配置（JSON格式）
	preview     string  // 预览图片URL
	category    string  // 组件分类
	tags        []string // 标签
	isPublic    bool    // 是否公开
	tenantID    string  // 租户ID
	version     string  // 版本号
	createdAt   time.Time
	updatedAt   time.Time
	
	// 领域事件
	domainEvents []event.DomainEvent
}

// NewComponent 创建新组件
func NewComponent(
	id valueobject.ComponentID,
	name string,
	componentType valueobject.ComponentType,
	template string,
	tenantID string,
) (*Component, error) {
	if id.IsEmpty() {
		return nil, errors.New("组件ID不能为空")
	}
	if name == "" {
		return nil, errors.New("组件名称不能为空")
	}
	if template == "" {
		return nil, errors.New("组件模板不能为空")
	}
	if tenantID == "" {
		return nil, errors.New("租户ID不能为空")
	}
	
	now := time.Now()
	component := &Component{
		id:            id,
		name:          name,
		componentType: componentType,
		template:      template,
		tenantID:      tenantID,
		version:       "1.0.0",
		isPublic:      false,
		createdAt:     now,
		updatedAt:     now,
		domainEvents:  make([]event.DomainEvent, 0),
	}
	
	// 添加组件创建事件
	component.addDomainEvent(NewComponentCreatedEvent(component))
	
	return component, nil
}

// Getters
func (c *Component) ID() valueobject.ComponentID {
	return c.id
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) Description() string {
	return c.description
}

func (c *Component) ComponentType() valueobject.ComponentType {
	return c.componentType
}

func (c *Component) Template() string {
	return c.template
}

func (c *Component) Config() string {
	return c.config
}

func (c *Component) Preview() string {
	return c.preview
}

func (c *Component) Category() string {
	return c.category
}

func (c *Component) Tags() []string {
	return c.tags
}

func (c *Component) IsPublic() bool {
	return c.isPublic
}

func (c *Component) TenantID() string {
	return c.tenantID
}

func (c *Component) Version() string {
	return c.version
}

func (c *Component) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Component) UpdatedAt() time.Time {
	return c.updatedAt
}

// UpdateName 更新组件名称
func (c *Component) UpdateName(name string) error {
	if name == "" {
		return errors.New("组件名称不能为空")
	}
	
	if c.name == name {
		return nil // 名称未变化
	}
	
	oldName := c.name
	c.name = name
	c.updatedAt = time.Now()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "name", oldName, name))
	
	return nil
}

// UpdateDescription 更新描述
func (c *Component) UpdateDescription(description string) {
	if c.description == description {
		return // 描述未变化
	}
	
	oldDescription := c.description
	c.description = description
	c.updatedAt = time.Now()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "description", oldDescription, description))
}

// UpdateTemplate 更新模板
func (c *Component) UpdateTemplate(template string) error {
	if template == "" {
		return errors.New("组件模板不能为空")
	}
	
	if c.template == template {
		return nil // 模板未变化
	}
	
	c.template = template
	c.updatedAt = time.Now()
	
	// 模板更新需要增加版本号
	c.incrementVersion()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "template", "", ""))
	
	return nil
}

// UpdateConfig 更新配置
func (c *Component) UpdateConfig(config string) {
	if c.config == config {
		return // 配置未变化
	}
	
	c.config = config
	c.updatedAt = time.Now()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "config", "", ""))
}

// UpdatePreview 更新预览图
func (c *Component) UpdatePreview(preview string) {
	if c.preview == preview {
		return // 预览图未变化
	}
	
	oldPreview := c.preview
	c.preview = preview
	c.updatedAt = time.Now()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "preview", oldPreview, preview))
}

// UpdateCategory 更新分类
func (c *Component) UpdateCategory(category string) {
	if c.category == category {
		return // 分类未变化
	}
	
	oldCategory := c.category
	c.category = category
	c.updatedAt = time.Now()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "category", oldCategory, category))
}

// UpdateTags 更新标签
func (c *Component) UpdateTags(tags []string) {
	c.tags = tags
	c.updatedAt = time.Now()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "tags", "", ""))
}

// MakePublic 设为公开
func (c *Component) MakePublic() {
	if c.isPublic {
		return // 已经是公开的
	}
	
	c.isPublic = true
	c.updatedAt = time.Now()
	
	// 添加组件公开事件
	c.addDomainEvent(NewComponentMadePublicEvent(c))
}

// MakePrivate 设为私有
func (c *Component) MakePrivate() {
	if !c.isPublic {
		return // 已经是私有的
	}
	
	c.isPublic = false
	c.updatedAt = time.Now()
	
	// 添加组件更新事件
	c.addDomainEvent(NewComponentUpdatedEvent(c, "isPublic", "true", "false"))
}

// IsOwnedBy 检查是否属于指定租户
func (c *Component) IsOwnedBy(tenantID string) bool {
	return c.tenantID == tenantID
}

// CanBeAccessed 检查是否可以被访问
func (c *Component) CanBeAccessed(tenantID string) bool {
	return c.isPublic || c.tenantID == tenantID
}

// incrementVersion 增加版本号（简单实现）
func (c *Component) incrementVersion() {
	// 简单的版本号递增逻辑，实际应该使用语义化版本控制
	switch c.version {
	case "1.0.0":
		c.version = "1.0.1"
	case "1.0.1":
		c.version = "1.0.2"
	default:
		c.version = "1.0.1"
	}
}

// GetDomainEvents 获取领域事件
func (c *Component) GetDomainEvents() []event.DomainEvent {
	return c.domainEvents
}

// ClearDomainEvents 清除领域事件
func (c *Component) ClearDomainEvents() {
	c.domainEvents = make([]event.DomainEvent, 0)
}

// addDomainEvent 添加领域事件
func (c *Component) addDomainEvent(domainEvent event.DomainEvent) {
	c.domainEvents = append(c.domainEvents, domainEvent)
} 