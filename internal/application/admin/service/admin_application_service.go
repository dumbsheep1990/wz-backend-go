package service

import (
	"context"
	"errors"
	"time"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/repository"
	domainService "wz-backend-go/internal/domain/admin/service"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// AdminApplicationService coordinates domain services for admin-related operations
type AdminApplicationService struct {
	adminDomainService *domainService.AdminDomainService
	roleDomainService  *domainService.RoleDomainService
	adminRepo          repository.AdminRepository
	roleRepo           repository.RoleRepository
	jwtService         JWTService
}

// NewAdminApplicationService creates a new AdminApplicationService
func NewAdminApplicationService(
	adminDomainService *domainService.AdminDomainService,
	roleDomainService *domainService.RoleDomainService,
	adminRepo repository.AdminRepository,
	roleRepo repository.RoleRepository,
	jwtService JWTService,
) *AdminApplicationService {
	return &AdminApplicationService{
		adminDomainService: adminDomainService,
		roleDomainService:  roleDomainService,
		adminRepo:          adminRepo,
		roleRepo:           roleRepo,
		jwtService:         jwtService,
	}
}

// CreateAdmin creates a new admin
func (s *AdminApplicationService) CreateAdmin(
	ctx context.Context,
	req dto.AdminCreateRequest,
) (*dto.AdminDTO, error) {
	// Convert request to domain objects
	username, err := valueobject.NewUsername(req.Username)
	if err != nil {
		return nil, err
	}

	roleID, err := valueobject.NewRoleID(req.RoleID)
	if err != nil {
		return nil, err
	}

	// Use domain service to create the admin
	admin, err := s.adminDomainService.CreateAdmin(ctx, username, req.Password, roleID)
	if err != nil {
		return nil, err
	}

	// Convert domain object to DTO
	adminDTO := dto.MapAdminToDTO(admin)
	return &adminDTO, nil
}

// UpdateAdmin updates an existing admin
func (s *AdminApplicationService) UpdateAdmin(
	ctx context.Context,
	id int64,
	req dto.AdminUpdateRequest,
) (*dto.AdminDTO, error) {
	adminID, err := valueobject.NewAdminID(id)
	if err != nil {
		return nil, err
	}

	// Retrieve the admin
	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("admin not found")
	}

	// Update username if provided
	if req.Username != "" {
		username, err := valueobject.NewUsername(req.Username)
		if err != nil {
			return nil, err
		}
		admin.UpdateUsername(username)
	}

	// Update role if provided
	if req.RoleID != "" {
		roleID, err := valueobject.NewRoleID(req.RoleID)
		if err != nil {
			return nil, err
		}
		if err := s.adminDomainService.UpdateAdminRole(ctx, adminID, roleID); err != nil {
			return nil, err
		}
	}

	// Update status if provided
	if req.Status != nil {
		status, err := valueobject.NewAdminStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		if err := s.adminDomainService.UpdateAdminStatus(ctx, adminID, status); err != nil {
			return nil, err
		}
	}

	// Persist changes if username was updated
	if req.Username != "" {
		if err := s.adminRepo.Update(ctx, admin); err != nil {
			return nil, err
		}
	}

	// Retrieve the updated admin
	updatedAdmin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, err
	}

	// Convert domain object to DTO
	adminDTO := dto.MapAdminToDTO(updatedAdmin)
	return &adminDTO, nil
}

// DeleteAdmin deletes an admin
func (s *AdminApplicationService) DeleteAdmin(ctx context.Context, id int64) error {
	adminID, err := valueobject.NewAdminID(id)
	if err != nil {
		return err
	}

	return s.adminDomainService.DeleteAdmin(ctx, adminID)
}

// GetAdmin retrieves an admin by ID
func (s *AdminApplicationService) GetAdmin(ctx context.Context, id int64) (*dto.AdminDTO, error) {
	adminID, err := valueobject.NewAdminID(id)
	if err != nil {
		return nil, err
	}

	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("admin not found")
	}

	adminDTO := dto.MapAdminToDTO(admin)
	return &adminDTO, nil
}

