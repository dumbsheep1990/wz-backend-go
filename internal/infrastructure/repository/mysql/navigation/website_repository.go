package navigation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"wz-backend-go/internal/domain/navigation/entity"
	"wz-backend-go/internal/domain/navigation/repository"
)

// MySQLWebsiteRepository implements the WebsiteRepository interface using MySQL
type MySQLWebsiteRepository struct {
	db *sqlx.DB
}

// NewMySQLWebsiteRepository creates a new instance of MySQLWebsiteRepository
func NewMySQLWebsiteRepository(db *sqlx.DB) repository.WebsiteRepository {
	return &MySQLWebsiteRepository{
		db: db,
	}
}

// Save persists a website to the database
func (r *MySQLWebsiteRepository) Save(ctx context.Context, website *entity.Website) error {
	tags, err := json.Marshal(website.Tags)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO navigation_websites (
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		) VALUES (
			:id, :category_id, :name, :url, :description, :icon_url, 
			:sort_order, :is_active, :is_new, :is_featured, :view_count,
			:tags, :created_at, :updated_at
		)
	`
	
	params := map[string]interface{}{
		"id":          website.ID,
		"category_id": website.CategoryID,
		"name":        website.Name,
		"url":         website.URL,
		"description": website.Description,
		"icon_url":    website.IconURL,
		"sort_order":  website.SortOrder,
		"is_active":   website.IsActive,
		"is_new":      website.IsNew,
		"is_featured": website.IsFeatured,
		"view_count":  website.ViewCount,
		"tags":        tags,
		"created_at":  website.CreatedAt,
		"updated_at":  website.UpdatedAt,
	}
	
	_, err = r.db.NamedExecContext(ctx, query, params)
	return err
}

// FindByID retrieves a website by its ID
func (r *MySQLWebsiteRepository) FindByID(ctx context.Context, id string) (*entity.Website, error) {
	query := `
		SELECT 
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		FROM navigation_websites
		WHERE id = ?
	`
	
	var result struct {
		entity.Website
		Tags []byte `db:"tags"`
	}
	
	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	var tags []string
	if err := json.Unmarshal(result.Tags, &tags); err != nil {
		return nil, err
	}
	
	website := result.Website
	website.Tags = tags
	
	return &website, nil
}

// FindByURL retrieves a website by its URL
func (r *MySQLWebsiteRepository) FindByURL(ctx context.Context, url string) (*entity.Website, error) {
	query := `
		SELECT 
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		FROM navigation_websites
		WHERE url = ?
	`
	
	var result struct {
		entity.Website
		Tags []byte `db:"tags"`
	}
	
	err := r.db.GetContext(ctx, &result, query, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	var tags []string
	if err := json.Unmarshal(result.Tags, &tags); err != nil {
		return nil, err
	}
	
	website := result.Website
	website.Tags = tags
	
	return &website, nil
}

// FindByCategory retrieves all websites in a category
func (r *MySQLWebsiteRepository) FindByCategory(ctx context.Context, categoryID string) ([]*entity.Website, error) {
	query := `
		SELECT 
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		FROM navigation_websites
		WHERE category_id = ?
	`
	
	var results []struct {
		entity.Website
		Tags []byte `db:"tags"`
	}
	
	err := r.db.SelectContext(ctx, &results, query, categoryID)
	if err != nil {
		return nil, err
	}
	
	websites := make([]*entity.Website, 0, len(results))
	for _, result := range results {
		var tags []string
		if err := json.Unmarshal(result.Tags, &tags); err != nil {
			return nil, err
		}
		
		website := result.Website
		website.Tags = tags
		websites = append(websites, &website)
	}
	
	return websites, nil
}

// FindByCategorySorted retrieves all websites in a category ordered by sortOrder
func (r *MySQLWebsiteRepository) FindByCategorySorted(ctx context.Context, categoryID string) ([]*entity.Website, error) {
	query := `
		SELECT 
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		FROM navigation_websites
		WHERE category_id = ?
		ORDER BY sort_order ASC
	`
	
	var results []struct {
		entity.Website
		Tags []byte `db:"tags"`
	}
	
	err := r.db.SelectContext(ctx, &results, query, categoryID)
	if err != nil {
		return nil, err
	}
	
	websites := make([]*entity.Website, 0, len(results))
	for _, result := range results {
		var tags []string
		if err := json.Unmarshal(result.Tags, &tags); err != nil {
			return nil, err
		}
		
		website := result.Website
		website.Tags = tags
		websites = append(websites, &website)
	}
	
	return websites, nil
}

// FindAll retrieves all websites
func (r *MySQLWebsiteRepository) FindAll(ctx context.Context) ([]*entity.Website, error) {
	query := `
		SELECT 
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		FROM navigation_websites
	`
	
	var results []struct {
		entity.Website
		Tags []byte `db:"tags"`
	}
	
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, err
	}
	
	websites := make([]*entity.Website, 0, len(results))
	for _, result := range results {
		var tags []string
		if err := json.Unmarshal(result.Tags, &tags); err != nil {
			return nil, err
		}
		
		website := result.Website
		website.Tags = tags
		websites = append(websites, &website)
	}
	
	return websites, nil
}

// FindFeatured retrieves all featured websites
func (r *MySQLWebsiteRepository) FindFeatured(ctx context.Context) ([]*entity.Website, error) {
	query := `
		SELECT 
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		FROM navigation_websites
		WHERE is_featured = true AND is_active = true
	`
	
	var results []struct {
		entity.Website
		Tags []byte `db:"tags"`
	}
	
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, err
	}
	
	websites := make([]*entity.Website, 0, len(results))
	for _, result := range results {
		var tags []string
		if err := json.Unmarshal(result.Tags, &tags); err != nil {
			return nil, err
		}
		
		website := result.Website
		website.Tags = tags
		websites = append(websites, &website)
	}
	
	return websites, nil
}

// FindPopular retrieves popular websites by view count
func (r *MySQLWebsiteRepository) FindPopular(ctx context.Context, limit int) ([]*entity.Website, error) {
	query := `
		SELECT 
			id, category_id, name, url, description, icon_url, 
			sort_order, is_active, is_new, is_featured, view_count,
			tags, created_at, updated_at
		FROM navigation_websites
		WHERE is_active = true
		ORDER BY view_count DESC
		LIMIT ?
	`
	
	var results []struct {
		entity.Website
		Tags []byte `db:"tags"`
	}
	
	err := r.db.SelectContext(ctx, &results, query, limit)
	if err != nil {
		return nil, err
	}
	
	websites := make([]*entity.Website, 0, len(results))
	for _, result := range results {
		var tags []string
		if err := json.Unmarshal(result.Tags, &tags); err != nil {
			return nil, err
		}
		
		website := result.Website
		website.Tags = tags
		websites = append(websites, &website)
	}
	
	return websites, nil
}

// FindByTags finds websites by tags
func (r *MySQLWebsiteRepository) FindByTags(ctx context.Context, tags []string) ([]*entity.Website, error) {
	// In a real implementation, this would use a more efficient query
	// For simplicity, we'll fetch all and filter in memory
	websites, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(tags) == 0 {
		return websites, nil
	}
	
	// Create a map for faster lookups
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag] = true
	}
	
	// Filter websites that have any of the requested tags
	var filtered []*entity.Website
	for _, website := range websites {
		for _, tag := range website.Tags {
			if tagMap[tag] {
				filtered = append(filtered, website)
				break
			}
		}
	}
	
	return filtered, nil
}

// Update updates an existing website
func (r *MySQLWebsiteRepository) Update(ctx context.Context, website *entity.Website) error {
	tags, err := json.Marshal(website.Tags)
	if err != nil {
		return err
	}

	query := `
		UPDATE navigation_websites
		SET 
			category_id = :category_id,
			name = :name,
			url = :url,
			description = :description,
			icon_url = :icon_url,
			sort_order = :sort_order,
			is_active = :is_active,
			is_new = :is_new,
			is_featured = :is_featured,
			view_count = :view_count,
			tags = :tags,
			updated_at = :updated_at
		WHERE id = :id
	`
	
	params := map[string]interface{}{
		"id":          website.ID,
		"category_id": website.CategoryID,
		"name":        website.Name,
		"url":         website.URL,
		"description": website.Description,
		"icon_url":    website.IconURL,
		"sort_order":  website.SortOrder,
		"is_active":   website.IsActive,
		"is_new":      website.IsNew,
		"is_featured": website.IsFeatured,
		"view_count":  website.ViewCount,
		"tags":        tags,
		"updated_at":  time.Now(),
	}
	
	_, err = r.db.NamedExecContext(ctx, query, params)
	return err
}

// Delete removes a website
func (r *MySQLWebsiteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM navigation_websites WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// IncrementViewCount increments the view count for a website
func (r *MySQLWebsiteRepository) IncrementViewCount(ctx context.Context, id string) error {
	query := `
		UPDATE navigation_websites
		SET 
			view_count = view_count + 1,
			updated_at = ?
		WHERE id = ?
	`
	
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// CountByCategory counts websites in a category
func (r *MySQLWebsiteRepository) CountByCategory(ctx context.Context, categoryID string) (int, error) {
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
