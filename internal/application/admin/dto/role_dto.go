package dto

import (
	"time"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// RoleDTO represents the data transfer object for a role
type RoleDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    string    `json:"parentId,omitempty"`
	ParentName  string    `json:"parentName,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// RoleWithPermissionsDTO extends RoleDTO with permissions
type RoleWithPermissionsDTO struct {
	RoleDTO
	Permissions []PermissionDTO `json:"permissions"`
}

// PermissionDTO represents a permission in the system
type PermissionDTO struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Full     string `json:"full"`
}

// RoleCreateRequest represents a request to create a role
type RoleCreateRequest struct {
	ID          string   `json:"id" binding:"required"`
	Name        string   `json:"name" binding:"required,min=2,max=32"`
	Description string   `json:"description"`
	ParentID    string   `json:"parentId,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// RoleUpdateRequest represents a request to update a role
type RoleUpdateRequest struct {
	Name        string   `json:"name,omitempty" binding:"omitempty,min=2,max=32"`
	Description string   `json:"description,omitempty"`
	ParentID    string   `json:"parentId,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// RoleListResponse represents a paginated list of roles
type RoleListResponse struct {
	Total int64     `json:"total"`
	Items []RoleDTO `json:"items"`
}

// RoleDetailResponse represents a detailed view of a role
type RoleDetailResponse struct {
	Role        RoleDTO        `json:"role"`
	Permissions []PermissionDTO `json:"permissions"`
	AdminCount  int64          `json:"adminCount"`
}

// MapRoleToDTO maps a Role entity to a RoleDTO
func MapRoleToDTO(role *entity.Role) RoleDTO {
	return RoleDTO{
		ID:          role.ID.Value(),
		Name:        role.Name.Value(),
		Description: role.Description,
		ParentID:    role.ParentID.Value(),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// MapRoleWithPermissionsToDTO maps a Role entity with permissions to a RoleWithPermissionsDTO
func MapRoleWithPermissionsToDTO(role *entity.Role) RoleWithPermissionsDTO {
	roleDTO := MapRoleToDTO(role)
	
	permissionDTOs := make([]PermissionDTO, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissionDTOs = append(permissionDTOs, PermissionDTO{
			Resource: permission.GetResourceType(),
			Action:   permission.GetActionType(),
			Full:     permission.String(),
		})
	}
	
	return RoleWithPermissionsDTO{
		RoleDTO:     roleDTO,
		Permissions: permissionDTOs,
	}
}

// MapPermissionToDTO maps a Permission value object to a PermissionDTO
func MapPermissionToDTO(permission valueobject.Permission) PermissionDTO {
	return PermissionDTO{
		Resource: permission.GetResourceType(),
		Action:   permission.GetActionType(),
		Full:     permission.String(),
	}
}

// MapPermissionsToDTO maps a slice of Permission value objects to a slice of PermissionDTO
func MapPermissionsToDTO(permissions []valueobject.Permission) []PermissionDTO {
	dtos := make([]PermissionDTO, 0, len(permissions))
	for _, permission := range permissions {
		dtos = append(dtos, MapPermissionToDTO(permission))
	}
	return dtos
}
