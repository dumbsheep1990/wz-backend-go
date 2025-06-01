package factory

import (
	"context"
	"log"
	"time"
	"github.com/gin-gonic/gin"
	appService "wz-backend-go/internal/application/render/service"
	domainService "wz-backend-go/internal/domain/render/service"
	"wz-backend-go/internal/infrastructure/event"
	"wz-backend-go/internal/infrastructure/persistence"
	"wz-backend-go/internal/interfaces/http/controller"
)

// RenderServiceFactory 渲染服务工厂
type RenderServiceFactory struct {
	renderResultRepo *persistence.RenderResultRepositoryImpl
	templateRepo     *persistence.TemplateRepositoryImpl
	eventBus         *event.EventBus
	
	// 依赖的其他服务，实际应用中需要注入
	siteService      SiteServiceClient
	pageService      PageServiceClient
}

// SiteServiceClient 站点服务客户端接口
type SiteServiceClient interface {
	appService.SiteService
}

// PageServiceClient 页面服务客户端接口
type PageServiceClient interface {
	appService.PageService
}

// NewRenderServiceFactory 创建渲染服务工厂
func NewRenderServiceFactory(
	siteService SiteServiceClient,
	pageService PageServiceClient,
) *RenderServiceFactory {
	return &RenderServiceFactory{
		renderResultRepo: persistence.NewRenderResultRepository(),
		templateRepo:     persistence.NewTemplateRepository(),
		eventBus:         event.NewEventBus(),
		siteService:      siteService,
		pageService:      pageService,
	}
}

// CreateRenderDomainService 创建渲染领域服务
func (f *RenderServiceFactory) CreateRenderDomainService() *domainService.RenderDomainService {
	return domainService.NewRenderDomainService(
		f.renderResultRepo,
		f.templateRepo,
		f.eventBus,
	)
}

// CreateRenderApplicationService 创建渲染应用服务
func (f *RenderServiceFactory) CreateRenderApplicationService() *appService.RenderApplicationService {
	renderDomainService := f.CreateRenderDomainService()
	
	return appService.NewRenderApplicationService(
		renderDomainService,
		f.siteService,
		f.pageService,
	)
}

// CreateRenderController 创建渲染控制器
func (f *RenderServiceFactory) CreateRenderController() *controller.RenderController {
	renderAppService := f.CreateRenderApplicationService()
	
	return controller.NewRenderController(renderAppService)
}

// RegisterRoutes 注册路由
func (f *RenderServiceFactory) RegisterRoutes(router *gin.Engine) {
	renderController := f.CreateRenderController()
	
	// 公共渲染API
	router.GET("/render/site", renderController.RenderSiteByDomain)
	router.GET("/render/sites/:siteId/:slug", renderController.RenderPageBySlug)
	
	// 内部API
	api := router.Group("/api/v1")
	{
		// 模板渲染
		api.POST("/render/template", renderController.RenderTemplate)
		
		// 预览API
		api.GET("/preview/sites/:siteId", renderController.PreviewSite)
		api.GET("/preview/sites/:siteId/pages/:pageId", renderController.PreviewPage)
		
		// 模板管理API
		templates := api.Group("/templates")
		{
			templates.POST("", renderController.CreateTemplate)
			templates.GET("/:id", renderController.GetTemplate)
			templates.PUT("/:id", renderController.UpdateTemplate)
			templates.DELETE("/:id", renderController.DeleteTemplate)
		}
	}
}

// StartCacheCleaner 启动缓存清理任务
func (f *RenderServiceFactory) StartCacheCleaner(ctx context.Context) {
	// 在实际应用中，这里应该启动一个定时任务来清理过期的缓存
	// 例如使用cron包或定时器
	go func() {
		// 定期清理过期缓存
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				err := f.renderResultRepo.DeleteExpired(ctx)
				if err != nil {
					// 记录日志
					log.Printf("清理过期缓存失败: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
