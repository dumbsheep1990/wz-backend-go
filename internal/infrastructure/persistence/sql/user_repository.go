package sql

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"

	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/repository"
	"wz-backend-go/internal/domain/user/valueobject"
)

// UserPO 用户持久化对象
type UserPO struct {
	ID                int64        `gorm:"primaryKey;column:id"`
	Username          string       `gorm:"uniqueIndex;column:username;size:50;not null"`
	Password          string       `gorm:"column:password;size:255;not null"`
	Email             string       `gorm:"uniqueIndex;column:email;size:100;not null"`
	Phone             string       `gorm:"uniqueIndex;column:phone;size:20"`
	Status            int32        `gorm:"column:status;not null;default:1"`
	IsVerified        bool         `gorm:"column:is_verified;not null;default:false"`
	IsCompanyVerified bool         `gorm:"column:is_company_verified;not null;default:false"`
	CreatedAt         time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time    `gorm:"column:updated_at;not null"`
	DeletedAt         sql.NullTime `gorm:"column:deleted_at"`
}

// TableName 表名
func (UserPO) TableName() string {
	return "users"
}

// UserBehaviorPO 用户行为持久化对象
type UserBehaviorPO struct {
	ID           int64     `gorm:"primaryKey;column:id"`
	UserID       int64     `gorm:"index;column:user_id;not null"`
	Action       string    `gorm:"column:action;size:50;not null"`
	ResourceType string    `gorm:"column:resource_type;size:50;not null"`
	ResourceID   int64     `gorm:"column:resource_id;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
}

// TableName 表名
func (UserBehaviorPO) TableName() string {
	return "user_behaviors"
}

// UserRepository 用户仓储实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Save 保存用户
func (r *UserRepository) Save(user *entity.User) error {
	// 如果用户ID为0，则为新用户，需要插入
	if user.ID().Value() == 0 {
		po := &UserPO{
			Username:          user.Username().Value(),
			Password:          user.Password(),
			Email:             user.Email().Value(),
			Phone:             user.Phone().Value(),
			Status:            user.Status().Value(),
			IsVerified:        user.IsVerified(),
			IsCompanyVerified: user.IsCompanyVerified(),
			CreatedAt:         user.CreatedAt(),
			UpdatedAt:         user.UpdatedAt(),
		}

		if err := r.db.Create(po).Error; err != nil {
			return err
		}

		// 设置用户ID
		userID := valueobject.NewUserID(po.ID)
		user.SetID(userID)

		return nil
	}

	// 否则为更新用户
	po := &UserPO{
		ID:                user.ID().Value(),
		Username:          user.Username().Value(),
		Password:          user.Password(),
		Email:             user.Email().Value(),
		Phone:             user.Phone().Value(),
		Status:            user.Status().Value(),
		IsVerified:        user.IsVerified(),
		IsCompanyVerified: user.IsCompanyVerified(),
		UpdatedAt:         user.UpdatedAt(),
	}

	return r.db.Model(&UserPO{}).Where("id = ?", po.ID).Updates(map[string]interface{}{
		"username":            po.Username,
		"password":            po.Password,
		"email":               po.Email,
		"phone":               po.Phone,
		"status":              po.Status,
		"is_verified":         po.IsVerified,
		"is_company_verified": po.IsCompanyVerified,
		"updated_at":          po.UpdatedAt,
	}).Error
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id valueobject.UserID) (*entity.User, error) {
	var po UserPO
	if err := r.db.Where("id = ?", id.Value()).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomainEntity(&po)
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username valueobject.Username) (*entity.User, error) {
	var po UserPO
	if err := r.db.Where("username = ?", username.Value()).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomainEntity(&po)
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(email valueobject.Email) (*entity.User, error) {
	var po UserPO
	if err := r.db.Where("email = ?", email.Value()).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomainEntity(&po)
}

// FindByPhone 根据手机号查找用户
func (r *UserRepository) FindByPhone(phone valueobject.Phone) (*entity.User, error) {
	var po UserPO
	if err := r.db.Where("phone = ?", phone.Value()).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomainEntity(&po)
}

// FindAll 分页查询用户列表
func (r *UserRepository) FindAll(page, pageSize int) ([]*entity.User, int64, error) {
	var count int64
	if err := r.db.Model(&UserPO{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var pos []UserPO
	if err := r.db.Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(pos))
	for i, po := range pos {
		user, err := r.toDomainEntity(&po)
		if err != nil {
			return nil, 0, err
		}
		users[i] = user
	}

	return users, count, nil
}

// SaveBehavior 保存用户行为
func (r *UserRepository) SaveBehavior(behavior *entity.UserBehavior) error {
	// 如果行为ID为0，则为新行为，需要插入
	if behavior.ID().Value() == 0 {
		po := &UserBehaviorPO{
			UserID:       behavior.UserID().Value(),
			Action:       behavior.Action(),
			ResourceType: behavior.ResourceType(),
			ResourceID:   behavior.ResourceID().Value(),
			CreatedAt:    behavior.CreatedAt(),
		}

		if err := r.db.Create(po).Error; err != nil {
			return err
		}

		// 设置行为ID
		behaviorID := valueobject.NewUserID(po.ID)
		behavior.SetID(behaviorID)

		return nil
	}

	// 行为一般不需要更新
	return nil
}

// FindBehaviorsByUserID 查询用户行为列表
func (r *UserRepository) FindBehaviorsByUserID(userID valueobject.UserID, page, pageSize int) ([]*entity.UserBehavior, int64, error) {
	var count int64
	if err := r.db.Model(&UserBehaviorPO{}).Where("user_id = ?", userID.Value()).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var pos []UserBehaviorPO
	if err := r.db.Where("user_id = ?", userID.Value()).Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	behaviors := make([]*entity.UserBehavior, len(pos))
	for i, po := range pos {
		behavior, err := r.toBehaviorDomainEntity(&po)
		if err != nil {
			return nil, 0, err
		}
		behaviors[i] = behavior
	}

	return behaviors, count, nil
}

// 将持久化对象转换为领域实体
func (r *UserRepository) toDomainEntity(po *UserPO) (*entity.User, error) {
	userID := valueobject.NewUserID(po.ID)

	username, err := valueobject.NewUsername(po.Username)
	if err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(po.Email)
	if err != nil {
		return nil, err
	}

	phone, err := valueobject.NewPhone(po.Phone)
	if err != nil {
		return nil, err
	}

	status, err := valueobject.NewUserStatus(po.Status)
	if err != nil {
		return nil, err
	}

	return entity.ReconstructUser(
		userID,
		username,
		po.Password,
		email,
		phone,
		status,
		po.IsVerified,
		po.IsCompanyVerified,
		po.CreatedAt,
		po.UpdatedAt,
	), nil
}

// 将行为持久化对象转换为领域实体
func (r *UserRepository) toBehaviorDomainEntity(po *UserBehaviorPO) (*entity.UserBehavior, error) {
	behaviorID := valueobject.NewUserID(po.ID)
	userID := valueobject.NewUserID(po.UserID)
	resourceID := valueobject.NewUserID(po.ResourceID)

	return entity.ReconstructUserBehavior(
		behaviorID,
		userID,
		po.Action,
		po.ResourceType,
		resourceID,
		po.CreatedAt,
	), nil
}
