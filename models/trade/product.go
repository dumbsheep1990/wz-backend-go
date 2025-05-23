package trade

import (
	"time"

	"wz-project/wz-backend-go/models/common"
)

// Product 产品模型
type Product struct {
	common.BaseIDModel
	common.BaseTimeModel
	ProductID    int64     `json:"productId" db:"product_id" gorm:"uniqueIndex;not null;comment:产品ID，业务唯一标识"`
	Name         string    `json:"name" db:"name" gorm:"not null;comment:产品名称"`
	CompanyID    int64     `json:"companyId" db:"company_id" gorm:"index;not null;comment:公司ID"`
	CompanyName  string    `json:"companyName" db:"company_name" gorm:"not null;comment:公司名称"`
	CategoryID   int64     `json:"categoryId" db:"category_id" gorm:"index;comment:分类ID"`
	CategoryName string    `json:"categoryName" db:"category_name" gorm:"comment:分类名称"`
	Price        float64   `json:"price" db:"price" gorm:"not null;comment:产品价格"`
	Currency     string    `json:"currency" db:"currency" gorm:"default:CNY;comment:货币类型"`
	Specifications interface{} `json:"specifications" db:"specifications" gorm:"type:json;comment:规格,JSON格式"`
	Material      interface{} `json:"material" db:"material" gorm:"type:json;comment:材质,JSON格式"`
	StockQuantity int       `json:"stockQuantity" db:"stock_quantity" gorm:"not null;default:0;comment:库存数量"`
	MinimumOrderQuantity int `json:"minimumOrderQuantity" db:"minimum_order_quantity" gorm:"default:1;comment:最小订购量"`
	Description  string    `json:"description" db:"description" gorm:"type:text;comment:产品描述"`
	Images       []string  `json:"images" db:"images" gorm:"type:json;comment:产品图片"`
	ContactPerson string   `json:"contactPerson" db:"contact_person" gorm:"comment:联系人"`
	ContactPhone string    `json:"contactPhone" db:"contact_phone" gorm:"comment:联系电话"`
	ContactEmail string    `json:"contactEmail" db:"contact_email" gorm:"comment:联系邮箱"`
	Address      string    `json:"address" db:"address" gorm:"comment:地址"`
	ViewCount    int       `json:"viewCount" db:"view_count" gorm:"default:0;comment:浏览量"`
	SaleCount    int       `json:"saleCount" db:"sale_count" gorm:"default:0;comment:销量"`
	Status       int       `json:"status" db:"status" gorm:"default:1;comment:状态:1-上架,0-下架"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ProductCategory 产品分类
type ProductCategory struct {
	common.BaseIDModel
	common.BaseTimeModel
	ParentID     int64     `json:"parentId" db:"parent_id" gorm:"index;default:0;comment:父分类ID"`
	Name         string    `json:"name" db:"name" gorm:"not null;comment:分类名称"`
	Code         string    `json:"code" db:"code" gorm:"uniqueIndex:idx_tenant_code;not null;comment:分类编码"`
	Level        int       `json:"level" db:"level" gorm:"default:1;comment:层级，1为顶级"`
	Path         string    `json:"path" db:"path" gorm:"comment:分类路径，例如：1,2,3"`
	Icon         string    `json:"icon" db:"icon" gorm:"comment:分类图标"`
	Banner       string    `json:"banner" db:"banner" gorm:"comment:分类banner"`
	SortOrder    int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	IsVisible    bool      `json:"isVisible" db:"is_visible" gorm:"default:true;comment:是否可见"`
	Description  string    `json:"description" db:"description" gorm:"comment:描述"`
	SeoTitle     string    `json:"seoTitle" db:"seo_title" gorm:"comment:SEO标题"`
	SeoKeywords  string    `json:"seoKeywords" db:"seo_keywords" gorm:"comment:SEO关键词"`
	SeoDescription string  `json:"seoDescription" db:"seo_description" gorm:"comment:SEO描述"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"uniqueIndex:idx_tenant_code;index;comment:租户ID"`
}

// ProductAttribute 产品属性
type ProductAttribute struct {
	common.BaseIDModel
	common.BaseTimeModel
	Name         string    `json:"name" db:"name" gorm:"not null;comment:属性名称"`
	Code         string    `json:"code" db:"code" gorm:"uniqueIndex:idx_tenant_code;not null;comment:属性编码"`
	AttributeType string   `json:"attributeType" db:"attribute_type" gorm:"not null;comment:属性类型：text-文本 number-数字 enum-枚举 date-日期"`
	InputType    string    `json:"inputType" db:"input_type" gorm:"not null;comment:输入类型：input-输入框 select-下拉 radio-单选 checkbox-多选 date-日期"`
	IsRequired   bool      `json:"isRequired" db:"is_required" gorm:"default:false;comment:是否必填"`
	IsMultiple   bool      `json:"isMultiple" db:"is_multiple" gorm:"default:false;comment:是否多选"`
	IsFilterable bool      `json:"isFilterable" db:"is_filterable" gorm:"default:false;comment:是否可筛选"`
	IsSearchable bool      `json:"isSearchable" db:"is_searchable" gorm:"default:false;comment:是否可搜索"`
	IsComparable bool      `json:"isComparable" db:"is_comparable" gorm:"default:false;comment:是否可比较"`
	SortOrder    int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	DefaultValue string    `json:"defaultValue" db:"default_value" gorm:"comment:默认值"`
	Description  string    `json:"description" db:"description" gorm:"comment:描述"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"uniqueIndex:idx_tenant_code;index;comment:租户ID"`
}

// ProductAttributeValue 产品属性值
type ProductAttributeValue struct {
	common.BaseIDModel
	AttributeID  int64     `json:"attributeId" db:"attribute_id" gorm:"index;not null;comment:属性ID"`
	Value        string    `json:"value" db:"value" gorm:"not null;comment:属性值"`
	SortOrder    int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ProductSpecification 产品规格
type ProductSpecification struct {
	common.BaseIDModel
	common.BaseTimeModel
	ProductID    int64     `json:"productId" db:"product_id" gorm:"index;not null;comment:产品ID"`
	SpecCode     string    `json:"specCode" db:"spec_code" gorm:"uniqueIndex:idx_product_spec_code;not null;comment:规格编码，如SKU"`
	SpecValues   string    `json:"specValues" db:"spec_values" gorm:"type:json;not null;comment:规格值，JSON格式，如：{\"颜色\":\"红色\",\"尺寸\":\"XL\"}"`
	Price        float64   `json:"price" db:"price" gorm:"not null;comment:价格"`
	OriginalPrice float64  `json:"originalPrice" db:"original_price" gorm:"comment:原价"`
	CostPrice    float64   `json:"costPrice" db:"cost_price" gorm:"comment:成本价"`
	Weight       float64   `json:"weight" db:"weight" gorm:"comment:重量(g)"`
	Volume       float64   `json:"volume" db:"volume" gorm:"comment:体积(cm³)"`
	Stock        int       `json:"stock" db:"stock" gorm:"default:0;comment:库存"`
	Sold         int       `json:"sold" db:"sold" gorm:"default:0;comment:已售"`
	IsDefault    bool      `json:"isDefault" db:"is_default" gorm:"default:false;comment:是否默认"`
	IsEnabled    bool      `json:"isEnabled" db:"is_enabled" gorm:"default:true;comment:是否启用"`
	Images       string    `json:"images" db:"images" gorm:"type:json;comment:图片，JSON格式"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ProductCategoryRelation 产品分类关联
type ProductCategoryRelation struct {
	common.BaseIDModel
	ProductID    int64     `json:"productId" db:"product_id" gorm:"uniqueIndex:idx_product_category;index;not null;comment:产品ID"`
	CategoryID   int64     `json:"categoryId" db:"category_id" gorm:"uniqueIndex:idx_product_category;index;not null;comment:分类ID"`
	IsPrimary    bool      `json:"isPrimary" db:"is_primary" gorm:"default:false;comment:是否主分类"`
	SortOrder    int       `json:"sortOrder" db:"sort_order" gorm:"default:0;comment:排序"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ProductAttributeRelation 产品属性关联
type ProductAttributeRelation struct {
	common.BaseIDModel
	ProductID     int64     `json:"productId" db:"product_id" gorm:"uniqueIndex:idx_product_attribute;index;not null;comment:产品ID"`
	AttributeID   int64     `json:"attributeId" db:"attribute_id" gorm:"uniqueIndex:idx_product_attribute;index;not null;comment:属性ID"`
	AttributeValue string   `json:"attributeValue" db:"attribute_value" gorm:"comment:属性值"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at" gorm:"autoCreateTime"`
	TenantID      int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}

// ProductReview 商品评价
type ProductReview struct {
	common.BaseIDModel
	common.BaseTimeModel
	UserID       int64     `json:"userId" db:"user_id" gorm:"index;not null;comment:用户ID"`
	ProductID    int64     `json:"productId" db:"product_id" gorm:"index;not null;comment:产品ID"`
	OrderID      string    `json:"orderId" db:"order_id" gorm:"index;comment:订单ID"`
	SpecID       int64     `json:"specId" db:"spec_id" gorm:"index;comment:规格ID"`
	Rating       int       `json:"rating" db:"rating" gorm:"not null;comment:评分(1-5星)"`
	Content      string    `json:"content" db:"content" gorm:"type:text;comment:评价内容"`
	Images       string    `json:"images" db:"images" gorm:"type:json;comment:图片，JSON格式"`
	IsAnonymous  bool      `json:"isAnonymous" db:"is_anonymous" gorm:"default:false;comment:是否匿名"`
	Status       string    `json:"status" db:"status" gorm:"default:pending;comment:状态：pending-待审核 approved-已通过 rejected-已拒绝"`
	ReplyContent string    `json:"replyContent" db:"reply_content" gorm:"type:text;comment:回复内容"`
	ReplyTime    *time.Time `json:"replyTime" db:"reply_time" gorm:"comment:回复时间"`
	Likes        int       `json:"likes" db:"likes" gorm:"default:0;comment:点赞数"`
	TenantID     int64     `json:"tenantId" db:"tenant_id" gorm:"index;comment:租户ID"`
}
