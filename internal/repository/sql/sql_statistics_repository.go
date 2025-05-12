package sql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain"
)

// StatisticsRepository 统计仓储SQL实现
type StatisticsRepository struct {
	db *sqlx.DB
}

// NewStatisticsRepository 创建统计仓储实例
func NewStatisticsRepository(db *sqlx.DB) domain.StatisticsRepository {
	return &StatisticsRepository{
		db: db,
	}
}

// RecordStatistics 记录统计数据
func (r *StatisticsRepository) RecordStatistics(data *domain.StatisticsData) error {
	// 检查是否已经存在当天的数据
	var exists bool
	var existingID int64
	checkQuery := `SELECT EXISTS(SELECT 1 FROM statistics_data 
        WHERE date = DATE(?) AND type = ? AND item_id = ? AND item_type = ? AND tenant_id = ?)`

	err := r.db.Get(&exists, checkQuery,
		data.Date, data.Type, data.ItemID, data.ItemType, data.TenantID)
	if err != nil {
		return fmt.Errorf("检查统计数据失败: %w", err)
	}

	now := time.Now()

	if exists {
		// 如果存在则更新
		idQuery := `SELECT id, value FROM statistics_data 
            WHERE date = DATE(?) AND type = ? AND item_id = ? AND item_type = ? AND tenant_id = ?
            LIMIT 1`

		var currentValue int64
		err = r.db.QueryRow(idQuery,
			data.Date, data.Type, data.ItemID, data.ItemType, data.TenantID).Scan(&existingID, &currentValue)
		if err != nil {
			return fmt.Errorf("获取现有统计数据失败: %w", err)
		}

		updateQuery := `UPDATE statistics_data SET value = value + ?, updated_at = ? WHERE id = ?`
		_, err = r.db.Exec(updateQuery, data.Value, now, existingID)
		if err != nil {
			return fmt.Errorf("更新统计数据失败: %w", err)
		}
	} else {
		// 否则插入新记录
		data.CreatedAt = now
		data.UpdatedAt = now

		insertQuery := `INSERT INTO statistics_data (
            date, type, value, item_id, item_type, tenant_id, created_at, updated_at
        ) VALUES (
            :date, :type, :value, :item_id, :item_type, :tenant_id, :created_at, :updated_at
        )`

		_, err = r.db.NamedExec(insertQuery, data)
		if err != nil {
			return fmt.Errorf("插入统计数据失败: %w", err)
		}
	}

	return nil
}

// GetOverview 获取概览数据
func (r *StatisticsRepository) GetOverview(tenantID int64) (*domain.OverviewData, error) {
	overview := &domain.OverviewData{}
	today := time.Now().Format("2006-01-02")

	// 查询总用户数
	userQuery := `SELECT COUNT(*) FROM users WHERE tenant_id = ?`
	err := r.db.Get(&overview.TotalUsers, userQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询用户总数失败: %w", err)
	}

	// 查询总内容数
	contentQuery := `SELECT COUNT(*) FROM content WHERE tenant_id = ?`
	err = r.db.Get(&overview.TotalContent, contentQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询内容总数失败: %w", err)
	}

	// 查询总PV
	pvQuery := `SELECT COALESCE(SUM(value), 0) FROM statistics_data 
        WHERE type = 'pv' AND tenant_id = ?`
	err = r.db.Get(&overview.TotalPV, pvQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询总PV失败: %w", err)
	}

	// 查询总UV
	uvQuery := `SELECT COALESCE(SUM(value), 0) FROM statistics_data 
        WHERE type = 'uv' AND tenant_id = ?`
	err = r.db.Get(&overview.TotalUV, uvQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询总UV失败: %w", err)
	}

	// 查询今日PV
	todayPVQuery := `SELECT COALESCE(SUM(value), 0) FROM statistics_data 
        WHERE type = 'pv' AND date = ? AND tenant_id = ?`
	err = r.db.Get(&overview.TodayPV, todayPVQuery, today, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询今日PV失败: %w", err)
	}

	// 查询今日UV
	todayUVQuery := `SELECT COALESCE(SUM(value), 0) FROM statistics_data 
        WHERE type = 'uv' AND date = ? AND tenant_id = ?`
	err = r.db.Get(&overview.TodayUV, todayUVQuery, today, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询今日UV失败: %w", err)
	}

	// 查询今日新增用户
	todayUserQuery := `SELECT COUNT(*) FROM users 
        WHERE DATE(created_at) = ? AND tenant_id = ?`
	err = r.db.Get(&overview.TodayNewUsers, todayUserQuery, today, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询今日新增用户失败: %w", err)
	}

	// 查询今日新增内容
	todayContentQuery := `SELECT COUNT(*) FROM content 
        WHERE DATE(created_at) = ? AND tenant_id = ?`
	err = r.db.Get(&overview.TodayNewContent, todayContentQuery, today, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询今日新增内容失败: %w", err)
	}

	return overview, nil
}

// GetStatisticsByType 根据类型获取统计数据
func (r *StatisticsRepository) GetStatisticsByType(tenantID int64, type_ string, startDate, endDate time.Time) ([]*domain.StatisticsData, error) {
	statsList := []*domain.StatisticsData{}

	query := `SELECT * FROM statistics_data 
        WHERE type = ? AND tenant_id = ? AND date BETWEEN ? AND ? 
        ORDER BY date ASC`

	err := r.db.Select(&statsList, query, type_, tenantID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("查询统计数据失败: %w", err)
	}

	return statsList, nil
}

// GetContentRanking 获取内容排行榜
func (r *StatisticsRepository) GetContentRanking(tenantID int64, limit int) ([]*domain.ContentRankingItem, error) {
	rankingList := []*domain.ContentRankingItem{}

	// 查询浏览量排行
	query := `
        SELECT 
            c.id, 
            c.title, 
            c.cover, 
            c.type, 
            COALESCE(view_stats.total_views, 0) as view_count,
            COALESCE(like_stats.total_likes, 0) as like_count,
            c.url
        FROM 
            content c
        LEFT JOIN (
            SELECT 
                item_id, 
                SUM(value) as total_views
            FROM 
                statistics_data
            WHERE 
                type = 'content_view' 
                AND tenant_id = ?
            GROUP BY 
                item_id
        ) view_stats ON c.id = view_stats.item_id
        LEFT JOIN (
            SELECT 
                item_id, 
                SUM(value) as total_likes
            FROM 
                statistics_data
            WHERE 
                type = 'content_like' 
                AND tenant_id = ?
            GROUP BY 
                item_id
        ) like_stats ON c.id = like_stats.item_id
        WHERE 
            c.tenant_id = ?
        ORDER BY 
            view_count DESC, like_count DESC
        LIMIT ?
    `

	err := r.db.Select(&rankingList, query, tenantID, tenantID, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("查询内容排行失败: %w", err)
	}

	return rankingList, nil
}
