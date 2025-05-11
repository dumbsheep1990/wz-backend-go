package service

import (
	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserMessageServiceImpl 用户消息服务实现
type UserMessageServiceImpl struct {
	repo domain.UserMessageRepository
}

// NewUserMessageService 创建用户消息服务
func NewUserMessageService(repo domain.UserMessageRepository) UserMessageService {
	return &UserMessageServiceImpl{
		repo: repo,
	}
}

// 通过依赖注入SQL连接创建服务
func NewUserMessageServiceWithConn(conn interface{}) UserMessageService {
	return &UserMessageServiceImpl{
		repo: sql.NewUserMessageRepository(conn),
	}
}

// CreateMessage 创建用户消息
func (s *UserMessageServiceImpl) CreateMessage(message *domain.UserMessage) (int64, error) {
	logx.Infof("创建用户消息: %+v", message)

	// 业务规则验证
	if err := s.validateMessage(message); err != nil {
		return 0, err
	}

	// 设置默认值
	if message.Status == 0 {
		message.Status = 0 // 默认未读
	}

	// 调用仓储层创建消息
	id, err := s.repo.Create(message)
	if err != nil {
		logx.Errorf("创建用户消息失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetMessageById 获取消息详情
func (s *UserMessageServiceImpl) GetMessageById(id int64) (*domain.UserMessage, error) {
	logx.Infof("获取消息详情: %d", id)

	message, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("获取消息详情失败: %v, id: %d", err, id)
		return nil, err
	}

	return message, nil
}

// ListMessagesByUser 获取用户消息列表
func (s *UserMessageServiceImpl) ListMessagesByUser(userID int64, page, pageSize int, query map[string]interface{}) ([]*domain.UserMessage, int64, error) {
	logx.Infof("获取用户消息列表: userID=%d, page=%d, pageSize=%d, query=%v", userID, page, pageSize, query)

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
	messages, total, err := s.repo.ListByUser(userID, page, pageSize, query)
	if err != nil {
		logx.Errorf("获取用户消息列表失败: %v", err)
		return nil, 0, err
	}

	return messages, total, nil
}

// MarkAsRead 标记消息为已读
func (s *UserMessageServiceImpl) MarkAsRead(id int64, userID int64) error {
	logx.Infof("标记消息为已读: id=%d, userID=%d", id, userID)

	// 参数验证
	if id <= 0 || userID <= 0 {
		logx.Error("无效的参数")
		return domain.ErrInvalidParam
	}

	// 首先检查消息是否存在且属于该用户
	message, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("消息不存在: %v, id: %d", err, id)
		return err
	}

	if message.UserID != userID {
		logx.Errorf("消息不属于该用户: messageID=%d, messageUserID=%d, specifiedUserID=%d",
			id, message.UserID, userID)
		return domain.ErrOperationForbidden
	}

	// 如果消息已读，则不需要再次标记
	if message.Status == 1 {
		return nil
	}

	// 调用仓储层标记消息为已读
	err = s.repo.MarkAsRead(id, userID)
	if err != nil {
		logx.Errorf("标记消息为已读失败: %v, id: %d, userID: %d", err, id, userID)
		return err
	}

	return nil
}

// MarkAllAsRead 标记用户所有消息为已读
func (s *UserMessageServiceImpl) MarkAllAsRead(userID int64) error {
	logx.Infof("标记所有消息为已读: userID=%d", userID)

	// 参数验证
	if userID <= 0 {
		logx.Error("无效的用户ID")
		return domain.ErrInvalidParam
	}

	// 调用仓储层标记所有消息为已读
	err := s.repo.MarkAllAsRead(userID)
	if err != nil {
		logx.Errorf("标记所有消息为已读失败: %v, userID: %d", err, userID)
		return err
	}

	return nil
}

// DeleteMessage 删除用户消息
func (s *UserMessageServiceImpl) DeleteMessage(id int64, userID int64) error {
	logx.Infof("删除用户消息: id=%d, userID=%d", id, userID)

	// 参数验证
	if id <= 0 || userID <= 0 {
		logx.Error("无效的参数")
		return domain.ErrInvalidParam
	}

	// 首先检查消息是否存在且属于该用户
	message, err := s.repo.GetByID(id)
	if err != nil {
		logx.Errorf("消息不存在: %v, id: %d", err, id)
		return err
	}

	if message.UserID != userID {
		logx.Errorf("消息不属于该用户: messageID=%d, messageUserID=%d, specifiedUserID=%d",
			id, message.UserID, userID)
		return domain.ErrOperationForbidden
	}

	// 调用仓储层删除消息
	err = s.repo.Delete(id, userID)
	if err != nil {
		logx.Errorf("删除用户消息失败: %v, id: %d, userID: %d", err, id, userID)
		return err
	}

	return nil
}

// CountUnread 统计用户未读消息数量
func (s *UserMessageServiceImpl) CountUnread(userID int64) (int64, error) {
	logx.Infof("统计用户未读消息数量: userID=%d", userID)

	// 参数验证
	if userID <= 0 {
		logx.Error("无效的用户ID")
		return 0, domain.ErrInvalidParam
	}

	// 调用仓储层统计未读消息
	count, err := s.repo.CountUnread(userID)
	if err != nil {
		logx.Errorf("统计用户未读消息数量失败: %v, userID: %d", err, userID)
		return 0, err
	}

	return count, nil
}

// validateMessage 验证消息
func (s *UserMessageServiceImpl) validateMessage(message *domain.UserMessage) error {
	// 这里可以添加业务规则验证
	// 例如：消息类型必须有效、内容不能为空等
	return nil
}
