package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/repository"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// AdminRepositoryImpl implements the AdminRepository interface using SQL database
type AdminRepositoryImpl struct {
	db *sql.DB
}

// NewAdminRepository creates a new AdminRepositoryImpl
func NewAdminRepository(db *sql.DB) *AdminRepositoryImpl {
	return &AdminRepositoryImpl{
		db: db,
	}
}

// Ensure AdminRepositoryImpl implements the AdminRepository interface
var _ repository.AdminRepository = (*AdminRepositoryImpl)(nil)

// FindByID finds an admin by ID
func (r *AdminRepositoryImpl) FindByID(ctx context.Context, id valueobject.AdminID) (*entity.Admin, error) {
	query := `
		SELECT id, username, password_hash, role_id, status, last_login, created_at, updated_at
		FROM admins
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id.Value())

	return r.scanAdmin(row)
}

// FindByUsername finds an admin by username
func (r *AdminRepositoryImpl) FindByUsername(ctx context.Context, username valueobject.Username) (*entity.Admin, error) {
	query := `
		SELECT id, username, password_hash, role_id, status, last_login, created_at, updated_at
		FROM admins
		WHERE username = ?
	`

	row := r.db.QueryRowContext(ctx, query, username.Value())

	return r.scanAdmin(row)
}

// Save persists an admin entity (creates or updates)
func (r *AdminRepositoryImpl) Save(ctx context.Context, admin *entity.Admin) error {
	// Check if admin exists
	existingAdmin, err := r.FindByID(ctx, admin.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingAdmin == nil {
		return r.Create(ctx, admin)
	} else {
		return r.Update(ctx, admin)
	}
}

// Create creates a new admin entity
func (r *AdminRepositoryImpl) Create(ctx context.Context, admin *entity.Admin) error {
	query := `
		INSERT INTO admins (
			id, username, password_hash, role_id, status, last_login, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		admin.ID.Value(),
		admin.Username.Value(),
		admin.PasswordHash,
		admin.Role.Value(),
		admin.Status.Value(),
		admin.LastLogin,
		admin.CreatedAt,
		admin.UpdatedAt,
	)

	return err
}

// Update updates an existing admin entity
func (r *AdminRepositoryImpl) Update(ctx context.Context, admin *entity.Admin) error {
	query := `
		UPDATE admins
		SET username = ?, 
			password_hash = ?, 
			role_id = ?, 
			status = ?, 
			last_login = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		admin.Username.Value(),
		admin.PasswordHash,
		admin.Role.Value(),
		admin.Status.Value(),
		admin.LastLogin,
		time.Now(),
		admin.ID.Value(),
	)

	return err
}

// Delete deletes an admin entity
func (r *AdminRepositoryImpl) Delete(ctx context.Context, id valueobject.AdminID) error {
	query := `DELETE FROM admins WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id.Value())
	return err
}

// List lists admin entities with pagination and filters
func (r *AdminRepositoryImpl) List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*entity.Admin, int64, error) {
	// Build the WHERE clause based on filters
	whereClause, args := r.buildWhereClause(filters)

	// Count total matching records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM admins %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query with pagination
	query := fmt.Sprintf(`
		SELECT id, username, password_hash, role_id, status, last_login, created_at, updated_at
		FROM admins
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// Add pagination params to args
	args = append(args, pageSize, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Scan results
	var admins []*entity.Admin
	for rows.Next() {
		admin, err := r.scanAdminRow(rows)
		if err != nil {
			return nil, 0, err
		}
		admins = append(admins, admin)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

// UpdateLastLogin updates the last login time for an admin
func (r *AdminRepositoryImpl) UpdateLastLogin(ctx context.Context, id valueobject.AdminID) error {
	query := `
		UPDATE admins
		SET last_login = ?,
			updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id.Value())
	return err
}

// CountAdmins counts the total number of admins, optionally filtered
func (r *AdminRepositoryImpl) CountAdmins(ctx context.Context, filters map[string]interface{}) (int64, error) {
	whereClause, args := r.buildWhereClause(filters)
	query := fmt.Sprintf("SELECT COUNT(*) FROM admins %s", whereClause)

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// Helper methods

// scanAdmin scans a single admin row
func (r *AdminRepositoryImpl) scanAdmin(row *sql.Row) (*entity.Admin, error) {
	var (
		id           int64
		username     string
		passwordHash string
		roleID       string
		status       int
		lastLogin    sql.NullTime
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := row.Scan(
		&id,
		&username,
		&passwordHash,
		&roleID,
		&status,
		&lastLogin,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert database values to domain value objects
	adminID, err := valueobject.NewAdminID(id)
	if err != nil {
		return nil, err
	}

	adminUsername, err := valueobject.NewUsername(username)
	if err != nil {
		return nil, err
	}

	adminRoleID, err := valueobject.NewRoleID(roleID)
	if err != nil {
		return nil, err
	}

	adminStatus, err := valueobject.NewAdminStatus(status)
	if err != nil {
		return nil, err
	}

	// Create admin entity
	admin := &entity.Admin{
		ID:           adminID,
		Username:     adminUsername,
		PasswordHash: passwordHash,
		Role:         adminRoleID,
		Status:       adminStatus,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	// Set last login if it's not null
	if lastLogin.Valid {
		admin.LastLogin = lastLogin.Time
	}

	return admin, nil
}

// scanAdminRow scans a row from sql.Rows
func (r *AdminRepositoryImpl) scanAdminRow(rows *sql.Rows) (*entity.Admin, error) {
	var (
		id           int64
		username     string
		passwordHash string
		roleID       string
		status       int
		lastLogin    sql.NullTime
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := rows.Scan(
		&id,
		&username,
		&passwordHash,
		&roleID,
		&status,
		&lastLogin,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Convert database values to domain value objects
	adminID, err := valueobject.NewAdminID(id)
	if err != nil {
		return nil, err
	}

	adminUsername, err := valueobject.NewUsername(username)
	if err != nil {
		return nil, err
	}

	adminRoleID, err := valueobject.NewRoleID(roleID)
	if err != nil {
		return nil, err
	}

	adminStatus, err := valueobject.NewAdminStatus(status)
	if err != nil {
		return nil, err
	}

	// Create admin entity
	admin := &entity.Admin{
		ID:           adminID,
		Username:     adminUsername,
		PasswordHash: passwordHash,
		Role:         adminRoleID,
		Status:       adminStatus,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	// Set last login if it's not null
	if lastLogin.Valid {
		admin.LastLogin = lastLogin.Time
	}

	return admin, nil
}

// buildWhereClause builds a SQL WHERE clause from a map of filters
func (r *AdminRepositoryImpl) buildWhereClause(filters map[string]interface{}) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	whereClause := "WHERE "
	args := make([]interface{}, 0, len(filters))
	first := true

	for key, value := range filters {
		if !first {
			whereClause += " AND "
		}
		first = false

		// Handle different filter keys
		switch key {
		case "username":
			whereClause += "username LIKE ?"
			args = append(args, fmt.Sprintf("%%%s%%", value))
		case "status":
			whereClause += "status = ?"
			args = append(args, value)
		case "role_id":
			whereClause += "role_id = ?"
			args = append(args, value)
		default:
			// Default to exact match
			whereClause += fmt.Sprintf("%s = ?", key)
			args = append(args, value)
		}
	}

	return whereClause, args
}
