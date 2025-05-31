package entity

import (
	"time"

	"wz-backend-go/internal/domain/community/valueobject"
	"wz-backend-go/internal/domain/shared/event"
)

const (
	EventTypeCommunityCreated   = "community.created"
	EventTypeCommunityUpdated   = "community.updated"
	EventTypeCommunityDeleted   = "community.deleted"
	EventTypeCommunityActivated = "community.activated"
)

// CommunityCreatedEvent 表示社区创建时的事件
type CommunityCreatedEvent struct {
	event.BaseEvent
	CommunityID   string            `json:"community_id"`
	CommunityName string            `json:"community_name"`
	OwnerID       string            `json:"owner_id"`
	OwnerName     string            `json:"owner_name"`
	Description   string            `json:"description"`
	Tags          []string          `json:"tags"`
	Location      string            `json:"location"`
	Metadata      map[string]string `json:"metadata"`
}

// NewCommunityCreatedEvent 创建一个新的社区创建事件
func NewCommunityCreatedEvent(community *Community) *CommunityCreatedEvent {
	return &CommunityCreatedEvent{
		BaseEvent: event.BaseEvent{
			ID:          event.NewEventID(),
			AggregateID: community.id.Value(),
			EventType:   "CommunityCreated",
			Version:     1,
			OccurredOn:  time.Now(),
		},
		CommunityID:   community.id.Value(),
		CommunityName: community.name.Value(),
		OwnerID:       community.ownerID,
		OwnerName:     community.ownerName,
		Description:   community.description,
		Tags:          community.tags.Values(),
		Location:      community.location,
		Metadata: map[string]string{
			"initial_member_count": "1",
			"initial_group_count":  "0",
			"initial_post_count":   "0",
		},
	}
}

// CommunityUpdatedEvent 表示社区更新时的事件
type CommunityUpdatedEvent struct {
	event.BaseEvent
	CommunityID string `json:"community_id"`
	FieldName   string `json:"field_name"`
	OldValue    string `json:"old_value"`
	NewValue    string `json:"new_value"`
	UpdatedBy   string `json:"updated_by,omitempty"`
}

// NewCommunityUpdatedEvent 创建一个新的社区更新事件
func NewCommunityUpdatedEvent(community *Community, fieldName, oldValue, newValue string) *CommunityUpdatedEvent {
	return &CommunityUpdatedEvent{
		BaseEvent: event.BaseEvent{
			ID:          event.NewEventID(),
			AggregateID: community.id.Value(),
			EventType:   "CommunityUpdated",
			Version:     1,
			OccurredOn:  time.Now(),
		},
		CommunityID: community.id.Value(),
		FieldName:   fieldName,
		OldValue:    oldValue,
		NewValue:    newValue,
		UpdatedBy:   community.ownerID, // 可以后续改为实际操作人
	}
}

// CommunityDeletedEvent 表示社区删除时的事件
type CommunityDeletedEvent struct {
	event.BaseEvent
	CommunityID   string            `json:"community_id"`
	CommunityName string            `json:"community_name"`
	OwnerID       string            `json:"owner_id"`
	OwnerName     string            `json:"owner_name"`
	Reason        string            `json:"reason"`
	DeletedBy     string            `json:"deleted_by,omitempty"`
	Stats         map[string]int32  `json:"stats"`
	Metadata      map[string]string `json:"metadata"`
}

// NewCommunityDeletedEvent 创建一个新的社区删除事件
func NewCommunityDeletedEvent(community *Community, reason string) *CommunityDeletedEvent {
	return &CommunityDeletedEvent{
		BaseEvent: event.BaseEvent{
			ID:          event.NewEventID(),
			AggregateID: community.id.Value(),
			EventType:   "CommunityDeleted",
			Version:     1,
			OccurredOn:  time.Now(),
		},
		CommunityID:   community.id.Value(),
		CommunityName: community.name.Value(),
		OwnerID:       community.ownerID,
		OwnerName:     community.ownerName,
		Reason:        reason,
		DeletedBy:     community.ownerID,
		Stats: map[string]int32{
			"final_member_count": community.memberCount,
			"final_group_count":  community.groupCount,
			"final_post_count":   community.postCount,
		},
		Metadata: map[string]string{
			"tags":     community.tags.ToString(),
			"location": community.location,
		},
	}
}

// CommunityActivatedEvent 表示社区激活时的事件
type CommunityActivatedEvent struct {
	event.BaseDomainEvent
}

// NewCommunityActivatedEvent 创建一个新的社区激活事件
func NewCommunityActivatedEvent(community *Community) CommunityActivatedEvent {
	return CommunityActivatedEvent{
		BaseDomainEvent: event.NewBaseDomainEvent(
			EventTypeCommunityActivated,
			community.ID().String(),
			"Community",
			map[string]interface{}{
				"id":          community.ID().String(),
				"activatedAt": community.UpdatedAt(),
			},
		),
	}
}

