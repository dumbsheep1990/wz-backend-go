package bootstrap

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	productAppService "wz-backend-go/internal/application/product/service"
	userAppService "wz-backend-go/internal/application/user/service"
	productRepo "wz-backend-go/internal/domain/product/repository"
	productDomainService "wz-backend-go/internal/domain/product/service"
	userRepo "wz-backend-go/internal/domain/user/repository"
	userDomainService "wz-backend-go/internal/domain/user/service"
	"wz-backend-go/internal/infrastructure/eventbus"
	"wz-backend-go/internal/infrastructure/persistence/product"
	"wz-backend-go/internal/infrastructure/persistence/sql"
	"wz-backend-go/internal/infrastructure/transaction"
	"wz-backend-go/internal/interfaces/http/handler"
	productHandler "wz-backend-go/internal/interfaces/http/product"
)

// App 应用程序
type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

// NewApp 创建应用程序
func NewApp(db *gorm.DB) *App {
	router := gin.Default()

	app := &App{
		Router: router,
		DB:     db,
	}

	app.setupRoutes()

	return app
}

// setupRoutes 设置路由
func (a *App) setupRoutes() {
	// 创建日志事件发布器
	eventPublisher := eventbus.NewLogEventPublisher()

	// 创建工作单元
	unitOfWork := transaction.NewGormUnitOfWork(a.DB)

	// 注册用户模块
	a.registerUserModule(unitOfWork.DB(), eventPublisher)

	// 注册商品模块
	a.registerProductModule(unitOfWork.DB(), eventPublisher)

	// 注册其他模块...
}

// registerUserModule 注册用户模块
func (a *App) registerUserModule(db *gorm.DB, eventPublisher userRepo.EventPublisher) {
	// 创建用户仓储
	userRepository := sql.NewUserRepository(db)

	// 创建用户领域服务
	userDomainService := userDomainService.NewUserDomainService(userRepository, eventPublisher)

	// 创建用户应用服务
	userAppSvc := userAppService.NewUserApplicationService(userDomainService)

	// 创建用户HTTP处理器
	userHandler := handler.NewUserHandler(userAppSvc)

	// 注册用户路由
	userHandler.RegisterRoutes(a.Router)
}

// registerProductModule 注册商品模块
func (a *App) registerProductModule(db *gorm.DB, eventPublisher productRepo.EventPublisher) {
	// 创建商品仓储
	productRepository := product.NewProductRepository(db)

	// 创建商品领域服务
	productDomainSvc := productDomainService.NewProductDomainService(productRepository, eventPublisher)

	// 创建商品应用服务
	productAppSvc := productAppService.NewProductApplicationService(productDomainSvc)

	// 创建商品HTTP处理器
	productHdl := productHandler.NewProductHandler(productAppSvc)

	// 注册商品路由
	productHdl.RegisterRoutes(a.Router)
}

// Run 运行应用程序
func (a *App) Run(addr string) error {
	log.Printf("Server is running on %s", addr)
	return a.Router.Run(addr)
}
