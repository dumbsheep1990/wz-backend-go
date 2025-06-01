package service

import (
	"database/sql"

	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"

	adminAppService "wz-backend-go/internal/application/admin/service"
	adminDomainService "wz-backend-go/internal/domain/admin/service"
	"wz-backend-go/internal/domain/shared/event"
	adminRepository "wz-backend-go/internal/infrastructure/persistence/admin"
)

// ServiceContext contains all services and repositories needed by the application
type ServiceContext struct {
	// Infrastructure
	DB       *sql.DB
	Redis    *redis.Client
	EventBus event.EventBus
	Enforcer *casbin.Enforcer

	// Domain Repositories
	AdminRepository adminRepository.AdminRepositoryImpl
	RoleRepository  adminRepository.RoleRepositoryImpl

	// Domain Services
	AdminDomainService adminDomainService.AdminDomainService
	RoleDomainService  adminDomainService.RoleDomainService

	// Application Services
	JWTService         *JWTServiceImpl
	AdminAppService    *adminAppService.AdminApplicationService
	RoleAppService     *adminAppService.RoleApplicationService
}

// NewServiceContext creates a new ServiceContext
func NewServiceContext(
	db *sql.DB,
	redis *redis.Client,
	eventBus event.EventBus,
	enforcer *casbin.Enforcer,
	jwtConfig JWTConfig,
) *ServiceContext {
	// Create repositories
	adminRepo := adminRepository.NewAdminRepository(db)
	roleRepo := adminRepository.NewRoleRepository(db)

	// Create domain services
	adminDomainSvc := adminDomainService.NewAdminDomainService(adminRepo, roleRepo, eventBus)
	roleDomainSvc := adminDomainService.NewRoleDomainService(roleRepo, adminRepo, eventBus)

	// Create JWT service
	jwtSvc := NewJWTService(jwtConfig, redis)

	// Create application services
	adminAppSvc := adminAppService.NewAdminApplicationService(
		adminDomainSvc,
		roleDomainSvc,
		jwtSvc,
		adminRepo,
		roleRepo,
	)

	roleAppSvc := adminAppService.NewRoleApplicationService(
		roleDomainSvc,
		roleRepo,
	)

	return &ServiceContext{
		// Infrastructure
		DB:       db,
		Redis:    redis,
		EventBus: eventBus,
		Enforcer: enforcer,

		// Domain Repositories
		AdminRepository: *adminRepo,
		RoleRepository:  *roleRepo,

		// Domain Services
		AdminDomainService: *adminDomainSvc,
		RoleDomainService:  *roleDomainSvc,

		// Application Services
		JWTService:      jwtSvc,
		AdminAppService: adminAppSvc,
		RoleAppService:  roleAppSvc,
	}
}
