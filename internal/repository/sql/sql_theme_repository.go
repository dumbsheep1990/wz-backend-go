package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain"
)

// ThemeRepository 主题仓储SQL实现
type ThemeRepository struct {
	db *sqlx.DB
}

// NewThemeRepository 创建主题仓储实例
func NewThemeRepository(db *sqlx.DB) domain.ThemeRepository {
	return &ThemeRepository{
		db: db,
	}
}

// Create 创建主题
func (r *ThemeRepository) Create(theme *domain.Theme) (int64, error) {
	now := time.Now()
	theme.CreatedAt = now
	theme.UpdatedAt = now

	query := `INSERT INTO themes (
        name, code, preview, description, status, 
        is_default, config, tenant_id, created_at, updated_at
    ) VALUES (
        :name, :code, :preview, :description, :status, 
        :is_default, :config, :tenant_id, :created_at, :updated_at
    )`

	// 如果设置为默认主题，先重置所有主题的默认状态
	if theme.IsDefault == 1 {
		resetQuery := `UPDATE themes SET is_default = 0 WHERE tenant_id = ?`
		_, err := r.db.Exec(resetQuery, theme.TenantID)
		if err != nil {
			return 0, fmt.Errorf("重置默认主题状态失败: %w", err)
		}
	}

	result, err := r.db.NamedExec(query, theme)
	if err != nil {
		return 0, fmt.Errorf("创建主题失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}

	return id, nil
}

// GetByID 根据ID获取主题
func (r *ThemeRepository) GetByID(id int64) (*domain.Theme, error) {
	var theme domain.Theme
	query := `SELECT * FROM themes WHERE id = ?`

	err := r.db.Get(&theme, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("获取主题失败: %w", err)
	}

	return &theme, nil
}

// List 获取主题列表
func (r *ThemeRepository) List(page, pageSize int, query map[string]interface{}) ([]*domain.Theme, int64, error) {
	themes := []*domain.Theme{}
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

	if code, ok := query["code"]; ok && code != "" {
		conditions = append(conditions, "code LIKE ?")
		args = append(args, "%"+code.(string)+"%")
	}

	// 构建WHERE子句
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM themes %s", whereClause)
	err := r.db.Get(&count, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("统计主题数量失败: %w", err)
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf("SELECT * FROM themes %s ORDER BY is_default DESC, id DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, pageSize, offset)

	err = r.db.Select(&themes, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询主题列表失败: %w", err)
	}

	return themes, count, nil
}

// Update 更新主题
func (r *ThemeRepository) Update(theme *domain.Theme) error {
	theme.UpdatedAt = time.Now()

	// 如果设置为默认主题，先重置所有主题的默认状态
	if theme.IsDefault == 1 {
		resetQuery := `UPDATE themes SET is_default = 0 WHERE tenant_id = ?`
		_, err := r.db.Exec(resetQuery, theme.TenantID)
		if err != nil {
			return fmt.Errorf("重置默认主题状态失败: %w", err)
		}
	}

	query := `UPDATE themes SET 
        name = :name, 
        code = :code, 
        preview = :preview, 
        description = :description, 
        status = :status, 
        is_default = :is_default, 
        config = :config, 
        updated_at = :updated_at 
    WHERE id = :id`

	_, err := r.db.NamedExec(query, theme)
	if err != nil {
		return fmt.Errorf("更新主题失败: %w", err)
	}

	return nil
}

// Delete 删除主题
func (r *ThemeRepository) Delete(id int64) error {
	// 先检查是否是默认主题
	var isDefault int
	var tenantID int64
	err := r.db.Get(&isDefault, "SELECT is_default, tenant_id FROM themes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("检查主题状态失败: %w", err)
	}

	// 不允许删除默认主题
	if isDefault == 1 {
		return fmt.Errorf("不能删除默认主题，请先设置其他主题为默认")
	}

	query := `DELETE FROM themes WHERE id = ?`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除主题失败: %w", err)
	}

	return nil
}

// SetDefault 设置默认主题
func (r *ThemeRepository) SetDefault(id int64, tenantID int64) error {
	// 先重置所有主题的默认状态
	resetQuery := `UPDATE themes SET is_default = 0 WHERE tenant_id = ?`
	_, err := r.db.Exec(resetQuery, tenantID)
	if err != nil {
		return fmt.Errorf("重置默认主题状态失败: %w", err)
	}

	// 再设置指定主题为默认
	query := `UPDATE themes SET is_default = 1 WHERE id = ? AND tenant_id = ?`
	result, err := r.db.Exec(query, id, tenantID)
	if err != nil {
		return fmt.Errorf("设置默认主题失败: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("主题不存在或不属于指定租户")
	}

	return nil
}

// GetDefault 获取默认主题
func (r *ThemeRepository) GetDefault(tenantID int64) (*domain.Theme, error) {
	var theme domain.Theme
	query := `SELECT * FROM themes WHERE is_default = 1 AND tenant_id = ? LIMIT 1`

	err := r.db.Get(&theme, query, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有默认主题，查找第一个可用主题
			fallbackQuery := `SELECT * FROM themes WHERE status = 1 AND tenant_id = ? LIMIT 1`
			fallbackErr := r.db.Get(&theme, fallbackQuery, tenantID)
			if fallbackErr != nil {
				if fallbackErr == sql.ErrNoRows {
					return nil, nil // 没有找到任何主题
				}
				return nil, fmt.Errorf("查询可用主题失败: %w", fallbackErr)
			}
			return &theme, nil
		}
		return nil, fmt.Errorf("获取默认主题失败: %w", err)
	}

	return &theme, nil
}
