package repository

import (
	"context"
	"time"
	"wz-backend-go/internal/domain/render/entity"
	"wz-backend-go/internal/domain/render/valueobject"
)

// RenderResultRepository 渲染结果仓储接口
type RenderResultRepository interface {
	// Save 保存渲染结果
	Save(ctx context.Context, result *entity.RenderResult) error
	
	// FindByID 根据ID查找渲染结果
	FindByID(ctx context.Context, id valueobject.RenderID) (*entity.RenderResult, error)
	
	// FindByCacheKey 根据缓存键查找渲染结果
	FindByCacheKey(ctx context.Context, cacheKey string) (*entity.RenderResult, error)
	
	// Delete 删除渲染结果
	Delete(ctx context.Context, id valueobject.RenderID) error
	
	// DeleteBySiteID 删除指定站点的所有渲染结果
	DeleteBySiteID(ctx context.Context, siteID string) error
	
	// DeleteByGroup 删除指定组的所有渲染结果
	DeleteByGroup(ctx context.Context, group string) error
	
	// DeleteExpired 删除所有过期的渲染结果
	DeleteExpired(ctx context.Context) error
	
	// FindExpiring 查找即将过期的渲染结果
	FindExpiring(ctx context.Context, within time.Duration) ([]*entity.RenderResult, error)
}

// TemplateRepository 模板仓储接口
type TemplateRepository interface {
	// Save 保存模板
	Save(ctx context.Context, template *entity.Template) error
	
	// FindByID 根据ID查找模板
	FindByID(ctx context.Context, id string) (*entity.Template, error)
	
	// FindByName 根据名称查找模板
	FindByName(ctx context.Context, name string, siteID string) (*entity.Template, error)
	
	// FindBySiteIDAndType 根据站点ID和类型查找模板
	FindBySiteIDAndType(ctx context.Context, siteID string, templateType string) ([]*entity.Template, error)
	
	// FindAll 查找所有模板
	FindAll(ctx context.Context, offset int, limit int) ([]*entity.Template, int, error)
	
	// Delete 删除模板
	Delete(ctx context.Context, id string) error
	
	// DeleteBySiteID 删除指定站点的所有模板
	DeleteBySiteID(ctx context.Context, siteID string) error
	
	// FindTemplatesByNamePattern 根据名称模式查找模板
	FindTemplatesByNamePattern(ctx context.Context, pattern string, siteID string) ([]*entity.Template, error)
}
