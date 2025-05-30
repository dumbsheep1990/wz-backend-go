package repository

import (
	"context"
	"wz-backend-go/internal/domain/site/entity"
	"wz-backend-go/internal/domain/site/valueobject"
)

// SiteRepository 站点仓储接口
type SiteRepository interface {
	// Save 保存站点
	Save(ctx context.Context, site *entity.Site) error
	
	// FindByID 根据ID查找站点
	FindByID(ctx context.Context, id valueobject.SiteID) (*entity.Site, error)
	
	// FindByIDAndTenant 根据ID和租户ID查找站点
	FindByIDAndTenant(ctx context.Context, id valueobject.SiteID, tenantID string) (*entity.Site, error)
	
	// FindByTenant 根据租户ID查找站点列表
	FindByTenant(ctx context.Context, tenantID string, filters SiteFilters) ([]*entity.Site, error)
	
	// FindByDomain 根据域名查找站点
	FindByDomain(ctx context.Context, domain valueobject.Domain) (*entity.Site, error)
	
	// Delete 删除站点
	Delete(ctx context.Context, id valueobject.SiteID) error
	
	// ExistsByDomain 检查域名是否已存在
	ExistsByDomain(ctx context.Context, domain valueobject.Domain) (bool, error)
	
	// ExistsByDomainExcludeID 检查域名是否已存在（排除指定ID）
	ExistsByDomainExcludeID(ctx context.Context, domain valueobject.Domain, excludeID valueobject.SiteID) (bool, error)
	
	// Count 统计站点数量
	Count(ctx context.Context, tenantID string, filters SiteFilters) (int64, error)
}

// SiteFilters 站点查询过滤器
type SiteFilters struct {
	Status string // 状态过滤
	Search string // 搜索关键词
	Limit  int    // 限制数量
	Offset int    // 偏移量
} 