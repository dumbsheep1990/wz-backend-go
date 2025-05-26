package service

import (
	"context"
	"errors"
	"time"

	"github.com/wz-project/wz-backend-go/internal/domain"
	"github.com/wz-project/wz-backend-go/internal/domain/model"
)

// TemplateService 实现模板服务接口
type TemplateService struct {
	repo domain.TemplateRepository
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(repo domain.TemplateRepository) *TemplateService {
	return &TemplateService{
		repo: repo,
	}
}

// GetTemplates 获取模板列表
func (s *TemplateService) GetTemplates(ctx context.Context, userID int64, page, pageSize int) ([]*model.Template, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.FindAll(ctx, userID, page, pageSize)
}

// GetTemplate 获取单个模板
func (s *TemplateService) GetTemplate(ctx context.Context, templateID int64) (*model.Template, error) {
	return s.repo.FindByID(ctx, templateID)
}

// CreateTemplate 创建模板
func (s *TemplateService) CreateTemplate(ctx context.Context, template *model.Template) error {
	if template.Name == "" {
		return errors.New("模板名称不能为空")
	}

	// 默认值设置
	if template.Type == "" {
		template.Type = model.TemplateBanner
	}
	template.Enabled = true
	template.IsNew = true
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	return s.repo.Create(ctx, template)
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(ctx context.Context, template *model.Template) error {
	if template.ID <= 0 {
		return errors.New("无效的模板ID")
	}

	if template.Name == "" {
		return errors.New("模板名称不能为空")
	}

	// 检查模板是否存在
	existingTemplate, err := s.repo.FindByID(ctx, template.ID)
	if err != nil {
		return err
	}

	// 检查是否有权限更新
	if existingTemplate.UserID != template.UserID {
		return errors.New("无权操作此模板")
	}

	template.UpdatedAt = time.Now()
	return s.repo.Update(ctx, template)
}

// DeleteTemplate 删除模板
func (s *TemplateService) DeleteTemplate(ctx context.Context, templateID int64) error {
	if templateID <= 0 {
		return errors.New("无效的模板ID")
	}

	return s.repo.Delete(ctx, templateID)
}

// UpdateTemplateStatus 更新模板状态
func (s *TemplateService) UpdateTemplateStatus(ctx context.Context, templateID int64, enabled bool) error {
	if templateID <= 0 {
		return errors.New("无效的模板ID")
	}

	// 检查模板是否存在
	template, err := s.repo.FindByID(ctx, templateID)
	if err != nil {
		return err
	}

	// 更新状态
	template.Enabled = enabled
	template.UpdatedAt = time.Now()

	return s.repo.Update(ctx, template)
}
