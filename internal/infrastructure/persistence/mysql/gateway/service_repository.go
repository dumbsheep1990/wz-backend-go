package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"wz-backend-go/internal/domain/gateway/entity"
	"wz-backend-go/internal/domain/gateway/repository"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

// MySQLServiceRepository implements the ServiceRepository interface using MySQL
type MySQLServiceRepository struct {
	db *sql.DB
}

// NewMySQLServiceRepository creates a new MySQLServiceRepository
func NewMySQLServiceRepository(db *sql.DB) *MySQLServiceRepository {
	return &MySQLServiceRepository{
		db: db,
	}
}

// Save saves a service to the database
func (r *MySQLServiceRepository) Save(ctx context.Context, service *entity.Service) error {
	// Check if service already exists
	existingService, err := r.FindByName(ctx, service.Name())
	if err != nil {
		return fmt.Errorf("failed to check if service exists: %w", err)
	}

	// Extract health URL string, handling nil URL
	var healthURLStr sql.NullString
	if url := service.HealthURL().URL(); url != nil {
		healthURLStr = sql.NullString{
			String: service.HealthURL().String(),
			Valid:  true,
		}
	}

	// Extract last healthy time, handling nil time
	var lastHealthyTime sql.NullTime
	if lastHealthy := service.LastHealthy(); lastHealthy != nil {
		lastHealthyTime = sql.NullTime{
			Time:  *lastHealthy,
			Valid: true,
		}
	}

	if existingService == nil {
		// Insert new service
		query := `
			INSERT INTO gateway_services (
				name, 
				base_url, 
				health_url, 
				default_auth, 
				rate_limit_type, 
				rate_limit_requests, 
				rate_limit_duration_seconds, 
				is_active, 
				is_healthy,
				last_healthy,
				error_message,
				created_at, 
				updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			service.Name().String(),
			service.BaseURL().String(),
			healthURLStr,
			service.DefaultAuth().String(),
			service.RateLimit().Type().String(),
			service.RateLimit().Requests(),
			int(service.RateLimit().Duration().Seconds()),
			service.IsActive(),
			service.IsHealthy(),
			lastHealthyTime,
			service.ErrorMessage(),
			service.CreatedAt(),
			service.UpdatedAt(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert service: %w", err)
		}
	} else {
		// Update existing service
		query := `
			UPDATE gateway_services SET 
				base_url = ?, 
				health_url = ?, 
				default_auth = ?, 
				rate_limit_type = ?, 
				rate_limit_requests = ?, 
				rate_limit_duration_seconds = ?, 
				is_active = ?, 
				is_healthy = ?,
				last_healthy = ?,
				error_message = ?,
				updated_at = ?
			WHERE name = ?
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			service.BaseURL().String(),
			healthURLStr,
			service.DefaultAuth().String(),
			service.RateLimit().Type().String(),
			service.RateLimit().Requests(),
			int(service.RateLimit().Duration().Seconds()),
			service.IsActive(),
			service.IsHealthy(),
			lastHealthyTime,
			service.ErrorMessage(),
			time.Now(),
			service.Name().String(),
		)
		if err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}
	}

	return nil
}

// FindByName finds a service by name
func (r *MySQLServiceRepository) FindByName(ctx context.Context, name valueobject.ServiceName) (*entity.Service, error) {
	query := `
		SELECT 
			name, 
			base_url, 
			health_url, 
			default_auth, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			is_healthy,
			last_healthy,
			error_message,
			created_at, 
			updated_at
		FROM gateway_services
		WHERE name = ?
	`

	row := r.db.QueryRowContext(ctx, query, name.String())
	return r.scanService(row)
}

// FindActive finds all active services
func (r *MySQLServiceRepository) FindActive(ctx context.Context) ([]*entity.Service, error) {
	query := `
		SELECT 
			name, 
			base_url, 
			health_url, 
			default_auth, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			is_healthy,
			last_healthy,
			error_message,
			created_at, 
			updated_at
		FROM gateway_services
		WHERE is_active = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active services: %w", err)
	}
	defer rows.Close()

	return r.scanServices(rows)
}

// FindHealthy finds all healthy services
func (r *MySQLServiceRepository) FindHealthy(ctx context.Context) ([]*entity.Service, error) {
	query := `
		SELECT 
			name, 
			base_url, 
			health_url, 
			default_auth, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			is_healthy,
			last_healthy,
			error_message,
			created_at, 
			updated_at
		FROM gateway_services
		WHERE is_healthy = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query healthy services: %w", err)
	}
	defer rows.Close()

	return r.scanServices(rows)
}

// FindAll finds all services with pagination
func (r *MySQLServiceRepository) FindAll(ctx context.Context, offset, limit int) ([]*entity.Service, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM gateway_services`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count services: %w", err)
	}

	// Get services with pagination
	query := `
		SELECT 
			name, 
			base_url, 
			health_url, 
			default_auth, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			is_healthy,
			last_healthy,
			error_message,
			created_at, 
			updated_at
		FROM gateway_services
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	services, err := r.scanServices(rows)
	if err != nil {
		return nil, 0, err
	}

	return services, totalCount, nil
}

// Delete deletes a service by name
func (r *MySQLServiceRepository) Delete(ctx context.Context, name valueobject.ServiceName) error {
	query := `DELETE FROM gateway_services WHERE name = ?`
	_, err := r.db.ExecContext(ctx, query, name.String())
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}

