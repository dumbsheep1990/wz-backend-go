package mock

import (
	"context"
	"errors"
	"time"
)

// MockSiteService 模拟的站点服务
type MockSiteService struct {
	sites map[string]interface{}
	domainToSiteID map[string]string
}

// NewMockSiteService 创建一个新的模拟站点服务
func NewMockSiteService() *MockSiteService {
	service := &MockSiteService{
		sites: make(map[string]interface{}),
		domainToSiteID: make(map[string]string),
	}

	// 添加一些模拟数据
	siteID := "site-001"
	service.sites[siteID] = map[string]interface{}{
		"id":        siteID,
		"name":      "万知网",
		"domain":    "wz.qq.com",
		"published": true,
		"created_at": time.Now().Add(-30 * 24 * time.Hour),
		"updated_at": time.Now().Add(-2 * 24 * time.Hour),
		"logo":       "https://wz.qq.com/logo.png",
		"theme":      "default",
		"settings": map[string]interface{}{
			"colors": map[string]string{
				"primary":   "#4CAF50",
				"secondary": "#FFC107",
				"accent":    "#2196F3",
			},
		},
		"footer": map[string]interface{}{
			"copyright": "版权所有©南京万知园信息工程有限公司 万知网",
			"address":   "接洽：江苏省盐城市通港路11号综合楼一楼1号",
			"icp":       "苏ICP备19049108号 | 苏公网安备 32010402000846号",
			"service":   "客户服务热线：0515-88200000",
			"report":    "违法和不良信息举报电话：0515-88200000 | 举报邮箱：2895915959@qq.com",
		},
		"categories": []map[string]interface{}{
			{"id": "cat-001", "name": "同用", "slug": "tongyong"},
			{"id": "cat-002", "name": "同好", "slug": "tonghao"},
			{"id": "cat-003", "name": "同购", "slug": "tonggou"},
			{"id": "cat-004", "name": "同年", "slug": "tongnian"},
			{"id": "cat-005", "name": "同游", "slug": "tongyou"},
			{"id": "cat-006", "name": "同在", "slug": "tongzai"},
			{"id": "cat-007", "name": "同市", "slug": "tongshi"},
			{"id": "cat-008", "name": "同企", "slug": "tongqi"},
			{"id": "cat-009", "name": "同亲", "slug": "tongqin"},
			{"id": "cat-010", "name": "同班", "slug": "tongban"},
			{"id": "cat-011", "name": "同师", "slug": "tongshi"},
			{"id": "cat-012", "name": "同业", "slug": "tongye"},
			{"id": "cat-013", "name": "同网", "slug": "tongwang"},
			{"id": "cat-014", "name": "同工", "slug": "tonggong"},
			{"id": "cat-015", "name": "同务", "slug": "tongwu"},
			{"id": "cat-016", "name": "同艺", "slug": "tongyi"},
			{"id": "cat-017", "name": "同玩", "slug": "tongwan"},
			{"id": "cat-018", "name": "同闲", "slug": "tongxian"},
			{"id": "cat-019", "name": "同拍", "slug": "tongpai"},
			{"id": "cat-020", "name": "同乡", "slug": "tongxiang"},
			{"id": "cat-021", "name": "同学", "slug": "tongxue"},
		},
	}

	// 添加域名映射
	service.domainToSiteID["wz.qq.com"] = siteID

	return service
}

// GetSiteByID 根据ID获取站点
func (s *MockSiteService) GetSiteByID(ctx context.Context, siteID string) (interface{}, error) {
	if site, exists := s.sites[siteID]; exists {
		return site, nil
	}
	return nil, errors.New("站点不存在")
}

// GetSiteByDomain 根据域名获取站点
func (s *MockSiteService) GetSiteByDomain(ctx context.Context, domain string) (interface{}, error) {
	if siteID, exists := s.domainToSiteID[domain]; exists {
		return s.GetSiteByID(ctx, siteID)
	}
	return nil, errors.New("站点不存在")
}

// IsSitePublished 检查站点是否已发布
func (s *MockSiteService) IsSitePublished(ctx context.Context, siteID string) bool {
	if site, exists := s.sites[siteID]; exists {
		if siteMap, ok := site.(map[string]interface{}); ok {
			if published, ok := siteMap["published"].(bool); ok {
				return published
			}
		}
	}
	return false
}

// MockPageService 模拟的页面服务
type MockPageService struct {
	pages map[string]interface{}
	sitePages map[string][]string
	slugToPageID map[string]map[string]string
	siteHomepages map[string]string
}

