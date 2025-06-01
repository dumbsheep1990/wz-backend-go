package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/repository"
	"wz-backend-go/internal/domain/admin/valueobject"
	"wz-backend-go/internal/infrastructure/event"
)

// RoleDomainService encapsulates the domain logic for role operations
type RoleDomainService struct {
	roleRepo  repository.RoleRepository
	adminRepo repository.AdminRepository
	eventBus  event.EventBus
}

// NewRoleDomainService creates a new RoleDomainService
func NewRoleDomainService(
	roleRepo repository.RoleRepository,
	adminRepo repository.AdminRepository,
	eventBus event.EventBus,
) *RoleDomainService {
	return &RoleDomainService{
		roleRepo:  roleRepo,
		adminRepo: adminRepo,
		eventBus:  eventBus,
	}
}

// CreateRole creates a new role with the specified permissions
func (s *RoleDomainService) CreateRole(
	ctx context.Context,
	id valueobject.RoleID,
	name valueobject.RoleName,
	description string,
	parentID valueobject.RoleID,
	permissions []valueobject.Permission,
) (*entity.Role, error) {
	// Check if the role ID already exists
	existingRole, err := s.roleRepo.FindByID(ctx, id)
	if err == nil && existingRole != nil {
		return nil, errors.New("role ID already exists")
	}

	// Check if the role name already exists
	existingRole, err = s.roleRepo.FindByName(ctx, name)
	if err == nil && existingRole != nil {
		return nil, errors.New("role name already exists")
	}

	// If a parent role is specified, check that it exists
	if parentID.Value() != "" {
		parentRole, err := s.roleRepo.FindByID(ctx, parentID)
		if err != nil {
			return nil, err
		}
		if parentRole == nil {
			return nil, errors.New("parent role not found")
		}
	}

	// Create the role entity
	role, err := entity.NewRole(id, name, description, parentID)
	if err != nil {
		return nil, err
	}

	// Add permissions if provided
	if len(permissions) > 0 {
		role.ReplacePermissions(permissions)
	}

	// Persist the role entity
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	// Publish the role created event
	roleCreatedEvent := entity.NewRoleCreatedEvent(role)
	s.eventBus.Publish(roleCreatedEvent)

	return role, nil
}

// UpdateRole updates an existing role's basic information
func (s *RoleDomainService) UpdateRole(
	ctx context.Context,
	id valueobject.RoleID,
	name valueobject.RoleName,
	description string,
	parentID valueobject.RoleID,
) (*entity.Role, error) {
	// Retrieve the role
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Check if the new name already exists (if different from current)
	if !role.Name.Equals(name) {
		existingRole, err := s.roleRepo.FindByName(ctx, name)
		if err == nil && existingRole != nil && !existingRole.ID.Equals(id) {
			return nil, errors.New("role name already exists")
		}
	}

	// If a parent role is specified, check that it exists
	if parentID.Value() != "" && !parentID.Equals(role.ParentID) {
		parentRole, err := s.roleRepo.FindByID(ctx, parentID)
		if err != nil {
			return nil, err
		}
		if parentRole == nil {
			return nil, errors.New("parent role not found")
		}

		// Prevent circular dependencies in hierarchy
		if parentID.Equals(id) {
			return nil, errors.New("a role cannot be its own parent")
		}
	}

	// Update role properties
	role.UpdateName(name)
	role.UpdateDescription(description)
	role.SetParent(parentID)

	// Persist the changes
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	// Publish the role updated event
	roleUpdatedEvent := entity.NewRoleUpdatedEvent(role)
	s.eventBus.Publish(roleUpdatedEvent)

	return role, nil
}

// AddPermission adds a permission to a role
func (s *RoleDomainService) AddPermission(
	ctx context.Context,
	roleID valueobject.RoleID,
	permission valueobject.Permission,
) error {
	// Retrieve the role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Add the permission
	if err := role.AddPermission(permission); err != nil {
		return err
	}

	// Persist the changes
	if err := s.roleRepo.AddPermission(ctx, roleID, permission); err != nil {
		return err
	}

	// Publish the permissions changed event
	permissionsChangedEvent := entity.NewRolePermissionsChangedEvent(role)
	s.eventBus.Publish(permissionsChangedEvent)

	return nil
}

