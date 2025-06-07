package entity

import (
	"time"

	"github.com/google/uuid"
)

// Category 学习分类实体
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`         // 分类名称
	Description string    `json:"description"`  // 分类描述
	Icon        string    `json:"icon"`         // 分类图标
	ParentID    *string   `json:"parentId"`     // 父级分类ID，允许为空
	Level       int       `json:"level"`        // 分类层级，从1开始
	Order       int       `json:"order"`        // 排序顺序
	CoursesCount int      `json:"coursesCount"` // 课程数量
	IsActive    bool      `json:"isActive"`     // 是否激活
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewCategory 创建新分类
func NewCategory(name string, parentID *string, level, order int) *Category {
	now := time.Now()
	return &Category{
		ID:           uuid.New().String(),
		Name:         name,
		ParentID:     parentID,
		Level:        level,
		Order:        order,
		IsActive:     true,
		CoursesCount: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Update 更新分类信息
func (c *Category) Update(name, description, icon string, order int) {
	c.Name = name
	c.Description = description
	c.Icon = icon
	c.Order = order
	c.UpdatedAt = time.Now()
}

// Activate 激活分类
func (c *Category) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Time{}
}

// Deactivate 停用分类
func (c *Category) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// IncrementCourseCount 增加课程计数
func (c *Category) IncrementCourseCount() {
	c.CoursesCount++
	c.UpdatedAt = time.Now()
}

// DecrementCourseCount 减少课程计数
func (c *Category) DecrementCourseCount() {
	if c.CoursesCount > 0 {
		c.CoursesCount--
		c.UpdatedAt = time.Now()
	}
}
