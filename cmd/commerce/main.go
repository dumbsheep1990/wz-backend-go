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
	
	"wz-backend-go/internal/application/commerce"
	"wz-backend-go/internal/delivery/http/router"
	commerceHandler "wz-backend-go/internal/delivery/http/handler/commerce"
	"wz-backend-go/internal/domain/commerce/service"
	mysqlcommerce "wz-backend-go/internal/infrastructure/persistence/mysql/commerce"
)

func main() {
	// Setup database connection
	dbConn, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer dbConn.Close()

	// Create repositories
	productRepo := mysqlcommerce.NewMySQLProductRepository(dbConn)
	storeRepo := mysqlcommerce.NewMySQLStoreRepository(dbConn)
	categoryRepo := mysqlcommerce.NewMySQLCategoryRepository(dbConn)

	// Create domain service
	commerceService := service.NewCommerceService(productRepo, storeRepo, categoryRepo)

	// Create application services
	productAppService := commerce.NewProductAppService(commerceService, productRepo, storeRepo, categoryRepo)
	storeAppService := commerce.NewStoreAppService(commerceService, storeRepo, productRepo)
	categoryAppService := commerce.NewCategoryAppService(commerceService, categoryRepo)

	// Create HTTP handlers
	productHandler := commerceHandler.NewProductHandler(productAppService)
	storeHandler := commerceHandler.NewStoreHandler(storeAppService)
	categoryHandler := commerceHandler.NewCategoryHandler(categoryAppService)

	// Setup Gin router
	ginRouter := gin.Default()
	ginRouter.Use(gin.Recovery())

	// Register routes
	router.RegisterCommerceRoutes(ginRouter, productHandler, storeHandler, categoryHandler)

	// Start server
	srv := &http.Server{
		Addr:    ":8081",
		Handler: ginRouter,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Starting commerce microservice on port 8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down commerce microservice...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Commerce microservice exited")
}

// setupDatabase establishes a connection to the MySQL database
func setupDatabase() (*sqlx.DB, error) {
	// Get database connection details from environment variables
	host := getEnvWithDefault("DB_HOST", "localhost")
	port := getEnvWithDefault("DB_PORT", "3306")
	user := getEnvWithDefault("DB_USER", "root")
	password := getEnvWithDefault("DB_PASSWORD", "password")
	dbname := getEnvWithDefault("DB_NAME", "wanzhi")

	// Create DSN string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	// Open database connection
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

// getEnvWithDefault returns the value of an environment variable or a default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
