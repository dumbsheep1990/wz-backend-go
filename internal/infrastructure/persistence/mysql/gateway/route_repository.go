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

// MySQLRouteRepository implements the RouteRepository interface using MySQL
type MySQLRouteRepository struct {
	db *sql.DB
}

// NewMySQLRouteRepository creates a new MySQLRouteRepository
func NewMySQLRouteRepository(db *sql.DB) *MySQLRouteRepository {
	return &MySQLRouteRepository{
		db: db,
	}
}

// Save saves a route to the database
func (r *MySQLRouteRepository) Save(ctx context.Context, route *entity.Route) error {
	// Check if route already exists
	existingRoute, err := r.FindByID(ctx, route.ID())
	if err != nil {
		return fmt.Errorf("failed to check if route exists: %w", err)
	}

	if existingRoute == nil {
		// Insert new route
		query := `
			INSERT INTO gateway_routes (
				id, 
				path, 
				service_name, 
				target_url, 
				auth_type, 
				rate_limit_type, 
				rate_limit_requests, 
				rate_limit_duration_seconds, 
				is_active, 
				created_at, 
				updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			route.ID().String(),
			route.Path().String(),
			route.ServiceName().String(),
			route.TargetURL().String(),
			route.AuthType().String(),
			route.RateLimit().Type().String(),
			route.RateLimit().Requests(),
			int(route.RateLimit().Duration().Seconds()),
			route.IsActive(),
			route.CreatedAt(),
			route.UpdatedAt(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert route: %w", err)
		}
	} else {
		// Update existing route
		query := `
			UPDATE gateway_routes SET 
				path = ?, 
				service_name = ?, 
				target_url = ?, 
				auth_type = ?, 
				rate_limit_type = ?, 
				rate_limit_requests = ?, 
				rate_limit_duration_seconds = ?, 
				is_active = ?, 
				updated_at = ?
			WHERE id = ?
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			route.Path().String(),
			route.ServiceName().String(),
			route.TargetURL().String(),
			route.AuthType().String(),
			route.RateLimit().Type().String(),
			route.RateLimit().Requests(),
			int(route.RateLimit().Duration().Seconds()),
			route.IsActive(),
			time.Now(),
			route.ID().String(),
		)
		if err != nil {
			return fmt.Errorf("failed to update route: %w", err)
		}
	}

	return nil
}

// FindByID finds a route by ID
func (r *MySQLRouteRepository) FindByID(ctx context.Context, id valueobject.RouteID) (*entity.Route, error) {
	query := `
		SELECT 
			id, 
			path, 
			service_name, 
			target_url, 
			auth_type, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			created_at, 
			updated_at
		FROM gateway_routes
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id.String())
	return r.scanRoute(row)
}

// FindByPath finds routes by path
func (r *MySQLRouteRepository) FindByPath(ctx context.Context, path valueobject.Path) ([]*entity.Route, error) {
	query := `
		SELECT 
			id, 
			path, 
			service_name, 
			target_url, 
			auth_type, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			created_at, 
			updated_at
		FROM gateway_routes
		WHERE path = ? AND is_active = true
	`

	rows, err := r.db.QueryContext(ctx, query, path.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query routes by path: %w", err)
	}
	defer rows.Close()

	return r.scanRoutes(rows)
}

// FindByServiceName finds routes by service name
func (r *MySQLRouteRepository) FindByServiceName(ctx context.Context, serviceName valueobject.ServiceName) ([]*entity.Route, error) {
	query := `
		SELECT 
			id, 
			path, 
			service_name, 
			target_url, 
			auth_type, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			created_at, 
			updated_at
		FROM gateway_routes
		WHERE service_name = ?
	`

	rows, err := r.db.QueryContext(ctx, query, serviceName.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query routes by service name: %w", err)
	}
	defer rows.Close()

	return r.scanRoutes(rows)
}

// FindActive finds all active routes
func (r *MySQLRouteRepository) FindActive(ctx context.Context) ([]*entity.Route, error) {
	query := `
		SELECT 
			id, 
			path, 
			service_name, 
			target_url, 
			auth_type, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			created_at, 
			updated_at
		FROM gateway_routes
		WHERE is_active = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active routes: %w", err)
	}
	defer rows.Close()

	return r.scanRoutes(rows)
}

// FindAll finds all routes with pagination
func (r *MySQLRouteRepository) FindAll(ctx context.Context, offset, limit int) ([]*entity.Route, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM gateway_routes`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count routes: %w", err)
	}

	// Get routes with pagination
	query := `
		SELECT 
			id, 
			path, 
			service_name, 
			target_url, 
			auth_type, 
			rate_limit_type, 
			rate_limit_requests, 
			rate_limit_duration_seconds, 
			is_active, 
			created_at, 
			updated_at
		FROM gateway_routes
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query routes: %w", err)
	}
	defer rows.Close()

	routes, err := r.scanRoutes(rows)
	if err != nil {
		return nil, 0, err
	}

	return routes, totalCount, nil
}

