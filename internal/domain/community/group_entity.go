package community

import (
	"errors"
	"time"
)

// Group 小组聚合根
type Group struct {
	id          ID
	name        CommunityName
	description Description
	communityID ID
	ownerID     ID
	ownerName   string
	status      GroupStatus
	tags        Tags
	createdAt   Timestamp
	updatedAt   Timestamp
	
	// 成员列表（仅存储ID，成员详情由用户服务管理）
	members     []ID
	
	// 聚合计数器
	memberCount int
	postCount   int
}

// NewGroup 创建小组
func NewGroup(
	name CommunityName,
	description Description,
	communityID ID,
	ownerID ID,
	ownerName string,
	tags Tags,
) (*Group, error) {
	now := NewTimestamp()
	
	return &Group{
		id:          NewID(),
		name:        name,
		description: description,
		communityID: communityID,
		ownerID:     ownerID,
		ownerName:   ownerName,
		status:      GroupStatusActive,
		tags:        tags,
		createdAt:   now,
		updatedAt:   now,
		members:     []ID{ownerID}, // 创建者自动成为成员
		memberCount: 1,
		postCount:   0,
	}, nil
}

// ReconstructGroup 从数据库重构小组实体
func ReconstructGroup(
	id string,
	name string,
	description string,
	communityID string,
	ownerID string,
	ownerName string,
	status int,
	members []string,
	tags []string,
	createdAt time.Time,
	updatedAt time.Time,
	memberCount int,
	postCount int,
) (*Group, error) {
	groupID, err := NewIDFromString(id)
	if err != nil {
		return nil, err
	}
	
	groupName, err := NewCommunityName(name)
	if err != nil {
		return nil, err
	}
	
	groupDesc, err := NewDescription(description)
	if err != nil {
		return nil, err
	}
	
	commID, err := NewIDFromString(communityID)
	if err != nil {
		return nil, err
	}
	
	ownID, err := NewIDFromString(ownerID)
	if err != nil {
		return nil, err
	}
	
	groupTags, err := NewTags(tags)
	if err != nil {
		return nil, err
	}
	
	// 转换成员ID列表
	memberIDs := make([]ID, 0, len(members))
	for _, m := range members {
		memberID, err := NewIDFromString(m)
		if err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, memberID)
	}
	
	return &Group{
		id:          groupID,
		name:        groupName,
		description: groupDesc,
		communityID: commID,
		ownerID:     ownID,
		ownerName:   ownerName,
		status:      GroupStatus(status),
		tags:        groupTags,
		createdAt:   NewTimestampFromTime(createdAt),
		updatedAt:   NewTimestampFromTime(updatedAt),
		members:     memberIDs,
		memberCount: memberCount,
		postCount:   postCount,
	}, nil
}

// ID 获取小组ID
func (g *Group) ID() ID {
	return g.id
}

// Name 获取小组名称
func (g *Group) Name() CommunityName {
	return g.name
}

// Description 获取小组描述
func (g *Group) Description() Description {
	return g.description
}

// CommunityID 获取所属社区ID
func (g *Group) CommunityID() ID {
	return g.communityID
}

// OwnerID 获取拥有者ID
func (g *Group) OwnerID() ID {
	return g.ownerID
}

// OwnerName 获取拥有者名称
func (g *Group) OwnerName() string {
	return g.ownerName
}

// Status 获取小组状态
func (g *Group) Status() GroupStatus {
	return g.status
}

// Tags 获取小组标签
func (g *Group) Tags() Tags {
	return g.tags
}

// CreatedAt 获取创建时间
func (g *Group) CreatedAt() Timestamp {
	return g.createdAt
}

// UpdatedAt 获取更新时间
func (g *Group) UpdatedAt() Timestamp {
	return g.updatedAt
}

// Members 获取成员ID列表
func (g *Group) Members() []ID {
	return g.members
}

// MemberCount 获取成员数量
func (g *Group) MemberCount() int {
	return g.memberCount
}

// PostCount 获取帖子数量
func (g *Group) PostCount() int {
	return g.postCount
}

// UpdateName 更新小组名称
func (g *Group) UpdateName(name CommunityName) {
	g.name = name
	g.updatedAt = NewTimestamp()
}

// UpdateDescription 更新小组描述
func (g *Group) UpdateDescription(description Description) {
	g.description = description
	g.updatedAt = NewTimestamp()
}

// UpdateTags 更新小组标签
func (g *Group) UpdateTags(tags Tags) {
	g.tags = tags
	g.updatedAt = NewTimestamp()
}

// UpdateStatus 更新小组状态
func (g *Group) UpdateStatus(status GroupStatus) {
	g.status = status
	g.updatedAt = NewTimestamp()
}

// Activate 激活小组
func (g *Group) Activate() {
	g.status = GroupStatusActive
	g.updatedAt = NewTimestamp()
}

// Deactivate 停用小组
func (g *Group) Deactivate() {
	g.status = GroupStatusInactive
	g.updatedAt = NewTimestamp()
}

// Delete 删除小组
func (g *Group) Delete() {
	g.status = GroupStatusDeleted
	g.updatedAt = NewTimestamp()
}

// AddMember 添加成员
func (g *Group) AddMember(memberID ID) error {
	// 检查小组状态
	if g.status != GroupStatusActive {
		return errors.New("只有活跃的小组才能添加成员")
	}
	
	// 检查成员是否已存在
	for _, id := range g.members {
		if id.Value() == memberID.Value() {
			return errors.New("成员已在小组中")
		}
	}
	
	// 添加成员
	g.members = append(g.members, memberID)
	g.memberCount++
	g.updatedAt = NewTimestamp()
	
	return nil
}

// RemoveMember 移除成员
func (g *Group) RemoveMember(memberID ID) error {
	// 检查是否为拥有者
	if g.ownerID.Value() == memberID.Value() {
		return errors.New("无法移除小组拥有者")
	}
	
	// 查找并移除成员
	for i, id := range g.members {
		if id.Value() == memberID.Value() {
			// 移除成员
			g.members = append(g.members[:i], g.members[i+1:]...)
			g.memberCount--
			g.updatedAt = NewTimestamp()
			return nil
		}
	}
	
	return errors.New("成员不在小组中")
}

// IncrementPostCount 增加帖子数量
func (g *Group) IncrementPostCount() {
	g.postCount++
	g.updatedAt = NewTimestamp()
}

// DecrementPostCount 减少帖子数量
func (g *Group) DecrementPostCount() {
	if g.postCount > 0 {
		g.postCount--
		g.updatedAt = NewTimestamp()
	}
}

// IsMember 检查用户是否为小组成员
func (g *Group) IsMember(userID ID) bool {
	for _, id := range g.members {
		if id.Value() == userID.Value() {
			return true
		}
	}
	return false
}

// IsOwner 检查用户是否为小组拥有者
func (g *Group) IsOwner(userID ID) bool {
	return g.ownerID.Value() == userID.Value()
}

// IsActive 检查小组是否激活
func (g *Group) IsActive() bool {
	return g.status == GroupStatusActive
}

// IsDeleted 检查小组是否已删除
func (g *Group) IsDeleted() bool {
	return g.status == GroupStatusDeleted
}

// CanPost 检查用户是否可以在小组中发帖
func (g *Group) CanPost(userID ID) bool {
	return g.status == GroupStatusActive && g.IsMember(userID)
}
