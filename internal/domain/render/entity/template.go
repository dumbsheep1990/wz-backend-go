package entity

import (
	"errors"
	"time"
	"wz-backend-go/internal/domain/render/valueobject"
)

// Template 表示渲染模板
type Template struct {
	id          string
	name        string
	description string
	content     string
	type_       string
	version     string
	createdAt   time.Time
	updatedAt   time.Time
	metadata    map[string]string
	siteID      string
}

// 模板类型常量
const (
	TemplateTypePage     = "page"
	TemplateTypeSection  = "section"
	TemplateTypeEmail    = "email"
	TemplateTypePartial  = "partial"
	TemplateTypeLayout   = "layout"
	TemplateTypeComponent = "component"
)

// NewTemplate 创建一个新的模板
func NewTemplate(
	id string,
	name string,
	description string,
	content string,
	type_ string,
	version string,
	siteID string,
	metadata map[string]string,
) (*Template, error) {
	if name == "" {
		return nil, errors.New("模板名称不能为空")
	}

	if content == "" {
		return nil, errors.New("模板内容不能为空")
	}

	if type_ == "" {
		return nil, errors.New("模板类型不能为空")
	}

	now := time.Now()

	if metadata == nil {
		metadata = make(map[string]string)
	}

	return &Template{
		id:          id,
		name:        name,
		description: description,
		content:     content,
		type_:       type_,
		version:     version,
		createdAt:   now,
		updatedAt:   now,
		metadata:    metadata,
		siteID:      siteID,
	}, nil
}

// ID 返回模板ID
func (t *Template) ID() string {
	return t.id
}

// Name 返回模板名称
func (t *Template) Name() string {
	return t.name
}

// Description 返回模板描述
func (t *Template) Description() string {
	return t.description
}

// Content 返回模板内容
func (t *Template) Content() string {
	return t.content
}

// Type 返回模板类型
func (t *Template) Type() string {
	return t.type_
}

// Version 返回模板版本
func (t *Template) Version() string {
	return t.version
}

// CreatedAt 返回创建时间
func (t *Template) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt 返回更新时间
func (t *Template) UpdatedAt() time.Time {
	return t.updatedAt
}

// Metadata 返回元数据
func (t *Template) Metadata() map[string]string {
	return t.metadata
}

// SiteID 返回站点ID
func (t *Template) SiteID() string {
	return t.siteID
}

// SetName 设置模板名称
func (t *Template) SetName(name string) error {
	if name == "" {
		return errors.New("模板名称不能为空")
	}
	t.name = name
	t.updatedAt = time.Now()
	return nil
}

// SetDescription 设置模板描述
func (t *Template) SetDescription(description string) {
	t.description = description
	t.updatedAt = time.Now()
}

// SetContent 设置模板内容
func (t *Template) SetContent(content string) error {
	if content == "" {
		return errors.New("模板内容不能为空")
	}
	t.content = content
	t.updatedAt = time.Now()
	return nil
}

// SetType 设置模板类型
func (t *Template) SetType(type_ string) error {
	if type_ == "" {
		return errors.New("模板类型不能为空")
	}
	t.type_ = type_
	t.updatedAt = time.Now()
	return nil
}

// SetVersion 设置模板版本
func (t *Template) SetVersion(version string) {
	t.version = version
	t.updatedAt = time.Now()
}

// SetMetadata 设置元数据
func (t *Template) SetMetadata(key string, value string) {
	t.metadata[key] = value
	t.updatedAt = time.Now()
}

// Render 渲染模板
func (t *Template) Render(context valueobject.TemplateContext, format valueobject.RenderFormat) (*RenderResult, error) {
	// 在真实实现中，这里将使用模板引擎渲染内容
	// 这里是一个简化的示例
	
	// 创建渲染结果
	renderID := valueobject.NewRenderID()
	
	// 创建缓存策略（示例使用一小时缓存）
	cacheKey := t.id + ":" + format.Format()
	cacheStrategy, _ := valueobject.NewCacheStrategy(
		true,
		time.Hour,
		valueobject.CacheLevelMemory,
		cacheKey,
		[]string{t.siteID, t.type_},
	)
	
	// 这里是一个模拟的渲染结果
	// 实际实现中，会基于模板内容和上下文数据进行渲染
	content := "渲染的" + t.name + "模板"
	
	return NewRenderResult(
		renderID,
		content,
		format,
		cacheStrategy,
		context,
	), nil
}
