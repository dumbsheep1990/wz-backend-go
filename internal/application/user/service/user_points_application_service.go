package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"wz-backend-go/internal/application/user/dto"
	"wz-backend-go/internal/domain/shared/event"
	"wz-backend-go/internal/domain/user/entity"
	"wz-backend-go/internal/domain/user/repository"
	"wz-backend-go/internal/domain/user/valueobject"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	FindUserNameByID(ctx context.Context, id int64) (string, error)
}

// UserPointsApplicationService 用户积分应用服务
type UserPointsApplicationService struct {
	userPointsRepo  repository.UserPointsRepository
	pointsRulesRepo repository.PointsRulesRepository
	userRepo        UserRepository
	eventBus        event.EventBus
	validator       *validator.Validate
	unitOfWork      database.UnitOfWork
}

// NewUserPointsApplicationService 创建用户积分应用服务
func NewUserPointsApplicationService(
	userPointsRepo repository.UserPointsRepository,
	pointsRulesRepo repository.PointsRulesRepository,
	userRepo UserRepository,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) *UserPointsApplicationService {
	return &UserPointsApplicationService{
		userPointsRepo:  userPointsRepo,
		pointsRulesRepo: pointsRulesRepo,
		userRepo:        userRepo,
		eventBus:        eventBus,
		validator:       validator.New(),
		unitOfWork:      unitOfWork,
	}
}

// CreatePoints 创建积分记录
func (s *UserPointsApplicationService) CreatePoints(ctx context.Context, req dto.CreatePointsRequest) (*dto.PointsDTO, error) {
	// 验证请求数据
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// 转换请求为值对象
	userID := valueobject.NewUserID(req.UserID)

	points, err := valueobject.NewPoints(req.Points)
	if err != nil {
		return nil, err
	}

	pointsType, err := valueobject.NewPointsType(req.Type)
	if err != nil {
		return nil, err
	}

	source, err := valueobject.NewSource(req.Source)
	if err != nil {
		return nil, err
	}

	description, err := valueobject.NewDescription(req.Description)
	if err != nil {
		return nil, err
	}

	relatedType := valueobject.NewRelatedType(req.RelatedType)

	operatorID := valueobject.NewUserID(req.OperatorID)

	tenantID := valueobject.NewTenantID(req.TenantID)

	// 如果是积分减少，检查用户积分是否足够
	if pointsType == valueobject.PointsTypeDecrease {
		totalPoints, err := s.userPointsRepo.GetTotalPointsByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if totalPoints < points.Value() {
			return nil, entity.ErrInsufficientPoints
		}
	}

	// 检查每日积分上限（仅对增加积分操作）
	if pointsType == valueobject.PointsTypeIncrease && source != valueobject.SourceAdmin {
		// 获取积分规则
		rules, err := s.pointsRulesRepo.FindByTenantID(ctx, tenantID)
		if err != nil {
			return nil, err
		}

		if rules != nil && rules.MaxDailyPoints() > 0 {
			// 获取用户今日已获取的积分
			today := time.Now()
			dailyPoints, err := s.pointsRulesRepo.GetUserDailyPoints(ctx, userID, today)
			if err != nil {
				return nil, err
			}

			// 如果加上本次积分将超过每日上限，则拒绝
			if dailyPoints+points.Value() > rules.MaxDailyPoints() {
				return nil, fmt.Errorf("超过每日积分上限 %d，今日已获取 %d", rules.MaxDailyPoints(), dailyPoints)
			}
		}
	}

	// 创建积分实体
	userPoints, err := entity.NewUserPoints(
		userID,
		points,
		pointsType,
		source,
		description,
		req.RelatedID,
		relatedType,
		operatorID,
		tenantID,
	)
	if err != nil {
		return nil, err
	}

	// 使用工作单元保存
	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		// 保存积分记录
		if err := s.userPointsRepo.Save(ctx, userPoints); err != nil {
			return err
		}

		// 发布领域事件
		events := userPoints.GetDomainEvents()
		for _, event := range events {
			s.eventBus.Publish(ctx, event)
		}
		userPoints.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 获取用户名和操作员名称
	username, _ := s.userRepo.FindUserNameByID(ctx, req.UserID)
	operatorName := ""
	if req.OperatorID > 0 {
		operatorName, _ = s.userRepo.FindUserNameByID(ctx, req.OperatorID)
	}

	// 返回DTO
	return dto.NewPointsDTOFromEntity(userPoints, username, operatorName), nil
}

