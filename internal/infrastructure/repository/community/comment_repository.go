package community

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domain "wz-backend-go/internal/domain/community"
)

// CommentRepository 评论仓储实现
type CommentRepository struct {
	db *sql.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// Save 保存评论
func (r *CommentRepository) Save(ctx context.Context, comment *domain.Comment) error {
	// 检查评论是否已存在
	existingComment, err := r.FindByID(ctx, comment.ID())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingComment == nil {
		// 新增评论
		query := `
			INSERT INTO community_comments (
				id, content, author_id, author_name, post_id, parent_id,
				status, created_at, like_count
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			comment.ID().Value(),
			comment.Content().Value(),
			comment.AuthorID().Value(),
			comment.AuthorName(),
			comment.PostID().Value(),
			comment.ParentID().Value(),
			int(comment.Status()),
			comment.CreatedAt().Value(),
			comment.LikeCount(),
		)

		return err
	} else {
		// 更新评论
		query := `
			UPDATE community_comments SET
				status = ?,
				like_count = ?
			WHERE id = ?
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			int(comment.Status()),
			comment.LikeCount(),
			comment.ID().Value(),
		)

		return err
	}
}

// FindByID 根据ID查找评论
func (r *CommentRepository) FindByID(ctx context.Context, id domain.ID) (*domain.Comment, error) {
	query := `
		SELECT 
			id, content, author_id, author_name, post_id, parent_id,
			status, created_at, like_count
		FROM community_comments
		WHERE id = ?
	`

	var (
		idStr       string
		content     string
		authorID    string
		authorName  string
		postID      string
		parentID    sql.NullString
		status      int
		createdAt   time.Time
		likeCount   int
	)

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idStr, &content, &authorID, &authorName, &postID, &parentID,
		&status, &createdAt, &likeCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}

	// 处理可能为空的父评论ID
	parentIDValue := ""
	if parentID.Valid {
		parentIDValue = parentID.String
	}

	// 重构评论实体
	return domain.ReconstructComment(
		idStr, content, authorID, authorName, postID, parentIDValue,
		status, createdAt, likeCount,
	)
}

// FindByPost 查找帖子下的评论（分页）
func (r *CommentRepository) FindByPost(ctx context.Context, postID domain.ID, offset, limit int) ([]*domain.Comment, error) {
	query := `
		SELECT 
			id, content, author_id, author_name, post_id, parent_id,
			status, created_at, like_count
		FROM community_comments
		WHERE post_id = ? AND parent_id IS NULL AND status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, postID.Value(), int(domain.CommentStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// CountByPost 统计帖子下的评论数量
func (r *CommentRepository) CountByPost(ctx context.Context, postID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM community_comments WHERE post_id = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, postID.Value(), int(domain.CommentStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByParent 查找父评论下的回复（分页）
func (r *CommentRepository) FindByParent(ctx context.Context, parentID domain.ID, offset, limit int) ([]*domain.Comment, error) {
	query := `
		SELECT 
			id, content, author_id, author_name, post_id, parent_id,
			status, created_at, like_count
		FROM community_comments
		WHERE parent_id = ? AND status != ?
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, parentID.Value(), int(domain.CommentStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// CountByParent 统计父评论下的回复数量
func (r *CommentRepository) CountByParent(ctx context.Context, parentID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM community_comments WHERE parent_id = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, parentID.Value(), int(domain.CommentStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByAuthor 查找用户发布的评论（分页）
func (r *CommentRepository) FindByAuthor(ctx context.Context, authorID domain.ID, offset, limit int) ([]*domain.Comment, error) {
	query := `
		SELECT 
			id, content, author_id, author_name, post_id, parent_id,
			status, created_at, like_count
		FROM community_comments
		WHERE author_id = ? AND status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, authorID.Value(), int(domain.CommentStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// CountByAuthor 统计用户发布的评论数量
func (r *CommentRepository) CountByAuthor(ctx context.Context, authorID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM community_comments WHERE author_id = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, authorID.Value(), int(domain.CommentStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindRecent 查找最近的评论（分页）
func (r *CommentRepository) FindRecent(ctx context.Context, offset, limit int) ([]*domain.Comment, error) {
	query := `
		SELECT 
			id, content, author_id, author_name, post_id, parent_id,
			status, created_at, like_count
		FROM community_comments
		WHERE status = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, int(domain.CommentStatusActive), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// Delete 删除评论
func (r *CommentRepository) Delete(ctx context.Context, id domain.ID) error {
	// 注意：这里只是逻辑删除，而不是物理删除
	query := `UPDATE community_comments SET status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, int(domain.CommentStatusDeleted), id.Value())
	return err
}

// scanComments 扫描评论结果集
func (r *CommentRepository) scanComments(rows *sql.Rows) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	for rows.Next() {
		var (
			idStr       string
			content     string
			authorID    string
			authorName  string
			postID      string
			parentID    sql.NullString
			status      int
			createdAt   time.Time
			likeCount   int
		)

		if err := rows.Scan(
			&idStr, &content, &authorID, &authorName, &postID, &parentID,
			&status, &createdAt, &likeCount,
		); err != nil {
			return nil, err
		}

		// 处理可能为空的父评论ID
		parentIDValue := ""
		if parentID.Valid {
			parentIDValue = parentID.String
		}

		// 重构评论实体
		comment, err := domain.ReconstructComment(
			idStr, content, authorID, authorName, postID, parentIDValue,
			status, createdAt, likeCount,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