// NewMockPageService 创建一个新的模拟页面服务
func NewMockPageService() *MockPageService {
	service := &MockPageService{
		pages:         make(map[string]interface{}),
		sitePages:     make(map[string][]string),
		slugToPageID:  make(map[string]map[string]string),
		siteHomepages: make(map[string]string),
	}

	siteID := "site-001"

	// 初始化站点的slug到页面ID映射
	service.slugToPageID[siteID] = make(map[string]string)

	// 添加首页
	homepageID := "page-001"
	service.pages[homepageID] = map[string]interface{}{
		"id":        homepageID,
		"siteId":    siteID,
		"title":     "万知网 - 首页",
		"slug":      "home",
		"templateId": "template-001",
		"created_at": time.Now().Add(-30 * 24 * time.Hour),
		"updated_at": time.Now().Add(-1 * 24 * time.Hour),
		"published": true,
		"isHomepage": true,
		"sections": []map[string]interface{}{
			{
				"id": "section-001",
				"type": "hero",
				"content": "<h1>欢迎来到万知网</h1><p>连接同类，发现世界</p>",
				"style": "primary",
			},
			{
				"id": "section-002",
				"type": "category-grid",
				"title": "热门分类",
				"categories": []string{"同用", "同好", "同购", "同年", "同游", "同在"},
			},
			{
				"id": "section-003",
				"type": "post-list",
				"title": "最新帖子",
				"postIds": []string{"post-001", "post-002", "post-003"},
			},
		},
		"meta": map[string]string{
			"description": "万知网 - 连接同类，发现世界",
			"keywords":    "万知网,同类,社区,交流",
		},
	}

	// 添加类别页面
	categoryPageID := "page-002"
	service.pages[categoryPageID] = map[string]interface{}{
		"id":        categoryPageID,
		"siteId":    siteID,
		"title":     "万知网 - 分类导航",
		"slug":      "categories",
		"templateId": "template-002",
		"created_at": time.Now().Add(-25 * 24 * time.Hour),
		"updated_at": time.Now().Add(-1 * 24 * time.Hour),
		"published": true,
		"sections": []map[string]interface{}{
			{
				"id": "section-004",
				"type": "category-list",
				"title": "所有分类",
				"showDescription": true,
			},
		},
		"meta": map[string]string{
			"description": "万知网分类导航 - 探索所有分类",
			"keywords":    "万知网,分类,导航,同类",
		},
	}

	// 添加"同乡"页面
	tongxiangPageID := "page-003"
	service.pages[tongxiangPageID] = map[string]interface{}{
		"id":        tongxiangPageID,
		"siteId":    siteID,
		"title":     "万知网 - 同乡圈",
		"slug":      "tongxiang",
		"templateId": "template-003",
		"created_at": time.Now().Add(-20 * 24 * time.Hour),
		"updated_at": time.Now().Add(-2 * 24 * time.Hour),
		"published": true,
		"sections": []map[string]interface{}{
			{
				"id": "section-005",
				"type": "post-grid",
				"title": "同乡帖子",
				"description": "查找来自同一地区的老乡",
				"postIds": []string{"post-004", "post-005", "post-006"},
			},
		},
		"meta": map[string]string{
			"description": "万知网同乡圈 - 与您的老乡建立联系",
			"keywords":    "万知网,同乡,老乡,地区,联系",
		},
	}

	// 添加"同学"页面
	tongxuePageID := "page-004"
	service.pages[tongxuePageID] = map[string]interface{}{
		"id":        tongxuePageID,
		"siteId":    siteID,
		"title":     "万知网 - 同学圈",
		"slug":      "tongxue",
		"templateId": "template-003",
		"created_at": time.Now().Add(-18 * 24 * time.Hour),
		"updated_at": time.Now().Add(-3 * 24 * time.Hour),
		"published": true,
		"sections": []map[string]interface{}{
			{
				"id": "section-006",
				"type": "post-grid",
				"title": "同学帖子",
				"description": "找到您的同学和校友",
				"postIds": []string{"post-007", "post-008", "post-009"},
			},
		},
		"meta": map[string]string{
			"description": "万知网同学圈 - 与您的同学和校友建立联系",
			"keywords":    "万知网,同学,校友,学校,联系",
		},
	}

	// 添加关于页面
	aboutPageID := "page-005"
	service.pages[aboutPageID] = map[string]interface{}{
		"id":        aboutPageID,
		"siteId":    siteID,
		"title":     "关于万知网",
		"slug":      "about",
		"templateId": "template-004",
		"created_at": time.Now().Add(-15 * 24 * time.Hour),
		"updated_at": time.Now().Add(-5 * 24 * time.Hour),
		"published": true,
		"sections": []map[string]interface{}{
			{
				"id": "section-007",
				"type": "content",
				"content": "<h1>关于万知网</h1><p>万知网是一个连接同类人的平台，帮助用户找到具有相同背景、兴趣和经历的人。</p>",
			},
		},
		"meta": map[string]string{
			"description": "了解万知网 - 我们的使命和价值观",
			"keywords":    "万知网,关于,使命,价值观",
		},
	}

	// 设置站点的页面
	service.sitePages[siteID] = []string{homepageID, categoryPageID, tongxiangPageID, tongxuePageID, aboutPageID}

	// 设置slug到页面ID的映射
	service.slugToPageID[siteID]["home"] = homepageID
	service.slugToPageID[siteID]["categories"] = categoryPageID
	service.slugToPageID[siteID]["tongxiang"] = tongxiangPageID
	service.slugToPageID[siteID]["tongxue"] = tongxuePageID
	service.slugToPageID[siteID]["about"] = aboutPageID

	// 设置站点的首页
	service.siteHomepages[siteID] = homepageID

	return service
}

// GetPageByID 根据ID获取页面
func (s *MockPageService) GetPageByID(ctx context.Context, pageID string) (interface{}, error) {
	if page, exists := s.pages[pageID]; exists {
		return page, nil
	}
	return nil, errors.New("页面不存在")
}

// GetPageBySlug 根据slug获取页面
func (s *MockPageService) GetPageBySlug(ctx context.Context, siteID string, slug string) (interface{}, error) {
	if siteSlugMap, exists := s.slugToPageID[siteID]; exists {
		if pageID, exists := siteSlugMap[slug]; exists {
			return s.GetPageByID(ctx, pageID)
		}
	}
	return nil, errors.New("页面不存在")
}

// GetHomePage 获取站点的首页
func (s *MockPageService) GetHomePage(ctx context.Context, siteID string) (interface{}, error) {
	if homepageID, exists := s.siteHomepages[siteID]; exists {
		return s.GetPageByID(ctx, homepageID)
	}
	return nil, errors.New("首页不存在")
}
