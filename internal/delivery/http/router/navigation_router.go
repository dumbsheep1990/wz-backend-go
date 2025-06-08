package router

import (
	"github.com/gin-gonic/gin"
	navigationHandler "wz-backend-go/internal/delivery/http/handler/navigation"
)

// SetupNavigationRoutes configures all routes for the navigation microservice
func SetupNavigationRoutes(
	router *gin.RouterGroup,
	categoryHandler *navigationHandler.CategoryHandler,
	websiteHandler *navigationHandler.WebsiteHandler,
) {
	// Navigation API group
	navigationRouter := router.Group("/navigation")
	{
		// Category routes
		categoryRouter := navigationRouter.Group("/categories")
		{
			categoryRouter.GET("", categoryHandler.GetCategories)
			categoryRouter.POST("", categoryHandler.CreateCategory)
			categoryRouter.GET("/:id", categoryHandler.GetCategoryByID)
			categoryRouter.PUT("/:id", categoryHandler.UpdateCategory)
			categoryRouter.DELETE("/:id", categoryHandler.DeleteCategory)
			categoryRouter.POST("/reorder", categoryHandler.ReorderCategories)
		}

		// Website routes
		websiteRouter := navigationRouter.Group("/websites")
		{
			websiteRouter.GET("", websiteHandler.GetWebsites)
			websiteRouter.POST("", websiteHandler.CreateWebsite)
			websiteRouter.PUT("/:id", websiteHandler.UpdateWebsite)
			websiteRouter.DELETE("/:id", websiteHandler.DeleteWebsite)
			websiteRouter.POST("/reorder", websiteHandler.ReorderWebsites)
			websiteRouter.GET("/popular", websiteHandler.GetPopularWebsites)
			websiteRouter.POST("/:id/track-view", websiteHandler.TrackWebsiteView)
		}
	}
}
