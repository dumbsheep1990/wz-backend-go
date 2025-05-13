package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"wz-backend-go/internal/domain"
)

// SQLUserFavoriteRepository 用户收藏的SQL实现
type SQLUserFavoriteRepository struct {
	db *sqlx.DB
}

// NewSQLUserFavoriteRepository 创建用户收藏仓储实例
func NewSQLUserFavoriteRepository(db *sqlx.DB) domain.UserFavoriteRepository {
	return &SQLUserFavoriteRepository{
		db: db,
	}
}

// GetByID 根据ID获取收藏记录
func (r *SQLUserFavoriteRepository) GetByID(id int64) (*domain.UserFavorite, error) {
	var favorite domain.UserFavorite
	query := `SELECT id, user_id, content_id, content_type, title, 
              cover_image, description, created_at, updated_at 
              FROM user_favorites WHERE id = ? AND deleted_at IS NULL`
	err := r.db.Get(&favorite, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "收藏记录不存在")
		}
		return nil, errors.Wrap(err, "获取收藏记录失败")
	}
	return &favorite, nil
}

// Create 创建收藏记录
func (r *SQLUserFavoriteRepository) Create(favorite *domain.UserFavorite) error {
	query := `INSERT INTO user_favorites (user_id, content_id, content_type, title, 
              cover_image, description, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	favorite.CreatedAt = now
	favorite.UpdatedAt = now

	result, err := r.db.Exec(query, favorite.UserID, favorite.ContentID, favorite.ContentType,
		favorite.Title, favorite.CoverImage, favorite.Description, favorite.CreatedAt, favorite.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "创建收藏记录失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "获取收藏记录ID失败")
	}

	favorite.ID = id
	return nil
}

// ListWithConditions 根据条件查询收藏记录列表
func (r *SQLUserFavoriteRepository) ListWithConditions(conditions map[string]interface{}, offset, limit int64) ([]*domain.UserFavorite, error) {
	var favorites []*domain.UserFavorite

	query := `SELECT f.id, f.user_id, f.content_id, f.content_type, f.title, 
              f.cover_image, f.description, f.created_at, f.updated_at 
              FROM user_favorites f 
              WHERE f.deleted_at IS NULL`

	// 构建条件查询
	args := []interface{}{}
	if len(conditions) > 0 {
		condStrs := []string{}

		if v, ok := conditions["user_id"]; ok {
			condStrs = append(condStrs, "f.user_id = ?")
			args = append(args, v)
		}

		if v, ok := conditions["content_type"]; ok {
			condStrs = append(condStrs, "f.content_type = ?")
			args = append(args, v)
		}

		if v, ok := conditions["content_id"]; ok {
			condStrs = append(condStrs, "f.content_id = ?")
			args = append(args, v)
		}

		if v, ok := conditions["title_like"]; ok {
			condStrs = append(condStrs, "f.title LIKE ?")
			args = append(args, "%"+v.(string)+"%")
		}

		if v, ok := conditions["start_time"]; ok {
			condStrs = append(condStrs, "f.created_at >= ?")
			args = append(args, v)
		}

		if v, ok := conditions["end_time"]; ok {
			condStrs = append(condStrs, "f.created_at <= ?")
			args = append(args, v)
		}

		if len(condStrs) > 0 {
			query += " AND " + strings.Join(condStrs, " AND ")
		}
	}

	// 添加排序和分页
	query += " ORDER BY f.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := r.db.Select(&favorites, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "查询收藏记录列表失败")
	}

	return favorites, nil
}

// CountWithConditions 根据条件统计收藏记录数量
func (r *SQLUserFavoriteRepository) CountWithConditions(conditions map[string]interface{}) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM user_favorites f WHERE f.deleted_at IS NULL`

	// 构建条件查询
	args := []interface{}{}
	if len(conditions) > 0 {
		condStrs := []string{}

		if v, ok := conditions["user_id"]; ok {
			condStrs = append(condStrs, "f.user_id = ?")
			args = append(args, v)
		}

		if v, ok := conditions["content_type"]; ok {
			condStrs = append(condStrs, "f.content_type = ?")
			args = append(args, v)
		}

		if v, ok := conditions["content_id"]; ok {
			condStrs = append(condStrs, "f.content_id = ?")
			args = append(args, v)
		}

		if v, ok := conditions["title_like"]; ok {
			condStrs = append(condStrs, "f.title LIKE ?")
			args = append(args, "%"+v.(string)+"%")
		}

		if v, ok := conditions["start_time"]; ok {
			condStrs = append(condStrs, "f.created_at >= ?")
			args = append(args, v)
		}

		if v, ok := conditions["end_time"]; ok {
			condStrs = append(condStrs, "f.created_at <= ?")
			args = append(args, v)
		}

		if len(condStrs) > 0 {
			query += " AND " + strings.Join(condStrs, " AND ")
		}
	}

	err := r.db.Get(&count, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "统计收藏记录数量失败")
	}

	return count, nil
}

