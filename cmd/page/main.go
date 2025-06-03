package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"

	pb "wz-backend-go/api/page"
	"wz-backend-go/internal/infrastructure/assembly"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// SQLTransactionManager implements the database.TransactionManager interface
type SQLTransactionManager struct {
	db *sql.DB
}

// BeginTx starts a new transaction
func (tm *SQLTransactionManager) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return ctx, err
	}
	// Store the transaction in the context
	return context.WithValue(ctx, "tx", tx), nil
}

// CommitTx commits a transaction
func (tm *SQLTransactionManager) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Commit()
}

// RollbackTx rolls back a transaction
func (tm *SQLTransactionManager) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Rollback()
}

func main() {
	// Get database connection details from environment variables
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "wz")
	
	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// Check database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	
	// Create transaction manager
	txManager := &SQLTransactionManager{db: db}
	
	// Assemble the page service
	pageAssembly := assembly.NewPageServiceAssembly(db, txManager)
	
	// Start gRPC server
	port := getEnv("PORT", "50052")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterPageServiceServer(grpcServer, pageAssembly.GRPCAdapter)
	
	// Handle graceful shutdown
	go func() {
		log.Printf("Starting page gRPC server on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down page gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Page server exited properly")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
