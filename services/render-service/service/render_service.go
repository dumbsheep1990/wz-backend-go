package service

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"wz-backend-go/models"
)

// 模拟站点数据
var sites = []models.Site{
	{
		ID:          "1",
		Name:        "企业展示站点",
		Description: "适合企业官网使用的模板",
		Domain:      "company.wanzhimarket.com",
		Logo:        "/img/logo1.png",
		Favicon:     "/img/favicon1.ico",
		TenantID:    "tenant_456",
		Theme: models.ThemeConfig{
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
		Navigation: models.Navigation{
			Type:  "horizontal",
			Items: []models.NavigationItem{},
			Style: map[string]string{},
		},
		Status: "published",
	},
}

// 模拟页面数据
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
		},
	},
}

// 区块数据
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

// 组件数据
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

// 模拟站点租户权限
var siteTenants = map[string]string{
	"1": "tenant_456",
}

// CheckSiteAccess 检查站点访问权限
func CheckSiteAccess(siteID string, tenantID string) bool {
	if tenant, exists := siteTenants[siteID]; exists {
		return tenant == tenantID
	}
	return false
}

// GetSiteWithAllPages 获取站点及其所有页面
func GetSiteWithAllPages(siteID string) (models.Site, error) {
	for _, site := range sites {
		if site.ID == siteID {
			// 找到站点后，加载所有页面
			if sitePages, exists := pages[siteID]; exists {
				// 加载页面中的区块和组件
				for i, page := range sitePages {
					if pageSections, exists := sections[page.ID]; exists {
						// 加载每个区块的组件
						for j, section := range pageSections {
							if sectionComponents, exists := components[section.ID]; exists {
								pageSections[j].Components = sectionComponents
							}
						}
						sitePages[i].Sections = pageSections
					}
				}
				site.Pages = sitePages
			}
			return site, nil
		}
	}
	return models.Site{}, errors.New("站点不存在")
}

// GetSiteAndPage 获取站点和特定页面
func GetSiteAndPage(siteID string, pageID string) (models.Site, models.Page, error) {
	// 获取站点
	site, err := GetSite(siteID)
	if err != nil {
		return models.Site{}, models.Page{}, err
	}

	// 获取页面
	for _, page := range pages[siteID] {
		if page.ID == pageID {
			// 加载页面的区块和组件
			if pageSections, exists := sections[page.ID]; exists {
				// 加载每个区块的组件
				for j, section := range pageSections {
					if sectionComponents, exists := components[section.ID]; exists {
						pageSections[j].Components = sectionComponents
					}
				}
				page.Sections = pageSections
			}
			return site, page, nil
		}
	}

	return models.Site{}, models.Page{}, errors.New("页面不存在")
}

// GenerateSitePreview 生成站点预览HTML
func GenerateSitePreview(site models.Site, device string) (string, error) {
	// 找到首页
	var homepage models.Page
	for _, page := range site.Pages {
		if page.IsHomepage {
			homepage = page
			break
		}
	}

	if homepage.ID == "" {
		return "", errors.New("找不到首页")
	}

	// 使用首页生成预览
	return GeneratePagePreview(site, homepage, device)
}