// RemovePermission removes a permission from a role
func (s *RoleDomainService) RemovePermission(
	ctx context.Context,
	roleID valueobject.RoleID,
	permission valueobject.Permission,
) error {
	// Retrieve the role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Remove the permission
	if err := role.RemovePermission(permission); err != nil {
		return err
	}

	// Persist the changes
	if err := s.roleRepo.RemovePermission(ctx, roleID, permission); err != nil {
		return err
	}

	// Publish the permissions changed event
	permissionsChangedEvent := entity.NewRolePermissionsChangedEvent(role)
	s.eventBus.Publish(permissionsChangedEvent)

	return nil
}

// SetPermissions replaces all permissions for a role
func (s *RoleDomainService) SetPermissions(
	ctx context.Context,
	roleID valueobject.RoleID,
	permissions []valueobject.Permission,
) error {
	// Retrieve the role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Replace the permissions
	role.ReplacePermissions(permissions)

	// Persist the changes
	if err := s.roleRepo.SetPermissions(ctx, roleID, permissions); err != nil {
		return err
	}

	// Publish the permissions changed event
	permissionsChangedEvent := entity.NewRolePermissionsChangedEvent(role)
	s.eventBus.Publish(permissionsChangedEvent)

	return nil
}

// DeleteRole deletes a role if it's not assigned to any admins
func (s *RoleDomainService) DeleteRole(
	ctx context.Context,
	roleID valueobject.RoleID,
) error {
	// Retrieve the role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Check if any admins are using this role
	// This would require extending the AdminRepository to add a method to count admins by role
	adminsCount, err := s.countAdminsByRole(ctx, roleID)
	if err != nil {
		return err
	}
	if adminsCount > 0 {
		return fmt.Errorf("cannot delete role: it is assigned to %d admins", adminsCount)
	}

	// Delete the role
	if err := s.roleRepo.Delete(ctx, roleID); err != nil {
		return err
	}

	return nil
}

// countAdminsByRole counts the number of admins assigned to a specific role
// This is a helper method that would need to be implemented in AdminRepository
func (s *RoleDomainService) countAdminsByRole(
	ctx context.Context,
	roleID valueobject.RoleID,
) (int64, error) {
	// This would need a specific method in the admin repository
	// For now, we'll use a filter on the CountAdmins method
	filters := map[string]interface{}{
		"role_id": roleID.Value(),
	}
	return s.adminRepo.CountAdmins(ctx, filters)
}

// GetRoleWithPermissions retrieves a role with its permissions
func (s *RoleDomainService) GetRoleWithPermissions(
	ctx context.Context,
	roleID valueobject.RoleID,
) (*entity.Role, error) {
	// Retrieve the role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Retrieve permissions for the role
	permissions, err := s.roleRepo.GetPermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Update the role with its permissions
	role.ReplacePermissions(permissions)

	return role, nil
}

// GetRoleHierarchy retrieves a role with its entire hierarchy (parent roles)
func (s *RoleDomainService) GetRoleHierarchy(
	ctx context.Context,
	roleID valueobject.RoleID,
) ([]*entity.Role, error) {
	// Retrieve the role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Start the hierarchy with the current role
	hierarchy := []*entity.Role{role}

	// Recursively add parent roles
	currentRoleID := role.ParentID
	visited := map[string]bool{role.ID.Value(): true}

	for currentRoleID.Value() != "" {
		// Check for circular dependencies
		if visited[currentRoleID.Value()] {
			return nil, errors.New("circular dependency detected in role hierarchy")
		}

		parentRole, err := s.roleRepo.FindByID(ctx, currentRoleID)
		if err != nil {
			return nil, err
		}
		if parentRole == nil {
			break
		}

		// Add to hierarchy and mark as visited
		hierarchy = append(hierarchy, parentRole)
		visited[parentRole.ID.Value()] = true

		// Move up the hierarchy
		currentRoleID = parentRole.ParentID
	}

	return hierarchy, nil
}
