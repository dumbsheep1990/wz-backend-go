package service

import (
	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// LinkServiceImpl 友情链接服务实现
type LinkServiceImpl struct {
	repo domain.LinkRepository
}

// NewLinkService 创建友情链接服务
func NewLinkService(repo domain.LinkRepository) LinkService {
	return &LinkServiceImpl{
		repo: repo,
	}
}

// 通过依赖注入SQL连接创建服务
func NewLinkServiceWithConn(conn interface{}) LinkService {
	return &LinkServiceImpl{
		repo: sql.NewLinkRepository(conn),
	}
}

// CreateLink 创建友情链接
func (s *LinkServiceImpl) CreateLink(link *domain.Link) (int64, error) {
	logx.Infof("创建友情链接: %+v", link)

	// 业务规则验证
	if err := s.validateLink(link); err != nil {
		return 0, err
	}

	// 调用仓储层创建链接
	id, err := s.repo.Create(link)
	if err != nil {
		logx.Errorf("创建友情链接失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetLinkById 获取友情链接详情
func (s *LinkServiceImpl) GetLinkById(id int64) (*domain.Link, error) {
	logx.Infof("获取友情链接详情: %d", id)

	link, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("获取友情链接详情失败: %v, id: %d", err, id)
		return nil, err
	}

	return link, nil
}

// ListLinks 获取友情链接列表
func (s *LinkServiceImpl) ListLinks(page, pageSize int, query map[string]interface{}) ([]*domain.Link, int64, error) {
	logx.Infof("获取友情链接列表: page=%d, pageSize=%d, query=%v", page, pageSize, query)

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用仓储层获取列表
	links, total, err := s.repo.List(page, pageSize, query)
	if err != nil {
		logx.Errorf("获取友情链接列表失败: %v", err)
		return nil, 0, err
	}

	return links, total, nil
}

// UpdateLink 更新友情链接
func (s *LinkServiceImpl) UpdateLink(link *domain.Link) error {
	logx.Infof("更新友情链接: %+v", link)

	// 首先检查链接是否存在
	_, err := s.repo.GetByID(link.ID)
	if err != nil {
		logx.Errorf("友情链接不存在: %v, id: %d", err, link.ID)
		return err
	}

	// 业务规则验证
	if err := s.validateLink(link); err != nil {
		return err
	}

	// 调用仓储层更新链接
	err = s.repo.Update(link)
	if err != nil {
		logx.Errorf("更新友情链接失败: %v", err)
		return err
	}

	return nil
}

// DeleteLink 删除友情链接
func (s *LinkServiceImpl) DeleteLink(id int64) error {
	logx.Infof("删除友情链接: %d", id)

	// 首先检查链接是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("友情链接不存在: %v, id: %d", err, id)
		return err
	}

	// 调用仓储层删除链接
	err = s.repo.Delete(id)
	if err != nil {
		logx.Errorf("删除友情链接失败: %v, id: %d", err, id)
		return err
	}

	return nil
}

// validateLink 验证友情链接
func (s *LinkServiceImpl) validateLink(link *domain.Link) error {
	// 这里可以添加业务规则验证
	// 例如：链接名称不能为空、URL必须有效等
	return nil
}
