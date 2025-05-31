package service

import (
	"context"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"

	"github.com/yourusername/wz-backend-go/internal/application/community/dto"
	"github.com/yourusername/wz-backend-go/internal/domain/community/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/community/repository"
	"github.com/yourusername/wz-backend-go/internal/domain/community/service"
	"github.com/yourusername/wz-backend-go/internal/domain/community/valueobject"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/persistence/database"
)

// CommunityApplicationService 处理社区相关的应用层逻辑
type CommunityApplicationService struct {
	communityRepo      repository.CommunityRepository
	userRepo           UserRepository // Interface to interact with users
	eventDispatcher    event.EventDispatcher
	validator          *validator.Validate
	unitOfWork         database.UnitOfWork
	domainService      *service.CommunityDomainService
}

// UserRepository 定义与用户交互的接口
type UserRepository interface {
	FindUserNameByID(ctx context.Context, id string) (string, error)
}

// NewCommunityApplicationService 创建一个新的社区应用服务
func NewCommunityApplicationService(
	communityRepo repository.CommunityRepository,
	userRepo UserRepository,
	eventDispatcher event.EventDispatcher,
	unitOfWork database.UnitOfWork,
	domainService *service.CommunityDomainService,
) *CommunityApplicationService {
	return &CommunityApplicationService{
		communityRepo:   communityRepo,
		userRepo:        userRepo,
		eventDispatcher: eventDispatcher,
		validator:       validator.New(),
		unitOfWork:      unitOfWork,
		domainService:   domainService,
	}
}

// CreateCommunity 创建一个新社区
func (s *CommunityApplicationService) CreateCommunity(ctx context.Context, cmd *dto.CreateCommunityCommand) (*dto.CommunityResponse, error) {
	if err := s.validator.Struct(cmd); err != nil {
		return nil, err
	}

	// 使用领域服务创建社区
	community, err := s.domainService.CreateCommunity(
		cmd.Name,
		cmd.Description,
		cmd.OwnerID,
		cmd.OwnerName,
		cmd.Tags,
		cmd.Location,
	)
	if err != nil {
		return nil, err
	}

	// 检查名称唯一性
	exists, err := s.communityRepo.ExistsByName(ctx, community.Name())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("社区名称已存在")
	}

	// Use unit of work to ensure atomicity
	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		// 持久化保存社区
		if err := s.communityRepo.Save(ctx, community); err != nil {
			return err
		}

		// Publish domain events
		s.publishDomainEvents(ctx, community)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get owner name for the response DTO
	ownerName, _ := s.userRepo.FindUserNameByID(ctx, cmd.OwnerID)

	// 返回DTO
	return s.toDTO(community, ownerName), nil
}

// GetCommunity 通过ID获取社区
func (s *CommunityApplicationService) GetCommunity(ctx context.Context, query *dto.GetCommunityQuery) (*dto.CommunityResponse, error) {
	if err := s.validator.Struct(query); err != nil {
		return nil, err
	}

	communityID, err := valueobject.NewCommunityID(query.ID)
	if err != nil {
		return nil, err
	}

	community, err := s.communityRepo.FindByID(ctx, communityID)
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

	return s.toDTO(community, ownerName), nil
}

// UpdateCommunity 更新现有社区
func (s *CommunityApplicationService) UpdateCommunity(ctx context.Context, cmd *dto.UpdateCommunityCommand) (*dto.CommunityResponse, error) {
	if err := s.validator.Struct(cmd); err != nil {
		return nil, err
	}

	// 获取社区
	communityID, err := valueobject.NewCommunityID(cmd.ID)
	if err != nil {
		return nil, err
	}

	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return nil, err
	}

	if community == nil {
		return nil, errors.New("community not found")
	}

	// 构建更新映射
	updates := make(map[string]interface{})
	if cmd.Name != "" {
		updates["name"] = cmd.Name
	}
	if cmd.Description != "" {
		updates["description"] = cmd.Description
	}
	if len(cmd.Tags) > 0 {
		updates["tags"] = cmd.Tags
	}
	if cmd.Location != "" {
		updates["location"] = cmd.Location
	}

	// 验证更新操作
	if err := s.domainService.ValidateCommunityUpdate(community, cmd.OperatorID, updates); err != nil {
		return nil, err
	}

	// 执行更新
	for field, value := range updates {
		switch field {
		case "name":
			if name, ok := value.(string); ok {
				communityName, err := valueobject.NewCommunityName(name)
				if err != nil {
					return nil, err
				}
				if err := community.UpdateName(communityName); err != nil {
					return nil, err
				}
			}
		case "description":
			if desc, ok := value.(string); ok {
				community.UpdateDescription(desc)
			}
		case "tags":
			if tags, ok := value.([]string); ok {
				communityTags, err := valueobject.NewTags(tags)
				if err != nil {
					return nil, err
				}
				community.UpdateTags(communityTags)
			}
		case "location":
			if loc, ok := value.(string); ok {
				community.UpdateLocation(loc)
			}
		}
	}

	// 保存更新
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return nil, err
	}

	// 发布领域事件
	s.publishDomainEvents(ctx, community)

	// Get owner name for the response DTO
	ownerName, _ := s.userRepo.FindUserNameByID(ctx, community.OwnerID().String())

	// TODO: Add logic to get group count, member count, and post count
	groupCount := 0
	memberCount := 1
	postCount := 0

	return s.toDTO(community, ownerName), nil
}

