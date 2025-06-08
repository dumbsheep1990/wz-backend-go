package commerce

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	
	"wz-backend-go/internal/domain/commerce/entity"
	"wz-backend-go/internal/domain/commerce/repository"
	"wz-backend-go/internal/domain/commerce/service"
)

// MySQLStoreRepository implements the StoreRepository interface for MySQL
type MySQLStoreRepository struct {
	db *sqlx.DB
}

// NewMySQLStoreRepository creates a new instance of MySQLStoreRepository
func NewMySQLStoreRepository(db *sqlx.DB) *MySQLStoreRepository {
	return &MySQLStoreRepository{
		db: db,
	}
}

// Create adds a new store to the database
func (r *MySQLStoreRepository) Create(ctx context.Context, store *entity.Store) error {
	query := `
		INSERT INTO commerce_stores (
			id, name, description, owner_id, logo_url, province, city, 
			district, address, contact_name, phone, rating, is_active, 
			created_at, updated_at
		) VALUES (
			:id, :name, :description, :owner_id, :logo_url, :province, :city, 
			:district, :address, :contact_name, :phone, :rating, :is_active, 
			:created_at, :updated_at
		)
	`
	
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":           store.ID,
		"name":         store.Name,
		"description":  store.Description,
		"owner_id":     store.OwnerID,
		"logo_url":     store.LogoURL,
		"province":     store.Province,
		"city":         store.City,
		"district":     store.District,
		"address":      store.Address,
		"contact_name": store.ContactName,
		"phone":        store.Phone,
		"rating":       store.Rating,
		"is_active":    store.IsActive,
		"created_at":   store.CreatedAt,
		"updated_at":   store.UpdatedAt,
	})
	
	return err
}

// Update updates an existing store in the database
func (r *MySQLStoreRepository) Update(ctx context.Context, store *entity.Store) error {
	query := `
		UPDATE commerce_stores
		SET name = :name,
			description = :description,
			logo_url = :logo_url,
			province = :province,
			city = :city,
			district = :district,
			address = :address,
			contact_name = :contact_name,
			phone = :phone,
			rating = :rating,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`
	
	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":           store.ID,
		"name":         store.Name,
		"description":  store.Description,
		"logo_url":     store.LogoURL,
		"province":     store.Province,
		"city":         store.City,
		"district":     store.District,
		"address":      store.Address,
		"contact_name": store.ContactName,
		"phone":        store.Phone,
		"rating":       store.Rating,
		"is_active":    store.IsActive,
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
		return service.ErrStoreNotFound
	}
	
	return nil
}

// Delete removes a store from the database
func (r *MySQLStoreRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM commerce_stores WHERE id = ?"
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrStoreNotFound
	}
	
	return nil
}

// FindByID retrieves a store by its ID
func (r *MySQLStoreRepository) FindByID(ctx context.Context, id string) (*entity.Store, error) {
	query := "SELECT * FROM commerce_stores WHERE id = ?"
	
	var store entity.Store
	err := r.db.GetContext(ctx, &store, query, id)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &store, nil
}

// FindByOwner retrieves all stores for a given owner
func (r *MySQLStoreRepository) FindByOwner(ctx context.Context, ownerID string) ([]*entity.Store, error) {
	query := "SELECT * FROM commerce_stores WHERE owner_id = ? ORDER BY name"
	
	var stores []*entity.Store
	err := r.db.SelectContext(ctx, &stores, query, ownerID)
	
	if err != nil {
		return nil, err
	}
	
	return stores, nil
}

// FindAll retrieves all stores with pagination and sorting
func (r *MySQLStoreRepository) FindAll(ctx context.Context, filters repository.StoreFilters) ([]*entity.Store, error) {
	query := "SELECT * FROM commerce_stores"
	
	var conditions []string
	var params []interface{}
	
	if filters.ActiveOnly {
		conditions = append(conditions, "is_active = ?")
		params = append(params, true)
	}
	
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	if filters.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, filters.SortOrder)
	} else {
		query += " ORDER BY name ASC"
	}
	
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d, %d", filters.Offset, filters.Limit)
	}
	
	var stores []*entity.Store
	err := r.db.SelectContext(ctx, &stores, query, params...)
	
	if err != nil {
		return nil, err
	}
	
	return stores, nil
}

// FindByRegion retrieves stores from a specific region
func (r *MySQLStoreRepository) FindByRegion(ctx context.Context, province, city, district string) ([]*entity.Store, error) {
	var query string
	var params []interface{}
	
	// Build query based on which location parameters are provided
	if district != "" {
		query = "SELECT * FROM commerce_stores WHERE province = ? AND city = ? AND district = ? ORDER BY name"
		params = append(params, province, city, district)
	} else if city != "" {
		query = "SELECT * FROM commerce_stores WHERE province = ? AND city = ? ORDER BY name"
		params = append(params, province, city)
	} else {
		query = "SELECT * FROM commerce_stores WHERE province = ? ORDER BY name"
		params = append(params, province)
	}
	
	var stores []*entity.Store
	err := r.db.SelectContext(ctx, &stores, query, params...)
	
	if err != nil {
		return nil, err
	}
	
	return stores, nil
}

// Search searches for stores based on a query string
func (r *MySQLStoreRepository) Search(ctx context.Context, query string) ([]*entity.Store, error) {
	searchQuery := `
		SELECT * FROM commerce_stores 
		WHERE (name LIKE ? OR description LIKE ?)
		ORDER BY name ASC
	`
	
	// Prepare search pattern
	searchPattern := "%" + query + "%"
	
	var stores []*entity.Store
	err := r.db.SelectContext(ctx, &stores, searchQuery, searchPattern, searchPattern)
	
	if err != nil {
		return nil, err
	}
	
	return stores, nil
}

// CountAll returns the total number of stores
func (r *MySQLStoreRepository) CountAll(ctx context.Context, activeOnly bool) (int, error) {
	var query string
	var count int
	
	if activeOnly {
		query = "SELECT COUNT(*) FROM commerce_stores WHERE is_active = ?"
		err := r.db.QueryRowContext(ctx, query, true).Scan(&count)
		if err != nil {
			return 0, err
		}
	} else {
		query = "SELECT COUNT(*) FROM commerce_stores"
		err := r.db.QueryRowContext(ctx, query).Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	
	return count, nil
}

// CountProducts returns the number of products for a given store
func (r *MySQLStoreRepository) CountProducts(ctx context.Context, storeID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM commerce_products WHERE store_id = ?"
	err := r.db.QueryRowContext(ctx, query, storeID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// UpdateRating updates the rating for a store
func (r *MySQLStoreRepository) UpdateRating(ctx context.Context, storeID string, newRating float64) error {
	query := "UPDATE commerce_stores SET rating = ?, updated_at = ? WHERE id = ?"
	
	result, err := r.db.ExecContext(ctx, query, newRating, time.Now(), storeID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return service.ErrStoreNotFound
	}
	
	return nil
}

// GenerateStoreID generates a new unique store ID
func (r *MySQLStoreRepository) GenerateStoreID(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}
