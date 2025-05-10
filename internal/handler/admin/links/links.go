package links

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/service"
	"wz-backend-go/internal/types"
)

// LinksHandler 处理友情链接相关请求
type LinksHandler struct {
	linkService service.LinkService
	logger      logx.Logger
}

// NewLinksHandler 创建友情链接处理器
func NewLinksHandler(linkService service.LinkService, logger logx.Logger) *LinksHandler {
	return &LinksHandler{
		linkService: linkService,
		logger:      logger,
	}
}

// ListLinks 获取友情链接列表
func (h *LinksHandler) ListLinks(c *gin.Context) {
	var req types.LinkListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误:" + err.Error()})
		return
	}

	// 默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 获取租户ID
	tenantID, _ := strconv.ParseInt(c.GetHeader("TenantID"), 10, 64)

	// 查询条件
	query := map[string]interface{}{
		"tenant_id": tenantID,
	}
	if req.Status != nil {
		query["status"] = *req.Status
	}
	if req.Name != "" {
		query["name"] = req.Name
	}

	links, total, err := h.linkService.ListLinks(req.Page, req.PageSize, query)
	if err != nil {
		h.logger.Errorf("获取友情链接列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取友情链接列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  links,
			"total": total,
			"page":  req.Page,
			"size":  req.PageSize,
		},
	})
}

// GetLink 获取友情链接详情
func (h *LinksHandler) GetLink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID参数错误"})
		return
	}

	link, err := h.linkService.GetLinkById(id)
	if err != nil {
		h.logger.Errorf("获取友情链接详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取友情链接详情失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": link,
	})
}

// CreateLink 创建友情链接
func (h *LinksHandler) CreateLink(c *gin.Context) {
	var req types.LinkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误:" + err.Error()})
		return
	}

	// 获取租户ID
	tenantID, _ := strconv.ParseInt(c.GetHeader("TenantID"), 10, 64)

	link := &domain.Link{
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Sort:        req.Sort,
		Status:      req.Status,
		Description: req.Description,
		TenantID:    tenantID,
	}

	id, err := h.linkService.CreateLink(link)
	if err != nil {
		h.logger.Errorf("创建友情链接失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建友情链接失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": gin.H{"id": id},
	})
}

// UpdateLink 更新友情链接
func (h *LinksHandler) UpdateLink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID参数错误"})
		return
	}

	var req types.LinkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误:" + err.Error()})
		return
	}

	// 获取租户ID
	tenantID, _ := strconv.ParseInt(c.GetHeader("TenantID"), 10, 64)

	link := &domain.Link{
		ID:          id,
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Sort:        req.Sort,
		Status:      req.Status,
		Description: req.Description,
		TenantID:    tenantID,
	}

	err = h.linkService.UpdateLink(link)
	if err != nil {
		h.logger.Errorf("更新友情链接失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新友情链接失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteLink 删除友情链接
func (h *LinksHandler) DeleteLink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID参数错误"})
		return
	}

	err = h.linkService.DeleteLink(id)
	if err != nil {
		h.logger.Errorf("删除友情链接失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除友情链接失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
