package service

import (
	"errors"
	"fmt"
	"wz-backend-go/models"
)

// 模拟组件分类数据
var componentCategories = []models.ComponentCategory{
	{
		ID:   "basic",
		Name: "基础组件",
		Components: []models.ComponentDefinition{
			{
				Type:        "text",
				Name:        "文本",
				Icon:        "text",
				Description: "文本内容块",
				DefaultSettings: map[string]interface{}{
					"textAlign": "left",
					"fontSize":  "16px",
				},
			},
			{
				Type:        "heading",
				Name:        "标题",
				Icon:        "heading",
				Description: "标题文本",
				DefaultSettings: map[string]interface{}{
					"level":     "h2",
					"textAlign": "left",
				},
			},
			{
				Type:        "button",
				Name:        "按钮",
				Icon:        "button",
				Description: "可点击的按钮",
				DefaultSettings: map[string]interface{}{
					"style":   "filled",
					"size":    "medium",
					"rounded": true,
				},
			},
			{
				Type:        "divider",
				Name:        "分隔线",
				Icon:        "divider",
				Description: "水平分隔线",
				DefaultSettings: map[string]interface{}{
					"style": "solid",
					"width": "100%",
				},
			},
		},
	},
	{
		ID:   "layout",
		Name: "布局组件",
		Components: []models.ComponentDefinition{
			{
				Type:        "container",
				Name:        "容器",
				Icon:        "container",
				Description: "内容容器",
				DefaultSettings: map[string]interface{}{
					"width":  "100%",
					"height": "auto",
				},
			},
			{
				Type:        "row",
				Name:        "行",
				Icon:        "row",
				Description: "水平行",
				DefaultSettings: map[string]interface{}{
					"gutter": 16,
				},
			},
			{
				Type:        "column",
				Name:        "列",
				Icon:        "column",
				Description: "垂直列",
				DefaultSettings: map[string]interface{}{
					"span": 12,
				},
			},
			{
				Type:        "card",
				Name:        "卡片",
				Icon:        "card",
				Description: "卡片容器",
				DefaultSettings: map[string]interface{}{
					"shadow":  "medium",
					"padding": 16,
				},
			},
		},
	},
	{
		ID:   "media",
		Name: "媒体组件",
		Components: []models.ComponentDefinition{
			{
				Type:        "image",
				Name:        "图片",
				Icon:        "image",
				Description: "图片展示",
				DefaultSettings: map[string]interface{}{
					"width":     "100%",
					"height":    "auto",
					"objectFit": "cover",
				},
			},
			{
				Type:        "video",
				Name:        "视频",
				Icon:        "video",
				Description: "视频播放器",
				DefaultSettings: map[string]interface{}{
					"autoplay": false,
					"controls": true,
					"loop":     false,
				},
			},
			{
				Type:        "carousel",
				Name:        "轮播图",
				Icon:        "carousel",
				Description: "图片轮播",
				DefaultSettings: map[string]interface{}{
					"autoplay":   true,
					"interval":   3000,
					"indicators": true,
					"arrows":     true,
				},
			},
		},
	},
}

// 模拟组件数据
var components = map[string][]models.Component{
	"section1": {
		{
			ID:        "comp1",
			SectionID: "section1",
			Type:      "heading",
			Name:      "网站标题",
			Settings: map[string]interface{}{
				"level":     "h1",
				"textAlign": "center",
			},
			Content: map[string]interface{}{
				"text": "企业官网",
			},
			Style: map[string]interface{}{
				"marginBottom": "20px",
			},
			SortOrder: 0,
		},
		{
			ID:        "comp2",
			SectionID: "section1",
			Type:      "text",
			Name:      "欢迎文本",
			Settings: map[string]interface{}{
				"textAlign": "center",
				"fontSize":  "18px",
			},
			Content: map[string]interface{}{
				"text": "欢迎访问我们的企业官网",
			},
			Style: map[string]interface{}{
				"color": "#666",
			},
			SortOrder: 1,
		},
	},
	"section2": {
		{
			ID:        "comp3",
			SectionID: "section2",
			Type:      "image",
			Name:      "宣传图片",
			Settings: map[string]interface{}{
				"width":     "100%",
				"objectFit": "cover",
			},
			Content: map[string]interface{}{
				"src": "/img/banner.jpg",
				"alt": "企业宣传图",
			},
			Style: map[string]interface{}{
				"borderRadius": "8px",
			},
			SortOrder: 0,
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

// ListComponentCategories 获取组件分类列表
func ListComponentCategories() ([]models.ComponentCategory, error) {
	return componentCategories, nil
}

// GetComponentDefinition 获取组件定义
func GetComponentDefinition(componentType string) (models.ComponentDefinition, error) {
	for _, category := range componentCategories {
		for _, definition := range category.Components {
			if definition.Type == componentType {
				return definition, nil
			}
		}
	}

	return models.ComponentDefinition{}, errors.New("组件类型不存在")
}

// AddComponent 添加组件到区块
func AddComponent(siteID string, pageID string, sectionID string, component models.Component) (models.Component, error) {
	// 生成一个简单的ID
	component.ID = fmt.Sprintf("comp%d", len(components[sectionID])+1)

	// 设置排序顺序
	if sectionComponents, exists := components[sectionID]; exists {
		component.SortOrder = len(sectionComponents)
	} else {
		component.SortOrder = 0
		components[sectionID] = []models.Component{}
	}

	// 添加到组件集合
	components[sectionID] = append(components[sectionID], component)

	return component, nil
}

// UpdateComponent 更新组件
func UpdateComponent(siteID string, pageID string, sectionID string, component models.Component) (models.Component, error) {
	if sectionComponents, exists := components[sectionID]; exists {
		for i, c := range sectionComponents {
			if c.ID == component.ID {
				// 保持一些不应该被客户端更新的字段
				component.SortOrder = c.SortOrder

				sectionComponents[i] = component
				components[sectionID] = sectionComponents

				return component, nil
			}
		}
		return models.Component{}, errors.New("组件不存在")
	}

	return models.Component{}, errors.New("区块没有组件")
}

// DeleteComponent 删除组件
func DeleteComponent(siteID string, pageID string, sectionID string, componentID string) error {
	if sectionComponents, exists := components[sectionID]; exists {
		for i, component := range sectionComponents {
			if component.ID == componentID {
				// 从切片中删除
				components[sectionID] = append(sectionComponents[:i], sectionComponents[i+1:]...)

				// 重新调整排序
				for j := i; j < len(components[sectionID]); j++ {
					components[sectionID][j].SortOrder = j
				}

				return nil
			}
		}
		return errors.New("组件不存在")
	}

	return errors.New("区块没有组件")
}

// ReorderComponents 重新排序组件
func ReorderComponents(siteID string, pageID string, sectionID string, componentOrder []string) error {
	if sectionComponents, exists := components[sectionID]; exists {
		if len(componentOrder) != len(sectionComponents) {
			return errors.New("组件数量不匹配")
		}

		// 创建新的排序组件切片
		newOrder := make([]models.Component, len(sectionComponents))
		for i, componentID := range componentOrder {
			found := false
			for _, component := range sectionComponents {
				if component.ID == componentID {
					component.SortOrder = i
					newOrder[i] = component
					found = true
					break
				}
			}
			if !found {
				return errors.New("无效的组件ID")
			}
		}

		components[sectionID] = newOrder
		return nil
	}

	return errors.New("区块没有组件")
}
