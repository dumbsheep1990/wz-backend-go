package service

import (
	"time"
	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

// StatisticsServiceImpl 统计服务实现
type StatisticsServiceImpl struct {
	repo domain.StatisticsRepository
}

// NewStatisticsService 创建统计服务
func NewStatisticsService(repo domain.StatisticsRepository) StatisticsService {
	return &StatisticsServiceImpl{
		repo: repo,
	}
}

// 通过依赖注入SQL连接创建服务
func NewStatisticsServiceWithConn(conn interface{}) StatisticsService {
	return &StatisticsServiceImpl{
		repo: sql.NewStatisticsRepository(conn),
	}
}

// RecordStatistics 记录统计数据
func (s *StatisticsServiceImpl) RecordStatistics(data *domain.StatisticsData) error {
	logx.Infof("记录统计数据: %+v", data)

	// 业务规则验证
	if err := s.validateStatisticsData(data); err != nil {
		return err
	}

	// 如果没有设置日期，默认使用当前日期
	if data.Date.IsZero() {
		data.Date = time.Now().Truncate(24 * time.Hour) // 截断到当天0点
	}

	// 调用仓储层记录统计数据
	err := s.repo.RecordStatistics(data)
	if err != nil {
		logx.Errorf("记录统计数据失败: %v", err)
		return err
	}

	return nil
}

// GetOverview 获取站点概览数据
func (s *StatisticsServiceImpl) GetOverview(tenantID int64) (*domain.OverviewData, error) {
	logx.Infof("获取站点概览数据: tenantID=%d", tenantID)

	// 参数验证
	if tenantID <= 0 {
		logx.Error("无效的租户ID")
		return nil, domain.ErrInvalidParam
	}

	// 调用仓储层获取概览数据
	overview, err := s.repo.GetOverview(tenantID)
	if err != nil {
		logx.Errorf("获取站点概览数据失败: %v", err)
		return nil, err
	}

	return overview, nil
}

// GetStatisticsByType 获取特定类型的统计数据
func (s *StatisticsServiceImpl) GetStatisticsByType(tenantID int64, type_ string, startDate, endDate time.Time) ([]*domain.StatisticsData, error) {
	logx.Infof("获取特定类型的统计数据: tenantID=%d, type=%s, startDate=%v, endDate=%v",
		tenantID, type_, startDate, endDate)

	// 参数验证
	if tenantID <= 0 {
		logx.Error("无效的租户ID")
		return nil, domain.ErrInvalidParam
	}

	if type_ == "" {
		logx.Error("统计类型不能为空")
		return nil, domain.ErrInvalidParam
	}

	// 如果开始日期为空，默认为30天前
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)
	}

	// 如果结束日期为空，默认为今天
	if endDate.IsZero() {
		endDate = time.Now().Truncate(24 * time.Hour)
	}

	// 验证日期范围
	if startDate.After(endDate) {
		logx.Error("开始日期不能晚于结束日期")
		return nil, domain.ErrInvalidParam
	}

	// 调用仓储层获取统计数据
	stats, err := s.repo.GetStatisticsByType(tenantID, type_, startDate, endDate)
	if err != nil {
		logx.Errorf("获取统计数据失败: %v", err)
		return nil, err
	}

	return stats, nil
}

// GetContentRanking 获取内容排行榜
func (s *StatisticsServiceImpl) GetContentRanking(tenantID int64, limit int) ([]*domain.ContentRankingItem, error) {
	logx.Infof("获取内容排行榜: tenantID=%d, limit=%d", tenantID, limit)

	// 参数验证
	if tenantID <= 0 {
		logx.Error("无效的租户ID")
		return nil, domain.ErrInvalidParam
	}

	if limit <= 0 || limit > 100 {
		limit = 10 // 默认获取前10条
	}

	// 调用仓储层获取内容排行榜
	ranking, err := s.repo.GetContentRanking(tenantID, limit)
	if err != nil {
		logx.Errorf("获取内容排行榜失败: %v", err)
		return nil, err
	}

	return ranking, nil
}

// validateStatisticsData 验证统计数据
func (s *StatisticsServiceImpl) validateStatisticsData(data *domain.StatisticsData) error {
	// 参数验证
	if data.Type == "" {
		logx.Error("统计类型不能为空")
		return domain.ErrInvalidParam
	}

	if data.Value <= 0 {
		logx.Error("统计值必须大于0")
		return domain.ErrInvalidParam
	}

	if data.TenantID <= 0 {
		logx.Error("无效的租户ID")
		return domain.ErrInvalidParam
	}

	// 如果有关联项目，则项目ID和类型都必须有值
	if data.ItemID > 0 && data.ItemType == "" {
		logx.Error("关联项目类型不能为空")
		return domain.ErrInvalidParam
	}

	if data.ItemID <= 0 && data.ItemType != "" {
		logx.Error("关联项目ID无效")
		return domain.ErrInvalidParam
	}

	// 这里可以添加更多业务规则验证
	return nil
}
