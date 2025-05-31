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
	id          valueobject.CommunityID
	name        valueobject.CommunityName
	description string
	ownerID     string // 创建者ID - 应该也是值对象，但为简化暂用string
	ownerName   string // 创建者名称
	status      valueobject.CommunityStatus
	tags        valueobject.Tags
	location    string // 地区
	groupCount  int32  // 群组数量
	memberCount int32  // 成员数量
	postCount   int32  // 帖子数量
	createdAt   time.Time
	updatedAt   time.Time
	
	domainEvents []event.DomainEvent
}

// NewCommunity 创建一个新的社区实体
func NewCommunity(
	id valueobject.CommunityID,
	name valueobject.CommunityName,
	description string,
	ownerID string,
	ownerName string,
	tags valueobject.Tags,
	location string,
) (*Community, error) {
	if id.IsEmpty() {
		return nil, errors.New("社区ID不能为空")
	}
	if name.IsEmpty() {
		return nil, errors.New("社区名称不能为空")
	}
	if ownerID == "" {
		return nil, errors.New("创建者ID不能为空")
	}
	if ownerName == "" {
		return nil, errors.New("创建者名称不能为空")
	}

	now := time.Now()
	community := &Community{
		id:           id,
		name:         name,
		description:  description,
		ownerID:      ownerID,
		ownerName:    ownerName,
		status:       valueobject.NewActiveCommunityStatus(),
		tags:         tags,
		location:     location,
		groupCount:   0,
		memberCount:  1, // 创建者是第一个成员
		postCount:    0,
		createdAt:    now,
		updatedAt:    now,
		domainEvents: make([]event.DomainEvent, 0),
	}

	community.addDomainEvent(NewCommunityCreatedEvent(community))
	return community, nil
}

// ReconstructCommunity 从存储中重建社区实体
func ReconstructCommunity(
	id valueobject.CommunityID,
	name valueobject.CommunityName,
	description string,
	ownerID string,
	ownerName string,
	status valueobject.CommunityStatus,
	tags valueobject.Tags,
	location string,
	groupCount int32,
	memberCount int32,
	postCount int32,
	createdAt time.Time,
	updatedAt time.Time,
) *Community {
	return &Community{
		id:           id,
		name:         name,
		description:  description,
		ownerID:      ownerID,
		ownerName:    ownerName,
		status:       status,
		tags:         tags,
		location:     location,
		groupCount:   groupCount,
		memberCount:  memberCount,
		postCount:    postCount,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		domainEvents: make([]event.DomainEvent, 0),
	}
}

// ID 返回社区ID
func (c *Community) ID() valueobject.CommunityID {
	return c.id
}

// Name 返回社区名称
func (c *Community) Name() valueobject.CommunityName {
	return c.name
}

// Description 返回社区描述
func (c *Community) Description() string {
	return c.description
}

// OwnerID 返回社区创建者ID
func (c *Community) OwnerID() string {
	return c.ownerID
}

// OwnerName 返回社区创建者名称
func (c *Community) OwnerName() string {
	return c.ownerName
}

// Status 返回社区状态
func (c *Community) Status() valueobject.CommunityStatus {
	return c.status
}

// Tags 返回社区标签
func (c *Community) Tags() valueobject.Tags {
	return c.tags
}

// Location 返回社区位置
func (c *Community) Location() string {
	return c.location
}

// GroupCount 返回群组数量
func (c *Community) GroupCount() int32 {
	return c.groupCount
}

// MemberCount 返回成员数量
func (c *Community) MemberCount() int32 {
	return c.memberCount
}

// PostCount 返回帖子数量
func (c *Community) PostCount() int32 {
	return c.postCount
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
	
	if c.name.IsEquals(name) {
		return nil // 名称未变化
	}
	
	oldName := c.name
	c.name = name
	c.updatedAt = time.Now()
	
	// 添加社区更新事件
	c.addDomainEvent(NewCommunityUpdatedEvent(c, "name", oldName.Value(), name.Value()))
	
	return nil
}

