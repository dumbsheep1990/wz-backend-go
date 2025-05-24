package community

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	domain "wz-backend-go/internal/domain/community"
)

// CommunityRepository 社区仓储实现
type CommunityRepository struct {
	db *sql.DB
}

// NewCommunityRepository 创建社区仓储
func NewCommunityRepository(db *sql.DB) *CommunityRepository {
	return &CommunityRepository{
		db: db,
	}
}

// Save 保存社区
func (r *CommunityRepository) Save(ctx context.Context, community *domain.Community) error {
	// 检查社区是否已存在
	existingCommunity, err := r.FindByID(ctx, community.ID())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingCommunity == nil {
		// 新增社区
		query := `
			INSERT INTO communities (
				id, name, description, owner_id, owner_name, status, 
				community_type, location, created_at, updated_at, 
				group_count, member_count, post_count
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		// 将标签列表序列化为JSON
		tags := strings.Join(community.Tags().Values(), ",")

		_, err = r.db.ExecContext(
			ctx,
			query,
			community.ID().Value(),
			community.Name().Value(),
			community.Description().Value(),
			community.OwnerID().Value(),
			community.OwnerName(),
			int(community.Status()),
			community.Type().String(),
			community.Location().Value(),
			community.CreatedAt().Value(),
			community.UpdatedAt().Value(),
			community.GroupCount(),
			community.MemberCount(),
			community.PostCount(),
		)

		if err != nil {
			return err
		}

		// 存储标签
		if len(community.Tags().Values()) > 0 {
			if err := r.saveTags(ctx, community.ID().Value(), community.Tags().Values()); err != nil {
				return err
			}
		}

		return nil
	} else {
		// 更新社区
		query := `
			UPDATE communities SET
				name = ?,
				description = ?,
				status = ?,
				location = ?,
				updated_at = ?,
				group_count = ?,
				member_count = ?,
				post_count = ?
			WHERE id = ?
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			community.Name().Value(),
			community.Description().Value(),
			int(community.Status()),
			community.Location().Value(),
			community.UpdatedAt().Value(),
			community.GroupCount(),
			community.MemberCount(),
			community.PostCount(),
			community.ID().Value(),
		)

		if err != nil {
			return err
		}

		// 更新标签（先删除旧标签，再添加新标签）
		if err := r.deleteTags(ctx, community.ID().Value()); err != nil {
			return err
		}

		if len(community.Tags().Values()) > 0 {
			if err := r.saveTags(ctx, community.ID().Value(), community.Tags().Values()); err != nil {
				return err
			}
		}

		return nil
	}
}

