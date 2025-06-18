package admin

import (
	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/admin"
	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/middleware"
	"wz-backend-go/internal/handler/admin/navigation"
	"wz-backend-go/internal/handler/admin/page"
	"wz-backend-go/internal/handler/admin/link"
	"wz-backend-go/internal/handler/admin/component"
	"wz-backend-go/internal/handler/admin/theme"
	"wz-backend-go/internal/client"
)

// RegisterRoutes 注册所有admin相关路由
func RegisterRoutes(
	r *gin.Engine,
	adminService *service.AdminApplicationService,
	roleService *service.RoleApplicationService,
	menuService *service.MenuApplicationService,
	systemService *service.SystemApplicationService,
	dictionaryService *service.DictionaryApplicationService,
	operationRecordService *service.OperationRecordApplicationService,
	paramsService *service.ParamsApplicationService,
	dashboardService *admin.DashboardService,
	navigationClient client.NavigationClient,
	jwtMiddleware *middleware.JWTMiddleware,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	// 创建处理程序
	loginHandler := NewLoginHandler(adminService)
	userHandler := NewUserHandler(adminService, roleService)
	authorityHandler := NewAuthorityHandler(roleService)
	menuHandler := NewMenuHandler(menuService)
	systemHandler := NewSystemHandler(systemService)
	dictionaryHandler := NewDictionaryHandler(dictionaryService)
	operationRecordHandler := NewOperationRecordHandler(operationRecordService)
	paramsHandler := NewParamsHandler(paramsService)
	dashboardHandler := NewDashboardHandler(dashboardService)
	navigationHandler := navigation.NewNavigationAdminHandler(navigationClient)

	// API路径前缀
	apiV1 := r.Group("/api/v1")

	// 公共接口（不需要认证）
	{
		// 登录相关
		apiV1.POST("/admin/login", loginHandler.Login)
		apiV1.POST("/admin/refresh_token", loginHandler.RefreshToken)
		apiV1.GET("/admin/captcha", loginHandler.Captcha)
	}

	// 需要JWT认证的接口
	adminGroup := apiV1.Group("/admin")
	adminGroup.Use(jwtMiddleware.Handler())
	{
		// 需要额外权限控制的接口
		authGroup := adminGroup.Group("")
		authGroup.Use(casbinMiddleware.Handler())
		{
			// 用户管理
			userRoutes := authGroup.Group("/user")
			{
				userRoutes.POST("/getUserList", userHandler.GetUserList)
				userRoutes.GET("/getUserInfo", userHandler.GetUserInfo)
				userRoutes.POST("/admin_register", userHandler.CreateUser)
				userRoutes.PUT("/setUserInfo", userHandler.UpdateUser)
				userRoutes.PUT("/setSelfInfo", userHandler.UpdateSelfInfo)
				userRoutes.DELETE("/deleteUser", userHandler.DeleteUser)
				userRoutes.POST("/changePassword", userHandler.ChangePassword)
				userRoutes.POST("/resetPassword", userHandler.ResetPassword)
				userRoutes.POST("/setUserAuthority", userHandler.SetUserAuthority)
				userRoutes.POST("/setUserAuthorities", userHandler.SetUserAuthorities)
			}

			// 角色权限管理
			authorityRoutes := authGroup.Group("/authority")
			{
				authorityRoutes.POST("/getAuthorityList", authorityHandler.GetAuthorityList)
				authorityRoutes.POST("/createAuthority", authorityHandler.CreateAuthority)
				authorityRoutes.PUT("/updateAuthority", authorityHandler.UpdateAuthority)
				authorityRoutes.POST("/deleteAuthority", authorityHandler.DeleteAuthority)
				authorityRoutes.POST("/setDataAuthority", authorityHandler.SetDataAuthority)
			}

			// 按钮权限管理
			authorityBtnRoutes := authGroup.Group("/authorityBtn")
			{
				// 在原API中有以下端点，但在应用服务中缺少实现
				// authorityBtnRoutes.POST("/getAuthorityBtn", authorityHandler.GetAuthorityBtnList)
				// authorityBtnRoutes.POST("/setAuthorityBtn", authorityHandler.SetAuthorityBtn)
			}

			// Casbin策略管理
			casbinRoutes := authGroup.Group("/casbin")
			{
				casbinRoutes.POST("/getPolicyPathByAuthorityId", authorityHandler.GetCasbinPolicy)
				casbinRoutes.POST("/updateCasbin", authorityHandler.UpdateCasbinPolicy)
			}

			// 菜单管理
			menuRoutes := authGroup.Group("/menu")
			{
				menuRoutes.POST("/getMenuList", menuHandler.GetMenuList)
				menuRoutes.POST("/addMenu", menuHandler.AddMenu)
				menuRoutes.POST("/updateMenu", menuHandler.UpdateMenu)
				menuRoutes.POST("/deleteMenu", menuHandler.DeleteMenu)
				menuRoutes.POST("/getMenuAuthority", menuHandler.GetMenusByAuthority)
				menuRoutes.POST("/addMenuAuthority", menuHandler.AddMenuAuthority)
			}

			// 导航管理
			navigationRoutes := authGroup.Group("/navigation")
			{
				// 主导航管理
				navigationRoutes.GET("/main", navigationHandler.GetMainNavigation)
				navigationRoutes.POST("/main", navigationHandler.SaveMainNavigation)
				navigationRoutes.DELETE("/main/:id", navigationHandler.DeleteMainNavigationItem)
				
				// 底部导航管理
				navigationRoutes.GET("/footer", navigationHandler.GetFooterNavigation)
				navigationRoutes.POST("/footer", navigationHandler.SaveFooterNavigation)
				navigationRoutes.DELETE("/footer/:id", navigationHandler.DeleteFooterNavigationItem)
				
				// 侧边导航管理
				navigationRoutes.GET("/side", navigationHandler.GetSideNavigation)
				navigationRoutes.POST("/side", navigationHandler.SaveSideNavigation)
				navigationRoutes.DELETE("/side/:id", navigationHandler.DeleteSideNavigationItem)
			}

			// 系统管理
			systemRoutes := authGroup.Group("/system")
			{
				systemRoutes.POST("/getSystemConfig", systemHandler.GetSystemConfig)
				systemRoutes.POST("/setSystemConfig", systemHandler.SetSystemConfig)
				systemRoutes.POST("/getServerInfo", systemHandler.GetServerInfo)
				systemRoutes.POST("/reloadSystem", systemHandler.ReloadSystem)
			}

			// 系统字典管理
			dictionaryRoutes := authGroup.Group("/sysDictionary")
			{
				dictionaryRoutes.POST("/getSysDictionaryList", dictionaryHandler.GetDictionaryList)
				dictionaryRoutes.POST("/createSysDictionary", dictionaryHandler.CreateDictionary)
				dictionaryRoutes.PUT("/updateSysDictionary", dictionaryHandler.UpdateDictionary)
				dictionaryRoutes.DELETE("/deleteSysDictionary", dictionaryHandler.DeleteDictionary)
			}

			// 系统字典详情管理
			dictionaryDetailRoutes := authGroup.Group("/sysDictionaryDetail")
			{
				dictionaryDetailRoutes.POST("/getSysDictionaryDetailList", dictionaryHandler.GetDictionaryDetailList)
				dictionaryDetailRoutes.POST("/createSysDictionaryDetail", dictionaryHandler.CreateDictionaryDetail)
				dictionaryDetailRoutes.PUT("/updateSysDictionaryDetail", dictionaryHandler.UpdateDictionaryDetail)
				dictionaryDetailRoutes.DELETE("/deleteSysDictionaryDetail", dictionaryHandler.DeleteDictionaryDetail)
			}

			// 操作日志管理
			operationRecordRoutes := authGroup.Group("/sysOperationRecord")
			{
				operationRecordRoutes.POST("/getSysOperationRecordList", operationRecordHandler.GetOperationRecordList)
				operationRecordRoutes.DELETE("/deleteSysOperationRecord", operationRecordHandler.DeleteOperationRecord)
				operationRecordRoutes.DELETE("/deleteSysOperationRecordByIds", operationRecordHandler.DeleteOperationRecordsByIds)
			}

			// 系统参数管理
			paramsRoutes := authGroup.Group("/params")
			{
				paramsRoutes.POST("/getParamsList", paramsHandler.GetParamsList)
				paramsRoutes.POST("/createParams", paramsHandler.CreateParams)
				paramsRoutes.PUT("/updateParams", paramsHandler.UpdateParams)
				paramsRoutes.DELETE("/deleteParams", paramsHandler.DeleteParams)
			}

			// 仪表盘数据
			dashboardRoutes := authGroup.Group("/dashboard")
			{
				dashboardRoutes.GET("/stats", dashboardHandler.GetDashboardStats)
				dashboardRoutes.GET("/userStats", dashboardHandler.GetUserStats)
				dashboardRoutes.GET("/contentStats", dashboardHandler.GetContentStats)
				dashboardRoutes.GET("/tradeStats", dashboardHandler.GetTradeStats)
				dashboardRoutes.GET("/siteStats", dashboardHandler.GetSiteStats)
				dashboardRoutes.GET("/communityStats", dashboardHandler.GetCommunityStats)
				dashboardRoutes.GET("/componentStats", dashboardHandler.GetComponentStats)
				dashboardRoutes.GET("/renderStats", dashboardHandler.GetRenderStats)
			}
		}
	}
}

