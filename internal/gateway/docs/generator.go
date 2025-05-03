package docs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"wz-backend-go/internal/gateway/config"
)

// APIDocGenerator API文档生成器
type APIDocGenerator struct {
	swagger    *openapi3.T
	config     config.Config
	outputPath string
	mu         sync.RWMutex
}

// NewAPIDocGenerator 创建新的API文档生成器
func NewAPIDocGenerator(conf config.Config, outputPath string) *APIDocGenerator {
	// 创建基础Swagger文档结构
	swagger := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       "微众后端网关API",
			Description: "微众后端服务API文档，自动从网关路由配置生成",
			Version:     "1.0.0",
			Contact: &openapi3.Contact{
				Name:  "微众技术团队",
				Email: "tech@example.com",
				URL:   "https://example.com",
			},
		},
		Paths: openapi3.Paths{},
		Components: &openapi3.Components{
			Schemas:         make(map[string]*openapi3.SchemaRef),
			SecuritySchemes: make(map[string]*openapi3.SecuritySchemeRef),
		},
	}

	// 添加认证方案
	if conf.Security.JwtSecret != "" {
		swagger.Components.SecuritySchemes["bearerAuth"] = &openapi3.SecuritySchemeRef{
			Value: &openapi3.SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "使用JWT令牌进行认证，通过'Authorization: Bearer {token}'头传递",
			},
		}
	}

	if conf.Security.APIKey.Enabled {
		headerName := conf.Security.APIKey.HeaderName
		if headerName == "" {
			headerName = "X-API-Key"
		}
		swagger.Components.SecuritySchemes["apiKeyAuth"] = &openapi3.SecuritySchemeRef{
			Value: &openapi3.SecurityScheme{
				Type: "apiKey",
				In:   "header",
				Name: headerName,
				Description: "使用API Key进行认证，通过'" + headerName + ": {key}'头传递，" +
					"或者通过查询参数'" + conf.Security.APIKey.QueryParamName + "={key}'传递",
			},
		}
	}

	if conf.Security.OAuth2.Enabled {
		swagger.Components.SecuritySchemes["oauth2"] = &openapi3.SecuritySchemeRef{
			Value: &openapi3.SecurityScheme{
				Type: "oauth2",
				Flows: &openapi3.OAuthFlows{
					AuthorizationCode: &openapi3.OAuthFlow{
						AuthorizationURL: conf.Security.OAuth2.AuthorizationURL,
						TokenURL:         conf.Security.OAuth2.TokenURL,
						Scopes:           parseScopes(conf.Security.OAuth2.Scope),
					},
				},
				Description: "使用OAuth2.0进行认证",
			},
		}
	}

	return &APIDocGenerator{
		swagger:    swagger,
		config:     conf,
		outputPath: outputPath,
	}
}

// parseScopes 解析OAuth2作用域字符串
func parseScopes(scopeStr string) map[string]string {
	scopes := make(map[string]string)
	for _, scope := range strings.Split(scopeStr, " ") {
		scope = strings.TrimSpace(scope)
		if scope != "" {
			scopes[scope] = fmt.Sprintf("'%s'权限", scope)
		}
	}
	return scopes
}

// GenerateFromServices 从服务配置生成API文档
func (g *APIDocGenerator) GenerateFromServices() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 清除现有路径
	g.swagger.Paths = openapi3.Paths{}

	// 注册所有服务的API
	for _, service := range g.config.Services {
		// 处理HTTP路由
		if service.Type == "http" {
			for _, route := range service.Routes {
				g.addHTTPRoute(service, route)
			}
		}

		// 处理gRPC路由
		if service.Type == "grpc" && len(service.GrpcOptions.Methods) > 0 {
			for _, method := range service.GrpcOptions.Methods {
				g.addGRPCMethod(service, method)
			}
		}
	}

	// 保存Swagger文档
	return g.saveSwagger()
}

