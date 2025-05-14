package service

import (
	"errors"
	"strings"
	"wz-backend-go/models"
)

// 模拟模板数据
var templates = []models.SiteTemplate{
	{
		ID:          "t1",
		Name:        "企业展示",
		Thumbnail:   "/img/template1.jpg",
		Description: "适合企业官网使用的模板",
		Config:      `{"pages":[{"name":"首页","layout":"default"},{"name":"关于我们","layout":"default"},{"name":"产品服务","layout":"default"},{"name":"联系我们","layout":"default"}]}`,
	},
	{
		ID:          "t2",
		Name:        "产品展示",
		Thumbnail:   "/img/template2.jpg",
		Description: "重点突出产品的模板",
		Config:      `{"pages":[{"name":"首页","layout":"full-width"},{"name":"产品目录","layout":"sidebar"},{"name":"产品详情","layout":"default"},{"name":"联系我们","layout":"default"}]}`,
	},
	{
		ID:          "t3",
		Name:        "简约风格",
		Thumbnail:   "/img/template3.jpg",
		Description: "简洁大方的设计风格",
		Config:      `{"pages":[{"name":"首页","layout":"default"},{"name":"博客","layout":"sidebar"},{"name":"作品集","layout":"full-width"},{"name":"关于","layout":"default"}]}`,
	},
}

// ListTemplates 获取模板列表
func ListTemplates(category string) ([]models.SiteTemplate, error) {
	if category == "" {
		return templates, nil
	}

	var result []models.SiteTemplate
	categoryLower := strings.ToLower(category)

	for _, template := range templates {
		// 简单过滤，实际中可能需要在模板中添加分类字段
		if strings.Contains(strings.ToLower(template.Name), categoryLower) ||
			strings.Contains(strings.ToLower(template.Description), categoryLower) {
			result = append(result, template)
		}
	}

	return result, nil
}

// GetTemplate 获取模板详情
func GetTemplate(templateID string) (models.SiteTemplate, error) {
	for _, template := range templates {
		if template.ID == templateID {
			return template, nil
		}
	}

	return models.SiteTemplate{}, errors.New("模板不存在")
}
