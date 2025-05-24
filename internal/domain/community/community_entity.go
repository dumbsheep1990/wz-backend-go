package community

import (
	"errors"
	"time"
)

// Community 社区聚合根
type Community struct {
	id          ID
	name        CommunityName
	description Description
	ownerID     ID
	ownerName   string
	status      CommunityStatus
	tags        Tags
	location    Location
	type_       CommunityType
	createdAt   Timestamp
	updatedAt   Timestamp
	
	// 聚合计数器
	groupCount  int
	memberCount int
	postCount   int
}

// NewCommunity 创建社区
func NewCommunity(
	name CommunityName,
	description Description,
	ownerID ID,
	ownerName string,
	tags Tags,
	location Location,
	type_ CommunityType,
) (*Community, error) {
	// 验证社区类型
	if !type_.IsValid() {
		return nil, errors.New("无效的社区类型")
	}
	
	now := NewTimestamp()
	
	return &Community{
		id:          NewID(),
		name:        name,
		description: description,
		ownerID:     ownerID,
		ownerName:   ownerName,
		status:      CommunityStatusActive,
		tags:        tags,
		location:    location,
		type_:       type_,
		createdAt:   now,
		updatedAt:   now,
		groupCount:  0,
		memberCount: 0,
		postCount:   0,
	}, nil
}

// ReconstructCommunity 从数据库重构社区实体
func ReconstructCommunity(
	id string,
	name string,
	description string,
	ownerID string,
	ownerName string,
	status int,
	tags []string,
	location string,
	type_ string,
	createdAt time.Time,
	updatedAt time.Time,
	groupCount int,
	memberCount int,
	postCount int,
) (*Community, error) {
	communityID, err := NewIDFromString(id)
	if err != nil {
		return nil, err
	}
	
	communityName, err := NewCommunityName(name)
	if err != nil {
		return nil, err
	}
	
	communityDesc, err := NewDescription(description)
	if err != nil {
		return nil, err
	}
	
	ownerIDObj, err := NewIDFromString(ownerID)
	if err != nil {
		return nil, err
	}
	
	communityTags, err := NewTags(tags)
	if err != nil {
		return nil, err
	}
	
	communityLocation, err := NewLocation(location)
	if err != nil {
		return nil, err
	}
	
	communityType := CommunityType(type_)
	if !communityType.IsValid() {
		return nil, errors.New("无效的社区类型")
	}
	
	return &Community{
		id:          communityID,
		name:        communityName,
		description: communityDesc,
		ownerID:     ownerIDObj,
		ownerName:   ownerName,
		status:      CommunityStatus(status),
		tags:        communityTags,
		location:    communityLocation,
		type_:       communityType,
		createdAt:   NewTimestampFromTime(createdAt),
		updatedAt:   NewTimestampFromTime(updatedAt),
		groupCount:  groupCount,
		memberCount: memberCount,
		postCount:   postCount,
	}, nil
}

// ID 获取社区ID
func (c *Community) ID() ID {
	return c.id
}

// Name 获取社区名称
func (c *Community) Name() CommunityName {
	return c.name
}

// Description 获取社区描述
func (c *Community) Description() Description {
	return c.description
}

// OwnerID 获取拥有者ID
func (c *Community) OwnerID() ID {
	return c.ownerID
}

// OwnerName 获取拥有者名称
func (c *Community) OwnerName() string {
	return c.ownerName
}

// Status 获取社区状态
func (c *Community) Status() CommunityStatus {
	return c.status
}

// Tags 获取社区标签
func (c *Community) Tags() Tags {
	return c.tags
}

// Location 获取社区位置
func (c *Community) Location() Location {
	return c.location
}

// Type 获取社区类型
func (c *Community) Type() CommunityType {
	return c.type_
}

// CreatedAt 获取创建时间
func (c *Community) CreatedAt() Timestamp {
	return c.createdAt
}

// UpdatedAt 获取更新时间
func (c *Community) UpdatedAt() Timestamp {
	return c.updatedAt
}

// GroupCount 获取小组数量
func (c *Community) GroupCount() int {
	return c.groupCount
}

// MemberCount 获取成员数量
func (c *Community) MemberCount() int {
	return c.memberCount
}

// PostCount 获取帖子数量
func (c *Community) PostCount() int {
	return c.postCount
}

// UpdateName 更新社区名称
func (c *Community) UpdateName(name CommunityName) {
	c.name = name
	c.updatedAt = NewTimestamp()
}

// UpdateDescription 更新社区描述
func (c *Community) UpdateDescription(description Description) {
	c.description = description
	c.updatedAt = NewTimestamp()
}

// UpdateTags 更新社区标签
func (c *Community) UpdateTags(tags Tags) {
	c.tags = tags
	c.updatedAt = NewTimestamp()
}

// UpdateLocation 更新社区位置
func (c *Community) UpdateLocation(location Location) {
	c.location = location
	c.updatedAt = NewTimestamp()
}

// UpdateStatus 更新社区状态
func (c *Community) UpdateStatus(status CommunityStatus) {
	c.status = status
	c.updatedAt = NewTimestamp()
}

// Activate 激活社区
func (c *Community) Activate() {
	c.status = CommunityStatusActive
	c.updatedAt = NewTimestamp()
}

// Deactivate 停用社区
func (c *Community) Deactivate() {
	c.status = CommunityStatusInactive
	c.updatedAt = NewTimestamp()
}

// Delete 删除社区
func (c *Community) Delete() {
	c.status = CommunityStatusDeleted
	c.updatedAt = NewTimestamp()
}

// IncrementGroupCount 增加小组数量
func (c *Community) IncrementGroupCount() {
	c.groupCount++
	c.updatedAt = NewTimestamp()
}

// DecrementGroupCount 减少小组数量
func (c *Community) DecrementGroupCount() {
	if c.groupCount > 0 {
		c.groupCount--
		c.updatedAt = NewTimestamp()
	}
}

// IncrementMemberCount 增加成员数量
func (c *Community) IncrementMemberCount() {
	c.memberCount++
	c.updatedAt = NewTimestamp()
}

// DecrementMemberCount 减少成员数量
func (c *Community) DecrementMemberCount() {
	if c.memberCount > 0 {
		c.memberCount--
		c.updatedAt = NewTimestamp()
	}
}

// IncrementPostCount 增加帖子数量
func (c *Community) IncrementPostCount() {
	c.postCount++
	c.updatedAt = NewTimestamp()
}

// DecrementPostCount 减少帖子数量
func (c *Community) DecrementPostCount() {
	if c.postCount > 0 {
		c.postCount--
		c.updatedAt = NewTimestamp()
	}
}

// CanAddGroup 检查是否可以添加小组
func (c *Community) CanAddGroup() bool {
	return c.status == CommunityStatusActive
}

// IsOwner 检查用户是否为社区拥有者
func (c *Community) IsOwner(userID ID) bool {
	return c.ownerID.Value() == userID.Value()
}

// IsActive 检查社区是否激活
func (c *Community) IsActive() bool {
	return c.status == CommunityStatusActive
}

// IsDeleted 检查社区是否已删除
func (c *Community) IsDeleted() bool {
	return c.status == CommunityStatusDeleted
}
