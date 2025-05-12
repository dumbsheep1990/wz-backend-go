package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"wz-backend-go/internal/domain"
)

// UserFavoriteRepository 用户收藏仓储SQL实现
type UserFavoriteRepository struct {
	db *sqlx.DB
}

// NewUserFavoriteRepository 创建用户收藏仓储实例
func NewUserFavoriteRepository(db *sqlx.DB) domain.UserFavoriteRepository {
	return &UserFavoriteRepository{
		db: db,
	}
}

// Create 创建用户收藏
func (r *UserFavoriteRepository) Create(favorite *domain.UserFavorite) (int64, error) {
	// 先检查是否已经收藏过
	exists, err := r.CheckFavorite(favorite.UserID, favorite.ItemID, favorite.ItemType)
	if err != nil {
		return 0, fmt.Errorf("检查收藏状态失败: %w", err)
	}

	if exists {
		return 0, fmt.Errorf("已经收藏过该内容")
	}

	// 设置创建时间
	now := time.Now()
	favorite.CreatedAt = now
	favorite.UpdatedAt = now

	query := `INSERT INTO user_favorites (
        user_id, item_id, item_type, title, cover, 
        summary, url, remark, tenant_id, created_at, updated_at
    ) VALUES (
        :user_id, :item_id, :item_type, :title, :cover, 
        :summary, :url, :remark, :tenant_id, :created_at, :updated_at
    )`

	result, err := r.db.NamedExec(query, favorite)
	if err != nil {
		return 0, fmt.Errorf("创建收藏失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}

	return id, nil
}

// GetByID 根据ID获取收藏记录
func (r *UserFavoriteRepository) GetByID(id int64) (*domain.UserFavorite, error) {
	var favorite domain.UserFavorite
	query := `SELECT * FROM user_favorites WHERE id = ?`

	err := r.db.Get(&favorite, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("获取收藏记录失败: %w", err)
	}

	return &favorite, nil
}

// ListByUser 获取用户收藏列表
func (r *UserFavoriteRepository) ListByUser(userID int64, page, pageSize int, itemType string) ([]*domain.UserFavorite, int64, error) {
	favorites := []*domain.UserFavorite{}
	var count int64

	// 构建查询条件
	conditions := []string{"user_id = ?"}
	args := []interface{}{userID}

	if itemType != "" {
		conditions = append(conditions, "item_type = ?")
		args = append(args, itemType)
	}

	// 构建WHERE子句
	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM user_favorites %s", whereClause)
	err := r.db.Get(&count, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("统计收藏记录数量失败: %w", err)
	}

	// 查询列表
	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf("SELECT * FROM user_favorites %s ORDER BY id DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, pageSize, offset)

	err = r.db.Select(&favorites, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询收藏列表失败: %w", err)
	}

	return favorites, count, nil
}

// Delete 删除用户收藏
func (r *UserFavoriteRepository) Delete(id int64, userID int64) error {
	query := `DELETE FROM user_favorites WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("删除收藏失败: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("收藏不存在或不属于该用户")
	}

	return nil
}

// CheckFavorite 检查是否已收藏
func (r *UserFavoriteRepository) CheckFavorite(userID int64, itemID int64, itemType string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_favorites WHERE user_id = ? AND item_id = ? AND item_type = ?)`
	err := r.db.Get(&exists, query, userID, itemID, itemType)
	if err != nil {
		return false, fmt.Errorf("检查收藏状态失败: %w", err)
	}

	return exists, nil
}
