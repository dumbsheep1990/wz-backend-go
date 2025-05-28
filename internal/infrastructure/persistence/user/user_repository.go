package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yourusername/wz-backend-go/internal/application/community/service"
)

// UserRepository implements service.UserRepository interface for the community service
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindUserNameByID finds a user's name by their ID
func (r *UserRepository) FindUserNameByID(ctx context.Context, id string) (string, error) {
	var name string
	err := r.db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // User not found
		}
		return "", err
	}
	return name, nil
}

// Ensure UserRepository implements service.UserRepository
var _ service.UserRepository = (*UserRepository)(nil)