// GetPointByID 根据ID获取积分记录
func (s *UserPointsApplicationService) GetPointByID(ctx context.Context, id string) (*dto.PointsDTO, error) {
	// 将ID转换为值对象
	pointID := valueobject.NewID(id)

	// 查询积分记录
	userPoints, err := s.userPointsRepo.FindByID(ctx, pointID)
	if err != nil {
		return nil, err
	}

	if userPoints == nil {
		return nil, errors.New("积分记录不存在")
	}

	// 获取用户名和操作员名称
	username, _ := s.userRepo.FindUserNameByID(ctx, userPoints.UserID().Value())
	operatorName := ""
	if userPoints.OperatorID().Value() > 0 {
		operatorName, _ = s.userRepo.FindUserNameByID(ctx, userPoints.OperatorID().Value())
	}

	// 返回DTO
	return dto.NewPointsDTOFromEntity(userPoints, username, operatorName), nil
}

// GetTotalPointsByUserID 获取用户总积分
func (s *UserPointsApplicationService) GetTotalPointsByUserID(ctx context.Context, userID int64) (int, error) {
	userIDVO := valueobject.NewUserID(userID)
	return s.userPointsRepo.GetTotalPointsByUserID(ctx, userIDVO)
}

// ListPointsByUserID 获取用户积分记录列表
func (s *UserPointsApplicationService) ListPointsByUserID(ctx context.Context, userID int64, page, pageSize int64) (*dto.ListPointsResponse, error) {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	userIDVO := valueobject.NewUserID(userID)
	offset := (page - 1) * pageSize

	// 查询积分记录
	points, err := s.userPointsRepo.FindByUserID(ctx, userIDVO, offset, pageSize)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := s.userPointsRepo.CountByUserID(ctx, userIDVO)
	if err != nil {
		return nil, err
	}

	// 获取用户名
	username, _ := s.userRepo.FindUserNameByID(ctx, userID)

	// 构建返回DTO
	items := make([]*dto.PointsDTO, 0, len(points))
	for _, point := range points {
		operatorName := ""
		if point.OperatorID().Value() > 0 {
			operatorName, _ = s.userRepo.FindUserNameByID(ctx, point.OperatorID().Value())
		}
		items = append(items, dto.NewPointsDTOFromEntity(point, username, operatorName))
	}

	return &dto.ListPointsResponse{
		Items: items,
		Total: total,
	}, nil
}

// ListPointsWithConditions 根据条件查询积分记录
func (s *UserPointsApplicationService) ListPointsWithConditions(ctx context.Context, req dto.ListPointsRequest) (*dto.ListPointsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	} else if req.PageSize > 100 {
		req.PageSize = 100
	}

	offset := (req.Page - 1) * req.PageSize

	// 构建查询条件
	conditions := make(map[string]interface{})
	if req.UserID > 0 {
		conditions["user_id"] = req.UserID
	}
	if req.Username != "" {
		conditions["username"] = req.Username
	}
	if req.Type > 0 {
		conditions["type"] = req.Type
	}
	if req.Source != "" {
		conditions["source"] = req.Source
	}
	if req.StartDate != "" {
		conditions["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		conditions["end_date"] = req.EndDate
	}

	// 查询积分记录
	points, err := s.userPointsRepo.FindWithConditions(ctx, conditions, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := s.userPointsRepo.CountWithConditions(ctx, conditions)
	if err != nil {
		return nil, err
	}

	// 构建返回DTO
	items := make([]*dto.PointsDTO, 0, len(points))
	for _, point := range points {
		username, _ := s.userRepo.FindUserNameByID(ctx, point.UserID().Value())
		operatorName := ""
		if point.OperatorID().Value() > 0 {
			operatorName, _ = s.userRepo.FindUserNameByID(ctx, point.OperatorID().Value())
		}
		items = append(items, dto.NewPointsDTOFromEntity(point, username, operatorName))
	}

	return &dto.ListPointsResponse{
		Items: items,
		Total: total,
	}, nil
}

// RevokePoint 撤销积分记录
func (s *UserPointsApplicationService) RevokePoint(ctx context.Context, id string, operatorID int64) error {
	// 将ID转换为值对象
	pointID := valueobject.NewID(id)
	operatorIDVO := valueobject.NewUserID(operatorID)

	// 查询积分记录
	userPoints, err := s.userPointsRepo.FindByID(ctx, pointID)
	if err != nil {
		return err
	}

	if userPoints == nil {
		return errors.New("积分记录不存在")
	}

	// 只允许撤销管理员创建的积分记录
	if userPoints.Source() != valueobject.SourceAdmin {
		return errors.New("只能撤销管理员调整的积分记录")
	}

	// 使用工作单元撤销
	return s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		// 撤销积分记录
		if err := userPoints.Revoke(operatorIDVO); err != nil {
			return err
		}

		// 保存变更
		if err := s.userPointsRepo.Save(ctx, userPoints); err != nil {
			return err
		}

		// 发布领域事件
		events := userPoints.GetDomainEvents()
		for _, event := range events {
			s.eventBus.Publish(ctx, event)
		}
		userPoints.ClearDomainEvents()

		return nil
	})
}

