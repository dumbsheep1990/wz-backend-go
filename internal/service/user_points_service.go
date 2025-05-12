package service

import (
	"errors"
	"fmt"
	"time"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserPointsServiceImpl 用户积分服务实现
type UserPointsServiceImpl struct {
	userPointsRepo domain.UserPointsRepository
	userRepo       repository.UserRepository
}

// NewUserPointsService 创建用户积分服务实例
func NewUserPointsService(userPointsRepo domain.UserPointsRepository, userRepo repository.UserRepository) UserPointsService {
	return &UserPointsServiceImpl{
		userPointsRepo: userPointsRepo,
		userRepo:       userRepo,
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
	id, err := s.userPointsRepo.Create(points)
	if err != nil {
		logx.Errorf("创建用户积分记录失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetPointByID 获取积分记录详情
func (s *UserPointsServiceImpl) GetPointByID(id int64) (*domain.UserPoints, error) {
	return s.userPointsRepo.GetByID(id)
}

// GetTotalPointsByUserID 获取用户总积分
func (s *UserPointsServiceImpl) GetTotalPointsByUserID(userID int64) (int, error) {
	return s.userPointsRepo.GetTotalPointsByUserID(userID)
}

// ListPointsByUserID 获取用户积分记录列表
func (s *UserPointsServiceImpl) ListPointsByUserID(userID int64, page, pageSize int64) ([]*domain.UserPoints, error) {
	offset := (page - 1) * pageSize
	return s.userPointsRepo.ListByUserID(userID, offset, pageSize)
}

// CountPointsByUserID 获取用户积分记录总数
func (s *UserPointsServiceImpl) CountPointsByUserID(userID int64) (int64, error) {
	return s.userPointsRepo.CountByUserID(userID)
}

// ListPointsWithTotal 获取积分记录列表及总数
func (s *UserPointsServiceImpl) ListPointsWithTotal(req *types.ListPointsRequest) ([]*domain.UserPoints, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	// 构造查询条件
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

	// 调用仓储层获取数据
	points, err := s.userPointsRepo.ListWithConditions(conditions, offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := s.userPointsRepo.CountWithConditions(conditions)
	if err != nil {
		return nil, 0, err
	}

	return points, total, nil
}

// DeletePoint 删除积分记录
func (s *UserPointsServiceImpl) DeletePoint(id int64) error {
	// 获取要删除的记录
	point, err := s.userPointsRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 只允许删除管理员创建的积分记录
	if point.Source != "admin" {
		return errors.New("只能删除管理员调整的积分记录")
	}

	// 创建抵消记录
	offsetPoint := &domain.UserPoints{
		UserID:      point.UserID,
		Points:      point.Points,
		Type:        (point.Type % 2) + 1, // 反转类型，如果原记录是增加积分，则抵消记录为减少积分
		Source:      "admin",
		Description: fmt.Sprintf("撤销积分调整 #%d: %s", point.ID, point.Description),
		RelatedID:   point.ID,
		RelatedType: "point_reversal",
		OperatorID:  point.OperatorID,
	}

	// 创建抵消记录
	_, err = s.userPointsRepo.Create(offsetPoint)
	if err != nil {
		return err
	}

	// 标记原记录为已撤销（可选，根据业务需求）
	return s.userPointsRepo.MarkAsRevoked(id)
}

// ExportPointsData 导出积分数据
func (s *UserPointsServiceImpl) ExportPointsData(req *types.ListPointsRequest) ([]byte, error) {
	// 不分页，获取所有符合条件的数据
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

	// 调用仓储层获取数据
	points, err := s.userPointsRepo.ListWithConditions(conditions, 0, 10000) // 限制最大导出数量
	if err != nil {
		return nil, err
	}

	// 实际实现中应该使用Excel处理库生成Excel文件
	// 这里简化处理，直接返回CSV格式的数据
	csvData := []byte("ID,用户ID,用户名,积分变动,变动后积分,类型,来源,描述,关联ID,关联类型,操作时间\n")

	for _, point := range points {
		// 格式化积分变动显示
		pointsStr := fmt.Sprintf("%d", point.Points)
		if point.Type == 1 {
			pointsStr = "+" + pointsStr
		} else {
			pointsStr = "-" + pointsStr
		}

		// 生成CSV行
		line := fmt.Sprintf("%d,%d,%s,%s,%d,%s,%s,%s,%d,%s,%s\n",
			point.ID,
			point.UserID,
			point.Username,
			pointsStr,
			point.TotalPoints,
			s.formatPointType(point.Type),
			s.formatPointSource(point.Source),
			point.Description,
			point.RelatedID,
			point.RelatedType,
			point.CreatedAt.Format("2006-01-02 15:04:05"))

		csvData = append(csvData, []byte(line)...)
	}

	return csvData, nil
}

// GetPointsStatistics 获取积分统计数据
func (s *UserPointsServiceImpl) GetPointsStatistics() (*types.PointsStatisticsResponse, error) {
	// 获取基础统计数据
	totalUsers, err := s.userPointsRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	totalPoints, err := s.userPointsRepo.SumPoints()
	if err != nil {
		return nil, err
	}

	avgPoints := int64(0)
	if totalUsers > 0 {
		avgPoints = totalPoints / totalUsers
	}

	maxPoints, err := s.userPointsRepo.MaxPoints()
	if err != nil {
		return nil, err
	}

	// 获取今日统计
	today := time.Now().Format("2006-01-02")
	todayIncrease, err := s.userPointsRepo.SumPointsByConditions(map[string]interface{}{
		"type":       1,
		"start_date": today,
		"end_date":   today,
	})
	if err != nil {
		return nil, err
	}

	todayDecrease, err := s.userPointsRepo.SumPointsByConditions(map[string]interface{}{
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

	monthIncrease, err := s.userPointsRepo.SumPointsByConditions(map[string]interface{}{
		"type":       1,
		"start_date": firstDay,
		"end_date":   lastDay,
	})
	if err != nil {
		return nil, err
	}

	monthDecrease, err := s.userPointsRepo.SumPointsByConditions(map[string]interface{}{
		"type":       2,
		"start_date": firstDay,
		"end_date":   lastDay,
	})
	if err != nil {
		return nil, err
	}

	// 获取来源分布
	sourceDistribution, err := s.userPointsRepo.GroupBySource()
	if err != nil {
		return nil, err
	}

	var sources []*types.SourceStats
	for _, src := range sourceDistribution {
		sources = append(sources, &types.SourceStats{
			Source: src.Source,
			Count:  src.Count,
		})
	}

	return &types.PointsStatisticsResponse{
		TotalUsers:         totalUsers,
		TotalPoints:        totalPoints,
		AvgPoints:          avgPoints,
		MaxPoints:          maxPoints,
		TodayIncrease:      todayIncrease,
		TodayDecrease:      todayDecrease,
		MonthIncrease:      monthIncrease,
		MonthDecrease:      monthDecrease,
		SourceDistribution: sources,
	}, nil
}

// GetPointsRules 获取积分规则
func (s *UserPointsServiceImpl) GetPointsRules() (*types.PointsRulesResponse, error) {
	rules, err := s.userPointsRepo.GetPointsRules()
	if err != nil {
		return nil, err
	}

	return &types.PointsRulesResponse{
		SignInPoints:      rules.SignInPoints,
		CommentPoints:     rules.CommentPoints,
		SharePoints:       rules.SharePoints,
		ArticlePoints:     rules.ArticlePoints,
		InvitePoints:      rules.InvitePoints,
		PurchaseRate:      rules.PurchaseRate,
		MaxDailyPoints:    rules.MaxDailyPoints,
		EnableExchange:    rules.EnableExchange,
		ExchangeRate:      rules.ExchangeRate,
		MinExchangePoints: rules.MinExchangePoints,
		UpdatedAt:         rules.UpdatedAt,
	}, nil
}

// UpdatePointsRules 更新积分规则
func (s *UserPointsServiceImpl) UpdatePointsRules(req *types.PointsRulesRequest) error {
	rules := &domain.PointsRules{
		SignInPoints:      req.SignInPoints,
		CommentPoints:     req.CommentPoints,
		SharePoints:       req.SharePoints,
		ArticlePoints:     req.ArticlePoints,
		InvitePoints:      req.InvitePoints,
		PurchaseRate:      req.PurchaseRate,
		MaxDailyPoints:    req.MaxDailyPoints,
		EnableExchange:    req.EnableExchange,
		ExchangeRate:      req.ExchangeRate,
		MinExchangePoints: req.MinExchangePoints,
		UpdatedAt:         time.Now(),
	}

	return s.userPointsRepo.UpdatePointsRules(rules)
}

// 验证积分规则
func (s *UserPointsServiceImpl) validatePoints(points *domain.UserPoints) error {
	if points.UserID <= 0 {
		return errors.New("用户ID无效")
	}

	if points.Points <= 0 {
		return errors.New("积分值必须大于0")
	}

	if points.Type != 1 && points.Type != 2 {
		return errors.New("积分类型无效，应为1(增加)或2(减少)")
	}

	if points.Type == 2 {
		// 积分减少时，检查用户积分是否足够
		totalPoints, err := s.userPointsRepo.GetTotalPointsByUserID(points.UserID)
		if err != nil {
			return err
		}

		if totalPoints < points.Points {
			return domain.ErrInsufficientPoints
		}
	}

	return nil
}

// 格式化积分类型
func (s *UserPointsServiceImpl) formatPointType(pointType int) string {
	switch pointType {
	case 1:
		return "增加"
	case 2:
		return "减少"
	default:
		return "未知"
	}
}

// 格式化积分来源
func (s *UserPointsServiceImpl) formatPointSource(source string) string {
	sourceMap := map[string]string{
		"sign":     "签到",
		"purchase": "购买",
		"activity": "活动",
		"comment":  "评论",
		"invite":   "邀请",
		"article":  "发布文章",
		"share":    "分享",
		"admin":    "管理员调整",
		"system":   "系统",
		"exchange": "积分兑换",
	}

	if name, ok := sourceMap[source]; ok {
		return name
	}

	return source
}
