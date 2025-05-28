package repository

import (
	"context"

	"github.com/yourusername/wz-backend-go/internal/domain/community/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/community/valueobject"
)

// CommunityRepository 定义社区持久化操作的接口
type CommunityRepository interface {
	// FindByID 通过ID查找社区
	FindByID(ctx context.Context, id valueobject.ID) (*entity.Community, error)
	
	// FindAll 获取所有社区，支持可选的分页
	FindAll(ctx context.Context, offset, limit int) ([]*entity.Community, int, error)
	
	// FindByOwnerID 通过创建者ID查找社区
	FindByOwnerID(ctx context.Context, ownerID valueobject.UserID) ([]*entity.Community, error)
	
	// FindByTags 通过标签查找社区
	FindByTags(ctx context.Context, tags []valueobject.Tag) ([]*entity.Community, error)
	
	// FindByLocation 通过位置查找社区
	FindByLocation(ctx context.Context, location valueobject.Location) ([]*entity.Community, error)
	
	// Save 持久化保存社区
	Save(ctx context.Context, community *entity.Community) error
	
	// Delete 删除社区
	Delete(ctx context.Context, id valueobject.ID) error
}
