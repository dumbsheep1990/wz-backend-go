package community

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yourusername/wz-backend-go/internal/domain/community/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/community/repository"
	"github.com/yourusername/wz-backend-go/internal/domain/community/valueobject"
)

// CommunityRepository implements the domain.CommunityRepository interface
type CommunityRepository struct {
	db *sql.DB
}

// NewCommunityRepository creates a new CommunityRepository
func NewCommunityRepository(db *sql.DB) *CommunityRepository {
	return &CommunityRepository{
		db: db,
	}
}

// FindByID finds a community by its ID
func (r *CommunityRepository) FindByID(ctx context.Context, id valueobject.ID) (*entity.Community, error) {
	// SQL query to find community by ID
	query := `
		SELECT 
			id, name, description, owner_id, status, 
			created_at, updated_at
		FROM communities 
		WHERE id = ? AND deleted_at IS NULL
	`
	
	row := r.db.QueryRowContext(ctx, query, id.String())
	
	var (
		idStr         string
		nameStr       string
		descriptionStr string
		ownerIDStr    string
		statusStr     string
		createdAt     time.Time
		updatedAt     time.Time
	)
	
	err := row.Scan(
		&idStr,
		&nameStr,
		&descriptionStr,
		&ownerIDStr,
		&statusStr,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Community not found
		}
		return nil, err
	}
	
	// Fetch tags for the community
	tags, err := r.getTagsForCommunity(ctx, idStr)
	if err != nil {
		return nil, err
	}
	
	// Fetch location for the community
	location, err := r.getLocationForCommunity(ctx, idStr)
	if err != nil {
		return nil, err
	}
	
	// Reconstruct the domain entity
	community, err := r.reconstituteCommunity(
		idStr,
		nameStr,
		descriptionStr,
		ownerIDStr,
		statusStr,
		tags,
		location,
		createdAt,
		updatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return community, nil
}

