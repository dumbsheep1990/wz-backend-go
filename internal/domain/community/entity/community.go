package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/yourusername/wz-backend-go/internal/domain/community/valueobject"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
)

// Community 表示领域中的社区实体
type Community struct {
	id          valueobject.ID
	name        valueobject.CommunityName
	description valueobject.Description
	ownerID     valueobject.UserID
	status      valueobject.CommunityStatus
	tags        []valueobject.Tag
	location    valueobject.Location
	createdAt   time.Time
	updatedAt   time.Time
	
	domainEvents []event.DomainEvent
}

// NewCommunity 创建一个新的社区实体
func NewCommunity(
	name valueobject.CommunityName,
	description valueobject.Description,
	ownerID valueobject.UserID,
	tags []valueobject.Tag,
	location valueobject.Location,
) (*Community, error) {
	if name.IsEmpty() {
		return nil, errors.New("社区名称不能为空")
	}

	if ownerID.IsEmpty() {
		return nil, errors.New("创建者ID不能为空")
	}

	now := time.Now()
	community := &Community{
		id:          valueobject.NewID(uuid.New().String()),
		name:        name,
		description: description,
		ownerID:     ownerID,
		status:      valueobject.StatusActive,
		tags:        tags,
		location:    location,
		createdAt:   now,
		updatedAt:   now,
	}

	community.addDomainEvent(NewCommunityCreatedEvent(community))
	return community, nil
}

// ID 返回社区ID
func (c *Community) ID() valueobject.ID {
	return c.id
}

// Name 返回社区名称
func (c *Community) Name() valueobject.CommunityName {
	return c.name
}

// Description 返回社区描述
func (c *Community) Description() valueobject.Description {
	return c.description
}

// OwnerID 返回社区创建者ID
func (c *Community) OwnerID() valueobject.UserID {
	return c.ownerID
}

// Status 返回社区状态
func (c *Community) Status() valueobject.CommunityStatus {
	return c.status
}

// Tags 返回社区标签
func (c *Community) Tags() []valueobject.Tag {
	return c.tags
}

// Location 返回社区位置
func (c *Community) Location() valueobject.Location {
	return c.location
}

// CreatedAt 返回创建时间
func (c *Community) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt 返回最后更新时间
func (c *Community) UpdatedAt() time.Time {
	return c.updatedAt
}

// UpdateName 更新社区名称
func (c *Community) UpdateName(name valueobject.CommunityName) error {
	if name.IsEmpty() {
		return errors.New("社区名称不能为空")
	}

	if c.status == valueobject.StatusDeleted {
		return errors.New("无法更新已删除的社区")
	}

	c.name = name
	c.updatedAt = time.Now()
	c.addDomainEvent(NewCommunityUpdatedEvent(c))
	return nil
}

// UpdateDescription 更新社区描述
func (c *Community) UpdateDescription(description valueobject.Description) error {
	if c.status == valueobject.StatusDeleted {
		return errors.New("无法更新已删除的社区")
	}

	c.description = description
	c.updatedAt = time.Now()
	c.addDomainEvent(NewCommunityUpdatedEvent(c))
	return nil
}

// UpdateTags 更新社区标签
func (c *Community) UpdateTags(tags []valueobject.Tag) error {
	if c.status == valueobject.StatusDeleted {
		return errors.New("无法更新已删除的社区")
	}

	c.tags = tags
	c.updatedAt = time.Now()
	c.addDomainEvent(NewCommunityUpdatedEvent(c))
	return nil
}

// UpdateLocation 更新社区位置
func (c *Community) UpdateLocation(location valueobject.Location) error {
	if c.status == valueobject.StatusDeleted {
		return errors.New("无法更新已删除的社区")
	}

	c.location = location
	c.updatedAt = time.Now()
	c.addDomainEvent(NewCommunityUpdatedEvent(c))
	return nil
}

// Delete 标记社区为已删除
func (c *Community) Delete() error {
	if c.status == valueobject.StatusDeleted {
		return errors.New("社区已经被删除")
	}

	c.status = valueobject.StatusDeleted
	c.updatedAt = time.Now()
	c.addDomainEvent(NewCommunityDeletedEvent(c))
	return nil
}

// Activate 标记社区为活跃状态
func (c *Community) Activate() error {
	if c.status == valueobject.StatusActive {
		return errors.New("社区已经是活跃状态")
	}

	c.status = valueobject.StatusActive
	c.updatedAt = time.Now()
	c.addDomainEvent(NewCommunityActivatedEvent(c))
	return nil
}

// addDomainEvent 添加领域事件到社区
func (c *Community) addDomainEvent(event event.DomainEvent) {
	c.domainEvents = append(c.domainEvents, event)
}

// GetDomainEvents 返回与社区关联的所有领域事件
func (c *Community) GetDomainEvents() []event.DomainEvent {
	return c.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (c *Community) ClearDomainEvents() {
	c.domainEvents = []event.DomainEvent{}
}