// GetAdminDetail retrieves detailed information about an admin including roles and permissions
func (s *AdminApplicationService) GetAdminDetail(ctx context.Context, id int64) (*dto.AdminDetailResponse, error) {
	adminID, err := valueobject.NewAdminID(id)
	if err != nil {
		return nil, err
	}

	// Retrieve admin with roles
	admin, roles, err := s.adminDomainService.GetAdminWithRoles(ctx, adminID)
	if err != nil {
		return nil, err
	}

	// Retrieve permissions
	_, permissions, err := s.adminDomainService.GetAdminWithPermissions(ctx, adminID)
	if err != nil {
		return nil, err
	}

	// Convert domain objects to DTOs
	adminDTO := dto.MapAdminToDTO(admin)
	
	roleDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dto.MapRoleToDTO(role))
	}
	
	permissionStrings := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		permissionStrings = append(permissionStrings, permission.String())
	}

	// Create the response
	response := &dto.AdminDetailResponse{
		Admin:       adminDTO,
		Roles:       roleDTOs,
		Permissions: permissionStrings,
	}

	return response, nil
}

// ListAdmins retrieves a paginated list of admins
func (s *AdminApplicationService) ListAdmins(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) (*dto.AdminListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Use repository to retrieve admins
	admins, total, err := s.adminRepo.List(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	// Convert domain objects to DTOs
	adminDTOs := make([]dto.AdminDTO, 0, len(admins))
	for _, admin := range admins {
		adminDTO := dto.MapAdminToDTO(admin)
		
		// Optionally fetch role name if needed
		if admin.Role.Value() != "" {
			role, err := s.roleRepo.FindByID(ctx, admin.Role)
			if err == nil && role != nil {
				adminDTO.RoleName = role.Name.Value()
			}
		}
		
		adminDTOs = append(adminDTOs, adminDTO)
	}

	// Create the response
	response := &dto.AdminListResponse{
		Total: total,
		Items: adminDTOs,
	}

	return response, nil
}

// ChangePassword changes an admin's password
func (s *AdminApplicationService) ChangePassword(
	ctx context.Context,
	id int64,
	req dto.AdminPasswordChangeRequest,
) error {
	adminID, err := valueobject.NewAdminID(id)
	if err != nil {
		return err
	}

	return s.adminDomainService.ChangeAdminPassword(
		ctx,
		adminID,
		req.CurrentPassword,
		req.NewPassword,
	)
}

// Login authenticates an admin and returns a JWT token
func (s *AdminApplicationService) Login(
	ctx context.Context,
	req dto.AdminLoginRequest,
	ip, userAgent string,
) (*dto.AdminLoginResponse, error) {
	// Convert request to domain objects
	username, err := valueobject.NewUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// Authenticate the admin
	admin, err := s.adminDomainService.AuthenticateAdmin(ctx, username, req.Password, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtService.GenerateToken(admin.ID.Value(), admin.Username.Value())
	if err != nil {
		return nil, err
	}

	// Retrieve admin roles
	_, roles, err := s.adminDomainService.GetAdminWithRoles(ctx, admin.ID)
	if err != nil {
		return nil, err
	}

	// Convert domain objects to DTOs
	adminDTO := dto.MapAdminToDTO(admin)
	
	roleDTOs := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dto.MapRoleToDTO(role))
	}

	// Create response with site info and categories
	response := &dto.AdminLoginResponse{
		Token:      token,
		ExpiresAt:  expiresAt,
		Admin:      &adminDTO,
		Roles:      roleDTOs,
		Categories: dto.GetDefaultCategories(),
		SiteInfo:   dto.NewDefaultSiteInfo(),
	}

	return response, nil
}

// GetCurrentAdmin retrieves the currently authenticated admin
func (s *AdminApplicationService) GetCurrentAdmin(ctx context.Context, adminID int64) (*dto.AdminDetailResponse, error) {
	return s.GetAdminDetail(ctx, adminID)
}

// HasPermission checks if an admin has a specific permission
func (s *AdminApplicationService) HasPermission(
	ctx context.Context,
	adminID int64,
	resource, action string,
) (bool, error) {
	id, err := valueobject.NewAdminID(adminID)
	if err != nil {
		return false, err
	}

	permission, err := valueobject.NewPermission(resource, action)
	if err != nil {
		return false, err
	}

	return s.adminDomainService.HasPermission(ctx, id, permission)
}
