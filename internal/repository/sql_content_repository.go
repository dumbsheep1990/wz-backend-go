package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Content 内容实体
type Content struct {
	ID           int64
	Type         string
	Title        string
	Content      string
	UserId       int64
	CategoryId   int64
	Status       int32
	ViewCount    int32
	LikeCount    int32
	CommentCount int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Category 分类实体
type Category struct {
	ID          int64
	Name        string
	Description string
	ParentId    int64
	SortOrder   int32
	Status      int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Recommendation 推荐内容实体
type Recommendation struct {
	ID         int64
	Type       string
	ResourceId int64
	SortOrder  int32
	OperatorId int64
	ExpireAt   time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ContentRepository 内容数据仓库接口
type ContentRepository interface {
	// 内容相关方法
	GetContentById(ctx context.Context, id int64) (*Content, error)
	GetContentList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Content, int64, error)
	UpdateContentStatus(ctx context.Context, id int64, status int32, reason string, operatorId int64) error
	DeleteContent(ctx context.Context, id int64) error

	// 分类相关方法
	GetCategoryById(ctx context.Context, id int64) (*Category, error)
	GetCategoryList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Category, int64, error)
	CreateCategory(ctx context.Context, category *Category) (int64, error)
	UpdateCategory(ctx context.Context, id int64, updates map[string]interface{}) error
	DeleteCategory(ctx context.Context, id int64) error

	// 推荐内容相关方法
	RecommendContent(ctx context.Context, recommendation *Recommendation) (int64, error)
	CancelRecommendation(ctx context.Context, id int64) error
	GetRecommendationList(ctx context.Context, page, pageSize int, contentType string) ([]*Recommendation, int64, error)
}

// SqlContentRepository SQL内容仓库实现
type SqlContentRepository struct {
	conn sqlx.SqlConn
}

// NewSqlContentRepository 创建内容仓库实例
func NewSqlContentRepository(conn sqlx.SqlConn) ContentRepository {
	return &SqlContentRepository{
		conn: conn,
	}
}

// GetContentById 获取指定ID的内容
func (r *SqlContentRepository) GetContentById(ctx context.Context, id int64) (*Content, error) {
	var content Content
	query := `SELECT id, type, title, content, user_id, category_id, status, 
	          view_count, like_count, comment_count, created_at, updated_at 
	          FROM contents WHERE id = ?`
	
	err := r.conn.QueryRowCtx(ctx, &content, query, id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &content, nil
}

// GetContentList 获取内容列表
func (r *SqlContentRepository) GetContentList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Content, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if filters != nil {
		if contentType, ok := filters["type"].(string); ok && contentType != "" {
			whereClause += " AND type = ?"
			args = append(args, contentType)
		}
		if status, ok := filters["status"].(int); ok && status > 0 {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if userId, ok := filters["userId"].(int64); ok && userId > 0 {
			whereClause += " AND user_id = ?"
			args = append(args, userId)
		}
		if categoryId, ok := filters["categoryId"].(int64); ok && categoryId > 0 {
			whereClause += " AND category_id = ?"
			args = append(args, categoryId)
		}
		if keyword, ok := filters["keyword"].(string); ok && keyword != "" {
			whereClause += " AND (title LIKE ? OR content LIKE ?)"
			keyword = "%" + keyword + "%"
			args = append(args, keyword, keyword)
		}
		if startTime, ok := filters["startTime"].(string); ok && startTime != "" {
			whereClause += " AND created_at >= ?"
			args = append(args, startTime)
		}
		if endTime, ok := filters["endTime"].(string); ok && endTime != "" {
			whereClause += " AND created_at <= ?"
			args = append(args, endTime)
		}
	}

	// 执行查询
	query := fmt.Sprintf(`SELECT id, type, title, content, user_id, category_id, status, 
	                     view_count, like_count, comment_count, created_at, updated_at 
	                     FROM contents %s ORDER BY id DESC LIMIT ? OFFSET ?`, whereClause)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM contents %s`, whereClause)
	
	// 添加分页参数
	queryArgs := append(args, pageSize, offset)
	
	var contents []*Content
	err := r.conn.QueryRowsCtx(ctx, &contents, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	var count int64
	err = r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return contents, count, nil
}

// UpdateContentStatus 更新内容状态
func (r *SqlContentRepository) UpdateContentStatus(ctx context.Context, id int64, status int32, reason string, operatorId int64) error {
	// 使用事务确保更新内容状态和记录状态变更历史的一致性
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 更新内容状态
	updateQuery := `UPDATE contents SET status = ?, updated_at = ? WHERE id = ?`
	_, err = tx.ExecCtx(ctx, updateQuery, status, time.Now(), id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 记录状态变更历史
	if operatorId > 0 {
		historyQuery := `INSERT INTO content_status_history (content_id, status, reason, operator_id, created_at) 
		                VALUES (?, ?, ?, ?, ?)`
		_, err = tx.ExecCtx(ctx, historyQuery, id, status, reason, operatorId, time.Now())
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeleteContent 删除内容
func (r *SqlContentRepository) DeleteContent(ctx context.Context, id int64) error {
	query := "DELETE FROM contents WHERE id = ?"
	_, err := r.conn.ExecCtx(ctx, query, id)
	return err
}

// GetCategoryById 获取指定ID的分类
func (r *SqlContentRepository) GetCategoryById(ctx context.Context, id int64) (*Category, error) {
	var category Category
	query := `SELECT id, name, description, parent_id, sort_order, status, created_at, updated_at 
	          FROM categories WHERE id = ?`
	
	err := r.conn.QueryRowCtx(ctx, &category, query, id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetCategoryList 获取分类列表
func (r *SqlContentRepository) GetCategoryList(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if filters != nil {
		if name, ok := filters["name"].(string); ok && name != "" {
			whereClause += " AND name LIKE ?"
			args = append(args, "%"+name+"%")
		}
		if status, ok := filters["status"].(int); ok && status > 0 {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if parentId, ok := filters["parentId"].(int64); ok {
			whereClause += " AND parent_id = ?"
			args = append(args, parentId)
		}
	}

	// 执行查询
	query := fmt.Sprintf(`SELECT id, name, description, parent_id, sort_order, status, created_at, updated_at 
	                     FROM categories %s ORDER BY sort_order ASC, id ASC LIMIT ? OFFSET ?`, whereClause)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM categories %s`, whereClause)
	
	// 添加分页参数
	queryArgs := append(args, pageSize, offset)
	
	var categories []*Category
	err := r.conn.QueryRowsCtx(ctx, &categories, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	var count int64
	err = r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

// CreateCategory 创建分类
func (r *SqlContentRepository) CreateCategory(ctx context.Context, category *Category) (int64, error) {
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	query := `INSERT INTO categories (name, description, parent_id, sort_order, status, created_at, updated_at) 
	         VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.conn.ExecCtx(ctx, query,
		category.Name,
		category.Description,
		category.ParentId,
		category.SortOrder,
		category.Status,
		category.CreatedAt,
		category.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateCategory 更新分类信息
func (r *SqlContentRepository) UpdateCategory(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	sets := []string{}
	args := []interface{}{}

	// 添加更新字段
	for k, v := range updates {
		// 转换键名为下划线形式
		fieldName := camelToSnake(k)
		sets = append(sets, fmt.Sprintf("%s = ?", fieldName))
		args = append(args, v)
	}

	// 添加更新时间
	sets = append(sets, "updated_at = ?")
	args = append(args, time.Now())

	// 添加ID条件
	args = append(args, id)

	query := fmt.Sprintf("UPDATE categories SET %s WHERE id = ?", strings.Join(sets, ", "))
	_, err := r.conn.ExecCtx(ctx, query, args...)
	return err
}

// DeleteCategory 删除分类
func (r *SqlContentRepository) DeleteCategory(ctx context.Context, id int64) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 检查是否有子分类
	var childCount int64
	childQuery := "SELECT COUNT(*) FROM categories WHERE parent_id = ?"
	err = tx.QueryRowCtx(ctx, &childCount, childQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if childCount > 0 {
		tx.Rollback()
		return fmt.Errorf("分类下存在子分类，无法删除")
	}

	// 检查分类下是否有内容
	var contentCount int64
	contentQuery := "SELECT COUNT(*) FROM contents WHERE category_id = ?"
	err = tx.QueryRowCtx(ctx, &contentCount, contentQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if contentCount > 0 {
		tx.Rollback()
		return fmt.Errorf("分类下存在内容，无法删除")
	}

	// 执行删除
	deleteQuery := "DELETE FROM categories WHERE id = ?"
	_, err = tx.ExecCtx(ctx, deleteQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// RecommendContent 推荐内容
func (r *SqlContentRepository) RecommendContent(ctx context.Context, recommendation *Recommendation) (int64, error) {
	now := time.Now()
	recommendation.CreatedAt = now
	recommendation.UpdatedAt = now

	// 如果未设置过期时间，默认7天后过期
	if recommendation.ExpireAt.IsZero() {
		recommendation.ExpireAt = now.AddDate(0, 0, 7)
	}

	query := `INSERT INTO recommendations (type, resource_id, sort_order, operator_id, expire_at, created_at, updated_at) 
	         VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.conn.ExecCtx(ctx, query,
		recommendation.Type,
		recommendation.ResourceId,
		recommendation.SortOrder,
		recommendation.OperatorId,
		recommendation.ExpireAt,
		recommendation.CreatedAt,
		recommendation.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CancelRecommendation 取消推荐
func (r *SqlContentRepository) CancelRecommendation(ctx context.Context, id int64) error {
	query := "DELETE FROM recommendations WHERE id = ?"
	_, err := r.conn.ExecCtx(ctx, query, id)
	return err
}

// GetRecommendationList 获取推荐内容列表
func (r *SqlContentRepository) GetRecommendationList(ctx context.Context, page, pageSize int, contentType string) ([]*Recommendation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if contentType != "" {
		whereClause += " AND type = ?"
		args = append(args, contentType)
	}
	
	// 只返回未过期的推荐
	whereClause += " AND (expire_at IS NULL OR expire_at > ?)"
	args = append(args, time.Now())

	// 执行查询
	query := fmt.Sprintf(`SELECT id, type, resource_id, sort_order, operator_id, expire_at, created_at, updated_at 
	                     FROM recommendations %s ORDER BY sort_order ASC, id DESC LIMIT ? OFFSET ?`, whereClause)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM recommendations %s`, whereClause)
	
	// 添加分页参数
	queryArgs := append(args, pageSize, offset)
	
	var recommendations []*Recommendation
	err := r.conn.QueryRowsCtx(ctx, &recommendations, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	var count int64
	err = r.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return recommendations, count, nil
} 