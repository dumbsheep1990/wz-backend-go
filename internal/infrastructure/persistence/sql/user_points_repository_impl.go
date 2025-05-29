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

// SQLUserPointsRepository 用户积分仓储的SQL实现
type SQLUserPointsRepository struct {
	db *sqlx.DB
}

// NewUserPointsRepository 创建用户积分仓储实例
func NewUserPointsRepository(db *sqlx.DB) repository.UserPointsRepository {
	return &SQLUserPointsRepository{
		db: db,
	}
}

// UserPointsDTO 数据库映射结构
type UserPointsDTO struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	Points      int       `db:"points"`
	TotalPoints int       `db:"total_points"`
	Type        int       `db:"type"`
	Source      string    `db:"source"`
	Description string    `db:"description"`
	RelatedID   int64     `db:"related_id"`
	RelatedType string    `db:"related_type"`
	TenantID    int64     `db:"tenant_id"`
	OperatorID  int64     `db:"operator_id"`
	IsRevoked   bool      `db:"is_revoked"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Save 保存用户积分记录
func (r *SQLUserPointsRepository) Save(ctx context.Context, userPoints *entity.UserPoints) error {
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
	if userPoints.ID().IsEmpty() {
		// 插入新记录
		query := `
			INSERT INTO user_points (
				user_id, points, total_points, type, source, description, 
				related_id, related_type, tenant_id, operator_id, is_revoked, 
				created_at, updated_at
			) VALUES (
				:user_id, :points, :total_points, :type, :source, :description, 
				:related_id, :related_type, :tenant_id, :operator_id, :is_revoked, 
				:created_at, :updated_at
			)
		`

		// 计算用户总积分
		totalPoints, err := r.calculateTotalPoints(ctx, valueobject.NewUserID(userPoints.UserID().Value()))
		if err != nil {
			if ok {
				return err
			}
			_ = tx.Rollback()
			return err
		}

		// 根据积分类型更新总积分
		if userPoints.PointsType() == valueobject.PointsTypeIncrease {
			totalPoints += userPoints.Points().Value()
		} else {
			totalPoints -= userPoints.Points().Value()
		}

		// 准备数据
		params := map[string]interface{}{
			"user_id":      userPoints.UserID().Value(),
			"points":       userPoints.Points().Value(),
			"total_points": totalPoints,
			"type":         userPoints.PointsType().Value(),
			"source":       userPoints.Source().String(),
			"description":  userPoints.Description().String(),
			"related_id":   userPoints.RelatedID(),
			"related_type": userPoints.RelatedType().String(),
			"tenant_id":    userPoints.TenantID().Value(),
			"operator_id":  userPoints.OperatorID().Value(),
			"is_revoked":   userPoints.IsRevoked(),
			"created_at":   userPoints.CreatedAt(),
			"updated_at":   userPoints.UpdatedAt(),
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
		userPoints.SetID(valueobject.NewID(fmt.Sprintf("%d", id)))

		// 设置总积分
		userPoints.SetTotalPoints(valueobject.Points(totalPoints))

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
		UPDATE user_points SET
			user_id = :user_id,
			points = :points,
			total_points = :total_points,
			type = :type,
			source = :source,
			description = :description,
			related_id = :related_id,
			related_type = :related_type,
			tenant_id = :tenant_id,
			operator_id = :operator_id,
			is_revoked = :is_revoked,
			updated_at = :updated_at
		WHERE id = :id
	`

	// 准备数据
	params := map[string]interface{}{
		"id":           userPoints.ID().String(),
		"user_id":      userPoints.UserID().Value(),
		"points":       userPoints.Points().Value(),
		"total_points": userPoints.TotalPoints().Value(),
		"type":         userPoints.PointsType().Value(),
		"source":       userPoints.Source().String(),
		"description":  userPoints.Description().String(),
		"related_id":   userPoints.RelatedID(),
		"related_type": userPoints.RelatedType().String(),
		"tenant_id":    userPoints.TenantID().Value(),
		"operator_id":  userPoints.OperatorID().Value(),
		"is_revoked":   userPoints.IsRevoked(),
		"updated_at":   time.Now(),
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

// FindByID 根据ID查找积分记录
func (r *SQLUserPointsRepository) FindByID(ctx context.Context, id valueobject.ID) (*entity.UserPoints, error) {
	query := `
		SELECT id, user_id, points, total_points, type, source, description, 
		       related_id, related_type, tenant_id, operator_id, is_revoked, 
		       created_at, updated_at
		FROM user_points
		WHERE id = ?
	`

	var dto UserPointsDTO
	err := r.db.Get(&dto, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 转换为领域实体
	return r.dtoToEntity(&dto)
}

// FindByUserID 根据用户ID查找积分记录列表
func (r *SQLUserPointsRepository) FindByUserID(ctx context.Context, userID valueobject.UserID, offset, limit int64) ([]*entity.UserPoints, error) {
	query := `
		SELECT id, user_id, points, total_points, type, source, description, 
		       related_id, related_type, tenant_id, operator_id, is_revoked, 
		       created_at, updated_at
		FROM user_points
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	var dtos []UserPointsDTO
	err := r.db.Select(&dtos, query, userID.Value(), limit, offset)
	if err != nil {
		return nil, err
	}

	// 转换为领域实体列表
	result := make([]*entity.UserPoints, 0, len(dtos))
	for _, dto := range dtos {
		entity, err := r.dtoToEntity(&dto)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
	}

	return result, nil
}

// CountByUserID 统计用户积分记录数量
func (r *SQLUserPointsRepository) CountByUserID(ctx context.Context, userID valueobject.UserID) (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM user_points
		WHERE user_id = ?
	`

	var count int64
	err := r.db.Get(&count, query, userID.Value())
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindWithConditions 根据条件查询积分记录
func (r *SQLUserPointsRepository) FindWithConditions(ctx context.Context, conditions map[string]interface{}, offset, limit int64) ([]*entity.UserPoints, error) {
	query := `
		SELECT id, user_id, points, total_points, type, source, description, 
		       related_id, related_type, tenant_id, operator_id, is_revoked, 
		       created_at, updated_at
		FROM user_points
		WHERE 1=1
	`

	var args []interface{}
	query, args = r.buildQueryConditions(query, conditions, args)

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dtos []UserPointsDTO
	err := r.db.Select(&dtos, query, args...)
	if err != nil {
		return nil, err
	}

	// 转换为领域实体列表
	result := make([]*entity.UserPoints, 0, len(dtos))
	for _, dto := range dtos {
		entity, err := r.dtoToEntity(&dto)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
	}

	return result, nil
}

// CountWithConditions 根据条件统计积分记录数量
func (r *SQLUserPointsRepository) CountWithConditions(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM user_points
		WHERE 1=1
	`

	var args []interface{}
	query, args = r.buildQueryConditions(query, conditions, args)

	var count int64
	err := r.db.Get(&count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetTotalPointsByUserID 获取用户总积分
func (r *SQLUserPointsRepository) GetTotalPointsByUserID(ctx context.Context, userID valueobject.UserID) (int, error) {
	return r.calculateTotalPoints(ctx, userID)
}

// CountUsers 统计有积分的用户数量
func (r *SQLUserPointsRepository) CountUsers(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(DISTINCT user_id) FROM user_points"

	var count int64
	err := r.db.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SumPoints 统计总积分
func (r *SQLUserPointsRepository) SumPoints(ctx context.Context) (int64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN type = 1 AND is_revoked = 0 THEN points 
				WHEN type = 2 AND is_revoked = 0 THEN -points
				ELSE 0 
			END
		), 0) AS total_points
		FROM user_points
	`

	var sum int64
	err := r.db.Get(&sum, query)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

// MaxPoints 获取最高积分
func (r *SQLUserPointsRepository) MaxPoints(ctx context.Context) (int64, error) {
	query := `
		SELECT COALESCE(MAX(total), 0) FROM (
			SELECT user_id, SUM(
				CASE 
					WHEN type = 1 AND is_revoked = 0 THEN points 
					WHEN type = 2 AND is_revoked = 0 THEN -points
					ELSE 0 
				END
			) AS total
			FROM user_points
			GROUP BY user_id
		) AS user_totals
	`

	var max int64
	err := r.db.Get(&max, query)
	if err != nil {
		return 0, err
	}

	return max, nil
}

// SumPointsByConditions 根据条件统计积分总和
func (r *SQLUserPointsRepository) SumPointsByConditions(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN type = 1 AND is_revoked = 0 THEN points 
				WHEN type = 2 AND is_revoked = 0 THEN -points
				ELSE 0 
			END
		), 0) AS total_points
		FROM user_points
		WHERE 1=1
	`

	var args []interface{}
	query, args = r.buildQueryConditions(query, conditions, args)

	var sum int64
	err := r.db.Get(&sum, query, args...)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

// GroupBySource 按来源分组统计
func (r *SQLUserPointsRepository) GroupBySource(ctx context.Context) ([]*repository.SourceStats, error) {
	query := `
		SELECT source, COUNT(*) as count
		FROM user_points
		WHERE is_revoked = 0
		GROUP BY source
		ORDER BY count DESC
	`

	var results []struct {
		Source string `db:"source"`
		Count  int64  `db:"count"`
	}

	err := r.db.Select(&results, query)
	if err != nil {
		return nil, err
	}

	stats := make([]*repository.SourceStats, 0, len(results))
	for _, result := range results {
		stats = append(stats, &repository.SourceStats{
			Source: result.Source,
			Count:  result.Count,
		})
	}

	return stats, nil
}

// 辅助方法 - 构建查询条件
func (r *SQLUserPointsRepository) buildQueryConditions(query string, conditions map[string]interface{}, args []interface{}) (string, []interface{}) {
	if userID, ok := conditions["user_id"]; ok {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	if username, ok := conditions["username"]; ok {
		// 用户名需要先查询用户ID
		query += " AND user_id IN (SELECT id FROM users WHERE username LIKE ?)"
		args = append(args, "%"+username.(string)+"%")
	}

	if pointsType, ok := conditions["type"]; ok {
		query += " AND type = ?"
		args = append(args, pointsType)
	}

	if source, ok := conditions["source"]; ok {
		query += " AND source = ?"
		args = append(args, source)
	}

	if startDate, ok := conditions["start_date"]; ok {
		query += " AND DATE(created_at) >= ?"
		args = append(args, startDate)
	}

	if endDate, ok := conditions["end_date"]; ok {
		query += " AND DATE(created_at) <= ?"
		args = append(args, endDate)
	}

	if isRevoked, ok := conditions["is_revoked"]; ok {
		query += " AND is_revoked = ?"
		args = append(args, isRevoked)
	}

	return query, args
}

// 辅助方法 - 计算用户总积分
func (r *SQLUserPointsRepository) calculateTotalPoints(ctx context.Context, userID valueobject.UserID) (int, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN type = 1 AND is_revoked = 0 THEN points 
				WHEN type = 2 AND is_revoked = 0 THEN -points
				ELSE 0 
			END
		), 0) AS total_points
		FROM user_points
		WHERE user_id = ?
	`

	var total int
	err := r.db.Get(&total, query, userID.Value())
	if err != nil {
		return 0, err
	}

	return total, nil
}

// 辅助方法 - 将DTO转换为实体
func (r *SQLUserPointsRepository) dtoToEntity(dto *UserPointsDTO) (*entity.UserPoints, error) {
	// 转换值对象
	id := valueobject.NewID(fmt.Sprintf("%d", dto.ID))
	userID := valueobject.NewUserID(dto.UserID)

	points, err := valueobject.NewPoints(dto.Points)
	if err != nil {
		return nil, err
	}

	totalPoints := valueobject.Points(dto.TotalPoints)

	pointsType, err := valueobject.NewPointsType(dto.Type)
	if err != nil {
		return nil, err
	}

	source, err := valueobject.NewSource(dto.Source)
	if err != nil {
		return nil, err
	}

	description, err := valueobject.NewDescription(dto.Description)
	if err != nil {
		return nil, err
	}

	relatedType := valueobject.NewRelatedType(dto.RelatedType)
	operatorID := valueobject.NewUserID(dto.OperatorID)
	tenantID := valueobject.NewTenantID(dto.TenantID)

	// 创建实体
	userPoints, err := entity.ReconstructUserPoints(
		id,
		userID,
		points,
		totalPoints,
		pointsType,
		source,
		description,
		dto.RelatedID,
		relatedType,
		tenantID,
		operatorID,
		dto.IsRevoked,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return userPoints, nil
}
