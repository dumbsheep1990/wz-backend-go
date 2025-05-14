package service

import (
	"errors"
	"fmt"
	"time"
	"wz-backend-go/models"
)

// 模拟数据
var pages = map[string][]models.Page{
	"1": {
		{
			ID:          "page1",
			SiteID:      "1",
			Name:        "首页",
			Slug:        "home",
			Title:       "企业官网首页",
			Description: "欢迎访问我们的企业官网",
			Keywords:    []string{"企业", "官网", "首页"},
			IsHomepage:  true,
			Layout:      "default",
			Sections:    []models.Section{},
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
			SortOrder:   0,
		},
		{
			ID:          "page2",
			SiteID:      "1",
			Name:        "关于我们",
			Slug:        "about",
			Title:       "关于我们 - 企业官网",
			Description: "了解我们的企业文化和团队",
			Keywords:    []string{"关于", "企业文化", "团队"},
			IsHomepage:  false,
			Layout:      "default",
			Sections:    []models.Section{},
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
			SortOrder:   1,
		},
	},
}

// 模拟站点权限检查
var siteTenants = map[string]string{
	"1": "tenant_456",
}

// 检查站点访问权限
func CheckSiteAccess(siteID string, tenantID string) bool {
	if tenant, exists := siteTenants[siteID]; exists {
		return tenant == tenantID
	}
	return false
}

// ListPages 获取站点的所有页面
func ListPages(siteID string) ([]models.Page, error) {
	if sitePages, exists := pages[siteID]; exists {
		return sitePages, nil
	}
	return []models.Page{}, nil
}

// GetPage 获取单个页面
func GetPage(siteID string, pageID string) (models.Page, error) {
	if sitePages, exists := pages[siteID]; exists {
		for _, page := range sitePages {
			if page.ID == pageID {
				return page, nil
			}
		}
	}
	return models.Page{}, errors.New("页面不存在")
}

// CreatePage 创建新页面
func CreatePage(page models.Page) (models.Page, error) {
	// 生成一个简单的ID
	page.ID = fmt.Sprintf("page%d", len(pages[page.SiteID])+1)

	// 设置排序顺序
	if sitePages, exists := pages[page.SiteID]; exists {
		page.SortOrder = len(sitePages)
	} else {
		page.SortOrder = 0
		pages[page.SiteID] = []models.Page{}
	}

	// 添加到页面集合
	pages[page.SiteID] = append(pages[page.SiteID], page)

	return page, nil
}

// UpdatePage 更新页面
func UpdatePage(page models.Page) (models.Page, error) {
	if sitePages, exists := pages[page.SiteID]; exists {
		for i, p := range sitePages {
			if p.ID == page.ID {
				// 保留一些不应该被客户端更新的字段
				page.CreatedAt = p.CreatedAt
				page.SortOrder = p.SortOrder

				sitePages[i] = page
				pages[page.SiteID] = sitePages
				return page, nil
			}
		}
	}
	return models.Page{}, errors.New("页面不存在")
}

// DeletePage 删除页面
func DeletePage(siteID string, pageID string) error {
	if sitePages, exists := pages[siteID]; exists {
		for i, page := range sitePages {
			if page.ID == pageID {
				// 从切片中删除
				pages[siteID] = append(sitePages[:i], sitePages[i+1:]...)

				// 重新调整排序
				for j := i; j < len(pages[siteID]); j++ {
					pages[siteID][j].SortOrder = j
				}

				return nil
			}
		}
	}
	return errors.New("页面不存在")
}

// UnsetOtherHomepages 将其他页面设置为非首页
func UnsetOtherHomepages(siteID string, exceptPageID ...string) error {
	if sitePages, exists := pages[siteID]; exists {
		for i, page := range sitePages {
			// 如果页面ID不在例外列表中，则设置为非首页
			isExcepted := false
			for _, exceptID := range exceptPageID {
				if page.ID == exceptID {
					isExcepted = true
					break
				}
			}

			if !isExcepted && page.IsHomepage {
				sitePages[i].IsHomepage = false
				sitePages[i].UpdatedAt = time.Now()
			}
		}
		pages[siteID] = sitePages
	}
	return nil
}

// SetHomepage 设置页面为首页
func SetHomepage(siteID string, pageID string) (models.Page, error) {
	if sitePages, exists := pages[siteID]; exists {
		for i, page := range sitePages {
			if page.ID == pageID {
				sitePages[i].IsHomepage = true
				sitePages[i].UpdatedAt = time.Now()
				pages[siteID] = sitePages
				return sitePages[i], nil
			}
		}
	}
	return models.Page{}, errors.New("页面不存在")
}

// ReorderPages 重新排序页面
func ReorderPages(siteID string, pageOrder []string) error {
	if sitePages, exists := pages[siteID]; exists {
		if len(pageOrder) != len(sitePages) {
			return errors.New("页面数量不匹配")
		}

		// 创建新的排序页面切片
		newOrder := make([]models.Page, len(sitePages))
		for i, pageID := range pageOrder {
			found := false
			for _, page := range sitePages {
				if page.ID == pageID {
					page.SortOrder = i
					page.UpdatedAt = time.Now()
					newOrder[i] = page
					found = true
					break
				}
			}
			if !found {
				return errors.New("无效的页面ID")
			}
		}

		pages[siteID] = newOrder
		return nil
	}
	return nil
}
