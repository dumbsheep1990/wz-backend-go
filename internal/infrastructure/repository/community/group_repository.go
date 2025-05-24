package community

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	domain "wz-backend-go/internal/domain/community"
)

// GroupRepository 小组仓储实现
type GroupRepository struct {
	db *sql.DB
}

// NewGroupRepository 创建小组仓储
func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

// Save 保存小组
func (r *GroupRepository) Save(ctx context.Context, group *domain.Group) error {
	// 检查小组是否已存在
	existingGroup, err := r.FindByID(ctx, group.ID())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if existingGroup == nil {
		// 新增小组
		query := `
			INSERT INTO community_groups (
				id, name, description, community_id, owner_id, owner_name, 
				status, created_at, updated_at, member_count, post_count
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			group.ID().Value(),
			group.Name().Value(),
			group.Description().Value(),
			group.CommunityID().Value(),
			group.OwnerID().Value(),
			group.OwnerName(),
			int(group.Status()),
			group.CreatedAt().Value(),
			group.UpdatedAt().Value(),
			group.MemberCount(),
			group.PostCount(),
		)

		if err != nil {
			return err
		}

		// 存储标签
		if len(group.Tags().Values()) > 0 {
			if err := r.saveTags(ctx, group.ID().Value(), group.Tags().Values()); err != nil {
				return err
			}
		}

		// 存储成员列表
		if len(group.Members()) > 0 {
			if err := r.saveMembers(ctx, group.ID().Value(), group.Members()); err != nil {
				return err
			}
		}

		return nil
	} else {
		// 更新小组
		query := `
			UPDATE community_groups SET
				name = ?,
				description = ?,
				status = ?,
				updated_at = ?,
				member_count = ?,
				post_count = ?
			WHERE id = ?
		`

		_, err = r.db.ExecContext(
			ctx,
			query,
			group.Name().Value(),
			group.Description().Value(),
			int(group.Status()),
			group.UpdatedAt().Value(),
			group.MemberCount(),
			group.PostCount(),
			group.ID().Value(),
		)

		if err != nil {
			return err
		}

		// 更新标签（先删除旧标签，再添加新标签）
		if err := r.deleteTags(ctx, group.ID().Value()); err != nil {
			return err
		}

		if len(group.Tags().Values()) > 0 {
			if err := r.saveTags(ctx, group.ID().Value(), group.Tags().Values()); err != nil {
				return err
			}
		}

		// 更新成员列表（先删除旧成员，再添加新成员）
		if err := r.deleteMembers(ctx, group.ID().Value()); err != nil {
			return err
		}

		if len(group.Members()) > 0 {
			if err := r.saveMembers(ctx, group.ID().Value(), group.Members()); err != nil {
				return err
			}
		}

		return nil
	}
}

// 保存标签
func (r *GroupRepository) saveTags(ctx context.Context, groupID string, tags []string) error {
	query := `INSERT INTO group_tags (group_id, tag) VALUES (?, ?)`
	
	for _, tag := range tags {
		_, err := r.db.ExecContext(ctx, query, groupID, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除标签
func (r *GroupRepository) deleteTags(ctx context.Context, groupID string) error {
	query := `DELETE FROM group_tags WHERE group_id = ?`
	_, err := r.db.ExecContext(ctx, query, groupID)
	return err
}

// 保存成员
func (r *GroupRepository) saveMembers(ctx context.Context, groupID string, members []domain.ID) error {
	query := `INSERT INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)`
	
	for _, memberID := range members {
		_, err := r.db.ExecContext(ctx, query, groupID, memberID.Value(), time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除成员
func (r *GroupRepository) deleteMembers(ctx context.Context, groupID string) error {
	query := `DELETE FROM group_members WHERE group_id = ?`
	_, err := r.db.ExecContext(ctx, query, groupID)
	return err
}

// FindByID 根据ID查找小组
func (r *GroupRepository) FindByID(ctx context.Context, id domain.ID) (*domain.Group, error) {
	query := `
		SELECT 
			id, name, description, community_id, owner_id, owner_name, 
			status, created_at, updated_at, member_count, post_count
		FROM community_groups
		WHERE id = ?
	`

	var (
		idStr        string
		name         string
		description  string
		communityID  string
		ownerID      string
		ownerName    string
		status       int
		createdAt    time.Time
		updatedAt    time.Time
		memberCount  int
		postCount    int
	)

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&idStr, &name, &description, &communityID, &ownerID, &ownerName,
		&status, &createdAt, &updatedAt, &memberCount, &postCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("小组不存在")
		}
		return nil, err
	}

	// 查询标签
	tags, err := r.findTags(ctx, idStr)
	if err != nil {
		return nil, err
	}

	// 查询成员
	members, err := r.findMembers(ctx, idStr)
	if err != nil {
		return nil, err
	}

	// 重构小组实体
	return domain.ReconstructGroup(
		idStr, name, description, communityID, ownerID, ownerName,
		status, members, tags, createdAt, updatedAt, memberCount, postCount,
	)
}

// 查询标签
func (r *GroupRepository) findTags(ctx context.Context, groupID string) ([]string, error) {
	query := `SELECT tag FROM group_tags WHERE group_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, groupID)
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

