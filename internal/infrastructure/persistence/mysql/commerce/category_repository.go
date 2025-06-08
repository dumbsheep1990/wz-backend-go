package commerce

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	
	"wz-backend-go/internal/domain/commerce/entity"
	"wz-backend-go/internal/domain/commerce/service"
)

// MySQLCategoryRepository implements the CategoryRepository interface for MySQL
type MySQLCategoryRepository struct {
	db *sqlx.DB
}

// NewMySQLCategoryRepository creates a new instance of MySQLCategoryRepository
func NewMySQLCategoryRepository(db *sqlx.DB) *MySQLCategoryRepository {
	return &MySQLCategoryRepository{
		db: db,
	}
}

// Create adds a new category to the database
func (r *MySQLCategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	query := `
		INSERT INTO commerce_categories (
			id, name, display_name, description, parent_id, icon_url,
			sort_order, level, is_active, created_at, updated_at
		) VALUES (
			:id, :name, :display_name, :description, :parent_id, :icon_url,
			:sort_order, :level, :is_active, :created_at, :updated_at
		)
	`
	
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":           category.ID,
		"name":         category.Name,
		"display_name": category.DisplayName,
		"description":  category.Description,
		"parent_id":    category.ParentID,
		"icon_url":     category.IconURL,
		"sort_order":   category.SortOrder,
		"level":        category.Level,
		"is_active":    category.IsActive,
		"created_at":   category.CreatedAt,
		"updated_at":   category.UpdatedAt,
	})
	
	return err
}

// Update updates an existing category in the database
func (r *MySQLCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	query := `
		UPDATE commerce_categories
		SET name = :name,
			display_name = :display_name,
			description = :description,
			parent_id = :parent_id,
			icon_url = :icon_url,
			sort_order = :sort_order,
			level = :level,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`
	
	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":           category.ID,
		"name":         category.Name,
		"display_name": category.DisplayName,
		"description":  category.Description,
		"parent_id":    category.ParentID,
		"icon_url":     category.IconURL,
		"sort_order":   category.SortOrder,
		"level":        category.Level,
		"is_active":    category.IsActive,
		"updated_at":   time.Now(),
	})
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrCategoryNotFound
	}
	
	return nil
}

// Delete removes a category from the database
func (r *MySQLCategoryRepository) Delete(ctx context.Context, id string) error {
	// First check if there are any child categories
	var childCount int
	childQuery := "SELECT COUNT(*) FROM commerce_categories WHERE parent_id = ?"
	err := r.db.QueryRowContext(ctx, childQuery, id).Scan(&childCount)
	if err != nil {
		return err
	}
	
	if childCount > 0 {
		return fmt.Errorf("cannot delete category with child categories")
	}
	
	// Then check if there are any products using this category
	var productCount int
	productQuery := "SELECT COUNT(*) FROM commerce_products WHERE category_id = ?"
	err = r.db.QueryRowContext(ctx, productQuery, id).Scan(&productCount)
	if err != nil {
		return err
	}
	
	if productCount > 0 {
		return fmt.Errorf("cannot delete category with associated products")
	}
	
	// If no children or products, proceed with deletion
	query := "DELETE FROM commerce_categories WHERE id = ?"
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrCategoryNotFound
	}
	
	return nil
}

// FindByID retrieves a category by its ID
func (r *MySQLCategoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	query := "SELECT * FROM commerce_categories WHERE id = ?"
	
	var category entity.Category
	err := r.db.GetContext(ctx, &category, query, id)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &category, nil
}

