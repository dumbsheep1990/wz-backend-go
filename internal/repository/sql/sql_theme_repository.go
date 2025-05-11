package sql

import (
	"fmt"
	"strings"
	"time"

	"wz-backend-go/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ThemeRepository 主题/模板SQL仓储实现
type ThemeRepository struct {
	conn sqlx.SqlConn
}

// NewThemeRepository 创建主题/模板仓储实现
func NewThemeRepository(conn sqlx.SqlConn) *ThemeRepository {
	return &ThemeRepository{
		conn: conn,
	}
}

// Create 创建主题/模板
func (r *ThemeRepository) Create(theme *domain.Theme) (int64, error) {
	query := `
		INSERT INTO themes (
			name, code, preview, description, status, is_default, 
			config, tenant_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	theme.CreatedAt = now
	theme.UpdatedAt = now

	// 如果设置为默认主题，需要先重置其他主题的默认状态
	if theme.IsDefault == 1 {
		err := r.resetDefaultThemes(theme.TenantID)
		if err != nil {
			return 0, err
		}
	}

	result, err := r.conn.Exec(query,
		theme.Name, theme.Code, theme.Preview, theme.Description,
		theme.Status, theme.IsDefault, theme.Config, theme.TenantID,
		theme.CreatedAt, theme.UpdatedAt,
	)
	if err != nil {
		logx.Errorf("创建主题失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取新创建主题ID失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetByID 根据ID获取主题/模板
func (r *ThemeRepository) GetByID(id int64) (*domain.Theme, error) {
	var theme domain.Theme
	query := `SELECT 
		id, name, code, preview, description, status, 
		is_default, config, tenant_id, created_at, updated_at
	FROM themes 
	WHERE id = ?`

	err := r.conn.QueryRow(&theme, query, id)
	if err != nil {
		logx.Errorf("根据ID查询主题失败: %v, id: %d", err, id)
		return nil, err
	}

	return &theme, nil
}

// List 获取主题/模板列表
func (r *ThemeRepository) List(page, pageSize int, query map[string]interface{}) ([]*domain.Theme, int64, error) {
	// 构建查询条件
	whereClause := "1=1"
	args := []interface{}{}

	// 动态添加查询条件
	if query != nil {
		if name, ok := query["name"].(string); ok && name != "" {
			whereClause += " AND name LIKE ?"
			args = append(args, fmt.Sprintf("%%%s%%", name))
		}
		if code, ok := query["code"].(string); ok && code != "" {
			whereClause += " AND code LIKE ?"
			args = append(args, fmt.Sprintf("%%%s%%", code))
		}
		if status, ok := query["status"].(int); ok {
			whereClause += " AND status = ?"
			args = append(args, status)
		}
		if isDefault, ok := query["is_default"].(int); ok {
			whereClause += " AND is_default = ?"
			args = append(args, isDefault)
		}
		if tenantID, ok := query["tenant_id"].(int64); ok && tenantID > 0 {
			whereClause += " AND tenant_id = ?"
			args = append(args, tenantID)
		}
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM themes WHERE %s", whereClause)
	var count int64
	err := r.conn.QueryRow(&count, countQuery, args...)
	if err != nil {
		logx.Errorf("查询主题总数失败: %v", err)
		return nil, 0, err
	}

	// 计算分页
	offset := (page - 1) * pageSize

	// 查询列表
	listQuery := fmt.Sprintf(`
		SELECT 
			id, name, code, preview, description, status, 
			is_default, config, tenant_id, created_at, updated_at
		FROM themes 
		WHERE %s 
		ORDER BY is_default DESC, id DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, pageSize, offset)
	var themes []*domain.Theme
	err = r.conn.QueryRows(&themes, listQuery, args...)
	if err != nil {
		logx.Errorf("查询主题列表失败: %v", err)
		return nil, 0, err
	}

	return themes, count, nil
}

// Update 更新主题/模板
func (r *ThemeRepository) Update(theme *domain.Theme) error {
	// 构建动态更新SQL
	var setValues []string
	var args []interface{}

	// 只更新非空字段
	if theme.Name != "" {
		setValues = append(setValues, "name = ?")
		args = append(args, theme.Name)
	}

	if theme.Code != "" {
		setValues = append(setValues, "code = ?")
		args = append(args, theme.Code)
	}

	if theme.Preview != "" {
		setValues = append(setValues, "preview = ?")
		args = append(args, theme.Preview)
	}

	if theme.Description != "" {
		setValues = append(setValues, "description = ?")
		args = append(args, theme.Description)
	}

	// 状态字段可以为0，需要特殊处理
	setValues = append(setValues, "status = ?")
	args = append(args, theme.Status)

	// 设置是否默认主题
	setValues = append(setValues, "is_default = ?")
	args = append(args, theme.IsDefault)

	// 如果设置为默认主题，需要先重置其他主题的默认状态
	if theme.IsDefault == 1 {
		err := r.resetDefaultThemes(theme.TenantID)
		if err != nil {
			return err
		}
	}

	if theme.Config != "" {
		setValues = append(setValues, "config = ?")
		args = append(args, theme.Config)
	}

	// 更新时间始终更新
	now := time.Now()
	setValues = append(setValues, "updated_at = ?")
	args = append(args, now)

	// 构建更新语句
	query := fmt.Sprintf("UPDATE themes SET %s WHERE id = ?", strings.Join(setValues, ", "))
	args = append(args, theme.ID)

	// 执行更新
	_, err := r.conn.Exec(query, args...)
	if err != nil {
		logx.Errorf("更新主题失败: %v, id: %d", err, theme.ID)
		return err
	}

	return nil
}

