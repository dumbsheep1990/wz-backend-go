package community

import (
	"context"
	"database/sql"
	"time"

	domain "wz-backend-go/internal/domain/community"
)

// ViewRepository 浏览记录仓储实现
type ViewRepository struct {
	db *sql.DB
}

// NewViewRepository 创建浏览记录仓储
func NewViewRepository(db *sql.DB) *ViewRepository {
	return &ViewRepository{
		db: db,
	}
}

// SaveView 保存浏览记录
func (r *ViewRepository) SaveView(ctx context.Context, postID domain.ID, userID domain.ID) error {
	// 浏览记录使用 REPLACE INTO 语法，确保每个用户对每个帖子只有一条最新的浏览记录
	query := `
		REPLACE INTO post_views (
			post_id, user_id, viewed_at
		) VALUES (?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, postID.Value(), userID.Value(), time.Now())
	return err
}

// CheckView 检查用户是否浏览过帖子
func (r *ViewRepository) CheckView(ctx context.Context, postID domain.ID, userID domain.ID) (bool, error) {
	query := `SELECT COUNT(*) FROM post_views WHERE post_id = ? AND user_id = ?`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, postID.Value(), userID.Value()).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// CountViews 统计帖子浏览数量
func (r *ViewRepository) CountViews(ctx context.Context, postID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM post_views WHERE post_id = ?`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, postID.Value()).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// FindPopularByViews 查找浏览量最高的帖子（分页）
func (r *ViewRepository) FindPopularByViews(ctx context.Context, offset, limit int) ([]*domain.Post, error) {
	query := `
		SELECT 
			p.id, p.title, p.content, p.author_id, p.author_name, p.community_id, p.group_id,
			p.status, p.created_at, p.updated_at, p.like_count, p.view_count, p.comment_count
		FROM community_posts p
		WHERE p.status = ?
		ORDER BY p.view_count DESC, p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, int(domain.PostStatusActive), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 使用PostRepository的方法来构建结果
	postRepo := &PostRepository{db: r.db}
	return postRepo.scanPosts(ctx, rows)
}
