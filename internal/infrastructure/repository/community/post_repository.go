package community

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domain "wz-backend-go/internal/domain/community"
)

// PostRepository 帖子仓储实现
type PostRepository struct {
	db *sql.DB
}

// NewPostRepository 创建帖子仓储
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// Save 保存帖子
func (r *PostRepository) Save(ctx context.Context, post *domain.Post) error {
	// 检查帖子是否已存在
	existingPost, err := r.FindByID(ctx, post.ID())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingPost == nil {
		// 新增帖子
		query := `
			INSERT INTO community_posts (
				id, title, content, author_id, author_name, community_id, group_id,
				status, created_at, updated_at, like_count, view_count, comment_count
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			post.ID().Value(),
			post.Title().Value(),
			post.Content().Value(),
			post.AuthorID().Value(),
			post.AuthorName(),
			post.CommunityID().Value(),
			post.GroupID().Value(),
			int(post.Status()),
			post.CreatedAt().Value(),
			post.UpdatedAt().Value(),
			post.LikeCount(),
			post.ViewCount(),
			post.CommentCount(),
		)

		if err != nil {
			return err
		}

		// 存储标签
		if len(post.Tags().Values()) > 0 {
			if err := r.saveTags(ctx, post.ID().Value(), post.Tags().Values()); err != nil {
				return err
			}
		}

		// 存储图片
		if len(post.Images().Values()) > 0 {
			if err := r.saveImages(ctx, post.ID().Value(), post.Images().Values()); err != nil {
				return err
			}
		}

		return nil
	} else {
		// 更新帖子
		query := `
			UPDATE community_posts SET
				title = ?,
				content = ?,
				status = ?,
				updated_at = ?,
				like_count = ?,
				view_count = ?,
				comment_count = ?
			WHERE id = ?
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			post.Title().Value(),
			post.Content().Value(),
			int(post.Status()),
			post.UpdatedAt().Value(),
			post.LikeCount(),
			post.ViewCount(),
			post.CommentCount(),
			post.ID().Value(),
		)

		if err != nil {
			return err
		}

		// 更新标签（先删除旧标签，再添加新标签）
		if err := r.deleteTags(ctx, post.ID().Value()); err != nil {
			return err
		}

		if len(post.Tags().Values()) > 0 {
			if err := r.saveTags(ctx, post.ID().Value(), post.Tags().Values()); err != nil {
				return err
			}
		}

		// 更新图片（先删除旧图片，再添加新图片）
		if err := r.deleteImages(ctx, post.ID().Value()); err != nil {
			return err
		}

		if len(post.Images().Values()) > 0 {
			if err := r.saveImages(ctx, post.ID().Value(), post.Images().Values()); err != nil {
				return err
			}
		}

		return nil
	}
}

