package commerce

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	
	"wz-backend-go/internal/domain/commerce/entity"
	"wz-backend-go/internal/domain/commerce/repository"
	"wz-backend-go/internal/domain/commerce/service"
)

// MySQLProductRepository implements the ProductRepository interface for MySQL
type MySQLProductRepository struct {
	db *sqlx.DB
}

// NewMySQLProductRepository creates a new instance of MySQLProductRepository
func NewMySQLProductRepository(db *sqlx.DB) *MySQLProductRepository {
	return &MySQLProductRepository{
		db: db,
	}
}

// Create adds a new product to the database
func (r *MySQLProductRepository) Create(ctx context.Context, product *entity.Product) error {
	// Convert complex types to JSON
	imagesBytes, err := json.Marshal(product.Images)
	if err != nil {
		return err
	}
	
	specsBytes, err := json.Marshal(product.Specifications)
	if err != nil {
		return err
	}
	
	tagsBytes, err := json.Marshal(product.Tags)
	if err != nil {
		return err
	}
	
	query := `
		INSERT INTO commerce_products (
			id, name, description, price, original_price, stock, featured, 
			store_id, category_id, thumbnail, images, specifications, 
			tags, view_count, sales_count, is_active, created_at, updated_at
		) VALUES (
			:id, :name, :description, :price, :original_price, :stock, :featured, 
			:store_id, :category_id, :thumbnail, :images, :specifications, 
			:tags, :view_count, :sales_count, :is_active, :created_at, :updated_at
		)
	`
	
	_, err = r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":             product.ID,
		"name":           product.Name,
		"description":    product.Description,
		"price":          product.Price,
		"original_price": product.OriginalPrice,
		"stock":          product.Stock,
		"featured":       product.Featured,
		"store_id":       product.StoreID,
		"category_id":    product.CategoryID,
		"thumbnail":      product.Thumbnail,
		"images":         imagesBytes,
		"specifications": specsBytes,
		"tags":           tagsBytes,
		"view_count":     product.ViewCount,
		"sales_count":    product.SalesCount,
		"is_active":      product.IsActive,
		"created_at":     product.CreatedAt,
		"updated_at":     product.UpdatedAt,
	})
	
	return err
}

// Update updates an existing product in the database
func (r *MySQLProductRepository) Update(ctx context.Context, product *entity.Product) error {
	// Convert complex types to JSON
	imagesBytes, err := json.Marshal(product.Images)
	if err != nil {
		return err
	}
	
	specsBytes, err := json.Marshal(product.Specifications)
	if err != nil {
		return err
	}
	
	tagsBytes, err := json.Marshal(product.Tags)
	if err != nil {
		return err
	}
	
	query := `
		UPDATE commerce_products
		SET name = :name,
			description = :description,
			price = :price,
			original_price = :original_price,
			stock = :stock,
			featured = :featured,
			store_id = :store_id,
			category_id = :category_id,
			thumbnail = :thumbnail,
			images = :images,
			specifications = :specifications,
			tags = :tags,
			view_count = :view_count,
			sales_count = :sales_count,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`
	
	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":             product.ID,
		"name":           product.Name,
		"description":    product.Description,
		"price":          product.Price,
		"original_price": product.OriginalPrice,
		"stock":          product.Stock,
		"featured":       product.Featured,
		"store_id":       product.StoreID,
		"category_id":    product.CategoryID,
		"thumbnail":      product.Thumbnail,
		"images":         imagesBytes,
		"specifications": specsBytes,
		"tags":           tagsBytes,
		"view_count":     product.ViewCount,
		"sales_count":    product.SalesCount,
		"is_active":      product.IsActive,
		"updated_at":     time.Now(),
	})
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrProductNotFound
	}
	
	return nil
}

// Delete removes a product from the database
func (r *MySQLProductRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM commerce_products WHERE id = ?"
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrProductNotFound
	}
	
	return nil
}

// FindByID retrieves a product by its ID
func (r *MySQLProductRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	query := "SELECT * FROM commerce_products WHERE id = ?"
	
	var product entity.Product
	var imagesBytes, specsBytes, tagsBytes []byte
	
	row := r.db.QueryRowxContext(ctx, query, id)
	
	// Map all fields except the JSON ones
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.OriginalPrice,
		&product.Stock,
		&product.Featured,
		&product.StoreID,
		&product.CategoryID,
		&product.Thumbnail,
		&imagesBytes,
		&specsBytes,
		&tagsBytes,
		&product.ViewCount,
		&product.SalesCount,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	// Unmarshal JSON fields
	if err = json.Unmarshal(imagesBytes, &product.Images); err != nil {
		return nil, err
	}
	
	if err = json.Unmarshal(specsBytes, &product.Specifications); err != nil {
		return nil, err
	}
	
	if err = json.Unmarshal(tagsBytes, &product.Tags); err != nil {
		return nil, err
	}
	
	return &product, nil
}

// FindByStoreID retrieves all products for a given store
func (r *MySQLProductRepository) FindByStoreID(ctx context.Context, storeID string, filters repository.ProductFilters) ([]*entity.Product, error) {
	query := "SELECT * FROM commerce_products WHERE store_id = ? AND is_active = ?"
	
	if filters.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, filters.SortOrder)
	} else {
		query += " ORDER BY created_at DESC"
	}
	
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d, %d", filters.Offset, filters.Limit)
	}
	
	return r.findProductsByQuery(ctx, query, storeID, filters.ActiveOnly)
}

// FindByCategoryID retrieves all products for a given category
func (r *MySQLProductRepository) FindByCategoryID(ctx context.Context, categoryID string, filters repository.ProductFilters) ([]*entity.Product, error) {
	query := "SELECT * FROM commerce_products WHERE category_id = ? AND is_active = ?"
	
	if filters.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, filters.SortOrder)
	} else {
		query += " ORDER BY created_at DESC"
	}
	
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d, %d", filters.Offset, filters.Limit)
	}
	
	return r.findProductsByQuery(ctx, query, categoryID, filters.ActiveOnly)
}

