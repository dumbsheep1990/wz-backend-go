package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"wz-backend-go/api/rpc/ad"
	"wz-backend-go/api/rpc/ai"
	"wz-backend-go/api/rpc/content"
	"wz-backend-go/api/rpc/file"
	"wz-backend-go/api/rpc/interaction"
	"wz-backend-go/api/rpc/notification"
	"wz-backend-go/api/rpc/recommend"
	"wz-backend-go/api/rpc/statistics"
	"wz-backend-go/api/rpc/trade"
	"wz-backend-go/api/rpc/user"
	"wz-backend-go/services/admin-service/config"
	"wz-backend-go/services/admin-service/internal/middleware"
	"wz-backend-go/services/admin-service/internal/repository"
)

// ServiceContext Admin服务上下文，包含所有依赖
type ServiceContext struct {
	Config              config.Config
	AdminAuthMiddleware rest.Middleware
	Enforcer            *casbin.Enforcer // 权限管理
	Redis               *redis.Redis     // Redis客户端
	DB                  sqlx.SqlConn     // 数据库连接

	// 数据仓库
	UserRepo         repository.UserRepository
	TenantRepo       repository.TenantRepository
	ContentRepo      repository.ContentRepository
	TradeRepo        repository.TradeRepository
	SettingsRepo     repository.SettingsRepository
	AdminRepo        repository.AdminRepository
	RoleRepo         repository.RoleRepository
	OperationLogRepo repository.OperationLogRepository

	// 各微服务的客户端
	UserClient         user.UserService
	ContentClient      content.ContentService
	InteractionClient  interaction.InteractionService
	AIClient           ai.AIService
	NotificationClient notification.NotificationService
	TradeClient        trade.TradeService
	FileClient         file.FileService
	StatisticsClient   statistics.StatisticsService
	AdClient           ad.AdService               // 广告服务客户端
	RecommendClient    recommend.RecommendService // 推荐服务客户端
}

// NewServiceContext 创建新的服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 数据库连接
	db := sqlx.NewMysql(c.DB.DataSource)

	// Redis客户端
	rds := redis.New(c.Cache[0].Host, redis.WithPass(c.Cache[0].Pass))

	// 创建各微服务的客户端
	userRPC := zrpc.MustNewClient(c.UserRPC)
	contentRPC := zrpc.MustNewClient(c.ContentRPC)
	interactionRPC := zrpc.MustNewClient(c.InteractionRPC)
	aiRPC := zrpc.MustNewClient(c.AIRPC)
	notificationRPC := zrpc.MustNewClient(c.NotificationRPC)
	tradeRPC := zrpc.MustNewClient(c.TradeRPC)
	fileRPC := zrpc.MustNewClient(c.FileRPC)
	statisticsRPC := zrpc.MustNewClient(c.StatisticsRPC)
	adRPC := zrpc.MustNewClient(c.AdRPC)
	recommendRPC := zrpc.MustNewClient(c.RecommendRPC)

	// 创建仓库实例
	userRepo := repository.NewUserRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	contentRepo := repository.NewContentRepository(db)
	tradeRepo := repository.NewTradeRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	operationLogRepo := repository.NewOperationLogRepository(db)

	// 权限管理器
	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", "configs/rbac_policy.csv")
	if err != nil {
		panic(err)
	}

	// 创建中间件
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(enforcer, rds).Handle

	return &ServiceContext{
		Config:              c,
		AdminAuthMiddleware: adminAuthMiddleware,
		Enforcer:            enforcer,
		Redis:               rds,
		DB:                  db,

		// 注入仓库
		UserRepo:         userRepo,
		TenantRepo:       tenantRepo,
		ContentRepo:      contentRepo,
		TradeRepo:        tradeRepo,
		SettingsRepo:     settingsRepo,
		AdminRepo:        adminRepo,
		RoleRepo:         roleRepo,
		OperationLogRepo: operationLogRepo,

		// 注入微服务客户端
		UserClient:         user.NewUserService(userRPC.Conn()),
		ContentClient:      content.NewContentService(contentRPC.Conn()),
		InteractionClient:  interaction.NewInteractionService(interactionRPC.Conn()),
		AIClient:           ai.NewAIService(aiRPC.Conn()),
		NotificationClient: notification.NewNotificationService(notificationRPC.Conn()),
		TradeClient:        trade.NewTradeService(tradeRPC.Conn()),
		FileClient:         file.NewFileService(fileRPC.Conn()),
		StatisticsClient:   statistics.NewStatisticsService(statisticsRPC.Conn()),
		AdClient:           ad.NewAdService(adRPC.Conn()),
		RecommendClient:    recommend.NewRecommendService(recommendRPC.Conn()),
	}
}
