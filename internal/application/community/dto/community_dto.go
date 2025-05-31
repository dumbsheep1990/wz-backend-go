package dto

import (
	"time"

	"github.com/yourusername/wz-backend-go/internal/domain/community/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/community/valueobject"
)

// CreateCommunityRequest 表示创建社区的请求
type CreateCommunityRequest struct {
	Name        string   `json:"name" validate:"required,max=100"`
	Description string   `json:"description"`
	OwnerID     string   `json:"owner_id" validate:"required"`
	OwnerName   string   `json:"owner_name"`
	Tags        []string `json:"tags"`
	Location    string   `json:"location"`
}

// UpdateCommunityRequest 表示更新社区的请求
type UpdateCommunityRequest struct {
	ID          string   `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"omitempty,max=100"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Location    string   `json:"location"`
	Status      string   `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE DELETED"`
}

// CommunityDTO 表示社区数据传输对象
type CommunityDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	OwnerName   string    `json:"owner_name,omitempty"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags,omitempty"`
	Location    string    `json:"location,omitempty"`
	GroupCount  int       `json:"group_count"`
	MemberCount int       `json:"member_count"`
	PostCount   int       `json:"post_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCommunityDTOFromEntity 从社区实体创建一个新的CommunityDTO
func NewCommunityDTOFromEntity(community *entity.Community, ownerName string, groupCount, memberCount, postCount int) *CommunityDTO {
	tags := make([]string, 0)
	for _, tag := range community.Tags() {
		tags = append(tags, tag.String())
	}

	return &CommunityDTO{
		ID:          community.ID().String(),
		Name:        community.Name().String(),
		Description: community.Description().String(),
		OwnerID:     community.OwnerID().String(),
		OwnerName:   ownerName,
		Status:      string(community.Status()),
		Tags:        tags,
		Location:    community.Location().String(),
		GroupCount:  groupCount,
		MemberCount: memberCount,
		PostCount:   postCount,
		CreatedAt:   community.CreatedAt(),
		UpdatedAt:   community.UpdatedAt(),
	}
}

// ListCommunitiesRequest 表示列出社区的请求
type ListCommunitiesRequest struct {
	OwnerID  string `json:"owner_id"`
	Tag      string `json:"tag"`
	Location string `json:"location"`
	Status   string `json:"status"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

// GetCommunityRequest 表示获取社区的请求
type GetCommunityRequest struct {
	ID string `json:"id" validate:"required"`
}

// DeleteCommunityRequest 表示删除社区的请求
type DeleteCommunityRequest struct {
	ID string `json:"id" validate:"required"`
}

// CommunityResponse 表示包含社区的响应
type CommunityResponse struct {
	Community *CommunityDTO `json:"community"`
}

// CommunitiesResponse 表示包含多个社区的响应
type CommunitiesResponse struct {
	Communities []*CommunityDTO `json:"communities"`
	Total       int             `json:"total"`
}

// CreateCommunityCommand 创建社区命令
type CreateCommunityCommand struct {
	Name        string   `json:"name" validate:"required,min=2,max=50"`
	Description string   `json:"description" validate:"max=1000"`
	OwnerID     string   `json:"owner_id" validate:"required"`
	OwnerName   string   `json:"owner_name" validate:"required"`
	Tags        []string `json:"tags" validate:"max=10,dive,min=1,max=20"`
	Location    string   `json:"location" validate:"max=100"`
}

// UpdateCommunityCommand 更新社区命令
type UpdateCommunityCommand struct {
	ID          string   `json:"id" validate:"required"`
	Name        string   `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Description string   `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        []string `json:"tags,omitempty" validate:"omitempty,max=10,dive,min=1,max=20"`
	Location    string   `json:"location,omitempty" validate:"omitempty,max=100"`
	OperatorID  string   `json:"operator_id" validate:"required"`
}

// DeleteCommunityCommand 删除社区命令
type DeleteCommunityCommand struct {
	ID         string `json:"id" validate:"required"`
	OperatorID string `json:"operator_id" validate:"required"`
	Reason     string `json:"reason" validate:"required,min=1,max=200"`
}

// ChangeCommunityStatusCommand 变更社区状态命令
type ChangeCommunityStatusCommand struct {
	ID         string `json:"id" validate:"required"`
	Status     int32  `json:"status" validate:"required,min=0,max=5"`
	OperatorID string `json:"operator_id" validate:"required"`
	Reason     string `json:"reason,omitempty" validate:"omitempty,max=200"`
}

