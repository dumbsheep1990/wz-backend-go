package handlers

import (
	"net/http"
	"wz-backend-go/services/render-service/service"

	"github.com/gin-gonic/gin"
)

// RenderSiteByDomain 根据域名渲染站点首页
func RenderSiteByDomain(c *gin.Context) {
	// 从查询参数或请求头获取域名
	domain := c.Query("domain")
	if domain == "" {
		domain = c.Request.Host
	}

	// 通过域名获取站点信息
	site, err := service.GetSiteByDomain(domain)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "站点不存在"})
		return
	}

	// 获取首页
	homepage, err := service.GetHomePage(site.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "首页不存在"})
		return
	}

	// 生成HTML
	htmlContent, err := service.GeneratePageHTML(site, homepage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回HTML内容
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}

// RenderPageBySlug 根据站点ID和页面Slug渲染特定页面
func RenderPageBySlug(c *gin.Context) {
	siteID := c.Param("siteId")
	slug := c.Param("slug")

	// 检查站点状态是否为已发布
	if !service.IsSitePublished(siteID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "站点不存在或未发布"})
		return
	}

	// 获取站点信息
	site, err := service.GetSite(siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "站点不存在"})
		return
	}

	// 根据slug获取页面
	page, err := service.GetPageBySlug(siteID, slug)
	if err != nil {
		// 如果没找到页面，尝试渲染首页
		if slug == "" || slug == "index" || slug == "home" {
			homepage, err := service.GetHomePage(siteID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "页面不存在"})
				return
			}

			// 生成首页HTML
			htmlContent, err := service.GeneratePageHTML(site, homepage)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, htmlContent)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "页面不存在"})
		return
	}

	// 生成页面HTML
	htmlContent, err := service.GeneratePageHTML(site, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回HTML内容
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}
