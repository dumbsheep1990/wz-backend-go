package service

import (
	"context"
	"fmt"
	"time"
	"github.com/go-playground/validator/v10"
	"wz-backend-go/internal/application/component/dto"
	"wz-backend-go/internal/domain/component/entity"
	"wz-backend-go/internal/domain/component/repository"
	componentService "wz-backend-go/internal/domain/component/service"
	"wz-backend-go/internal/domain/component/valueobject"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/database"
)

// ComponentApplicationService 组件应用服务
type ComponentApplicationService struct {
	componentRepo repository.ComponentRepository
	domainService *componentService.ComponentDomainService
	eventBus      event.EventBus
	validator     *validator.Validate
	unitOfWork    database.UnitOfWork
}

// NewComponentApplicationService 创建组件应用服务
func NewComponentApplicationService(
	componentRepo repository.ComponentRepository,
	domainService *componentService.ComponentDomainService,
	eventBus event.EventBus,
	validator *validator.Validate,
	unitOfWork database.UnitOfWork,
) *ComponentApplicationService {
	return &ComponentApplicationService{
		componentRepo: componentRepo,
		domainService: domainService,
		eventBus:      eventBus,
		validator:     validator,
		unitOfWork:    unitOfWork,
	}
}

// CreateComponent 创建组件
func (s *ComponentApplicationService) CreateComponent(ctx context.Context, req dto.CreateComponentRequest, tenantID string) (*dto.ComponentDTO, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 验证名称唯一性
	if err := s.domainService.ValidateUniqueName(ctx, req.Name, tenantID); err != nil {
		return nil, err
	}
	
	// 验证模板内容
	if err := s.domainService.ValidateTemplateContent(req.Template); err != nil {
		return nil, err
	}
	
	// 验证配置内容
	if err := s.domainService.ValidateConfigContent(req.Config); err != nil {
		return nil, err
	}
	
	// 创建值对象
	componentID, err := valueobject.NewComponentID(generateComponentID())
	if err != nil {
		return nil, fmt.Errorf("创建组件ID失败: %w", err)
	}
	
	componentType, err := valueobject.NewComponentType(req.ComponentType)
	if err != nil {
		return nil, fmt.Errorf("创建组件类型失败: %w", err)
	}
	
	// 创建组件聚合
	component, err := entity.NewComponent(componentID, req.Name, componentType, req.Template, tenantID)
	if err != nil {
		return nil, fmt.Errorf("创建组件失败: %w", err)
	}
	
	// 设置可选属性
	if req.Description != "" {
		component.UpdateDescription(req.Description)
	}
	if req.Config != "" {
		component.UpdateConfig(req.Config)
	}
	if req.Preview != "" {
		component.UpdatePreview(req.Preview)
	}
	if req.Category != "" {
		component.UpdateCategory(req.Category)
	}
	if len(req.Tags) > 0 {
		component.UpdateTags(req.Tags)
	}
	if req.IsPublic {
		component.MakePublic()
	}
	
	// 在事务中保存组件
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.componentRepo.Save(ctx, component); err != nil {
			return fmt.Errorf("保存组件失败: %w", err)
		}
		
		// 发布领域事件
		events := component.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		component.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return s.componentToDTO(component), nil
}

