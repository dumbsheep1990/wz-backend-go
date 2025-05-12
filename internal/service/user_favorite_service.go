package service

import (
	"fmt"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository"
	"wz-backend-go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserFavoriteServiceImpl 用户收藏服务实现
type UserFavoriteServiceImpl struct {
	favoriteRepo domain.UserFavoriteRepository
	userRepo     repository.UserRepository
}

// NewUserFavoriteService 创建用户收藏服务实例
func NewUserFavoriteService(favoriteRepo domain.UserFavoriteRepository, userRepo repository.UserRepository) UserFavoriteService {
	return &UserFavoriteServiceImpl{
		favoriteRepo: favoriteRepo,
		userRepo:     userRepo,
	}
}

// CreateFavorite 创建用户收藏
func (s *UserFavoriteServiceImpl) CreateFavorite(favorite *domain.UserFavorite) (int64, error) {
	logx.Infof("创建用户收藏: %+v", favorite)

	// 检查是否已收藏
	exists, err := s.favoriteRepo.CheckFavorite(favorite.UserID, favorite.ItemID, favorite.ItemType)
	if err != nil {
		logx.Errorf("检查收藏失败: %v", err)
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("已经收藏过该内容")
	}

	// 调用仓储层创建收藏记录
	id, err := s.favoriteRepo.Create(favorite)
	if err != nil {
		logx.Errorf("创建用户收藏失败: %v", err)
		return 0, err
	}

	return id, nil
}

// GetFavoriteByID 获取收藏详情
func (s *UserFavoriteServiceImpl) GetFavoriteByID(id int64) (*domain.UserFavorite, error) {
	return s.favoriteRepo.GetByID(id)
}

// ListFavoritesWithTotal 获取收藏列表及总数
func (s *UserFavoriteServiceImpl) ListFavoritesWithTotal(req *types.ListFavoritesRequest) ([]*domain.UserFavorite, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	// 构造查询条件
	conditions := make(map[string]interface{})
	if req.UserID > 0 {
		conditions["user_id"] = req.UserID
	}
	if req.Username != "" {
		conditions["username"] = req.Username
	}
	if req.Title != "" {
		conditions["title"] = req.Title
	}
	if req.ItemType != "" {
		conditions["item_type"] = req.ItemType
	}
	if req.StartDate != "" {
		conditions["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		conditions["end_date"] = req.EndDate
	}

	// 调用仓储层获取数据
	favorites, err := s.favoriteRepo.ListWithConditions(conditions, offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := s.favoriteRepo.CountWithConditions(conditions)
	if err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// DeleteFavorite 删除收藏记录
func (s *UserFavoriteServiceImpl) DeleteFavorite(id int64) error {
	// 获取要删除的记录
	_, err := s.favoriteRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 调用仓储层删除记录
	return s.favoriteRepo.DeleteByID(id)
}

// BatchDeleteFavorites 批量删除收藏记录
func (s *UserFavoriteServiceImpl) BatchDeleteFavorites(ids []int64) error {
	return s.favoriteRepo.BatchDelete(ids)
}

// CheckFavorite 检查是否已收藏
func (s *UserFavoriteServiceImpl) CheckFavorite(userID int64, itemID int64, itemType string) (bool, error) {
	return s.favoriteRepo.CheckFavorite(userID, itemID, itemType)
}

// ExportFavoritesData 导出收藏数据
func (s *UserFavoriteServiceImpl) ExportFavoritesData(req *types.ListFavoritesRequest) ([]byte, error) {
	// 不分页，获取所有符合条件的数据
	conditions := make(map[string]interface{})
	if req.UserID > 0 {
		conditions["user_id"] = req.UserID
	}
	if req.Username != "" {
		conditions["username"] = req.Username
	}
	if req.Title != "" {
		conditions["title"] = req.Title
	}
	if req.ItemType != "" {
		conditions["item_type"] = req.ItemType
	}
	if req.StartDate != "" {
		conditions["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		conditions["end_date"] = req.EndDate
	}

	// 调用仓储层获取数据
	favorites, err := s.favoriteRepo.ListWithConditions(conditions, 0, 10000) // 限制最大导出数量
	if err != nil {
		return nil, err
	}

	// 实际实现中应该使用Excel处理库生成Excel文件
	// 这里简化处理，直接返回CSV格式的数据
	csvData := []byte("ID,用户ID,用户名,内容ID,内容类型,标题,URL,收藏时间\n")

	for _, favorite := range favorites {
		// 生成CSV行
		line := fmt.Sprintf("%d,%d,%s,%d,%s,%s,%s,%s\n",
			favorite.ID,
			favorite.UserID,
			favorite.Username,
			favorite.ItemID,
			s.formatItemType(favorite.ItemType),
			favorite.Title,
			favorite.URL,
			favorite.CreatedAt.Format("2006-01-02 15:04:05"))

		csvData = append(csvData, []byte(line)...)
	}

	return csvData, nil
}

// GetFavoritesStatistics 获取收藏统计数据
func (s *UserFavoriteServiceImpl) GetFavoritesStatistics() (*types.FavoritesStatisticsResponse, error) {
	// 获取基础统计数据
	totalUsers, err := s.favoriteRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	totalFavorites, err := s.favoriteRepo.CountFavorites()
	if err != nil {
		return nil, err
	}

	todayFavorites, err := s.favoriteRepo.CountTodayFavorites()
	if err != nil {
		return nil, err
	}

	monthFavorites, err := s.favoriteRepo.CountMonthFavorites()
	if err != nil {
		return nil, err
	}

	// 获取类型分布
	typeDistribution, err := s.favoriteRepo.GroupByType()
	if err != nil {
		return nil, err
	}

	var typeStats []*types.TypeStats
	for _, item := range typeDistribution {
		typeStats = append(typeStats, &types.TypeStats{
			Type:  item.Type,
			Count: item.Count,
		})
	}

	return &types.FavoritesStatisticsResponse{
		TotalUsers:       totalUsers,
		TotalFavorites:   totalFavorites,
		TodayFavorites:   todayFavorites,
		MonthFavorites:   monthFavorites,
		TypeDistribution: typeStats,
	}, nil
}

// GetHotContent 获取热门收藏内容
func (s *UserFavoriteServiceImpl) GetHotContent() ([]*types.HotContentResponse, error) {
	// 获取热门内容数据
	hotContent, err := s.favoriteRepo.GetHotContent(10) // 获取前10个热门内容
	if err != nil {
		return nil, err
	}

	var result []*types.HotContentResponse
	for _, item := range hotContent {
		result = append(result, &types.HotContentResponse{
			ItemID:     item.ItemID,
			ItemType:   item.ItemType,
			Title:      item.Title,
			Cover:      item.Cover,
			Count:      item.Count,
			CreateDate: item.CreateDate,
		})
	}

	return result, nil
}

// GetFavoritesTrend 获取收藏趋势数据
func (s *UserFavoriteServiceImpl) GetFavoritesTrend(period string) ([]*types.TrendDataResponse, error) {
	// 验证周期参数
	if period != "week" && period != "month" && period != "year" {
		period = "month" // 默认为月
	}

	// 获取趋势数据
	trendData, err := s.favoriteRepo.GetTrend(period)
	if err != nil {
		return nil, err
	}

	var result []*types.TrendDataResponse
	for _, item := range trendData {
		result = append(result, &types.TrendDataResponse{
			Date:  item.Date,
			Count: item.Count,
		})
	}

	return result, nil
}

// 格式化收藏类型
func (s *UserFavoriteServiceImpl) formatItemType(itemType string) string {
	typeMap := map[string]string{
		"article": "文章",
		"product": "商品",
		"video":   "视频",
	}

	if name, ok := typeMap[itemType]; ok {
		return name
	}

	return itemType
}
