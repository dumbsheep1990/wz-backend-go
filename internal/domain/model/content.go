package model

import (
	"time"
)

// Category 内容分类
type Category struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	ParentID    int64     `db:"parent_id" json:"parent_id"`
	SortOrder   int32     `db:"sort_order" json:"sort_order"`
	Status      int32     `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Post 用户发布帖子
type Post struct {
	ID           int64     `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	Content      string    `db:"content" json:"content"`
	UserID       int64     `db:"user_id" json:"user_id"`
	CategoryID   int64     `db:"category_id" json:"category_id"`
	Status       int32     `db:"status" json:"status"`
	ViewCount    int32     `db:"view_count" json:"view_count"`
	LikeCount    int32     `db:"like_count" json:"like_count"`
	CommentCount int32     `db:"comment_count" json:"comment_count"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// Review 帖子评论
type Review struct {
	ID         int64     `db:"id" json:"id"`
	PostID     int64     `db:"post_id" json:"post_id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	Content    string    `db:"content" json:"content"`
	Status     int32     `db:"status" json:"status"`
	LikeCount  int32     `db:"like_count" json:"like_count"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// ContentStatusLog 内容状态变更日志
type ContentStatusLog struct {
	ID           int64     `db:"id" json:"id"`
	ResourceType string    `db:"resource_type" json:"resource_type"`
	ResourceID   int64     `db:"resource_id" json:"resource_id"`
	Status       int32     `db:"status" json:"status"`
	Reason       string    `db:"reason" json:"reason"`
	OperatorID   int64     `db:"operator_id" json:"operator_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// HotContent 热门/推荐内容项
type HotContent struct {
	ID           int64     `db:"id" json:"id"`
	ResourceType string    `db:"resource_type" json:"resource_type"`
	ResourceID   int64     `db:"resource_id" json:"resource_id"`
	SortOrder    int32     `db:"sort_order" json:"sort_order"`
	OperatorID   int64     `db:"operator_id" json:"operator_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// CreateCategoryRequest 创建新分类的请求
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	ParentID    int64  `json:"parent_id,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
}

// UpdateCategoryRequest 更新分类的请求
type UpdateCategoryRequest struct {
	CategoryID  int64  `json:"-"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    int64  `json:"parent_id,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
	Status      int32  `json:"status,omitempty"`
}

// CreatePostRequest 表示创建新帖子的请求
type CreatePostRequest struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	UserID     int64  `json:"-"`
	CategoryID int64  `json:"category_id" validate:"required"`
}

// UpdatePostRequest 更新帖子的请求
type UpdatePostRequest struct {
	PostID     int64  `json:"-"`
	UserID     int64  `json:"-"` // 用于授权检查
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	CategoryID int64  `json:"category_id,omitempty"`
}

// CreateReviewRequest 创建新评论的请求
type CreateReviewRequest struct {
	PostID  int64  `json:"post_id" validate:"required"`
	UserID  int64  `json:"-"`
	Content string `json:"content" validate:"required"`
}

// UpdateReviewRequest 更新评论的请求
type UpdateReviewRequest struct {
	ReviewID int64  `json:"-"`
	UserID   int64  `json:"-"` // 用于授权检查
	Content  string `json:"content" validate:"required"`
}

// UpdateContentStatusRequest 更新内容状态的请求
type UpdateContentStatusRequest struct {
	ResourceType string `json:"resource_type" validate:"required"`
	ResourceID   int64  `json:"resource_id" validate:"required"`
	Status       int32  `json:"status" validate:"required"`
	Reason       string `json:"reason,omitempty"`
	OperatorID   int64  `json:"-"`
}

// SetHotContentRequest 设置热门/推荐内容的请求
type SetHotContentRequest struct {
	ResourceType string `json:"resource_type" validate:"required"`
	ResourceID   int64  `json:"resource_id" validate:"required"`
	SortOrder    int32  `json:"sort_order,omitempty"`
	OperatorID   int64  `json:"-"`
}
