package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"wz-backend-go/models/common"
	"wz-backend-go/models/site"
)

// 模拟数据库
var sites = []site.Site{
	{
		ID:          "1",
		Name:        "企业展示站点",
		Description: "适合企业官网使用的模板",
		Domain:      "company.wanzhimarket.com",
		Logo:        "/img/logo1.png",
		Favicon:     "/img/favicon1.ico",
		TenantID:    "tenant_456",
		Theme: site.ThemeConfig{
			PrimaryColor:    "#FF5722",
			SecondaryColor:  "#2196F3",
			AccentColor:     "#4CAF50",
			TextColor:       "#333333",
			BackgroundColor: "#FFFFFF",
			FontFamily:      "Arial, sans-serif",
			HeaderStyle:     "standard",
			BorderRadius:    "medium",
			CustomCSS:       "",
		},
		Navigation: site.Navigation{
			Type:  "horizontal",
			Items: []site.NavigationItem{},
			Style: map[string]string{},
		},
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-12 * time.Hour),
		PublishedAt: nil,
		Status:      "draft",
		Thumbnail:   "/img/site-thumbnail1.jpg",
	},
}

// ListSites 获取站点列表
func ListSites(tenantID string, status string, search string) ([]site.Site, error) {
	var result []site.Site

	for _, site := range sites {
		// 过滤租户ID
		if site.TenantID != tenantID {
			continue
		}

		// 状态过滤
		if status != "" && site.Status != status {
			continue
		}

		// 搜索过滤
		if search != "" {
			searchLower := strings.ToLower(search)
			nameLower := strings.ToLower(site.Name)
			descLower := strings.ToLower(site.Description)

			if !strings.Contains(nameLower, searchLower) && !strings.Contains(descLower, searchLower) {
				continue
			}
		}

		result = append(result, site)
	}

	return result, nil
}

// GetSite 获取单个站点
func GetSite(siteID string, tenantID string) (site.Site, error) {
	for _, site := range sites {
		if site.ID == siteID && site.TenantID == tenantID {
			return site, nil
		}
	}

	return site.Site{}, errors.New("站点不存在")
}

// CreateSite 创建新站点
func CreateSite(site site.Site) (site.Site, error) {
	// 生成一个简单的ID (在实际应用中应该使用UUID或其他唯一标识符)
	site.ID = fmt.Sprintf("%d", len(sites)+1)

	// 添加到集合
	sites = append(sites, site)

	return site, nil
}

// UpdateSite 更新站点
func UpdateSite(site site.Site) (site.Site, error) {
	for i, s := range sites {
		if s.ID == site.ID {
			// 保留一些不应该被客户端更新的字段
			site.TenantID = s.TenantID
			site.CreatedAt = s.CreatedAt
			site.PublishedAt = s.PublishedAt
			site.Status = s.Status

			sites[i] = site
			return site, nil
		}
	}

	return site.Site{}, errors.New("站点不存在")
}

// DeleteSite 删除站点
func DeleteSite(siteID string) error {
	for i, site := range sites {
		if site.ID == siteID {
			// 从切片中删除
			sites = append(sites[:i], sites[i+1:]...)
			return nil
		}
	}

	return errors.New("站点不存在")
}

// PublishSite 发布站点
func PublishSite(siteID string) (site.Site, error) {
	for i, site := range sites {
		if site.ID == siteID {
			now := time.Now()
			sites[i].Status = "published"
			sites[i].PublishedAt = &now
			sites[i].UpdatedAt = now

			return sites[i], nil
		}
	}

	return site.Site{}, errors.New("站点不存在")
}

// CheckSiteOwnership 检查站点所有权
func CheckSiteOwnership(siteID string, tenantID string) (bool, error) {
	for _, site := range sites {
		if site.ID == siteID {
			return site.TenantID == tenantID, nil
		}
	}

	return false, errors.New("站点不存在")
}