// GeneratePagePreview 生成页面预览HTML
func GeneratePagePreview(site models.Site, page models.Page, device string) (string, error) {
	// 通用HTML模板
	htmlTemplate := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Page.Title }}</title>
    <meta name="description" content="{{ .Page.Description }}">
    <link rel="icon" href="{{ .Site.Favicon }}" type="image/x-icon">
    <style>
        /* 预览模式样式 */
        body {
            font-family: {{ .Site.Theme.FontFamily }};
            color: {{ .Site.Theme.TextColor }};
            background-color: {{ .Site.Theme.BackgroundColor }};
            margin: 0;
            padding: 0;
        }
        .preview-bar {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            background-color: #333;
            color: white;
            padding: 10px;
            text-align: center;
            z-index: 1000;
        }
        .preview-bar span {
            margin-right: 15px;
        }
        .content {
            margin-top: 50px;
            padding: 20px;
        }
        /* 响应式样式 */
        {{ if eq .Device "mobile" }}
        .content {
            max-width: 360px;
            margin-left: auto;
            margin-right: auto;
        }
        {{ else if eq .Device "tablet" }}
        .content {
            max-width: 768px;
            margin-left: auto;
            margin-right: auto;
        }
        {{ else }}
        .content {
            max-width: 1200px;
            margin-left: auto;
            margin-right: auto;
        }
        {{ end }}
        
        /* 站点主题样式 */
        a {
            color: {{ .Site.Theme.PrimaryColor }};
        }
        .btn {
            background-color: {{ .Site.Theme.PrimaryColor }};
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: {{ .Site.Theme.BorderRadius }};
            cursor: pointer;
        }
        .section {
            margin-bottom: 30px;
        }
        /* 自定义CSS */
        {{ .Site.Theme.CustomCSS }}
    </style>
</head>
<body>
    <div class="preview-bar">
        <span>预览模式: {{ .Device }}</span>
        <span>页面: {{ .Page.Name }}</span>
    </div>
    
    <div class="content">
        <!-- 导航 -->
        <header>
            <div class="logo">
                <img src="{{ .Site.Logo }}" alt="{{ .Site.Name }}" style="max-height: 50px;">
            </div>
            <nav>
                <ul style="display: flex; list-style: none; gap: 20px;">
                    {{ range .Site.Pages }}
                    <li><a href="#{{ .Slug }}">{{ .Name }}</a></li>
                    {{ end }}
                </ul>
            </nav>
        </header>
        
        <!-- 页面内容 -->
        {{ range .Page.Sections }}
        <div class="section" id="{{ .ID }}">
            <h2>{{ .Title }}</h2>
            {{ range .Components }}
            <div class="component">
                {{ if eq .Type "heading" }}
                <{{ .Settings.level }}>{{ .Content.text }}</{{ .Settings.level }}>
                {{ else if eq .Type "text" }}
                <p>{{ .Content.text }}</p>
                {{ else if eq .Type "image" }}
                <img src="{{ .Content.src }}" alt="{{ .Content.alt }}" style="max-width: 100%;">
                {{ else if eq .Type "button" }}
                <button class="btn">{{ .Content.text }}</button>
                {{ end }}
            </div>
            {{ end }}
        </div>
        {{ end }}
        
        <!-- 页脚 -->
        <footer>
            <p>© {{ .Site.Name }}</p>
        </footer>
    </div>