// DeleteByID 根据ID删除收藏记录
func (r *SQLUserFavoriteRepository) DeleteByID(id int64) error {
	query := `UPDATE user_favorites SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return errors.Wrap(err, "删除收藏记录失败")
	}

	return nil
}

// BatchDelete 批量删除收藏记录
func (r *SQLUserFavoriteRepository) BatchDelete(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE user_favorites SET deleted_at = ? WHERE id IN (?) AND deleted_at IS NULL`

	// 构建IN参数
	inQuery, args, err := sqlx.In(query, time.Now(), ids)
	if err != nil {
		return errors.Wrap(err, "构建批量删除SQL失败")
	}

	// 转换为DB特定语法
	inQuery = r.db.Rebind(inQuery)

	_, err = r.db.Exec(inQuery, args...)
	if err != nil {
		return errors.Wrap(err, "批量删除收藏记录失败")
	}

	return nil
}

// CountUsers 统计有收藏的用户数量
func (r *SQLUserFavoriteRepository) CountUsers() (int64, error) {
	var count int64

	query := `SELECT COUNT(DISTINCT user_id) FROM user_favorites WHERE deleted_at IS NULL`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, errors.Wrap(err, "统计用户数量失败")
	}

	return count, nil
}

// CountFavorites 统计总收藏数量
func (r *SQLUserFavoriteRepository) CountFavorites() (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM user_favorites WHERE deleted_at IS NULL`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, errors.Wrap(err, "统计收藏数量失败")
	}

	return count, nil
}

// GetTypeDistribution 获取收藏类型分布
func (r *SQLUserFavoriteRepository) GetTypeDistribution() ([]*domain.TypeDistributionItem, error) {
	var items []*domain.TypeDistributionItem

	query := `SELECT content_type as type, COUNT(*) as count 
              FROM user_favorites 
              WHERE deleted_at IS NULL 
              GROUP BY content_type 
              ORDER BY count DESC`

	err := r.db.Select(&items, query)
	if err != nil {
		return nil, errors.Wrap(err, "获取类型分布失败")
	}

	return items, nil
}

// GetHotContent 获取热门收藏内容
func (r *SQLUserFavoriteRepository) GetHotContent(limit int) ([]*domain.HotContentItem, error) {
	var items []*domain.HotContentItem

	query := `SELECT content_id, content_type, title, COUNT(*) as favorite_count 
              FROM user_favorites 
              WHERE deleted_at IS NULL 
              GROUP BY content_id, content_type 
              ORDER BY favorite_count DESC 
              LIMIT ?`

	err := r.db.Select(&items, query, limit)
	if err != nil {
		return nil, errors.Wrap(err, "获取热门内容失败")
	}

	return items, nil
}

// GetFavoriteTrend 获取收藏趋势数据
func (r *SQLUserFavoriteRepository) GetFavoriteTrend(period string, limit int) ([]*domain.TrendItem, error) {
	var items []*domain.TrendItem

	var timeFormat string
	switch period {
	case "day":
		timeFormat = "%Y-%m-%d"
	case "week":
		timeFormat = "%Y-%u" // 年-周数
	case "month":
		timeFormat = "%Y-%m"
	default:
		timeFormat = "%Y-%m-%d"
	}

	query := `SELECT 
              DATE_FORMAT(created_at, ?) as time_unit, 
              COUNT(*) as count 
              FROM user_favorites 
              WHERE deleted_at IS NULL 
              GROUP BY time_unit 
              ORDER BY time_unit DESC 
              LIMIT ?`

	err := r.db.Select(&items, query, timeFormat, limit)
	if err != nil {
		return nil, errors.Wrap(err, "获取趋势数据失败")
	}

	return items, nil
}
