package component

import (
	"time"

	"wz-backend-go/models/common"
)

// ComponentCategory 组件分类
type ComponentCategory struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string `json:"name" db:"name" gorm:"not null;comment:分类名称"`
	Code        string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:分类编码"`
	Icon        string `json:"icon" db:"icon" gorm:"comment:图标"`
	Description string `json:"description" db:"description" gorm:"comment:描述"`
	SortOrder   int    `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	ParentID    int64  `json:"parentId" db:"parent_id" gorm:"index;default:0;comment:父分类ID"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ComponentTemplate 组件模板
type ComponentTemplate struct {
	common.BaseIDModel
	common.BaseTimeModel
	CategoryID  int64  `json:"categoryId" db:"category_id" gorm:"index;not null;comment:分类ID"`
	Name        string `json:"name" db:"name" gorm:"not null;comment:组件名称"`
	Code        string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:组件编码"`
	Version     string `json:"version" db:"version" gorm:"not null;comment:版本号"`
	Icon        string `json:"icon" db:"icon" gorm:"comment:图标"`
	Screenshot  string `json:"screenshot" db:"screenshot" gorm:"comment:截图"`
	Description string `json:"description" db:"description" gorm:"comment:描述"`
	Author      string `json:"author" db:"author" gorm:"comment:作者"`
	Tags        string `json:"tags" db:"tags" gorm:"type:json;comment:标签,JSON数组"`
	Type        string `json:"type" db:"type" gorm:"comment:类型:basic-基础,business-业务,layout-布局,container-容器"`
	Framework   string `json:"framework" db:"framework" gorm:"comment:框架:react,vue,angular"`
	UseCount    int    `json:"useCount" db:"use_count" gorm:"default:0;comment:使用次数"`
	IsSystem    bool   `json:"isSystem" db:"is_system" gorm:"default:false;comment:是否系统组件"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ComponentResource 组件资源
type ComponentResource struct {
	common.BaseIDModel
	common.BaseTimeModel
	TemplateID   int64  `json:"templateId" db:"template_id" gorm:"index;not null;comment:组件模板ID"`
	ResourceType string `json:"resourceType" db:"resource_type" gorm:"not null;comment:资源类型:js,css,html,img,other"`
	ResourcePath string `json:"resourcePath" db:"resource_path" gorm:"not null;comment:资源路径"`
	IsMain       bool   `json:"isMain" db:"is_main" gorm:"default:false;comment:是否主资源"`
	Size         int64  `json:"size" db:"size" gorm:"comment:大小(字节)"`
	MD5          string `json:"md5" db:"md5" gorm:"comment:MD5哈希值"`
	Status       int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	TenantID     int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ComponentMetadata 组件元数据
type ComponentMetadata struct {
	common.BaseIDModel
	TemplateID      int64  `json:"templateId" db:"template_id" gorm:"uniqueIndex;not null;comment:组件模板ID"`
	SchemaVersion   string `json:"schemaVersion" db:"schema_version" gorm:"comment:Schema版本"`
	PropsSchema     string `json:"propsSchema" db:"props_schema" gorm:"type:json;comment:属性Schema,JSON格式"`
	EventsSchema    string `json:"eventsSchema" db:"events_schema" gorm:"type:json;comment:事件Schema,JSON格式"`
	SlotsSchema     string `json:"slotsSchema" db:"slots_schema" gorm:"type:json;comment:插槽Schema,JSON格式"`
	StylesSchema    string `json:"stylesSchema" db:"styles_schema" gorm:"type:json;comment:样式Schema,JSON格式"`
	DefaultProps    string `json:"defaultProps" db:"default_props" gorm:"type:json;comment:默认属性值,JSON格式"`
	DefaultStyles   string `json:"defaultStyles" db:"default_styles" gorm:"type:json;comment:默认样式值,JSON格式"`
	Dependencies    string `json:"dependencies" db:"dependencies" gorm:"type:json;comment:依赖项,JSON格式"`
	ConfigInterface string `json:"configInterface" db:"config_interface" gorm:"comment:配置界面路径"`
	TenantID        int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ComponentInstance 组件实例
type ComponentInstance struct {
	common.BaseIDModel
	common.BaseTimeModel
	TemplateID  int64  `json:"templateId" db:"template_id" gorm:"index;not null;comment:组件模板ID"`
	PageID      string `json:"pageId" db:"page_id" gorm:"index;comment:页面ID"`
	SectionID   string `json:"sectionId" db:"section_id" gorm:"index;comment:区块ID"`
	InstanceKey string `json:"instanceKey" db:"instance_key" gorm:"uniqueIndex;not null;comment:实例唯一标识"`
	Props       string `json:"props" db:"props" gorm:"type:json;comment:属性值,JSON格式"`
	Events      string `json:"events" db:"events" gorm:"type:json;comment:事件绑定,JSON格式"`
	Styles      string `json:"styles" db:"styles" gorm:"type:json;comment:样式值,JSON格式"`
	DataSource  string `json:"dataSource" db:"data_source" gorm:"type:json;comment:数据源配置,JSON格式"`
	Position    string `json:"position" db:"position" gorm:"type:json;comment:位置信息,JSON格式"`
	ParentKey   string `json:"parentKey" db:"parent_key" gorm:"index;comment:父组件实例Key"`
	SortOrder   int    `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ComponentDataSource 组件数据源
type ComponentDataSource struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name          string `json:"name" db:"name" gorm:"not null;comment:数据源名称"`
	Code          string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:数据源编码"`
	Type          string `json:"type" db:"type" gorm:"not null;comment:类型:api,database,static,function"`
	Configuration string `json:"configuration" db:"configuration" gorm:"type:json;not null;comment:配置,JSON格式"`
	ReturnSchema  string `json:"returnSchema" db:"return_schema" gorm:"type:json;comment:返回数据结构,JSON格式"`
	Description   string `json:"description" db:"description" gorm:"comment:描述"`
	IsSystem      bool   `json:"isSystem" db:"is_system" gorm:"default:false;comment:是否系统数据源"`
	Status        int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	CreatedBy     int64  `json:"createdBy" db:"created_by" gorm:"comment:创建人"`
	UpdatedBy     int64  `json:"updatedBy" db:"updated_by" gorm:"comment:更新人"`
	TenantID      int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ComponentEffect 组件特效
type ComponentEffect struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name        string `json:"name" db:"name" gorm:"not null;comment:特效名称"`
	Code        string `json:"code" db:"code" gorm:"uniqueIndex;not null;comment:特效编码"`
	Type        string `json:"type" db:"type" gorm:"not null;comment:类型:animation,interaction,transition"`
	Preview     string `json:"preview" db:"preview" gorm:"comment:预览图/GIF"`
	Configuration string `json:"configuration" db:"configuration" gorm:"type:json;not null;comment:配置,JSON格式"`
	Description string `json:"description" db:"description" gorm:"comment:描述"`
	IsSystem    bool   `json:"isSystem" db:"is_system" gorm:"default:false;comment:是否系统特效"`
	Status      int    `json:"status" db:"status" gorm:"default:1;comment:状态:1-正常,0-禁用"`
	TenantID    int64  `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
