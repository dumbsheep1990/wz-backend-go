ffpackage service

import (
	"database/sql"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
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
	Redis               *redis.Client     // Redis客户端
	DB                  sqlx.SqlConn     // 数据库连接
	RestServer          rest.Server       // REST服务器

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

	// 新增的存储库
	AuthorityRepository   repository.AuthorityRepository
	MenuRepository        repository.MenuRepository
	ApiRepository         repository.ApiRepository
	DictionaryRepository  repository.DictionaryRepository
	OperationRepository   repository.OperationRepository
	ParamsRepository      repository.ParamsRepository
	SystemConfigRepository repository.SystemConfigRepository
	JwtRepository         repository.JwtRepository
}

// NewServiceContext 创建新的服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 数据库连接
	db := sqlx.NewMysql(c.DB.DataSource)

	// Redis客户端
	rds := redis.NewClient(&redis.Options{
		Addr:     c.Cache[0].Host,
		Password: c.Cache[0].Pass,
		DB:       0,
	})

	// 初始化Casbin
	enforcer, err := initCasbin(c.DB.DataSource)
	if err != nil {
		log.Fatalf("初始化Casbin失败: %v", err)
	}

	// 创建REST服务器
	restServer := rest.MustNewServer(c.RestConf)

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

	// 创建中间件
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(enforcer, rds).Handle

	// 创建服务上下文
	ctx := &ServiceContext{
		Config:              c,
		AdminAuthMiddleware: adminAuthMiddleware,
		Enforcer:            enforcer,
		Redis:               rds,
		DB:                  db,
		RestServer:          restServer,

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

		// 初始化新增的存储库
		AuthorityRepository:   repository.NewAuthorityRepository(db),
		MenuRepository:        repository.NewMenuRepository(db),
		ApiRepository:         repository.NewApiRepository(db),
		DictionaryRepository:  repository.NewDictionaryRepository(db),
		OperationRepository:   repository.NewOperationRepository(db),
		ParamsRepository:      repository.NewParamsRepository(db),
		SystemConfigRepository: repository.NewSystemConfigRepository(db),
		JwtRepository:         repository.NewJwtRepository(db),
	}

	// 初始化存储库
	ctx.initRepositories()

	return ctx
}

// 初始化Casbin
func initCasbin(dsn string) (*casbin.Enforcer, error) {
	// 定义模型
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		return nil, err
	}
	
	// 使用GORM适配器
	adapter, err := gormadapter.NewAdapter("mysql", dsn, true)
	if err != nil {
		return nil, err
	}
	
	// 创建执行器
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	
	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	
	return enforcer, nil
}

// 初始化存储库
func (s *ServiceContext) initRepositories() {
	s.UserRepo = repository.NewUserRepository(s.DB)
	s.TenantRepo = repository.NewTenantRepository(s.DB)
	s.ContentRepo = repository.NewContentRepository(s.DB)
	s.TradeRepo = repository.NewTradeRepository(s.DB)
	s.SettingsRepo = repository.NewSettingsRepository(s.DB)
	s.AdminRepo = repository.NewAdminRepository(s.DB)
	s.RoleRepo = repository.NewRoleRepository(s.DB)
	s.OperationLogRepo = repository.NewOperationLogRepository(s.DB)
	s.AuthorityRepository = repository.NewAuthorityRepository(s.DB)
	s.MenuRepository = repository.NewMenuRepository(s.DB)
	s.ApiRepository = repository.NewApiRepository(s.DB)
	s.DictionaryRepository = repository.NewDictionaryRepository(s.DB)
	s.OperationRepository = repository.NewOperationRepository(s.DB)
	s.ParamsRepository = repository.NewParamsRepository(s.DB)
	s.SystemConfigRepository = repository.NewSystemConfigRepository(s.DB)
	s.JwtRepository = repository.NewJwtRepository(s.DB)
}
