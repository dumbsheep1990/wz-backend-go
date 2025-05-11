package sql

import (
	"time"

	"wz-backend-go/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// StatisticsRepository 统计数据SQL仓储实现
type StatisticsRepository struct {
	conn sqlx.SqlConn
}

// NewStatisticsRepository 创建统计数据仓储实现
func NewStatisticsRepository(conn sqlx.SqlConn) *StatisticsRepository {
	return &StatisticsRepository{
		conn: conn,
	}
}

// RecordStatistics 记录统计数据
func (r *StatisticsRepository) RecordStatistics(data *domain.StatisticsData) error {
	// 检查是否已存在，避免重复记录
	checkQuery := `
		SELECT id FROM statistics 
		WHERE date = ? AND type = ? AND tenant_id = ?
	`
	if data.ItemID > 0 {
		checkQuery += " AND item_id = ? AND item_type = ?"
	} else {
		checkQuery += " AND item_id = 0"
	}

	var id int64
	var err error
	var args []interface{}

	args = append(args, data.Date, data.Type, data.TenantID)
	if data.ItemID > 0 {
		args = append(args, data.ItemID, data.ItemType)
	}

	err = r.conn.QueryRow(&id, checkQuery, args...)

	now := time.Now()

	// 如果找到记录，更新value累加
	if err == nil {
		updateQuery := "UPDATE statistics SET value = value + ?, updated_at = ? WHERE id = ?"
		_, err = r.conn.Exec(updateQuery, data.Value, now, id)
		if err != nil {
			logx.Errorf("更新统计数据失败: %v, id: %d", err, id)
			return err
		}
		return nil
	}

	// 如果没有记录，则插入新记录
	if err == sqlx.ErrNotFound {
		query := `
			INSERT INTO statistics (
				date, type, value, item_id, item_type, tenant_id, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`

		data.CreatedAt = now
		data.UpdatedAt = now

		_, err = r.conn.Exec(query,
			data.Date, data.Type, data.Value, data.ItemID,
			data.ItemType, data.TenantID, data.CreatedAt, data.UpdatedAt,
		)
		if err != nil {
			logx.Errorf("插入统计数据失败: %v", err)
			return err
		}
		return nil
	}

	// 其他错误
	logx.Errorf("检查统计数据记录失败: %v", err)
	return err
}

// GetOverview 获取概览数据
func (r *StatisticsRepository) GetOverview(tenantID int64) (*domain.OverviewData, error) {
	var overview domain.OverviewData

	// 获取总用户数
	userCountQuery := "SELECT COUNT(*) FROM users WHERE tenant_id = ?"
	err := r.conn.QueryRow(&overview.TotalUsers, userCountQuery, tenantID)
	if err != nil {
		logx.Errorf("获取总用户数失败: %v", err)
		return nil, err
	}

	// 获取总内容数
	contentCountQuery := "SELECT COUNT(*) FROM contents WHERE tenant_id = ?"
	err = r.conn.QueryRow(&overview.TotalContent, contentCountQuery, tenantID)
	if err != nil {
		logx.Errorf("获取总内容数失败: %v", err)
		return nil, err
	}

	// 获取总PV
	pvQuery := "SELECT IFNULL(SUM(value), 0) FROM statistics WHERE type = 'pv' AND tenant_id = ?"
	err = r.conn.QueryRow(&overview.TotalPV, pvQuery, tenantID)
	if err != nil {
		logx.Errorf("获取总PV失败: %v", err)
		return nil, err
	}

	// 获取总UV
	uvQuery := "SELECT IFNULL(SUM(value), 0) FROM statistics WHERE type = 'uv' AND tenant_id = ?"
	err = r.conn.QueryRow(&overview.TotalUV, uvQuery, tenantID)
	if err != nil {
		logx.Errorf("获取总UV失败: %v", err)
		return nil, err
	}

	// 今日日期
	today := time.Now().Format("2006-01-02")

	// 获取今日PV
	todayPVQuery := "SELECT IFNULL(SUM(value), 0) FROM statistics WHERE type = 'pv' AND date = ? AND tenant_id = ?"
	err = r.conn.QueryRow(&overview.TodayPV, todayPVQuery, today, tenantID)
	if err != nil {
		logx.Errorf("获取今日PV失败: %v", err)
		return nil, err
	}

	// 获取今日UV
	todayUVQuery := "SELECT IFNULL(SUM(value), 0) FROM statistics WHERE type = 'uv' AND date = ? AND tenant_id = ?"
	err = r.conn.QueryRow(&overview.TodayUV, todayUVQuery, today, tenantID)
	if err != nil {
		logx.Errorf("获取今日UV失败: %v", err)
		return nil, err
	}

	// 获取今日新增用户
	todayNewUsersQuery := "SELECT COUNT(*) FROM users WHERE DATE(created_at) = ? AND tenant_id = ?"
	err = r.conn.QueryRow(&overview.TodayNewUsers, todayNewUsersQuery, today, tenantID)
	if err != nil {
		logx.Errorf("获取今日新增用户失败: %v", err)
		return nil, err
	}

	// 获取今日新增内容
	todayNewContentQuery := "SELECT COUNT(*) FROM contents WHERE DATE(created_at) = ? AND tenant_id = ?"
	err = r.conn.QueryRow(&overview.TodayNewContent, todayNewContentQuery, today, tenantID)
	if err != nil {
		logx.Errorf("获取今日新增内容失败: %v", err)
		return nil, err
	}

	return &overview, nil
}

// GetStatisticsByType 获取特定类型的统计数据
func (r *StatisticsRepository) GetStatisticsByType(tenantID int64, type_ string, startDate, endDate time.Time) ([]*domain.StatisticsData, error) {
	query := `
		SELECT 
			id, date, type, value, item_id, item_type, 
			tenant_id, created_at, updated_at
		FROM statistics 
		WHERE type = ? AND tenant_id = ? AND date BETWEEN ? AND ?
		ORDER BY date ASC
	`

	var stats []*domain.StatisticsData
	err := r.conn.QueryRows(&stats, query, type_, tenantID, startDate, endDate)
	if err != nil {
		logx.Errorf("获取统计数据失败: %v, type: %s", err, type_)
		return nil, err
	}

	return stats, nil
}

// GetContentRanking 获取内容排行榜
func (r *StatisticsRepository) GetContentRanking(tenantID int64, limit int) ([]*domain.ContentRankingItem, error) {
	// 首先查询访问量最高的内容ID
	query := `
		SELECT 
			s.item_id as id,
			c.title,
			c.cover,
			c.type,
			IFNULL(SUM(CASE WHEN s.type = 'content_view' THEN s.value ELSE 0 END), 0) as view_count,
			IFNULL(SUM(CASE WHEN s.type = 'content_like' THEN s.value ELSE 0 END), 0) as like_count,
			CONCAT('/content/', c.id) as url
		FROM statistics s
		JOIN contents c ON s.item_id = c.id
		WHERE s.tenant_id = ? AND s.item_type = 'content' AND c.status = 1
		GROUP BY s.item_id, c.title, c.cover, c.type
		ORDER BY view_count DESC
		LIMIT ?
	`

	var ranking []*domain.ContentRankingItem
	err := r.conn.QueryRows(&ranking, query, tenantID, limit)
	if err != nil {
		logx.Errorf("获取内容排行榜失败: %v", err)
		return nil, err
	}

	return ranking, nil
}
