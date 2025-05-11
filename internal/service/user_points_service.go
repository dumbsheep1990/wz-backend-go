package service

import (
	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserPointsServiceImpl 用户积分服务实现
type UserPointsServiceImpl struct {
	repo domain.UserPointsRepository
}

// NewUserPointsService 创建用户积分服务
func NewUserPointsService(repo domain.UserPointsRepository) UserPointsService {
	return &UserPointsServiceImpl{
		repo: repo,
	}
}

// 通过依赖注入SQL连接创建服务
func NewUserPointsServiceWithConn(conn interface{}) UserPointsService {
	return &UserPointsServiceImpl{
		repo: sql.NewUserPointsRepository(conn),
	}
}

// CreatePoints 创建用户积分记录
func (s *UserPointsServiceImpl) CreatePoints(points *domain.UserPoints) (int64, error) {
	logx.Infof("创建用户积分记录: %+v", points)

	// 业务规则验证
	if err := s.validatePoints(points); err != nil {
		return 0, err
	}

	// 调用仓储层创建积分记录
	id, err := s.repo.Create(points)
	if err != nil {
		logx.Errorf("创建用户积分记录失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetPointsById 获取积分记录详情
func (s *UserPointsServiceImpl) GetPointsById(id int64) (*domain.UserPoints, error) {
	logx.Infof("获取积分记录详情: %d", id)

	points, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("获取积分记录详情失败: %v, id: %d", err, id)
		return nil, err
	}

	return points, nil
}

// ListPointsByUser 获取用户积分记录列表
func (s *UserPointsServiceImpl) ListPointsByUser(userID int64, page, pageSize int) ([]*domain.UserPoints, int64, error) {
	logx.Infof("获取用户积分记录列表: userID=%d, page=%d, pageSize=%d", userID, page, pageSize)

	// 参数验证
	if userID <= 0 {
		logx.Error("无效的用户ID")
		return nil, 0, domain.ErrInvalidParam
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用仓储层获取列表
	points, total, err := s.repo.ListByUser(userID, page, pageSize)
	if err != nil {
		logx.Errorf("获取用户积分记录列表失败: %v", err)
		return nil, 0, err
	}

	return points, total, nil
}

// GetUserTotalPoints 获取用户总积分
func (s *UserPointsServiceImpl) GetUserTotalPoints(userID int64) (int, error) {
	logx.Infof("获取用户总积分: userID=%d", userID)

	// 参数验证
	if userID <= 0 {
		logx.Error("无效的用户ID")
		return 0, domain.ErrInvalidParam
	}

	// 调用仓储层获取用户总积分
	total, err := s.repo.GetUserTotalPoints(userID)
	if err != nil {
		logx.Errorf("获取用户总积分失败: %v, userID: %d", err, userID)
		return 0, err
	}

	return total, nil
}

// validatePoints 验证积分记录
func (s *UserPointsServiceImpl) validatePoints(points *domain.UserPoints) error {
	// 参数验证
	if points.UserID <= 0 {
		logx.Error("无效的用户ID")
		return domain.ErrInvalidParam
	}

	if points.Points <= 0 {
		logx.Error("积分值必须大于0")
		return domain.ErrInvalidParam
	}

	if points.Type != 1 && points.Type != 2 {
		logx.Errorf("无效的积分类型: %d", points.Type)
		return domain.ErrInvalidParam
	}

	// 这里可以添加更多业务规则验证
	return nil
}
