package valueobject

import (
	"encoding/json"
	"errors"
)

// TemplateContext 表示模板渲染的上下文
type TemplateContext struct {
	data     map[string]interface{}
	metadata map[string]string
}

// NewTemplateContext 创建一个新的模板上下文
func NewTemplateContext(data map[string]interface{}, metadata map[string]string) TemplateContext {
	if data == nil {
		data = make(map[string]interface{})
	}

	if metadata == nil {
		metadata = make(map[string]string)
	}

	return TemplateContext{
		data:     data,
		metadata: metadata,
	}
}

// Data 返回模板上下文数据
func (tc TemplateContext) Data() map[string]interface{} {
	return tc.data
}

// Metadata 返回模板上下文元数据
func (tc TemplateContext) Metadata() map[string]string {
	return tc.metadata
}

// Get 获取指定键的数据
func (tc TemplateContext) Get(key string) (interface{}, bool) {
	value, exists := tc.data[key]
	return value, exists
}

// GetMetadata 获取指定键的元数据
func (tc TemplateContext) GetMetadata(key string) (string, bool) {
	value, exists := tc.metadata[key]
	return value, exists
}

// SetData 设置指定键的数据
func (tc *TemplateContext) SetData(key string, value interface{}) {
	tc.data[key] = value
}

// SetMetadata 设置指定键的元数据
func (tc *TemplateContext) SetMetadata(key string, value string) {
	tc.metadata[key] = value
}

// Clone 克隆模板上下文
func (tc TemplateContext) Clone() TemplateContext {
	newData := make(map[string]interface{})
	for k, v := range tc.data {
		newData[k] = v
	}

	newMetadata := make(map[string]string)
	for k, v := range tc.metadata {
		newMetadata[k] = v
	}

	return TemplateContext{
		data:     newData,
		metadata: newMetadata,
	}
}

// Merge 合并另一个模板上下文
func (tc *TemplateContext) Merge(other TemplateContext) {
	for k, v := range other.data {
		tc.data[k] = v
	}

	for k, v := range other.metadata {
		tc.metadata[k] = v
	}
}

// ToJSON 将模板上下文转换为JSON字符串
func (tc TemplateContext) ToJSON() (string, error) {
	type contextJSON struct {
		Data     map[string]interface{} `json:"data"`
		Metadata map[string]string      `json:"metadata"`
	}

	ctx := contextJSON{
		Data:     tc.data,
		Metadata: tc.metadata,
	}

	bytes, err := json.Marshal(ctx)
	if err != nil {
		return "", errors.New("无法将模板上下文转换为JSON: " + err.Error())
	}

	return string(bytes), nil
}

// FromJSON 从JSON字符串创建模板上下文
func FromJSON(jsonStr string) (TemplateContext, error) {
	type contextJSON struct {
		Data     map[string]interface{} `json:"data"`
		Metadata map[string]string      `json:"metadata"`
	}

	var ctx contextJSON
	err := json.Unmarshal([]byte(jsonStr), &ctx)
	if err != nil {
		return TemplateContext{}, errors.New("无法从JSON解析模板上下文: " + err.Error())
	}

	return NewTemplateContext(ctx.Data, ctx.Metadata), nil
}