// addHTTPRoute 添加HTTP路由到文档
func (g *APIDocGenerator) addHTTPRoute(service config.ServiceConfig, route config.RouteConfig) {
	// 构建完整路径
	fullPath := service.Prefix + route.Path

	// 构建路径项
	pathItem := &openapi3.PathItem{
		Summary:     fmt.Sprintf("%s服务路由", service.Name),
		Description: fmt.Sprintf("服务: %s, 目标: %s", service.Name, route.Target),
	}

	// 创建操作
	operation := &openapi3.Operation{
		Tags:        []string{service.Name},
		Summary:     fmt.Sprintf("%s服务 - %s", service.Name, route.Path),
		Description: fmt.Sprintf("服务: %s, 路径: %s, 目标: %s", service.Name, route.Path, route.Target),
		Responses:   openapi3.NewResponses(),
	}

	// 添加请求和响应示例
	operation.Responses["200"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: openapi3.Stringf("成功响应"),
			Content: openapi3.NewContentWithJSONSchema(
				&openapi3.Schema{
					Type:        "object",
					Description: "响应数据",
				},
			),
		},
	}

	// 添加认证需求
	if service.Authentication {
		operation.Security = []openapi3.SecurityRequirement{}
		
		for _, method := range g.config.Security.AuthMethods {
			switch method {
			case "jwt":
				operation.Security = append(operation.Security, 
					openapi3.SecurityRequirement{"bearerAuth": []string{}})
			case "apikey":
				operation.Security = append(operation.Security, 
					openapi3.SecurityRequirement{"apiKeyAuth": []string{}})
			case "oauth2":
				scopes := strings.Split(g.config.Security.OAuth2.Scope, " ")
				operation.Security = append(operation.Security, 
					openapi3.SecurityRequirement{"oauth2": scopes})
			}
		}
	}

	// 根据HTTP方法设置操作
	method := strings.ToUpper(route.Method)
	switch method {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	case "PATCH":
		pathItem.Patch = operation
	case "HEAD":
		pathItem.Head = operation
	case "OPTIONS":
		pathItem.Options = operation
	case "TRACE":
		pathItem.Trace = operation
	}

	// 添加到路径
	g.swagger.Paths[fullPath] = pathItem
}

// addGRPCMethod 添加gRPC方法到文档
func (g *APIDocGenerator) addGRPCMethod(service config.ServiceConfig, method config.GrpcMethodConfig) {
	// 构建完整路径
	fullPath := service.Prefix + method.Path

	// 构建路径项
	pathItem := &openapi3.PathItem{
		Summary:     fmt.Sprintf("%s gRPC服务方法", service.Name),
		Description: fmt.Sprintf("服务: %s, gRPC方法: %s", service.Name, method.Name),
	}

	// 创建操作
	operation := &openapi3.Operation{
		Tags:        []string{service.Name + " (gRPC)"},
		Summary:     fmt.Sprintf("%s服务 - %s方法", service.Name, method.Name),
		Description: fmt.Sprintf("gRPC服务: %s, 方法: %s", service.Name, method.Name),
		Responses:   openapi3.NewResponses(),
	}

	// 添加请求体
	requestBody := &openapi3.RequestBody{
		Description: "gRPC请求参数",
		Content: openapi3.NewContentWithJSONSchema(
			&openapi3.Schema{
				Type:        "object",
				Description: "请求参数JSON",
			},
		),
		Required: true,
	}
	operation.RequestBody = &openapi3.RequestBodyRef{Value: requestBody}

	// 添加响应
	operation.Responses["200"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: openapi3.Stringf("成功响应"),
			Content: openapi3.NewContentWithJSONSchema(
				&openapi3.Schema{
					Type:        "object",
					Description: "响应数据JSON",
				},
			),
		},
	}

	// 添加常见错误响应
	operation.Responses["400"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: openapi3.Stringf("请求参数错误"),
		},
	}
	operation.Responses["500"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: openapi3.Stringf("服务器内部错误"),
		},
	}

	// 添加认证需求
	if service.Authentication {
		operation.Security = []openapi3.SecurityRequirement{}
		
		for _, authMethod := range g.config.Security.AuthMethods {
			switch authMethod {
			case "jwt":
				operation.Security = append(operation.Security, 
					openapi3.SecurityRequirement{"bearerAuth": []string{}})
			case "apikey":
				operation.Security = append(operation.Security, 
					openapi3.SecurityRequirement{"apiKeyAuth": []string{}})
			case "oauth2":
				scopes := strings.Split(g.config.Security.OAuth2.Scope, " ")
				operation.Security = append(operation.Security, 
					openapi3.SecurityRequirement{"oauth2": scopes})
			}
		}
	}

	// 根据HTTP方法设置操作
	httpMethod := strings.ToUpper(method.HTTPMethod)
	switch httpMethod {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	case "PATCH":
		pathItem.Patch = operation
	default:
		// 默认为POST
		pathItem.Post = operation
	}

	// 添加到路径
	g.swagger.Paths[fullPath] = pathItem
}

// saveSwagger 保存Swagger文档
func (g *APIDocGenerator) saveSwagger() error {
	// 创建输出目录
	if err := ensureDir(g.outputPath); err != nil {
		return err
	}

	// 保存为JSON文件
	jsonPath := filepath.Join(g.outputPath, "swagger.json")
	jsonData, err := json.MarshalIndent(g.swagger, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化Swagger文档失败: %w", err)
	}

	if err := ioutil.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("保存Swagger JSON文件失败: %w", err)
	}

	return nil
}

// GetSwaggerJSON 获取Swagger JSON
func (g *APIDocGenerator) GetSwaggerJSON() ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	return json.MarshalIndent(g.swagger, "", "  ")
}

// ensureDir 确保目录存在
func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