// Delete deletes a route by ID
func (r *MySQLRouteRepository) Delete(ctx context.Context, id valueobject.RouteID) error {
	query := `DELETE FROM gateway_routes WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	return nil
}

// scanRoute scans a single route from a row
func (r *MySQLRouteRepository) scanRoute(row *sql.Row) (*entity.Route, error) {
	var (
		id                     string
		pathStr                string
		serviceNameStr         string
		targetURLStr           string
		authTypeStr            string
		rateLimitTypeStr       string
		rateLimitRequests      int
		rateLimitDurationSecs  int
		isActive               bool
		createdAt              time.Time
		updatedAt              time.Time
	)

	err := row.Scan(
		&id,
		&pathStr,
		&serviceNameStr,
		&targetURLStr,
		&authTypeStr,
		&rateLimitTypeStr,
		&rateLimitRequests,
		&rateLimitDurationSecs,
		&isActive,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan route: %w", err)
	}

	// Create value objects
	routeID, err := valueobject.NewRouteID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid route ID: %w", err)
	}

	path, err := valueobject.NewPath(pathStr)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	serviceName, err := valueobject.NewServiceName(serviceNameStr)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}

	targetURL, err := valueobject.NewTargetURL(targetURLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}

	authType := valueobject.AuthTypeFromString(authTypeStr)

	rateLimitType := valueobject.RateLimitTypeFromString(rateLimitTypeStr)
	rateLimit, err := valueobject.NewRateLimit(
		rateLimitType,
		rateLimitRequests,
		rateLimitDurationSecs,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid rate limit: %w", err)
	}

	// Create route entity
	routeBuilder := entity.NewRouteBuilder()
	route := routeBuilder.
		WithID(routeID).
		WithPath(path).
		WithServiceName(serviceName).
		WithTargetURL(targetURL).
		WithAuthType(authType).
		WithRateLimit(rateLimit).
		WithActive(isActive).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()

	return route, nil
}

// scanRoutes scans multiple routes from rows
func (r *MySQLRouteRepository) scanRoutes(rows *sql.Rows) ([]*entity.Route, error) {
	var routes []*entity.Route

	for rows.Next() {
		var (
			id                     string
			pathStr                string
			serviceNameStr         string
			targetURLStr           string
			authTypeStr            string
			rateLimitTypeStr       string
			rateLimitRequests      int
			rateLimitDurationSecs  int
			isActive               bool
			createdAt              time.Time
			updatedAt              time.Time
		)

		err := rows.Scan(
			&id,
			&pathStr,
			&serviceNameStr,
			&targetURLStr,
			&authTypeStr,
			&rateLimitTypeStr,
			&rateLimitRequests,
			&rateLimitDurationSecs,
			&isActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route: %w", err)
		}

		// Create value objects
		routeID, err := valueobject.NewRouteID(id)
		if err != nil {
			return nil, fmt.Errorf("invalid route ID: %w", err)
		}

		path, err := valueobject.NewPath(pathStr)
		if err != nil {
			return nil, fmt.Errorf("invalid path: %w", err)
		}

		serviceName, err := valueobject.NewServiceName(serviceNameStr)
		if err != nil {
			return nil, fmt.Errorf("invalid service name: %w", err)
		}

		targetURL, err := valueobject.NewTargetURL(targetURLStr)
		if err != nil {
			return nil, fmt.Errorf("invalid target URL: %w", err)
		}

		authType := valueobject.AuthTypeFromString(authTypeStr)

		rateLimitType := valueobject.RateLimitTypeFromString(rateLimitTypeStr)
		rateLimit, err := valueobject.NewRateLimit(
			rateLimitType,
			rateLimitRequests,
			rateLimitDurationSecs,
		)
		if err != nil {
			return nil, fmt.Errorf("invalid rate limit: %w", err)
		}

		// Create route entity
		routeBuilder := entity.NewRouteBuilder()
		route := routeBuilder.
			WithID(routeID).
			WithPath(path).
			WithServiceName(serviceName).
			WithTargetURL(targetURL).
			WithAuthType(authType).
			WithRateLimit(rateLimit).
			WithActive(isActive).
			WithCreatedAt(createdAt).
			WithUpdatedAt(updatedAt).
			Build()

		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return routes, nil
}

// EnsureTable creates the gateway_routes table if it doesn't exist
func (r *MySQLRouteRepository) EnsureTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS gateway_routes (
			id VARCHAR(100) PRIMARY KEY,
			path VARCHAR(500) NOT NULL,
			service_name VARCHAR(100) NOT NULL,
			target_url VARCHAR(500) NOT NULL,
			auth_type VARCHAR(20) NOT NULL,
			rate_limit_type VARCHAR(20) NOT NULL,
			rate_limit_requests INT NOT NULL,
			rate_limit_duration_seconds INT NOT NULL,
			is_active BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			INDEX idx_service_name (service_name),
			INDEX idx_path (path),
			INDEX idx_is_active (is_active)
		)
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create gateway_routes table: %w", err)
	}
	return nil
}
