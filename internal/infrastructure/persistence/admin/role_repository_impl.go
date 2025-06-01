package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"wz-backend-go/internal/domain/admin/entity"
	"wz-backend-go/internal/domain/admin/repository"
	"wz-backend-go/internal/domain/admin/valueobject"
)

// RoleRepositoryImpl implements the RoleRepository interface using SQL database
type RoleRepositoryImpl struct {
	db *sql.DB
}

// NewRoleRepository creates a new RoleRepositoryImpl
func NewRoleRepository(db *sql.DB) *RoleRepositoryImpl {
	return &RoleRepositoryImpl{
		db: db,
	}
}

// Ensure RoleRepositoryImpl implements the RoleRepository interface
var _ repository.RoleRepository = (*RoleRepositoryImpl)(nil)

// FindByID finds a role by ID
func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id valueobject.RoleID) (*entity.Role, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM roles
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id.Value())

	role, err := r.scanRole(row)
	if err != nil {
		return nil, err
	}

	if role != nil {
		// Load permissions for this role
		permissions, err := r.getPermissionsForRole(ctx, id)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions
	}

	return role, nil
}

// FindByName finds a role by name
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name valueobject.RoleName) (*entity.Role, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM roles
		WHERE name = ?
	`

	row := r.db.QueryRowContext(ctx, query, name.Value())

	role, err := r.scanRole(row)
	if err != nil {
		return nil, err
	}

	if role != nil {
		// Load permissions for this role
		permissions, err := r.getPermissionsForRole(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions
	}

	return role, nil
}

// Save persists a role entity (creates or updates)
func (r *RoleRepositoryImpl) Save(ctx context.Context, role *entity.Role) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if role exists
	existingRole, err := r.FindByID(ctx, role.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingRole == nil {
		// Create new role
		query := `
			INSERT INTO roles (
				id, name, description, parent_id, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?)
		`

		_, err = tx.ExecContext(
			ctx,
			query,
			role.ID.Value(),
			role.Name.Value(),
			role.Description,
			nullableParentID(role.ParentID),
			role.CreatedAt,
			role.UpdatedAt,
		)
		if err != nil {
			return err
		}
	} else {
		// Update existing role
		query := `
			UPDATE roles
			SET name = ?, 
				description = ?, 
				parent_id = ?,
				updated_at = ?
			WHERE id = ?
		`

		_, err = tx.ExecContext(
			ctx,
			query,
			role.Name.Value(),
			role.Description,
			nullableParentID(role.ParentID),
			time.Now(),
			role.ID.Value(),
		)
		if err != nil {
			return err
		}

		// Delete all existing permissions for this role
		_, err = tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE role_id = ?", role.ID.Value())
		if err != nil {
			return err
		}
	}

	// Insert permissions
	if len(role.Permissions) > 0 {
		// Prepare bulk insert query
		valueStrings := make([]string, 0, len(role.Permissions))
		valueArgs := make([]interface{}, 0, len(role.Permissions)*2)

		for _, permission := range role.Permissions {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, role.ID.Value(), permission.Value())
		}

		query := fmt.Sprintf(
			"INSERT INTO role_permissions (role_id, permission) VALUES %s",
			strings.Join(valueStrings, ", "),
		)

		_, err = tx.ExecContext(ctx, query, valueArgs...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Delete deletes a role entity
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id valueobject.RoleID) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete all permissions for this role
	_, err = tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE role_id = ?", id.Value())
	if err != nil {
		return err
	}

	// Delete the role
	_, err = tx.ExecContext(ctx, "DELETE FROM roles WHERE id = ?", id.Value())
	if err != nil {
		return err
	}

	return tx.Commit()
}

// List lists role entities with pagination
func (r *RoleRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error) {
	// Count total roles
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM roles").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query roles with pagination
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM roles
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Scan results
	var roles []*entity.Role
	for rows.Next() {
		var (
			id          string
			name        string
			description string
			parentID    sql.NullString
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(&id, &name, &description, &parentID, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}

		// Convert database values to domain value objects
		roleID, err := valueobject.NewRoleID(id)
		if err != nil {
			return nil, 0, err
		}

		roleName, err := valueobject.NewRoleName(name)
		if err != nil {
			return nil, 0, err
		}

		// Create role entity
		role := &entity.Role{
			ID:          roleID,
			Name:        roleName,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		// Set parent ID if not null
		if parentID.Valid {
			parentRoleID, err := valueobject.NewRoleID(parentID.String)
			if err != nil {
				return nil, 0, err
			}
			role.ParentID = parentRoleID
		}

		// Load permissions for this role
		permissions, err := r.getPermissionsForRole(ctx, roleID)
		if err != nil {
			return nil, 0, err
		}
		role.Permissions = permissions

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// FindByAdminID finds roles by admin ID
func (r *RoleRepositoryImpl) FindByAdminID(ctx context.Context, adminID valueobject.AdminID) (*entity.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.parent_id, r.created_at, r.updated_at
		FROM roles r
		JOIN admins a ON r.id = a.role_id
		WHERE a.id = ?
	`

	row := r.db.QueryRowContext(ctx, query, adminID.Value())

	role, err := r.scanRole(row)
	if err != nil {
		return nil, err
	}

	if role != nil {
		// Load permissions for this role
		permissions, err := r.getPermissionsForRole(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions
	}

	return role, nil
}

// GetAllRoles gets all roles (for smaller systems where pagination is not needed)
func (r *RoleRepositoryImpl) GetAllRoles(ctx context.Context) ([]*entity.Role, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM roles
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		var (
			id          string
			name        string
			description string
			parentID    sql.NullString
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(&id, &name, &description, &parentID, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		// Convert database values to domain value objects
		roleID, err := valueobject.NewRoleID(id)
		if err != nil {
			return nil, err
		}

		roleName, err := valueobject.NewRoleName(name)
		if err != nil {
			return nil, err
		}

		// Create role entity
		role := &entity.Role{
			ID:          roleID,
			Name:        roleName,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		// Set parent ID if not null
		if parentID.Valid {
			parentRoleID, err := valueobject.NewRoleID(parentID.String)
			if err != nil {
				return nil, err
			}
			role.ParentID = parentRoleID
		}

		// Load permissions for this role
		permissions, err := r.getPermissionsForRole(ctx, roleID)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// HasPermission checks if a role has a specific permission
func (r *RoleRepositoryImpl) HasPermission(ctx context.Context, roleID valueobject.RoleID, permission valueobject.Permission) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM role_permissions 
		WHERE role_id = ? AND permission = ?
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, roleID.Value(), permission.Value()).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPermissions gets all permissions for a role
func (r *RoleRepositoryImpl) GetPermissions(ctx context.Context, roleID valueobject.RoleID) ([]valueobject.Permission, error) {
	return r.getPermissionsForRole(ctx, roleID)
}

// Helper methods

// scanRole scans a single role row
func (r *RoleRepositoryImpl) scanRole(row *sql.Row) (*entity.Role, error) {
	var (
		id          string
		name        string
		description string
		parentID    sql.NullString
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := row.Scan(&id, &name, &description, &parentID, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert database values to domain value objects
	roleID, err := valueobject.NewRoleID(id)
	if err != nil {
		return nil, err
	}

	roleName, err := valueobject.NewRoleName(name)
	if err != nil {
		return nil, err
	}

	// Create role entity
	role := &entity.Role{
		ID:          roleID,
		Name:        roleName,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Set parent ID if not null
	if parentID.Valid {
		parentRoleID, err := valueobject.NewRoleID(parentID.String)
		if err != nil {
			return nil, err
		}
		role.ParentID = parentRoleID
	}

	return role, nil
}

// getPermissionsForRole gets all permissions for a role
func (r *RoleRepositoryImpl) getPermissionsForRole(ctx context.Context, roleID valueobject.RoleID) ([]valueobject.Permission, error) {
	query := `
		SELECT permission
		FROM role_permissions
		WHERE role_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, roleID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []valueobject.Permission
	for rows.Next() {
		var permissionStr string
		if err := rows.Scan(&permissionStr); err != nil {
			return nil, err
		}

		permission, err := valueobject.NewPermission(permissionStr)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// nullableParentID handles nullable parent ID for SQL queries
func nullableParentID(parentID valueobject.RoleID) interface{} {
	if parentID.Value() == "" {
		return nil
	}
	return parentID.Value()
}
