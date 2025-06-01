package service

import (
	"context"
	"errors"
	"time"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/repository"
	"wz-backend-go/internal/domain/admin/valueobject"
	"wz-backend-go/internal/infrastructure/event"
)

// AdminDomainService encapsulates the domain logic for admin operations
type AdminDomainService struct {
	adminRepo repository.AdminRepository
	roleRepo  repository.RoleRepository
	eventBus  event.EventBus
}

// NewAdminDomainService creates a new AdminDomainService
func NewAdminDomainService(
	adminRepo repository.AdminRepository,
	roleRepo repository.RoleRepository,
	eventBus event.EventBus,
) *AdminDomainService {
	return &AdminDomainService{
		adminRepo: adminRepo,
		roleRepo:  roleRepo,
		eventBus:  eventBus,
	}
}

// CreateAdmin creates a new admin with the specified role
func (s *AdminDomainService) CreateAdmin(
	ctx context.Context,
	username valueobject.Username,
	password string,
	roleID valueobject.RoleID,
) (*entity.Admin, error) {
	// Check if username already exists
	existingAdmin, err := s.adminRepo.FindByUsername(ctx, username)
	if err == nil && existingAdmin != nil {
		return nil, errors.New("username already exists")
	}

	// Verify that the role exists
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Generate a new AdminID
	adminID, err := valueobject.NewAdminID(time.Now().UnixNano())
	if err != nil {
		return nil, err
	}

	// Create the admin entity
	admin, err := entity.NewAdmin(
		adminID,
		username,
		password,
		roleID,
		valueobject.AdminStatusActive,
	)
	if err != nil {
		return nil, err
	}

	// Persist the admin entity
	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return nil, err
	}

	// Publish the admin created event
	adminCreatedEvent := entity.NewAdminCreatedEvent(admin)
	s.eventBus.Publish(adminCreatedEvent)

	return admin, nil
}

// UpdateAdminRole updates an admin's role
func (s *AdminDomainService) UpdateAdminRole(
	ctx context.Context,
	adminID valueobject.AdminID,
	newRoleID valueobject.RoleID,
) error {
	// Retrieve the admin
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return errors.New("admin not found")
	}

	// Verify that the new role exists
	role, err := s.roleRepo.FindByID(ctx, newRoleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Store the old role for event publishing
	oldRole := admin.Role

	// Update the admin's role
	admin.UpdateRole(newRoleID)

	// Persist the changes
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return err
	}

	// Publish the admin role changed event
	roleChangedEvent := entity.NewAdminRoleChangedEvent(admin, oldRole)
	s.eventBus.Publish(roleChangedEvent)

	return nil
}

// UpdateAdminStatus updates an admin's status
func (s *AdminDomainService) UpdateAdminStatus(
	ctx context.Context,
	adminID valueobject.AdminID,
	newStatus valueobject.AdminStatus,
) error {
	// Retrieve the admin
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return errors.New("admin not found")
	}

	// Store the old status for event publishing
	oldStatus := admin.Status

	// Update the admin's status
	if err := admin.UpdateStatus(newStatus); err != nil {
		return err
	}

	// Persist the changes
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return err
	}

	// Publish the admin status changed event
	statusChangedEvent := entity.NewAdminStatusChangedEvent(admin, oldStatus)
	s.eventBus.Publish(statusChangedEvent)

	return nil
}

// AuthenticateAdmin authenticates an admin using username and password
func (s *AdminDomainService) AuthenticateAdmin(
	ctx context.Context,
	username valueobject.Username,
	password string,
	ip, userAgent string,
) (*entity.Admin, error) {
	// Retrieve the admin by username
	admin, err := s.adminRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("admin not found")
	}

	// Check if the admin is active
	if !admin.IsActive() {
		return nil, errors.New("admin account is not active")
	}

	// Verify password
	if !admin.VerifyPassword(password) {
		return nil, errors.New("invalid password")
	}

	// Record the login
	admin.RecordLogin()
	if err := s.adminRepo.UpdateLastLogin(ctx, admin.ID); err != nil {
		return nil, err
	}

	// Publish login event
	loginEvent := entity.NewAdminLoggedInEvent(admin, ip, userAgent)
	s.eventBus.Publish(loginEvent)

	return admin, nil
}

// ChangeAdminPassword changes an admin's password
func (s *AdminDomainService) ChangeAdminPassword(
	ctx context.Context,
	adminID valueobject.AdminID,
	currentPassword, newPassword string,
) error {
	// Retrieve the admin
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return errors.New("admin not found")
	}

	// Change password
	if err := admin.ChangePassword(currentPassword, newPassword); err != nil {
		return err
	}

	// Persist the changes
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return err
	}

	// Publish password changed event
	passwordChangedEvent := entity.NewAdminPasswordChangedEvent(admin)
	s.eventBus.Publish(passwordChangedEvent)

	return nil
}

// DeleteAdmin deletes an admin
func (s *AdminDomainService) DeleteAdmin(
	ctx context.Context,
	adminID valueobject.AdminID,
) error {
	// Verify that the admin exists
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return errors.New("admin not found")
	}

	// Delete the admin
	if err := s.adminRepo.Delete(ctx, adminID); err != nil {
		return err
	}

	return nil
}

// GetAdminWithRoles retrieves an admin with their roles
func (s *AdminDomainService) GetAdminWithRoles(
	ctx context.Context,
	adminID valueobject.AdminID,
) (*entity.Admin, []*entity.Role, error) {
	// Retrieve the admin
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, nil, err
	}
	if admin == nil {
		return nil, nil, errors.New("admin not found")
	}

	// Retrieve the roles for this admin
	roles, err := s.roleRepo.FindByAdminID(ctx, adminID)
	if err != nil {
		return nil, nil, err
	}

	return admin, roles, nil
}

// GetAdminWithPermissions retrieves an admin with their permissions
func (s *AdminDomainService) GetAdminWithPermissions(
	ctx context.Context,
	adminID valueobject.AdminID,
) (*entity.Admin, []valueobject.Permission, error) {
	// Retrieve the admin
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, nil, err
	}
	if admin == nil {
		return nil, nil, errors.New("admin not found")
	}

	// Retrieve the admin's role
	role, err := s.roleRepo.FindByID(ctx, admin.Role)
	if err != nil {
		return nil, nil, err
	}
	if role == nil {
		return nil, nil, errors.New("role not found for admin")
	}

	// Retrieve permissions for the role
	permissions, err := s.roleRepo.GetPermissions(ctx, role.ID)
	if err != nil {
		return nil, nil, err
	}

	return admin, permissions, nil
}

// HasPermission checks if an admin has a specific permission
func (s *AdminDomainService) HasPermission(
	ctx context.Context,
	adminID valueobject.AdminID,
	permission valueobject.Permission,
) (bool, error) {
	// Retrieve the admin with permissions
	_, permissions, err := s.GetAdminWithPermissions(ctx, adminID)
	if err != nil {
		return false, err
	}

	// Check if the admin has the required permission
	for _, p := range permissions {
		if p.Equals(permission) {
			return true, nil
		}
		
		// Check for wildcard permissions
		if p.GetResourceType() == permission.GetResourceType() && 
		   (p.GetActionType() == "*" || p.GetActionType() == permission.GetActionType()) {
			return true, nil
		}
	}

	return false, nil
}
