package repository

import (
	"context"
	"wz-backend-go/internal/domain/component/entity"
	"wz-backend-go/internal/domain/component/valueobject"
)

// ComponentRepository 组件仓储接口
type ComponentRepository interface {
	// Save 保存组件
	Save(ctx context.Context, component *entity.Component) error
	
	// FindByID 根据ID查找组件
	FindByID(ctx context.Context, id valueobject.ComponentID) (*entity.Component, error)
	
	// FindByIDAndTenant 根据ID和租户ID查找组件
	FindByIDAndTenant(ctx context.Context, id valueobject.ComponentID, tenantID string) (*entity.Component, error)
	
	// FindByTenant 根据租户ID查找组件列表
	FindByTenant(ctx context.Context, tenantID string, filters ComponentFilters) ([]*entity.Component, error)
	
	// FindPublicComponents 查找公开组件
	FindPublicComponents(ctx context.Context, filters ComponentFilters) ([]*entity.Component, error)
	
	// FindByType 根据组件类型查找
	FindByType(ctx context.Context, componentType valueobject.ComponentType, filters ComponentFilters) ([]*entity.Component, error)
	
	// FindByCategory 根据分类查找组件
	FindByCategory(ctx context.Context, category string, filters ComponentFilters) ([]*entity.Component, error)
	
	// Delete 删除组件
	Delete(ctx context.Context, id valueobject.ComponentID) error
	
	// ExistsByName 检查组件名称是否已存在（在租户内）
	ExistsByName(ctx context.Context, name string, tenantID string) (bool, error)
	
	// ExistsByNameExcludeID 检查组件名称是否已存在（排除指定ID）
	ExistsByNameExcludeID(ctx context.Context, name string, tenantID string, excludeID valueobject.ComponentID) (bool, error)
	
	// Count 统计组件数量
	Count(ctx context.Context, tenantID string, filters ComponentFilters) (int64, error)
	
	// CountPublic 统计公开组件数量
	CountPublic(ctx context.Context, filters ComponentFilters) (int64, error)
}

// ComponentFilters 组件查询过滤器
type ComponentFilters struct {
	ComponentType string   // 组件类型过滤
	Category      string   // 分类过滤
	IsPublic      *bool    // 是否公开过滤
	Search        string   // 搜索关键词
	Tags          []string // 标签过滤
	TenantID      string   // 租户ID过滤
	Limit         int      // 限制数量
	Offset        int      // 偏移量
	SortBy        string   // 排序字段 (name, created_at, updated_at)
	SortOrder     string   // 排序顺序 (asc, desc)
} 