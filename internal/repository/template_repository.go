package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wz-project/wz-backend-go/internal/domain/model"
)

// TemplateRepository 实现模板仓储接口
type TemplateRepository struct {
	db *sql.DB
}

// NewTemplateRepository 创建模板仓储实例
func NewTemplateRepository(db *sql.DB) *TemplateRepository {
	return &TemplateRepository{
		db: db,
	}
}

// FindAll 获取模板列表
func (r *TemplateRepository) FindAll(ctx context.Context, userID int64, page, pageSize int) ([]*model.Template, int, error) {
	offset := (page - 1) * pageSize

	// 查询总数
	var total int
	countQuery := "SELECT COUNT(*) FROM templates WHERE user_id = ?"
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询分页数据
	query := `
		SELECT id, user_id, name, type, preview, content, enabled, is_new, public_share, created_at, updated_at
		FROM templates 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	templates := make([]*model.Template, 0)
	for rows.Next() {
		template := &model.Template{}
		err := rows.Scan(
			&template.ID,
			&template.UserID,
			&template.Name,
			&template.Type,
			&template.Preview,
			&template.Content,
			&template.Enabled,
			&template.IsNew,
			&template.PublicShare,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// FindByID 获取单个模板
func (r *TemplateRepository) FindByID(ctx context.Context, templateID int64) (*model.Template, error) {
	query := `
		SELECT id, user_id, name, type, preview, content, enabled, is_new, public_share, created_at, updated_at
		FROM templates 
		WHERE id = ?
	`

	template := &model.Template{}
	err := r.db.QueryRowContext(ctx, query, templateID).Scan(
		&template.ID,
		&template.UserID,
		&template.Name,
		&template.Type,
		&template.Preview,
		&template.Content,
		&template.Enabled,
		&template.IsNew,
		&template.PublicShare,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("模板不存在，ID: %d", templateID)
		}
		return nil, err
	}

	return template, nil
}

// Create 创建模板
func (r *TemplateRepository) Create(ctx context.Context, template *model.Template) error {
	query := `
		INSERT INTO templates (user_id, name, type, preview, content, enabled, is_new, public_share, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now
	template.Enabled = true  // 默认启用
	template.IsNew = true    // 新创建的默认为新模板

	result, err := r.db.ExecContext(
		ctx,
		query,
		template.UserID,
		template.Name,
		template.Type,
		template.Preview,
		template.Content,
		template.Enabled,
		template.IsNew,
		template.PublicShare,
		template.CreatedAt,
		template.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	template.ID = id

	return nil
}

// Update 更新模板
func (r *TemplateRepository) Update(ctx context.Context, template *model.Template) error {
	query := `
		UPDATE templates
		SET name = ?, type = ?, preview = ?, content = ?, enabled = ?, is_new = ?, public_share = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	template.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		template.Name,
		template.Type,
		template.Preview,
		template.Content,
		template.Enabled,
		template.IsNew,
		template.PublicShare,
		template.UpdatedAt,
		template.ID,
		template.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("模板不存在或无权操作，ID: %d", template.ID)
	}

	return nil
}

// Delete 删除模板
func (r *TemplateRepository) Delete(ctx context.Context, templateID int64) error {
	query := "DELETE FROM templates WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, templateID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("模板不存在，ID: %d", templateID)
	}

	return nil
}
