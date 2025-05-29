package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"wz-backend-go/internal/application/user/dto"
	"wz-backend-go/internal/domain/shared/event"
	domainEntity "wz-backend-go/internal/domain/user/entity"
	favoriteEvent "wz-backend-go/internal/domain/user/event"
	"wz-backend-go/internal/domain/user/repository"
	"wz-backend-go/internal/domain/user/valueobject"
	"wz-backend-go/internal/infrastructure/persistence/database"
)

// UserFavoriteApplicationService 用户收藏应用服务
type UserFavoriteApplicationService struct {
	favoriteRepo repository.UserFavoriteRepository
	eventBus     event.EventBus
	uow          database.UnitOfWork
}

// NewUserFavoriteApplicationService 创建用户收藏应用服务
func NewUserFavoriteApplicationService(
	favoriteRepo repository.UserFavoriteRepository,
	eventBus event.EventBus,
	uow database.UnitOfWork,
) *UserFavoriteApplicationService {
	return &UserFavoriteApplicationService{
		favoriteRepo: favoriteRepo,
		eventBus:     eventBus,
		uow:          uow,
	}
}

// CreateFavorite 创建收藏记录
func (s *UserFavoriteApplicationService) CreateFavorite(ctx context.Context, req *dto.CreateFavoriteRequest) (*dto.FavoriteDTO, error) {
	// 验证必填参数
	if req.UserID <= 0 || req.ItemID <= 0 || req.Title == "" || req.ItemType == "" {
		return nil, errors.New("缺少必要参数")
	}

	// 创建值对象
	userID, err := valueobject.NewFavoriteUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	itemID, err := valueobject.NewFavoriteItemID(req.ItemID)
	if err != nil {
		return nil, err
	}

	itemType, err := valueobject.NewFavoriteItemType(req.ItemType)
	if err != nil {
		return nil, err
	}

	title, err := valueobject.NewFavoriteTitle(req.Title)
	if err != nil {
		return nil, err
	}

	url, err := valueobject.NewFavoriteURL(req.URL)
	if err != nil {
		return nil, err
	}

	tenantID, err := valueobject.NewFavoriteTenantID(req.TenantID)
	if err != nil {
		return nil, err
	}

	// 检查是否已收藏
	exists, err := s.favoriteRepo.CheckFavorite(ctx, userID, itemID, itemType)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("已经收藏过该内容")
	}

	// 使用工作单元开始事务
	var favorite *domainEntity.UserFavorite
	var favoriteID valueobject.FavoriteID

	err = s.uow.DoInTx(ctx, func(ctx context.Context) error {
		// 创建实体
		tempID, _ := valueobject.NewFavoriteID(0) // 临时ID，数据库会生成真正的ID
		favorite, err = domainEntity.NewUserFavorite(
			tempID,
			userID,
			itemID,
			itemType,
			title,
			req.Cover,
			req.Summary,
			url,
			req.Remark,
			tenantID,
		)
		if err != nil {
			return err
		}

		// 保存到数据库
		favoriteID, err = s.favoriteRepo.Create(ctx, favorite)
		if err != nil {
			return err
		}

		// 发布领域事件
		newFavorite, _ := s.favoriteRepo.GetByID(ctx, favoriteID)
		createdEvent := favoriteEvent.NewUserFavoriteCreatedEvent(newFavorite)
		return s.eventBus.Publish(ctx, createdEvent)
	})

	if err != nil {
		logx.Errorf("创建收藏失败: %v", err)
		return nil, err
	}

	// 获取创建后的记录
	createdFavorite, err := s.favoriteRepo.GetByID(ctx, favoriteID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return &dto.FavoriteDTO{
		ID:        createdFavorite.ID().Value(),
		UserID:    createdFavorite.UserID().Value(),
		Username:  createdFavorite.Username(),
		ItemID:    createdFavorite.ItemID().Value(),
		ItemType:  createdFavorite.ItemType().Value(),
		Title:     createdFavorite.Title().Value(),
		Cover:     createdFavorite.Cover(),
		Summary:   createdFavorite.Summary(),
		URL:       createdFavorite.URL().Value(),
		Remark:    createdFavorite.Remark(),
		TenantID:  createdFavorite.TenantID().Value(),
		CreatedAt: createdFavorite.CreatedAt(),
		UpdatedAt: createdFavorite.UpdatedAt(),
	}, nil
}

