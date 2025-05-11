package sql

import (
	"time"

	"wz-backend-go/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserPointsRepository 用户积分SQL仓储实现
type UserPointsRepository struct {
	conn sqlx.SqlConn
}

// NewUserPointsRepository 创建用户积分仓储实现
func NewUserPointsRepository(conn sqlx.SqlConn) *UserPointsRepository {
	return &UserPointsRepository{
		conn: conn,
	}
}

// Create 创建用户积分记录
func (r *UserPointsRepository) Create(points *domain.UserPoints) (int64, error) {
	// 首先获取用户当前总积分
	currentTotal, err := r.GetUserTotalPoints(points.UserID)
	if err != nil {
		logx.Errorf("获取用户当前总积分失败: %v, userID: %d", err, points.UserID)
		return 0, err
	}

	// 根据积分变动类型计算新的总积分
	newTotal := currentTotal
	if points.Type == 1 { // 增加积分
		newTotal += points.Points
	} else if points.Type == 2 { // 减少积分
		newTotal -= points.Points
		// 防止总积分变为负数
		if newTotal < 0 {
			logx.Warnf("用户积分不足, userID: %d, 当前积分: %d, 尝试扣减: %d", points.UserID, currentTotal, points.Points)
			return 0, sqlx.ErrNotFound
		}
	}

	// 设置总积分
	points.TotalPoints = newTotal

	// 插入积分记录
	query := `
		INSERT INTO user_points (
			user_id, points, total_points, type, source, description,
			related_id, related_type, tenant_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	points.CreatedAt = now
	points.UpdatedAt = now

	// 开始数据库事务
	tx, err := r.conn.Begin()
	if err != nil {
		logx.Errorf("开始事务失败: %v", err)
		return 0, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logx.Errorf("回滚事务失败: %v", rollbackErr)
			}
		}
	}()

	// 执行插入操作
	result, err := tx.Exec(query,
		points.UserID, points.Points, points.TotalPoints, points.Type,
		points.Source, points.Description, points.RelatedID,
		points.RelatedType, points.TenantID, points.CreatedAt, points.UpdatedAt,
	)
	if err != nil {
		logx.Errorf("创建用户积分记录失败: %v", err)
		return 0, err
	}

	// 更新用户表中的总积分字段（如果有）
	updateUserPointsQuery := "UPDATE users SET total_points = ? WHERE id = ?"
	_, err = tx.Exec(updateUserPointsQuery, newTotal, points.UserID)
	if err != nil {
		logx.Errorf("更新用户总积分失败: %v, userID: %d", err, points.UserID)
		return 0, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		logx.Errorf("提交事务失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取新创建积分记录ID失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetByID 根据ID获取积分记录
func (r *UserPointsRepository) GetByID(id int64) (*domain.UserPoints, error) {
	var points domain.UserPoints
	query := `SELECT 
		id, user_id, points, total_points, type, source, description, 
		related_id, related_type, tenant_id, created_at, updated_at
	FROM user_points 
	WHERE id = ?`

	err := r.conn.QueryRow(&points, query, id)
	if err != nil {
		logx.Errorf("根据ID查询积分记录失败: %v, id: %d", err, id)
		return nil, err
	}

	return &points, nil
}

// ListByUser 获取用户积分记录列表
func (r *UserPointsRepository) ListByUser(userID int64, page, pageSize int) ([]*domain.UserPoints, int64, error) {
	// 查询总数
	countQuery := "SELECT COUNT(*) FROM user_points WHERE user_id = ?"
	var count int64
	err := r.conn.QueryRow(&count, countQuery, userID)
	if err != nil {
		logx.Errorf("查询用户积分记录总数失败: %v", err)
		return nil, 0, err
	}

	// 计算分页
	offset := (page - 1) * pageSize

	// 查询列表
	listQuery := `
		SELECT 
			id, user_id, points, total_points, type, source, description, 
			related_id, related_type, tenant_id, created_at, updated_at
		FROM user_points 
		WHERE user_id = ? 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	var pointsList []*domain.UserPoints
	err = r.conn.QueryRows(&pointsList, listQuery, userID, pageSize, offset)
	if err != nil {
		logx.Errorf("查询用户积分记录列表失败: %v", err)
		return nil, 0, err
	}

	return pointsList, count, nil
}

// GetUserTotalPoints 获取用户总积分
func (r *UserPointsRepository) GetUserTotalPoints(userID int64) (int, error) {
	// 优先从users表获取（更高效）
	var totalPoints int
	userQuery := "SELECT total_points FROM users WHERE id = ?"
	err := r.conn.QueryRow(&totalPoints, userQuery, userID)
	if err == nil {
		return totalPoints, nil
	}

	// 如果从用户表获取失败，则从积分记录表计算
	calcQuery := "SELECT IFNULL(MAX(total_points), 0) FROM user_points WHERE user_id = ? ORDER BY id DESC LIMIT 1"
	err = r.conn.QueryRow(&totalPoints, calcQuery, userID)
	if err != nil {
		// 如果没有找到记录，默认为0积分
		if err == sqlx.ErrNotFound {
			return 0, nil
		}
		logx.Errorf("计算用户总积分失败: %v, userID: %d", err, userID)
		return 0, err
	}

	return totalPoints, nil
}
