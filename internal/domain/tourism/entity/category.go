package entity

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a tourism category
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    string    `json:"parent_id"` // For hierarchical categories
	Icon        string    `json:"icon"`      // URL or path to category icon
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCategory creates a new category entity
func NewCategory(name, description, parentID, icon string, sortOrder int) *Category {
	now := time.Now()
	return &Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		ParentID:    parentID,
		Icon:        icon,
		SortOrder:   sortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the category information
func (c *Category) Update(name, description, parentID, icon string, sortOrder int) {
	c.Name = name
	c.Description = description
	c.ParentID = parentID
	c.Icon = icon
	c.SortOrder = sortOrder
	c.UpdatedAt = time.Now()
}
