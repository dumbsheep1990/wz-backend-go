package service

import (
	"context"
	"errors"
	"fmt"
	"wz-backend-go/internal/domain/component/entity"
	"wz-backend-go/internal/domain/component/repository"
	"wz-backend-go/internal/domain/component/valueobject"
)

// ComponentDomainService 组件领域服务
type ComponentDomainService struct {
	componentRepo repository.ComponentRepository
}

// NewComponentDomainService 创建组件领域服务
func NewComponentDomainService(componentRepo repository.ComponentRepository) *ComponentDomainService {
	return &ComponentDomainService{
		componentRepo: componentRepo,
	}
}

// ValidateUniqueName 验证组件名称唯一性（在租户内）
func (s *ComponentDomainService) ValidateUniqueName(ctx context.Context, name, tenantID string) error {
	exists, err := s.componentRepo.ExistsByName(ctx, name, tenantID)
	if err != nil {
		return fmt.Errorf("检查组件名称唯一性失败: %w", err)
	}
	if exists {
		return errors.New("组件名称已被使用")
	}
	return nil
}

// ValidateUniqueNameForUpdate 验证更新时的组件名称唯一性
func (s *ComponentDomainService) ValidateUniqueNameForUpdate(ctx context.Context, name, tenantID string, componentID valueobject.ComponentID) error {
	exists, err := s.componentRepo.ExistsByNameExcludeID(ctx, name, tenantID, componentID)
	if err != nil {
		return fmt.Errorf("检查组件名称唯一性失败: %w", err)
	}
	if exists {
		return errors.New("组件名称已被使用")
	}
	return nil
}

// ValidateOwnership 验证组件所有权
func (s *ComponentDomainService) ValidateOwnership(ctx context.Context, componentID valueobject.ComponentID, tenantID string) (*entity.Component, error) {
	component, err := s.componentRepo.FindByIDAndTenant(ctx, componentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查找组件失败: %w", err)
	}
	if component == nil {
		return nil, errors.New("组件不存在或无权限")
	}
	
	if !component.IsOwnedBy(tenantID) {
		return nil, errors.New("无权限操作此组件")
	}
	
	return component, nil
}

// ValidateAccess 验证组件访问权限（包括公开组件）
func (s *ComponentDomainService) ValidateAccess(ctx context.Context, componentID valueobject.ComponentID, tenantID string) (*entity.Component, error) {
	component, err := s.componentRepo.FindByID(ctx, componentID)
	if err != nil {
		return nil, fmt.Errorf("查找组件失败: %w", err)
	}
	if component == nil {
		return nil, errors.New("组件不存在")
	}
	
	if !component.CanBeAccessed(tenantID) {
		return nil, errors.New("无权限访问此组件")
	}
	
	return component, nil
}

// CanDeleteComponent 检查是否可以删除组件
func (s *ComponentDomainService) CanDeleteComponent(ctx context.Context, componentID valueobject.ComponentID, tenantID string) error {
	component, err := s.ValidateOwnership(ctx, componentID, tenantID)
	if err != nil {
		return err
	}
	
	// 检查组件是否被使用（这里可以扩展检查是否有页面在使用此组件）
	// TODO: 实现组件使用情况检查
	_ = component
	
	return nil
}

// ValidatePublicAccess 验证公开组件的访问
func (s *ComponentDomainService) ValidatePublicAccess(ctx context.Context, componentID valueobject.ComponentID) (*entity.Component, error) {
	component, err := s.componentRepo.FindByID(ctx, componentID)
	if err != nil {
		return nil, fmt.Errorf("查找组件失败: %w", err)
	}
	if component == nil {
		return nil, errors.New("组件不存在")
	}
	
	if !component.IsPublic() {
		return nil, errors.New("组件不是公开的")
	}
	
	return component, nil
}

// GetAvailableComponents 获取可用组件（自有 + 公开）
func (s *ComponentDomainService) GetAvailableComponents(ctx context.Context, tenantID string, componentType valueobject.ComponentType, filters repository.ComponentFilters) ([]*entity.Component, error) {
	// 获取自有组件
	filters.TenantID = tenantID
	ownComponents, err := s.componentRepo.FindByType(ctx, componentType, filters)
	if err != nil {
		return nil, fmt.Errorf("查询自有组件失败: %w", err)
	}
	
	// 获取公开组件
	publicFilters := filters
	publicFilters.TenantID = ""
	isPublic := true
	publicFilters.IsPublic = &isPublic
	publicComponents, err := s.componentRepo.FindByType(ctx, componentType, publicFilters)
	if err != nil {
		return nil, fmt.Errorf("查询公开组件失败: %w", err)
	}
	
	// 合并结果，去重
	componentMap := make(map[string]*entity.Component)
	
	// 先添加自有组件（优先级高）
	for _, component := range ownComponents {
		componentMap[component.ID().Value()] = component
	}
	
	// 再添加公开组件（不覆盖自有的）
	for _, component := range publicComponents {
		if _, exists := componentMap[component.ID().Value()]; !exists {
			componentMap[component.ID().Value()] = component
		}
	}
	
	// 转换为切片
	result := make([]*entity.Component, 0, len(componentMap))
	for _, component := range componentMap {
		result = append(result, component)
	}
	
	return result, nil
}

// ValidateTemplateContent 验证模板内容
func (s *ComponentDomainService) ValidateTemplateContent(template string) error {
	if template == "" {
		return errors.New("模板内容不能为空")
	}
	
	// 这里可以添加更复杂的模板验证逻辑
	// 比如检查HTML语法、安全性等
	if len(template) > 100000 { // 100KB限制
		return errors.New("模板内容过大，不能超过100KB")
	}
	
	return nil
}

// ValidateConfigContent 验证配置内容
func (s *ComponentDomainService) ValidateConfigContent(config string) error {
	if config == "" {
		return nil // 配置可以为空
	}
	
	// 这里可以添加JSON格式验证
	if len(config) > 50000 { // 50KB限制
		return errors.New("配置内容过大，不能超过50KB")
	}
	
	return nil
} 