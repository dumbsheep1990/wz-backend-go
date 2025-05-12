package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain"
)

// UserPointsRepository 用户积分仓储SQL实现
type UserPointsRepository struct {
	db *sqlx.DB
}

// NewUserPointsRepository 创建用户积分仓储实例
func NewUserPointsRepository(db *sqlx.DB) domain.UserPointsRepository {
	return &UserPointsRepository{
		db: db,
	}
}

// Create 创建用户积分记录
func (r *UserPointsRepository) Create(points *domain.UserPoints) (int64, error) {
	// 开启事务
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("开启事务失败: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 获取用户当前总积分
	var currentTotal int
	query := `SELECT COALESCE(SUM(points), 0) FROM user_points WHERE user_id = ? AND type = 1`
	err = tx.Get(&currentTotal, query, points.UserID)
	if err != nil {
		return 0, fmt.Errorf("获取用户积分失败: %w", err)
	}

	// 获取用户当前总减少积分
	var currentDecrease int
	query = `SELECT COALESCE(SUM(points), 0) FROM user_points WHERE user_id = ? AND type = 2`
	err = tx.Get(&currentDecrease, query, points.UserID)
	if err != nil {
		return 0, fmt.Errorf("获取用户减少积分失败: %w", err)
	}

	// 计算新的总积分
	newTotal := currentTotal - currentDecrease
	if points.Type == 1 {
		// 增加积分
		newTotal += points.Points
	} else if points.Type == 2 {
		// 减少积分
		newTotal -= points.Points
		if newTotal < 0 {
			return 0, fmt.Errorf("用户积分不足")
		}
	}

	// 设置记录创建时间和总积分
	now := time.Now()
	points.CreatedAt = now
	points.UpdatedAt = now
	points.TotalPoints = newTotal

	// 插入积分记录
	insertQuery := `INSERT INTO user_points (
        user_id, points, total_points, type, source, 
        description, related_id, related_type, tenant_id, 
        created_at, updated_at
    ) VALUES (
        :user_id, :points, :total_points, :type, :source, 
        :description, :related_id, :related_type, :tenant_id, 
        :created_at, :updated_at
    )`

	result, err := tx.NamedExec(insertQuery, points)
	if err != nil {
		return 0, fmt.Errorf("创建积分记录失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("提交事务失败: %w", err)
	}

	return id, nil
}

// GetByID 根据ID获取积分记录
func (r *UserPointsRepository) GetByID(id int64) (*domain.UserPoints, error) {
	var points domain.UserPoints
	query := `SELECT * FROM user_points WHERE id = ?`

	err := r.db.Get(&points, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("获取积分记录失败: %w", err)
	}

	return &points, nil
}

// ListByUser 获取用户积分记录列表
func (r *UserPointsRepository) ListByUser(userID int64, page, pageSize int) ([]*domain.UserPoints, int64, error) {
	pointsList := []*domain.UserPoints{}
	var count int64

	// 查询总数
	countQuery := `SELECT COUNT(*) FROM user_points WHERE user_id = ?`
	err := r.db.Get(&count, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("统计积分记录数量失败: %w", err)
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listQuery := `SELECT * FROM user_points WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`
	err = r.db.Select(&pointsList, listQuery, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询积分记录列表失败: %w", err)
	}

	return pointsList, count, nil
}

// GetUserTotalPoints 获取用户总积分
func (r *UserPointsRepository) GetUserTotalPoints(userID int64) (int, error) {
	// 获取用户总增加积分
	var totalIncrease int
	query := `SELECT COALESCE(SUM(points), 0) FROM user_points WHERE user_id = ? AND type = 1`
	err := r.db.Get(&totalIncrease, query, userID)
	if err != nil {
		return 0, fmt.Errorf("获取用户增加积分失败: %w", err)
	}

	// 获取用户总减少积分
	var totalDecrease int
	query = `SELECT COALESCE(SUM(points), 0) FROM user_points WHERE user_id = ? AND type = 2`
	err = r.db.Get(&totalDecrease, query, userID)
	if err != nil {
		return 0, fmt.Errorf("获取用户减少积分失败: %w", err)
	}

	// 计算总积分
	totalPoints := totalIncrease - totalDecrease
	if totalPoints < 0 {
		// 积分不应该为负数
		return 0, nil
	}

	return totalPoints, nil
}
