package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"

	"wz-backend-go/internal/application/admin/dto"
	"wz-backend-go/internal/application/admin/service"
	"wz-backend-go/internal/delivery/http/internal/types"
)

type LoginHandler struct {
	adminService *service.AdminApplicationService
}

func NewLoginHandler(adminService *service.AdminApplicationService) *LoginHandler {
	return &LoginHandler{
		adminService: adminService,
	}
}

// Login 管理员登录
func (h *LoginHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用服务进行登录
	loginReq := dto.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
		Captcha:  req.Captcha,
		// CaptchaId 可能需要从上下文中获取
		CaptchaId: req.CaptchaId,
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	resp, err := h.adminService.Login(c.Request.Context(), loginReq, ip, userAgent)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusUnauthorized, "登录失败: "+err.Error())
		return
	}

	// 构造响应
	response := types.LoginResponse{
		Code:    http.StatusOK,
		Message: "登录成功",
		Data: types.LoginResponseData{
			Token:     resp.Token,
			ExpiresAt: resp.ExpiresAt,
			User: types.AdminInfo{
				Id:        resp.Admin.ID,
				Username:  resp.Admin.Username,
				Role:      resp.Admin.RoleID,
				RoleName:  "",  // 可能需要从角色信息中获取
				Avatar:    resp.Admin.Avatar,
				Status:    int(resp.Admin.Status),
				LastLogin: resp.Admin.LastLoginAt.Format(time.RFC3339),
				CreatedAt: resp.Admin.CreatedAt.Format(time.RFC3339),
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken 刷新token
func (h *LoginHandler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.ErrorCtx(c, http.StatusBadRequest, err.Error())
		return
	}

	// 这里需要调用 JWTService 进行 token 刷新
	// 假设我们有一个方法可以获取 JWTService
	jwtService := GetJWTService() // 需要实现此方法
	
	token, expiresAt, err := jwtService.RefreshToken(c.Request.Context(), req.Token)
	if err != nil {
		httpx.ErrorCtx(c, http.StatusUnauthorized, "刷新令牌失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, types.RefreshTokenResponse{
		Code:    http.StatusOK,
		Message: "刷新令牌成功",
		Data: types.RefreshTokenResponseData{
			Token:     token,
			ExpiresAt: expiresAt,
		},
	})
}

// Captcha 获取验证码
func (h *LoginHandler) Captcha(c *gin.Context) {
	// 获取验证码逻辑
	// 这里需要一个生成验证码的服务
	captchaService := GetCaptchaService() // 需要实现此方法
	
	captchaId, picPath := captchaService.GenerateCaptcha()

	c.JSON(http.StatusOK, types.CaptchaResponse{
		Code:    http.StatusOK,
		Message: "获取验证码成功",
		Data: types.CaptchaResponseData{
			CaptchaId: captchaId,
			PicPath:   picPath,
		},
	})
}

// 临时帮助函数，实际实现需要根据项目的依赖注入方式
func GetJWTService() *service.JWTService {
	// 这只是一个示例，实际实现应根据项目的依赖注入机制
	panic("需要实现 GetJWTService")
}

func GetCaptchaService() interface{} {
	// 这只是一个示例，实际实现应根据项目的依赖注入机制
	panic("需要实现 GetCaptchaService")
}
