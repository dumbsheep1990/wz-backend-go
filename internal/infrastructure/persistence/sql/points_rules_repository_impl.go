package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/repository"
	"wz-backend-go/internal/domain/user/valueobject"
)

// SQLPointsRulesRepository 积分规则仓储的SQL实现
type SQLPointsRulesRepository struct {
	db *sqlx.DB
}

// NewPointsRulesRepository 创建积分规则仓储实例
func NewPointsRulesRepository(db *sqlx.DB) repository.PointsRulesRepository {
	return &SQLPointsRulesRepository{
		db: db,
	}
}

// PointsRulesDTO 数据库映射结构
type PointsRulesDTO struct {
	ID                int64     `db:"id"`
	SignInPoints      int       `db:"sign_in_points"`
	CommentPoints     int       `db:"comment_points"`
	SharePoints       int       `db:"share_points"`
	ArticlePoints     int       `db:"article_points"`
	InvitePoints      int       `db:"invite_points"`
	PurchaseRate      int       `db:"purchase_rate"`
	MaxDailyPoints    int       `db:"max_daily_points"`
	EnableExchange    bool      `db:"enable_exchange"`
	ExchangeRate      int       `db:"exchange_rate"`
	MinExchangePoints int       `db:"min_exchange_points"`
	TenantID          int64     `db:"tenant_id"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// Save 保存积分规则
func (r *SQLPointsRulesRepository) Save(ctx context.Context, rules *entity.PointsRules) error {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		// 如果上下文中没有事务，开始一个新事务
		var err error
		tx, err = r.db.Beginx()
		if err != nil {
			return err
		}
		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			}
		}()
	}

	// 检查是否为新记录
	if rules.ID().IsEmpty() {
		// 插入新记录
		query := `
			INSERT INTO points_rules (
				sign_in_points, comment_points, share_points, article_points, 
				invite_points, purchase_rate, max_daily_points, enable_exchange, 
				exchange_rate, min_exchange_points, tenant_id, updated_at
			) VALUES (
				:sign_in_points, :comment_points, :share_points, :article_points, 
				:invite_points, :purchase_rate, :max_daily_points, :enable_exchange, 
				:exchange_rate, :min_exchange_points, :tenant_id, :updated_at
			)
		`

		// 准备数据
		params := map[string]interface{}{
			"sign_in_points":      rules.SignInPoints(),
			"comment_points":      rules.CommentPoints(),
			"share_points":        rules.SharePoints(),
			"article_points":      rules.ArticlePoints(),
			"invite_points":       rules.InvitePoints(),
			"purchase_rate":       rules.PurchaseRate(),
			"max_daily_points":    rules.MaxDailyPoints(),
			"enable_exchange":     rules.EnableExchange(),
			"exchange_rate":       rules.ExchangeRate(),
			"min_exchange_points": rules.MinExchangePoints(),
			"tenant_id":           rules.TenantID().Value(),
			"updated_at":          rules.UpdatedAt(),
		}

		// 执行SQL
		result, err := tx.NamedExec(query, params)
		if err != nil {
			if ok {
				return err
			}
			_ = tx.Rollback()
			return err
		}

		// 获取新记录ID
		id, err := result.LastInsertId()
		if err != nil {
			if ok {
				return err
			}
			_ = tx.Rollback()
			return err
		}

		// 设置ID
		rules.SetID(valueobject.NewID(fmt.Sprintf("%d", id)))

		if !ok {
			// 如果是新创建的事务，提交
			if err := tx.Commit(); err != nil {
				return err
			}
		}

		return nil
	}

	// 更新现有记录
	query := `
		UPDATE points_rules SET
			sign_in_points = :sign_in_points,
			comment_points = :comment_points,
			share_points = :share_points,
			article_points = :article_points,
			invite_points = :invite_points,
			purchase_rate = :purchase_rate,
			max_daily_points = :max_daily_points,
			enable_exchange = :enable_exchange,
			exchange_rate = :exchange_rate,
			min_exchange_points = :min_exchange_points,
			updated_at = :updated_at
		WHERE id = :id
	`

	// 准备数据
	params := map[string]interface{}{
		"id":                  rules.ID().String(),
		"sign_in_points":      rules.SignInPoints(),
		"comment_points":      rules.CommentPoints(),
		"share_points":        rules.SharePoints(),
		"article_points":      rules.ArticlePoints(),
		"invite_points":       rules.InvitePoints(),
		"purchase_rate":       rules.PurchaseRate(),
		"max_daily_points":    rules.MaxDailyPoints(),
		"enable_exchange":     rules.EnableExchange(),
		"exchange_rate":       rules.ExchangeRate(),
		"min_exchange_points": rules.MinExchangePoints(),
		"updated_at":          time.Now(),
	}

	// 执行SQL
	_, err := tx.NamedExec(query, params)
	if err != nil {
		if ok {
			return err
		}
		_ = tx.Rollback()
		return err
	}

	if !ok {
		// 如果是新创建的事务，提交
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// FindByTenantID 根据租户ID查找积分规则
func (r *SQLPointsRulesRepository) FindByTenantID(ctx context.Context, tenantID valueobject.TenantID) (*entity.PointsRules, error) {
	query := `
		SELECT id, sign_in_points, comment_points, share_points, article_points, 
		       invite_points, purchase_rate, max_daily_points, enable_exchange, 
		       exchange_rate, min_exchange_points, tenant_id, updated_at
		FROM points_rules
		WHERE tenant_id = ?
		ORDER BY id DESC
		LIMIT 1
	`

	var dto PointsRulesDTO
	err := r.db.Get(&dto, query, tenantID.Value())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 转换为领域实体
	return r.dtoToEntity(&dto)
}

// GetUserDailyPoints 获取用户当日获取的积分总和
func (r *SQLPointsRulesRepository) GetUserDailyPoints(ctx context.Context, userID valueobject.UserID, date time.Time) (int, error) {
	// 计算日期范围
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN type = 1 AND is_revoked = 0 THEN points 
				ELSE 0 
			END
		), 0) AS daily_points
		FROM user_points
		WHERE user_id = ? 
		  AND created_at BETWEEN ? AND ?
		  AND source != 'admin'
	`

	var dailyPoints int
	err := r.db.Get(&dailyPoints, query, userID.Value(), startDate, endDate)
	if err != nil {
		return 0, err
	}

	return dailyPoints, nil
}

// 辅助方法 - 将DTO转换为实体
func (r *SQLPointsRulesRepository) dtoToEntity(dto *PointsRulesDTO) (*entity.PointsRules, error) {
	// 转换值对象
	id := valueobject.NewID(fmt.Sprintf("%d", dto.ID))
	tenantID := valueobject.NewTenantID(dto.TenantID)

	// 创建实体
	rules, err := entity.ReconstructPointsRules(
		id,
		dto.SignInPoints,
		dto.CommentPoints,
		dto.SharePoints,
		dto.ArticlePoints,
		dto.InvitePoints,
		dto.PurchaseRate,
		dto.MaxDailyPoints,
		dto.EnableExchange,
		dto.ExchangeRate,
		dto.MinExchangePoints,
		tenantID,
		dto.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return rules, nil
}