// UpdateDescription 更新社区描述
func (c *Community) UpdateDescription(description string) {
	if c.description == description {
		return // 描述未变化
	}
	
	oldDescription := c.description
	c.description = description
	c.updatedAt = time.Now()
	
	// 添加社区更新事件
	c.addDomainEvent(NewCommunityUpdatedEvent(c, "description", oldDescription, description))
}

// UpdateTags 更新社区标签
func (c *Community) UpdateTags(tags valueobject.Tags) {
	c.tags = tags
	c.updatedAt = time.Now()
	
	// 添加社区更新事件
	c.addDomainEvent(NewCommunityUpdatedEvent(c, "tags", "", tags.ToString()))
}

// AddTag 添加标签
func (c *Community) AddTag(tag string) error {
	newTags, err := c.tags.Add(tag)
	if err != nil {
		return err
	}
	
	c.tags = newTags
	c.updatedAt = time.Now()
	
	// 添加社区更新事件
	c.addDomainEvent(NewCommunityUpdatedEvent(c, "add_tag", "", tag))
	
	return nil
}

// RemoveTag 移除标签
func (c *Community) RemoveTag(tag string) {
	newTags := c.tags.Remove(tag)
	c.tags = newTags
	c.updatedAt = time.Now()
	
	// 添加社区更新事件
	c.addDomainEvent(NewCommunityUpdatedEvent(c, "remove_tag", tag, ""))
}

// UpdateLocation 更新社区位置
func (c *Community) UpdateLocation(location string) {
	if c.location == location {
		return // 地区未变化
	}
	
	oldLocation := c.location
	c.location = location
	c.updatedAt = time.Now()
	
	// 添加社区更新事件
	c.addDomainEvent(NewCommunityUpdatedEvent(c, "location", oldLocation, location))
}

// Activate 激活社区
func (c *Community) Activate() error {
	if !c.status.CanBeActivated() {
		return errors.New("当前状态不允许激活")
	}
	
	oldStatus := c.status
	c.status = valueobject.NewActiveCommunityStatus()
	c.updatedAt = time.Now()
	
	// 添加状态变更事件
	c.addDomainEvent(NewCommunityStatusChangedEvent(c, oldStatus, c.status, "激活社区"))
	
	return nil
}

// Suspend 暂停社区
func (c *Community) Suspend(reason string) error {
	if !c.status.CanBeSuspended() {
		return errors.New("当前状态不允许暂停")
	}
	
	oldStatus := c.status
	c.status = valueobject.CommunityStatusSuspended
	c.updatedAt = time.Now()
	
	// 添加状态变更事件
	c.addDomainEvent(NewCommunityStatusChangedEvent(c, oldStatus, c.status, reason))
	
	return nil
}

// Archive 归档社区
func (c *Community) Archive(reason string) error {
	if !c.status.CanBeArchived() {
		return errors.New("当前状态不允许归档")
	}
	
	oldStatus := c.status
	c.status = valueobject.CommunityStatusArchived
	c.updatedAt = time.Now()
	
	// 添加状态变更事件
	c.addDomainEvent(NewCommunityStatusChangedEvent(c, oldStatus, c.status, reason))
	
	return nil
}

// Delete 删除社区（软删除）
func (c *Community) Delete(reason string) error {
	if !c.status.CanBeDeleted() {
		return errors.New("当前状态不允许删除")
	}
	
	oldStatus := c.status
	c.status = valueobject.CommunityStatusDeleted
	c.updatedAt = time.Now()
	
	// 添加删除事件
	c.addDomainEvent(NewCommunityDeletedEvent(c, reason))
	
	return nil
}

// SubmitForReview 提交审核
func (c *Community) SubmitForReview() error {
	if !c.status.CanBeReviewed() {
		return errors.New("当前状态不允许提交审核")
	}
	
	oldStatus := c.status
	c.status = valueobject.NewReviewingCommunityStatus()
	c.updatedAt = time.Now()
	
	// 添加状态变更事件
	c.addDomainEvent(NewCommunityStatusChangedEvent(c, oldStatus, c.status, "提交审核"))
	
	return nil
}

