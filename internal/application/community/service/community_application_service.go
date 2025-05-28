package service

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/yourusername/wz-backend-go/internal/application/community/dto"
	"github.com/yourusername/wz-backend-go/internal/domain/community/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/community/repository"
	"github.com/yourusername/wz-backend-go/internal/domain/community/valueobject"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/persistence/database"
)

// CommunityApplicationService 处理社区相关的应用层逻辑
type CommunityApplicationService struct {
	communityRepo repository.CommunityRepository
	userRepo      UserRepository // Interface to interact with users
	eventBus      event.EventBus
	validator     *validator.Validate
	unitOfWork    database.UnitOfWork
}

// UserRepository 定义与用户交互的接口
type UserRepository interface {
	FindUserNameByID(ctx context.Context, id string) (string, error)
}

// NewCommunityApplicationService 创建一个新的社区应用服务
func NewCommunityApplicationService(
	communityRepo repository.CommunityRepository,
	userRepo UserRepository,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) *CommunityApplicationService {
	return &CommunityApplicationService{
		communityRepo: communityRepo,
		userRepo:      userRepo,
		eventBus:      eventBus,
		validator:     validator.New(),
		unitOfWork:    unitOfWork,
	}
}

// CreateCommunity 创建一个新社区
func (s *CommunityApplicationService) CreateCommunity(ctx context.Context, req dto.CreateCommunityRequest) (*dto.CommunityDTO, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// 将请求值转换为领域值对象
	name, err := valueobject.NewCommunityName(req.Name)
	if err != nil {
		return nil, err
	}

	description := valueobject.NewDescription(req.Description)
	ownerID := valueobject.NewUserID(req.OwnerID)
	location := valueobject.NewLocation(req.Location)

	// 转换标签
	tags := make([]valueobject.Tag, 0, len(req.Tags))
	for _, tag := range req.Tags {
		tagVO, err := valueobject.NewTag(tag)
		if err != nil {
			continue // Skip invalid tags
		}
		tags = append(tags, tagVO)
	}

	// 创建领域实体
	community, err := entity.NewCommunity(name, description, ownerID, tags, location)
	if err != nil {
		return nil, err
	}

	// Use unit of work to ensure atomicity
	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		// 持久化保存社区
		if err := s.communityRepo.Save(ctx, community); err != nil {
			return err
		}

		// Publish domain events
		events := community.GetDomainEvents()
		for _, event := range events {
			s.eventBus.Publish(ctx, event)
		}
		community.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get owner name for the response DTO
	ownerName, _ := s.userRepo.FindUserNameByID(ctx, req.OwnerID)

	// 返回DTO
	return dto.NewCommunityDTOFromEntity(community, ownerName, 0, 1, 0), nil
}

// GetCommunity 通过ID获取社区
func (s *CommunityApplicationService) GetCommunity(ctx context.Context, req dto.GetCommunityRequest) (*dto.CommunityDTO, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	id := valueobject.NewID(req.ID)
	community, err := s.communityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if community == nil {
		return nil, errors.New("community not found")
	}

	// Get owner name
	ownerName, _ := s.userRepo.FindUserNameByID(ctx, community.OwnerID().String())

	// TODO: Add logic to get group count, member count, and post count
	// 目前我们使用占位符
	groupCount := 0
	memberCount := 1
	postCount := 0

	return dto.NewCommunityDTOFromEntity(community, ownerName, groupCount, memberCount, postCount), nil
}

// UpdateCommunity 更新现有社区
func (s *CommunityApplicationService) UpdateCommunity(ctx context.Context, req dto.UpdateCommunityRequest) (*dto.CommunityDTO, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	id := valueobject.NewID(req.ID)
	community, err := s.communityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if community == nil {
		return nil, errors.New("community not found")
	}

	// 根据请求应用更新
	if req.Name != "" {
		name, err := valueobject.NewCommunityName(req.Name)
		if err != nil {
			return nil, err
		}
		if err := community.UpdateName(name); err != nil {
			return nil, err
		}
	}

	if req.Description != "" {
		description := valueobject.NewDescription(req.Description)
		if err := community.UpdateDescription(description); err != nil {
			return nil, err
		}
	}

	if len(req.Tags) > 0 {
		tags := make([]valueobject.Tag, 0, len(req.Tags))
		for _, tag := range req.Tags {
			tagVO, err := valueobject.NewTag(tag)
			if err != nil {
				continue // Skip invalid tags
			}
			tags = append(tags, tagVO)
		}
		if err := community.UpdateTags(tags); err != nil {
			return nil, err
		}
	}

	if req.Location != "" {
		location := valueobject.NewLocation(req.Location)
		if err := community.UpdateLocation(location); err != nil {
			return nil, err
		}
	}

	// 如果指定了状态，处理状态变更
	if req.Status != "" {
		switch valueobject.CommunityStatus(req.Status) {
		case valueobject.StatusActive:
			if err := community.Activate(); err != nil && err.Error() != "社区已经是活跃状态" {
				return nil, err
			}
		case valueobject.StatusDeleted:
			if err := community.Delete(); err != nil && err.Error() != "社区已经被删除" {
				return nil, err
			}
		}
	}

	// Use unit of work to ensure atomicity
	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		// 持久化保存更新后的社区
		if err := s.communityRepo.Save(ctx, community); err != nil {
			return err
		}

		// Publish domain events
		events := community.GetDomainEvents()
		for _, event := range events {
			s.eventBus.Publish(ctx, event)
		}
		community.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get owner name for the response DTO
	ownerName, _ := s.userRepo.FindUserNameByID(ctx, community.OwnerID().String())

	// TODO: Add logic to get group count, member count, and post count
	groupCount := 0
	memberCount := 1
	postCount := 0

	return dto.NewCommunityDTOFromEntity(community, ownerName, groupCount, memberCount, postCount), nil
}