// Delete 删除主题/模板
func (r *ThemeRepository) Delete(id int64) error {
	// 先检查是否是默认主题，不允许删除默认主题
	var isDefault int
	checkQuery := "SELECT is_default FROM themes WHERE id = ?"
	err := r.conn.QueryRow(&isDefault, checkQuery, id)
	if err != nil {
		logx.Errorf("查询主题是否为默认状态失败: %v, id: %d", err, id)
		return err
	}

	if isDefault == 1 {
		return fmt.Errorf("不能删除默认主题，请先设置其他主题为默认")
	}

	query := "DELETE FROM themes WHERE id = ?"
	_, err = r.conn.Exec(query, id)
	if err != nil {
		logx.Errorf("删除主题失败: %v, id: %d", err, id)
		return err
	}

	return nil
}

// SetDefault 设置默认主题
func (r *ThemeRepository) SetDefault(id int64, tenantID int64) error {
	// 先重置所有主题的默认状态
	err := r.resetDefaultThemes(tenantID)
	if err != nil {
		return err
	}

	// 然后设置指定主题为默认
	query := "UPDATE themes SET is_default = 1, updated_at = ? WHERE id = ? AND tenant_id = ?"
	now := time.Now()
	_, err = r.conn.Exec(query, now, id, tenantID)
	if err != nil {
		logx.Errorf("设置默认主题失败: %v, id: %d", err, id)
		return err
	}

	// 更新站点配置中的主题ID
	updateSiteConfigQuery := "UPDATE site_configs SET theme_id = ?, updated_at = ? WHERE tenant_id = ?"
	_, err = r.conn.Exec(updateSiteConfigQuery, id, now, tenantID)
	if err != nil {
		logx.Warnf("更新站点配置的主题ID失败: %v, themeID: %d, tenantID: %d", err, id, tenantID)
		// 这里不返回错误，因为主题已经设置为默认，站点配置更新失败不影响核心功能
	}

	return nil
}

// GetDefault 获取默认主题
func (r *ThemeRepository) GetDefault(tenantID int64) (*domain.Theme, error) {
	var theme domain.Theme
	query := `
		SELECT 
			id, name, code, preview, description, status, 
			is_default, config, tenant_id, created_at, updated_at
		FROM themes 
		WHERE is_default = 1 AND tenant_id = ?
		LIMIT 1
	`

	err := r.conn.QueryRow(&theme, query, tenantID)
	if err != nil {
		if err == sqlx.ErrNotFound {
			// 如果找不到默认主题，尝试获取任意一个启用的主题
			fallbackQuery := `
				SELECT 
					id, name, code, preview, description, status, 
					is_default, config, tenant_id, created_at, updated_at
				FROM themes 
				WHERE status = 1 AND tenant_id = ?
				LIMIT 1
			`
			fallbackErr := r.conn.QueryRow(&theme, fallbackQuery, tenantID)
			if fallbackErr != nil {
				logx.Errorf("查询租户默认主题和备选主题都失败: %v, tenantID: %d", fallbackErr, tenantID)
				return nil, fallbackErr
			}

			// 自动将这个主题设置为默认主题
			setErr := r.SetDefault(theme.ID, tenantID)
			if setErr != nil {
				logx.Warnf("自动设置默认主题失败: %v", setErr)
				// 这里不返回错误，因为已经找到了一个主题
			}

			return &theme, nil
		}

		logx.Errorf("查询默认主题失败: %v, tenantID: %d", err, tenantID)
		return nil, err
	}

	return &theme, nil
}

// resetDefaultThemes 重置指定租户的所有主题为非默认状态
func (r *ThemeRepository) resetDefaultThemes(tenantID int64) error {
	query := "UPDATE themes SET is_default = 0, updated_at = ? WHERE tenant_id = ?"
	_, err := r.conn.Exec(query, time.Now(), tenantID)
	if err != nil {
		logx.Errorf("重置默认主题状态失败: %v, tenantID: %d", err, tenantID)
		return err
	}
	return nil
}
