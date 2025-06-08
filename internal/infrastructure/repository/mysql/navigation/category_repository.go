package navigation

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"wz-backend-go/internal/domain/navigation/entity"
	"wz-backend-go/internal/domain/navigation/repository"
)

// MySQLCategoryRepository implements the CategoryRepository interface using MySQL
type MySQLCategoryRepository struct {
	db *sqlx.DB
}

// NewMySQLCategoryRepository creates a new instance of MySQLCategoryRepository
func NewMySQLCategoryRepository(db *sqlx.DB) repository.CategoryRepository {
	return &MySQLCategoryRepository{
		db: db,
	}
}

// Save persists a category to the database
func (r *MySQLCategoryRepository) Save(ctx context.Context, category *entity.Category) error {
	query := `
		INSERT INTO navigation_categories (
			id, name, display_name, description, icon_url, 
			sort_order, is_active, created_at, updated_at
		) VALUES (
			:id, :name, :display_name, :description, :icon_url, 
			:sort_order, :is_active, :created_at, :updated_at
		)
	`
	
	params := map[string]interface{}{
		"id":           category.ID,
		"name":         category.Name,
		"display_name": category.DisplayName,
		"description":  category.Description,
		"icon_url":     category.IconURL,
		"sort_order":   category.SortOrder,
		"is_active":    category.IsActive,
		"created_at":   category.CreatedAt,
		"updated_at":   category.UpdatedAt,
	}
	
	_, err := r.db.NamedExecContext(ctx, query, params)
	return err
}

// FindByID retrieves a category by its ID
func (r *MySQLCategoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	query := `
		SELECT 
			id, name, display_name, description, icon_url, 
			sort_order, is_active, created_at, updated_at
		FROM navigation_categories
		WHERE id = ?
	`
	
	var category entity.Category
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &category, nil
}

// FindByName retrieves a category by its name
func (r *MySQLCategoryRepository) FindByName(ctx context.Context, name string) (*entity.Category, error) {
	query := `
		SELECT 
			id, name, display_name, description, icon_url, 
			sort_order, is_active, created_at, updated_at
		FROM navigation_categories
		WHERE name = ?
	`
	
	var category entity.Category
	err := r.db.GetContext(ctx, &category, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &category, nil
}

// FindAll retrieves all categories
func (r *MySQLCategoryRepository) FindAll(ctx context.Context) ([]*entity.Category, error) {
	query := `
		SELECT 
			id, name, display_name, description, icon_url, 
			sort_order, is_active, created_at, updated_at
		FROM navigation_categories
	`
	
	var categories []*entity.Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

// FindActive retrieves all active categories
func (r *MySQLCategoryRepository) FindActive(ctx context.Context) ([]*entity.Category, error) {
	query := `
		SELECT 
			id, name, display_name, description, icon_url, 
			sort_order, is_active, created_at, updated_at
		FROM navigation_categories
		WHERE is_active = true
	`
	
	var categories []*entity.Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

// Update updates an existing category
func (r *MySQLCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	query := `
		UPDATE navigation_categories
		SET 
			name = :name,
			display_name = :display_name,
			description = :description,
			icon_url = :icon_url,
			sort_order = :sort_order,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`
	
	params := map[string]interface{}{
		"id":           category.ID,
		"name":         category.Name,
		"display_name": category.DisplayName,
		"description":  category.Description,
		"icon_url":     category.IconURL,
		"sort_order":   category.SortOrder,
		"is_active":    category.IsActive,
		"updated_at":   time.Now(),
	}
	
	_, err := r.db.NamedExecContext(ctx, query, params)
	return err
}

// Delete removes a category
func (r *MySQLCategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM navigation_categories WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// FindAllSorted retrieves all categories ordered by sortOrder
func (r *MySQLCategoryRepository) FindAllSorted(ctx context.Context) ([]*entity.Category, error) {
	query := `
		SELECT 
			id, name, display_name, description, icon_url, 
			sort_order, is_active, created_at, updated_at
		FROM navigation_categories
		ORDER BY sort_order ASC
	`
	
	var categories []*entity.Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

// CountWebsites counts websites in a category
func (r *MySQLCategoryRepository) CountWebsites(ctx context.Context, categoryID string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM navigation_websites
		WHERE category_id = ?
	`
	
	var count int
	err := r.db.GetContext(ctx, &count, query, categoryID)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// GenerateID generates a new unique ID
func (r *MySQLCategoryRepository) GenerateID() string {
	return uuid.New().String()
}