// JoinCommunityCommand 加入社区命令
type JoinCommunityCommand struct {
	CommunityID string `json:"community_id" validate:"required"`
	MemberID    string `json:"member_id" validate:"required"`
	MemberName  string `json:"member_name" validate:"required"`
	JoinMethod  string `json:"join_method" validate:"required,oneof=invite request public admin"`
}

// LeaveCommunityCommand 离开社区命令
type LeaveCommunityCommand struct {
	CommunityID  string `json:"community_id" validate:"required"`
	MemberID     string `json:"member_id" validate:"required"`
	MemberName   string `json:"member_name" validate:"required"`
	LeaveReason  string `json:"leave_reason,omitempty" validate:"omitempty,max=200"`
}

// GetCommunityQuery 获取社区查询
type GetCommunityQuery struct {
	ID string `json:"id" validate:"required"`
}

// ListCommunitiesQuery 获取社区列表查询
type ListCommunitiesQuery struct {
	Status   string   `json:"status,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Location string   `json:"location,omitempty"`
	Search   string   `json:"search,omitempty" validate:"omitempty,max=100"`
	Page     int      `json:"page" validate:"min=1"`
	Limit    int      `json:"limit" validate:"min=1,max=100"`
}

// GetCommunityHealthQuery 获取社区健康度查询
type GetCommunityHealthQuery struct {
	CommunityID string `json:"community_id" validate:"required"`
}

// RecommendTagsQuery 推荐标签查询
type RecommendTagsQuery struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Location    string `json:"location,omitempty" validate:"omitempty,max=100"`
}

// CommunityHealthResponse 社区健康度响应
type CommunityHealthResponse struct {
	CommunityID   string                 `json:"community_id"`
	OverallScore  int                    `json:"overall_score"`
	Indicators    map[string]interface{} `json:"indicators"`
	Suggestions   []string               `json:"suggestions"`
	LastCheckedAt time.Time              `json:"last_checked_at"`
}

// RecommendTagsResponse 推荐标签响应
type RecommendTagsResponse struct {
	RecommendedTags []string `json:"recommended_tags"`
}

// CommunityStatsResponse 社区统计响应
type CommunityStatsResponse struct {
	CommunityID       string `json:"community_id"`
	TotalGroups       int32  `json:"total_groups"`
	TotalMembers      int32  `json:"total_members"`
	TotalPosts        int32  `json:"total_posts"`
	ActiveMembers     int32  `json:"active_members"`
	PostsThisWeek     int32  `json:"posts_this_week"`
	PostsThisMonth    int32  `json:"posts_this_month"`
	GrowthRate        string `json:"growth_rate"`
	EngagementRate    string `json:"engagement_rate"`
}

// CommunityEventResponse 社区事件响应
type CommunityEventResponse struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	CommunityID string                 `json:"community_id"`
	AggregateID string                 `json:"aggregate_id"`
	Data        map[string]interface{} `json:"data"`
	OccurredAt  time.Time              `json:"occurred_at"`
	Version     int                    `json:"version"`
}

// CommunitySearchResponse 社区搜索响应
type CommunitySearchResponse struct {
	Communities  []*CommunityResponse `json:"communities"`
	Total        int                  `json:"total"`
	SearchTerms  []string             `json:"search_terms"`
	Suggestions  []string             `json:"suggestions"`
	Facets       map[string][]string  `json:"facets"`
}

// CommunityMemberResponse 社区成员响应
type CommunityMemberResponse struct {
	MemberID     string    `json:"member_id"`
	MemberName   string    `json:"member_name"`
	Role         string    `json:"role"`
	JoinedAt     time.Time `json:"joined_at"`
	LastActiveAt time.Time `json:"last_active_at"`
	PostCount    int32     `json:"post_count"`
	IsActive     bool      `json:"is_active"`
}

// CommunityMembersResponse 社区成员列表响应
type CommunityMembersResponse struct {
	CommunityID string                     `json:"community_id"`
	Members     []*CommunityMemberResponse `json:"members"`
	Total       int                        `json:"total"`
	Page        int                        `json:"page"`
	Limit       int                        `json:"limit"`
}
