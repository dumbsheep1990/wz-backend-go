package service

import (
	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserFavoriteServiceImpl 用户收藏服务实现
type UserFavoriteServiceImpl struct {
	repo domain.UserFavoriteRepository
}

// NewUserFavoriteService 创建用户收藏服务
func NewUserFavoriteService(repo domain.UserFavoriteRepository) UserFavoriteService {
	return &UserFavoriteServiceImpl{
		repo: repo,
	}
}

// 通过依赖注入SQL连接创建服务
func NewUserFavoriteServiceWithConn(conn interface{}) UserFavoriteService {
	return &UserFavoriteServiceImpl{
		repo: sql.NewUserFavoriteRepository(conn),
	}
}

// CreateFavorite 创建用户收藏
func (s *UserFavoriteServiceImpl) CreateFavorite(favorite *domain.UserFavorite) (int64, error) {
	logx.Infof("创建用户收藏: %+v", favorite)

	// 业务规则验证
	if err := s.validateFavorite(favorite); err != nil {
		return 0, err
	}

	// 检查是否已收藏
	exists, err := s.repo.CheckFavorite(favorite.UserID, favorite.ItemID, favorite.ItemType)
	if err != nil {
		logx.Errorf("检查是否已收藏失败: %v", err)
		return 0, err
	}

	if exists {
		logx.Infof("用户已收藏该项目: userID=%d, itemID=%d, itemType=%s",
			favorite.UserID, favorite.ItemID, favorite.ItemType)
		return 0, domain.ErrAlreadyExists
	}

	// 调用仓储层创建收藏
	id, err := s.repo.Create(favorite)
	if err != nil {
		logx.Errorf("创建用户收藏失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetFavoriteById 获取收藏详情
func (s *UserFavoriteServiceImpl) GetFavoriteById(id int64) (*domain.UserFavorite, error) {
	logx.Infof("获取收藏详情: %d", id)

	favorite, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("获取收藏详情失败: %v, id: %d", err, id)
		return nil, err
	}

	return favorite, nil
}

// ListFavoritesByUser 获取用户收藏列表
func (s *UserFavoriteServiceImpl) ListFavoritesByUser(userID int64, page, pageSize int, itemType string) ([]*domain.UserFavorite, int64, error) {
	logx.Infof("获取用户收藏列表: userID=%d, page=%d, pageSize=%d, itemType=%s", userID, page, pageSize, itemType)

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
	favorites, total, err := s.repo.ListByUser(userID, page, pageSize, itemType)
	if err != nil {
		logx.Errorf("获取用户收藏列表失败: %v", err)
		return nil, 0, err
	}

	return favorites, total, nil
}

// DeleteFavorite 删除用户收藏
func (s *UserFavoriteServiceImpl) DeleteFavorite(id int64, userID int64) error {
	logx.Infof("删除用户收藏: id=%d, userID=%d", id, userID)

	// 参数验证
	if id <= 0 || userID <= 0 {
		logx.Error("无效的参数")
		return domain.ErrInvalidParam
	}

	// 首先检查收藏是否存在且属于该用户
	favorite, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("收藏不存在: %v, id: %d", err, id)
		return err
	}

	if favorite.UserID != userID {
		logx.Errorf("收藏不属于该用户: favoriteID=%d, favoriteUserID=%d, specifiedUserID=%d",
			id, favorite.UserID, userID)
		return domain.ErrOperationForbidden
	}

	// 调用仓储层删除收藏
	err = s.repo.Delete(id, userID)
	if err != nil {
		logx.Errorf("删除用户收藏失败: %v, id: %d, userID: %d", err, id, userID)
		return err
	}

	return nil
}

// CheckFavorite 检查是否已收藏
func (s *UserFavoriteServiceImpl) CheckFavorite(userID int64, itemID int64, itemType string) (bool, error) {
	logx.Infof("检查是否已收藏: userID=%d, itemID=%d, itemType=%s", userID, itemID, itemType)

	// 参数验证
	if userID <= 0 || itemID <= 0 || itemType == "" {
		logx.Error("无效的参数")
		return false, domain.ErrInvalidParam
	}

	// 调用仓储层检查是否已收藏
	exists, err := s.repo.CheckFavorite(userID, itemID, itemType)
	if err != nil {
		logx.Errorf("检查是否已收藏失败: %v", err)
		return false, err
	}

	return exists, nil
}

// validateFavorite 验证收藏
func (s *UserFavoriteServiceImpl) validateFavorite(favorite *domain.UserFavorite) error {
	// 参数验证
	if favorite.UserID <= 0 {
		logx.Error("无效的用户ID")
		return domain.ErrInvalidParam
	}

	if favorite.ItemID <= 0 {
		logx.Error("无效的项目ID")
		return domain.ErrInvalidParam
	}

	if favorite.ItemType == "" {
		logx.Error("项目类型不能为空")
		return domain.ErrInvalidParam
	}

	if favorite.Title == "" {
		logx.Error("项目标题不能为空")
		return domain.ErrInvalidParam
	}

	// 这里可以添加更多业务规则验证
	return nil
}
