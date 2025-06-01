package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wz-backend-go/internal/application/render/dto"
	appService "wz-backend-go/internal/application/render/service"
)

// RenderController 渲染控制器
type RenderController struct {
	renderAppService *appService.RenderApplicationService
}

// NewRenderController 创建一个新的渲染控制器
func NewRenderController(renderAppService *appService.RenderApplicationService) *RenderController {
	return &RenderController{
		renderAppService: renderAppService,
	}
}

// RenderTemplate 渲染模板
// @Summary 渲染模板
// @Description 根据模板ID和上下文数据渲染内容
// @Tags 渲染
// @Accept json
// @Produce json
// @Param request body dto.RenderRequestDTO true "渲染请求"
// @Success 200 {object} dto.RenderResponseDTO
// @Failure 400 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /api/v1/render/template [post]
func (c *RenderController) RenderTemplate(ctx *gin.Context) {
	var request dto.RenderRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.RenderErrorDTO{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
			Details: err.Error(),
		})
		return
	}

	result, err := c.renderAppService.RenderTemplate(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.RenderErrorDTO{
			Code:    http.StatusInternalServerError,
			Message: "渲染失败",
			Details: err.Error(),
		})
		return
	}

	// 根据格式返回不同的响应
	if request.Format == "json" {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.Header("Content-Type", result.ContentType)
		ctx.String(http.StatusOK, result.Content)
	}
}

// RenderPageBySlug 根据页面Slug渲染页面
// @Summary 渲染页面
// @Description 根据站点ID和页面Slug渲染页面
// @Tags 渲染
// @Accept json
// @Produce html
// @Param siteId path string true "站点ID"
// @Param slug path string true "页面Slug"
// @Success 200 {string} string "渲染的HTML内容"
// @Failure 400 {object} dto.RenderErrorDTO
// @Failure 404 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /render/sites/{siteId}/{slug} [get]
func (c *RenderController) RenderPageBySlug(ctx *gin.Context) {
	siteID := ctx.Param("siteId")
	slug := ctx.Param("slug")

	request := dto.PageRenderRequestDTO{
		SiteID: siteID,
		Slug:   slug,
		Format: ctx.Query("format"),
	}

	result, err := c.renderAppService.RenderPage(ctx, request)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "页面不存在" || err.Error() == "站点不存在" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.RenderErrorDTO{
			Code:    status,
			Message: "渲染失败",
			Details: err.Error(),
		})
		return
	}

	ctx.Header("Content-Type", result.ContentType)
	ctx.String(http.StatusOK, result.Content)
}

// RenderSiteByDomain 根据域名渲染站点
// @Summary 根据域名渲染站点
// @Description 根据域名渲染站点首页
// @Tags 渲染
// @Accept json
// @Produce html
// @Param domain query string false "域名"
// @Success 200 {string} string "渲染的HTML内容"
// @Failure 400 {object} dto.RenderErrorDTO
// @Failure 404 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /render/site [get]
func (c *RenderController) RenderSiteByDomain(ctx *gin.Context) {
	domain := ctx.Query("domain")
	if domain == "" {
		domain = ctx.Request.Host
	}

	request := dto.PageRenderRequestDTO{
		Domain: domain,
		Format: ctx.Query("format"),
	}

	result, err := c.renderAppService.RenderPage(ctx, request)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "站点不存在" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.RenderErrorDTO{
			Code:    status,
			Message: "渲染失败",
			Details: err.Error(),
		})
		return
	}

	ctx.Header("Content-Type", result.ContentType)
	ctx.String(http.StatusOK, result.Content)
}

// PreviewSite 预览站点
// @Summary 预览站点
// @Description 预览站点首页
// @Tags 预览
// @Accept json
// @Produce html
// @Param siteId path string true "站点ID"
// @Success 200 {string} string "渲染的HTML内容"
// @Failure 400 {object} dto.RenderErrorDTO
// @Failure 404 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /api/v1/preview/sites/{siteId} [get]
func (c *RenderController) PreviewSite(ctx *gin.Context) {
	siteID := ctx.Param("siteId")
	deviceType := ctx.Query("deviceType")
	if deviceType == "" {
		deviceType = "desktop"
	}

	request := dto.PreviewRequestDTO{
		SiteID:     siteID,
		DeviceType: deviceType,
	}

	result, err := c.renderAppService.PreviewPage(ctx, request)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "站点不存在" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.RenderErrorDTO{
			Code:    status,
			Message: "预览失败",
			Details: err.Error(),
		})
		return
	}

	ctx.Header("Content-Type", result.ContentType)
	ctx.String(http.StatusOK, result.Content)
}

