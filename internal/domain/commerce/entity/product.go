package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Error declarations
var (
	ErrInsufficientStock = errors.New("insufficient stock available")
)

// Product represents a product in the commerce system
type Product struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Price        float64   `db:"price"`
	CategoryID   string    `db:"category_id"`
	StoreID      string    `db:"store_id"`
	ImageURLs    []string  `db:"-"`
	ImageURLsRaw string    `db:"image_urls"` // JSON string of image URLs
	Region       string    `db:"region"`
	District     string    `db:"district"`
	Specifications map[string]string `db:"-"`
	SpecsRaw     string    `db:"specifications"` // JSON string of specs
	Stock        int       `db:"stock"`
	ViewCount    int64     `db:"view_count"`
	SoldCount    int       `db:"sold_count"`
	IsActive     bool      `db:"is_active"`
	IsFeatured   bool      `db:"is_featured"`
	IsNew        bool      `db:"is_new"`
	Tags         []string  `db:"-"`
	TagsRaw      string    `db:"tags"` // JSON string of tags
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// NewProduct creates a new product
func NewProduct(id, name, description string, price float64, categoryID, storeID, region, district string, 
	imageURLs []string, specifications map[string]string, stock int, tags []string) *Product {
	
	now := time.Now()
	return &Product{
		ID:            id,
		Name:          name,
		Description:   description,
		Price:         price,
		CategoryID:    categoryID,
		StoreID:       storeID,
		ImageURLs:     imageURLs,
		Region:        region,
		District:      district,
		Specifications: specifications,
		Stock:         stock,
		ViewCount:     0,
		SoldCount:     0,
		IsActive:      true,
		IsNew:         true,
		IsFeatured:    false,
		Tags:          tags,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// Activate sets the product as active
func (p *Product) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// Deactivate sets the product as inactive
func (p *Product) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// MarkAsFeatured sets the product as featured
func (p *Product) MarkAsFeatured() {
	p.IsFeatured = true
	p.UpdatedAt = time.Now()
}

// UnmarkAsFeatured removes the featured status
func (p *Product) UnmarkAsFeatured() {
	p.IsFeatured = false
	p.UpdatedAt = time.Now()
}

// IncrementViewCount increments the view count
func (p *Product) IncrementViewCount() {
	p.ViewCount++
	p.UpdatedAt = time.Now()
}

// DecrementStock decreases the stock by the specified quantity
func (p *Product) DecrementStock(quantity int) error {
	if p.Stock < quantity {
		return ErrInsufficientStock
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

// UpdateDetails updates the product details
func (p *Product) UpdateDetails(name, description string, price float64, 
	imageURLs []string, specifications map[string]string, stock int, tags []string) {
	p.Name = name
	p.Description = description
	p.Price = price
	p.ImageURLs = imageURLs
	p.Specifications = specifications
	p.Stock = stock
	p.Tags = tags
	p.UpdatedAt = time.Now()
}

// MarshalJSON serializes Product to JSON
func (p *Product) MarshalJSON() ([]byte, error) {
	type Alias Product
	return json.Marshal(&struct {
		*Alias
		Price      string `json:"price"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}{
		Alias:     (*Alias)(p),
		Price:     formatPrice(p.Price),
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	})
}

// Unmarshal processes raw JSON data from the database
func (p *Product) Unmarshal() error {
	if p.ImageURLsRaw != "" {
		if err := json.Unmarshal([]byte(p.ImageURLsRaw), &p.ImageURLs); err != nil {
			return err
		}
	}

	if p.SpecsRaw != "" {
		if err := json.Unmarshal([]byte(p.SpecsRaw), &p.Specifications); err != nil {
			return err
		}
	}

	if p.TagsRaw != "" {
		if err := json.Unmarshal([]byte(p.TagsRaw), &p.Tags); err != nil {
			return err
		}
	}
	
	return nil
}

// Marshal processes entity data for database storage
func (p *Product) Marshal() error {
	imageURLsJSON, err := json.Marshal(p.ImageURLs)
	if err != nil {
		return err
	}
	p.ImageURLsRaw = string(imageURLsJSON)

	specsJSON, err := json.Marshal(p.Specifications)
	if err != nil {
		return err
	}
	p.SpecsRaw = string(specsJSON)

	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return err
	}
	p.TagsRaw = string(tagsJSON)
	
	return nil
}

// Helper function to format price
func formatPrice(price float64) string {
	return "¥" + formatFloat(price, 2)
}

// Helper function to format float with precision
func formatFloat(num float64, precision int) string {
	format := "%." + string(rune('0'+precision)) + "f"
	return fmt.Sprintf(format, num)
}
