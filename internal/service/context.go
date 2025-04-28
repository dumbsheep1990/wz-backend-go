package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/casbin/casbin/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"wz-backend-go/internal/repository"
	"wz-backend-go/api/rpc/user"
	"wz-backend-go/api/rpc/content"
	"wz-backend-go/api/rpc/interaction"
	"wz-backend-go/api/rpc/notification"
	"wz-backend-go/api/rpc/ai"
	"wz-backend-go/api/rpc/file"
	"wz-backend-go/api/rpc/statistics"
)

// ServiceContext 服务上下文，包含所有服务依赖
type ServiceContext struct {
	// 数据库连接
	DB *gorm.DB
	
	// 缓存客户端
	Redis *redis.Client
	
	// 消息队列
	UserEventProducer    *kafka.Writer
	ContentEventProducer *kafka.Writer
	UserEventConsumer    *kafka.Reader
	ContentEventConsumer *kafka.Reader
	
	// 权限管理
	Enforcer *casbin.Enforcer
	
	// gRPC服务客户端
	UserClient         user.UserServiceClient
	ContentClient      content.ContentServiceClient
	InteractionClient  interaction.InteractionServiceClient
	NotificationClient notification.NotificationServiceClient
	AIClient           ai.AIServiceClient
	FileClient         file.FileServiceClient
	StatisticsClient   statistics.StatisticsServiceClient
	
	// 仓储层
	TenantRepo repository.TenantRepository
	UserRepo   repository.UserRepository
	// 其他仓储...
}

// NewServiceContext 创建服务上下文
func NewServiceContext(
	db *gorm.DB,
	redis *redis.Client,
	userEventProducer *kafka.Writer,
	contentEventProducer *kafka.Writer,
	userEventConsumer *kafka.Reader,
	contentEventConsumer *kafka.Reader,
	enforcer *casbin.Enforcer,
	userConn *grpc.ClientConn,
	contentConn *grpc.ClientConn,
	interactionConn *grpc.ClientConn,
	notificationConn *grpc.ClientConn,
	aiConn *grpc.ClientConn,
	fileConn *grpc.ClientConn,
	statisticsConn *grpc.ClientConn,
	tenantRepo repository.TenantRepository,
	userRepo repository.UserRepository,
) *ServiceContext {
	return &ServiceContext{
		DB:                  db,
		Redis:               redis,
		UserEventProducer:   userEventProducer,
		ContentEventProducer: contentEventProducer,
		UserEventConsumer:    userEventConsumer,
		ContentEventConsumer: contentEventConsumer,
		Enforcer:            enforcer,
		UserClient:          user.NewUserServiceClient(userConn),
		ContentClient:       content.NewContentServiceClient(contentConn),
		InteractionClient:   interaction.NewInteractionServiceClient(interactionConn),
		NotificationClient:  notification.NewNotificationServiceClient(notificationConn),
		AIClient:            ai.NewAIServiceClient(aiConn),
		FileClient:          file.NewFileServiceClient(fileConn),
		StatisticsClient:    statistics.NewStatisticsServiceClient(statisticsConn),
		TenantRepo:          tenantRepo,
		UserRepo:            userRepo,
	}
}
