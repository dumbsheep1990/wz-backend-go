package router

import (
	"github.com/gin-gonic/gin"
	
	commerceHandler "wz-backend-go/internal/delivery/http/handler/commerce"
)

// RegisterCommerceRoutes registers all commerce module routes
func RegisterCommerceRoutes(
	router *gin.Engine,
	productHandler *commerceHandler.ProductHandler,
	storeHandler *commerceHandler.StoreHandler,
	categoryHandler *commerceHandler.CategoryHandler,
) {
	api := router.Group("/api/v1/commerce")
	
	// Product routes
	product := api.Group("/products")
	{
		product.POST("", productHandler.CreateProduct)
		product.PUT("/:id", productHandler.UpdateProduct)
		product.GET("/:id", productHandler.GetProductByID)
		product.GET("", productHandler.GetProducts)
		product.GET("/featured", productHandler.GetFeaturedProducts)
		product.GET("/new", productHandler.GetNewProducts)
		product.GET("/popular", productHandler.GetPopularProducts)
		product.POST("/:id/view", productHandler.TrackProductView)
		product.GET("/search", productHandler.SearchProducts)
	}
	
	// Store routes
	store := api.Group("/stores")
	{
		store.POST("", storeHandler.CreateStore)
		store.PUT("/:id", storeHandler.UpdateStore)
		store.GET("/:id", storeHandler.GetStoreByID)
		store.GET("/owner/:ownerID", storeHandler.GetStoresByOwner)
		store.GET("", storeHandler.GetStores)
		store.GET("/search", storeHandler.SearchStores)
	}
	
	// Category routes
	category := api.Group("/categories")
	{
		category.POST("", categoryHandler.CreateCategory)
		category.PUT("/:id", categoryHandler.UpdateCategory)
		category.GET("/:id", categoryHandler.GetCategoryByID)
		category.GET("/root", categoryHandler.GetRootCategories)
		category.GET("", categoryHandler.GetCategories)
		category.POST("/reorder", categoryHandler.ReorderCategories)
		category.GET("/hierarchy", categoryHandler.GetCategoryHierarchy)
	}
}
