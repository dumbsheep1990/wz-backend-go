package dto

import (
	"time"

	"wz-backend-go/internal/domain/user/entity"
)

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Status            int32     `json:"status"`
	StatusDesc        string    `json:"statusDesc"`
	IsVerified        bool      `json:"isVerified"`
	IsCompanyVerified bool      `json:"isCompanyVerified"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// UserBehaviorDTO 用户行为数据传输对象
type UserBehaviorDTO struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resourceType"`
	ResourceID   int64     `json:"resourceId"`
	CreatedAt    time.Time `json:"createdAt"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=30"`
}

// UserBehaviorRequest 记录用户行为请求
type UserBehaviorRequest struct {
	Action       string `json:"action" binding:"required"`
	ResourceType string `json:"resourceType" binding:"required"`
	ResourceID   int64  `json:"resourceId" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    *UserDTO `json:"data,omitempty"`
}

// UsersResponse 用户列表响应
type UsersResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    *UsersData `json:"data,omitempty"`
}

// UsersData 用户列表数据
type UsersData struct {
	Total int64      `json:"total"`
	List  []*UserDTO `json:"list"`
}

// UserBehaviorsResponse 用户行为列表响应
type UserBehaviorsResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    *UserBehaviorsData `json:"data,omitempty"`
}

// UserBehaviorsData 用户行为列表数据
type UserBehaviorsData struct {
	Total int64              `json:"total"`
	List  []*UserBehaviorDTO `json:"list"`
}

// 将领域实体转换为DTO
func ToUserDTO(user *entity.User) *UserDTO {
	return &UserDTO{
		ID:                user.ID().Value(),
		Username:          user.Username().Value(),
		Email:             user.Email().Value(),
		Phone:             user.Phone().Value(),
		Status:            user.Status().Value(),
		StatusDesc:        user.Status().String(),
		IsVerified:        user.IsVerified(),
		IsCompanyVerified: user.IsCompanyVerified(),
		CreatedAt:         user.CreatedAt(),
		UpdatedAt:         user.UpdatedAt(),
	}
}

// 将用户行为领域实体转换为DTO
func ToUserBehaviorDTO(behavior *entity.UserBehavior) *UserBehaviorDTO {
	return &UserBehaviorDTO{
		ID:           behavior.ID().Value(),
		UserID:       behavior.UserID().Value(),
		Action:       behavior.Action(),
		ResourceType: behavior.ResourceType(),
		ResourceID:   behavior.ResourceID().Value(),
		CreatedAt:    behavior.CreatedAt(),
	}
}

// 将用户DTO列表转换为响应
func ToUsersResponse(users []*entity.User, total int64) *UsersResponse {
	userDTOs := make([]*UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = ToUserDTO(user)
	}

	return &UsersResponse{
		Code:    0,
		Message: "success",
		Data: &UsersData{
			Total: total,
			List:  userDTOs,
		},
	}
}

// 将用户行为DTO列表转换为响应
func ToUserBehaviorsResponse(behaviors []*entity.UserBehavior, total int64) *UserBehaviorsResponse {
	behaviorDTOs := make([]*UserBehaviorDTO, len(behaviors))
	for i, behavior := range behaviors {
		behaviorDTOs[i] = ToUserBehaviorDTO(behavior)
	}

	return &UserBehaviorsResponse{
		Code:    0,
		Message: "success",
		Data: &UserBehaviorsData{
			Total: total,
			List:  behaviorDTOs,
		},
	}
}
