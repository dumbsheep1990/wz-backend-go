package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	
	"wz-backend-go/internal/application/navigation"
	"wz-backend-go/internal/delivery/http/handler/navigation"
	"wz-backend-go/internal/delivery/http/router"
	"wz-backend-go/internal/domain/navigation/service"
	navigationRepo "wz-backend-go/internal/infrastructure/repository/mysql/navigation"
)

func main() {
	// Initialize database connection
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// Initialize repositories
	categoryRepo := navigationRepo.NewMySQLCategoryRepository(db)
	websiteRepo := navigationRepo.NewMySQLWebsiteRepository(db)
	
	// Initialize domain services
	navigationService := service.NewNavigationService(categoryRepo, websiteRepo)
	
	// Initialize application services
	navigationAppService := navigation.NewNavigationAppService(navigationService, categoryRepo, websiteRepo)
	
	// Initialize HTTP handlers
	categoryHandler := navigationHandler.NewCategoryHandler(navigationAppService)
	websiteHandler := navigationHandler.NewWebsiteHandler(navigationAppService)
	
	// Setup Gin router
	r := gin.Default()
	
	// API versioning
	v1 := r.Group("/api/v1")
	
	// Register routes
	router.SetupNavigationRoutes(v1, categoryHandler, websiteHandler)
	
	// Start HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	
	// Start server in a goroutine so we can gracefully shut down
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	log.Println("Navigation microservice started on :8080")
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Create a deadline context for the server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	
	log.Println("Server exited properly")
}

// setupDatabase initializes the database connection
func setupDatabase() (*sqlx.DB, error) {
	// Get database configuration from environment variables
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "wanzhidb")
	
	// Build connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// Connect to the database
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}

// getEnv gets environment variable or returns fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
