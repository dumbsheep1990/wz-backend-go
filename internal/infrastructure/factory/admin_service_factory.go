package factory

import (
	"database/sql"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	adminAppService "wz-backend-go/internal/application/admin/service"
	adminDomainService "wz-backend-go/internal/domain/admin/service"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/infrastructure/persistence/admin"
	infraService "wz-backend-go/internal/infrastructure/service"
	"wz-backend-go/internal/interfaces/http/controller"
	"wz-backend-go/internal/interfaces/http/middleware"
)

// AdminServiceFactory 管理员服务工厂
type AdminServiceFactory struct {
	db       *sql.DB
	redis    *redis.Client
	eventBus event.EventBus
	enforcer *casbin.Enforcer
	
	// JWT配置
	jwtConfig infraService.JWTConfig
	
	// 仓储实现
	adminRepository *admin.AdminRepositoryImpl
	roleRepository  *admin.RoleRepositoryImpl
	
	// 服务实例
	adminDomainService  *adminDomainService.AdminDomainService
	roleDomainService   *adminDomainService.RoleDomainService
	jwtService          *infraService.JWTServiceImpl
	adminAppService     *adminAppService.AdminApplicationService
	roleAppService      *adminAppService.RoleApplicationService
	
	// 控制器
	adminController *controller.AdminController
	roleController  *controller.RoleController
	
	// 中间件
	authMiddleware *middleware.AuthMiddleware
}

// NewAdminServiceFactory 创建管理员服务工厂
func NewAdminServiceFactory(
	db *sql.DB,
	redis *redis.Client,
	eventBus event.EventBus,
	enforcer *casbin.Enforcer,
	jwtConfig infraService.JWTConfig,
) *AdminServiceFactory {
	factory := &AdminServiceFactory{
		db:       db,
		redis:    redis,
		eventBus: eventBus,
		enforcer: enforcer,
		jwtConfig: jwtConfig,
	}
	
	// 初始化各层组件
	factory.initRepositories()
	factory.initDomainServices()
	factory.initApplicationServices()
	factory.initControllers()
	factory.initMiddlewares()
	
	return factory
}

// initRepositories 初始化仓储层
func (f *AdminServiceFactory) initRepositories() {
	f.adminRepository = admin.NewAdminRepository(f.db)
	f.roleRepository = admin.NewRoleRepository(f.db)
}

// initDomainServices 初始化领域服务
func (f *AdminServiceFactory) initDomainServices() {
	f.adminDomainService = adminDomainService.NewAdminDomainService(
		f.adminRepository,
		f.roleRepository,
		f.eventBus,
	)
	
	f.roleDomainService = adminDomainService.NewRoleDomainService(
		f.roleRepository,
		f.adminRepository,
		f.eventBus,
	)
}

// initApplicationServices 初始化应用服务
func (f *AdminServiceFactory) initApplicationServices() {
	f.jwtService = infraService.NewJWTService(f.jwtConfig, f.redis)
	
	f.adminAppService = adminAppService.NewAdminApplicationService(
		f.adminDomainService,
		f.roleDomainService,
		f.jwtService,
		f.adminRepository,
		f.roleRepository,
	)
	
	f.roleAppService = adminAppService.NewRoleApplicationService(
		f.roleDomainService,
		f.roleRepository,
	)
}

// initControllers 初始化控制器
func (f *AdminServiceFactory) initControllers() {
	f.adminController = controller.NewAdminController(f.adminAppService)
	f.roleController = controller.NewRoleController(f.roleAppService)
}

// initMiddlewares 初始化中间件
func (f *AdminServiceFactory) initMiddlewares() {
	f.authMiddleware = middleware.NewAuthMiddleware(f.jwtService, f.adminAppService)
}

