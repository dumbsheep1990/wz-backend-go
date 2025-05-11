package sql

import (
	"fmt"
	"strings"
	"time"

	"wz-backend-go/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// LinkRepository 友情链接SQL仓储实现
type LinkRepository struct {
	conn sqlx.SqlConn
}

// NewLinkRepository 创建友情链接仓储实现
func NewLinkRepository(conn sqlx.SqlConn) *LinkRepository {
	return &LinkRepository{
		conn: conn,
	}
}

// Create 创建友情链接
func (r *LinkRepository) Create(link *domain.Link) (int64, error) {
	query := `
		INSERT INTO links (
			name, url, logo, sort, status, description, tenant_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now

	result, err := r.conn.Exec(query,
		link.Name, link.URL, link.Logo, link.Sort, link.Status,
		link.Description, link.TenantID, link.CreatedAt, link.UpdatedAt,
	)
	if err != nil {
		logx.Errorf("创建友情链接失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取新创建友情链接ID失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetByID 根据ID获取友情链接
func (r *LinkRepository) GetByID(id int64) (*domain.Link, error) {
	var link domain.Link
	query := `SELECT 
		id, name, url, logo, sort, status, description, tenant_id, created_at, updated_at
	FROM links 
	WHERE id = ?`

	err := r.conn.QueryRow(&link, query, id)
	if err != nil {
		logx.Errorf("根据ID查询友情链接失败: %v, id: %d", err, id)
		return nil, err
	}

	return &link, nil
}

// List 获取友情链接列表
func (r *LinkRepository) List(page, pageSize int, query map[string]interface{}) ([]*domain.Link, int64, error) {
	// 构建查询条件
	whereClause := "1=1"
	args := []interface{}{}

	// 动态添加查询条件
	if query != nil {
		if name, ok := query["name"].(string); ok && name != "" {
			whereClause += " AND name LIKE ?"
			args = append(args, fmt.Sprintf("%%%s%%", name))
		}
		if status, ok := query["status"].(int); ok {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if tenantID, ok := query["tenant_id"].(int64); ok && tenantID > 0 {
			whereClause += " AND tenant_id = ?"
			args = append(args, tenantID)
		}
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM links WHERE %s", whereClause)
	var count int64
	err := r.conn.QueryRow(&count, countQuery, args...)
	if err != nil {
		logx.Errorf("查询友情链接总数失败: %v", err)
		return nil, 0, err
	}

	// 计算分页
	offset := (page - 1) * pageSize

	// 查询列表
	listQuery := fmt.Sprintf(`
		SELECT 
			id, name, url, logo, sort, status, description, tenant_id, created_at, updated_at
		FROM links 
		WHERE %s 
		ORDER BY sort ASC, id DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)
	var links []*domain.Link
	err = r.conn.QueryRows(&links, listQuery, args...)
	if err != nil {
		logx.Errorf("查询友情链接列表失败: %v", err)
		return nil, 0, err
	}

	return links, count, nil
}

// Update 更新友情链接
func (r *LinkRepository) Update(link *domain.Link) error {
	// 构建动态更新SQL
	var setValues []string
	var args []interface{}

	// 只更新非空字段
	if link.Name != "" {
		setValues = append(setValues, "name = ?")
		args = append(args, link.Name)
	}

	if link.URL != "" {
		setValues = append(setValues, "url = ?")
		args = append(args, link.URL)
	}

	if link.Logo != "" {
		setValues = append(setValues, "logo = ?")
		args = append(args, link.Logo)
	}

	// 排序和状态字段可以为0，需要特殊处理
	setValues = append(setValues, "sort = ?")
	args = append(args, link.Sort)

	setValues = append(setValues, "status = ?")
	args = append(args, link.Status)

	if link.Description != "" {
		setValues = append(setValues, "description = ?")
		args = append(args, link.Description)
	}

	// 更新时间始终更新
	now := time.Now()
	setValues = append(setValues, "updated_at = ?")
	args = append(args, now)

	// 构建更新语句
	query := fmt.Sprintf("UPDATE links SET %s WHERE id = ?", strings.Join(setValues, ", "))
	args = append(args, link.ID)

	// 执行更新
	_, err := r.conn.Exec(query, args...)
	if err != nil {
		logx.Errorf("更新友情链接失败: %v, id: %d", err, link.ID)
		return err
	}

	return nil
}

// Delete 删除友情链接
func (r *LinkRepository) Delete(id int64) error {
	query := "DELETE FROM links WHERE id = ?"
	_, err := r.conn.Exec(query, id)
	if err != nil {
		logx.Errorf("删除友情链接失败: %v, id: %d", err, id)
		return err
	}

	return nil
}
