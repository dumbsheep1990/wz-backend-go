package valueobject

import "errors"

// RenderFormat 表示渲染输出的格式
type RenderFormat struct {
	format      string
	contentType string
	extension   string
}

// 预定义的渲染格式
var (
	// FormatHTML HTML渲染格式
	FormatHTML = RenderFormat{
		format:      "html",
		contentType: "text/html; charset=utf-8",
		extension:   ".html",
	}

	// FormatJSON JSON渲染格式
	FormatJSON = RenderFormat{
		format:      "json",
		contentType: "application/json; charset=utf-8",
		extension:   ".json",
	}

	// FormatAMP AMP渲染格式
	FormatAMP = RenderFormat{
		format:      "amp",
		contentType: "text/html; charset=utf-8",
		extension:   ".amp.html",
	}

	// FormatXML XML渲染格式
	FormatXML = RenderFormat{
		format:      "xml",
		contentType: "application/xml; charset=utf-8",
		extension:   ".xml",
	}

	// FormatText 纯文本渲染格式
	FormatText = RenderFormat{
		format:      "text",
		contentType: "text/plain; charset=utf-8",
		extension:   ".txt",
	}
)

// NewRenderFormat 创建一个新的渲染格式
func NewRenderFormat(format, contentType, extension string) (RenderFormat, error) {
	if format == "" {
		return RenderFormat{}, errors.New("渲染格式不能为空")
	}

	if contentType == "" {
		return RenderFormat{}, errors.New("内容类型不能为空")
	}

	return RenderFormat{
		format:      format,
		contentType: contentType,
		extension:   extension,
	}, nil
}

// Format 返回渲染格式
func (rf RenderFormat) Format() string {
	return rf.format
}

// ContentType 返回内容类型
func (rf RenderFormat) ContentType() string {
	return rf.contentType
}

// Extension 返回文件扩展名
func (rf RenderFormat) Extension() string {
	return rf.extension
}

// Equals 比较两个渲染格式是否相等
func (rf RenderFormat) Equals(other RenderFormat) bool {
	return rf.format == other.format &&
		rf.contentType == other.contentType &&
		rf.extension == other.extension
}

// FromString 从字符串创建渲染格式
func FormatFromString(format string) (RenderFormat, error) {
	switch format {
	case "html":
		return FormatHTML, nil
	case "json":
		return FormatJSON, nil
	case "amp":
		return FormatAMP, nil
	case "xml":
		return FormatXML, nil
	case "text":
		return FormatText, nil
	default:
		return RenderFormat{}, errors.New("不支持的渲染格式: " + format)
	}
}
