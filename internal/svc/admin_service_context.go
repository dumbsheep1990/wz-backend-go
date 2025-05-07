package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	"wz-backend-go/internal/config"
	"wz-backend-go/internal/repository"
)

// AdminServiceContext 后台管理服务上下文
type AdminServiceContext struct {
	Config config.AdminConfig

	// 数据库连接
	SqlConn sqlx.SqlConn

	// Redis客户端
	RedisClient *redis.Redis

	// 数据仓库
	UserRepo     repository.UserRepository
	TenantRepo   repository.TenantRepository
	ContentRepo  repository.ContentRepository
	TradeRepo    repository.TradeRepository
	SettingsRepo repository.SettingsRepository

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

	return &AdminServiceContext{
		Config:      c,
		SqlConn:     sqlConn,
		RedisClient: rc,

		// 初始化数据仓库
		UserRepo:     repository.NewSqlUserRepository(sqlConn),
		TenantRepo:   repository.NewSqlTenantRepository(sqlConn),
		ContentRepo:  repository.NewSqlContentRepository(sqlConn),
		TradeRepo:    repository.NewSqlTradeRepository(sqlConn),
		SettingsRepo: repository.NewSqlSettingsRepository(sqlConn),

		// 初始化RPC客户端
		UserClient:    NewUserRpcClient(zrpc.MustNewClient(c.UserRpc)),
		ContentClient: NewContentRpcClient(zrpc.MustNewClient(c.ContentRpc)),
		TradeClient:   NewTradeRpcClient(zrpc.MustNewClient(c.TradeRpc)),
		SearchClient:  NewSearchRpcClient(zrpc.MustNewClient(c.SearchRpc)),
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