// SetupSiteManagementRoutes 设置站点管理相关路由
func SetupSiteManagementRoutes(
	r *gin.Engine,
	pageClient client.PageClientInterface,
	linkClient client.LinkClientInterface,
	componentClient client.ComponentClientInterface,
	themeClient client.ThemeClientInterface,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	// 导入新的处理器包
	pageHandler := page.NewPageHandler(pageClient)
	linkHandler := link.NewLinkHandler(linkClient)
	componentHandler := component.NewComponentHandler(componentClient)
	themeHandler := theme.NewThemeHandler(themeClient)

	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(jwtMiddleware.Handler())

	// 页面管理路由
	pageGroup := adminGroup.Group("/pages")
	{
		pageGroup.GET("", pageHandler.GetPageList)
		pageGroup.POST("", pageHandler.CreatePage)
		pageGroup.GET("/:id", pageHandler.GetPageDetail)
		pageGroup.PUT("/:id", pageHandler.UpdatePage)
		pageGroup.DELETE("/:id", pageHandler.DeletePage)
		pageGroup.POST("/:id/toggle-status", pageHandler.TogglePageStatus)
		pageGroup.GET("/:id/preview", pageHandler.PreviewPage)
		pageGroup.POST("/batch-update", pageHandler.BatchUpdate)
	}

	// 链接管理路由
	linkGroup := adminGroup.Group("/links")
	{
		linkGroup.GET("", linkHandler.GetLinkList)
		linkGroup.POST("", linkHandler.CreateLink)
		linkGroup.GET("/:id", linkHandler.GetLinkDetail)
		linkGroup.PUT("/:id", linkHandler.UpdateLink)
		linkGroup.DELETE("/:id", linkHandler.DeleteLink)
		linkGroup.POST("/:id/verify", linkHandler.VerifyLink)
		linkGroup.POST("/batch-verify", linkHandler.BatchVerifyLinks)
		linkGroup.GET("/categories", linkHandler.GetLinkCategories)
		linkGroup.PUT("/sort", linkHandler.UpdateLinkSort)
	}

	// 组件管理路由
	componentGroup := adminGroup.Group("/components")
	{
		componentGroup.GET("", componentHandler.GetComponentList)
		componentGroup.POST("", componentHandler.CreateComponent)
		componentGroup.GET("/:id", componentHandler.GetComponentDetail)
		componentGroup.PUT("/:id", componentHandler.UpdateComponent)
		componentGroup.DELETE("/:id", componentHandler.DeleteComponent)
		componentGroup.GET("/:id/preview", componentHandler.PreviewComponent)
		componentGroup.POST("/import", componentHandler.ImportComponent)
		componentGroup.GET("/types", componentHandler.GetComponentTypes)
	}

	// 主题管理路由
	themeGroup := adminGroup.Group("/themes")
	{
		themeGroup.GET("", themeHandler.GetThemeList)
		themeGroup.POST("", themeHandler.CreateTheme)
		themeGroup.GET("/:id", themeHandler.GetThemeDetail)
		themeGroup.PUT("/:id", themeHandler.UpdateTheme)
		themeGroup.DELETE("/:id", themeHandler.DeleteTheme)
		themeGroup.POST("/:id/apply", themeHandler.ApplyTheme)
		themeGroup.GET("/:id/preview", themeHandler.PreviewTheme)
		themeGroup.GET("/:id/export", themeHandler.ExportTheme)
		themeGroup.POST("/import", themeHandler.ImportTheme)
		themeGroup.GET("/current", themeHandler.GetCurrentTheme)
	}
}

