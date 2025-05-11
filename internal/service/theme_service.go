package service

import (
	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// ThemeServiceImpl 主题服务实现
type ThemeServiceImpl struct {
	repo domain.ThemeRepository
}

// NewThemeService 创建主题服务
func NewThemeService(repo domain.ThemeRepository) ThemeService {
	return &ThemeServiceImpl{
		repo: repo,
	}
}

// 通过依赖注入SQL连接创建服务
func NewThemeServiceWithConn(conn interface{}) ThemeService {
	return &ThemeServiceImpl{
		repo: sql.NewThemeRepository(conn),
	}
}

// CreateTheme 创建主题
func (s *ThemeServiceImpl) CreateTheme(theme *domain.Theme) (int64, error) {
	logx.Infof("创建主题: %+v", theme)

	// 业务规则验证
	if err := s.validateTheme(theme); err != nil {
		return 0, err
	}

	// 调用仓储层创建主题
	id, err := s.repo.Create(theme)
	if err != nil {
		logx.Errorf("创建主题失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetThemeById 获取主题详情
func (s *ThemeServiceImpl) GetThemeById(id int64) (*domain.Theme, error) {
	logx.Infof("获取主题详情: %d", id)

	theme, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("获取主题详情失败: %v, id: %d", err, id)
		return nil, err
	}

	return theme, nil
}

// ListThemes 获取主题列表
func (s *ThemeServiceImpl) ListThemes(page, pageSize int, query map[string]interface{}) ([]*domain.Theme, int64, error) {
	logx.Infof("获取主题列表: page=%d, pageSize=%d, query=%v", page, pageSize, query)

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用仓储层获取列表
	themes, total, err := s.repo.List(page, pageSize, query)
	if err != nil {
		logx.Errorf("获取主题列表失败: %v", err)
		return nil, 0, err
	}

	return themes, total, nil
}

// UpdateTheme 更新主题
func (s *ThemeServiceImpl) UpdateTheme(theme *domain.Theme) error {
	logx.Infof("更新主题: %+v", theme)

	// 首先检查主题是否存在
	_, err := s.repo.GetByID(theme.ID)
	if err != nil {
		logx.Errorf("主题不存在: %v, id: %d", err, theme.ID)
		return err
	}

	// 业务规则验证
	if err := s.validateTheme(theme); err != nil {
		return err
	}

	// 调用仓储层更新主题
	err = s.repo.Update(theme)
	if err != nil {
		logx.Errorf("更新主题失败: %v", err)
		return err
	}

	return nil
}

// DeleteTheme 删除主题
func (s *ThemeServiceImpl) DeleteTheme(id int64) error {
	logx.Infof("删除主题: %d", id)

	// 首先检查主题是否存在
	theme, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("主题不存在: %v, id: %d", err, id)
		return err
	}

	// 不允许删除默认主题
	if theme.IsDefault == 1 {
		logx.Errorf("不能删除默认主题: id=%d", id)
		return domain.ErrOperationForbidden
	}

	// 调用仓储层删除主题
	err = s.repo.Delete(id)
	if err != nil {
		logx.Errorf("删除主题失败: %v, id: %d", err, id)
		return err
	}

	return nil
}

// SetDefaultTheme 设置默认主题
func (s *ThemeServiceImpl) SetDefaultTheme(id int64, tenantID int64) error {
	logx.Infof("设置默认主题: id=%d, tenantID=%d", id, tenantID)

	// 首先检查主题是否存在
	theme, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("主题不存在: %v, id: %d", err, id)
		return err
	}

	// 检查主题是否属于指定租户
	if theme.TenantID != tenantID {
		logx.Errorf("主题不属于该租户: themeID=%d, themeTenantID=%d, specifiedTenantID=%d",
			id, theme.TenantID, tenantID)
		return domain.ErrOperationForbidden
	}

	// 检查主题状态是否为启用
	if theme.Status != 1 {
		logx.Errorf("无法将禁用状态的主题设为默认: id=%d", id)
		return domain.ErrOperationForbidden
	}

	// 调用仓储层设置默认主题
	err = s.repo.SetDefault(id, tenantID)
	if err != nil {
		logx.Errorf("设置默认主题失败: %v", err)
		return err
	}

	return nil
}

// GetDefaultTheme 获取默认主题
func (s *ThemeServiceImpl) GetDefaultTheme(tenantID int64) (*domain.Theme, error) {
	logx.Infof("获取默认主题: tenantID=%d", tenantID)

	// 调用仓储层获取默认主题
	theme, err := s.repo.GetDefault(tenantID)
	if err != nil {
		logx.Errorf("获取默认主题失败: %v, tenantID: %d", err, tenantID)
		return nil, err
	}

	return theme, nil
}

// validateTheme 验证主题
func (s *ThemeServiceImpl) validateTheme(theme *domain.Theme) error {
	// 这里可以添加业务规则验证
	// 例如：主题名称不能为空、主题代码必须唯一等
	return nil
}
