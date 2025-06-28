package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wanzhitouzi/wz-backend-go/internal/application/tourism"
	tourismHandler "github.com/wanzhitouzi/wz-backend-go/internal/delivery/http/handler/tourism"
	"github.com/wanzhitouzi/wz-backend-go/internal/delivery/http/middleware"
)

// TourismRouter handles routing for tourism-related endpoints
type TourismRouter struct {
	scenicSpotAppService *tourism.ScenicSpotAppService
	categoryAppService   *tourism.CategoryAppService
	reviewAppService     *tourism.ReviewAppService
}

// NewTourismRouter creates a new TourismRouter instance
func NewTourismRouter(
	scenicSpotAppService *tourism.ScenicSpotAppService,
	categoryAppService *tourism.CategoryAppService,
	reviewAppService *tourism.ReviewAppService,
) *TourismRouter {
	return &TourismRouter{
		scenicSpotAppService: scenicSpotAppService,
		categoryAppService:   categoryAppService,
		reviewAppService:     reviewAppService,
	}
}

// RegisterRoutes registers all tourism-related routes
func (r *TourismRouter) RegisterRoutes(router *gin.RouterGroup) {
	// Create handlers
	scenicSpotHandler := tourismHandler.NewScenicSpotHandler(r.scenicSpotAppService)
	categoryHandler := tourismHandler.NewCategoryHandler(r.categoryAppService)
	reviewHandler := tourismHandler.NewReviewHandler(r.reviewAppService)

	// Setup tourism routes group
	tourismRoutes := router.Group("/tourism")
	{
		// Scenic spot routes
		scenicSpotRoutes := tourismRoutes.Group("/scenic-spots")
		{
			// Public routes
			scenicSpotRoutes.GET("", scenicSpotHandler.List)
			scenicSpotRoutes.GET("/:id", scenicSpotHandler.Get)
			scenicSpotRoutes.GET("/search", scenicSpotHandler.Search)

			// Admin-only routes
			adminScenicSpotRoutes := scenicSpotRoutes.Group("")
			adminScenicSpotRoutes.Use(middleware.AdminOnly())
			{
				adminScenicSpotRoutes.POST("", scenicSpotHandler.Create)
				adminScenicSpotRoutes.PUT("/:id", scenicSpotHandler.Update)
				adminScenicSpotRoutes.DELETE("/:id", scenicSpotHandler.Delete)
			}
		}

		// Category routes
		categoryRoutes := tourismRoutes.Group("/categories")
		{
			// Public routes
			categoryRoutes.GET("", categoryHandler.List)
			categoryRoutes.GET("/tree", categoryHandler.GetTree)
			categoryRoutes.GET("/:id", categoryHandler.Get)
			categoryRoutes.GET("/:category_id/scenic-spots", scenicSpotHandler.ListByCategory)

			// Admin-only routes
			adminCategoryRoutes := categoryRoutes.Group("")
			adminCategoryRoutes.Use(middleware.AdminOnly())
			{
				adminCategoryRoutes.POST("", categoryHandler.Create)
				adminCategoryRoutes.PUT("/:id", categoryHandler.Update)
				adminCategoryRoutes.DELETE("/:id", categoryHandler.Delete)
			}
		}

		// Area routes
		areaRoutes := tourismRoutes.Group("/areas")
		{
			areaRoutes.GET("/:area/scenic-spots", scenicSpotHandler.ListByArea)
		}

		// Review routes
		reviewRoutes := tourismRoutes.Group("/reviews")
		{
			// Public routes
			reviewRoutes.GET("/:id", reviewHandler.Get)

			// Authenticated user routes
			authReviewRoutes := reviewRoutes.Group("")
			authReviewRoutes.Use(middleware.Authentication())
			{
				authReviewRoutes.POST("", reviewHandler.Create)
				authReviewRoutes.PUT("/:id", reviewHandler.Update)
				authReviewRoutes.DELETE("/:id", reviewHandler.Delete)
				authReviewRoutes.POST("/:id/like", reviewHandler.Like)
				authReviewRoutes.POST("/:id/unlike", reviewHandler.Unlike)
			}
		}

		// Scenic spot reviews routes
		tourismRoutes.GET("/scenic-spots/:scenic_spot_id/reviews", reviewHandler.ListByScenicSpot)
		
		// User reviews routes
		tourismRoutes.GET("/users/:user_id/reviews", reviewHandler.ListByUser)
	}
}
