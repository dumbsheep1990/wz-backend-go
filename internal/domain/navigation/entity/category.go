package entity

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a navigation category as displayed in the 万知导航 UI
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`         // 类别名称，如 "主流网站", "平台网站", etc.
	DisplayName string    `json:"display_name"` // 显示的名称
	Description string    `json:"description"`  // 类别描述
	IconURL     string    `json:"icon_url"`     // 类别图标URL
	SortOrder   int       `json:"sort_order"`   // 排序顺序
	IsActive    bool      `json:"is_active"`    // 是否激活
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCategory creates a new category entity
func NewCategory(name, displayName, description, iconURL string, sortOrder int) *Category {
	now := time.Now()
	return &Category{
		ID:          uuid.New().String(),
		Name:        name,
		DisplayName: displayName,
		Description: description,
		IconURL:     iconURL,
		SortOrder:   sortOrder,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Activate activates the category
func (c *Category) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the category
func (c *Category) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// UpdateDetails updates the category details
func (c *Category) UpdateDetails(displayName, description, iconURL string) {
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
