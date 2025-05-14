package service

import (
	"errors"
	"fmt"
	"time"
	"wz-backend-go/models"
)

// 模拟数据 - 区块
var sections = map[string][]models.Section{
	"page1": {
		{
			ID:         "section1",
			PageID:     "page1",
			Type:       "header",
			Title:      "页面头部",
			Settings:   map[string]interface{}{"backgroundColor": "#f5f5f5"},
			Components: []models.Component{},
			Style:      map[string]interface{}{"padding": "20px"},
			SortOrder:  0,
		},
		{
			ID:         "section2",
			PageID:     "page1",
			Type:       "content",
			Title:      "主要内容",
			Settings:   map[string]interface{}{"columns": 1},
			Components: []models.Component{},
			Style:      map[string]interface{}{"margin": "20px 0"},
			SortOrder:  1,
		},
	},
}

// ListSections 获取页面下的所有区块
func ListSections(siteID string, pageID string) ([]models.Section, error) {
	// 先检查页面是否存在
	if _, err := GetPage(siteID, pageID); err != nil {
		return nil, err
	}

	if pageSections, exists := sections[pageID]; exists {
		return pageSections, nil
	}

	return []models.Section{}, nil
}

// AddSection 添加新区块
func AddSection(siteID string, pageID string, section models.Section) (models.Section, error) {
	// 先检查页面是否存在
	if _, err := GetPage(siteID, pageID); err != nil {
		return models.Section{}, err
	}

	// 生成一个简单的ID
	section.ID = fmt.Sprintf("section%d", len(sections[pageID])+1)

	// 设置排序顺序
	if pageSections, exists := sections[pageID]; exists {
		section.SortOrder = len(pageSections)
	} else {
		section.SortOrder = 0
		sections[pageID] = []models.Section{}
	}

	// 添加到区块集合
	sections[pageID] = append(sections[pageID], section)

	// 更新页面的时间戳
	UpdatePageTimestamp(siteID, pageID)

	return section, nil
}

// UpdateSection 更新区块
func UpdateSection(siteID string, pageID string, section models.Section) (models.Section, error) {
	// 先检查页面是否存在
	if _, err := GetPage(siteID, pageID); err != nil {
		return models.Section{}, err
	}

	if pageSections, exists := sections[pageID]; exists {
		for i, s := range pageSections {
			if s.ID == section.ID {
				// 保持一些不应该被客户端更新的字段
				section.SortOrder = s.SortOrder

				pageSections[i] = section
				sections[pageID] = pageSections

				return section, nil
			}
		}
		return models.Section{}, errors.New("区块不存在")
	}

	return models.Section{}, errors.New("页面没有区块")
}

// DeleteSection 删除区块
func DeleteSection(siteID string, pageID string, sectionID string) error {
	// 先检查页面是否存在
	if _, err := GetPage(siteID, pageID); err != nil {
		return err
	}

	if pageSections, exists := sections[pageID]; exists {
		for i, section := range pageSections {
			if section.ID == sectionID {
				// 从切片中删除
				sections[pageID] = append(pageSections[:i], pageSections[i+1:]...)

				// 重新调整排序
				for j := i; j < len(sections[pageID]); j++ {
					sections[pageID][j].SortOrder = j
				}

				return nil
			}
		}
		return errors.New("区块不存在")
	}

	return errors.New("页面没有区块")
}

// ReorderSections 重新排序区块
func ReorderSections(siteID string, pageID string, sectionOrder []string) error {
	// 先检查页面是否存在
	if _, err := GetPage(siteID, pageID); err != nil {
		return err
	}

	if pageSections, exists := sections[pageID]; exists {
		if len(sectionOrder) != len(pageSections) {
			return errors.New("区块数量不匹配")
		}

		// 创建新的排序区块切片
		newOrder := make([]models.Section, len(pageSections))
		for i, sectionID := range sectionOrder {
			found := false
			for _, section := range pageSections {
				if section.ID == sectionID {
					section.SortOrder = i
					newOrder[i] = section
					found = true
					break
				}
			}
			if !found {
				return errors.New("无效的区块ID")
			}
		}

		sections[pageID] = newOrder
		return nil
	}

	return errors.New("页面没有区块")
}

// UpdatePageTimestamp 更新页面时间戳
func UpdatePageTimestamp(siteID string, pageID string) {
	if sitePages, exists := pages[siteID]; exists {
		for i, page := range sitePages {
			if page.ID == pageID {
				sitePages[i].UpdatedAt = time.Now()
				pages[siteID] = sitePages
				break
			}
		}
	}
}
