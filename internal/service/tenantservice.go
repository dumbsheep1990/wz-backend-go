package service

import (
	"context"
	"strconv"

	"wz-backend-go/internal/domain/model"
	"wz-backend-go/internal/repository"
)

// TenantService 租户服务接口
type TenantService interface {
	// 创建新租户
	CreateTenant(ctx context.Context, req *model.CreateTenantRequest, creatorUserID int64) (*model.Tenant, error)
	// 获取租户详情
	GetTenantByID(ctx context.Context, id int64) (*model.Tenant, error)
	// 获取租户详情（字符串ID）
	GetTenantByIDString(ctx context.Context, id string) (*model.Tenant, error)
	// 通过子域名获取租户
	GetTenantBySubdomain(ctx context.Context, subdomain string) (*model.Tenant, error)
	// 更新租户信息
	UpdateTenant(ctx context.Context, req *model.UpdateTenantRequest) error
	// 列出所有活跃租户
	ListActiveTenants(ctx context.Context) ([]*model.Tenant, error)
	// 添加用户到租户
	AddUserToTenant(ctx context.Context, tenantID, userID int64, role string) error
	// 获取租户下的用户列表
	GetTenantUsers(ctx context.Context, tenantID int64) ([]*model.User, error)
	// 从租户中移除用户
	RemoveUserFromTenant(ctx context.Context, tenantID, userID int64) error
	// 检查用户是否属于租户
	CheckUserInTenant(ctx context.Context, tenantID, userID int64) (bool, string, error)
}

// tenantService 租户服务实现
type tenantService struct {
	tenantRepo repository.TenantRepository
}

// NewTenantService 创建租户服务
func NewTenantService(tenantRepo repository.TenantRepository) TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
	}
}

// CreateTenant 创建新租户
func (s *tenantService) CreateTenant(ctx context.Context, req *model.CreateTenantRequest, creatorUserID int64) (*model.Tenant, error) {
	tenant := &model.Tenant{
		Name:          req.Name,
		Subdomain:     req.Subdomain,
		TenantType:    req.TenantType,
		Description:   req.Description,
		Logo:          req.Logo,
		CreatorUserID: creatorUserID,
		Status:        1, // 默认状态为正常
	}
	
	return s.tenantRepo.CreateTenant(ctx, tenant)
}

// GetTenantByID 获取租户详情
func (s *tenantService) GetTenantByID(ctx context.Context, id int64) (*model.Tenant, error) {
	return s.tenantRepo.GetTenantByID(ctx, id)
}

// GetTenantByIDString 获取租户详情（字符串ID）
func (s *tenantService) GetTenantByIDString(ctx context.Context, id string) (*model.Tenant, error) {
	tenantID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.GetTenantByID(ctx, tenantID)
}

// GetTenantBySubdomain 通过子域名获取租户
func (s *tenantService) GetTenantBySubdomain(ctx context.Context, subdomain string) (*model.Tenant, error) {
	return s.tenantRepo.GetTenantBySubdomain(ctx, subdomain)
}

// UpdateTenant 更新租户信息
func (s *tenantService) UpdateTenant(ctx context.Context, req *model.UpdateTenantRequest) error {
	tenant, err := s.GetTenantByID(ctx, req.ID)
	if err != nil {
		return err
	}
	
	// 更新字段
	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Description != "" {
		tenant.Description = req.Description
	}
	if req.Logo != "" {
		tenant.Logo = req.Logo
	}
	if req.Status != 0 {
		tenant.Status = req.Status
	}
	
	return s.tenantRepo.UpdateTenant(ctx, tenant)
}

// ListActiveTenants 列出所有活跃租户
func (s *tenantService) ListActiveTenants(ctx context.Context) ([]*model.Tenant, error) {
	return s.tenantRepo.ListActiveTenants(ctx)
}

// AddUserToTenant 添加用户到租户
func (s *tenantService) AddUserToTenant(ctx context.Context, tenantID, userID int64, role string) error {
	tenantUser := &model.TenantUser{
		TenantID:  tenantID,
		UserID:    userID,
		Role:      role,
		Status:    1, // 默认状态为正常
	}
	
	return s.tenantRepo.AddUserToTenant(ctx, tenantUser)
}

// GetTenantUsers 获取租户下的用户列表
func (s *tenantService) GetTenantUsers(ctx context.Context, tenantID int64) ([]*model.User, error) {
	return s.tenantRepo.GetTenantUsers(ctx, tenantID)
}

// RemoveUserFromTenant 从租户中移除用户
func (s *tenantService) RemoveUserFromTenant(ctx context.Context, tenantID, userID int64) error {
	return s.tenantRepo.RemoveUserFromTenant(ctx, tenantID, userID)
}

// CheckUserInTenant 检查用户是否属于租户
func (s *tenantService) CheckUserInTenant(ctx context.Context, tenantID, userID int64) (bool, string, error) {
	return s.tenantRepo.CheckUserInTenant(ctx, tenantID, userID)
}
