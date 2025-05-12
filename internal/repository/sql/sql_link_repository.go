package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain"
)

// LinkRepository 友情链接SQL仓储实现
type LinkRepository struct {
	db *sqlx.DB
}

// NewLinkRepository 创建友情链接仓储实现
func NewLinkRepository(db *sqlx.DB) domain.LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

// Create 创建友情链接
func (r *LinkRepository) Create(link *domain.Link) (int64, error) {
	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now

	query := `INSERT INTO links (name, url, logo, sort, status, description, tenant_id, created_at, updated_at) 
              VALUES (:name, :url, :logo, :sort, :status, :description, :tenant_id, :created_at, :updated_at)`

	result, err := r.db.NamedExec(query, link)
	if err != nil {
		return 0, fmt.Errorf("创建友情链接失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}

	return id, nil
}

// GetByID 根据ID获取友情链接
func (r *LinkRepository) GetByID(id int64) (*domain.Link, error) {
	var link domain.Link
	query := `SELECT * FROM links WHERE id = ?`

	err := r.db.Get(&link, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("获取友情链接失败: %w", err)
	}

	return &link, nil
}

// List 获取友情链接列表
func (r *LinkRepository) List(page, pageSize int, query map[string]interface{}) ([]*domain.Link, int64, error) {
	links := []*domain.Link{}
	var count int64

	// 构建查询条件
	conditions := []string{}
	args := []interface{}{}

	if tenantID, ok := query["tenant_id"]; ok {
		conditions = append(conditions, "tenant_id = ?")
		args = append(args, tenantID)
	}

	if status, ok := query["status"]; ok {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	if name, ok := query["name"]; ok && name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+name.(string)+"%")
	}

	// 构建WHERE子句
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM links %s", whereClause)
	err := r.db.Get(&count, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("统计友情链接数量失败: %w", err)
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf("SELECT * FROM links %s ORDER BY sort ASC, id DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, pageSize, offset)

	err = r.db.Select(&links, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询友情链接列表失败: %w", err)
	}

	return links, count, nil
}

// Update 更新友情链接
func (r *LinkRepository) Update(link *domain.Link) error {
	link.UpdatedAt = time.Now()

	query := `UPDATE links SET 
                name = :name, 
                url = :url, 
                logo = :logo, 
                sort = :sort, 
                status = :status, 
                description = :description, 
                updated_at = :updated_at 
              WHERE id = :id`

	_, err := r.db.NamedExec(query, link)
	if err != nil {
		return fmt.Errorf("更新友情链接失败: %w", err)
	}

	return nil
}

// Delete 删除友情链接
func (r *LinkRepository) Delete(id int64) error {
	query := `DELETE FROM links WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除友情链接失败: %w", err)
	}

	return nil
}
