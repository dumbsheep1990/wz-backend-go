package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ProductName 商品名称值对象
type ProductName string

// NewProductName 创建商品名称值对象
func NewProductName(name string) (ProductName, error) {
	if name == "" {
		return "", errors.New("商品名称不能为空")
	}

	if len(name) > 100 {
		return "", errors.New("商品名称不能超过100个字符")
	}

	return ProductName(name), nil
}

// Value 获取名称值
func (n ProductName) Value() string {
	return string(n)
}

// String 字符串表示
func (n ProductName) String() string {
	return string(n)
}

// ProductDescription 商品描述值对象
type ProductDescription string

// NewProductDescription 创建商品描述值对象
func NewProductDescription(description string) (ProductDescription, error) {
	if len(description) > 2000 {
		return "", errors.New("商品描述不能超过2000个字符")
	}

	return ProductDescription(description), nil
}

// Value 获取描述值
func (d ProductDescription) Value() string {
	return string(d)
}

// String 字符串表示
func (d ProductDescription) String() string {
	return string(d)
}

// ProductSKU 商品SKU值对象
type ProductSKU string

// SKU格式正则表达式: 字母、数字和连字符组合，长度为6-20个字符
var skuRegex = regexp.MustCompile(`^[a-zA-Z0-9-]{6,20}$`)

// NewProductSKU 创建商品SKU值对象
func NewProductSKU(sku string) (ProductSKU, error) {
	if sku == "" {
		return "", errors.New("SKU不能为空")
	}

	// 转换为大写并去除空格
	sku = strings.ToUpper(strings.TrimSpace(sku))

	if !skuRegex.MatchString(sku) {
		return "", errors.New("SKU格式不正确，必须是6-20个字母、数字或连字符的组合")
	}

	return ProductSKU(sku), nil
}

// Value 获取SKU值
func (s ProductSKU) Value() string {
	return string(s)
}

// String 字符串表示
func (s ProductSKU) String() string {
	return string(s)
}

// Price 价格值对象
type Price int64

// NewPrice 创建价格值对象（以分为单位）
func NewPrice(price int64) (Price, error) {
	if price < 0 {
		return 0, errors.New("价格不能为负数")
	}

	return Price(price), nil
}

// Value 获取价格值（以分为单位）
func (p Price) Value() int64 {
	return int64(p)
}

// Float 获取价格的浮点数表示（以元为单位）
func (p Price) Float() float64 {
	return float64(p) / 100.0
}

// String 价格的字符串表示（以元为单位，保留两位小数）
func (p Price) String() string {
	return fmt.Sprintf("%.2f", p.Float())
}

// ProductStatus 商品状态值对象
type ProductStatus int32

const (
	ProductStatusDraft    ProductStatus = 0 // 草稿
	ProductStatusActive   ProductStatus = 1 // 上架
	ProductStatusInactive ProductStatus = 2 // 下架
	ProductStatusSoldOut  ProductStatus = 3 // 售罄
	ProductStatusDeleted  ProductStatus = 4 // 已删除
)

// NewProductStatus 创建商品状态值对象
func NewProductStatus(status int32) (ProductStatus, error) {
	switch status {
	case int32(ProductStatusDraft), int32(ProductStatusActive), int32(ProductStatusInactive), int32(ProductStatusSoldOut), int32(ProductStatusDeleted):
		return ProductStatus(status), nil
	default:
		return 0, errors.New("无效的商品状态")
	}
}

// Value 获取状态值
func (s ProductStatus) Value() int32 {
	return int32(s)
}

// String 状态的字符串表示
func (s ProductStatus) String() string {
	switch s {
	case ProductStatusDraft:
		return "草稿"
	case ProductStatusActive:
		return "上架"
	case ProductStatusInactive:
		return "下架"
	case ProductStatusSoldOut:
		return "售罄"
	case ProductStatusDeleted:
		return "已删除"
	default:
		return "未知状态"
	}
}