// FindByParentID retrieves categories by parent ID
func (r *MySQLCategoryRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Category, error) {
	query := "SELECT * FROM commerce_categories WHERE parent_id = ? ORDER BY sort_order, name"
	
	var categories []*entity.Category
	err := r.db.SelectContext(ctx, &categories, query, parentID)
	
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

// FindRootCategories retrieves all root categories (those with no parent)
func (r *MySQLCategoryRepository) FindRootCategories(ctx context.Context) ([]*entity.Category, error) {
	query := "SELECT * FROM commerce_categories WHERE parent_id IS NULL OR parent_id = '' ORDER BY sort_order, name"
	
	var categories []*entity.Category
	err := r.db.SelectContext(ctx, &categories, query)
	
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

// FindByLevel retrieves categories by level
func (r *MySQLCategoryRepository) FindByLevel(ctx context.Context, level int) ([]*entity.Category, error) {
	query := "SELECT * FROM commerce_categories WHERE level = ? ORDER BY sort_order, name"
	
	var categories []*entity.Category
	err := r.db.SelectContext(ctx, &categories, query, level)
	
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

// FindAll retrieves all categories
func (r *MySQLCategoryRepository) FindAll(ctx context.Context, activeOnly bool) ([]*entity.Category, error) {
	var query string
	var categories []*entity.Category
	
	if activeOnly {
		query = "SELECT * FROM commerce_categories WHERE is_active = ? ORDER BY sort_order, name"
		err := r.db.SelectContext(ctx, &categories, query, true)
		if err != nil {
			return nil, err
		}
	} else {
		query = "SELECT * FROM commerce_categories ORDER BY sort_order, name"
		err := r.db.SelectContext(ctx, &categories, query)
		if err != nil {
			return nil, err
		}
	}
	
	return categories, nil
}

// CountProducts returns the number of products for a given category
func (r *MySQLCategoryRepository) CountProducts(ctx context.Context, categoryID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM commerce_products WHERE category_id = ?"
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// GetNextSortOrder gets the next available sort order for categories with the same parent
func (r *MySQLCategoryRepository) GetNextSortOrder(ctx context.Context, parentID string) (int, error) {
	var query string
	var maxSortOrder int
	
	if parentID == "" {
		query = "SELECT COALESCE(MAX(sort_order), 0) FROM commerce_categories WHERE parent_id IS NULL OR parent_id = ''"
		err := r.db.QueryRowContext(ctx, query).Scan(&maxSortOrder)
		if err != nil {
			return 0, err
		}
	} else {
		query = "SELECT COALESCE(MAX(sort_order), 0) FROM commerce_categories WHERE parent_id = ?"
		err := r.db.QueryRowContext(ctx, query, parentID).Scan(&maxSortOrder)
		if err != nil {
			return 0, err
		}
	}
	
	return maxSortOrder + 1, nil
}

// UpdateParent updates the parent of a category
func (r *MySQLCategoryRepository) UpdateParent(ctx context.Context, categoryID, newParentID string, newLevel int) error {
	query := "UPDATE commerce_categories SET parent_id = ?, level = ?, updated_at = ? WHERE id = ?"
	
	var parentIDParam interface{} = sql.NullString{String: newParentID, Valid: newParentID != ""}
	if newParentID == "" {
		parentIDParam = nil
	}
	
	result, err := r.db.ExecContext(ctx, query, parentIDParam, newLevel, time.Now(), categoryID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrCategoryNotFound
	}
	
	return nil
}

// UpdateSortOrder updates the sort order of a single category
func (r *MySQLCategoryRepository) UpdateSortOrder(ctx context.Context, categoryID string, sortOrder int) error {
	query := "UPDATE commerce_categories SET sort_order = ?, updated_at = ? WHERE id = ?"
	
	result, err := r.db.ExecContext(ctx, query, sortOrder, time.Now(), categoryID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrCategoryNotFound
	}
	
	return nil
}

// BatchUpdateSortOrder updates sort orders for multiple categories in a transaction
func (r *MySQLCategoryRepository) BatchUpdateSortOrder(ctx context.Context, categoryIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	
	// Update sort orders in transaction
	for i, id := range categoryIDs {
		query := "UPDATE commerce_categories SET sort_order = ?, updated_at = ? WHERE id = ?"
		_, err = tx.ExecContext(ctx, query, i+1, time.Now(), id)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

// GenerateCategoryID generates a new unique category ID
func (r *MySQLCategoryRepository) GenerateCategoryID(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}
