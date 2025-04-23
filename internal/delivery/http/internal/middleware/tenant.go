package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"wz-backend-go/internal/domain/model"
)

// TenantMiddleware 租户中间件，从子域名解析租户ID
func TenantMiddleware(mainDomain string, tenantService TenantService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			// 检查请求头是否包含X-Tenant-ID
			tenantIDHeader := r.Header.Get("X-Tenant-ID")
			if tenantIDHeader != "" {
				// 如果请求头中包含租户ID，则直接使用
				tenant, err := tenantService.GetTenantByIDString(r.Context(), tenantIDHeader)
				if err == nil && tenant != nil {
					ctx := context.WithValue(r.Context(), CtxTenantIDKey, tenant.ID)
					// 在这个分支中，我们从头部获取了租户信息，继续请求
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				// 如果租户ID无效，继续尝试从子域名解析
				logx.Infof("无效的租户ID头部: %s, 尝试从子域名解析", tenantIDHeader)
			}

			// 检查是否是子域名
			if host != mainDomain && strings.Contains(host, ".") {
				// 提取子域名部分
				subdomain := strings.Split(host, ".")[0]
				
				// 通过子域名查找租户
				tenant, err := tenantService.GetTenantBySubdomain(r.Context(), subdomain)
				if err == nil && tenant != nil {
					// 如果找到租户，将租户ID设置到上下文中
					ctx := context.WithValue(r.Context(), CtxTenantIDKey, tenant.ID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				
				// 无效的子域名，可能返回404或重定向到主域名
				logx.Infof("无效的子域名: %s", subdomain)
				http.Error(w, "未找到租户站点", http.StatusNotFound)
				return
			}
			
			// 如果是主域名，继续处理请求
			next.ServeHTTP(w, r)
		})
	}
}

// TenantService 租户服务接口
type TenantService interface {
	GetTenantBySubdomain(ctx context.Context, subdomain string) (*model.Tenant, error)
	GetTenantByIDString(ctx context.Context, id string) (*model.Tenant, error)
}

// ValidateTenantAccess 验证租户访问权限
func ValidateTenantAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从上下文中获取当前请求的租户ID
		requestTenantID, hasTenant := GetTenantIDFromContext(r.Context())
		
		// 从JWT令牌中获取用户关联的租户ID
		tokenTenantID, hasTokenTenant := GetTenantIDFromContext(r.Context())
		
		// 如果请求中有租户ID，同时用户令牌中也有租户ID，则需要验证是否一致
		if hasTenant && hasTokenTenant && requestTenantID != tokenTenantID {
			// 用户角色
			role, hasRole := GetUserRoleFromContext(r.Context())
			
			// 如果是平台管理员，则允许访问任何租户
			if hasRole && role == model.RolePlatformAdmin {
				next.ServeHTTP(w, r)
				return
			}
			
			// 否则，仅允许访问用户所属租户
			http.Error(w, "无权访问此租户的资源", http.StatusForbidden)
			return
		}
		
		// 允许请求继续
		next.ServeHTTP(w, r)
	})
}