// UpdateComponent 更新组件
func (s *ComponentApplicationService) UpdateComponent(ctx context.Context, componentID string, req dto.UpdateComponentRequest, tenantID string) (*dto.ComponentDTO, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 创建组件ID
	id, err := valueobject.NewComponentID(componentID)
	if err != nil {
		return nil, fmt.Errorf("无效的组件ID: %w", err)
	}
	
	// 验证组件所有权
	component, err := s.domainService.ValidateOwnership(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	
	// 更新组件信息
	if req.Name != nil {
		// 验证名称唯一性
		if err := s.domainService.ValidateUniqueNameForUpdate(ctx, *req.Name, tenantID, id); err != nil {
			return nil, err
		}
		if err := component.UpdateName(*req.Name); err != nil {
			return nil, fmt.Errorf("更新组件名称失败: %w", err)
		}
	}
	
	if req.Description != nil {
		component.UpdateDescription(*req.Description)
	}
	
	if req.Template != nil {
		// 验证模板内容
		if err := s.domainService.ValidateTemplateContent(*req.Template); err != nil {
			return nil, err
		}
		if err := component.UpdateTemplate(*req.Template); err != nil {
			return nil, fmt.Errorf("更新模板失败: %w", err)
		}
	}
	
	if req.Config != nil {
		// 验证配置内容
		if err := s.domainService.ValidateConfigContent(*req.Config); err != nil {
			return nil, err
		}
		component.UpdateConfig(*req.Config)
	}
	
	if req.Preview != nil {
		component.UpdatePreview(*req.Preview)
	}
	
	if req.Category != nil {
		component.UpdateCategory(*req.Category)
	}
	
	if req.Tags != nil {
		component.UpdateTags(req.Tags)
	}
	
	// 在事务中保存组件
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.componentRepo.Save(ctx, component); err != nil {
			return fmt.Errorf("保存组件失败: %w", err)
		}
		
		// 发布领域事件
		events := component.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		component.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return s.componentToDTO(component), nil
}

// GetComponent 获取组件详情
func (s *ComponentApplicationService) GetComponent(ctx context.Context, componentID, tenantID string) (*dto.ComponentDTO, error) {
	id, err := valueobject.NewComponentID(componentID)
	if err != nil {
		return nil, fmt.Errorf("无效的组件ID: %w", err)
	}
	
	component, err := s.domainService.ValidateAccess(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	
	return s.componentToDTO(component), nil
}

// ListComponents 获取组件列表
func (s *ComponentApplicationService) ListComponents(ctx context.Context, req dto.ComponentListRequest, tenantID string) (*dto.ComponentListResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}
	
	filters := repository.ComponentFilters{
		ComponentType: req.ComponentType,
		Category:      req.Category,
		IsPublic:      req.IsPublic,
		Search:        req.Search,
		Tags:          req.Tags,
		TenantID:      tenantID,
		Limit:         req.Size,
		Offset:        (req.Page - 1) * req.Size,
		SortBy:        req.SortBy,
		SortOrder:     req.SortOrder,
	}
	
	components, err := s.componentRepo.FindByTenant(ctx, tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("查询组件列表失败: %w", err)
	}
	
	total, err := s.componentRepo.Count(ctx, tenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("统计组件数量失败: %w", err)
	}
	
	componentDTOs := make([]dto.ComponentDTO, 0, len(components))
	for _, component := range components {
		componentDTOs = append(componentDTOs, *s.componentToDTO(component))
	}
	
	return &dto.ComponentListResponse{
		Components: componentDTOs,
		Total:      total,
		Page:       req.Page,
		Size:       req.Size,
	}, nil
}

// ListPublicComponents 获取公开组件列表
func (s *ComponentApplicationService) ListPublicComponents(ctx context.Context, req dto.PublicComponentListRequest) (*dto.PublicComponentListResponse, error) {
	// 验证请求参数
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}
	
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}
	
	isPublic := true
	filters := repository.ComponentFilters{
		ComponentType: req.ComponentType,
		Category:      req.Category,
		IsPublic:      &isPublic,
		Search:        req.Search,
		Tags:          req.Tags,
		Limit:         req.Size,
		Offset:        (req.Page - 1) * req.Size,
		SortBy:        req.SortBy,
		SortOrder:     req.SortOrder,
	}
	
	components, err := s.componentRepo.FindPublicComponents(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("查询公开组件列表失败: %w", err)
	}
	
	total, err := s.componentRepo.CountPublic(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("统计公开组件数量失败: %w", err)
	}
	
	componentDTOs := make([]dto.ComponentDTO, 0, len(components))
	for _, component := range components {
		componentDTOs = append(componentDTOs, *s.componentToDTO(component))
	}
	
	return &dto.PublicComponentListResponse{
		Components: componentDTOs,
		Total:      total,
		Page:       req.Page,
		Size:       req.Size,
	}, nil
}

// MakeComponentPublic 设置组件为公开
func (s *ComponentApplicationService) MakeComponentPublic(ctx context.Context, componentID, tenantID string) (*dto.MakePublicResponse, error) {
	id, err := valueobject.NewComponentID(componentID)
	if err != nil {
		return nil, fmt.Errorf("无效的组件ID: %w", err)
	}
	
	// 验证组件所有权
	component, err := s.domainService.ValidateOwnership(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	
	// 设置为公开
	component.MakePublic()
	
	// 在事务中保存组件
	err = s.unitOfWork.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.componentRepo.Save(ctx, component); err != nil {
			return fmt.Errorf("保存组件失败: %w", err)
		}
		
		// 发布领域事件
		events := component.GetDomainEvents()
		for _, event := range events {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布事件失败: %w", err)
			}
		}
		component.ClearDomainEvents()
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &dto.MakePublicResponse{
		Component: *s.componentToDTO(component),
		Message:   "组件已设置为公开",
	}, nil
}

// DeleteComponent 删除组件
func (s *ComponentApplicationService) DeleteComponent(ctx context.Context, componentID, tenantID string) error {
	id, err := valueobject.NewComponentID(componentID)
	if err != nil {
		return fmt.Errorf("无效的组件ID: %w", err)
	}
	
	// 验证是否可以删除
	if err := s.domainService.CanDeleteComponent(ctx, id, tenantID); err != nil {
		return err
	}
	
	// 删除组件
	if err := s.componentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除组件失败: %w", err)
	}
	
	return nil
}

// componentToDTO 将组件实体转换为DTO
func (s *ComponentApplicationService) componentToDTO(component *entity.Component) *dto.ComponentDTO {
	return &dto.ComponentDTO{
		ID:            component.ID().Value(),
		Name:          component.Name(),
		Description:   component.Description(),
		ComponentType: component.ComponentType().Value(),
		Template:      component.Template(),
		Config:        component.Config(),
		Preview:       component.Preview(),
		Category:      component.Category(),
		Tags:          component.Tags(),
		IsPublic:      component.IsPublic(),
		TenantID:      component.TenantID(),
		Version:       component.Version(),
		CreatedAt:     component.CreatedAt(),
		UpdatedAt:     component.UpdatedAt(),
	}
}

// generateComponentID 生成组件ID (简单实现，实际应使用UUID)
func generateComponentID() string {
	// 这里应该使用更好的ID生成策略，比如UUID
	return fmt.Sprintf("comp_%d", time.Now().UnixNano())
} 