// scanService scans a single service from a row
func (r *MySQLServiceRepository) scanService(row *sql.Row) (*entity.Service, error) {
	var (
		nameStr                string
		baseURLStr             string
		healthURLStr           sql.NullString
		defaultAuthStr         string
		rateLimitTypeStr       string
		rateLimitRequests      int
		rateLimitDurationSecs  int
		isActive               bool
		isHealthy              bool
		lastHealthy            sql.NullTime
		errorMessage           string
		createdAt              time.Time
		updatedAt              time.Time
	)

	err := row.Scan(
		&nameStr,
		&baseURLStr,
		&healthURLStr,
		&defaultAuthStr,
		&rateLimitTypeStr,
		&rateLimitRequests,
		&rateLimitDurationSecs,
		&isActive,
		&isHealthy,
		&lastHealthy,
		&errorMessage,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan service: %w", err)
	}

	// Create value objects
	serviceName, err := valueobject.NewServiceName(nameStr)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}

	baseURL, err := valueobject.NewTargetURL(baseURLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	var healthURL valueobject.TargetURL
	if healthURLStr.Valid {
		healthURL, err = valueobject.NewTargetURL(healthURLStr.String)
		if err != nil {
			return nil, fmt.Errorf("invalid health URL: %w", err)
		}
	} else {
		healthURL, _ = valueobject.NewEmptyTargetURL()
	}

	defaultAuth := valueobject.AuthTypeFromString(defaultAuthStr)

	rateLimitType := valueobject.RateLimitTypeFromString(rateLimitTypeStr)
	rateLimit, err := valueobject.NewRateLimit(
		rateLimitType,
		rateLimitRequests,
		rateLimitDurationSecs,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid rate limit: %w", err)
	}

	// Create service entity
	serviceBuilder := entity.NewServiceBuilder()
	serviceBuilder = serviceBuilder.
		WithName(serviceName).
		WithBaseURL(baseURL).
		WithHealthURL(healthURL).
		WithDefaultAuth(defaultAuth).
		WithRateLimit(rateLimit).
		WithActive(isActive).
		WithHealthy(isHealthy).
		WithErrorMessage(errorMessage).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt)

	// Add last healthy time if available
	if lastHealthy.Valid {
		lastHealthyTime := lastHealthy.Time
		serviceBuilder = serviceBuilder.WithLastHealthy(&lastHealthyTime)
	}

	service := serviceBuilder.Build()
	return service, nil
}

// scanServices scans multiple services from rows
func (r *MySQLServiceRepository) scanServices(rows *sql.Rows) ([]*entity.Service, error) {
	var services []*entity.Service

	for rows.Next() {
		var (
			nameStr                string
			baseURLStr             string
			healthURLStr           sql.NullString
			defaultAuthStr         string
			rateLimitTypeStr       string
			rateLimitRequests      int
			rateLimitDurationSecs  int
			isActive               bool
			isHealthy              bool
			lastHealthy            sql.NullTime
			errorMessage           string
			createdAt              time.Time
			updatedAt              time.Time
		)

		err := rows.Scan(
			&nameStr,
			&baseURLStr,
			&healthURLStr,
			&defaultAuthStr,
			&rateLimitTypeStr,
			&rateLimitRequests,
			&rateLimitDurationSecs,
			&isActive,
			&isHealthy,
			&lastHealthy,
			&errorMessage,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}

		// Create value objects
		serviceName, err := valueobject.NewServiceName(nameStr)
		if err != nil {
			return nil, fmt.Errorf("invalid service name: %w", err)
		}

		baseURL, err := valueobject.NewTargetURL(baseURLStr)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %w", err)
		}

		var healthURL valueobject.TargetURL
		if healthURLStr.Valid {
			healthURL, err = valueobject.NewTargetURL(healthURLStr.String)
			if err != nil {
				return nil, fmt.Errorf("invalid health URL: %w", err)
			}
		} else {
			healthURL, _ = valueobject.NewEmptyTargetURL()
		}

		defaultAuth := valueobject.AuthTypeFromString(defaultAuthStr)

		rateLimitType := valueobject.RateLimitTypeFromString(rateLimitTypeStr)
		rateLimit, err := valueobject.NewRateLimit(
			rateLimitType,
			rateLimitRequests,
			rateLimitDurationSecs,
		)
		if err != nil {
			return nil, fmt.Errorf("invalid rate limit: %w", err)
		}

		// Create service entity
		serviceBuilder := entity.NewServiceBuilder()
		serviceBuilder = serviceBuilder.
			WithName(serviceName).
			WithBaseURL(baseURL).
			WithHealthURL(healthURL).
			WithDefaultAuth(defaultAuth).
			WithRateLimit(rateLimit).
			WithActive(isActive).
			WithHealthy(isHealthy).
			WithErrorMessage(errorMessage).
			WithCreatedAt(createdAt).
			WithUpdatedAt(updatedAt)

		// Add last healthy time if available
		if lastHealthy.Valid {
			lastHealthyTime := lastHealthy.Time
			serviceBuilder = serviceBuilder.WithLastHealthy(&lastHealthyTime)
		}

		service := serviceBuilder.Build()
		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return services, nil
}

// EnsureTable creates the gateway_services table if it doesn't exist
func (r *MySQLServiceRepository) EnsureTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS gateway_services (
			name VARCHAR(100) PRIMARY KEY,
			base_url VARCHAR(500) NOT NULL,
			health_url VARCHAR(500),
			default_auth VARCHAR(20) NOT NULL,
			rate_limit_type VARCHAR(20) NOT NULL,
			rate_limit_requests INT NOT NULL,
			rate_limit_duration_seconds INT NOT NULL,
			is_active BOOLEAN NOT NULL,
			is_healthy BOOLEAN NOT NULL,
			last_healthy TIMESTAMP NULL,
			error_message TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			INDEX idx_is_active (is_active),
			INDEX idx_is_healthy (is_healthy)
		)
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create gateway_services table: %w", err)
	}
	return nil
}