// 查询成员
func (r *GroupRepository) findMembers(ctx context.Context, groupID string) ([]string, error) {
	query := `SELECT user_id FROM group_members WHERE group_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var members []string
	for rows.Next() {
		var memberID string
		if err := rows.Scan(&memberID); err != nil {
			return nil, err
		}
		members = append(members, memberID)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return members, nil
}

// FindByCommunity 查找社区下的小组（分页）
func (r *GroupRepository) FindByCommunity(ctx context.Context, communityID domain.ID, offset, limit int) ([]*domain.Group, error) {
	query := `
		SELECT 
			id, name, description, community_id, owner_id, owner_name, 
			status, created_at, updated_at, member_count, post_count
		FROM community_groups
		WHERE community_id = ? AND status != ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, communityID.Value(), int(domain.GroupStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.Group
	for rows.Next() {
		var (
			idStr        string
			name         string
			description  string
			communityIDStr string
			ownerID      string
			ownerName    string
			status       int
			createdAt    time.Time
			updatedAt    time.Time
			memberCount  int
			postCount    int
		)

		if err := rows.Scan(
			&idStr, &name, &description, &communityIDStr, &ownerID, &ownerName,
			&status, &createdAt, &updatedAt, &memberCount, &postCount,
		); err != nil {
			return nil, err
		}

		// 查询标签
		tags, err := r.findTags(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 查询成员
		members, err := r.findMembers(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 重构小组实体
		group, err := domain.ReconstructGroup(
			idStr, name, description, communityIDStr, ownerID, ownerName,
			status, members, tags, createdAt, updatedAt, memberCount, postCount,
		)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// CountByCommunity 统计社区下的小组数量
func (r *GroupRepository) CountByCommunity(ctx context.Context, communityID domain.ID) (int, error) {
	query := `SELECT COUNT(*) FROM community_groups WHERE community_id = ? AND status != ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, communityID.Value(), int(domain.GroupStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByMember 查找用户加入的小组
func (r *GroupRepository) FindByMember(ctx context.Context, memberID domain.ID, offset, limit int) ([]*domain.Group, error) {
	query := `
		SELECT 
			g.id, g.name, g.description, g.community_id, g.owner_id, g.owner_name, 
			g.status, g.created_at, g.updated_at, g.member_count, g.post_count
		FROM community_groups g
		JOIN group_members m ON g.id = m.group_id
		WHERE m.user_id = ? AND g.status != ?
		ORDER BY g.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, memberID.Value(), int(domain.GroupStatusDeleted), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.Group
	for rows.Next() {
		var (
			idStr          string
			name           string
			description    string
			communityID    string
			ownerID        string
			ownerName      string
			status         int
			createdAt      time.Time
			updatedAt      time.Time
			memberCount    int
			postCount      int
		)

		if err := rows.Scan(
			&idStr, &name, &description, &communityID, &ownerID, &ownerName,
			&status, &createdAt, &updatedAt, &memberCount, &postCount,
		); err != nil {
			return nil, err
		}

		// 查询标签
		tags, err := r.findTags(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 查询成员
		members, err := r.findMembers(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 重构小组实体
		group, err := domain.ReconstructGroup(
			idStr, name, description, communityID, ownerID, ownerName,
			status, members, tags, createdAt, updatedAt, memberCount, postCount,
		)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// CountByMember 统计用户加入的小组数量
func (r *GroupRepository) CountByMember(ctx context.Context, memberID domain.ID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM community_groups g
		JOIN group_members m ON g.id = m.group_id
		WHERE m.user_id = ? AND g.status != ?
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, memberID.Value(), int(domain.GroupStatusDeleted)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindByOwner 查找用户拥有的小组
func (r *GroupRepository) FindByOwner(ctx context.Context, ownerID domain.ID) ([]*domain.Group, error) {
	query := `
		SELECT 
			id, name, description, community_id, owner_id, owner_name, 
			status, created_at, updated_at, member_count, post_count
		FROM community_groups
		WHERE owner_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.Group
	for rows.Next() {
		var (
			idStr        string
			name         string
			description  string
			communityID  string
			ownerIDStr   string
			ownerName    string
			status       int
			createdAt    time.Time
			updatedAt    time.Time
			memberCount  int
			postCount    int
		)

		if err := rows.Scan(
			&idStr, &name, &description, &communityID, &ownerIDStr, &ownerName,
			&status, &createdAt, &updatedAt, &memberCount, &postCount,
		); err != nil {
			return nil, err
		}

		// 查询标签
		tags, err := r.findTags(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 查询成员
		members, err := r.findMembers(ctx, idStr)
		if err != nil {
			return nil, err
		}

		// 重构小组实体
		group, err := domain.ReconstructGroup(
			idStr, name, description, communityID, ownerIDStr, ownerName,
			status, members, tags, createdAt, updatedAt, memberCount, postCount,
		)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// Delete 删除小组
func (r *GroupRepository) Delete(ctx context.Context, id domain.ID) error {
	// 注意：这里只是逻辑删除，而不是物理删除
	query := `UPDATE community_groups SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, int(domain.GroupStatusDeleted), time.Now(), id.Value())
	return err
}
