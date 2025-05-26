package domain

import (
	"context"
	"time"

	"github.com/wz-project/wz-backend-go/internal/domain/model"
)

// TemplateService 定义模板服务接口
type TemplateService interface {
	// 获取模板列表
	GetTemplates(ctx context.Context, userID int64, page, pageSize int) ([]*model.Template, int, error)
	
	// 获取单个模板
	GetTemplate(ctx context.Context, templateID int64) (*model.Template, error)
	
	// 创建模板
	CreateTemplate(ctx context.Context, template *model.Template) error
	
	// 更新模板
	UpdateTemplate(ctx context.Context, template *model.Template) error
	
	// 删除模板
	DeleteTemplate(ctx context.Context, templateID int64) error
	
	// 更新模板状态
	UpdateTemplateStatus(ctx context.Context, templateID int64, enabled bool) error
}

// TemplateRepository 定义模板仓储接口
type TemplateRepository interface {
	// 获取模板列表
	FindAll(ctx context.Context, userID int64, page, pageSize int) ([]*model.Template, int, error)
	
	// 获取单个模板
	FindByID(ctx context.Context, templateID int64) (*model.Template, error)
	
	// 创建模板
	Create(ctx context.Context, template *model.Template) error
	
	// 更新模板
	Update(ctx context.Context, template *model.Template) error
	
	// 删除模板
	Delete(ctx context.Context, templateID int64) error
}
