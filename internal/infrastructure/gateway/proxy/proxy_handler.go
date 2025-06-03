package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"wz-backend-go/internal/domain/gateway/entity"
)

// ProxyHandler 代理处理器实现
type ProxyHandler struct {
	httpClient *http.Client
}

// NewProxyHandler 创建新的代理处理器
func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			// 不要跟随重定向，我们需要将重定向响应直接返回给客户端
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// ProxyRequest 代理请求到目标服务
func (p *ProxyHandler) ProxyRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, route *entity.Route, service *entity.Service) error {
	// 记录请求开始时间
	startTime := time.Now()
	
	// 构建目标URL
	targetURL, err := p.buildTargetURL(r, route, service)
	if err != nil {
		http.Error(w, "无法构建目标URL", http.StatusInternalServerError)
		return fmt.Errorf("构建目标URL失败: %w", err)
	}
	
	// 创建代理请求
	proxyReq, err := p.createProxyRequest(ctx, r, targetURL)
	if err != nil {
		http.Error(w, "创建代理请求失败", http.StatusInternalServerError)
		return fmt.Errorf("创建代理请求失败: %w", err)
	}
	
	// 转发请求头
	p.copyRequestHeaders(r, proxyReq, service)
	
	// 添加代理信息头
	p.addProxyHeaders(proxyReq, r)
	
	// 发送代理请求
	resp, err := p.httpClient.Do(proxyReq)
	if err != nil {
		log.Printf("代理请求失败: %v", err)
		http.Error(w, "服务暂时不可用", http.StatusBadGateway)
		return fmt.Errorf("发送代理请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 复制响应头
	p.copyResponseHeaders(resp, w)
	
	// 设置代理响应头
	p.addProxyResponseHeaders(w, route, service, startTime)
	
	// 写入状态码
	w.WriteHeader(resp.StatusCode)
	
	// 复制响应体
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("复制响应体失败: %v", err)
		return fmt.Errorf("复制响应体失败: %w", err)
	}
	
	// 记录访问日志
	duration := time.Since(startTime)
	log.Printf("[%s] %s %s -> %s %d (%s)", r.Method, r.URL.Path, r.RemoteAddr, targetURL, resp.StatusCode, duration)
	
	return nil
}

// 构建目标URL
func (p *ProxyHandler) buildTargetURL(r *http.Request, route *entity.Route, service *entity.Service) (string, error) {
	// 从目标服务基础URL开始
	baseURL := service.BaseURL
	
	// 删除末尾斜杠
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	
	// 确定目标路径
	var targetPath string
	if route.StripPrefix && strings.HasPrefix(r.URL.Path, route.Path) {
		// 如果设置了StripPrefix，则从路径中移除前缀
		targetPath = strings.TrimPrefix(r.URL.Path, route.Path)
		if !strings.HasPrefix(targetPath, "/") {
			targetPath = "/" + targetPath
		}
	} else {
		// 否则使用完整路径
		targetPath = r.URL.Path
	}
	
	// 如果指定了重写路径，应用重写规则
	if route.RewritePath != "" {
		// 简单替换 - 在实际实现中，可能需要支持正则表达式重写
		targetPath = route.RewritePath
	}
	
	// 组合URL，包括查询字符串
	targetURL := fmt.Sprintf("%s%s", baseURL, targetPath)
	if r.URL.RawQuery != "" {
		targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
	}
	
	// 验证URL是否有效
	_, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("无效的目标URL: %w", err)
	}
	
	return targetURL, nil
}

// 创建代理请求
func (p *ProxyHandler) createProxyRequest(ctx context.Context, originalReq *http.Request, targetURL string) (*http.Request, error) {
	// 创建新的请求
	var body io.Reader
	if originalReq.Body != nil {
		// 读取原始请求体
		bodyBytes, err := io.ReadAll(originalReq.Body)
		if err != nil {
			return nil, fmt.Errorf("读取请求体失败: %w", err)
		}
		originalReq.Body.Close()
		
		// 创建新的请求体
		body = bytes.NewBuffer(bodyBytes)
		
		// 为原始请求重新创建Body以便后续处理
		originalReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	
	// 创建代理请求
	proxyReq, err := http.NewRequestWithContext(ctx, originalReq.Method, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建代理请求失败: %w", err)
	}
	
	return proxyReq, nil
}

// 复制请求头
func (p *ProxyHandler) copyRequestHeaders(originalReq *http.Request, proxyReq *http.Request, service *entity.Service) {
	// 跳过特定头部（将由代理处理）
	skipHeaders := map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"TE":                  true,
		"Trailers":            true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}
	
	// 复制所有其他头部
	for name, values := range originalReq.Header {
		if !skipHeaders[name] {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}
	}
	
	// 处理Host头部
	if service.PreserveHost && originalReq.Host != "" {
		proxyReq.Host = originalReq.Host
	}
}

// 添加代理信息头
func (p *ProxyHandler) addProxyHeaders(proxyReq *http.Request, originalReq *http.Request) {
	// 添加代理信息
	proxyReq.Header.Set("X-Forwarded-For", getClientIP(originalReq))
	proxyReq.Header.Set("X-Forwarded-Proto", getScheme(originalReq))
	proxyReq.Header.Set("X-Forwarded-Host", originalReq.Host)
	
	// 添加原始URL
	proxyReq.Header.Set("X-Original-URL", originalReq.URL.String())
	
	// 添加代理标识
	proxyReq.Header.Set("X-Proxy-ID", "wanzhi-gateway")
}

// 复制响应头
func (p *ProxyHandler) copyResponseHeaders(resp *http.Response, w http.ResponseWriter) {
	// 跳过特定头部
	skipHeaders := map[string]bool{
		"Connection":        true,
		"Keep-Alive":        true,
		"Transfer-Encoding": true,
	}
	
	// 复制所有其他头部
	for name, values := range resp.Header {
		if !skipHeaders[name] {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
	}
}

// 添加代理响应头
func (p *ProxyHandler) addProxyResponseHeaders(w http.ResponseWriter, route *entity.Route, service *entity.Service, startTime time.Time) {
	// 添加代理信息
	w.Header().Set("X-Proxy-By", "wanzhi-gateway")
	w.Header().Set("X-Proxy-Service", service.Name)
	
	// 添加路由信息
	w.Header().Set("X-Route-ID", route.ID)
	
	// 添加性能信息
	duration := time.Since(startTime)
	w.Header().Set("X-Response-Time", fmt.Sprintf("%d", duration.Milliseconds()))
}

// 获取客户端IP
func getClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// 获取链中的第一个IP
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// 检查X-Real-IP头
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	
	// 从RemoteAddr获取IP
	ip := r.RemoteAddr
	// 移除端口部分
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	
	return ip
}

// 获取请求协议
func getScheme(r *http.Request) string {
	// 检查X-Forwarded-Proto头
	forwardedProto := r.Header.Get("X-Forwarded-Proto")
	if forwardedProto != "" {
		return forwardedProto
	}
	
	// 如果是HTTPS请求
	if r.TLS != nil {
		return "https"
	}
	
	// 默认为HTTP
	return "http"
}