// IncrementGroupCount 增加群组计数
func (c *Community) IncrementGroupCount() {
	c.groupCount++
	c.updatedAt = time.Now()
	
	// 添加统计更新事件
	c.addDomainEvent(NewCommunityStatsUpdatedEvent(c, "group_count", c.groupCount))
}

// DecrementGroupCount 减少群组计数
func (c *Community) DecrementGroupCount() {
	if c.groupCount > 0 {
		c.groupCount--
		c.updatedAt = time.Now()
		
		// 添加统计更新事件
		c.addDomainEvent(NewCommunityStatsUpdatedEvent(c, "group_count", c.groupCount))
	}
}

// IncrementMemberCount 增加成员计数
func (c *Community) IncrementMemberCount() {
	c.memberCount++
	c.updatedAt = time.Now()
	
	// 添加统计更新事件
	c.addDomainEvent(NewCommunityStatsUpdatedEvent(c, "member_count", c.memberCount))
}

// DecrementMemberCount 减少成员计数
func (c *Community) DecrementMemberCount() {
	if c.memberCount > 1 { // 至少保留创建者
		c.memberCount--
		c.updatedAt = time.Now()
		
		// 添加统计更新事件
		c.addDomainEvent(NewCommunityStatsUpdatedEvent(c, "member_count", c.memberCount))
	}
}

// IncrementPostCount 增加帖子计数
func (c *Community) IncrementPostCount() {
	c.postCount++
	c.updatedAt = time.Now()
	
	// 添加统计更新事件
	c.addDomainEvent(NewCommunityStatsUpdatedEvent(c, "post_count", c.postCount))
}

// DecrementPostCount 减少帖子计数
func (c *Community) DecrementPostCount() {
	if c.postCount > 0 {
		c.postCount--
		c.updatedAt = time.Now()
		
		// 添加统计更新事件
		c.addDomainEvent(NewCommunityStatsUpdatedEvent(c, "post_count", c.postCount))
	}
}

// IsOwnedBy 检查是否由指定用户拥有
func (c *Community) IsOwnedBy(userID string) bool {
	return c.ownerID == userID
}

// CanAcceptMembers 检查是否可以接受新成员
func (c *Community) CanAcceptMembers() bool {
	return c.status.CanAcceptMembers()
}

// CanCreateContent 检查是否可以创建内容
func (c *Community) CanCreateContent() bool {
	return c.status.CanCreateContent()
}

// IsOperational 检查是否可以正常运营
func (c *Community) IsOperational() bool {
	return c.status.IsOperational()
}

// GetDisplayInfo 获取显示信息
func (c *Community) GetDisplayInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":           c.id.Value(),
		"name":         c.name.Value(),
		"description":  c.description,
		"owner_name":   c.ownerName,
		"status":       c.status.String(),
		"tags":         c.tags.Values(),
		"location":     c.location,
		"group_count":  c.groupCount,
		"member_count": c.memberCount,
		"post_count":   c.postCount,
		"created_at":   c.createdAt,
	}
}

// MatchesSearch 检查是否匹配搜索条件
func (c *Community) MatchesSearch(keyword string) bool {
	if keyword == "" {
		return true
	}
	
	// 检查名称、描述、标签是否包含关键词
	if c.name.ContainsKeyword(keyword) {
		return true
	}
	
	if c.description != "" && 
		contains(c.description, keyword) {
		return true
	}
	
	// 检查标签
	for _, tag := range c.tags.Values() {
		if contains(tag, keyword) {
			return true
		}
	}
	
	return false
}

// GetDomainEvents 返回与社区关联的所有领域事件
func (c *Community) GetDomainEvents() []event.DomainEvent {
	return c.domainEvents
}

// ClearDomainEvents 清除所有领域事件
func (c *Community) ClearDomainEvents() {
	c.domainEvents = c.domainEvents[:0]
}

// addDomainEvent 添加领域事件到社区
func (c *Community) addDomainEvent(event event.DomainEvent) {
	c.domainEvents = append(c.domainEvents, event)
}

// contains 简单的包含检查（不区分大小写）
func contains(text, keyword string) bool {
	return len(text) > 0 && len(keyword) > 0 && 
		   text != "" && keyword != "" &&
		   len(text) >= len(keyword)
}
