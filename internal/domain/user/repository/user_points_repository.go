package repository

import (
	"context"
	"time"

	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/valueobject"
)

// SourceStats 表示来源统计
type SourceStats struct {
	Source string
	Count  int64
}

// UserPointsRepository 定义用户积分仓储接口
type UserPointsRepository interface {
	// 基本增删改查
	Save(ctx context.Context, userPoints *entity.UserPoints) error
	FindByID(ctx context.Context, id valueobject.ID) (*entity.UserPoints, error)
	FindByUserID(ctx context.Context, userID valueobject.UserID, offset, limit int64) ([]*entity.UserPoints, error)
	CountByUserID(ctx context.Context, userID valueobject.UserID) (int64, error)

	// 业务查询
	FindWithConditions(ctx context.Context, conditions map[string]interface{}, offset, limit int64) ([]*entity.UserPoints, error)
	CountWithConditions(ctx context.Context, conditions map[string]interface{}) (int64, error)
	GetTotalPointsByUserID(ctx context.Context, userID valueobject.UserID) (int, error)

	// 统计相关
	CountUsers(ctx context.Context) (int64, error)
	SumPoints(ctx context.Context) (int64, error)
	MaxPoints(ctx context.Context) (int64, error)
	SumPointsByConditions(ctx context.Context, conditions map[string]interface{}) (int64, error)
	GroupBySource(ctx context.Context) ([]*SourceStats, error)
}

// PointsRulesRepository 定义积分规则仓储接口
type PointsRulesRepository interface {
	// 基本增删改查
	Save(ctx context.Context, rules *entity.PointsRules) error
	FindByTenantID(ctx context.Context, tenantID valueobject.TenantID) (*entity.PointsRules, error)

	// 检查用户今日获取的积分
	GetUserDailyPoints(ctx context.Context, userID valueobject.UserID, date time.Time) (int, error)
}
