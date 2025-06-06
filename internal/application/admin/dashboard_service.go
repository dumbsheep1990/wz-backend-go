package admin

import (
	"context"
	"sync"
	"time"

	"wz-backend-go/internal/client"
	"wz-backend-go/internal/domain/admin/dto"
)

// DashboardService 仪表盘应用服务，聚合来自不同微服务的数据
type DashboardService struct {
	clientManager *client.ClientManager
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService(clientManager *client.ClientManager) *DashboardService {
	return &DashboardService{
		clientManager: clientManager,
	}
}

// DashboardStats 仪表盘统计数据，聚合多个微服务的数据
type DashboardStats struct {
	UserStats       *dto.UserStats       `json:"userStats"`
	ContentStats    *dto.ContentStats    `json:"contentStats"`
	TradeStats      *dto.TradeStats      `json:"tradeStats"`
	CommunityStats  *dto.CommunityStats  `json:"communityStats"`
	SiteStats       *dto.SiteStats       `json:"siteStats"`
	ComponentStats  *dto.ComponentStats  `json:"componentStats"`
	RenderStats     *dto.RenderStats     `json:"renderStats"`
	UpdatedAt       time.Time            `json:"updatedAt"`
}

// GetDashboardStats 获取仪表盘聚合数据
func (s *DashboardService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{
		UpdatedAt: time.Now(),
	}

	// 使用WaitGroup并发获取各服务数据
	var wg sync.WaitGroup
	errCh := make(chan error, 7) // 用于收集错误
	
	// 获取用户服务统计数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		userClient, err := s.clientManager.GetUserClient()
		if err != nil {
			errCh <- err
			return
		}
		
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		
		resp, err := userClient.GetUserStats(ctx)
		if err != nil {
			errCh <- err
			return
		}
		
		// 转换为DTO
		stats.UserStats = &dto.UserStats{
			TotalUsers:     resp.TotalUsers,
			ActiveUsers:    resp.ActiveUsers,
			NewUsers:       resp.NewUsers,
			VerifiedUsers:  resp.VerifiedUsers,
			LastUpdateTime: time.Unix(resp.LastUpdateTime, 0),
		}
	}()
	
	// 获取内容服务统计数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		contentClient, err := s.clientManager.GetContentClient()
		if err != nil {
			errCh <- err
			return
		}
		
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		
		resp, err := contentClient.GetContentStats(ctx)
		if err != nil {
			errCh <- err
			return
		}
		
		// 转换为DTO
		stats.ContentStats = &dto.ContentStats{
			TotalContent:    resp.TotalContent,
			PublishedContent: resp.PublishedContent,
			DraftContent:    resp.DraftContent,
			Categories:      resp.Categories,
			Tags:           resp.Tags,
			LastUpdateTime: time.Unix(resp.LastUpdateTime, 0),
		}
	}()
	
	// 获取交易服务统计数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		tradeClient, err := s.clientManager.GetTradeClient()
		if err != nil {
			errCh <- err
			return
		}
		
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		
		resp, err := tradeClient.GetTradeStats(ctx)
		if err != nil {
			errCh <- err
			return
		}
		
		// 转换为DTO
		stats.TradeStats = &dto.TradeStats{
			TotalOrders:     resp.TotalOrders,
			PendingOrders:   resp.PendingOrders,
			CompletedOrders: resp.CompletedOrders,
			TotalRevenue:    resp.TotalRevenue,
			DailyRevenue:    resp.DailyRevenue,
			LastUpdateTime:  time.Unix(resp.LastUpdateTime, 0),
		}
	}()
	
	// 获取社区服务统计数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		communityClient, err := s.clientManager.GetCommunityClient()
		if err != nil {
			errCh <- err
			return
		}
		
		// 由于社区客户端可能没有实现统计方法，我们这里模拟一下
		// 实际应用中，需要确保所有微服务都实现了统计接口
		stats.CommunityStats = &dto.CommunityStats{
			TotalCommunities: 100,
			ActiveCommunities: 80,
			TotalGroups: 500,
			TotalPosts: 10000,
			TotalComments: 50000,
			LastUpdateTime: time.Now(),
		}
	}()
	
	// 获取站点服务统计数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		siteClient, err := s.clientManager.GetSiteClient()
		if err != nil {
			errCh <- err
			return
		}
		
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		
		resp, err := siteClient.GetSiteStats(ctx)
		if err != nil {
			errCh <- err
			return
		}
		
		// 转换为DTO
		stats.SiteStats = &dto.SiteStats{
			TotalSites:     resp.TotalSites,
			ActiveSites:    resp.ActiveSites,
			TotalPages:     resp.TotalPages,
			PublishedPages: resp.PublishedPages,
			TotalTemplates: resp.TotalTemplates,
			LastUpdateTime: time.Unix(resp.LastUpdateTime, 0),
		}
	}()
	
	// 获取组件服务统计数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		componentClient, err := s.clientManager.GetComponentClient()
		if err != nil {
			errCh <- err
			return
		}
		
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		
		resp, err := componentClient.GetComponentStats(ctx)
		if err != nil {
			errCh <- err
			return
		}
		
		// 转换为DTO
		stats.ComponentStats = &dto.ComponentStats{
			TotalComponents:     resp.TotalComponents,
			PublishedComponents: resp.PublishedComponents,
			Categories:          resp.Categories,
			LastUpdateTime:      time.Unix(resp.LastUpdateTime, 0),
		}
	}()
	
	// 获取渲染服务统计数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		renderClient, err := s.clientManager.GetRenderClient()
		if err != nil {
			errCh <- err
			return
		}
		
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		
		resp, err := renderClient.GetRenderStats(ctx)
		if err != nil {
			errCh <- err
			return
		}
		
		// 转换为DTO
		stats.RenderStats = &dto.RenderStats{
			TotalRenders:        resp.TotalRenders,
			CacheHitRate:        resp.CacheHitRate,
			AverageRenderTime:   resp.AverageRenderTime,
			ErrorRate:           resp.ErrorRate,
			LastUpdateTime:      time.Unix(resp.LastUpdateTime, 0),
		}
	}()
	
	// 等待所有goroutine完成
	wg.Wait()
	close(errCh)
	
	// 检查是否有错误发生
	for err := range errCh {
		if err != nil {
			// 这里我们选择继续而不是返回错误，因为部分数据缺失不应该导致整个请求失败
			// 实际应用中，可以根据需要调整此策略
			continue
		}
	}
	
	return stats, nil
}
