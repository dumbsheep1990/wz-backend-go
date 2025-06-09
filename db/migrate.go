package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationOptions holds configuration options for migrations
type MigrationOptions struct {
	DB         *sql.DB
	SchemaName string
	TableName  string // Migration table name
	Directory  string // Base migrations directory
	Domains    []string // Specific domains to migrate, empty means all
}

// RunMigrations executes migrations from the specified domains
func RunMigrations(opts MigrationOptions) error {
	if opts.DB == nil {
		return fmt.Errorf("database connection is required")
	}

	if opts.TableName == "" {
		opts.TableName = "schema_migrations"
	}

	if opts.Directory == "" {
		opts.Directory = "db/migrations"
	}

	// If no domains specified, get all directories under migrations folder
	var domains []string
	if len(opts.Domains) == 0 {
		entries, err := os.ReadDir(opts.Directory)
		if err != nil {
			return fmt.Errorf("failed to read migrations directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() && entry.Name() != "schema" {
				domains = append(domains, entry.Name())
			}
		}
	} else {
		domains = opts.Domains
	}

	// Configure the MySQL driver for migrations
	driver, err := mysql.WithInstance(opts.DB, &mysql.Config{
		MigrationsTable: opts.TableName,
		DatabaseName:    opts.SchemaName,
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Run migrations for each domain
	for _, domain := range domains {
		log.Printf("Running migrations for domain: %s", domain)
		
		sourceURL := fmt.Sprintf("file://%s", filepath.Join(opts.Directory, domain))
		m, err := migrate.NewWithDatabaseInstance(
			sourceURL,
			opts.SchemaName,
			driver,
		)
		if err != nil {
			return fmt.Errorf("failed to create migration instance for domain %s: %w", domain, err)
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migration failed for domain %s: %w", domain, err)
		}
		
		log.Printf("Migration completed for domain: %s", domain)
	}

	return nil
}

// CreateNewMigration creates a new migration file in the specified domain directory
func CreateNewMigration(domain, name string) error {
	if domain == "" || name == "" {
		return fmt.Errorf("domain and name are required")
	}
	
	// Clean name to be valid as a filename
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)
	
	domainDir := filepath.Join("db/migrations", domain)
	
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
				seqStr := parts[0]
				if seq, err := fmt.Sscanf(seqStr, "%03d", &nextSeq); err == nil && seq > 0 {
					nextSeq = seq + 1
				}
			}
		}
	}
	
	// Create the migration files
	upFile := filepath.Join(domainDir, fmt.Sprintf("%03d_%s.up.sql", nextSeq, name))
	downFile := filepath.Join(domainDir, fmt.Sprintf("%03d_%s.down.sql", nextSeq, name))
	
	// Create empty migration files
	if err := os.WriteFile(upFile, []byte("-- Migration Up\n\n"), 0644); err != nil {
		return fmt.Errorf("failed to create up migration: %w", err)
	}
	
	if err := os.WriteFile(downFile, []byte("-- Migration Down\n\n"), 0644); err != nil {
		return fmt.Errorf("failed to create down migration: %w", err)
	}
	
	log.Printf("Created new migration files:\n%s\n%s", upFile, downFile)
	return nil
}
