package entity

import (
	"time"

	"github.com/google/uuid"
)

// Website represents a website entry in the navigation system
type Website struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"category_id"` // 所属类别ID
	Name        string    `json:"name"`        // 网站名称
	URL         string    `json:"url"`         // 网站URL
	Description string    `json:"description"` // 网站描述
	IconURL     string    `json:"icon_url"`    // 网站图标URL
	SortOrder   int       `json:"sort_order"`  // 排序顺序
	IsActive    bool      `json:"is_active"`   // 是否激活
	IsNew       bool      `json:"is_new"`      // 是否新添加
	IsFeatured  bool      `json:"is_featured"` // 是否推荐
	ViewCount   int64     `json:"view_count"`  // 查看次数
	Tags        []string  `json:"tags"`        // 标签列表
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewWebsite creates a new website entity
func NewWebsite(categoryID, name, url, description, iconURL string, sortOrder int, tags []string) *Website {
	now := time.Now()
	return &Website{
		ID:          uuid.New().String(),
		CategoryID:  categoryID,
		Name:        name,
		URL:         url,
		Description: description,
		IconURL:     iconURL,
		SortOrder:   sortOrder,
		IsActive:    true,
		IsNew:       true,
		IsFeatured:  false,
		ViewCount:   0,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IncrementViewCount increments the view count
func (w *Website) IncrementViewCount() {
	w.ViewCount++
	w.UpdatedAt = time.Now()
}

// Activate activates the website
func (w *Website) Activate() {
	w.IsActive = true
	w.UpdatedAt = time.Now()
}

// Deactivate deactivates the website
func (w *Website) Deactivate() {
	w.IsActive = false
	w.UpdatedAt = time.Now()
}

// MarkAsFeatured marks the website as featured
func (w *Website) MarkAsFeatured() {
	w.IsFeatured = true
	w.UpdatedAt = time.Now()
}

// UnmarkAsFeatured unmarks the website as featured
func (w *Website) UnmarkAsFeatured() {
	w.IsFeatured = false
	w.UpdatedAt = time.Now()
}

// UpdateDetails updates the website details
func (w *Website) UpdateDetails(name, url, description, iconURL string, tags []string) {
	w.Name = name
	w.URL = url
	w.Description = description
	w.IconURL = iconURL
	w.Tags = tags
	w.UpdatedAt = time.Now()
}

// UpdateSortOrder updates the website sort order
func (w *Website) UpdateSortOrder(sortOrder int) {
	w.SortOrder = sortOrder
	w.UpdatedAt = time.Now()
}

// MarkAsOld marks the website as not new
func (w *Website) MarkAsOld() {
	w.IsNew = false
	w.UpdatedAt = time.Now()
}
