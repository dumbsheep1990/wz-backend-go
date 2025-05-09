package svc

import (
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
	"wz-backend-go/internal/admin/config"
	"wz-backend-go/internal/middleware"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext Admin服务上下文，包含所有依赖
type ServiceContext struct {
	Config     config.Config
	AdminCheck rest.Middleware
	Enforcer   *casbin.Enforcer // 权限管理
	Redis      *redis.Redis     // Redis客户端
	DB         sqlx.SqlConn     // 数据库连接

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
	rds := redis.New(c.Cache[0].Host)

	// 权限管理器
	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", "configs/rbac_policy.csv")
	if err != nil {
		panic(err)
	}

	// 创建各微服务的客户端
	userRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.UserServiceRPC,
	})
	contentRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.ContentServiceRPC,
	})
	interactionRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.InteractionServiceRPC,
	})
	aiRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.AIServiceRPC,
	})
	notificationRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.NotificationServiceRPC,
	})
	tradeRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.TradeServiceRPC,
	})
	fileRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.FileServiceRPC,
	})
	statisticsRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.StatisticsServiceRPC,
	})
	adRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.AdServiceRPC,
	})
	recommendRPC := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.RecommendServiceRPC,
	})

	return &ServiceContext{
		Config:     c,
		AdminCheck: middleware.NewAdminCheckMiddleware(enforcer).Handle,
		Enforcer:   enforcer,
		Redis:      rds,
		DB:         db,

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
