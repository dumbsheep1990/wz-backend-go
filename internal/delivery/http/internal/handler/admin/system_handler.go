package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/types"
)

// SystemHandler 系统管理处理程序
type SystemHandler struct {
	systemService *service.SystemApplicationService
}

// NewSystemHandler 创建系统管理处理程序
func NewSystemHandler(systemService *service.SystemApplicationService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// GetSystemConfig 获取系统配置
func (h *SystemHandler) GetSystemConfig(c *gin.Context) {
	// 获取系统配置
	config, err := h.systemService.GetSystemConfig(c.Request.Context())
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取系统配置失败: "+err.Error())
		return
	}

	// 转换为API响应格式
	response := types.GetSystemConfigResponse{
		Code:    http.StatusOK,
		Message: "获取系统配置成功",
		Data: types.SystemConfig{
			Id:                config.ID,
			CreatedAt:         config.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:         config.UpdatedAt.Format(http.TimeFormat),
			SystemConfig:      config.SystemConfig,
			LogoUrl:           config.LogoUrl,
			ApiUrl:            config.ApiUrl,
			Name:              config.Name,
			Description:       config.Description,
			ContentUrl:        config.ContentUrl,
			UploadFileSize:    config.UploadFileSize,
			EmailFrom:         config.EmailFrom,
			EmailHost:         config.EmailHost,
			EmailPort:         config.EmailPort,
			EmailSecret:       config.EmailSecret,
			EmailIsSSL:        config.EmailIsSSL,
			EmailNickname:     config.EmailNickname,
			LdapOpen:          config.LdapOpen,
			LdapHost:          config.LdapHost,
			LdapPort:          config.LdapPort,
			LdapBindDn:        config.LdapBindDn,
			LdapPassword:      config.LdapPassword,
			LdapBaseDn:        config.LdapBaseDn,
			LdapUserField:     config.LdapUserField,
			LdapNickName:      config.LdapNickName,
			LdapEmails:        config.LdapEmails,
			LdapTls:           config.LdapTls,
			DingTalkOpen:      config.DingTalkOpen,
			DingTalkAppKey:    config.DingTalkAppKey,
			DingTalkAppSecret: config.DingTalkAppSecret,
			DingTalkAgentId:   config.DingTalkAgentId,
		},
	}

	c.JSON(http.StatusOK, response)
}

// SetSystemConfig 设置系统配置
func (h *SystemHandler) SetSystemConfig(c *gin.Context) {
	var req types.SystemConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转换为应用服务请求格式
	configReq := h.systemService.ConvertToSystemConfigRequest(req)

	// 更新系统配置
	err := h.systemService.SetSystemConfig(c.Request.Context(), configReq)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "设置系统配置失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "设置系统配置成功",
	})
}

// GetServerInfo 获取服务器信息
func (h *SystemHandler) GetServerInfo(c *gin.Context) {
	// 获取服务器信息
	info, err := h.systemService.GetServerInfo(c.Request.Context())
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "获取服务器信息失败: "+err.Error())
		return
	}

	// 转换为API响应格式
	response := types.GetServerInfoResponse{
		Code:    http.StatusOK,
		Message: "获取服务器信息成功",
		Data: types.ServerInfo{
			Os: types.OsInfo{
				Name:          info.Os.Name,
				Platform:      info.Os.Platform,
				Family:        info.Os.Family,
				Version:       info.Os.Version,
				KernelArch:    info.Os.KernelArch,
				KernelVersion: info.Os.KernelVersion,
			},
			Cpu: types.CpuInfo{
				Model:       info.Cpu.Model,
				Cores:       info.Cpu.Cores,
				UsedPercent: info.Cpu.UsedPercent,
			},
			Ram: types.RamInfo{
				Total:       info.Ram.Total,
				Used:        info.Ram.Used,
				Free:        info.Ram.Free,
				UsedPercent: info.Ram.UsedPercent,
			},
			Disk: types.DiskInfo{
				Total:       info.Disk.Total,
				Used:        info.Disk.Used,
				Free:        info.Disk.Free,
				UsedPercent: info.Disk.UsedPercent,
			},
			GoInfo: types.GoInfo{
				Version:      info.Go.Version,
				NumGoroutine: info.Go.NumGoroutine,
				NumCpu:       info.Go.NumCpu,
			},
			DbInfo: types.DbInfo{
				Type:        info.Db.Type,
				Version:     info.Db.Version,
				DbName:      info.Db.DbName,
				MaxOpenConn: info.Db.MaxOpenConn,
				MaxIdleConn: info.Db.MaxIdleConn,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// ReloadSystem 重启系统
func (h *SystemHandler) ReloadSystem(c *gin.Context) {
	// 调用重启系统服务
	err := h.systemService.ReloadSystem(c.Request.Context())
	if err != nil {
		httpx.ErrorCtx(c, http.StatusInternalServerError, "重启系统失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Code:    http.StatusOK,
		Message: "系统已重启",
	})
}