// DeleteCommunity 删除社区
func (s *CommunityApplicationService) DeleteCommunity(ctx context.Context, cmd *dto.DeleteCommunityCommand) error {
	if err := s.validator.Struct(cmd); err != nil {
		return err
	}

	// 获取社区
	communityID, err := valueobject.NewCommunityID(cmd.ID)
	if err != nil {
		return err
	}

	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return err
	}

	if community == nil {
		return errors.New("未找到社区")
	}

	// 检查权限
	if !community.IsOwnedBy(cmd.OperatorID) {
		return errors.New("只有社区创建者可以删除社区")
	}

	// 删除社区
	if err := community.Delete(cmd.Reason); err != nil {
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
		s.publishDomainEvents(ctx, community)

		return nil
	})
}

// ListCommunities 使用过滤和分页列出社区
func (s *CommunityApplicationService) ListCommunities(ctx context.Context, query *dto.ListCommunitiesQuery) (*dto.CommunitiesResponse, error) {
	// 构建过滤条件
	filter := &repository.CommunityFilter{
		Status:   query.Status,
		Tags:     query.Tags,
		Location: query.Location,
		Search:   query.Search,
	}

	// 构建分页参数
	pagination := &repository.Pagination{
		Page:  query.Page,
		Limit: query.Limit,
	}

	communities, total, err := s.communityRepo.FindWithFilter(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	// 将实体转换为DTO对象
	communityDTOs := make([]*dto.CommunityResponse, 0, len(communities))
	for _, community := range communities {
		ownerName, _ := s.userRepo.FindUserNameByID(ctx, community.OwnerID().String())
		
		// TODO: 添加获取群组数量、成员数量和帖子数量的逻辑
		groupCount := 0
		memberCount := 1
		postCount := 0
		
		communityDTOs = append(communityDTOs, s.toDTO(community, ownerName))
	}

	return &dto.CommunitiesResponse{
		Communities: communityDTOs,
		Total:       total,
		Page:        query.Page,
		Limit:       query.Limit,
	}, nil
}

// ChangeCommunityStatus 变更社区状态
func (s *CommunityApplicationService) ChangeCommunityStatus(ctx context.Context, cmd *dto.ChangeCommunityStatusCommand) (*dto.CommunityResponse, error) {
	// 获取社区
	communityID, err := valueobject.NewCommunityID(cmd.ID)
	if err != nil {
		return nil, err
	}

	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return nil, err
	}
	if community == nil {
		return nil, errors.New("社区不存在")
	}

	// 解析目标状态
	targetStatus, err := valueobject.NewCommunityStatus(cmd.Status)
	if err != nil {
		return nil, err
	}

	// 验证状态变更
	if err := s.domainService.ValidateCommunityStatusChange(community, targetStatus, cmd.OperatorID, cmd.Reason); err != nil {
		return nil, err
	}

	// 执行状态变更
	switch targetStatus {
	case valueobject.CommunityStatusActive:
		if err := community.Activate(); err != nil {
			return nil, err
		}
	case valueobject.CommunityStatusSuspended:
		if err := community.Suspend(cmd.Reason); err != nil {
			return nil, err
		}
	case valueobject.CommunityStatusArchived:
		if err := community.Archive(cmd.Reason); err != nil {
			return nil, err
		}
	case valueobject.CommunityStatusDeleted:
		if err := community.Delete(cmd.Reason); err != nil {
			return nil, err
		}
	case valueobject.CommunityStatusReviewing:
		if err := community.SubmitForReview(); err != nil {
			return nil, err
		}
	}

	// 保存状态变更
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return nil, err
	}

	// 发布领域事件
	s.publishDomainEvents(ctx, community)

	// Get owner name for the response DTO
	ownerName, _ := s.userRepo.FindUserNameByID(ctx, community.OwnerID().String())

	// TODO: Add logic to get group count, member count, and post count
	groupCount := 0
	memberCount := 1
	postCount := 0

	return s.toDTO(community, ownerName), nil
}