// 保存标签
func (r *PostRepository) saveTags(ctx context.Context, postID string, tags []string) error {
	query := `INSERT INTO post_tags (post_id, tag) VALUES (?, ?)`
	
	for _, tag := range tags {
		_, err := r.db.ExecContext(ctx, query, postID, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除标签
func (r *PostRepository) deleteTags(ctx context.Context, postID string) error {
	query := `DELETE FROM post_tags WHERE post_id = ?`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

// 保存图片
func (r *PostRepository) saveImages(ctx context.Context, postID string, images []string) error {
	query := `INSERT INTO post_images (post_id, image_url) VALUES (?, ?)`
	
	for _, imageURL := range images {
		_, err := r.db.ExecContext(ctx, query, postID, imageURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除图片
func (r *PostRepository) deleteImages(ctx context.Context, postID string) error {
	query := `DELETE FROM post_images WHERE post_id = ?`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

// FindByID 根据ID查找帖子
func (r *PostRepository) FindByID(ctx context.Context, id domain.ID) (*domain.Post, error) {
	query := `
		SELECT 
			id, title, content, author_id, author_name, community_id, group_id,
			status, created_at, updated_at, like_count, view_count, comment_count
		FROM community_posts
		WHERE id = ?
	`

	var (
		idStr        string
		title        string
		content      string
		authorID     string
		authorName   string
		communityID  string
		groupID      string
		status       int
		createdAt    time.Time
		updatedAt    time.Time
		likeCount    int
		viewCount    int
		commentCount int
	)

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idStr, &title, &content, &authorID, &authorName, &communityID, &groupID,
		&status, &createdAt, &updatedAt, &likeCount, &viewCount, &commentCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("帖子不存在")
		}
		return nil, err
	}

	// 查询标签
	tags, err := r.findTags(ctx, idStr)
	if err != nil {
		return nil, err
	}

	// 查询图片
	images, err := r.findImages(ctx, idStr)
	if err != nil {
		return nil, err
	}

	// 重构帖子实体
	return domain.ReconstructPost(
		idStr, title, content, authorID, authorName, communityID, groupID,
		status, tags, images, createdAt, updatedAt, likeCount, viewCount, commentCount,
	)
}

// 查询标签
func (r *PostRepository) findTags(ctx context.Context, postID string) ([]string, error) {
	query := `SELECT tag FROM post_tags WHERE post_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return tags, nil
}

// 查询图片
func (r *PostRepository) findImages(ctx context.Context, postID string) ([]string, error) {
	query := `SELECT image_url FROM post_images WHERE post_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var images []string
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return nil, err
		}
		images = append(images, imageURL)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return images, nil
}

// FindByCommunity 查找社区下的帖子（分页）
func (r *PostRepository) FindByCommunity(ctx context.Context, communityID domain.ID, offset, limit int) ([]*domain.Post, error) {
	query := `
		SELECT 
			id, title, content, author_id, author_name, community_id, group_id,
			status, created_at, updated_at, like_count, view_count, comment_count
		FROM community_posts
		WHERE community_id = ? AND status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, communityID.Value(), int(domain.PostStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(ctx, rows)
}

// CountByCommunity 统计社区下的帖子数量
func (r *PostRepository) CountByCommunity(ctx context.Context, communityID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM community_posts WHERE community_id = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, communityID.Value(), int(domain.PostStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByGroup 查找小组下的帖子（分页）
func (r *PostRepository) FindByGroup(ctx context.Context, groupID domain.ID, offset, limit int) ([]*domain.Post, error) {
	query := `
		SELECT 
			id, title, content, author_id, author_name, community_id, group_id,
			status, created_at, updated_at, like_count, view_count, comment_count
		FROM community_posts
		WHERE group_id = ? AND status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, groupID.Value(), int(domain.PostStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(ctx, rows)
}

// CountByGroup 统计小组下的帖子数量
func (r *PostRepository) CountByGroup(ctx context.Context, groupID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM community_posts WHERE group_id = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, groupID.Value(), int(domain.PostStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByAuthor 查找用户发布的帖子（分页）
func (r *PostRepository) FindByAuthor(ctx context.Context, authorID domain.ID, offset, limit int) ([]*domain.Post, error) {
	query := `
		SELECT 
			id, title, content, author_id, author_name, community_id, group_id,
			status, created_at, updated_at, like_count, view_count, comment_count
		FROM community_posts
		WHERE author_id = ? AND status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, authorID.Value(), int(domain.PostStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(ctx, rows)
}

// CountByAuthor 统计用户发布的帖子数量
func (r *PostRepository) CountByAuthor(ctx context.Context, authorID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM community_posts WHERE author_id = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, authorID.Value(), int(domain.PostStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByTags 根据标签查找帖子（分页）
func (r *PostRepository) FindByTags(ctx context.Context, tags []string, offset, limit int) ([]*domain.Post, error) {
	// 构建查询条件
	placeholders := make([]string, len(tags))
	args := make([]interface{}, 0, len(tags)+3)
	
	for i := range tags {
		placeholders[i] = "?"
		args = append(args, tags[i])
	}
	
	// 添加状态和分页参数
	args = append(args, int(domain.PostStatusDeleted), limit, offset)
	
	query := `
		SELECT 
			p.id, p.title, p.content, p.author_id, p.author_name, p.community_id, p.group_id,
			p.status, p.created_at, p.updated_at, p.like_count, p.view_count, p.comment_count
		FROM community_posts p
		JOIN post_tags t ON p.id = t.post_id
		WHERE t.tag IN (` + placeholders[0] + `) AND p.status != ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(ctx, rows)
}

// FindRecent 查找最近的帖子（分页）
func (r *PostRepository) FindRecent(ctx context.Context, offset, limit int) ([]*domain.Post, error) {
	query := `
		SELECT 
			id, title, content, author_id, author_name, community_id, group_id,
			status, created_at, updated_at, like_count, view_count, comment_count
		FROM community_posts
		WHERE status = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, int(domain.PostStatusActive), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(ctx, rows)
}

// FindPopular 查找热门帖子（分页）
func (r *PostRepository) FindPopular(ctx context.Context, offset, limit int) ([]*domain.Post, error) {
	query := `
		SELECT 
			id, title, content, author_id, author_name, community_id, group_id,
			status, created_at, updated_at, like_count, view_count, comment_count
		FROM community_posts
		WHERE status = ?
		ORDER BY (like_count * 3 + view_count + comment_count * 2) DESC, created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, int(domain.PostStatusActive), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(ctx, rows)
}

// Delete 删除帖子
func (r *PostRepository) Delete(ctx context.Context, id domain.ID) error {
	// 注意：这里只是逻辑删除，而不是物理删除
	query := `UPDATE community_posts SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, int(domain.PostStatusDeleted), time.Now(), id.Value())
	return err
}

// scanPosts 扫描帖子结果集
func (r *PostRepository) scanPosts(ctx context.Context, rows *sql.Rows) ([]*domain.Post, error) {
	var posts []*domain.Post
	for rows.Next() {
		var (
			idStr        string
			title        string
			content      string
			authorID     string
			authorName   string
			communityID  string
			groupID      string
			status       int
			createdAt    time.Time
			updatedAt    time.Time
			likeCount    int
			viewCount    int
			commentCount int
		)

		if err := rows.Scan(
			&idStr, &title, &content, &authorID, &authorName, &communityID, &groupID,
			&status, &createdAt, &updatedAt, &likeCount, &viewCount, &commentCount,
		); err != nil {
			return nil, err
		}

		// 查询标签
		tags, err := r.findTags(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 查询图片
		images, err := r.findImages(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 重构帖子实体
		post, err := domain.ReconstructPost(
			idStr, title, content, authorID, authorName, communityID, groupID,
			status, tags, images, createdAt, updatedAt, likeCount, viewCount, commentCount,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
