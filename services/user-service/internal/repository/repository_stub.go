package repository

import (
	"context"
	"time"

	"wz-backend-go/services/user-service/internal/model"
)

// NewUserRepository 创建用户存储库的存根实现
func NewUserRepository() UserRepository {
	return &stubUserRepository{}
}

// stubUserRepository 存根实现，用于编译通过
type stubUserRepository struct{}

func (r *stubUserRepository) Create(ctx context.Context, user model.User) (int64, error) {
	return 1, nil
}

func (r *stubUserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return &model.User{
		ID:        id,
		Username:  "test_user",
		Email:     "test@example.com",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *stubUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return &model.User{
		ID:        1,
		Username:  username,
		Email:     "test@example.com",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *stubUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return &model.User{
		ID:        1,
		Username:  "test_user",
		Email:     email,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *stubUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return false, nil
}

func (r *stubUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (r *stubUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (r *stubUserRepository) GetUserBehaviors(ctx context.Context, userID int64, startTime, endTime time.Time) ([]model.UserBehavior, error) {
	return []model.UserBehavior{}, nil
}
