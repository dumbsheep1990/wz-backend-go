package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

// Import your project's db package - adjust this path according to your Go module name
// For example: "github.com/your-username/wz-backend-go/db"
// Since I don't know your exact module path, you'll need to update this import statement
// after this file is created
var dbPackage = "your-module/db"

func main() {
	// Load environment variables from .env file if present
	_ = godotenv.Load()

	// Define command line flags
	createCmd := flag.Bool("create", false, "Create a new migration")
	migrateCmd := flag.Bool("migrate", false, "Run migrations")
	domainFlag := flag.String("domain", "", "Domain to create/run migrations for (comma-separated for multiple)")
	nameFlag := flag.String("name", "", "Name for the new migration")
	dbURLFlag := flag.String("db", os.Getenv("DATABASE_URL"), "Database connection URL")
	schemaFlag := flag.String("schema", "wz_backend", "Database schema name")

	flag.Parse()

	if !*createCmd && !*migrateCmd {
		fmt.Println("Please specify either -create or -migrate")
		flag.Usage()
		os.Exit(1)
	}

	if *createCmd {
		if *domainFlag == "" || *nameFlag == "" {
			fmt.Println("Domain and name are required for creating migrations")
			flag.Usage()
			os.Exit(1)
		}

		// For now, implement direct file creation functionality here
		// Later you can move this to the db package once you update its import path
		if err := createNewMigration(*domainFlag, *nameFlag); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		return
	}

	if *migrateCmd {
		if *dbURLFlag == "" {
			fmt.Println("Database URL is required for running migrations")
			fmt.Println("Set it via -db flag or DATABASE_URL environment variable")
			flag.Usage()
			os.Exit(1)
		}

		// Open database connection
		dbConn, err := sql.Open("mysql", *dbURLFlag)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer dbConn.Close()

		// Ping the database to verify connection
		if err := dbConn.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		var domains []string
		if *domainFlag != "" {
			domains = strings.Split(*domainFlag, ",")
		}

		fmt.Println("Running migrations for domains:", domains)
		fmt.Println("This functionality will be implemented in the db package")
		fmt.Println("Please update the import path in this file once ready")
		
		// Placeholder for when you update the import path to call your db package
		// err = db.RunMigrations(db.MigrationOptions{
		//     DB:         dbConn,
		//     SchemaName: *schemaFlag,
		//     TableName:  "schema_migrations",
		//     Directory:  "db/migrations",
		//     Domains:    domains,
		// })
		
		fmt.Println("Migration simulation completed successfully")
	}
}

// createNewMigration creates a new migration file in the specified domain directory
// This function can be moved to the db package once import path is updated
func createNewMigration(domain, name string) error {
	if domain == "" || name == "" {
		return fmt.Errorf("domain and name are required")
	}
	
	// Clean name to be valid as a filename
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)
	
	baseDir := "db/migrations"
	domainDir := fmt.Sprintf("%s/%s", baseDir, domain)
	
	// Ensure domain directory exists
	if _, err := os.Stat(domainDir); os.IsNotExist(err) {
		if err := os.MkdirAll(domainDir, 0755); err != nil {
			return fmt.Errorf("failed to create domain directory: %w", err)
		}
	}
	
	// Find the next sequence number
	entries, err := os.ReadDir(domainDir)
	if err != nil {
		return fmt.Errorf("failed to read domain directory: %w", err)
	}
	
	nextSeq := 1
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			parts := strings.Split(entry.Name(), "_")
			if len(parts) > 0 {
				var seq int
				if _, err := fmt.Sscanf(parts[0], "%03d", &seq); err == nil && seq >= nextSeq {
					nextSeq = seq + 1
				}
			}
		}
	}
	
	// Create the migration files
	upFile := fmt.Sprintf("%s/%03d_%s.up.sql", domainDir, nextSeq, name)
	downFile := fmt.Sprintf("%s/%03d_%s.down.sql", domainDir, nextSeq, name)
	
	// Create empty migration files
	if err := os.WriteFile(upFile, []byte("-- Migration Up for "+domain+"\n\n"), 0644); err != nil {
		return fmt.Errorf("failed to create up migration: %w", err)
	}
	
	if err := os.WriteFile(downFile, []byte("-- Migration Down for "+domain+"\n\n"), 0644); err != nil {
		return fmt.Errorf("failed to create down migration: %w", err)
	}
	
	log.Printf("Created new migration files:\n%s\n%s", upFile, downFile)
	return nil
}