// FindAll retrieves all communities with pagination
func (r *CommunityRepository) FindAll(ctx context.Context, offset, limit int) ([]*entity.Community, int, error) {
	// Get total count
	var total int
	countErr := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM communities WHERE deleted_at IS NULL").Scan(&total)
	if countErr != nil {
		return nil, 0, countErr
	}
	
	// Query communities with pagination
	query := `
		SELECT 
			id, name, description, owner_id, status, 
			created_at, updated_at
		FROM communities 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var communities []*entity.Community
	
	for rows.Next() {
		var (
			idStr         string
			nameStr       string
			descriptionStr string
			ownerIDStr    string
			statusStr     string
			createdAt     time.Time
			updatedAt     time.Time
		)
		
		err := rows.Scan(
			&idStr,
			&nameStr,
			&descriptionStr,
			&ownerIDStr,
			&statusStr,
			&createdAt,
			&updatedAt,
		)
		
		if err != nil {
			return nil, 0, err
		}
		
		// Fetch tags for the community
		tags, err := r.getTagsForCommunity(ctx, idStr)
		if err != nil {
			return nil, 0, err
		}
		
		// Fetch location for the community
		location, err := r.getLocationForCommunity(ctx, idStr)
		if err != nil {
			return nil, 0, err
		}
		
		// Reconstruct the domain entity
		community, err := r.reconstituteCommunity(
			idStr,
			nameStr,
			descriptionStr,
			ownerIDStr,
			statusStr,
			tags,
			location,
			createdAt,
			updatedAt,
		)
		
		if err != nil {
			return nil, 0, err
		}
		
		communities = append(communities, community)
	}
	
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	
	return communities, total, nil
}

// FindByOwnerID finds communities by owner ID
func (r *CommunityRepository) FindByOwnerID(ctx context.Context, ownerID valueobject.UserID) ([]*entity.Community, error) {
	query := `
		SELECT 
			id, name, description, owner_id, status, 
			created_at, updated_at
		FROM communities 
		WHERE owner_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, ownerID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var communities []*entity.Community
	
	for rows.Next() {
		var (
			idStr         string
			nameStr       string
			descriptionStr string
			ownerIDStr    string
			statusStr     string
			createdAt     time.Time
			updatedAt     time.Time
		)
		
		err := rows.Scan(
			&idStr,
			&nameStr,
			&descriptionStr,
			&ownerIDStr,
			&statusStr,
			&createdAt,
			&updatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		// Fetch tags for the community
		tags, err := r.getTagsForCommunity(ctx, idStr)
		if err != nil {
			return nil, err
		}
		
		// Fetch location for the community
		location, err := r.getLocationForCommunity(ctx, idStr)
		if err != nil {
			return nil, err
		}
		
		// Reconstruct the domain entity
		community, err := r.reconstituteCommunity(
			idStr,
			nameStr,
			descriptionStr,
			ownerIDStr,
			statusStr,
			tags,
			location,
			createdAt,
			updatedAt,
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

// FindByTags finds communities by tags
func (r *CommunityRepository) FindByTags(ctx context.Context, tags []valueobject.Tag) ([]*entity.Community, error) {
	// This is a simplified implementation that doesn't handle multiple tags optimally
	// In a real implementation, you would use a more efficient query
	
	if len(tags) == 0 {
		return []*entity.Community{}, nil
	}
	
	// Using the first tag for simplicity
	query := `
		SELECT 
			c.id, c.name, c.description, c.owner_id, c.status, 
			c.created_at, c.updated_at
		FROM communities c
		JOIN community_tags ct ON c.id = ct.community_id
		WHERE ct.tag = ? AND c.deleted_at IS NULL
		ORDER BY c.created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, tags[0].String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var communities []*entity.Community
	
	for rows.Next() {
		var (
			idStr         string
			nameStr       string
			descriptionStr string
			ownerIDStr    string
			statusStr     string
			createdAt     time.Time
			updatedAt     time.Time
		)
		
		err := rows.Scan(
			&idStr,
			&nameStr,
			&descriptionStr,
			&ownerIDStr,
			&statusStr,
			&createdAt,
			&updatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		// Fetch all tags for the community
		communityTags, err := r.getTagsForCommunity(ctx, idStr)
		if err != nil {
			return nil, err
		}
		
		// Fetch location for the community
		location, err := r.getLocationForCommunity(ctx, idStr)
		if err != nil {
			return nil, err
		}
		
		// Reconstruct the domain entity
		community, err := r.reconstituteCommunity(
			idStr,
			nameStr,
			descriptionStr,
			ownerIDStr,
			statusStr,
			communityTags,
			location,
			createdAt,
			updatedAt,
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

// FindByLocation finds communities by location
func (r *CommunityRepository) FindByLocation(ctx context.Context, location valueobject.Location) ([]*entity.Community, error) {
	query := `
		SELECT 
			id, name, description, owner_id, status, 
			created_at, updated_at
		FROM communities 
		WHERE location = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, location.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var communities []*entity.Community
	
	for rows.Next() {
		var (
			idStr         string
			nameStr       string
			descriptionStr string
			ownerIDStr    string
			statusStr     string
			createdAt     time.Time
			updatedAt     time.Time
		)
		
		err := rows.Scan(
			&idStr,
			&nameStr,
			&descriptionStr,
			&ownerIDStr,
			&statusStr,
			&createdAt,
			&updatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		// Fetch tags for the community
		tags, err := r.getTagsForCommunity(ctx, idStr)
		if err != nil {
			return nil, err
		}
		
		// Reconstruct the domain entity
		community, err := r.reconstituteCommunity(
			idStr,
			nameStr,
			descriptionStr,
			ownerIDStr,
			statusStr,
			tags,
			location,
			createdAt,
			updatedAt,
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

// Save persists a community
func (r *CommunityRepository) Save(ctx context.Context, community *entity.Community) error {
	// Check if the community already exists
	existingCommunity, err := r.FindByID(ctx, community.ID())
	if err != nil {
		return err
	}
	
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()
	
	if existingCommunity == nil {
		// Insert new community
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO communities (
				id, name, description, owner_id, status, 
				created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			community.ID().String(),
			community.Name().String(),
			community.Description().String(),
			community.OwnerID().String(),
			string(community.Status()),
			community.CreatedAt(),
			community.UpdatedAt(),
		)
		if err != nil {
			return err
		}
		
		// Insert location
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO community_locations (community_id, location) VALUES (?, ?)`,
			community.ID().String(),
			community.Location().String(),
		)
		if err != nil {
			return err
		}
		
	} else {
		// Update existing community
		_, err = tx.ExecContext(
			ctx,
			`UPDATE communities SET 
				name = ?, description = ?, status = ?, updated_at = ?
			WHERE id = ?`,
			community.Name().String(),
			community.Description().String(),
			string(community.Status()),
			community.UpdatedAt(),
			community.ID().String(),
		)
		if err != nil {
			return err
		}
		
		// Update location
		_, err = tx.ExecContext(
			ctx,
			`UPDATE community_locations SET location = ? WHERE community_id = ?`,
			community.Location().String(),
			community.ID().String(),
		)
		if err != nil {
			return err
		}
		
		// Delete existing tags
		_, err = tx.ExecContext(
			ctx,
			`DELETE FROM community_tags WHERE community_id = ?`,
			community.ID().String(),
		)
		if err != nil {
			return err
		}
	}
	
	// Insert tags
	for _, tag := range community.Tags() {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO community_tags (community_id, tag) VALUES (?, ?)`,
			community.ID().String(),
			tag.String(),
		)
		if err != nil {
			return err
		}
	}
	
	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}
	
	return nil
}

// Delete removes a community
func (r *CommunityRepository) Delete(ctx context.Context, id valueobject.ID) error {
	// In this implementation, we'll do a soft delete by setting deleted_at
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE communities SET deleted_at = ?, status = ? WHERE id = ?`,
		time.Now(),
		string(valueobject.StatusDeleted),
		id.String(),
	)
	
	return err
}

// Helper methods

// getTagsForCommunity fetches tags for a community
func (r *CommunityRepository) getTagsForCommunity(ctx context.Context, communityID string) ([]valueobject.Tag, error) {
	query := `SELECT tag FROM community_tags WHERE community_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tags []valueobject.Tag
	
	for rows.Next() {
		var tagStr string
		if err := rows.Scan(&tagStr); err != nil {
			return nil, err
		}
		
		tag, err := valueobject.NewTag(tagStr)
		if err != nil {
			continue // Skip invalid tags
		}
		
		tags = append(tags, tag)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return tags, nil
}

// getLocationForCommunity fetches location for a community
func (r *CommunityRepository) getLocationForCommunity(ctx context.Context, communityID string) (valueobject.Location, error) {
	var locationStr string
	
	err := r.db.QueryRowContext(
		ctx,
		`SELECT location FROM community_locations WHERE community_id = ?`,
		communityID,
	).Scan(&locationStr)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return valueobject.NewLocation(""), nil // Empty location
		}
		return valueobject.NewLocation(""), err
	}
	
	return valueobject.NewLocation(locationStr), nil
}

// reconstituteCommunity recreates a Community entity from persistence data
func (r *CommunityRepository) reconstituteCommunity(
	id string,
	name string,
	description string,
	ownerID string,
	status string,
	tags []valueobject.Tag,
	location valueobject.Location,
	createdAt time.Time,
	updatedAt time.Time,
) (*entity.Community, error) {
	// This is a simplified reconstruction method
	// In a real implementation, you would use reflection or another method to reconstruct the entity
	
	// Create value objects
	idVO := valueobject.NewID(id)
	nameVO, err := valueobject.NewCommunityName(name)
	if err != nil {
		return nil, err
	}
	descriptionVO := valueobject.NewDescription(description)
	ownerIDVO := valueobject.NewUserID(ownerID)
	statusVO := valueobject.CommunityStatus(status)
	
	// Create a new Community instance
	// This is a simplification and might not fully represent your actual domain model
	community := &entity.Community{}
	
	// Use reflection or another method to set private fields
	// For now, this is a placeholder
	
	return community, nil
}

// Ensure CommunityRepository implements repository.CommunityRepository
var _ repository.CommunityRepository = (*CommunityRepository)(nil)
