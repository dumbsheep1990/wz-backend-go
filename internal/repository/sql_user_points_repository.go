package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"wz-backend-go/internal/domain"
)

// SQLUserPointsRepository 用户积分的SQL实现
type SQLUserPointsRepository struct {
	db *sqlx.DB
}

// NewSQLUserPointsRepository 创建用户积分仓储实例
func NewSQLUserPointsRepository(db *sqlx.DB) domain.UserPointsRepository {
	return &SQLUserPointsRepository{
		db: db,
	}
}

// GetByID 根据ID获取积分记录
func (r *SQLUserPointsRepository) GetByID(id int64) (*domain.UserPoints, error) {
	var point domain.UserPoints
	query := `SELECT id, user_id, points, points_type, points_source, description, 
              created_at, updated_at, operator_id, is_revoked 
              FROM user_points WHERE id = ? AND deleted_at IS NULL`
	err := r.db.Get(&point, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "积分记录不存在")
		}
		return nil, errors.Wrap(err, "获取积分记录失败")
	}
	return &point, nil
}

// Create 创建积分记录
func (r *SQLUserPointsRepository) Create(points *domain.UserPoints) error {
	query := `INSERT INTO user_points (user_id, points, points_type, points_source, 
              description, operator_id, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	points.CreatedAt = now
	points.UpdatedAt = now
	
	result, err := r.db.Exec(query, points.UserID, points.Points, points.PointsType, 
		points.PointsSource, points.Description, points.OperatorID, points.CreatedAt, points.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "创建积分记录失败")
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "获取积分记录ID失败")
	}
	
	points.ID = id
	return nil
}

// MarkAsRevoked 标记积分记录为已撤销
func (r *SQLUserPointsRepository) MarkAsRevoked(id int64) error {
	query := `UPDATE user_points SET is_revoked = true, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
	
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return errors.Wrap(err, "撤销积分记录失败")
	}
	
	return nil
}

