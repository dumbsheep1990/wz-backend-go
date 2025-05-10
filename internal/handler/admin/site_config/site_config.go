package site_config

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/service"
	"wz-backend-go/internal/types"
)

// SiteConfigHandler 处理站点配置相关请求
type SiteConfigHandler struct {
	siteConfigService service.SiteConfigService
	logger            logx.Logger
}

// NewSiteConfigHandler 创建站点配置处理器
func NewSiteConfigHandler(siteConfigService service.SiteConfigService, logger logx.Logger) *SiteConfigHandler {
	return &SiteConfigHandler{
		siteConfigService: siteConfigService,
		logger:            logger,
	}
}

// GetSiteConfig 获取站点配置
func (h *SiteConfigHandler) GetSiteConfig(c *gin.Context) {
	// 获取租户ID
	tenantID, _ := strconv.ParseInt(c.GetHeader("TenantID"), 10, 64)

	config, err := h.siteConfigService.GetSiteConfig(tenantID)
	if err != nil {
		h.logger.Errorf("获取站点配置失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取站点配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": config,
	})
}

// UpdateSiteConfig 更新站点配置
func (h *SiteConfigHandler) UpdateSiteConfig(c *gin.Context) {
	var req types.SiteConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误:" + err.Error()})
		return
	}

	// 获取租户ID
	tenantID, _ := strconv.ParseInt(c.GetHeader("TenantID"), 10, 64)

	config := &domain.SiteConfig{
		SiteName:       req.SiteName,
		SiteLogo:       req.SiteLogo,
		SeoTitle:       req.SeoTitle,
		SeoKeywords:    req.SeoKeywords,
		SeoDescription: req.SeoDescription,
		IcpNumber:      req.IcpNumber,
		Copyright:      req.Copyright,
		ThemeID:        req.ThemeID,
		ContactEmail:   req.ContactEmail,
		ContactPhone:   req.ContactPhone,
		Address:        req.Address,
		TenantID:       tenantID,
	}

	err := h.siteConfigService.UpdateSiteConfig(config)
	if err != nil {
		h.logger.Errorf("更新站点配置失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新站点配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}
