package service

import (
	"errors"
	"time"

	"github.com/wz-project/wz-backend-go/services/community-service/models"
	"gorm.io/gorm"
)

// SimilarApplicationService 同乡申请服务
type SimilarApplicationService struct {
	DB *gorm.DB
}

// NewSimilarApplicationService 创建新的同乡申请服务
func NewSimilarApplicationService(db *gorm.DB) *SimilarApplicationService {
	return &SimilarApplicationService{
		DB: db,
	}
}

// CreateApplication 创建新的申请
func (s *SimilarApplicationService) CreateApplication(application *models.SimilarApplication) error {
	// 验证必填字段
	if application.Name == "" {
		return errors.New("姓名不能为空")
	}

	if application.ApplicationType == "" {
		return errors.New("申请类型不能为空")
	}

	// 设置默认状态
	application.Status = "pending"
	application.CreatedAt = time.Now()
	application.UpdatedAt = time.Now()

	// 保存到数据库
	return s.DB.Create(application).Error
}

// GetApplicationByID 通过ID获取申请
func (s *SimilarApplicationService) GetApplicationByID(id string) (*models.SimilarApplication, error) {
	var application models.SimilarApplication
	err := s.DB.Where("id = ?", id).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

// GetApplicationsByUserID 获取用户的所有申请
func (s *SimilarApplicationService) GetApplicationsByUserID(userID string) ([]models.SimilarApplication, error) {
	var applications []models.SimilarApplication
	err := s.DB.Where("user_id = ?", userID).Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}

// ListApplications 列出申请（支持分页和过滤）
func (s *SimilarApplicationService) ListApplications(page, pageSize int, filters map[string]interface{}) ([]models.SimilarApplication, int64, error) {
	var applications []models.SimilarApplication
	var total int64

	query := s.DB.Model(&models.SimilarApplication{})

	// 应用过滤条件
	for key, value := range filters {
		if value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&applications).Error
	if err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// UpdateApplication 更新申请
func (s *SimilarApplicationService) UpdateApplication(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return s.DB.Model(&models.SimilarApplication{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateApplicationStatus 更新申请状态
func (s *SimilarApplicationService) UpdateApplicationStatus(id, status string) error {
	return s.DB.Model(&models.SimilarApplication{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// DeleteApplication 删除申请
func (s *SimilarApplicationService) DeleteApplication(id string) error {
	return s.DB.Where("id = ?", id).Delete(&models.SimilarApplication{}).Error
}