// ListWithConditions 根据条件查询积分记录列表
func (r *SQLUserPointsRepository) ListWithConditions(conditions map[string]interface{}, offset, limit int64) ([]*domain.UserPoints, error) {
	var points []*domain.UserPoints
	
	query := `SELECT p.id, p.user_id, p.points, p.points_type, p.points_source, p.description, 
              p.created_at, p.updated_at, p.operator_id, p.is_revoked 
              FROM user_points p 
              WHERE p.deleted_at IS NULL`
	
	// 构建条件查询
	args := []interface{}{}
	if len(conditions) > 0 {
		condStrs := []string{}
		
		if v, ok := conditions["user_id"]; ok {
			condStrs = append(condStrs, "p.user_id = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["points_type"]; ok {
			condStrs = append(condStrs, "p.points_type = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["points_source"]; ok {
			condStrs = append(condStrs, "p.points_source = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["is_revoked"]; ok {
			condStrs = append(condStrs, "p.is_revoked = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["start_time"]; ok {
			condStrs = append(condStrs, "p.created_at >= ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["end_time"]; ok {
			condStrs = append(condStrs, "p.created_at <= ?")
			args = append(args, v)
		}
		
		if len(condStrs) > 0 {
			query += " AND " + strings.Join(condStrs, " AND ")
		}
	}
	
	// 添加排序和分页
	query += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	err := r.db.Select(&points, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "查询积分记录列表失败")
	}
	
	return points, nil
}

// CountWithConditions 根据条件统计积分记录数量
func (r *SQLUserPointsRepository) CountWithConditions(conditions map[string]interface{}) (int64, error) {
	var count int64
	
	query := `SELECT COUNT(*) FROM user_points p WHERE p.deleted_at IS NULL`
	
	// 构建条件查询
	args := []interface{}{}
	if len(conditions) > 0 {
		condStrs := []string{}
		
		if v, ok := conditions["user_id"]; ok {
			condStrs = append(condStrs, "p.user_id = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["points_type"]; ok {
			condStrs = append(condStrs, "p.points_type = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["points_source"]; ok {
			condStrs = append(condStrs, "p.points_source = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["is_revoked"]; ok {
			condStrs = append(condStrs, "p.is_revoked = ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["start_time"]; ok {
			condStrs = append(condStrs, "p.created_at >= ?")
			args = append(args, v)
		}
		
		if v, ok := conditions["end_time"]; ok {
			condStrs = append(condStrs, "p.created_at <= ?")
			args = append(args, v)
		}
		
		if len(condStrs) > 0 {
			query += " AND " + strings.Join(condStrs, " AND ")
		}
	}
	
	err := r.db.Get(&count, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "统计积分记录数量失败")
	}
	
	return count, nil
}

// CountUsers 统计有积分的用户数量
func (r *SQLUserPointsRepository) CountUsers() (int64, error) {
	var count int64
	
	query := `SELECT COUNT(DISTINCT user_id) FROM user_points WHERE deleted_at IS NULL`
	
	err := r.db.Get(&count, query)
	if err != nil {
		return 0, errors.Wrap(err, "统计用户数量失败")
	}
	
	return count, nil
}

// SumPoints 统计总积分
func (r *SQLUserPointsRepository) SumPoints() (int64, error) {
	var sum int64
	
	query := `SELECT COALESCE(SUM(
                CASE WHEN is_revoked = 0 THEN points ELSE 0 END
              ), 0) AS total_points 
              FROM user_points 
              WHERE deleted_at IS NULL`
	
	err := r.db.Get(&sum, query)
	if err != nil {
		return 0, errors.Wrap(err, "统计总积分失败")
	}
	
	return sum, nil
}

// GetUserPointsTotal 获取用户总积分
func (r *SQLUserPointsRepository) GetUserPointsTotal(userID int64) (int64, error) {
	var total int64
	
	query := `SELECT COALESCE(SUM(
                CASE WHEN is_revoked = 0 THEN points ELSE 0 END
              ), 0) AS total_points 
              FROM user_points 
              WHERE user_id = ? AND deleted_at IS NULL`
	
	err := r.db.Get(&total, query, userID)
	if err != nil {
		return 0, errors.Wrap(err, "获取用户总积分失败")
	}
	
	return total, nil
}

// GetPointsRules 获取积分规则
func (r *SQLUserPointsRepository) GetPointsRules() (*domain.PointsRules, error) {
	var rules domain.PointsRules
	
	query := `SELECT id, sign_in_points, comment_points, share_points, 
              purchase_points_ratio, login_points, created_at, updated_at 
              FROM points_rules ORDER BY id DESC LIMIT 1`
	
	err := r.db.Get(&rules, query)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有规则，返回默认规则
			return &domain.PointsRules{
				SignInPoints:      5,
				CommentPoints:     2,
				SharePoints:       3,
				PurchasePointsRatio: 0.01,
				LoginPoints:       1,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}, nil
		}
		return nil, errors.Wrap(err, "获取积分规则失败")
	}
	
	return &rules, nil
}

// UpdatePointsRules 更新积分规则
func (r *SQLUserPointsRepository) UpdatePointsRules(rules *domain.PointsRules) error {
	if rules.ID > 0 {
		// 更新现有规则
		query := `UPDATE points_rules SET 
                  sign_in_points = ?, 
                  comment_points = ?, 
                  share_points = ?, 
                  purchase_points_ratio = ?, 
                  login_points = ?, 
                  updated_at = ? 
                  WHERE id = ?`
				  
		rules.UpdatedAt = time.Now()
		
		_, err := r.db.Exec(query, rules.SignInPoints, rules.CommentPoints, 
			rules.SharePoints, rules.PurchasePointsRatio, rules.LoginPoints, 
			rules.UpdatedAt, rules.ID)
		if err != nil {
			return errors.Wrap(err, "更新积分规则失败")
		}
	} else {
		// 创建新规则
		query := `INSERT INTO points_rules (
                  sign_in_points, comment_points, share_points, 
                  purchase_points_ratio, login_points, created_at, updated_at
                  ) VALUES (?, ?, ?, ?, ?, ?, ?)`
				  
		now := time.Now()
		rules.CreatedAt = now
		rules.UpdatedAt = now
		
		result, err := r.db.Exec(query, rules.SignInPoints, rules.CommentPoints, 
			rules.SharePoints, rules.PurchasePointsRatio, rules.LoginPoints, 
			rules.CreatedAt, rules.UpdatedAt)
		if err != nil {
			return errors.Wrap(err, "创建积分规则失败")
		}
		
		id, err := result.LastInsertId()
		if err != nil {
			return errors.Wrap(err, "获取积分规则ID失败")
		}
		
		rules.ID = id
	}
	
	return nil
} 