// RegisterRoutes 注册路由
func (f *AdminServiceFactory) RegisterRoutes(router *gin.Engine) {
	// 获取认证中间件
	authHandler := f.authMiddleware.AdminAuthHandler()
	
	// 注册管理员控制器路由
	adminGroup := router.Group("/api/v1/admin")
	{
		// 公开API - 不需要认证
		adminGroup.POST("/login", f.adminController.Login)
		
		// 需要认证的API
		authorized := adminGroup.Group("")
		authorized.Use(authHandler)
		{
			// 管理员个人信息
			authorized.GET("/profile", f.adminController.GetCurrentAdmin)
			authorized.PUT("/password", f.adminController.ChangePassword)
			
			// 超级管理员权限API
			superAdmin := authorized.Group("")
			superAdmin.Use(f.authMiddleware.PermissionMiddleware("admin:manage"))
			{
				superAdmin.POST("/admins", f.adminController.CreateAdmin)
				superAdmin.GET("/admins", f.adminController.ListAdmins)
				superAdmin.GET("/admins/:id", f.adminController.GetAdmin)
				superAdmin.PUT("/admins/:id", f.adminController.UpdateAdmin)
				superAdmin.DELETE("/admins/:id", f.adminController.DeleteAdmin)
			}
		}
	}
	
	// 注册角色控制器路由
	f.roleController.RegisterRoutes(router, authHandler)
	
	// 记录路由注册完成
	log.Println("管理员服务路由注册完成")
}

// CreateServiceContext 创建服务上下文
func (f *AdminServiceFactory) CreateServiceContext() *infraService.ServiceContext {
	return &infraService.ServiceContext{
		// 基础设施
		DB:       f.db,
		Redis:    f.redis,
		EventBus: f.eventBus,
		Enforcer: f.enforcer,
		
		// 领域仓储
		AdminRepository: *f.adminRepository,
		RoleRepository:  *f.roleRepository,
		
		// 领域服务
		AdminDomainService: *f.adminDomainService,
		RoleDomainService:  *f.roleDomainService,
		
		// 应用服务
		JWTService:      f.jwtService,
		AdminAppService: f.adminAppService,
		RoleAppService:  f.roleAppService,
	}
}

// SetupEventHandlers 设置事件处理器
func (f *AdminServiceFactory) SetupEventHandlers() {
	// 订阅管理员事件
	f.eventBus.Subscribe("AdminCreated", func(e event.DomainEvent) error {
		// 这里可以处理管理员创建事件，例如发送通知邮件等
		log.Printf("事件处理: 管理员创建 - %v", e)
		return nil
	})
	
	f.eventBus.Subscribe("AdminLoggedIn", func(e event.DomainEvent) error {
		// 处理管理员登录事件，例如记录登录日志等
		log.Printf("事件处理: 管理员登录 - %v", e)
		return nil
	})
	
	f.eventBus.Subscribe("AdminRoleChanged", func(e event.DomainEvent) error {
		// 处理管理员角色变更事件
		log.Printf("事件处理: 管理员角色变更 - %v", e)
		return nil
	})
	
	f.eventBus.Subscribe("AdminStatusChanged", func(e event.DomainEvent) error {
		// 处理管理员状态变更事件
		log.Printf("事件处理: 管理员状态变更 - %v", e)
		return nil
	})
	
	// 订阅角色事件
	f.eventBus.Subscribe("RoleCreated", func(e event.DomainEvent) error {
		// 处理角色创建事件
		log.Printf("事件处理: 角色创建 - %v", e)
		return nil
	})
	
	f.eventBus.Subscribe("RoleUpdated", func(e event.DomainEvent) error {
		// 处理角色更新事件
		log.Printf("事件处理: 角色更新 - %v", e)
		return nil
	})
	
	f.eventBus.Subscribe("RolePermissionsChanged", func(e event.DomainEvent) error {
		// 处理角色权限变更事件
		log.Printf("事件处理: 角色权限变更 - %v", e)
		return nil
	})
	
	log.Println("管理员服务事件处理器设置完成")
}