// JoinCommunity 加入社区
func (s *CommunityApplicationService) JoinCommunity(ctx context.Context, cmd *dto.JoinCommunityCommand) error {
	// 获取社区
	communityID, err := valueobject.NewCommunityID(cmd.CommunityID)
	if err != nil {
		return err
	}

	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return err
	}
	if community == nil {
		return errors.New("社区不存在")
	}

	// 验证加入操作
	if err := s.domainService.ValidateMemberJoin(community, cmd.MemberID, cmd.JoinMethod); err != nil {
		return err
	}

	// 增加成员计数
	community.IncrementMemberCount()

	// 保存变更
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return err
	}

	// 发布成员加入事件
	joinEvent := entity.NewCommunityMemberJoinedEvent(community, cmd.MemberID, cmd.MemberName, cmd.JoinMethod)
	if err := s.eventDispatcher.Dispatch(ctx, joinEvent); err != nil {
		log.Printf("Failed to dispatch member joined event: %v", err)
	}

	// 发布领域事件
	s.publishDomainEvents(ctx, community)

	return nil
}

// LeaveCommunity 离开社区
func (s *CommunityApplicationService) LeaveCommunity(ctx context.Context, cmd *dto.LeaveCommunityCommand) error {
	// 获取社区
	communityID, err := valueobject.NewCommunityID(cmd.CommunityID)
	if err != nil {
		return err
	}

	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return err
	}
	if community == nil {
		return errors.New("社区不存在")
	}

	// 验证离开操作
	if err := s.domainService.ValidateMemberLeave(community, cmd.MemberID, cmd.LeaveReason); err != nil {
		return err
	}

	// 减少成员计数
	community.DecrementMemberCount()

	// 保存变更
	if err := s.communityRepo.Save(ctx, community); err != nil {
		return err
	}

	// 发布成员离开事件
	leaveEvent := entity.NewCommunityMemberLeftEvent(community, cmd.MemberID, cmd.MemberName, cmd.LeaveReason)
	if err := s.eventDispatcher.Dispatch(ctx, leaveEvent); err != nil {
		log.Printf("Failed to dispatch member left event: %v", err)
	}

	// 发布领域事件
	s.publishDomainEvents(ctx, community)

	return nil
}

// GetCommunityHealth 获取社区健康度
func (s *CommunityApplicationService) GetCommunityHealth(ctx context.Context, query *dto.GetCommunityHealthQuery) (*dto.CommunityHealthResponse, error) {
	// 获取社区
	communityID, err := valueobject.NewCommunityID(query.CommunityID)
	if err != nil {
		return nil, err
	}

	community, err := s.communityRepo.FindByID(ctx, communityID)
	if err != nil {
		return nil, err
	}
	if community == nil {
		return nil, errors.New("社区不存在")
	}

	// 计算健康度
	health := s.domainService.CheckCommunityHealth(community)

	return &dto.CommunityHealthResponse{
		CommunityID:   query.CommunityID,
		OverallScore:  health["overall_score"].(int),
		Indicators:    health["indicators"].(map[string]interface{}),
		Suggestions:   health["suggestions"].([]string),
	}, nil
}

// RecommendTags 推荐标签
func (s *CommunityApplicationService) RecommendTags(ctx context.Context, query *dto.RecommendTagsQuery) (*dto.RecommendTagsResponse, error) {
	tags := s.domainService.RecommendTags(query.Name, query.Description, query.Location)
	
	return &dto.RecommendTagsResponse{
		RecommendedTags: tags,
	}, nil
}

// publishDomainEvents 发布领域事件
func (s *CommunityApplicationService) publishDomainEvents(ctx context.Context, community *entity.Community) {
	events := community.GetDomainEvents()
	for _, domainEvent := range events {
		if err := s.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			log.Printf("Failed to dispatch domain event: %v", err)
		}
	}
	community.ClearDomainEvents()
}

// toDTO 转换为DTO
func (s *CommunityApplicationService) toDTO(community *entity.Community, ownerName string) *dto.CommunityResponse {
	return &dto.CommunityResponse{
		ID:          community.ID().Value(),
		Name:        community.Name().Value(),
		Description: community.Description(),
		OwnerID:     community.OwnerID(),
		OwnerName:   ownerName,
		Status:      community.Status().String(),
		StatusCode:  community.Status().Value(),
		Tags:        community.Tags().Values(),
		Location:    community.Location(),
		GroupCount:  community.GroupCount(),
		MemberCount: community.MemberCount(),
		PostCount:   community.PostCount(),
		CreatedAt:   community.CreatedAt(),
		UpdatedAt:   community.UpdatedAt(),
	}
}