// FindAll retrieves all products with pagination and sorting
func (r *MySQLProductRepository) FindAll(ctx context.Context, filters repository.ProductFilters) ([]*entity.Product, error) {
	query := "SELECT * FROM commerce_products"
	
	var conditions []string
	var params []interface{}
	
	if filters.ActiveOnly {
		conditions = append(conditions, "is_active = ?")
		params = append(params, true)
	}
	
	if filters.Featured {
		conditions = append(conditions, "featured = ?")
		params = append(params, true)
	}
	
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	if filters.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, filters.SortOrder)
	} else {
		query += " ORDER BY created_at DESC"
	}
	
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d, %d", filters.Offset, filters.Limit)
	}
	
	rows, err := r.db.QueryxContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanProductRows(rows)
}

// FindFeatured retrieves featured products
func (r *MySQLProductRepository) FindFeatured(ctx context.Context, limit int) ([]*entity.Product, error) {
	query := "SELECT * FROM commerce_products WHERE featured = ? AND is_active = ? ORDER BY created_at DESC LIMIT ?"
	
	rows, err := r.db.QueryxContext(ctx, query, true, true, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanProductRows(rows)
}

// FindPopular retrieves products sorted by view count
func (r *MySQLProductRepository) FindPopular(ctx context.Context, limit int) ([]*entity.Product, error) {
	query := "SELECT * FROM commerce_products WHERE is_active = ? ORDER BY view_count DESC LIMIT ?"
	
	rows, err := r.db.QueryxContext(ctx, query, true, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanProductRows(rows)
}

// FindNew retrieves the newest products
func (r *MySQLProductRepository) FindNew(ctx context.Context, limit int) ([]*entity.Product, error) {
	query := "SELECT * FROM commerce_products WHERE is_active = ? ORDER BY created_at DESC LIMIT ?"
	
	rows, err := r.db.QueryxContext(ctx, query, true, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanProductRows(rows)
}

// Search searches for products based on a query string
func (r *MySQLProductRepository) Search(ctx context.Context, query string, filters repository.ProductFilters) ([]*entity.Product, error) {
	searchQuery := `
		SELECT * FROM commerce_products 
		WHERE (name LIKE ? OR description LIKE ?) AND is_active = ?
	`
	
	// Add sorting
	if filters.SortBy != "" {
		searchQuery += fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, filters.SortOrder)
	} else {
		searchQuery += " ORDER BY created_at DESC"
	}
	
	// Add pagination
	if filters.Limit > 0 {
		searchQuery += fmt.Sprintf(" LIMIT %d, %d", filters.Offset, filters.Limit)
	}
	
	// Prepare search pattern
	searchPattern := "%" + query + "%"
	
	rows, err := r.db.QueryxContext(ctx, searchQuery, searchPattern, searchPattern, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanProductRows(rows)
}

// CountAll returns the total number of products
func (r *MySQLProductRepository) CountAll(ctx context.Context, activeOnly bool) (int, error) {
	var query string
	var count int
	
	if activeOnly {
		query = "SELECT COUNT(*) FROM commerce_products WHERE is_active = ?"
		err := r.db.QueryRowContext(ctx, query, true).Scan(&count)
		if err != nil {
			return 0, err
		}
	} else {
		query = "SELECT COUNT(*) FROM commerce_products"
		err := r.db.QueryRowContext(ctx, query).Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	
	return count, nil
}

// CountByStore returns the number of products for a given store
func (r *MySQLProductRepository) CountByStore(ctx context.Context, storeID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM commerce_products WHERE store_id = ?"
	err := r.db.QueryRowContext(ctx, query, storeID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// CountByCategory returns the number of products for a given category
func (r *MySQLProductRepository) CountByCategory(ctx context.Context, categoryID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM commerce_products WHERE category_id = ?"
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// IncrementViewCount increments the view count for a product
func (r *MySQLProductRepository) IncrementViewCount(ctx context.Context, id string) error {
	query := "UPDATE commerce_products SET view_count = view_count + 1 WHERE id = ?"
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrProductNotFound
	}
	
	return nil
}

// GenerateProductID generates a new unique product ID
func (r *MySQLProductRepository) GenerateProductID(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}

// Helper method to find products by a custom query
func (r *MySQLProductRepository) findProductsByQuery(ctx context.Context, query string, param interface{}, activeOnly bool) ([]*entity.Product, error) {
	rows, err := r.db.QueryxContext(ctx, query, param, activeOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanProductRows(rows)
}

// Helper method to scan product rows
func (r *MySQLProductRepository) scanProductRows(rows *sqlx.Rows) ([]*entity.Product, error) {
	products := []*entity.Product{}
	
	for rows.Next() {
		var product entity.Product
		var imagesBytes, specsBytes, tagsBytes []byte
		
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.OriginalPrice,
			&product.Stock,
			&product.Featured,
			&product.StoreID,
			&product.CategoryID,
			&product.Thumbnail,
			&imagesBytes,
			&specsBytes,
			&tagsBytes,
			&product.ViewCount,
			&product.SalesCount,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		// Unmarshal JSON fields
		if err = json.Unmarshal(imagesBytes, &product.Images); err != nil {
			return nil, err
		}
		
		if err = json.Unmarshal(specsBytes, &product.Specifications); err != nil {
			return nil, err
		}
		
		if err = json.Unmarshal(tagsBytes, &product.Tags); err != nil {
			return nil, err
		}
		
		products = append(products, &product)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return products, nil
}