// CommunityStatusChangedEvent 社区状态变更事件
type CommunityStatusChangedEvent struct {
	event.BaseEvent
	CommunityID   string `json:"community_id"`
	CommunityName string `json:"community_name"`
	OldStatus     string `json:"old_status"`
	NewStatus     string `json:"new_status"`
	Reason        string `json:"reason"`
	ChangedBy     string `json:"changed_by,omitempty"`
}

// NewCommunityStatusChangedEvent 创建社区状态变更事件
func NewCommunityStatusChangedEvent(community *Community, oldStatus, newStatus valueobject.CommunityStatus, reason string) *CommunityStatusChangedEvent {
	return &CommunityStatusChangedEvent{
		BaseEvent: event.BaseEvent{
			ID:          event.NewEventID(),
			AggregateID: community.id.Value(),
			EventType:   "CommunityStatusChanged",
			Version:     1,
			OccurredOn:  time.Now(),
		},
		CommunityID:   community.id.Value(),
		CommunityName: community.name.Value(),
		OldStatus:     oldStatus.String(),
		NewStatus:     newStatus.String(),
		Reason:        reason,
		ChangedBy:     community.ownerID,
	}
}

// CommunityStatsUpdatedEvent 社区统计更新事件
type CommunityStatsUpdatedEvent struct {
	event.BaseEvent
	CommunityID string `json:"community_id"`
	StatType    string `json:"stat_type"`
	NewValue    int32  `json:"new_value"`
	Change      int32  `json:"change"`
}

// NewCommunityStatsUpdatedEvent 创建社区统计更新事件
func NewCommunityStatsUpdatedEvent(community *Community, statType string, newValue int32) *CommunityStatsUpdatedEvent {
	var oldValue int32
	switch statType {
	case "member_count":
		oldValue = newValue - 1
	case "group_count":
		oldValue = newValue - 1
	case "post_count":
		oldValue = newValue - 1
	}

	return &CommunityStatsUpdatedEvent{
		BaseEvent: event.BaseEvent{
			ID:          event.NewEventID(),
			AggregateID: community.id.Value(),
			EventType:   "CommunityStatsUpdated",
			Version:     1,
			OccurredOn:  time.Now(),
		},
		CommunityID: community.id.Value(),
		StatType:    statType,
		NewValue:    newValue,
		Change:      newValue - oldValue,
	}
}

// CommunityMemberJoinedEvent 社区成员加入事件
type CommunityMemberJoinedEvent struct {
	event.BaseEvent
	CommunityID   string            `json:"community_id"`
	CommunityName string            `json:"community_name"`
	MemberID      string            `json:"member_id"`
	MemberName    string            `json:"member_name"`
	JoinMethod    string            `json:"join_method"`
	Metadata      map[string]string `json:"metadata"`
}

// NewCommunityMemberJoinedEvent 创建社区成员加入事件
func NewCommunityMemberJoinedEvent(community *Community, memberID, memberName, joinMethod string) *CommunityMemberJoinedEvent {
	return &CommunityMemberJoinedEvent{
		BaseEvent: event.BaseEvent{
			ID:          event.NewEventID(),
			AggregateID: community.id.Value(),
			EventType:   "CommunityMemberJoined",
			Version:     1,
			OccurredOn:  time.Now(),
		},
		CommunityID:   community.id.Value(),
		CommunityName: community.name.Value(),
		MemberID:      memberID,
		MemberName:    memberName,
		JoinMethod:    joinMethod,
		Metadata: map[string]string{
			"total_members": string(rune(community.memberCount)),
		},
	}
}

// CommunityMemberLeftEvent 社区成员离开事件
type CommunityMemberLeftEvent struct {
	event.BaseEvent
	CommunityID   string            `json:"community_id"`
	CommunityName string            `json:"community_name"`
	MemberID      string            `json:"member_id"`
	MemberName    string            `json:"member_name"`
	LeaveReason   string            `json:"leave_reason"`
	Metadata      map[string]string `json:"metadata"`
}

// NewCommunityMemberLeftEvent 创建社区成员离开事件
func NewCommunityMemberLeftEvent(community *Community, memberID, memberName, leaveReason string) *CommunityMemberLeftEvent {
	return &CommunityMemberLeftEvent{
		BaseEvent: event.BaseEvent{
			ID:          event.NewEventID(),
			AggregateID: community.id.Value(),
			EventType:   "CommunityMemberLeft",
			Version:     1,
			OccurredOn:  time.Now(),
		},
		CommunityID:   community.id.Value(),
		CommunityName: community.name.Value(),
		MemberID:      memberID,
		MemberName:    memberName,
		LeaveReason:   leaveReason,
		Metadata: map[string]string{
			"remaining_members": string(rune(community.memberCount)),
		},
	}
}