// DeleteCommunity 删除社区
func (s *CommunityApplicationService) DeleteCommunity(ctx context.Context, req dto.DeleteCommunityRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return err
	}

	id := valueobject.NewID(req.ID)
	community, err := s.communityRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if community == nil {
		return errors.New("未找到社区")
	}

	// 在领域模型中将社区标记为已删除
	if err := community.Delete(); err != nil {
		return err
	}

	// Use unit of work to ensure atomicity
	return s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		// 可以更新社区状态或物理删除它
		// 这里我们只是更新状态
		if err := s.communityRepo.Save(ctx, community); err != nil {
			return err
		}

		// Publish domain events
		events := community.GetDomainEvents()
		for _, event := range events {
			s.eventBus.Publish(ctx, event)
		}
		community.ClearDomainEvents()

		return nil
	})
}

// ListCommunities 使用过滤和分页列出社区
func (s *CommunityApplicationService) ListCommunities(ctx context.Context, req dto.ListCommunitiesRequest) (*dto.CommunitiesResponse, error) {
	// 默认分页值
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // 默认限制
	} else if limit > 100 {
		limit = 100 // 最大限制
	}

	var communities []*entity.Community
	var total int
	var err error

	// 根据请求应用过滤器
	if req.OwnerID != "" {
		ownerID := valueobject.NewUserID(req.OwnerID)
		communities, err = s.communityRepo.FindByOwnerID(ctx, ownerID)
		if err != nil {
			return nil, err
		}
		total = len(communities)
		// Simple pagination in memory for this example
		if offset < total {
			end := offset + limit
			if end > total {
				end = total
			}
			communities = communities[offset:end]
		} else {
			communities = []*entity.Community{}
		}
	} else if req.Tag != "" {
		tag, err := valueobject.NewTag(req.Tag)
		if err != nil {
			return nil, err
		}
		communities, err = s.communityRepo.FindByTags(ctx, []valueobject.Tag{tag})
		if err != nil {
			return nil, err
		}
		total = len(communities)
		// Simple pagination in memory for this example
		if offset < total {
			end := offset + limit
			if end > total {
				end = total
			}
			communities = communities[offset:end]
		} else {
			communities = []*entity.Community{}
		}
	} else if req.Location != "" {
		location := valueobject.NewLocation(req.Location)
		communities, err = s.communityRepo.FindByLocation(ctx, location)
		if err != nil {
			return nil, err
		}
		total = len(communities)
		// Simple pagination in memory for this example
		if offset < total {
			end := offset + limit
			if end > total {
				end = total
			}
			communities = communities[offset:end]
		} else {
			communities = []*entity.Community{}
		}
	} else {
		// 获取所有社区并进行分页
		communities, total, err = s.communityRepo.FindAll(ctx, offset, limit)
		if err != nil {
			return nil, err
		}
	}

	// 将实体转换为DTO对象
	dtos := make([]*dto.CommunityDTO, 0, len(communities))
	for _, community := range communities {
		ownerName, _ := s.userRepo.FindUserNameByID(ctx, community.OwnerID().String())
		
		// TODO: 添加获取群组数量、成员数量和帖子数量的逻辑
		groupCount := 0
		memberCount := 1
		postCount := 0
		
		dtos = append(dtos, dto.NewCommunityDTOFromEntity(community, ownerName, groupCount, memberCount, postCount))
	}

	return &dto.CommunitiesResponse{
		Communities: dtos,
		Total:       total,
	}, nil
}