// SetupTreeNavigationRoutes 设置树形导航管理路由
func SetupTreeNavigationRoutes(
	r *gin.Engine,
	treeNavigationClient client.TreeNavigationClientInterface,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	// 树形导航处理器
	treeNavHandler := navigation.NewTreeNavigationHandler(treeNavigationClient)

	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(jwtMiddleware.Handler())

	// 树形导航管理路由
	treeNavGroup := adminGroup.Group("/tree-navigation")
	{
		// 基础CRUD操作
		treeNavGroup.GET("", treeNavHandler.GetNavigationTree)               // 获取导航树 ?type=main|footer|sidebar
		treeNavGroup.POST("", treeNavHandler.CreateNavigationItem)           // 创建导航项
		treeNavGroup.GET("/:id", treeNavHandler.GetNavigationItem)           // 获取导航项详情
		treeNavGroup.PUT("/:id", treeNavHandler.UpdateNavigationItem)        // 更新导航项
		treeNavGroup.DELETE("/:id", treeNavHandler.DeleteNavigationItem)     // 删除导航项

		// 高级功能
		treeNavGroup.PUT("/order", treeNavHandler.UpdateNavigationOrder)               // 更新排序
		treeNavGroup.PUT("/:id/visibility", treeNavHandler.ToggleNavigationVisibility) // 切换可见性
		treeNavGroup.DELETE("/batch", treeNavHandler.BatchDeleteNavigationItems)       // 批量删除

		// 导入导出
		treeNavGroup.GET("/export", treeNavHandler.ExportNavigationTree)  // 导出导航树
		treeNavGroup.POST("/import", treeNavHandler.ImportNavigationTree) // 导入导航树
	}
}
