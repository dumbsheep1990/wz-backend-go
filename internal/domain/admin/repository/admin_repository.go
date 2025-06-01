package repository

import (
	"context"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// AdminRepository defines the interface for admin persistence operations
type AdminRepository interface {
	// FindByID finds an admin by ID
	FindByID(ctx context.Context, id valueobject.AdminID) (*entity.Admin, error)
	
	// FindByUsername finds an admin by username
	FindByUsername(ctx context.Context, username valueobject.Username) (*entity.Admin, error)
	
	// Save persists an admin entity
	Save(ctx context.Context, admin *entity.Admin) error
	
	// Create creates a new admin entity
	Create(ctx context.Context, admin *entity.Admin) error
	
	// Update updates an existing admin entity
	Update(ctx context.Context, admin *entity.Admin) error
	
	// Delete deletes an admin entity
	Delete(ctx context.Context, id valueobject.AdminID) error
	
	// List lists admin entities with pagination and filters
	List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*entity.Admin, int64, error)
	
	// UpdateLastLogin updates the last login time for an admin
	UpdateLastLogin(ctx context.Context, id valueobject.AdminID) error
	
	// CountAdmins counts the total number of admins, optionally filtered
	CountAdmins(ctx context.Context, filters map[string]interface{}) (int64, error)
}
