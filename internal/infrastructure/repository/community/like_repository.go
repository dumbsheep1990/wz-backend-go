package community

import (
	"context"
	"database/sql"
	"time"

	domain "wz-backend-go/internal/domain/community"
)

// LikeRepository 点赞仓储实现
type LikeRepository struct {
	db *sql.DB
}

// NewLikeRepository 创建点赞仓储
func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

// SavePostLike 保存帖子点赞
func (r *LikeRepository) SavePostLike(ctx context.Context, postID domain.ID, userID domain.ID) error {
	// 检查是否已经点赞
	liked, err := r.CheckPostLike(ctx, postID, userID)
	if err != nil {
		return err
	}

	if liked {
		return nil // 已经点赞过，不需要再次保存
	}

	query := `INSERT INTO post_likes (post_id, user_id, created_at) VALUES (?, ?, ?)`
	_, err = r.db.ExecContext(ctx, query, postID.Value(), userID.Value(), time.Now())
	return err
}

// RemovePostLike 移除帖子点赞
func (r *LikeRepository) RemovePostLike(ctx context.Context, postID domain.ID, userID domain.ID) error {
	query := `DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`
	_, err := r.db.ExecContext(ctx, query, postID.Value(), userID.Value())
	return err
}

// CheckPostLike 检查用户是否点赞过帖子
func (r *LikeRepository) CheckPostLike(ctx context.Context, postID domain.ID, userID domain.ID) (bool, error) {
	query := `SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND user_id = ?`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, postID.Value(), userID.Value()).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// CountPostLikes 统计帖子点赞数量
func (r *LikeRepository) CountPostLikes(ctx context.Context, postID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM post_likes WHERE post_id = ?`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, postID.Value()).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// SaveCommentLike 保存评论点赞
func (r *LikeRepository) SaveCommentLike(ctx context.Context, commentID domain.ID, userID domain.ID) error {
	// 检查是否已经点赞
	liked, err := r.CheckCommentLike(ctx, commentID, userID)
	if err != nil {
		return err
	}

	if liked {
		return nil // 已经点赞过，不需要再次保存
	}

	query := `INSERT INTO comment_likes (comment_id, user_id, created_at) VALUES (?, ?, ?)`
	_, err = r.db.ExecContext(ctx, query, commentID.Value(), userID.Value(), time.Now())
	return err
}

// RemoveCommentLike 移除评论点赞
func (r *LikeRepository) RemoveCommentLike(ctx context.Context, commentID domain.ID, userID domain.ID) error {
	query := `DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	_, err := r.db.ExecContext(ctx, query, commentID.Value(), userID.Value())
	return err
}

// CheckCommentLike 检查用户是否点赞过评论
func (r *LikeRepository) CheckCommentLike(ctx context.Context, commentID domain.ID, userID domain.ID) (bool, error) {
	query := `SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, commentID.Value(), userID.Value()).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// CountCommentLikes 统计评论点赞数量
func (r *LikeRepository) CountCommentLikes(ctx context.Context, commentID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, commentID.Value()).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}
