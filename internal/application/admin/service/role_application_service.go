package service

import (
	"context"
	"errors"
	"strings"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/repository"
	domainService "wz-backend-go/internal/domain/admin/service"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// RoleApplicationService coordinates domain services for role-related operations
type RoleApplicationService struct {
	roleDomainService *domainService.RoleDomainService
	roleRepo          repository.RoleRepository
	adminRepo         repository.AdminRepository
}

// NewRoleApplicationService creates a new RoleApplicationService
func NewRoleApplicationService(
	roleDomainService *domainService.RoleDomainService,
	roleRepo repository.RoleRepository,
	adminRepo repository.AdminRepository,
) *RoleApplicationService {
	return &RoleApplicationService{
		roleDomainService: roleDomainService,
		roleRepo:          roleRepo,
		adminRepo:         adminRepo,
	}
}

// CreateRole creates a new role
func (s *RoleApplicationService) CreateRole(
	ctx context.Context,
	req dto.RoleCreateRequest,
) (*dto.RoleDTO, error) {
	// Convert request to domain objects
	roleID, err := valueobject.NewRoleID(req.ID)
	if err != nil {
		return nil, err
	}

	roleName, err := valueobject.NewRoleName(req.Name)
	if err != nil {
		return nil, err
	}

	var parentID valueobject.RoleID
	if req.ParentID != "" {
		parentID, err = valueobject.NewRoleID(req.ParentID)
		if err != nil {
			return nil, err
		}
	}

	// Convert permission strings to Permission value objects
	permissions := make([]valueobject.Permission, 0, len(req.Permissions))
	for _, permStr := range req.Permissions {
		perm, err := valueobject.ParsePermission(permStr)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	// Use domain service to create the role
	role, err := s.roleDomainService.CreateRole(
		ctx,
		roleID,
		roleName,
		req.Description,
		parentID,
		permissions,
	)
	if err != nil {
		return nil, err
	}

	// Convert domain object to DTO
	roleDTO := dto.MapRoleToDTO(role)
	return &roleDTO, nil
}

// UpdateRole updates an existing role
func (s *RoleApplicationService) UpdateRole(
	ctx context.Context,
	id string,
	req dto.RoleUpdateRequest,
) (*dto.RoleDTO, error) {
	roleID, err := valueobject.NewRoleID(id)
	if err != nil {
		return nil, err
	}

	// Retrieve existing role to update
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Update role name if provided
	var roleName valueobject.RoleName
	if req.Name != "" {
		roleName, err = valueobject.NewRoleName(req.Name)
		if err != nil {
			return nil, err
		}
	} else {
		roleName = role.Name
	}

	// Update parent ID if provided
	var parentID valueobject.RoleID
	if req.ParentID != "" {
		parentID, err = valueobject.NewRoleID(req.ParentID)
		if err != nil {
			return nil, err
		}
	} else {
		parentID = role.ParentID
	}

	// Use description from request or existing role
	description := req.Description
	if description == "" {
		description = role.Description
	}

	// Update the role using domain service
	updatedRole, err := s.roleDomainService.UpdateRole(
		ctx,
		roleID,
		roleName,
		description,
		parentID,
	)
	if err != nil {
		return nil, err
	}

	// Update permissions if provided
	if len(req.Permissions) > 0 {
		permissions := make([]valueobject.Permission, 0, len(req.Permissions))
		for _, permStr := range req.Permissions {
			perm, err := valueobject.ParsePermission(permStr)
			if err != nil {
				return nil, err
			}
			permissions = append(permissions, perm)
		}

		if err := s.roleDomainService.SetPermissions(ctx, roleID, permissions); err != nil {
			return nil, err
		}

		// Reload the role with updated permissions
		updatedRole, err = s.roleDomainService.GetRoleWithPermissions(ctx, roleID)
		if err != nil {
			return nil, err
		}
	}

	// Convert domain object to DTO
	roleDTO := dto.MapRoleToDTO(updatedRole)
	return &roleDTO, nil
}

// DeleteRole deletes a role
func (s *RoleApplicationService) DeleteRole(ctx context.Context, id string) error {
	roleID, err := valueobject.NewRoleID(id)
	if err != nil {
		return err
	}

	return s.roleDomainService.DeleteRole(ctx, roleID)
}

// GetRole retrieves a role by ID
func (s *RoleApplicationService) GetRole(ctx context.Context, id string) (*dto.RoleDTO, error) {
	roleID, err := valueobject.NewRoleID(id)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	roleDTO := dto.MapRoleToDTO(role)
	return &roleDTO, nil
}

// GetRoleDetail retrieves detailed information about a role
func (s *RoleApplicationService) GetRoleDetail(ctx context.Context, id string) (*dto.RoleDetailResponse, error) {
	roleID, err := valueobject.NewRoleID(id)
	if err != nil {
		return nil, err
	}

	// Retrieve role with permissions
	role, err := s.roleDomainService.GetRoleWithPermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Count admins with this role
	filters := map[string]interface{}{
		"role_id": role.ID.Value(),
	}
	adminCount, err := s.adminRepo.CountAdmins(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Convert domain objects to DTOs
	roleDTO := dto.MapRoleToDTO(role)
	permissionDTOs := dto.MapPermissionsToDTO(role.Permissions)

	// Create response
	response := &dto.RoleDetailResponse{
		Role:        roleDTO,
		Permissions: permissionDTOs,
		AdminCount:  adminCount,
	}

	return response, nil
}

// ListRoles retrieves a paginated list of roles
func (s *RoleApplicationService) ListRoles(
	ctx context.Context,
	page, pageSize int,
) (*dto.RoleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Use repository to retrieve roles
	roles, total, err := s.roleRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert domain objects to DTOs
	roleDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		roleDTO := dto.MapRoleToDTO(role)
		
		// If the role has a parent, get the parent name
		if role.ParentID.Value() != "" {
			parentRole, err := s.roleRepo.FindByID(ctx, role.ParentID)
			if err == nil && parentRole != nil {
				roleDTO.ParentName = parentRole.Name.Value()
			}
		}
		
		roleDTOs = append(roleDTOs, roleDTO)
	}

	// Create response
	response := &dto.RoleListResponse{
		Total: total,
		Items: roleDTOs,
	}

	return response, nil
}

// AddPermission adds a permission to a role
func (s *RoleApplicationService) AddPermission(
	ctx context.Context,
	roleID string,
	permissionStr string,
) error {
	id, err := valueobject.NewRoleID(roleID)
	if err != nil {
		return err
	}

	permission, err := valueobject.ParsePermission(permissionStr)
	if err != nil {
		return err
	}

	return s.roleDomainService.AddPermission(ctx, id, permission)
}

// RemovePermission removes a permission from a role
func (s *RoleApplicationService) RemovePermission(
	ctx context.Context,
	roleID string,
	permissionStr string,
) error {
	id, err := valueobject.NewRoleID(roleID)
	if err != nil {
		return err
	}

	permission, err := valueobject.ParsePermission(permissionStr)
	if err != nil {
		return err
	}

	return s.roleDomainService.RemovePermission(ctx, id, permission)
}

// GetPermissions retrieves all permissions for a role
func (s *RoleApplicationService) GetPermissions(
	ctx context.Context,
	roleID string,
) ([]dto.PermissionDTO, error) {
	id, err := valueobject.NewRoleID(roleID)
	if err != nil {
		return nil, err
	}

	// Get role with permissions
	role, err := s.roleDomainService.GetRoleWithPermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	permissionDTOs := dto.MapPermissionsToDTO(role.Permissions)
	return permissionDTOs, nil
}

// SetPermissions sets all permissions for a role
func (s *RoleApplicationService) SetPermissions(
	ctx context.Context,
	roleID string,
	permissionStrs []string,
) error {
	id, err := valueobject.NewRoleID(roleID)
	if err != nil {
		return err
	}

	// Convert permission strings to Permission value objects
	permissions := make([]valueobject.Permission, 0, len(permissionStrs))
	for _, permStr := range permissionStrs {
		perm, err := valueobject.ParsePermission(permStr)
		if err != nil {
			return err
		}
		permissions = append(permissions, perm)
	}

	return s.roleDomainService.SetPermissions(ctx, id, permissions)
}

// GetDefaultPermissions returns the default permissions for new roles
func (s *RoleApplicationService) GetDefaultPermissions() []dto.PermissionDTO {
	// Define default permissions based on common admin operations
	defaults := []string{
		"dashboard:view",
		"profile:view", 
		"profile:edit",
	}
	
	permissions := make([]dto.PermissionDTO, 0, len(defaults))
	for _, permStr := range defaults {
		parts := strings.Split(permStr, ":")
		if len(parts) == 2 {
			permissions = append(permissions, dto.PermissionDTO{
				Resource: parts[0],
				Action:   parts[1],
				Full:     permStr,
			})
		}
	}
	
	return permissions
}

// GetAvailablePermissions returns all available permissions in the system
func (s *RoleApplicationService) GetAvailablePermissions() []dto.PermissionDTO {
	// Define all available permissions in the system
	// This could be dynamically generated from API endpoints or configuration
	resources := []string{
		"dashboard", "admin", "role", "user", "content", 
		"category", "site", "page", "template", "render",
		"system", "settings", "logs", "profile",
	}
	
	actions := []string{
		"view", "create", "edit", "delete", "list", "publish", 
		"unpublish", "approve", "reject", "download", "upload",
	}
	
	// Special permissions for 万知 categories
	categories := []string{
		"同用", "同好", "同购", "同年", "同游", "同在", "同市", "同企", 
		"同亲", "同班", "同师", "同业", "同网", "同工", "同务", "同艺", 
		"同玩", "同闲", "同拍", "同乡", "同学",
	}
	
	permissions := make([]dto.PermissionDTO, 0)
	
	// Add standard resource:action permissions
	for _, resource := range resources {
		for _, action := range actions {
			permStr := resource + ":" + action
			permissions = append(permissions, dto.PermissionDTO{
				Resource: resource,
				Action:   action,
				Full:     permStr,
			})
		}
		
		// Add wildcard action for each resource
		permissions = append(permissions, dto.PermissionDTO{
			Resource: resource,
			Action:   "*",
			Full:     resource + ":*",
		})
	}
	
	// Add category-specific permissions
	for _, category := range categories {
		permStr := "category:" + category
		permissions = append(permissions, dto.PermissionDTO{
			Resource: "category",
			Action:   category,
			Full:     permStr,
		})
	}
	
	return permissions
}