</body>
</html>
`

	// 准备模板数据
	templateData := map[string]interface{}{
		"Site":   site,
		"Page":   page,
		"Device": device,
	}

	// 解析模板
	tmpl, err := template.New("preview").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	// 执行模板
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// GetSiteByDomain 根据域名获取站点
func GetSiteByDomain(domain string) (models.Site, error) {
	for _, site := range sites {
		if site.Domain == domain && site.Status == "published" {
			return site, nil
		}
	}
	return models.Site{}, errors.New("站点不存在或未发布")
}

// GetSite 获取站点信息
func GetSite(siteID string) (models.Site, error) {
	for _, site := range sites {
		if site.ID == siteID {
			return site, nil
		}
	}
	return models.Site{}, errors.New("站点不存在")
}

// GetHomePage 获取站点首页
func GetHomePage(siteID string) (models.Page, error) {
	sitePages, exists := pages[siteID]
	if !exists {
		return models.Page{}, errors.New("站点没有页面")
	}

	for _, page := range sitePages {
		if page.IsHomepage {
			// 加载页面区块和组件
			if pageSections, exists := sections[page.ID]; exists {
				// 加载每个区块的组件
				for j, section := range pageSections {
					if sectionComponents, exists := components[section.ID]; exists {
						pageSections[j].Components = sectionComponents
					}
				}
				page.Sections = pageSections
			}
			return page, nil
		}
	}

	return models.Page{}, errors.New("找不到首页")
}

// GetPageBySlug 通过slug获取页面
func GetPageBySlug(siteID string, slug string) (models.Page, error) {
	sitePages, exists := pages[siteID]
	if !exists {
		return models.Page{}, errors.New("站点没有页面")
	}

	// 规范化slug
	slug = strings.ToLower(slug)

	for _, page := range sitePages {
		if strings.ToLower(page.Slug) == slug {
			// 加载页面区块和组件
			if pageSections, exists := sections[page.ID]; exists {
				// 加载每个区块的组件
				for j, section := range pageSections {
					if sectionComponents, exists := components[section.ID]; exists {
						pageSections[j].Components = sectionComponents
					}
				}
				page.Sections = pageSections
			}
			return page, nil
		}
	}

	return models.Page{}, fmt.Errorf("找不到slug为%s的页面", slug)
}

// IsSitePublished 检查站点是否已发布
func IsSitePublished(siteID string) bool {
	for _, site := range sites {
		if site.ID == siteID && site.Status == "published" {
			return true
		}
	}
	return false
}

// GeneratePageHTML 生成页面HTML
func GeneratePageHTML(site models.Site, page models.Page) (string, error) {
	// HTML模板 - 这里简化处理，实际中会更复杂
	htmlTemplate := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Page.Title }}</title>
    <meta name="description" content="{{ .Page.Description }}">
    <meta name="keywords" content="{{ join .Page.Keywords ", " }}">
    <link rel="icon" href="{{ .Site.Favicon }}" type="image/x-icon">
    <style>
        /* 站点样式 */
        body {
            font-family: {{ .Site.Theme.FontFamily }};
            color: {{ .Site.Theme.TextColor }};
            background-color: {{ .Site.Theme.BackgroundColor }};
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
        }
        a {
            color: {{ .Site.Theme.PrimaryColor }};
        }
        .btn {
            background-color: {{ .Site.Theme.PrimaryColor }};
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: {{ .Site.Theme.BorderRadius }};
            cursor: pointer;
        }
        .section {
            margin-bottom: 30px;
        }
        /* 自定义CSS */
        {{ .Site.Theme.CustomCSS }}
    </style>
</head>
<body>
    <div class="container">
        <!-- 导航 -->
        <header>
            <div class="logo">
                <img src="{{ .Site.Logo }}" alt="{{ .Site.Name }}" style="max-height: 50px;">
            </div>
            <nav>
                <ul style="display: flex; list-style: none; gap: 20px;">
                    {{ range .Site.Navigation.Items }}
                    <li><a href="{{ .Link }}">{{ .Label }}</a></li>
                    {{ end }}
                </ul>
            </nav>
        </header>
        
        <!-- 页面内容 -->
        {{ range .Page.Sections }}
        <div class="section" id="{{ .ID }}">
            <h2>{{ .Title }}</h2>
            {{ range .Components }}
            <div class="component">
                {{ if eq .Type "heading" }}
                <{{ .Settings.level }}>{{ .Content.text }}</{{ .Settings.level }}>
                {{ else if eq .Type "text" }}
                <p>{{ .Content.text }}</p>
                {{ else if eq .Type "image" }}
                <img src="{{ .Content.src }}" alt="{{ .Content.alt }}" style="max-width: 100%;">
                {{ else if eq .Type "button" }}
                <button class="btn">{{ .Content.text }}</button>
                {{ end }}
            </div>
            {{ end }}
        </div>
        {{ end }}
        
        <!-- 页脚 -->
        <footer>
            <p>© {{ .Site.Name }}</p>
        </footer>
    </div>
</body>
</html>
`

	// 注册自定义函数
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	// 准备模板数据
	templateData := map[string]interface{}{
		"Site": site,
		"Page": page,
	}

	// 解析模板
	tmpl, err := template.New("page").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	// 执行模板
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
