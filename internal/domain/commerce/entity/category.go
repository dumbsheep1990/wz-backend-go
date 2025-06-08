package entity

import (
	"time"
)

// Category represents a product category in the commerce system
type Category struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	DisplayName string    `db:"display_name"`
	Description string    `db:"description"`
	ParentID    string    `db:"parent_id"`
	IconURL     string    `db:"icon_url"`
	SortOrder   int       `db:"sort_order"`
	Level       int       `db:"level"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// NewCategory creates a new product category
func NewCategory(id, name, displayName, description, parentID, iconURL string, sortOrder, level int) *Category {
	now := time.Now()
	return &Category{
		ID:          id,
		Name:        name,
		DisplayName: displayName,
		Description: description,
		ParentID:    parentID,
		IconURL:     iconURL,
		SortOrder:   sortOrder,
		Level:       level,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Activate sets the category as active
func (c *Category) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate sets the category as inactive
func (c *Category) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// UpdateDetails updates category details
func (c *Category) UpdateDetails(name, displayName, description, iconURL string) {
	c.Name = name
	c.DisplayName = displayName
	c.Description = description
	c.IconURL = iconURL
	c.UpdatedAt = time.Now()
}

// UpdateSortOrder updates the category sort order
func (c *Category) UpdateSortOrder(sortOrder int) {
	c.SortOrder = sortOrder
	c.UpdatedAt = time.Now()
}

// UpdateParent updates the parent category
func (c *Category) UpdateParent(parentID string, level int) {
	c.ParentID = parentID
	c.Level = level
	c.UpdatedAt = time.Now()
}
