package sql

import (
	"time"

	"wz-backend-go/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserFavoriteRepository 用户收藏SQL仓储实现
type UserFavoriteRepository struct {
	conn sqlx.SqlConn
}

// NewUserFavoriteRepository 创建用户收藏仓储实现
func NewUserFavoriteRepository(conn sqlx.SqlConn) *UserFavoriteRepository {
	return &UserFavoriteRepository{
		conn: conn,
	}
}

// Create 创建用户收藏
func (r *UserFavoriteRepository) Create(favorite *domain.UserFavorite) (int64, error) {
	// 检查是否已经收藏过
	exists, err := r.CheckFavorite(favorite.UserID, favorite.ItemID, favorite.ItemType)
	if err != nil {
		logx.Errorf("检查收藏状态失败: %v", err)
		return 0, err
	}

	if exists {
		logx.Infof("用户已收藏该项目: userID=%d, itemID=%d, itemType=%s", favorite.UserID, favorite.ItemID, favorite.ItemType)
		return 0, sqlx.ErrNotFound
	}

	query := `
		INSERT INTO user_favorites (
			user_id, item_id, item_type, title, cover, summary,
			url, remark, tenant_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	favorite.CreatedAt = now
	favorite.UpdatedAt = now

	result, err := r.conn.Exec(query,
		favorite.UserID, favorite.ItemID, favorite.ItemType,
		favorite.Title, favorite.Cover, favorite.Summary,
		favorite.URL, favorite.Remark, favorite.TenantID,
		favorite.CreatedAt, favorite.UpdatedAt,
	)
	if err != nil {
		logx.Errorf("创建用户收藏失败: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取新创建收藏ID失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetByID 根据ID获取收藏
func (r *UserFavoriteRepository) GetByID(id int64) (*domain.UserFavorite, error) {
	var favorite domain.UserFavorite
	query := `SELECT 
		id, user_id, item_id, item_type, title, cover, summary,
		url, remark, tenant_id, created_at, updated_at
	FROM user_favorites 
	WHERE id = ?`

	err := r.conn.QueryRow(&favorite, query, id)
	if err != nil {
		logx.Errorf("根据ID查询收藏失败: %v, id: %d", err, id)
		return nil, err
	}

	return &favorite, nil
}

// ListByUser 获取用户收藏列表
func (r *UserFavoriteRepository) ListByUser(userID int64, page, pageSize int, itemType string) ([]*domain.UserFavorite, int64, error) {
	// 构建查询条件
	whereClause := "user_id = ?"
	args := []interface{}{userID}

	// 如果指定了收藏类型，添加类型过滤
	if itemType != "" {
		whereClause += " AND item_type = ?"
		args = append(args, itemType)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM user_favorites WHERE " + whereClause
	var count int64
	err := r.conn.QueryRow(&count, countQuery, args...)
	if err != nil {
		logx.Errorf("查询用户收藏总数失败: %v", err)
		return nil, 0, err
	}

	// 计算分页
	offset := (page - 1) * pageSize

	// 查询列表
	listQuery := `
		SELECT 
			id, user_id, item_id, item_type, title, cover, summary,
			url, remark, tenant_id, created_at, updated_at
		FROM user_favorites 
		WHERE ` + whereClause + ` 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	args = append(args, pageSize, offset)
	var favorites []*domain.UserFavorite
	err = r.conn.QueryRows(&favorites, listQuery, args...)
	if err != nil {
		logx.Errorf("查询用户收藏列表失败: %v", err)
		return nil, 0, err
	}

	return favorites, count, nil
}

// Delete 删除用户收藏
func (r *UserFavoriteRepository) Delete(id int64, userID int64) error {
	query := "DELETE FROM user_favorites WHERE id = ? AND user_id = ?"
	_, err := r.conn.Exec(query, id, userID)
	if err != nil {
		logx.Errorf("删除用户收藏失败: %v, id: %d, userID: %d", err, id, userID)
		return err
	}

	return nil
}

// CheckFavorite 检查是否已收藏
func (r *UserFavoriteRepository) CheckFavorite(userID int64, itemID int64, itemType string) (bool, error) {
	query := "SELECT COUNT(*) FROM user_favorites WHERE user_id = ? AND item_id = ? AND item_type = ?"
	var count int
	err := r.conn.QueryRow(&count, query, userID, itemID, itemType)
	if err != nil {
		logx.Errorf("检查收藏状态失败: %v", err)
		return false, err
	}

	return count > 0, nil
}