// 保存标签
func (r *CommunityRepository) saveTags(ctx context.Context, communityID string, tags []string) error {
	query := `INSERT INTO community_tags (community_id, tag) VALUES (?, ?)`
	
	for _, tag := range tags {
		_, err := r.db.ExecContext(ctx, query, communityID, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除标签
func (r *CommunityRepository) deleteTags(ctx context.Context, communityID string) error {
	query := `DELETE FROM community_tags WHERE community_id = ?`
	_, err := r.db.ExecContext(ctx, query, communityID)
	return err
}

// FindByID 根据ID查找社区
func (r *CommunityRepository) FindByID(ctx context.Context, id domain.ID) (*domain.Community, error) {
	query := `
		SELECT 
			id, name, description, owner_id, owner_name, status, 
			community_type, location, created_at, updated_at, 
			group_count, member_count, post_count
		FROM communities
		WHERE id = ?
	`

	var (
		idStr           string
		name            string
		description     string
		ownerID         string
		ownerName       string
		status          int
		communityType   string
		location        string
		createdAt       time.Time
		updatedAt       time.Time
		groupCount      int
		memberCount     int
		postCount       int
	)

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idStr, &name, &description, &ownerID, &ownerName, &status,
		&communityType, &location, &createdAt, &updatedAt,
		&groupCount, &memberCount, &postCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("社区不存在")
		}
		return nil, err
	}

	// 查询标签
	tags, err := r.findTags(ctx, idStr)
	if err != nil {
		return nil, err
	}

	// 重构社区实体
	return domain.ReconstructCommunity(
		idStr, name, description, ownerID, ownerName, status,
		tags, location, communityType,
		createdAt, updatedAt,
		groupCount, memberCount, postCount,
	)
}

// 查询标签
func (r *CommunityRepository) findTags(ctx context.Context, communityID string) ([]string, error) {
	query := `SELECT tag FROM community_tags WHERE community_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, communityID)
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

// FindByName 根据名称查找社区
func (r *CommunityRepository) FindByName(ctx context.Context, name domain.CommunityName) (*domain.Community, error) {
	query := `
		SELECT 
			id, name, description, owner_id, owner_name, status, 
			community_type, location, created_at, updated_at, 
			group_count, member_count, post_count
		FROM communities
		WHERE name = ?
	`

	var (
		idStr           string
		nameStr         string
		description     string
		ownerID         string
		ownerName       string
		status          int
		communityType   string
		location        string
		createdAt       time.Time
		updatedAt       time.Time
		groupCount      int
		memberCount     int
		postCount       int
	)

	err := r.db.QueryRowContext(ctx, query, name.Value()).Scan(
		&idStr, &nameStr, &description, &ownerID, &ownerName, &status,
		&communityType, &location, &createdAt, &updatedAt,
		&groupCount, &memberCount, &postCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("社区不存在")
		}
		return nil, err
	}

	// 查询标签
	tags, err := r.findTags(ctx, idStr)
	if err != nil {
		return nil, err
	}

	// 重构社区实体
	return domain.ReconstructCommunity(
		idStr, nameStr, description, ownerID, ownerName, status,
		tags, location, communityType,
		createdAt, updatedAt,
		groupCount, memberCount, postCount,
	)
}

// FindByOwner 查找用户拥有的社区
func (r *CommunityRepository) FindByOwner(ctx context.Context, ownerID domain.ID) ([]*domain.Community, error) {
	query := `
		SELECT 
			id, name, description, owner_id, owner_name, status, 
			community_type, location, created_at, updated_at, 
			group_count, member_count, post_count
		FROM communities
		WHERE owner_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []*domain.Community
	for rows.Next() {
		var (
			idStr           string
			name            string
			description     string
			ownerIDStr      string
			ownerName       string
			status          int
			communityType   string
			location        string
			createdAt       time.Time
			updatedAt       time.Time
			groupCount      int
			memberCount     int
			postCount       int
		)

		if err := rows.Scan(
			&idStr, &name, &description, &ownerIDStr, &ownerName, &status,
			&communityType, &location, &createdAt, &updatedAt,
			&groupCount, &memberCount, &postCount,
		); err != nil {
			return nil, err
		}

		// 查询标签
		tags, err := r.findTags(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 重构社区实体
		community, err := domain.ReconstructCommunity(
			idStr, name, description, ownerIDStr, ownerName, status,
			tags, location, communityType,
			createdAt, updatedAt,
			groupCount, memberCount, postCount,
		)

		if err != nil {
			return nil, err
		}

		communities = append(communities, community)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return communities, nil
}

// FindAll 查找所有社区（分页）
func (r *CommunityRepository) FindAll(ctx context.Context, offset, limit int) ([]*domain.Community, error) {
	query := `
		SELECT 
			id, name, description, owner_id, owner_name, status, 
			community_type, location, created_at, updated_at, 
			group_count, member_count, post_count
		FROM communities
		WHERE status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, int(domain.CommunityStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []*domain.Community
	for rows.Next() {
		var (
			idStr           string
			name            string
			description     string
			ownerID         string
			ownerName       string
			status          int
			communityType   string
			location        string
			createdAt       time.Time
			updatedAt       time.Time
			groupCount      int
			memberCount     int
			postCount       int
		)

		if err := rows.Scan(
			&idStr, &name, &description, &ownerID, &ownerName, &status,
			&communityType, &location, &createdAt, &updatedAt,
			&groupCount, &memberCount, &postCount,
		); err != nil {
			return nil, err
		}

		// 查询标签
		tags, err := r.findTags(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 重构社区实体
		community, err := domain.ReconstructCommunity(
			idStr, name, description, ownerID, ownerName, status,
			tags, location, communityType,
			createdAt, updatedAt,
			groupCount, memberCount, postCount,
		)

		if err != nil {
			return nil, err
		}

		communities = append(communities, community)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return communities, nil
}

// FindByType 根据类型查找社区
func (r *CommunityRepository) FindByType(ctx context.Context, type_ domain.CommunityType, offset, limit int) ([]*domain.Community, error) {
	query := `
		SELECT 
			id, name, description, owner_id, owner_name, status, 
			community_type, location, created_at, updated_at, 
			group_count, member_count, post_count
		FROM communities
		WHERE community_type = ? AND status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, type_.String(), int(domain.CommunityStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []*domain.Community
	for rows.Next() {
		var (
			idStr           string
			name            string
			description     string
			ownerID         string
			ownerName       string
			status          int
			communityType   string
			location        string
			createdAt       time.Time
			updatedAt       time.Time
			groupCount      int
			memberCount     int
			postCount       int
		)

		if err := rows.Scan(
			&idStr, &name, &description, &ownerID, &ownerName, &status,
			&communityType, &location, &createdAt, &updatedAt,
			&groupCount, &memberCount, &postCount,
		); err != nil {
			return nil, err
		}

		// 查询标签
		tags, err := r.findTags(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 重构社区实体
		community, err := domain.ReconstructCommunity(
			idStr, name, description, ownerID, ownerName, status,
			tags, location, communityType,
			createdAt, updatedAt,
			groupCount, memberCount, postCount,
		)

		if err != nil {
			return nil, err
		}

		communities = append(communities, community)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return communities, nil
}

// CountAll 统计所有社区数量
func (r *CommunityRepository) CountAll(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM communities WHERE status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, int(domain.CommunityStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountByType 统计特定类型的社区数量
func (r *CommunityRepository) CountByType(ctx context.Context, type_ domain.CommunityType) (int, error) {
	query := `SELECT COUNT(*) FROM communities WHERE community_type = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, type_.String(), int(domain.CommunityStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Delete 删除社区
func (r *CommunityRepository) Delete(ctx context.Context, id domain.ID) error {
	// 注意：这里只是逻辑删除，而不是物理删除
	query := `UPDATE communities SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, int(domain.CommunityStatusDeleted), time.Now(), id.Value())
	return err
}
