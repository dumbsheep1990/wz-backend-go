package repository

import (
	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/valueobject"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 保存用户
	Save(user *entity.User) error

	// 根据ID查找用户
	FindByID(id valueobject.UserID) (*entity.User, error)

	// 根据用户名查找用户
	FindByUsername(username valueobject.Username) (*entity.User, error)

	// 根据邮箱查找用户
	FindByEmail(email valueobject.Email) (*entity.User, error)

	// 根据手机号查找用户
	FindByPhone(phone valueobject.Phone) (*entity.User, error)

	// 分页查询用户列表
	FindAll(page, pageSize int) ([]*entity.User, int64, error)

	// 保存用户行为
	SaveBehavior(behavior *entity.UserBehavior) error

	// 查询用户行为列表
	FindBehaviorsByUserID(userID valueobject.UserID, page, pageSize int) ([]*entity.UserBehavior, int64, error)
}

// EventPublisher 事件发布接口
type EventPublisher interface {
	// 发布事件
	Publish(event interface{}) error
}
