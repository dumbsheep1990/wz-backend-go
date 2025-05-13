package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	"net/http"
	"wz-backend-go/internal/config"
	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminServiceContext 后台管理服务上下文
type AdminServiceContext struct {
	Config config.AdminConfig

	// 数据库连接
	SqlConn sqlx.SqlConn

	// Redis客户端
	RedisClient *redis.Redis

	// 数据仓库
	UserRepo         repository.UserRepository
	TenantRepo       repository.TenantRepository
	ContentRepo      repository.ContentRepository
	TradeRepo        repository.TradeRepository
	SettingsRepo     repository.SettingsRepository
	UserPointsRepo   repository.UserPointsRepository
	UserFavoriteRepo repository.UserFavoriteRepository

	// 服务
	UserPointsService   *service.UserPointsService
	UserFavoriteService *service.UserFavoriteService

	// RPC客户端
	UserClient    UserRpcClient    // 用户服务客户端
	ContentClient ContentRpcClient // 内容服务客户端
	TradeClient   TradeRpcClient   // 交易服务客户端
	SearchClient  SearchRpcClient  // 搜索服务客户端
}

// NewAdminServiceContext 创建后台管理服务上下文
func NewAdminServiceContext(c config.AdminConfig) *AdminServiceContext {
	sqlConn := sqlx.NewMysql(c.Database.DataSource)

	rc := redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = "node"
		r.Pass = c.Redis.Pass
		r.DB = c.Redis.DB
	})

	// 创建仓储实例
	userRepo := repository.NewSqlUserRepository(sqlConn)
	tenantRepo := repository.NewSqlTenantRepository(sqlConn)
	contentRepo := repository.NewSqlContentRepository(sqlConn)
	tradeRepo := repository.NewSqlTradeRepository(sqlConn)
	settingsRepo := repository.NewSqlSettingsRepository(sqlConn)
	userPointsRepo := repository.NewSQLUserPointsRepository(sqlConn)
	userFavoriteRepo := repository.NewSQLUserFavoriteRepository(sqlConn)

	// 创建服务实例
	userPointsService := service.NewUserPointsService(userPointsRepo, userRepo)
	userFavoriteService := service.NewUserFavoriteService(userFavoriteRepo, userRepo, contentRepo)

	return &AdminServiceContext{
		Config:      c,
		SqlConn:     sqlConn,
		RedisClient: rc,

		// 数据仓库
		UserRepo:         userRepo,
		TenantRepo:       tenantRepo,
		ContentRepo:      contentRepo,
		TradeRepo:        tradeRepo,
		SettingsRepo:     settingsRepo,
		UserPointsRepo:   userPointsRepo,
		UserFavoriteRepo: userFavoriteRepo,

		// 服务
		UserPointsService:   userPointsService,
		UserFavoriteService: userFavoriteService,

		// 初始化RPC客户端
		UserClient:    NewUserRpcClient(zrpc.MustNewClient(c.UserRpc)),
		ContentClient: NewContentRpcClient(zrpc.MustNewClient(c.ContentRpc)),
		TradeClient:   NewTradeRpcClient(zrpc.MustNewClient(c.TradeRpc)),
		SearchClient:  NewSearchRpcClient(zrpc.MustNewClient(c.SearchRpc)),
	}
}

// WrapHandler 封装Gin处理函数
func (ctx *AdminServiceContext) WrapHandler(handler func(*gin.Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ginCtx := &gin.Context{
			Request: r,
			Writer:  w,
		}

		handler(ginCtx)
	}
}

// JwtAuth JWT认证中间件
func (ctx *AdminServiceContext) JwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 验证JWT令牌
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "未授权", http.StatusUnauthorized)
			return
		}

		// TODO: 实现JWT验证逻辑

		// 验证通过，继续处理请求
		next(w, r)
	}
}

// 注意：下面是RPC客户端接口的定义
// 实际项目中，这些接口会由protoc自动生成，或者你需要按照项目的实际RPC定义来实现它们

type UserRpcClient interface {
	// 在这里定义用户服务的RPC方法
	// 例如：GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error)
}

type ContentRpcClient interface {
	// 在这里定义内容服务的RPC方法
}

type TradeRpcClient interface {
	// 在这里定义交易服务的RPC方法
}

type SearchRpcClient interface {
	// 在这里定义搜索服务的RPC方法
}