// PreviewPage 预览页面
// @Summary 预览页面
// @Description 预览特定站点的特定页面
// @Tags 预览
// @Accept json
// @Produce html
// @Param siteId path string true "站点ID"
// @Param pageId path string true "页面ID"
// @Success 200 {string} string "渲染的HTML内容"
// @Failure 400 {object} dto.RenderErrorDTO
// @Failure 404 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /api/v1/preview/sites/{siteId}/pages/{pageId} [get]
func (c *RenderController) PreviewPage(ctx *gin.Context) {
	siteID := ctx.Param("siteId")
	pageID := ctx.Param("pageId")
	deviceType := ctx.Query("deviceType")
	if deviceType == "" {
		deviceType = "desktop"
	}

	request := dto.PreviewRequestDTO{
		SiteID:     siteID,
		PageID:     pageID,
		DeviceType: deviceType,
	}

	result, err := c.renderAppService.PreviewPage(ctx, request)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "页面不存在" || err.Error() == "站点不存在" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.RenderErrorDTO{
			Code:    status,
			Message: "预览失败",
			Details: err.Error(),
		})
		return
	}

	ctx.Header("Content-Type", result.ContentType)
	ctx.String(http.StatusOK, result.Content)
}

// 以下是模板管理相关的API

// CreateTemplate 创建模板
// @Summary 创建模板
// @Description 创建一个新的模板
// @Tags 模板
// @Accept json
// @Produce json
// @Param request body dto.CreateTemplateDTO true "创建模板请求"
// @Success 201 {object} dto.TemplateDTO
// @Failure 400 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /api/v1/templates [post]
func (c *RenderController) CreateTemplate(ctx *gin.Context) {
	var request dto.CreateTemplateDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.RenderErrorDTO{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
			Details: err.Error(),
		})
		return
	}

	template, err := c.renderAppService.CreateTemplate(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.RenderErrorDTO{
			Code:    http.StatusInternalServerError,
			Message: "创建模板失败",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, template)
}

// UpdateTemplate 更新模板
// @Summary 更新模板
// @Description 更新现有模板
// @Tags 模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Param request body dto.UpdateTemplateDTO true "更新模板请求"
// @Success 200 {object} dto.TemplateDTO
// @Failure 400 {object} dto.RenderErrorDTO
// @Failure 404 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /api/v1/templates/{id} [put]
func (c *RenderController) UpdateTemplate(ctx *gin.Context) {
	templateID := ctx.Param("id")
	var request dto.UpdateTemplateDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.RenderErrorDTO{
			Code:    http.StatusBadRequest,
			Message: "无效的请求参数",
			Details: err.Error(),
		})
		return
	}

	template, err := c.renderAppService.UpdateTemplate(ctx, templateID, request)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "模板不存在" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.RenderErrorDTO{
			Code:    status,
			Message: "更新模板失败",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, template)
}

// GetTemplate 获取模板
// @Summary 获取模板
// @Description 获取特定模板的详细信息
// @Tags 模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 200 {object} dto.TemplateDTO
// @Failure 404 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /api/v1/templates/{id} [get]
func (c *RenderController) GetTemplate(ctx *gin.Context) {
	templateID := ctx.Param("id")

	template, err := c.renderAppService.GetTemplate(ctx, templateID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "模板不存在" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.RenderErrorDTO{
			Code:    status,
			Message: "获取模板失败",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, template)
}

// DeleteTemplate 删除模板
// @Summary 删除模板
// @Description 删除特定模板
// @Tags 模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 204 "无内容"
// @Failure 404 {object} dto.RenderErrorDTO
// @Failure 500 {object} dto.RenderErrorDTO
// @Router /api/v1/templates/{id} [delete]
func (c *RenderController) DeleteTemplate(ctx *gin.Context) {
	templateID := ctx.Param("id")

	err := c.renderAppService.DeleteTemplate(ctx, templateID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "模板不存在" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.RenderErrorDTO{
			Code:    status,
			Message: "删除模板失败",
			Details: err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
