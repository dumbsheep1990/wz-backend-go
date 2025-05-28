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