// GetPointsStatistics 获取积分统计数据
func (s *UserPointsApplicationService) GetPointsStatistics(ctx context.Context) (*dto.PointsStatisticsResponse, error) {
	// 获取基础统计数据
	totalUsers, err := s.userPointsRepo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}

	totalPoints, err := s.userPointsRepo.SumPoints(ctx)
	if err != nil {
		return nil, err
	}

	avgPoints := int64(0)
	if totalUsers > 0 {
		avgPoints = totalPoints / totalUsers
	}

	maxPoints, err := s.userPointsRepo.MaxPoints(ctx)
	if err != nil {
		return nil, err
	}

	// 获取今日统计
	today := time.Now().Format("2006-01-02")
	todayIncrease, err := s.userPointsRepo.SumPointsByConditions(ctx, map[string]interface{}{
		"type":       1,
		"start_date": today,
		"end_date":   today,
	})
	if err != nil {
		return nil, err
	}

	todayDecrease, err := s.userPointsRepo.SumPointsByConditions(ctx, map[string]interface{}{
		"type":       2,
		"start_date": today,
		"end_date":   today,
	})
	if err != nil {
		return nil, err
	}

	// 获取本月统计
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	monthIncrease, err := s.userPointsRepo.SumPointsByConditions(ctx, map[string]interface{}{
		"type":       1,
		"start_date": firstDay,
		"end_date":   lastDay,
	})
	if err != nil {
		return nil, err
	}

	monthDecrease, err := s.userPointsRepo.SumPointsByConditions(ctx, map[string]interface{}{
		"type":       2,
		"start_date": firstDay,
		"end_date":   lastDay,
	})
	if err != nil {
		return nil, err
	}

	// 获取来源分布
	sourceDistribution, err := s.userPointsRepo.GroupBySource(ctx)
	if err != nil {
		return nil, err
	}

	sourceStats := make([]*dto.SourceStats, 0, len(sourceDistribution))
	for _, src := range sourceDistribution {
		source, _ := valueobject.NewSource(src.Source)
		sourceStats = append(sourceStats, &dto.SourceStats{
			Source: src.Source,
			Name:   source.DisplayName(),
			Count:  src.Count,
		})
	}

	return &dto.PointsStatisticsResponse{
		TotalUsers:    totalUsers,
		TotalPoints:   totalPoints,
		AvgPoints:     avgPoints,
		MaxPoints:     maxPoints,
		TodayIncrease: todayIncrease,
		TodayDecrease: todayDecrease,
		MonthIncrease: monthIncrease,
		MonthDecrease: monthDecrease,
		SourceStats:   sourceStats,
	}, nil
}

// GetPointsRules 获取积分规则
func (s *UserPointsApplicationService) GetPointsRules(ctx context.Context, tenantID int64) (*dto.PointsRulesResponse, error) {
	tenantIDVO := valueobject.NewTenantID(tenantID)

	// 查询积分规则
	rules, err := s.pointsRulesRepo.FindByTenantID(ctx, tenantIDVO)
	if err != nil {
		return nil, err
	}

	if rules == nil {
		return nil, errors.New("积分规则不存在")
	}

	// 返回DTO
	return dto.NewPointsRulesDTOFromEntity(rules), nil
}

// UpdatePointsRules 更新积分规则
func (s *UserPointsApplicationService) UpdatePointsRules(ctx context.Context, req dto.PointsRulesRequest) (*dto.PointsRulesResponse, error) {
	// 验证请求数据
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	tenantIDVO := valueobject.NewTenantID(req.TenantID)

	// 查询现有规则
	rules, err := s.pointsRulesRepo.FindByTenantID(ctx, tenantIDVO)
	if err != nil {
		return nil, err
	}

	// 如果规则不存在，创建新规则
	if rules == nil {
		rules, err = entity.NewPointsRules(
			req.SignInPoints,
			req.CommentPoints,
			req.SharePoints,
			req.ArticlePoints,
			req.InvitePoints,
			req.PurchaseRate,
			req.MaxDailyPoints,
			req.EnableExchange,
			req.ExchangeRate,
			req.MinExchangePoints,
			tenantIDVO,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// 更新现有规则
		err = rules.Update(
			req.SignInPoints,
			req.CommentPoints,
			req.SharePoints,
			req.ArticlePoints,
			req.InvitePoints,
			req.PurchaseRate,
			req.MaxDailyPoints,
			req.EnableExchange,
			req.ExchangeRate,
			req.MinExchangePoints,
		)
		if err != nil {
			return nil, err
		}
	}

	// 使用工作单元保存
	err = s.unitOfWork.Execute(ctx, func(ctx context.Context) error {
		// 保存规则
		if err := s.pointsRulesRepo.Save(ctx, rules); err != nil {
			return err
		}

		// 发布领域事件
		events := rules.GetDomainEvents()
		for _, event := range events {
			s.eventBus.Publish(ctx, event)
		}
		rules.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 返回DTO
	return dto.NewPointsRulesDTOFromEntity(rules), nil
}