// GetFavoriteByID 获取收藏详情
func (s *UserFavoriteApplicationService) GetFavoriteByID(ctx context.Context, id int64) (*dto.FavoriteDTO, error) {
	favoriteID, err := valueobject.NewFavoriteID(id)
	if err != nil {
		return nil, err
	}

	favorite, err := s.favoriteRepo.GetByID(ctx, favoriteID)
	if err != nil {
		return nil, err
	}

	return &dto.FavoriteDTO{
		ID:        favorite.ID().Value(),
		UserID:    favorite.UserID().Value(),
		Username:  favorite.Username(),
		ItemID:    favorite.ItemID().Value(),
		ItemType:  favorite.ItemType().Value(),
		Title:     favorite.Title().Value(),
		Cover:     favorite.Cover(),
		Summary:   favorite.Summary(),
		URL:       favorite.URL().Value(),
		Remark:    favorite.Remark(),
		TenantID:  favorite.TenantID().Value(),
		CreatedAt: favorite.CreatedAt(),
		UpdatedAt: favorite.UpdatedAt(),
	}, nil
}

// ListFavorites 获取收藏列表
func (s *UserFavoriteApplicationService) ListFavorites(ctx context.Context, req *dto.ListFavoritesRequest) (*dto.ListFavoritesResponse, error) {
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
	if req.TenantID > 0 {
		conditions["tenant_id"] = req.TenantID
	}

	// 调用仓储层获取数据
	favorites, err := s.favoriteRepo.ListWithConditions(ctx, conditions, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := s.favoriteRepo.CountWithConditions(ctx, conditions)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	var favoriteDTOs []*dto.FavoriteDTO
	for _, favorite := range favorites {
		favoriteDTOs = append(favoriteDTOs, &dto.FavoriteDTO{
			ID:        favorite.ID().Value(),
			UserID:    favorite.UserID().Value(),
			Username:  favorite.Username(),
			ItemID:    favorite.ItemID().Value(),
			ItemType:  favorite.ItemType().Value(),
			Title:     favorite.Title().Value(),
			Cover:     favorite.Cover(),
			Summary:   favorite.Summary(),
			URL:       favorite.URL().Value(),
			Remark:    favorite.Remark(),
			TenantID:  favorite.TenantID().Value(),
			CreatedAt: favorite.CreatedAt(),
			UpdatedAt: favorite.UpdatedAt(),
		})
	}

	return &dto.ListFavoritesResponse{
		List:  favoriteDTOs,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// DeleteFavorite 删除收藏记录
func (s *UserFavoriteApplicationService) DeleteFavorite(ctx context.Context, id int64, currentUserID int64, isAdmin bool) error {
	favoriteID, err := valueobject.NewFavoriteID(id)
	if err != nil {
		return err
	}

	// 获取要删除的记录
	favorite, err := s.favoriteRepo.GetByID(ctx, favoriteID)
	if err != nil {
		return err
	}

	// 权限检查：只有管理员或记录创建者可以删除
	currentUserIDVO, _ := valueobject.NewFavoriteUserID(currentUserID)
	if !isAdmin && !favorite.IsBelongsToUser(currentUserIDVO) {
		return errors.New("无权限删除此收藏")
	}

	// 使用工作单元开始事务
	return s.uow.DoInTx(ctx, func(ctx context.Context) error {
		// 删除记录
		err = s.favoriteRepo.DeleteByID(ctx, favoriteID)
		if err != nil {
			return err
		}

		// 发布领域事件
		deletedEvent := favoriteEvent.NewUserFavoriteDeletedEvent(favorite)
		return s.eventBus.Publish(ctx, deletedEvent)
	})
}

// BatchDeleteFavorites 批量删除收藏记录
func (s *UserFavoriteApplicationService) BatchDeleteFavorites(ctx context.Context, req *dto.BatchDeleteFavoritesRequest, currentUserID int64, isAdmin bool) error {
	if len(req.IDs) == 0 {
		return errors.New("未指定要删除的记录ID")
	}

	// 转换ID列表
	var favoriteIDs []valueobject.FavoriteID
	for _, id := range req.IDs {
		favoriteID, err := valueobject.NewFavoriteID(id)
		if err != nil {
			return err
		}
		favoriteIDs = append(favoriteIDs, favoriteID)
	}

	// 对于非管理员，需要验证权限
	if !isAdmin {
		currentUserIDVO, _ := valueobject.NewFavoriteUserID(currentUserID)
		for _, id := range favoriteIDs {
			favorite, err := s.favoriteRepo.GetByID(ctx, id)
			if err != nil {
				return err
			}
			if !favorite.IsBelongsToUser(currentUserIDVO) {
				return fmt.Errorf("无权限删除ID为 %d 的收藏", id.Value())
			}
		}
	}

	// 使用工作单元开始事务
	return s.uow.DoInTx(ctx, func(ctx context.Context) error {
		// 批量删除记录
		return s.favoriteRepo.BatchDelete(ctx, favoriteIDs)
	})
}

// CheckFavorite 检查是否已收藏
func (s *UserFavoriteApplicationService) CheckFavorite(ctx context.Context, req *dto.CheckFavoriteRequest) (*dto.CheckFavoriteResponse, error) {
	userID, err := valueobject.NewFavoriteUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	itemID, err := valueobject.NewFavoriteItemID(req.ItemID)
	if err != nil {
		return nil, err
	}

	itemType, err := valueobject.NewFavoriteItemType(req.ItemType)
	if err != nil {
		return nil, err
	}

	exists, err := s.favoriteRepo.CheckFavorite(ctx, userID, itemID, itemType)
	if err != nil {
		return nil, err
	}

	return &dto.CheckFavoriteResponse{
		Exists: exists,
	}, nil
}

// GetFavoritesStatistics 获取收藏统计数据
func (s *UserFavoriteApplicationService) GetFavoritesStatistics(ctx context.Context) (*dto.FavoritesStatisticsResponse, error) {
	// 获取基础统计数据
	totalUsers, err := s.favoriteRepo.CountUsers(ctx)
	if err != nil {
		return nil, err
	}

	totalFavorites, err := s.favoriteRepo.CountFavorites(ctx)
	if err != nil {
		return nil, err
	}

	todayFavorites, err := s.favoriteRepo.CountTodayFavorites(ctx)
	if err != nil {
		return nil, err
	}

	monthFavorites, err := s.favoriteRepo.CountMonthFavorites(ctx)
	if err != nil {
		return nil, err
	}

	// 获取类型分布
	typeDistribution, err := s.favoriteRepo.GroupByType(ctx)
	if err != nil {
		return nil, err
	}

	var typeStats []*dto.TypeStatsDTO
	for _, item := range typeDistribution {
		typeStats = append(typeStats, &dto.TypeStatsDTO{
			Type:  item.Type,
			Count: item.Count,
		})
	}

	return &dto.FavoritesStatisticsResponse{
		TotalUsers:       totalUsers,
		TotalFavorites:   totalFavorites,
		TodayFavorites:   todayFavorites,
		MonthFavorites:   monthFavorites,
		TypeDistribution: typeStats,
	}, nil
}

// GetHotContent 获取热门收藏内容
func (s *UserFavoriteApplicationService) GetHotContent(ctx context.Context, limit int) ([]*dto.HotContentResponseDTO, error) {
	if limit <= 0 {
		limit = 10 // 默认获取10条
	}

	// 获取热门内容数据
	hotContent, err := s.favoriteRepo.GetHotContent(ctx, limit)
	if err != nil {
		return nil, err
	}

	var result []*dto.HotContentResponseDTO
	for _, item := range hotContent {
		result = append(result, &dto.HotContentResponseDTO{
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
func (s *UserFavoriteApplicationService) GetFavoritesTrend(ctx context.Context, req *dto.GetFavoritesTrendRequest) ([]*dto.TrendDataResponseDTO, error) {
	// 验证周期参数
	if req.Period != "week" && req.Period != "month" && req.Period != "year" {
		req.Period = "month" // 默认为月
	}

	// 获取趋势数据
	trendData, err := s.favoriteRepo.GetTrend(ctx, req.Period)
	if err != nil {
		return nil, err
	}

	var result []*dto.TrendDataResponseDTO
	for _, item := range trendData {
		result = append(result, &dto.TrendDataResponseDTO{
			Date:  item.Date,
			Count: item.Count,
		})
	}

	return result, nil
}

// ExportFavoritesData 导出收藏数据
func (s *UserFavoriteApplicationService) ExportFavoritesData(ctx context.Context, req *dto.ExportFavoritesDataRequest) ([]byte, error) {
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
	if req.TenantID > 0 {
		conditions["tenant_id"] = req.TenantID
	}

	// 调用仓储层获取数据
	favorites, err := s.favoriteRepo.ListWithConditions(ctx, conditions, 0, 10000) // 限制最大导出数量
	if err != nil {
		return nil, err
	}

	// 实际实现中应该使用Excel处理库生成Excel文件
	// 这里简化处理，直接返回CSV格式的数据
	csvData := []byte("ID,用户ID,用户名,内容ID,内容类型,标题,URL,收藏时间\n")

	for _, favorite := range favorites {
		// 生成CSV行
		line := fmt.Sprintf("%d,%d,%s,%d,%s,%s,%s,%s\n",
			favorite.ID().Value(),
			favorite.UserID().Value(),
			favorite.Username(),
			favorite.ItemID().Value(),
			favorite.ItemType().Value(),
			favorite.Title().Value(),
			favorite.URL().Value(),
			favorite.CreatedAt().Format("2006-01-02 15:04:05"))

		csvData = append(csvData, []byte(line)...)
	}

	return csvData, nil
}
