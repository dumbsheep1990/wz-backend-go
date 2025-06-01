package repository

import (
	"context"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// RoleRepository defines the interface for role persistence operations
type RoleRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id valueobject.RoleID) (*entity.Role, error)
	
	// FindByName finds a role by name
	FindByName(ctx context.Context, name valueobject.RoleName) (*entity.Role, error)
	
	// Save persists a role entity
	Save(ctx context.Context, role *entity.Role) error
	
	// Create creates a new role entity
	Create(ctx context.Context, role *entity.Role) error
	
	// Update updates an existing role entity
	Update(ctx context.Context, role *entity.Role) error
	
	// Delete deletes a role entity
	Delete(ctx context.Context, id valueobject.RoleID) error
	
	// List lists role entities with pagination
	List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error)
	
	// FindAllByIDs finds all roles by their IDs
	FindAllByIDs(ctx context.Context, ids []valueobject.RoleID) ([]*entity.Role, error)
	
	// FindByAdminID finds roles associated with an admin
	FindByAdminID(ctx context.Context, adminID valueobject.AdminID) ([]*entity.Role, error)
	
	// AddPermission adds a permission to a role
	AddPermission(ctx context.Context, roleID valueobject.RoleID, permission valueobject.Permission) error
	
	// RemovePermission removes a permission from a role
	RemovePermission(ctx context.Context, roleID valueobject.RoleID, permission valueobject.Permission) error
	
	// ClearPermissions removes all permissions from a role
	ClearPermissions(ctx context.Context, roleID valueobject.RoleID) error
	
	// SetPermissions sets the complete list of permissions for a role
	SetPermissions(ctx context.Context, roleID valueobject.RoleID, permissions []valueobject.Permission) error
	
	// GetPermissions gets all permissions for a role
	GetPermissions(ctx context.Context, roleID valueobject.RoleID) ([]valueobject.Permission, error)
}
