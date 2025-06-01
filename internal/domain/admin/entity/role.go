package entity

import (
	"errors"
	"time"

	"wz-backend-go/internal/domain/admin/valueobject"
)

// Role represents a role entity in the admin domain
type Role struct {
	ID          valueobject.RoleID
	Name        valueobject.RoleName
	Description string
	Permissions []valueobject.Permission
	ParentID    valueobject.RoleID // For hierarchical roles
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewRole creates a new Role entity
func NewRole(
	id valueobject.RoleID,
	name valueobject.RoleName,
	description string,
	parentID valueobject.RoleID,
) (*Role, error) {
	now := time.Now()
	
	role := &Role{
		ID:          id,
		Name:        name,
		Description: description,
		ParentID:    parentID,
		Permissions: make([]valueobject.Permission, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	return role, nil
}

// AddPermission adds a permission to the role
func (r *Role) AddPermission(permission valueobject.Permission) error {
	// Check if permission already exists
	for _, p := range r.Permissions {
		if p.Equals(permission) {
			return errors.New("permission already exists for this role")
		}
	}
	
	r.Permissions = append(r.Permissions, permission)
	r.UpdatedAt = time.Now()
	return nil
}

// RemovePermission removes a permission from the role
func (r *Role) RemovePermission(permission valueobject.Permission) error {
	for i, p := range r.Permissions {
		if p.Equals(permission) {
			// Remove the permission by replacing it with the last one and truncating
			lastIndex := len(r.Permissions) - 1
			r.Permissions[i] = r.Permissions[lastIndex]
			r.Permissions = r.Permissions[:lastIndex]
			r.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return errors.New("permission not found for this role")
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permission valueobject.Permission) bool {
	for _, p := range r.Permissions {
		if p.Equals(permission) {
			return true
		}
	}
	
	// Check for wildcard permissions (e.g., "resource:*")
	resourceType := permission.GetResourceType()
	wildcard, _ := valueobject.NewPermission(resourceType, "*")
	
	for _, p := range r.Permissions {
		if p.Equals(wildcard) || p.GetResourceType() == resourceType && p.GetActionType() == "*" {
			return true
		}
	}
	
	return false
}

// UpdateName updates the role name
func (r *Role) UpdateName(name valueobject.RoleName) {
	r.Name = name
	r.UpdatedAt = time.Now()
}

// UpdateDescription updates the role description
func (r *Role) UpdateDescription(description string) {
	r.Description = description
	r.UpdatedAt = time.Now()
}

// SetParent sets the parent role
func (r *Role) SetParent(parentID valueobject.RoleID) {
	r.ParentID = parentID
	r.UpdatedAt = time.Now()
}

// ClearPermissions removes all permissions from the role
func (r *Role) ClearPermissions() {
	r.Permissions = make([]valueobject.Permission, 0)
	r.UpdatedAt = time.Now()
}

// ReplacePermissions replaces all permissions with a new set
func (r *Role) ReplacePermissions(permissions []valueobject.Permission) {
	r.Permissions = permissions
	r.UpdatedAt = time.Now()
}
